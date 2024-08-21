package battle

import (
	"fmt"
	"sao/types"
	"sao/utils"
	"sao/world/calendar"
	"sao/world/location"
	"slices"

	"github.com/disgoorg/disgo/discord"
	"github.com/google/uuid"
)

type SummonEntityMeta struct {
	Owner uuid.UUID
	Type  uuid.UUID
}

type Fight struct {
	Entities        EntityMap
	ExpireMap       map[uuid.UUID]int
	SummonMap       map[uuid.UUID]SummonEntityMeta
	SpeedMap        map[uuid.UUID]int
	StartTime       *calendar.Calendar
	ExternalChannel chan FightEvent
	DiscordChannel  chan types.DiscordMessageStruct
	Effects         []ActionEffect
	Location        *location.Location
	Tournament      *TournamentData
	AdditionalLoot  []types.WithTarget[Loot]
	PlayerActions   chan Action
}

type TournamentData struct {
	Tournament uuid.UUID
	Location   string
}

type EntityMap map[uuid.UUID]EntityEntry

type EntityEntry struct {
	Entity Entity
	Side   int
}

func (mp EntityMap) SidesLeft() []int {
	sides := make([]int, 0)

	for _, entity := range mp {

		if entity.Entity.GetCurrentHP() <= 0 {
			continue
		}

		exists := false

		for _, side := range sides {
			if side == entity.Side {
				exists = true
			}
		}

		if !exists {
			sides = append(sides, entity.Side)
		}
	}

	return sides
}

func (mp EntityMap) FromSide(side int) []Entity {
	entities := make([]Entity, 0)

	for _, entity := range mp {
		if entity.Side == side {
			entities = append(entities, entity.Entity)
		}
	}

	return entities
}

func (f *Fight) AddAdditionalLoot(loot Loot, source uuid.UUID, teamWide bool) {
	if teamWide {
		for _, entity := range f.GetAlliesFor(source) {
			f.AdditionalLoot = append(
				f.AdditionalLoot,
				types.WithTarget[Loot]{
					Value:  loot,
					Target: entity.GetUUID(),
				},
			)
		}
	}

	f.AdditionalLoot = append(
		f.AdditionalLoot,
		types.WithTarget[Loot]{
			Value:  loot,
			Target: source,
		},
	)
}

func (f *Fight) IsFinished() bool {
	return len(f.Entities.SidesLeft()) <= 1
}

func (f *Fight) GetEnemiesFor(uuid uuid.UUID) []Entity {
	entitySide := f.Entities[uuid].Side

	enemiesList := make([]Entity, 0)

	for _, entry := range f.Entities {
		if entry.Side == entitySide {
			continue
		}

		if entry.Entity.GetCurrentHP() <= 0 {
			continue
		}

		enemiesList = append(enemiesList, entry.Entity)
	}

	return enemiesList
}

func (f *Fight) GetAlliesFor(uuid uuid.UUID) []Entity {
	entitySide := f.Entities[uuid].Side

	alliesList := make([]Entity, 0)

	for _, entry := range f.Entities {
		if entry.Side != entitySide {
			continue
		}

		if entry.Entity.GetCurrentHP() <= 0 {
			continue
		}

		alliesList = append(alliesList, entry.Entity)
	}

	return alliesList
}

func (f *Fight) HandleAction(act Action) {
	channelId := f.Location.CID

	if f.Tournament != nil {
		channelId = f.Tournament.Location
	}

	switch act.Event {
	case ACTION_ATTACK:
		sourceEntityEntry := f.Entities[act.Source]
		sourceEntity := sourceEntityEntry.Entity

		tempMeta := act.Meta

		if _, ok := tempMeta.(ActionDamage); !ok {
			tempMeta = ActionDamage{
				Damage: []Damage{
					{
						Value:    sourceEntity.GetStat(types.STAT_AD),
						Type:     types.DMG_PHYSICAL,
						CanDodge: true,
					},
				},
				CanDodge: true,
			}
		}

		overallDmg := []Damage{
			{Value: 0, Type: types.DMG_PHYSICAL},
			{Value: 0, Type: types.DMG_MAGICAL},
			{Value: 0, Type: types.DMG_TRUE},
		}

		meta := tempMeta.(ActionDamage)

		increaseValue := sourceEntity.TriggerEvent(types.TRIGGER_ATTACK_BEFORE, types.EventData{
			Source: sourceEntity,
			Target: f.Entities[act.Target].Entity,
			Fight:  &f,
		}, meta)

		for _, val := range increaseValue {
			for _, dmg := range val.(types.AttackTriggerMeta).Effects {
				meta.Damage = append(meta.Damage, Damage{
					Value:     dmg.Value,
					Type:      dmg.Type,
					IsPercent: dmg.Percent,
				})
			}
		}

		slices.SortFunc(meta.Damage, func(left Damage, right Damage) int {
			if left.IsPercent && !right.IsPercent {
				return -1
			}

			if !left.IsPercent && right.IsPercent {
				return 1
			}

			return 0
		})

		canDodge := meta.CanDodge && f.Entities[act.Target].Entity.CanDodge()

		dodged := false
		var dmgDealt []Damage

		if canDodge {
			for _, dmg := range meta.Damage {
				if dmg.IsPercent {
					continue
				}
				overallDmg[dmg.Type].Value += dmg.Value
			}

			dmgDealt, dodged = f.Entities[act.Target].Entity.(DodgeEntity).TakeDMGOrDodge(
				ActionDamage{overallDmg, canDodge},
			)
		} else {
			dmgDealt = f.Entities[act.Target].Entity.TakeDMG(meta)
		}

		tempEmbed := discord.NewEmbedBuilder().SetTitle("Atak")

		if !dodged {
			sourceEntity.TriggerEvent(types.TRIGGER_ATTACK_HIT, types.EventData{
				Source: sourceEntity,
				Target: f.Entities[act.Target].Entity,
				Fight:  &f,
			}, meta)
			f.Entities[act.Target].Entity.TriggerEvent(types.TRIGGER_ATTACK_GOT_HIT, types.EventData{
				Target: sourceEntity,
				Source: f.Entities[act.Target].Entity,
				Fight:  &f,
			}, nil)

			dmgSum := dmgDealt[0].Value + dmgDealt[1].Value + dmgDealt[2].Value

			tempEmbed.
				SetFooterTextf("%s zaatakował %s", sourceEntity.GetName(), f.Entities[act.Target].Entity.GetName()).
				SetDescriptionf("Zadano łacznie %d obrażeń", dmgSum)

			dmgText := ""

			for _, dmg := range meta.Damage {
				if dmg.Value == 0 {
					continue
				}

				dmgType := "fizycznych"

				switch dmg.Type {
				case types.DMG_MAGICAL:
					dmgType = "magicznych"
				case types.DMG_TRUE:
					dmgType = "prawdziwych"
				}

				if dmg.IsPercent {
					dmgText += fmt.Sprintf("- %d%% obrażeń %s\n", dmg.Value, dmgType)
				} else {
					dmgText += fmt.Sprintf("- %d obrażeń %s\n", dmg.Value, dmgType)
				}
			}

			tempEmbed.AddField("Obrażenia", dmgText, false)

			vampValue := sourceEntity.GetStat(types.STAT_ATK_VAMP)

			if vampValue > 0 {
				value := utils.PercentOf(dmgSum, vampValue)

				sourceEntity.Heal(value)

				sourceEntity.TriggerEvent(types.TRIGGER_HEAL_SELF, types.EventData{
					Source: sourceEntity,
					Target: sourceEntity,
					Fight:  &f,
				}, ActionEffectHeal{Value: value})

				tempEmbed.AddField("Wampiryzm!", fmt.Sprintf("%s dodatkowo wyleczył się o %d", sourceEntity.GetName(), value), false)
			}

		} else {
			sourceEntity.TriggerEvent(types.TRIGGER_ATTACK_MISS, types.EventData{
				Source: sourceEntity,
				Target: f.Entities[act.Target].Entity,
				Fight:  &f,
			}, nil)

			tempEmbed.SetDescriptionf("%s zaatakował %s, ale atak został uniknięty", sourceEntity.GetName(), f.Entities[act.Target].Entity.GetName())
		}

		targetEntity := f.Entities[act.Target]

		if targetEntity.Entity.GetFlags()&types.ENTITY_AUTO == 0 {
			if targetEntity.Entity.(PlayerEntity).GetDefendingState() {
				if utils.RandomNumber(0, 100) < targetEntity.Entity.GetStat(types.STAT_AGL) {
					counterDmg := utils.PercentOf(targetEntity.Entity.GetStat(types.STAT_AD), 70)

					counterDmg += utils.PercentOf(targetEntity.Entity.GetStat(types.STAT_DEF), 15)
					counterDmg += utils.PercentOf(targetEntity.Entity.GetStat(types.STAT_MR), 15)

					f.HandleAction(Action{
						Event:  ACTION_COUNTER,
						Source: act.Target,
						Target: act.Source,
						Meta: ActionDamage{
							Damage: []Damage{
								{
									Value:    counterDmg,
									Type:     types.DMG_PHYSICAL,
									CanDodge: true,
								},
							},
							CanDodge: true,
						},
					})
				}
			}
		}

		f.DiscordChannel <- types.DiscordMessageStruct{
			ChannelID:      channelId,
			MessageContent: discord.NewMessageCreateBuilder().AddEmbeds(tempEmbed.Build()).Build(),
		}
	case ACTION_EFFECT:
		meta := act.Meta.(ActionEffect)

		if meta.Duration == 0 {
			if meta.Effect == EFFECT_HEAL {
				healValue := meta.Value

				if act.Source != act.Target {
					healValue = utils.PercentOf(meta.Value, 100+f.Entities[act.Source].Entity.GetStat(types.STAT_HEAL_POWER))
				}

				f.Entities[act.Target].Entity.Heal(healValue)
				return
			}

			f.Entities[act.Target].Entity.ApplyEffect(meta)

			return
		}

		if meta.Effect == EFFECT_SHIELD {
			meta.Value = utils.PercentOf(meta.Value, 100+f.Entities[act.Source].Entity.GetStat(types.STAT_HEAL_POWER))

			return
		}

		if meta.Effect == EFFECT_TAUNT {
			for _, entity := range f.GetEnemiesFor(act.Target) {
				newEffect := ActionEffect{
					Effect:   EFFECT_TAUNTED,
					Duration: meta.Duration,
					Meta:     act.Target,
					Value:    0,
					Uuid:     uuid.New(),
				}

				entity.ApplyEffect(newEffect)
			}

			return
		}

		f.Entities[act.Target].Entity.ApplyEffect(meta)
	case ACTION_DEFEND:
		entity := f.Entities[act.Source]

		if entity.Entity.GetFlags()&types.ENTITY_AUTO == 0 {
			entity.Entity.(PlayerEntity).SetDefendingState(true)

			f.HandleAction(Action{
				Event:  ACTION_EFFECT,
				Source: act.Source,
				Target: act.Source,
				Meta: ActionEffect{
					Effect:   EFFECT_STAT_INC,
					Duration: 1,
					Meta: ActionEffectStat{
						Stat:      types.STAT_DEF,
						Value:     20,
						IsPercent: true,
					},
				},
			})

			f.HandleAction(Action{
				Event:  ACTION_EFFECT,
				Source: act.Source,
				Target: act.Source,
				Meta: ActionEffect{
					Effect:   EFFECT_STAT_INC,
					Duration: 0,
					Meta: ActionEffectStat{
						Stat:      types.STAT_MR,
						Value:     20,
						IsPercent: true,
					},
				},
			})
		}
	case ACTION_SKILL:
		sourceEntity := f.Entities[act.Source]

		sourceEntityFlags := sourceEntity.Entity.GetFlags()

		isPlayer := sourceEntityFlags&types.ENTITY_AUTO == 0
		isSummon := sourceEntityFlags&types.ENTITY_SUMMON != 0 || sourceEntityFlags&types.ENTITY_AUTO == 0

		if !isPlayer || !isSummon {
			return
		}

		skillUsageMeta := act.Meta.(ActionSkillMeta)

		skillUpgrades := 0

		if isPlayer {
			skillUpgrades = sourceEntity.Entity.(PlayerEntity).GetUpgrades(skillUsageMeta.Lvl)
		} else {
			ownerUuid := f.SummonMap[act.Source].Owner

			skillUpgrades = f.Entities[ownerUuid].Entity.(PlayerEntity).GetUpgrades(skillUsageMeta.Lvl)
		}

		skill := sourceEntity.Entity.(PlayerEntity).GetLvlSkill(skillUsageMeta.Lvl)

		trigger := skill.GetTrigger()

		skillCost := skill.GetCost()

		if skill.IsLevelSkill() {
			trigger = skill.(types.PlayerSkillUpgradable).GetUpgradableTrigger(skillUpgrades)

			skillCost = skill.(types.PlayerSkillUpgradable).GetUpgradableCost(skillUpgrades)
		}

		if sourceEntity.Entity.GetCurrentMana() < skillCost {
			f.DiscordChannel <- types.DiscordMessageStruct{
				ChannelID:      channelId,
				MessageContent: discord.NewMessageCreateBuilder().SetContent("Nie masz many na użycie tej umiejętności").Build(),
			}

			return
		}

		if trigger.Type != types.TRIGGER_ACTIVE {
			f.DiscordChannel <- types.DiscordMessageStruct{
				ChannelID:      channelId,
				MessageContent: discord.NewMessageCreateBuilder().SetContent("Nie można użyć tej umiejętności").Build(),
			}

			return
		}

		sourceEntity.Entity.RestoreMana(-skillCost)

		if skillUsageMeta.Lvl%10 == 0 {
			sourceEntity.Entity.TriggerEvent(types.TRIGGER_CAST_ULT, types.EventData{
				Source: sourceEntity.Entity,
				Target: f.Entities[act.Target].Entity,
				Fight:  &f,
			}, nil)
		}

		if len(skillUsageMeta.Targets) > 0 {
			for _, target := range skillUsageMeta.Targets {
				if skill.IsLevelSkill() {
					skill.(types.PlayerSkillUpgradable).UpgradableExecute(sourceEntity.Entity, f.Entities[target].Entity, f, nil, skillUpgrades)
				} else {
					skill.Execute(sourceEntity.Entity, f.Entities[target].Entity, f, nil)
				}
			}
		} else {
			if skill.IsLevelSkill() {
				skill.(types.PlayerSkillUpgradable).UpgradableExecute(sourceEntity.Entity, f.Entities[act.Target].Entity, f, nil, skillUpgrades)
			} else {
				skill.Execute(sourceEntity.Entity, f.Entities[act.Target].Entity, f, nil)
			}
		}

		f.DiscordChannel <- types.DiscordMessageStruct{
			ChannelID: channelId,
			MessageContent: discord.
				NewMessageCreateBuilder().
				AddEmbeds(
					discord.NewEmbedBuilder().
						SetTitle("Skill!").
						SetDescriptionf(
							"%s użył `%s`!\n",
							sourceEntity.Entity.GetName(),
							skill.GetName(),
						).
						Build(),
				).
				Build(),
		}

		return
	case ACTION_DMG:
		targetEntity := f.Entities[act.Target]
		meta := act.Meta.(ActionDamage)

		if targetEntity.Entity.GetCurrentHP() <= 0 {
			return
		}

		overallDmg := []Damage{
			{Value: 0, Type: types.DMG_PHYSICAL},
			{Value: 0, Type: types.DMG_MAGICAL},
			{Value: 0, Type: types.DMG_TRUE},
		}

		increaseValue := f.Entities[act.Target].Entity.TriggerEvent(types.TRIGGER_DAMAGE_BEFORE,
			types.EventData{
				Source: f.Entities[act.Source].Entity,
				Target: f.Entities[act.Target].Entity,
				Fight:  &f,
			}, nil)

		for _, val := range increaseValue {
			for _, dmg := range val.(types.DamageTriggerMeta).Effects {
				meta.Damage = append(meta.Damage, Damage{
					Value:     dmg.Value,
					Type:      dmg.Type,
					IsPercent: dmg.Percent,
				})
			}
		}

		slices.SortFunc(meta.Damage, func(left Damage, right Damage) int {
			if left.IsPercent && !right.IsPercent {
				return -1
			}

			if !left.IsPercent && right.IsPercent {
				return 1
			}

			return 0
		})

		canDodge := meta.CanDodge && f.Entities[act.Target].Entity.CanDodge()

		dodged := false
		var dmgDealt []Damage

		if canDodge {
			for _, dmg := range meta.Damage {
				if dmg.IsPercent {
					continue
				}
				overallDmg[dmg.Type].Value += dmg.Value
			}

			dmgDealt, dodged = f.Entities[act.Target].Entity.(DodgeEntity).TakeDMGOrDodge(
				ActionDamage{overallDmg, canDodge},
			)
		} else {
			dmgDealt = f.Entities[act.Target].Entity.TakeDMG(meta)
		}

		tempEmbed := discord.NewEmbedBuilder().SetTitle("Obrażenia!")

		sourceEntity := f.Entities[act.Source].Entity

		if !dodged {
			f.Entities[act.Target].Entity.TriggerEvent(types.TRIGGER_ATTACK_GOT_HIT, types.EventData{
				Target: sourceEntity,
				Source: f.Entities[act.Target].Entity,
				Fight:  &f,
			}, nil)

			dmgSum := dmgDealt[0].Value + dmgDealt[1].Value + dmgDealt[2].Value

			tempEmbed.
				SetFooterTextf("%s zadał obrażenia %s", sourceEntity.GetName(), f.Entities[act.Target].Entity.GetName()).
				SetDescriptionf("Zadano łacznie %d obrażeń", dmgSum)

			dmgText := ""

			for _, dmg := range meta.Damage {
				if dmg.Value == 0 {
					continue
				}

				dmgType := "fizycznych"

				switch dmg.Type {
				case types.DMG_MAGICAL:
					dmgType = "magicznych"
				case types.DMG_TRUE:
					dmgType = "prawdziwych"
				}

				if dmg.IsPercent {
					dmgText += fmt.Sprintf("- %d%% obrażeń %s\n", dmg.Value, dmgType)
				} else {
					dmgText += fmt.Sprintf("- %d obrażeń %s\n", dmg.Value, dmgType)
				}
			}

			tempEmbed.AddField("Obrażenia", dmgText, false)

			vampValue := sourceEntity.GetStat(types.STAT_OMNI_VAMP)

			if vampValue > 0 {
				value := utils.PercentOf(dmgSum, vampValue)

				sourceEntity.Heal(value)

				sourceEntity.TriggerEvent(types.TRIGGER_HEAL_SELF, types.EventData{
					Source: sourceEntity,
					Target: sourceEntity,
					Fight:  &f,
				}, ActionEffectHeal{Value: value})

				tempEmbed.AddField("Wampiryzm!", fmt.Sprintf("%s dodatkowo wyleczył się o %d", sourceEntity.GetName(), value), false)
			}

		} else {
			tempEmbed.SetDescriptionf("%s chciał zadać obrażenia %s, ale nie trafił", sourceEntity.GetName(), f.Entities[act.Target].Entity.GetName())
		}

		if targetEntity.Entity.GetCurrentHP() <= 0 {
			f.Entities[act.Source].Entity.TriggerEvent(types.TRIGGER_EXECUTE,
				types.EventData{
					Source: f.Entities[act.Source].Entity,
					Target: f.Entities[act.Target].Entity,
					Fight:  &f,
				}, nil)
		}

		if targetEntity.Entity.GetFlags()&types.ENTITY_AUTO == 0 && targetEntity.Entity.GetCurrentHP() > 0 {
			if targetEntity.Entity.(PlayerEntity).GetDefendingState() {
				if utils.RandomNumber(0, 100) < targetEntity.Entity.GetStat(types.STAT_AGL) {
					counterDmg := utils.PercentOf(targetEntity.Entity.GetStat(types.STAT_AD), 70)

					counterDmg += utils.PercentOf(targetEntity.Entity.GetStat(types.STAT_DEF), 15)
					counterDmg += utils.PercentOf(targetEntity.Entity.GetStat(types.STAT_MR), 15)

					f.HandleAction(Action{
						Event:  ACTION_COUNTER,
						Source: act.Target,
						Target: act.Source,
						Meta: ActionDamage{
							Damage: []Damage{
								{
									Value:    counterDmg,
									Type:     types.DMG_PHYSICAL,
									CanDodge: true,
								},
							},
							CanDodge: true,
						},
					})
				}
			}
		}

		f.DiscordChannel <- types.DiscordMessageStruct{
			ChannelID:      channelId,
			MessageContent: discord.NewMessageCreateBuilder().AddEmbeds(tempEmbed.Build()).Build(),
		}
	case ACTION_COUNTER:
		sourceEntityEntry := f.Entities[act.Source]
		sourceEntity := sourceEntityEntry.Entity

		overallDmg := []Damage{
			{Value: 0, Type: types.DMG_PHYSICAL},
			{Value: 0, Type: types.DMG_MAGICAL},
			{Value: 0, Type: types.DMG_TRUE},
		}

		meta := act.Meta.(ActionDamage)

		slices.SortFunc(meta.Damage, func(left Damage, right Damage) int {
			if left.IsPercent && !right.IsPercent {
				return -1
			}

			if !left.IsPercent && right.IsPercent {
				return 1
			}

			return 0
		})

		canDodge := meta.CanDodge && f.Entities[act.Target].Entity.CanDodge()

		dodged := false
		var dmgDealt []Damage

		if canDodge {
			for _, dmg := range meta.Damage {
				if dmg.IsPercent {
					continue
				}
				overallDmg[dmg.Type].Value += dmg.Value
			}

			dmgDealt, dodged = f.Entities[act.Target].Entity.(DodgeEntity).TakeDMGOrDodge(
				ActionDamage{overallDmg, canDodge},
			)
		} else {
			dmgDealt = f.Entities[act.Target].Entity.TakeDMG(meta)
		}

		tempEmbed := discord.NewEmbedBuilder().SetTitle("Kontra!")

		if !dodged {
			dmgSum := dmgDealt[0].Value + dmgDealt[1].Value + dmgDealt[2].Value

			tempEmbed.
				SetFooterTextf("%s zaatakował %s", sourceEntity.GetName(), f.Entities[act.Target].Entity.GetName()).
				SetDescriptionf("Zadano łącznie %d obrażeń", dmgSum)

			dmgText := ""

			for _, dmg := range meta.Damage {
				if dmg.Value == 0 {
					continue
				}

				dmgType := "fizycznych"

				switch dmg.Type {
				case types.DMG_MAGICAL:
					dmgType = "magicznych"
				case types.DMG_TRUE:
					dmgType = "prawdziwych"
				}

				if dmg.IsPercent {
					dmgText += fmt.Sprintf("- %d%% obrażeń %s\n", dmg.Value, dmgType)
				} else {
					dmgText += fmt.Sprintf("- %d obrażeń %s\n", dmg.Value, dmgType)
				}
			}

			tempEmbed.AddField("Obrażenia", dmgText, false)

			vampValue := sourceEntity.GetStat(types.STAT_ATK_VAMP)

			if vampValue > 0 {
				value := utils.PercentOf(dmgSum, vampValue)

				sourceEntity.Heal(value)

				sourceEntity.TriggerEvent(types.TRIGGER_HEAL_SELF,
					types.EventData{
						Source: sourceEntity,
						Target: f.Entities[act.Target].Entity,
						Fight:  &f,
					}, ActionEffectHeal{Value: value})

				tempEmbed.AddField("Wampiryzm!", fmt.Sprintf("%s dodatkowo wyleczył się o %d", sourceEntity.GetName(), value), false)
			}

		} else {
			sourceEntity.TriggerEvent(types.TRIGGER_ATTACK_MISS,
				types.EventData{
					Source: sourceEntity,
					Target: f.Entities[act.Target].Entity,
					Fight:  &f,
				}, nil)

			tempEmbed.SetDescriptionf("%s chciał skontrować ale nie trafił!", sourceEntity.GetName())
		}

		f.DiscordChannel <- types.DiscordMessageStruct{
			ChannelID:      channelId,
			MessageContent: discord.NewMessageCreateBuilder().AddEmbeds(tempEmbed.Build()).Build(),
		}

		targetEntity := f.Entities[act.Target]

		if targetEntity.Entity.GetFlags()&types.ENTITY_AUTO == 0 {
			if targetEntity.Entity.(PlayerEntity).GetDefendingState() {
				if utils.RandomNumber(0, 100) < targetEntity.Entity.GetStat(types.STAT_AGL) {

					counterDmg := utils.PercentOf(targetEntity.Entity.GetStat(types.STAT_AD), 70)

					counterDmg += utils.PercentOf(targetEntity.Entity.GetStat(types.STAT_DEF), 15)
					counterDmg += utils.PercentOf(targetEntity.Entity.GetStat(types.STAT_MR), 15)

					f.HandleAction(Action{
						Event:  ACTION_COUNTER,
						Source: act.Target,
						Target: act.Source,
						Meta: ActionDamage{
							Damage: []Damage{
								{
									Value:    counterDmg,
									Type:     types.DMG_PHYSICAL,
									CanDodge: true,
								},
							},
							CanDodge: true,
						},
					})
				}
			}
		}
	case ACTION_ITEM:
		sourceEntity := f.Entities[act.Source]

		itemMeta := act.Meta.(ActionItemMeta)

		var item *types.PlayerItem
		var itemIdx int

		for idx, invItem := range sourceEntity.Entity.(PlayerEntity).GetAllItems() {
			if invItem.UUID == itemMeta.Item {
				item = invItem
				itemIdx = idx
				break
			}
		}

		if act.Target == uuid.Nil {
			for _, target := range itemMeta.Targets {
				item.UseItem(sourceEntity.Entity, f.Entities[target].Entity, nil)
			}
		} else {
			item.UseItem(sourceEntity.Entity, f.Entities[act.Target].Entity, nil)
		}

		if item.Count == 0 && item.Consume {
			sourceEntity.Entity.(PlayerEntity).RemoveItem(itemIdx)
		}

		f.DiscordChannel <- types.DiscordMessageStruct{
			ChannelID: channelId,
			MessageContent: discord.
				NewMessageCreateBuilder().
				AddEmbeds(
					discord.NewEmbedBuilder().
						SetTitle("Przedmiot!").
						SetDescriptionf(
							"%s użył %s!\nEfekt: %s",
							sourceEntity.Entity.GetName(),
							item.Name, item.Description,
						).
						Build(),
				).Build(),
		}
	case ACTION_RUN:
		entity := f.Entities[act.Source].Entity
		side := f.Entities[act.Source].Side

		if utils.RandomNumber(0, 100) < entity.GetStat(types.STAT_AGL) {
			f.DiscordChannel <- types.DiscordMessageStruct{
				ChannelID: channelId,
				MessageContent: discord.
					NewMessageCreateBuilder().
					AddEmbeds(
						discord.NewEmbedBuilder().
							SetTitle("Ucieczka!").
							SetDescriptionf("%s próbował uciec, ale mu się to nie udało", entity.GetName()).
							SetColor(0xff0000).
							Build(),
					).Build(),
			}

			return
		}

		delete(f.Entities, act.Source)

		entities := f.Entities.FromSide(side)

		count := 0

		for _, entity := range entities {
			if entity.GetCurrentHP() > 0 && entity.GetFlags()&types.ENTITY_AUTO == 0 {
				count++
			}
		}

		f.DiscordChannel <- types.DiscordMessageStruct{
			ChannelID: channelId,
			MessageContent: discord.
				NewMessageCreateBuilder().
				AddEmbeds(
					discord.NewEmbedBuilder().
						SetTitle("Ucieczka!").
						SetDescriptionf("%s próbował uciec i mu się to udało", entity.GetName()).
						SetColor(0x00ff00).
						Build(),
				).Build(),
		}

		if count == 0 {
			f.ExternalChannel <- FightEndMsg{RunAway: true}
		}
	case ACTION_SUMMON:
		sourceEntity := f.Entities[act.Source]

		expires := act.Meta.(ActionSummon).Flags&SUMMON_FLAG_EXPIRE != 0

		newEntity := act.Meta.(ActionSummon).Entity

		f.DiscordChannel <- types.DiscordMessageStruct{
			ChannelID: channelId,
			MessageContent: discord.
				NewMessageCreateBuilder().
				AddEmbeds(
					discord.
						NewEmbedBuilder().
						SetTitle("Przywołanie!").
						SetDescriptionf("%s przywołał %s", sourceEntity.Entity.GetName(), newEntity.GetName()).
						SetColor(0x00ff00).
						Build(),
				).
				Build(),
		}

		if expires {
			f.ExpireMap[newEntity.GetUUID()] = act.Meta.(ActionSummon).ExpireTimer
		}

		f.SummonMap[newEntity.GetUUID()] = SummonEntityMeta{
			Owner: sourceEntity.Entity.GetUUID(),
			Type:  act.Meta.(ActionSummon).EntityType,
		}

		f.Entities[act.Meta.(ActionSummon).Entity.GetUUID()] = EntityEntry{
			Entity: newEntity,
			Side:   sourceEntity.Side,
		}
	}
}

func (f *Fight) CanSummon(entityType uuid.UUID, maxCount int) bool {
	if maxCount <= 0 || entityType == uuid.Nil {
		return true
	}

	count := 0

	for _, entity := range f.SummonMap {
		if entity.Type == entityType {
			count++
		}
	}

	return count < maxCount
}

func (f *Fight) Init() {
	f.SpeedMap = make(map[uuid.UUID]int)

	for _, entity := range f.Entities {
		f.SpeedMap[entity.Entity.GetUUID()] = entity.Entity.GetStat(types.STAT_SPD)
	}

	f.ExternalChannel = make(chan FightEvent, 10)
	f.PlayerActions = make(chan Action, 10)
	f.ExpireMap = make(map[uuid.UUID]int)
	f.SummonMap = make(map[uuid.UUID]SummonEntityMeta)
}

func (f *Fight) Run() {
	f.ExternalChannel <- FightStartMsg{}

	channelId := f.Location.CID

	if f.Tournament != nil {
		channelId = f.Tournament.Location
	}

	for len(f.Entities.SidesLeft()) > 1 {

		for entity, exp := range f.ExpireMap {
			f.ExpireMap[entity] = exp - 1

			if exp <= 0 {
				delete(f.ExpireMap, entity)

				f.ExternalChannel <- SummonExpired{Entity: entity}
			}
		}

		turnList := make([]uuid.UUID, 0)

		for uuid, speed := range f.SpeedMap {
			entity := f.Entities[uuid].Entity

			f.SpeedMap[uuid] = speed + entity.GetStat(types.STAT_SPD)

			turns := f.SpeedMap[uuid] / SPEED_GAUGE

			if turns == 0 {
				continue
			}

			f.SpeedMap[uuid] -= turns * SPEED_GAUGE

			for i := 0; i < turns; i++ {
				turnList = append(turnList, entity.GetUUID())
			}
		}

		for _, entityUuid := range turnList {
			entity := f.Entities[entityUuid].Entity

			if entity.GetCurrentHP() <= 0 {
				continue
			}

			entity.TriggerEvent(types.TRIGGER_TURN, types.EventData{
				Source: entity,
				Target: entity,
				Fight:  &f,
			}, nil)

			if entity.GetFlags()&types.ENTITY_AUTO == 0 {
				entity.(PlayerEntity).SetDefendingState(false)

				entity.(PlayerEntity).ReduceCooldowns(types.TRIGGER_TURN)

				f.ExternalChannel <- FightActionNeededMsg{Entity: entityUuid}

				tempAction := <-f.PlayerActions

				f.HandleAction(tempAction)

				for tempAction.ConsumeTurn != nil && !*tempAction.ConsumeTurn {
					f.ExternalChannel <- FightActionNeededMsg{Entity: entityUuid}

					tempAction = <-f.PlayerActions

					f.HandleAction(tempAction)
				}
			} else {
				if entity.GetEffectByType(EFFECT_STUN) == nil {
					for _, action := range entity.Action(f) {
						f.HandleAction(action)
					}
				} else {
					f.DiscordChannel <- types.DiscordMessageStruct{
						ChannelID: channelId,
						MessageContent: discord.
							NewMessageCreateBuilder().
							AddEmbeds(
								discord.NewEmbedBuilder().
									SetTitle("Efekt!").
									SetDescriptionf("%s jest unieruchomiony, pomijamy!", entity.GetName()).
									Build(),
							).Build(),
					}
				}
			}

			entity.TriggerAllEffects()
		}

		for _, entry := range f.Entities {
			if entry.Entity.GetCurrentHP() <= 0 {
				continue
			}
		}
	}

	f.ExternalChannel <- FightEndMsg{}
}

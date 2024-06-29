package battle

import (
	"fmt"
	"sao/types"
	"sao/utils"
	"sao/world/calendar"
	"sao/world/location"
	"slices"
	"sort"

	"github.com/disgoorg/disgo/discord"
	"github.com/google/uuid"
)

type Fight struct {
	Entities        EntityMap
	SpeedMap        map[uuid.UUID]int
	StartTime       *calendar.Calendar
	ExternalChannel chan FightEvent
	DiscordChannel  chan types.DiscordMessageStruct
	Effects         []ActionEffect
	Location        *location.Location
	Tournament      *TournamentData
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

func (f *Fight) TriggerPassive(entityUuid uuid.UUID, triggerType types.SkillTrigger, meta interface{}) {
	entityEntry, exists := f.Entities[entityUuid]

	if !exists {
		return
	}

	if entityEntry.Entity.IsAuto() {
		return
	}

	sourceEntity := entityEntry.Entity.(PlayerEntity)

	for _, skill := range sourceEntity.GetAllSkills() {
		if skill.GetTrigger().Type == types.TRIGGER_ACTIVE {
			continue
		}

		if skill.GetTrigger().Event.TriggerType != triggerType {
			continue
		}

		//TODO CD for skills that are not lvl bound
		if skill.IsLevelSkill() {
			sourceEntity.SetLvlCD(skill.(types.PlayerSkillLevel).GetLevel(), skill.GetCD())
		}

		if skill.GetCost() != 0 {
			sourceEntity.RestoreMana(-skill.GetCost())
		}

		targets := f.FindValidTargets(sourceEntity.GetUUID(), *skill.GetTrigger().Event)

		if skill.GetTrigger().Event.TargetCount != -1 {
			count := skill.GetTrigger().Event.TargetCount

			if count > len(targets) {
				count = len(targets)
			}

			targets = targets[:count]
		}

		var tempFight interface{} = f

		for _, target := range targets {
			skill.Execute(sourceEntity, f.Entities[target].Entity, &tempFight, meta)
		}
	}
}

func (f *Fight) TriggerPassiveWithCheck(entityUuid uuid.UUID, triggerType types.SkillTrigger, meta interface{}, additionalCheck func(Entity, types.PlayerSkill) bool) {
	entityEntry, exists := f.Entities[entityUuid]

	if !exists {
		return
	}

	if entityEntry.Entity.IsAuto() {
		return
	}

	sourceEntity := entityEntry.Entity.(PlayerEntity)

	for _, skill := range sourceEntity.GetAllSkills() {
		if skill.GetTrigger().Type == types.TRIGGER_ACTIVE {
			continue
		}

		if skill.GetTrigger().Event.TriggerType != triggerType {
			continue
		}

		if !additionalCheck(sourceEntity, skill) {
			continue
		}

		//TODO handle CD for not lvl skills
		if skill.IsLevelSkill() {
			sourceEntity.SetLvlCD(skill.(types.PlayerSkillLevel).GetLevel(), skill.GetCD())
		}

		if skill.GetCost() != 0 {
			sourceEntity.RestoreMana(-skill.GetCost())
		}

		targets := f.FindValidTargets(sourceEntity.GetUUID(), *skill.GetTrigger().Event)

		if skill.GetTrigger().Event.TargetCount != -1 {
			count := skill.GetTrigger().Event.TargetCount

			if count > len(targets) {
				count = len(targets)
			}

			targets = targets[:count]
		}

		var tempFight interface{} = f

		for _, target := range targets {
			skill.Execute(sourceEntity, f.Entities[target].Entity, &tempFight, meta)
		}
	}
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
						Value:    sourceEntity.GetATK(),
						Type:     DMG_PHYSICAL,
						CanDodge: true,
					},
				},
				CanDodge: true,
			}
		}

		overallDmg := []Damage{
			{Value: 0, Type: DMG_PHYSICAL},
			{Value: 0, Type: DMG_MAGICAL},
			{Value: 0, Type: DMG_TRUE},
		}

		meta := tempMeta.(ActionDamage)

		if !sourceEntity.IsAuto() {
			for _, skill := range sourceEntity.(PlayerEntity).GetAllSkills() {
				if skill.GetTrigger().Type == types.TRIGGER_ACTIVE {
					continue
				}

				if skill.GetTrigger().Event.TriggerType != types.TRIGGER_ATTACK_BEFORE {
					continue
				}

				//TODO CD for skills that are not lvl bound
				if skill.IsLevelSkill() {
					sourceEntity.(PlayerEntity).SetLvlCD(skill.(types.PlayerSkillLevel).GetLevel(), skill.GetCD())
				}

				if skill.GetCost() != 0 {
					sourceEntity.RestoreMana(-skill.GetCost())
				}

				eventData := skill.Execute(sourceEntity, f.Entities[act.Target].Entity, &f, meta)

				if eventData != nil {
					castMeta := make([]Damage, len(eventData.(types.AttackTriggerMeta).Effects))

					for idx, effect := range eventData.(types.AttackTriggerMeta).Effects {
						castMeta[idx] = Damage{
							Value:     effect.Value,
							Type:      DamageType(effect.Type),
							IsPercent: effect.Percent,
						}
					}
					meta.Damage = append(meta.Damage, castMeta...)
				}
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
			if !sourceEntity.IsAuto() {
				for _, skill := range sourceEntity.(PlayerEntity).GetAllSkills() {
					if skill.GetTrigger().Type == types.TRIGGER_ACTIVE {
						continue
					}

					if skill.GetTrigger().Event.TriggerType != types.TRIGGER_ATTACK_HIT {
						continue
					}

					//TODO CD for skills that are not lvl bound
					if skill.IsLevelSkill() {
						sourceEntity.(PlayerEntity).SetLvlCD(skill.(types.PlayerSkillLevel).GetLevel(), skill.GetCD())
					}

					if skill.GetCost() != 0 {
						sourceEntity.RestoreMana(-skill.GetCost())
					}

					skill.Execute(sourceEntity, f.Entities[act.Target].Entity, &f, meta)
				}
			}

			f.TriggerPassive(act.Target, types.TRIGGER_ATTACK_GOT_HIT, nil)

			dmgSum := dmgDealt[0].Value + dmgDealt[1].Value + dmgDealt[2].Value

			tempEmbed.SetFooterTextf("%s zaatakował %s", sourceEntity.GetName(), f.Entities[act.Target].Entity.GetName())

			tempEmbed.SetDescriptionf("Zadano łacznie %d obrażeń", dmgSum)

			dmgText := ""

			for _, dmg := range meta.Damage {
				if dmg.Value == 0 {
					continue
				}

				dmgType := "fizycznych"

				switch dmg.Type {
				case DMG_MAGICAL:
					dmgType = "magicznych"
				case DMG_TRUE:
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

				f.TriggerPassive(act.Source, types.TRIGGER_HEAL_SELF, ActionEffectHeal{Value: value})

				tempEmbed.AddField("Wampyryzm!", fmt.Sprintf("%s dodatkowo wyleczył się o %d", sourceEntity.GetName(), value), false)
			}

		} else {
			f.TriggerPassive(act.Source, types.TRIGGER_ATTACK_MISS, nil)
			f.TriggerPassive(act.Target, types.TRIGGER_DODGE, nil)

			tempEmbed.SetDescriptionf("%s zaatakował %s, ale atak został uniknięty", sourceEntity.GetName(), f.Entities[act.Target].Entity.GetName())
		}

		targetEntity := f.Entities[act.Target]

		if !targetEntity.Entity.IsAuto() {
			if targetEntity.Entity.(PlayerEntity).GetDefendingState() {
				if utils.RandomNumber(0, 100) < targetEntity.Entity.GetAGL() {
					counterDmg := utils.PercentOf(targetEntity.Entity.GetATK(), 70)

					counterDmg += utils.PercentOf(targetEntity.Entity.GetDEF(), 15)
					counterDmg += utils.PercentOf(targetEntity.Entity.GetMR(), 15)

					f.HandleAction(Action{
						Event:  ACTION_COUNTER,
						Source: act.Target,
						Target: act.Source,
						Meta: ActionDamage{
							Damage: []Damage{
								{
									Value:    counterDmg,
									Type:     DMG_PHYSICAL,
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
			if meta.Effect == EFFECT_HEAL_SELF {
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
		} else {
			f.Entities[act.Target].Entity.ApplyEffect(act.Meta.(ActionEffect))
		}
	case ACTION_DEFEND:
		entity := f.Entities[act.Source]

		f.TriggerPassive(act.Source, types.TRIGGER_DEFEND_START, nil)

		if !entity.Entity.IsAuto() {
			entity.Entity.(PlayerEntity).SetDefendingState(true)

			f.TriggerPassive(act.Source, types.TRIGGER_DEFEND_START, nil)

			f.HandleAction(Action{
				Event:  ACTION_EFFECT,
				Source: act.Source,
				Target: act.Source,
				Meta: ActionEffect{
					Effect:   EFFECT_STAT_INC,
					Duration: 1,
					Meta: &map[string]interface{}{
						"stat":     types.STAT_DEF,
						"percent":  20,
						"duration": 1,
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
					Meta: &map[string]interface{}{
						"stat":     types.STAT_MR,
						"percent":  20,
						"duration": 1,
					},
				},
			})
		}
	case ACTION_SKILL:
		sourceEntity := f.Entities[act.Source]

		if !sourceEntity.Entity.IsAuto() {
			for _, skill := range sourceEntity.Entity.(PlayerEntity).GetAllSkills() {
				if skill.GetTrigger().Type == types.TRIGGER_ACTIVE {
					continue
				}

				if skill.GetTrigger().Event.TriggerType != types.TRIGGER_MANA {
					continue
				}

				if skill.GetTrigger().Event.Meta["value"].(int) > sourceEntity.Entity.GetCurrentMana() {
					continue
				}

				targets := f.FindValidTargets(sourceEntity.Entity.GetUUID(), *skill.GetTrigger().Event)

				if skill.GetTrigger().Event.TargetCount != -1 {
					count := skill.GetTrigger().Event.TargetCount

					if count > len(targets) {
						count = len(targets)
					}

					targets = targets[:count]
				}

				for _, target := range targets {
					targetEntity := f.Entities[target]

					beforeSkillHP := targetEntity.Entity.GetCurrentHP()

					var tempFight interface{} = f

					skill.Execute(sourceEntity.Entity, targetEntity, &tempFight, nil)

					//Check if it's dmg skill so it doesn't trigger on heal/barrier etc
					if !targetEntity.Entity.IsAuto() && beforeSkillHP > targetEntity.Entity.GetCurrentHP() {
						if targetEntity.Entity.(PlayerEntity).GetDefendingState() {
							if utils.RandomNumber(0, 100) < targetEntity.Entity.GetAGL() {
								f.HandleAction(Action{
									Event:  ACTION_ATTACK,
									Source: act.Target,
									Target: act.Source,
									Meta: ActionDamage{
										Damage: []Damage{{
											Value:    targetEntity.Entity.GetATK(),
											Type:     DMG_PHYSICAL,
											CanDodge: true,
										}},
										CanDodge: true,
									},
								})

							}
						}
					}
				}
			}

			skillUsageMeta := act.Meta.(ActionSkillMeta)

			if skillUsageMeta.IsForLevel {
				skill := sourceEntity.Entity.(PlayerEntity).GetLvlSkill(skillUsageMeta.Lvl)

				if skill.GetTrigger().Type != types.TRIGGER_ACTIVE {
					return
				}

				if skillUsageMeta.Lvl%10 != 0 {
					f.TriggerPassive(act.Source, types.TRIGGER_CAST_LVL, nil)
				} else {
					f.TriggerPassive(act.Source, types.TRIGGER_CAST_ULT, nil)
				}

				var tempFight interface{} = f

				if act.Target == uuid.Nil {
					for _, target := range skillUsageMeta.Targets {
						skill.Execute(sourceEntity.Entity, f.Entities[target].Entity, &tempFight, nil)
					}
				} else {
					skill.Execute(sourceEntity.Entity, f.Entities[act.Target].Entity, &tempFight, nil)
				}

				f.DiscordChannel <- types.DiscordMessageStruct{
					ChannelID: channelId,
					MessageContent: discord.
						NewMessageCreateBuilder().
						AddEmbeds(
							discord.NewEmbedBuilder().
								SetTitle("Skill!").
								SetDescriptionf(
									"%s użył %s!\n",
									sourceEntity.Entity.GetName(),
									skill.GetName(),
								).
								Build(),
						).
						Build(),
				}
			}
		}
	case ACTION_DMG:
		targetEntity := f.Entities[act.Target]
		meta := act.Meta.(ActionDamage)

		if targetEntity.Entity.GetCurrentHP() <= 0 {
			return
		}

		f.TriggerPassive(act.Source, types.TRIGGER_DAMAGE_BEFORE, nil)

		if meta.CanDodge && targetEntity.Entity.CanDodge() {
			targetEntity.Entity.(DodgeEntity).TakeDMGOrDodge(meta)

			f.TriggerPassive(act.Target, types.TRIGGER_DAMAGE_AFTER, nil)
		} else {
			targetEntity.Entity.TakeDMG(meta)

			f.TriggerPassive(act.Target, types.TRIGGER_DODGE, nil)
		}

		if targetEntity.Entity.GetCurrentHP() <= 0 {
			f.TriggerPassive(act.Source, types.TRIGGER_EXECUTE, nil)
		}

		f.TriggerPassiveWithCheck(act.Source, types.TRIGGER_HEALTH, nil, func(e Entity, ps types.PlayerSkill) bool {
			hpValue := 0

			if ps.GetTrigger().Event.Meta["value"] != nil {
				hpValue = ps.GetTrigger().Event.Meta["value"].(int)
			} else {
				hpValue = (ps.GetTrigger().Event.Meta["percent"].(int) * e.GetMaxHP() / 100)
			}

			return hpValue < e.GetCurrentHP()
		})
	case ACTION_COUNTER:
		sourceEntityEntry := f.Entities[act.Source]
		sourceEntity := sourceEntityEntry.Entity

		overallDmg := []Damage{
			{Value: 0, Type: DMG_PHYSICAL},
			{Value: 0, Type: DMG_MAGICAL},
			{Value: 0, Type: DMG_TRUE},
		}

		meta := act.Meta.(ActionDamage)

		if !sourceEntity.IsAuto() {
			for _, skill := range sourceEntity.(PlayerEntity).GetAllSkills() {
				if skill.GetTrigger().Type == types.TRIGGER_ACTIVE {
					continue
				}

				if skill.GetTrigger().Event.TriggerType != types.TRIGGER_COUNTER_ATTEMPT {
					continue
				}

				//TODO CD for skills that are not lvl bound
				if skill.IsLevelSkill() {
					sourceEntity.(PlayerEntity).SetLvlCD(skill.(types.PlayerSkillLevel).GetLevel(), skill.GetCD())
				}

				if skill.GetCost() != 0 {
					sourceEntity.RestoreMana(-skill.GetCost())
				}

				eventData := skill.Execute(sourceEntity, f.Entities[act.Target].Entity, &f, meta)

				if eventData != nil {
					castMeta := make([]Damage, len(eventData.(types.AttackTriggerMeta).Effects))

					for idx, effect := range eventData.(types.AttackTriggerMeta).Effects {
						castMeta[idx] = Damage{
							Value:     effect.Value,
							Type:      DamageType(effect.Type),
							IsPercent: effect.Percent,
						}
					}
					meta.Damage = append(meta.Damage, castMeta...)
				}
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

		tempEmbed := discord.NewEmbedBuilder().SetTitle("Kontra!")

		if !dodged {
			if !sourceEntity.IsAuto() {
				for _, skill := range sourceEntity.(PlayerEntity).GetAllSkills() {
					if skill.GetTrigger().Type == types.TRIGGER_ACTIVE {
						continue
					}

					if skill.GetTrigger().Event.TriggerType != types.TRIGGER_ATTACK_HIT {
						continue
					}

					//TODO CD for skills that are not lvl bound
					if skill.IsLevelSkill() {
						sourceEntity.(PlayerEntity).SetLvlCD(skill.(types.PlayerSkillLevel).GetLevel(), skill.GetCD())
					}

					if skill.GetCost() != 0 {
						sourceEntity.RestoreMana(-skill.GetCost())
					}

					skill.Execute(sourceEntity, f.Entities[act.Target].Entity, &f, meta)
				}
			}

			dmgSum := dmgDealt[0].Value + dmgDealt[1].Value + dmgDealt[2].Value

			tempEmbed.SetFooterTextf("%s zaatakował %s", sourceEntity.GetName(), f.Entities[act.Target].Entity.GetName())

			tempEmbed.SetDescriptionf("Zadano łącznie %d obrażeń", dmgSum)

			dmgText := ""

			for _, dmg := range meta.Damage {
				if dmg.Value == 0 {
					continue
				}

				dmgType := "fizycznych"

				switch dmg.Type {
				case DMG_MAGICAL:
					dmgType = "magicznych"
				case DMG_TRUE:
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

				f.TriggerPassive(act.Source, types.TRIGGER_HEAL_SELF, ActionEffectHeal{Value: value})

				tempEmbed.AddField("Wampyryzm!", fmt.Sprintf("%s dodatkowo wyleczył się o %d", sourceEntity.GetName(), value), false)
			}

		} else {
			f.TriggerPassive(act.Source, types.TRIGGER_ATTACK_MISS, nil)
			f.TriggerPassive(act.Target, types.TRIGGER_DODGE, nil)

			tempEmbed.SetDescriptionf("%s chciał skontrować ale nie trafił!", sourceEntity.GetName())
		}

		f.DiscordChannel <- types.DiscordMessageStruct{
			ChannelID:      channelId,
			MessageContent: discord.NewMessageCreateBuilder().AddEmbeds(tempEmbed.Build()).Build(),
		}

		targetEntity := f.Entities[act.Target]

		if !targetEntity.Entity.IsAuto() {
			if targetEntity.Entity.(PlayerEntity).GetDefendingState() {
				if utils.RandomNumber(0, 100) < targetEntity.Entity.GetAGL() {

					counterDmg := utils.PercentOf(targetEntity.Entity.GetATK(), 70)

					counterDmg += utils.PercentOf(targetEntity.Entity.GetDEF(), 15)
					counterDmg += utils.PercentOf(targetEntity.Entity.GetMR(), 15)

					f.HandleAction(Action{
						Event:  ACTION_COUNTER,
						Source: act.Target,
						Target: act.Source,
						Meta: ActionDamage{
							Damage: []Damage{
								{
									Value:    counterDmg,
									Type:     DMG_PHYSICAL,
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

		if utils.RandomNumber(0, 100) < entity.GetAGL() {
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

		//TODO copy data lmao
		delete(f.Entities, act.Source)

		entities := f.Entities.FromSide(side)

		count := 0

		for _, entity := range entities {
			if entity.GetCurrentHP() > 0 && !entity.IsAuto() {
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
			f.ExternalChannel <- FightEndMsg{}
		}
	}
}

func (f *Fight) TriggerAll(triggerType types.SkillTrigger, meta interface{}) {
	for entityUuid := range f.Entities {
		f.TriggerPassive(entityUuid, triggerType, meta)
	}
}

func (f *Fight) Init() {
	f.SpeedMap = make(map[uuid.UUID]int)

	for _, entity := range f.Entities {
		f.SpeedMap[entity.Entity.GetUUID()] = entity.Entity.GetSPD()
	}

	f.ExternalChannel = make(chan FightEvent, 10)
	f.PlayerActions = make(chan Action, 10)

	f.TriggerAll(types.TRIGGER_FIGHT_START, nil)
}

func (f *Fight) FindValidTargets(source uuid.UUID, trigger types.EventTriggerDetails) []uuid.UUID {
	sourceEntity := f.Entities[source].Entity
	sourceSide := f.Entities[source].Side

	if len(trigger.TargetType) == 1 && trigger.TargetType[0] == types.TARGET_SELF {
		return []uuid.UUID{source}
	}

	targetEntities := make([]Entity, 0)

	for _, targetType := range trigger.TargetType {
		if targetType == types.TARGET_SELF {
			targetEntities = append(targetEntities, sourceEntity)
		}
	}

	isAllyValid := false

	for _, targetType := range trigger.TargetType {
		if targetType == types.TARGET_ALLY {
			isAllyValid = true
		}
	}

	isEnemyValid := false

	for _, targetType := range trigger.TargetType {
		if targetType == types.TARGET_ENEMY {
			isEnemyValid = true
		}
	}

	for _, entity := range f.Entities {
		if entity.Side == sourceSide && isAllyValid {
			targetEntities = append(targetEntities, entity.Entity)
		}

		if entity.Side != sourceSide && isEnemyValid {
			targetEntities = append(targetEntities, entity.Entity)
		}
	}

	sortInit := EntitySort{
		Entities: targetEntities,
		Order:    trigger.TargetDetails,
		Meta:     &trigger.Meta,
	}

	sort.Sort(sortInit)

	targets := make([]uuid.UUID, len(targetEntities))

	for i, entity := range sortInit.Entities {
		targets[i] = entity.GetUUID()
	}

	return targets
}

func (f *Fight) Run() {
	f.ExternalChannel <- FightStartMsg{}

	channelId := f.Location.CID

	if f.Tournament != nil {
		channelId = f.Tournament.Location
	}

	for len(f.Entities.SidesLeft()) > 1 {
		turnList := make([]uuid.UUID, 0)

		for uuid, speed := range f.SpeedMap {
			entity := f.Entities[uuid].Entity

			f.SpeedMap[uuid] = speed + entity.GetSPD()

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

			f.TriggerPassive(entityUuid, types.TRIGGER_TURN, nil)

			if !entity.IsAuto() {
				player := entity.(PlayerEntity)

				player.SetDefendingState(false)

				f.TriggerPassive(entityUuid, types.TRIGGER_DEFEND_END, nil)

				skills := player.GetSkillsCD()

				for skill, cd := range skills {
					if lvl, ok := skill.(int); ok {
						player.SetLvlCD(lvl, cd-1)
					} else {
						fmt.Println("TODO: non lvl skills")
						//TODO for non lvl skills
						// skillUuid := skill.(uuid.UUID)
					}
				}

				f.ExternalChannel <- FightActionNeededMsg{Entity: entityUuid}
				f.HandleAction(<-f.PlayerActions)
			} else {
				if !(entity.GetEffectByType(EFFECT_DISARM) != nil || entity.GetEffectByType(EFFECT_STUN) != nil || entity.GetEffectByType(EFFECT_STUN) != nil || entity.GetEffectByType(EFFECT_ROOT) != nil || entity.GetEffectByType(EFFECT_GROUND) != nil || entity.GetEffectByType(EFFECT_BLIND) != nil) {
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

	f.TriggerAll(types.TRIGGER_FIGHT_END, nil)
	f.ExternalChannel <- FightEndMsg{}
}

package battle

import (
	"fmt"
	"sao/data"
	"sao/types"
	"sao/utils"

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
	ExternalChannel chan FightEvent
	DiscordChannel  chan types.DiscordEvent
	Location        *types.Location
	Meta            *FightMeta
	PlayerActions   chan types.Action
	EventHandlers   map[uuid.UUID]EventHandler
}

func (f *Fight) Init() {
	f.ExternalChannel = make(chan FightEvent, 10)
	f.PlayerActions = make(chan types.Action, 10)
	f.ExpireMap = make(map[uuid.UUID]int)
	f.SummonMap = make(map[uuid.UUID]SummonEntityMeta)
	f.EventHandlers = make(map[uuid.UUID]EventHandler)
}

func (f *Fight) SidesLeft() []int {
	sides := make([]int, 0)

	for _, entity := range f.Entities {
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

func (f *Fight) AppendEventHandler(owner uuid.UUID, sTrigger types.SkillTrigger, handler func(source, target types.Entity, fightInstance types.FightInstance, meta any) any) uuid.UUID {
	handlerUuid := uuid.New()

	f.EventHandlers[handlerUuid] = EventHandler{Target: owner, Handler: handler, Trigger: sTrigger}

	return handlerUuid
}

func (f *Fight) TriggerEvent(source types.Entity, target types.Entity, event types.SkillTrigger, meta any) []any {
	returnValue := make([]any, 0)

	for _, handler := range f.EventHandlers {
		if handler.Trigger == event && handler.Target == target.GetUUID() {
			rValue := handler.Handler(source, target, f, meta)

			if rValue != nil {
				returnValue = append(returnValue, rValue)
			}
		}
	}

	rValue := source.TriggerEvent(event, types.EventData{Source: source, Target: target, Fight: f}, meta)

	if rValue != nil {
		returnValue = append(returnValue, rValue...)
	}

	if event == types.TRIGGER_ATTACK_GOT_HIT {
		rv := f.TriggerEvent(target, source, types.TRIGGER_DAMAGE_GOT_HIT, meta)

		if len(rv) > 0 {
			returnValue = append(returnValue, rv...)
		}
	}

	return returnValue
}

func (f *Fight) GetChannelId() string {
	if f.Meta.ThreadId != "" {
		return f.Meta.ThreadId
	}

	if f.Meta.Tournament != nil {
		return f.Meta.Tournament.Location
	}

	return f.Location.CID
}

func (f *Fight) GetEnemiesFor(uuid uuid.UUID) []types.Entity {
	entitySide := f.Entities[uuid].Side

	return f.FilteredEntities(func(entity EntityEntry) bool {
		return entity.Side != entitySide && entity.Entity.GetCurrentHP() > 0 && entity.Entity.GetUUID() != uuid
	})
}

func (f *Fight) GetAlliesFor(uuid uuid.UUID) []types.Entity {
	entitySide := f.Entities[uuid].Side

	return f.FilteredEntities(func(entity EntityEntry) bool {
		return entity.Side == entitySide && entity.Entity.GetCurrentHP() > 0 && entity.Entity.GetUUID() != uuid
	})
}

func (f *Fight) FilteredEntities(filter func(entity EntityEntry) bool) []types.Entity {
	entities := make([]types.Entity, 0)

	for _, entity := range f.Entities {
		if filter(*entity) {
			entities = append(entities, entity.Entity)
		}
	}

	return entities
}

func (f *Fight) HandleAction(act types.Action) {
	if act.Target == uuid.Nil {
		act.Target = act.Source
	}

	switch act.Event {
	case types.ACTION_ATTACK:
		f.HandleActionAttack(act)
	case types.ACTION_EFFECT:
		f.HandleActionEffect(act)
	case types.ACTION_DEFEND:
		entity := f.Entities[act.Source].Entity

		if entity.GetFlags()&types.ENTITY_AUTO == 0 {
			entity.(types.PlayerEntity).SetDefendingState(true)

			f.SendMessage(f.GetChannelId(),
				discord.MessageCreate{
					Embeds: []discord.Embed{{
						Title:       "Defensywa!",
						Description: fmt.Sprintf("%s przygotowuje się na nadchodzący atak!", entity.GetName()),
						Color:       0x00ff00,
					}},
				},
				false,
			)
		}
	case types.ACTION_SKILL:
		f.HandleActionSkill(act)
	case types.ACTION_DMG:
		f.HandleActionDamage(act)
	case types.ACTION_COUNTER:
		f.HandleActionCounter(act)
	case types.ACTION_ITEM:
		f.HandleActionItem(act)
	case types.ACTION_RUN:
		f.HandleActionRun(act)
	case types.ACTION_SUMMON:
		f.HandleActionSummon(act)
	}
}

func (f *Fight) HandleActionAttack(act types.Action) {
	targetEntity := f.Entities[act.Target].Entity
	var meta types.ActionDamage

	if tempMeta, ok := act.Meta.(types.ActionDamage); !ok {
		meta = types.ActionDamage{
			Damage: []types.Damage{
				{
					Value: f.Entities[act.Source].Entity.GetStat(types.STAT_AD),
					Type:  types.DMG_PHYSICAL,
				},
			},
			CanDodge: true,
		}
	} else {
		meta = tempMeta
	}

	if targetEntity.GetCurrentHP() <= 0 {
		return
	}

	attackEffects := f.TriggerEvent(f.Entities[act.Source].Entity, f.Entities[act.Target].Entity, types.TRIGGER_ATTACK_BEFORE, meta)

	tempEffects := make([]types.Damage, 0)

	tempEffects = append(tempEffects, meta.Damage...)

	for _, effect := range attackEffects {
		if attackTriggerMeta, ok := effect.(types.AttackTriggerMeta); ok {
			for _, tempEffect := range attackTriggerMeta.Effects {
				tempEffects = append(tempEffects, types.Damage{
					Value:     tempEffect.Value,
					Type:      tempEffect.Type,
					IsPercent: tempEffect.Percent,
				})
			}
		}
	}

	constDamage := OverallDamage(tempEffects)

	dmgDealt, dodged := targetEntity.TakeDMGOrDodge(types.ActionDamage{Damage: constDamage, CanDodge: meta.CanDodge})

	f.TriggerAttackEffect(dodged, dmgDealt, tempEffects, AttackEmbedMeta{
		Title:      "Atak!",
		TextIfHit:  "%s zaatakował %s.",
		TextIfMiss: "%s chciał zaatakować %s, ale nie trafił.",
	}, AttackMeta{
		Source:  f.Entities[act.Source].Entity,
		Target:  f.Entities[act.Target].Entity,
		IsSkill: false,

		EventHitAfterSource:  types.TRIGGER_ATTACK_HIT,
		EventHitAfterTarget:  types.TRIGGER_ATTACK_GOT_HIT,
		EventMissAfterSource: types.TRIGGER_ATTACK_MISS,
		EventMissAfterTarget: types.TRIGGER_NONE,
	})
}

func (f *Fight) HandleActionEffect(act types.Action) {
	meta := act.Meta.(types.ActionEffect)

	if meta.Duration == 0 {
		if meta.Effect == types.EFFECT_HEAL {
			healValue := meta.Value

			if act.Source != act.Target {
				healValue = utils.PercentOf(healValue, 100+f.Entities[act.Source].Entity.GetStat(types.STAT_HEAL_POWER))

				f.TriggerEvent(f.Entities[act.Source].Entity, f.Entities[act.Target].Entity, types.TRIGGER_HEAL_OTHER, healValue)
			} else {
				f.TriggerEvent(f.Entities[act.Source].Entity, f.Entities[act.Source].Entity, types.TRIGGER_HEAL_SELF, healValue)
			}

			f.Entities[act.Target].Entity.Heal(healValue)
			return
		}

		f.Entities[act.Target].Entity.ApplyEffect(meta)
		return
	}

	if meta.Effect == types.EFFECT_SHIELD && act.Source != act.Target {
		healValue := utils.PercentOf(meta.Value, 100+f.Entities[act.Source].Entity.GetStat(types.STAT_HEAL_POWER))

		f.TriggerEvent(f.Entities[act.Source].Entity, f.Entities[act.Target].Entity, types.TRIGGER_HEAL_OTHER, healValue)
	}

	if meta.Effect == types.EFFECT_TAUNT {
		for _, entity := range f.GetEnemiesFor(act.Target) {
			entity.ApplyEffect(
				types.ActionEffect{Effect: types.EFFECT_TAUNTED, Duration: meta.Duration, Meta: act.Target},
			)
		}

		return
	}

	f.Entities[act.Target].Entity.ApplyEffect(meta)
}

func (f *Fight) HandleActionSkill(act types.Action) {
	sourceEntity := f.Entities[act.Source].Entity
	sourceEntityFlags := sourceEntity.GetFlags()

	isPlayer := sourceEntityFlags&types.ENTITY_AUTO == 0
	isSummon := sourceEntityFlags&types.ENTITY_SUMMON != 0 || sourceEntityFlags&types.ENTITY_AUTO == 0

	if !isPlayer || !isSummon {
		return
	}

	skillUsageMeta := act.Meta.(types.ActionSkillMeta)

	skillUpgrades := 0

	if isPlayer {
		skillUpgrades = sourceEntity.(types.PlayerEntity).GetUpgrades(skillUsageMeta.Lvl)
	} else {
		ownerUuid := f.SummonMap[act.Source].Owner

		skillUpgrades = f.Entities[ownerUuid].Entity.(types.PlayerEntity).GetUpgrades(skillUsageMeta.Lvl)
	}

	skill := sourceEntity.(types.PlayerEntity).GetLvlSkill(skillUsageMeta.Lvl)

	trigger := skill.GetTrigger()
	skillCost := skill.GetCost()
	cooldown := skill.GetCD()

	if skill.IsLevelSkill() {
		trigger = skill.(types.PlayerSkillUpgradable).GetUpgradableTrigger(skillUpgrades)
		skillCost = skill.(types.PlayerSkillUpgradable).GetUpgradableCost(skillUpgrades)
		cooldown = skill.(types.PlayerSkillUpgradable).GetCooldown(skillUpgrades)
	}

	if sourceEntity.GetCurrentMana() < skillCost {
		f.SendMessage(
			f.GetChannelId(),
			discord.NewMessageCreateBuilder().SetContent("Nie masz many na użycie tej umiejętności").Build(),
			false,
		)

		return
	}

	if trigger.Type != types.TRIGGER_ACTIVE {
		f.SendMessage(
			f.GetChannelId(),
			discord.NewMessageCreateBuilder().SetContent("Nie można użyć tej umiejętności").Build(),
			false,
		)

		return
	}

	sourceEntity.RestoreMana(-skillCost)
	sourceEntity.(types.PlayerEntity).SetLvlCD(skillUsageMeta.Lvl, cooldown)

	if skillUsageMeta.Lvl%10 == 0 {
		sourceEntity.TriggerEvent(types.TRIGGER_CAST_ULT, types.EventData{
			Source: sourceEntity,
			Target: f.Entities[act.Target].Entity,
			Fight:  f,
		}, nil)
	}

	if len(skillUsageMeta.Targets) > 0 {
		for _, target := range skillUsageMeta.Targets {
			if skill.IsLevelSkill() {
				skill.(types.PlayerSkillUpgradable).UpgradableExecute(sourceEntity.(types.PlayerEntity), f.Entities[target].Entity, f, nil)
			} else {
				skill.Execute(sourceEntity.(types.PlayerEntity), f.Entities[target].Entity, f, nil)
			}
		}
	} else {
		if skill.IsLevelSkill() {
			skill.(types.PlayerSkillUpgradable).UpgradableExecute(sourceEntity.(types.PlayerEntity), f.Entities[act.Target].Entity, f, nil)
		} else {
			skill.Execute(sourceEntity.(types.PlayerEntity), f.Entities[act.Target].Entity, f, nil)
		}
	}

	f.SendMessage(
		f.GetChannelId(),
		discord.MessageCreate{
			Embeds: []discord.Embed{{
				Title:       "Skill!",
				Description: fmt.Sprintf("%s użył `%s`!", sourceEntity.GetName(), skill.GetName()),
				Color:       0x00ff00,
			}},
		},
		false,
	)
}

func (f *Fight) HandleActionDamage(act types.Action) {
	targetEntity := f.Entities[act.Target].Entity
	meta := act.Meta.(types.ActionDamage)

	if targetEntity.GetCurrentHP() <= 0 {
		return
	}

	damageEffects := f.TriggerEvent(f.Entities[act.Source].Entity, f.Entities[act.Target].Entity, types.TRIGGER_DAMAGE_BEFORE, meta)

	tempEffects := make([]types.Damage, 0)

	tempEffects = append(tempEffects, meta.Damage...)

	for _, effect := range damageEffects {
		if damageTriggerMeta, ok := effect.(types.DamageTriggerMeta); ok {
			for _, tempEffect := range damageTriggerMeta.Effects {
				tempEffects = append(tempEffects, types.Damage{
					Value:     tempEffect.Value,
					Type:      tempEffect.Type,
					IsPercent: tempEffect.Percent,
				})
			}
		}
	}

	constDamage := OverallDamage(tempEffects)

	dmgDealt, dodged := f.Entities[act.Target].Entity.TakeDMGOrDodge(
		types.ActionDamage{Damage: constDamage, CanDodge: meta.CanDodge},
	)

	f.TriggerAttackEffect(dodged, dmgDealt, tempEffects, AttackEmbedMeta{
		Title:      "Obrażenia!",
		TextIfHit:  "%s zadał obrażenia %s.",
		TextIfMiss: "%s chciał zadać obrażenia %s, ale nie trafił.",
	}, AttackMeta{Source: f.Entities[act.Source].Entity, Target: f.Entities[act.Target].Entity, IsSkill: true,
		EventHitAfterSource:  types.TRIGGER_DAMAGE,
		EventHitAfterTarget:  types.TRIGGER_DAMAGE_GOT_HIT,
		EventMissAfterSource: types.TRIGGER_NONE,
		EventMissAfterTarget: types.TRIGGER_NONE,
	})
}

func OverallDamage(damage []types.Damage) []types.Damage {
	constDamage := []types.Damage{
		{Value: 0, Type: types.DMG_PHYSICAL},
		{Value: 0, Type: types.DMG_MAGICAL},
		{Value: 0, Type: types.DMG_TRUE},
	}

	percentageDamage := []types.Damage{
		{Value: 0, Type: types.DMG_PHYSICAL},
		{Value: 0, Type: types.DMG_MAGICAL},
		{Value: 0, Type: types.DMG_TRUE},
	}

	for _, dmg := range damage {
		if dmg.IsPercent {
			percentageDamage[dmg.Type].Value += dmg.Value
		} else {
			constDamage[dmg.Type].Value += dmg.Value
		}
	}

	for _, elt := range percentageDamage {
		constDamage[elt.Type].Value = utils.PercentOf(constDamage[elt.Type].Value, 100+elt.Value)
	}

	return constDamage
}

func (f *Fight) HandleActionCounter(act types.Action) {
	sourceEntity := f.Entities[act.Source].Entity

	meta := act.Meta.(types.ActionDamage)

	attackEffects := f.TriggerEvent(f.Entities[act.Source].Entity, f.Entities[act.Target].Entity, types.TRIGGER_ATTACK_BEFORE, meta)

	tempEffects := make([]types.Damage, 0)

	tempEffects = append(tempEffects, meta.Damage...)

	for _, effect := range attackEffects {
		if attackTriggerMeta, ok := effect.(types.AttackTriggerMeta); ok {
			for _, tempEffect := range attackTriggerMeta.Effects {
				tempEffects = append(tempEffects, types.Damage{
					Value:     tempEffect.Value,
					Type:      tempEffect.Type,
					IsPercent: tempEffect.Percent,
				})
			}
		}
	}

	constDamage := OverallDamage(tempEffects)

	dmgDealt, dodged := f.Entities[act.Target].Entity.TakeDMGOrDodge(
		types.ActionDamage{Damage: constDamage, CanDodge: true},
	)

	f.TriggerAttackEffect(dodged, dmgDealt, tempEffects, AttackEmbedMeta{
		Title:      "Kontra!",
		TextIfHit:  "%s zaatakował %s.",
		TextIfMiss: "%s nie trafił.",
	}, AttackMeta{Source: sourceEntity, Target: f.Entities[act.Target].Entity, IsSkill: false,
		EventHitAfterSource:  types.TRIGGER_ATTACK_HIT,
		EventHitAfterTarget:  types.TRIGGER_ATTACK_GOT_HIT,
		EventMissAfterSource: types.TRIGGER_ATTACK_MISS,
		EventMissAfterTarget: types.TRIGGER_NONE,
	})
}

func (f *Fight) TriggerAttackEffect(dodged bool, damage map[types.DamageType]int, rawDamage []types.Damage, embedMeta AttackEmbedMeta, meta AttackMeta) {
	tempEmbed := discord.NewEmbedBuilder().SetTitle(embedMeta.Title)

	if !dodged {
		dmgSum := damage[0] + damage[1] + damage[2]

		vampType := types.STAT_ATK_VAMP

		if meta.IsSkill {
			vampType = types.STAT_OMNI_VAMP
		}

		f.TriggerVampEvent(meta.Source, vampType, dmgSum, tempEmbed)

		tempEmbed.
			SetFooterTextf(embedMeta.TextIfHit+"%s ma teraz %d HP", meta.Source.GetName(), meta.Target.GetName(), meta.Target.GetName(), meta.Target.GetCurrentHP()).
			SetDescriptionf("Zadano łącznie %d obrażeń", dmgSum).SetColor(0x00ff00).
			AddField("Obrażenia", DamageSummary(rawDamage), false)

		if meta.EventHitAfterSource != types.TRIGGER_NONE {
			f.TriggerEvent(meta.Source, meta.Target, meta.EventHitAfterSource, nil)
		}

		if meta.EventHitAfterTarget != types.TRIGGER_NONE {
			f.TriggerEvent(meta.Target, meta.Source, meta.EventHitAfterTarget, nil)
		}

		f.TriggerCounter(meta.Source, meta.Target)
	} else {
		if meta.EventMissAfterSource != types.TRIGGER_NONE {
			f.TriggerEvent(meta.Source, meta.Target, meta.EventMissAfterSource, nil)
		}

		if meta.EventMissAfterTarget != types.TRIGGER_NONE {
			f.TriggerEvent(meta.Target, meta.Source, meta.EventMissAfterTarget, nil)
		}

		tempEmbed.SetDescriptionf(embedMeta.TextIfMiss, meta.Source.GetName(), meta.Target.GetName()).SetColor(0xff0000)
	}

	tempEmbed.AddFields(embedMeta.AdditionalFields...)

	f.SendMessage(
		f.GetChannelId(),
		discord.MessageCreate{
			Embeds: []discord.Embed{tempEmbed.Build()},
		},
		false,
	)
}

type AttackEmbedMeta struct {
	Title string
	/*source <action> target.*/
	TextIfHit string
	/*source <action> target.*/
	TextIfMiss       string
	AdditionalFields []discord.EmbedField
}

type AttackMeta struct {
	Source  types.Entity
	Target  types.Entity
	IsSkill bool

	EventMissAfterSource types.SkillTrigger
	EventMissAfterTarget types.SkillTrigger
	EventHitAfterSource  types.SkillTrigger
	EventHitAfterTarget  types.SkillTrigger
}

func (f *Fight) TriggerVampEvent(source types.Entity, vampType types.Stat, dmg int, embed *discord.EmbedBuilder) {
	vampValue := source.GetStat(vampType)

	if vampValue <= 0 {
		return
	}

	value := utils.PercentOf(dmg, vampValue)

	source.Heal(value)

	source.TriggerEvent(
		types.TRIGGER_HEAL_SELF,
		types.EventData{Source: source, Target: source, Fight: f},
		value,
	)

	if embed != nil {
		embed.AddField("Wampiryzm!", fmt.Sprintf("%s dodatkowo wyleczył się o %d", source.GetName(), value), false)
	}
}

// Target as in target of the damage
func (f *Fight) TriggerCounter(source, target types.Entity) {
	if types.HasFlag(target.GetFlags(), types.ENTITY_AUTO) {
		return
	}

	if !target.(types.PlayerEntity).GetDefendingState() {
		return
	}

	if utils.RandomNumber(0, 100) < target.GetStat(types.STAT_AGL) {
		counterDmg := utils.PercentOf(source.GetStat(types.STAT_AD), 70)
		counterDmg += utils.PercentOf(source.GetStat(types.STAT_DEF), 15)
		counterDmg += utils.PercentOf(source.GetStat(types.STAT_MR), 15)

		f.HandleAction(types.Action{
			Event:  types.ACTION_COUNTER,
			Source: target.GetUUID(),
			Target: source.GetUUID(),
			Meta: types.ActionDamage{
				Damage:   []types.Damage{{Value: counterDmg, Type: types.DMG_PHYSICAL}},
				CanDodge: true,
			},
		})
	}
}

func DamageSummary(dmgList []types.Damage) string {
	dmgText := ""

	for _, dmg := range dmgList {
		if dmg.Value == 0 {
			continue
		}

		dmgType := "fizycznych"

		switch dmg.Type {
		case types.DMG_MAGICAL:
			dmgType = "magicznych"
		case types.DMG_TRUE:
			dmgType = "nieuchronnych"
		}

		if dmg.IsPercent {
			dmgText += fmt.Sprintf("- %d%% obrażeń %s\n", dmg.Value, dmgType)
		} else {
			dmgText += fmt.Sprintf("- %d obrażeń %s\n", dmg.Value, dmgType)
		}
	}

	return dmgText
}

func (f *Fight) HandleActionItem(act types.Action) {
	sourceEntity := f.Entities[act.Source].Entity.(types.PlayerEntity)

	itemMeta := act.Meta.(types.ActionItemMeta)

	var item *types.PlayerItem

	for _, invItem := range sourceEntity.GetAllItems() {
		if invItem.UUID == itemMeta.Item {
			item = invItem
			break
		}
	}

	if act.Target == uuid.Nil {
		for _, target := range itemMeta.Targets {
			sourceEntity.UseItem(item.UUID, f.Entities[target].Entity, f)
		}
	} else {
		sourceEntity.UseItem(item.UUID, f.Entities[act.Target].Entity, f)
	}

	f.SendMessage(
		f.GetChannelId(),
		discord.MessageCreate{Embeds: []discord.Embed{{
			Title:       "Przedmiot!",
			Description: fmt.Sprintf("%s użył %s!\nEfekt: %s", sourceEntity.GetName(), item.Name, item.Description),
		}}},
		false,
	)
}

func (f *Fight) HandleActionRun(act types.Action) {
	entity := f.Entities[act.Source].Entity
	side := f.Entities[act.Source].Side

	if utils.RandomNumber(0, 100) < entity.GetStat(types.STAT_AGL) {
		f.SendMessage(
			f.GetChannelId(),
			discord.MessageCreate{Embeds: []discord.Embed{{
				Title:       "Ucieczka!",
				Description: fmt.Sprintf("%s próbował uciec i mu się to nie udało", entity.GetName()),
				Color:       0xff0000,
			}}},
			false,
		)

		return
	}

	f.Entities[act.Source].Entity.(types.PlayerEntity).ClearFight()

	delete(f.Entities, act.Source)

	entities := f.FromSide(side)

	count := 0

	for _, entity := range entities {
		if entity.GetCurrentHP() > 0 && !types.HasFlag(entity.GetFlags(), types.ENTITY_AUTO) {
			count++
		}
	}

	f.SendMessage(
		f.GetChannelId(),
		discord.MessageCreate{Embeds: []discord.Embed{{
			Title:       "Ucieczka!",
			Description: fmt.Sprintf("%s próbował uciec i mu się to udało", entity.GetName()),
			Color:       0x00ff00,
		}}},
		false,
	)

	if count == 0 {
		f.ExternalChannel <- FightEndMsg{RunAway: true}
	}
}

func (f *Fight) HandleActionSummon(act types.Action) {
	actionMeta := act.Meta.(types.ActionSummon)
	sourceEntity := f.Entities[act.Source]

	f.SendMessage(
		f.GetChannelId(),
		discord.MessageCreate{Embeds: []discord.Embed{{
			Title:       "Przywołanie!",
			Description: fmt.Sprintf("%s przywołał %s", sourceEntity.Entity.GetName(), actionMeta.Entity.GetName()),
			Color:       0x00ff00,
		}}},
		false,
	)

	newEntityUUID := actionMeta.Entity.GetUUID()

	if types.HasFlag(actionMeta.Flags, types.SUMMON_FLAG_EXPIRE) {
		f.ExpireMap[newEntityUUID] = actionMeta.ExpireTimer
	}

	f.SummonMap[newEntityUUID] = SummonEntityMeta{Owner: act.Source, Type: actionMeta.EntityType}
	f.Entities[newEntityUUID] = &EntityEntry{Entity: actionMeta.Entity, Side: sourceEntity.Side}
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

func (f *Fight) ForwardSummons() {
	for entity, exp := range f.ExpireMap {
		f.ExpireMap[entity] = exp - 1

		entityData := f.Entities[entity].Entity

		if entityData.GetCurrentHP() <= 0 {
			delete(f.ExpireMap, entity)

			f.ExternalChannel <- SummonDied{Entity: entity}

			continue
		}

		if exp <= 0 {
			delete(f.ExpireMap, entity)

			f.ExternalChannel <- SummonExpired{Entity: entity}
		}
	}
}

func (f *Fight) RequestAction(uid uuid.UUID) {
	record, exists := f.Entities[uid]

	if !exists || record.Entity.GetCurrentHP() <= 0 {
		return
	}

	if f.IsFinished() {
		return
	}

	record.Entity.TriggerEvent(types.TRIGGER_TURN, types.EventData{
		Source: record.Entity,
		Target: record.Entity,
		Fight:  f,
	}, nil)

	record.Turn++

	if !types.HasFlag(record.Entity.GetFlags(), types.ENTITY_AUTO) {
		record.Entity.(types.PlayerEntity).SetDefendingState(false)

		record.Entity.(types.PlayerEntity).ReduceCooldowns(types.TRIGGER_TURN)

		f.ExternalChannel <- FightActionNeededMsg{Entity: uid}

		tempAction := <-f.PlayerActions

		f.HandleAction(tempAction)

		for tempAction.ConsumeTurn != nil && !*tempAction.ConsumeTurn {
			f.ExternalChannel <- FightActionNeededMsg{Entity: uid}

			tempAction = <-f.PlayerActions

			f.HandleAction(tempAction)
		}
	} else {
		if record.Entity.GetEffectByType(types.EFFECT_STUN) == nil {
			for _, action := range record.Entity.Action(f) {
				f.HandleAction(action)
			}

			return
		}

		f.SendMessage(
			f.GetChannelId(),
			discord.MessageCreate{Embeds: []discord.Embed{{
				Title:       "Efekt!",
				Description: fmt.Sprintf("%s jest ogłuszony, pomijamy!", record.Entity.GetName()),
			}}},
			false,
		)
	}

	record.Entity.TriggerAllEffects()
	record.Entity.TriggerTempSkills()
}

func (f *Fight) Run() {
	f.ExternalChannel <- FightStartMsg{}

	for len(f.SidesLeft()) > 1 {
		f.ForwardSummons()

		for uuid, val := range f.Entities {
			val.Speed += f.Entities[uuid].Entity.GetStat(types.STAT_SPD)

			for val.Speed >= data.WorldConfig.SpeedGauge {
				val.Speed -= data.WorldConfig.SpeedGauge

				f.RequestAction(uuid)
			}
		}
	}

	f.ExternalChannel <- FightEndMsg{}
}

func (f *Fight) SendMessage(channelId string, content discord.MessageCreate, dm bool) {
	f.DiscordChannel <- types.DiscordSendMsg{
		Data: types.DiscordMessageStruct{ChannelID: channelId, MessageContent: content, DM: dm},
	}
}

func (f *Fight) GetEntity(uuid uuid.UUID) types.Entity {
	return f.Entities[uuid].Entity
}

func (f *Fight) FromSide(side int) []types.Entity {
	return f.FilteredEntities(func(entity EntityEntry) bool { return entity.Side == side })
}

func (f *Fight) RemoveEventHandler(uuid uuid.UUID) {
	delete(f.EventHandlers, uuid)
}

func (f *Fight) GetTurnFor(uuid uuid.UUID) int {
	return f.Entities[uuid].Turn
}

func (f *Fight) IsFinished() bool {
	return len(f.SidesLeft()) <= 1
}

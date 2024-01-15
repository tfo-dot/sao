package battle

import (
	"fmt"
	"sao/types"
	"sao/utils"
	"sao/world/calendar"
	"sort"

	"github.com/google/uuid"
)

type Fight struct {
	Entities        EntityMap
	SpeedMap        map[uuid.UUID]int
	StartTime       *calendar.Calendar
	ActionChannel   chan Action
	ExternalChannel chan []byte
	Effects         []ActionEffect
}

type EntityMap map[uuid.UUID]EntityEntry

type EntityEntry struct {
	Entity Entity
	Side   int
}

func (f Fight) GetEntity(uuid uuid.UUID) (Entity, int) {
	entry, ok := f.Entities[uuid]

	if ok {
		return entry.Entity, entry.Side
	}

	return nil, -1
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

func (f *Fight) DispatchActionAttack(act Action) (int, bool) {
	sourceEntity := f.Entities[act.Source]

	tempMeta := act.Meta

	if _, ok := tempMeta.(ActionDamage); !ok {
		tempMeta = Damage{
			Value:    sourceEntity.Entity.GetATK(),
			Type:     DMG_PHYSICAL,
			CanDodge: true,
		}.ToActionMeta()
	}

	meta := tempMeta.(ActionDamage)

	if meta.CanDodge && f.Entities[act.Target].Entity.CanDodge() {
		atk, dodged := f.Entities[act.Target].Entity.(DodgeEntity).TakeDMGOrDodge(meta)

		if !dodged && !sourceEntity.Entity.IsAuto() {
			for _, skill := range sourceEntity.Entity.(PlayerEntity).GetAllSkills() {
				if skill.Trigger.Type != types.TRIGGER_PASSIVE || skill.Trigger.Event.TriggerType != types.TRIGGER_HIT {
					continue
				}

				targets := f.FindValidTargets(sourceEntity.Entity.GetUUID(), *skill.Trigger.Event)

				if skill.Trigger.Event.TargetCount != -1 {
					count := skill.Trigger.Event.TargetCount

					if count > len(targets) {
						count = len(targets)
					}

					targets = targets[:count]
				}

				for _, target := range targets {
					skill.Action(sourceEntity.Entity, f.Entities[target].Entity, f)
				}

			}
		}

		return atk, dodged
	} else {
		if !sourceEntity.Entity.IsAuto() {
			for _, skill := range sourceEntity.Entity.(PlayerEntity).GetAllSkills() {
				if skill.Trigger.Type == types.TRIGGER_ACTIVE {
					continue
				}

				if skill.Trigger.Event.TriggerType != types.TRIGGER_HIT {
					continue
				}

				targets := f.FindValidTargets(sourceEntity.Entity.GetUUID(), *skill.Trigger.Event)

				if skill.Trigger.Event.TargetCount != -1 {
					count := skill.Trigger.Event.TargetCount

					if count > len(targets) {
						count = len(targets)
					}

					targets = targets[:count]
				}

				for _, target := range targets {
					skill.Action(sourceEntity.Entity, f.Entities[target].Entity, f)
				}
			}
		}

		return f.Entities[act.Target].Entity.TakeDMG(meta), false
	}
}

func (f *Fight) HandleAction(act Action) {
	switch act.Event {
	case ACTION_ATTACK:
		dmgDealt, dodged := f.DispatchActionAttack(act)

		entity := f.Entities[act.Source]

		if entity.Entity.HasEffect(EFFECT_VAMP) && !dodged {
			effect := entity.Entity.GetEffect(EFFECT_VAMP)

			entity.Entity.Heal(utils.PercentOf(dmgDealt, effect.Value))
		}

		if !dodged && !entity.Entity.IsAuto() {
			for _, skill := range entity.Entity.(PlayerEntity).GetAllSkills() {
				if skill.Trigger.Type == types.TRIGGER_ACTIVE {
					continue
				}

				if skill.Trigger.Event.TriggerType != types.TRIGGER_ATTACK {
					continue
				}

				targets := f.FindValidTargets(entity.Entity.GetUUID(), *skill.Trigger.Event)

				if skill.Trigger.Event.TargetCount != -1 {
					count := skill.Trigger.Event.TargetCount

					if count > len(targets) {
						count = len(targets)
					}

					targets = targets[:count]
				}

				for _, target := range targets {
					skill.Action(entity.Entity, f.Entities[target].Entity, f)
				}
			}
		}

		targetEntity := f.Entities[act.Target]

		if !targetEntity.Entity.IsAuto() {
			if targetEntity.Entity.(PlayerEntity).GetDefendingState() {
				if utils.RandomNumber(0, 100) < targetEntity.Entity.GetAGL() {

					f.HandleAction(Action{
						Event:  ACTION_ATTACK,
						Source: act.Target,
						Target: act.Source,
						Meta: Damage{
							Value:    targetEntity.Entity.GetATK(),
							Type:     DMG_PHYSICAL,
							CanDodge: true,
						}.ToActionMeta(),
					})
				}
			}
		}

	case ACTION_EFFECT:
		meta := act.Meta.(ActionEffect)

		if meta.Duration == 0 {
			switch meta.Effect {
			case EFFECT_HEAL:
				f.Entities[act.Target].Entity.Heal(meta.Value)
			case EFFECT_STAT_INC:
				f.Entities[act.Target].Entity.ApplyEffect(meta)
			}
		}

		f.Entities[act.Target].Entity.ApplyEffect(act.Meta.(ActionEffect))
	case ACTION_DEFEND:
		entity := f.Entities[act.Source]

		if !entity.Entity.IsAuto() {
			for _, skill := range entity.Entity.(PlayerEntity).GetAllSkills() {
				if skill.Trigger.Type == types.TRIGGER_ACTIVE {
					continue
				}

				if skill.Trigger.Event.TriggerType != types.TRIGGER_DEFEND {
					continue
				}

				targets := f.FindValidTargets(entity.Entity.GetUUID(), *skill.Trigger.Event)

				if skill.Trigger.Event.TargetCount != -1 {
					count := skill.Trigger.Event.TargetCount

					if count > len(targets) {
						count = len(targets)
					}

					targets = targets[:count]
				}

				for _, target := range targets {
					skill.Action(entity.Entity, f.Entities[target].Entity, f)
				}
			}

			entity.Entity.(PlayerEntity).SetDefendingState(true)

			f.HandleAction(Action{
				Event:  ACTION_EFFECT,
				Source: act.Source,
				Target: act.Source,
				Meta: ActionEffect{
					Effect:   EFFECT_STAT_INC,
					Duration: 1,
					Meta: &map[string]interface{}{
						"stat":     STAT_DEF,
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
						"stat":     STAT_DEF,
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
				if skill.Trigger.Type == types.TRIGGER_ACTIVE {
					continue
				}

				if skill.Trigger.Event.TriggerType != types.TRIGGER_MANA {
					continue
				}

				if skill.Trigger.Event.Meta["value"].(int) > sourceEntity.Entity.GetCurrentMana() {
					continue
				}

				targets := f.FindValidTargets(sourceEntity.Entity.GetUUID(), *skill.Trigger.Event)

				if skill.Trigger.Event.TargetCount != -1 {
					count := skill.Trigger.Event.TargetCount

					if count > len(targets) {
						count = len(targets)
					}

					targets = targets[:count]
				}

				for _, target := range targets {
					targetEntity := f.Entities[target]

					beforeSkillHP := targetEntity.Entity.GetCurrentHP()

					skill.Action(sourceEntity.Entity, targetEntity, f)

					//Check if it's dmg skill so it doesn't trigger on heal/barrier etc
					if !targetEntity.Entity.IsAuto() && beforeSkillHP > targetEntity.Entity.GetCurrentHP() {
						if targetEntity.Entity.(PlayerEntity).GetDefendingState() {
							if utils.RandomNumber(0, 100) < targetEntity.Entity.GetAGL() {
								f.HandleAction(Action{
									Event:  ACTION_ATTACK,
									Source: act.Target,
									Target: act.Source,
									Meta: Damage{
										Value:    targetEntity.Entity.GetATK(),
										Type:     DMG_PHYSICAL,
										CanDodge: true,
									}.ToActionMeta(),
								})
							}
						}
					}
				}
			}
		}

	case ACTION_DMG:
		//TODO Redirect all the damage here
		sourceEntity := f.Entities[act.Source]
		targetEntity := f.Entities[act.Target]
		meta := act.Meta.(ActionDamage)

		if targetEntity.Entity.GetCurrentHP() <= 0 {
			return
		}

		if meta.CanDodge && targetEntity.Entity.CanDodge() {
			targetEntity.Entity.(DodgeEntity).TakeDMGOrDodge(meta)
		} else {
			targetEntity.Entity.TakeDMG(meta)
		}

		if targetEntity.Entity.GetCurrentHP() <= 0 {
			if !sourceEntity.Entity.IsAuto() {
				for _, skill := range sourceEntity.Entity.(PlayerEntity).GetAllSkills() {
					if skill.Trigger.Type == types.TRIGGER_ACTIVE {
						continue
					}

					if skill.Trigger.Event.TriggerType != types.TRIGGER_EXECUTE {
						continue
					}

					targets := f.FindValidTargets(sourceEntity.Entity.GetUUID(), *skill.Trigger.Event)

					if skill.Trigger.Event.TargetCount != -1 {
						count := skill.Trigger.Event.TargetCount

						if count > len(targets) {
							count = len(targets)
						}

						targets = targets[:count]
					}

					for _, target := range targets {
						if f.Entities[target].Entity.GetCurrentHP() <= 0 {
							skill.Action(sourceEntity.Entity, f.Entities[target].Entity, f)
						}
					}
				}
			}
		}

		if !targetEntity.Entity.IsAuto() {
			for _, skill := range targetEntity.Entity.(PlayerEntity).GetAllSkills() {
				if skill.Trigger.Type == types.TRIGGER_ACTIVE {
					continue
				}

				if skill.Trigger.Event.TriggerType != types.TRIGGER_HEALTH {
					continue
				}

				hpValue := 0

				if skill.Trigger.Event.Meta["value"] != nil {
					hpValue = skill.Trigger.Event.Meta["value"].(int)
				} else {
					hpValue = (skill.Trigger.Event.Meta["percent"].(int) * targetEntity.Entity.GetMaxHP() / 100)
				}

				if hpValue > targetEntity.Entity.GetCurrentHP() {
					continue
				}

				targets := f.FindValidTargets(sourceEntity.Entity.GetUUID(), *skill.Trigger.Event)

				if skill.Trigger.Event.TargetCount != -1 {
					count := skill.Trigger.Event.TargetCount

					if count > len(targets) {
						count = len(targets)
					}

					targets = targets[:count]
				}

				for _, target := range targets {
					if f.Entities[target].Entity.GetCurrentHP() <= 0 {
						skill.Action(sourceEntity.Entity, f.Entities[target].Entity, f)
					}
				}
			}
		}

	default:
		fmt.Printf("Unknown action %d\n", act.Event)
		panic("Not implemented (actions)")
	}
}

func (f *Fight) Init(currentTime *calendar.Calendar) {
	f.SpeedMap = make(map[uuid.UUID]int)
	f.ActionChannel = make(chan Action, 10)

	for _, entity := range f.Entities {
		f.SpeedMap[entity.Entity.GetUUID()] = entity.Entity.GetSPD()
	}

	f.StartTime = currentTime.Copy()
	f.ExternalChannel = make(chan []byte)

	//FIGHT START EVENT
	for _, entity := range f.Entities {
		if !entity.Entity.IsAuto() {
			for _, skill := range entity.Entity.(PlayerEntity).GetAllSkills() {
				if skill.Trigger.Type == types.TRIGGER_ACTIVE {
					continue
				}

				if skill.Trigger.Event.TriggerType != types.TRIGGER_FIGHT_START {
					continue
				}

				targets := f.FindValidTargets(entity.Entity.GetUUID(), *skill.Trigger.Event)

				if skill.Trigger.Event.TargetCount != -1 {
					count := skill.Trigger.Event.TargetCount

					if count > len(targets) {
						count = len(targets)
					}

					targets = targets[:count]
				}

				for _, target := range targets {
					skill.Action(entity.Entity, f.Entities[target].Entity, f)
				}

			}
		}
	}
}

func (f *Fight) FindValidTargets(source uuid.UUID, trigger types.EventTriggerDetails) []uuid.UUID {
	sourceEntity, sourceSide := f.GetEntity(source)

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
	for len(f.Entities.SidesLeft()) > 1 {
		turnList := make([]uuid.UUID, 0)

		for uuid, speed := range f.SpeedMap {
			entity, _ := f.GetEntity(uuid)

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

		for _, uuid := range turnList {
			entity, _ := f.GetEntity(uuid)

			if entity.GetCurrentHP() == 0 {
				continue
			}

			if !entity.IsAuto() {
				for _, skill := range entity.(PlayerEntity).GetAllSkills() {
					if skill.Trigger.Type == types.TRIGGER_ACTIVE {
						continue
					}

					if skill.Trigger.Event.TriggerType != types.TRIGGER_TURN {
						continue
					}

					targets := f.FindValidTargets(uuid, *skill.Trigger.Event)

					if skill.Trigger.Event.TargetCount != -1 {
						count := skill.Trigger.Event.TargetCount

						if count > len(targets) {
							count = len(targets)
						}

						targets = targets[:count]
					}

					for _, target := range targets {
						skill.Action(entity, f.Entities[target].Entity, f)
					}
				}

				fmt.Printf("Entity %s is taking action\n", entity.GetName())

				bytes, err := entity.GetUUID().MarshalBinary()

				if err != nil {
					panic(err)
				}

				packet := make([]byte, 1+len(bytes))
				packet[0] = byte(MSG_ACTION_NEEDED)
				copy(packet[1:], bytes)

				f.ExternalChannel <- packet

				f.HandleAction(<-f.ActionChannel)
			} else {
				actionNum := entity.Action(f)

				for i := 0; i < actionNum; i++ {
					f.HandleAction(<-f.ActionChannel)
				}
			}

			entity.TriggerAllEffects()
		}

		for _, entry := range f.Entities {
			if entry.Entity.GetCurrentHP() <= 0 {
				fmt.Printf("Entity %s died!\n", entry.Entity.GetName())
				continue
			}

			fmt.Printf("Entity %s has %d hp left\n", entry.Entity.GetName(), entry.Entity.GetCurrentHP())
		}
	}

	//FIGHT END EVENT
	for _, entity := range f.Entities {
		if !entity.Entity.IsAuto() {
			for _, skill := range entity.Entity.(PlayerEntity).GetAllSkills() {
				if skill.Trigger.Type == types.TRIGGER_ACTIVE {
					continue
				}

				if skill.Trigger.Event.TriggerType != types.TRIGGER_FIGHT_END {
					continue
				}

				targets := f.FindValidTargets(entity.Entity.GetUUID(), *skill.Trigger.Event)

				if skill.Trigger.Event.TargetCount != -1 {
					count := skill.Trigger.Event.TargetCount

					if count > len(targets) {
						count = len(targets)
					}

					targets = targets[:count]
				}

				for _, target := range targets {
					skill.Action(entity.Entity, f.Entities[target].Entity, f)
				}

			}
		}
	}

	f.ExternalChannel <- []byte{byte(MSG_FIGHT_END)}
	close(f.ExternalChannel)
}

func (f *Fight) IsFinished() bool {
	return len(f.Entities.SidesLeft()) <= 1
}

func (f *Fight) GetEnemiesFor(uuid uuid.UUID) []Entity {
	_, userSide := f.GetEntity(uuid)

	enemiesList := make([]Entity, 0)

	for _, entry := range f.Entities {
		if entry.Side == userSide {
			continue
		}

		enemiesList = append(enemiesList, entry.Entity)
	}

	return enemiesList
}

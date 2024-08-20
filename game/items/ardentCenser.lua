--ItemID
ReservedUIDs[0] = "00000000-0000-0000-0000-000000000015"
--SkillID
ReservedUIDs[2] = "00000000-0000-0001-0000-000000000015"
--EffectID
ReservedUIDs[2] = "00000000-0000-0001-0001-000000000015"

-- Meta
UUID = ReservedUIDs[0]
Name = "Ognisty trybularz"
Description = "Leczenie i tarcze zwiększają obrażenia i prędkość sojusznika."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  HEAL_POWER = 10,
  AP = 30,
  HP = 50
}

-- Effects
Effects[0] = {
  GetName = function() return "Ognisty trybularz" end,
  GetDescription = function() return "Leczenie i tarcze zwiększają obrażenia i prędkość sojusznika." end,
  GetTrigger = function()
    return {
      Type = "PASSIVE",
      Event = "HEAL_OTHER"
    }
  end,
  GetUUID = function() return ReservedUIDs[1] end,
  Execute = function(owner, target, fightInstance, meta)
    target:AppendTempSkill({
      Value = utils.BaseAttackIncreaseSkill(
        {
          Calculate = function(meta)
            return {
              Effects = {
                {
                  Value = utils.PercentOf(owner:GetStat("STAT_AP"), 25),
                  Percent = false,
                  Type = 1,
                },
              },
            }
          end
        }
      ),
    })

    fightInstance:HandleAction({
      Event = "ACTION_EFFECT",
      Source = owner:GetUUID(),
      Target = target:GetUUID(),
      Meta = {
        Effect = "EFFECT_STAT_INC",
        Value = 0,
        Duration = 1,
        Uuid = ReservedUIDs[2],
        Meta = {
          Stat = "STAT_SPD",
          Value = 10,
          IsPercent = false,
        },
        Caster = owner:GetUUID(),
        Source = "SOURCE_ITEM",
      },
    })

    return nil
  end,
  GetEvents = function() return nil end,
  GetCD = function() return 0 end,
  GetCost = function() return 0 end
}

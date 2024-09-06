ReservedUIDs = {
  "00000000-0000-0000-0000-000000000015",
  "00000000-0000-0001-0000-000000000015",
  "00000000-0000-0001-0001-000000000015",
}

-- Meta
UUID = ReservedUIDs[1]
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
Effects = { {
  Trigger = {
    Type = "PASSIVE",
    Event = "HEAL_OTHER"
  },
  UUID = ReservedUIDs[2],
  Execute = function(owner, target, fightInstance, meta)
    ---@diagnostic disable-next-line: undefined-global
    AppendTempSkill(target, {
      Value = {
        Trigger = {
          Type = "PASSIVE",
          Event = "ATTACK_BEFORE"
        },
        Execute = function(_owner, _target, _fightInstance, _meta)
          return {
            Effects = {
              {
                ---@diagnostic disable-next-line: undefined-global
                Value = utils.PercentOf(GetStat(owner, StatsConst.STAT_AP), 25),
                Percent = false,
                Type = 1,
              },
            },
          }
        end
      },
      Expire = 2,
      AfterUsage = false,
      Either = true
    })

    ---@diagnostic disable-next-line: undefined-global
    HandleAction(fightInstance, {
      Event = "ACTION_EFFECT",
      ---@diagnostic disable-next-line: undefined-global
      Source = GetUUID(owner),
      ---@diagnostic disable-next-line: undefined-global
      Target = GetUUID(target),
      Meta = {
        Effect = "EFFECT_STAT_INC",
        Value = 0,
        Duration = 1,
        Uuid = ReservedUIDs[3],
        Meta = {
          Stat = StatsConst.STAT_SPD,
          Value = 10,
          IsPercent = false,
        },
        ---@diagnostic disable-next-line: undefined-global
        Caster = GetUUID(owner),
        ---@diagnostic disable-next-line: undefined-global
        Target = GetUUID(target),
        Source = "SOURCE_ITEM",
      },
    })

    return nil
  end,
} }
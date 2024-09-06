ReservedUIDs = {
  "00000000-0000-0000-0000-00000000000F",
  "00000000-0000-0001-0000-00000000000F",
}

-- Meta
UUID = ReservedUIDs[1]
Name = "Bransoleta kontrolera"
Description = "Nałożenie efektu CC leczy ciebie i sojusznika."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  AP = 20,
  ATK = 10,
  SPD = 5,
}

-- Effects
Effects = { {
  Trigger = {
    Type = "PASSIVE",
    Event = "APPLY_CROWD_CONTROL",
  },
  UUID = ReservedUIDs[2],
  Execute = function(owner, target, fightInstance, meta)
    local healValue = utils.PercentOf(GetStat(owner, StatsConst.STAT_AD), 15) +
        utils.PercentOf(GetStat(owner, StatsConst.STAT_AP), 15)

    ---@diagnostic disable-next-line: undefined-global
    local allies = GetAlliesFor(fightInstance, owner)

    local healTarget = allies[math.random(#allies)]

    ---@diagnostic disable-next-line: undefined-global
    HandleAction(fightInstance, {
      Event = "ACTION_EFFECT",
      ---@diagnostic disable-next-line: undefined-global
      Source = GetUUID(owner),
      ---@diagnostic disable-next-line: undefined-global
      Target = GetUUID(owner),
      Meta = {
        Effect = "EFFECT_HEAL",
        Value = 0,
        Duration = 0,
        ---@diagnostic disable-next-line: undefined-global
        Target = GetUUID(owner),
        ---@diagnostic disable-next-line: undefined-global
        Caster = GetUUID(owner),
        Source = "SOURCE_ITEM",
        Meta = {
          ---@diagnostic disable-next-line: undefined-global
          Value = healValue
        }
      },
    })

    ---@diagnostic disable-next-line: undefined-global
    HandleAction(fightInstance, {
      Event = "ACTION_EFFECT",
      ---@diagnostic disable-next-line: undefined-global
      Source = GetUUID(owner),
      ---@diagnostic disable-next-line: undefined-global
      Target = GetUUID(healTarget),
      Meta = {
        Effect = "EFFECT_HEAL",
        Value = 0,
        Duration = 0,
        ---@diagnostic disable-next-line: undefined-global
        Target = GetUUID(healTarget),
        ---@diagnostic disable-next-line: undefined-global
        Caster = GetUUID(owner),
        Source = "SOURCE_ITEM",
        Meta = {
          ---@diagnostic disable-next-line: undefined-global
          Value = healValue
        }
      },
    })

    return nil
  end,
} }

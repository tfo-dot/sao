ReservedUIDs = {
  "00000000-0000-0000-0000-000000000005",
  "00000000-0000-0001-0000-000000000005",
}

-- Meta
UUID = ReservedUIDs[1]
Name = "Wodne ostrze"
Description = "Zadawanie obrażeń leczy o brakujące zdrowie."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  ATK = 25,
  VAMP = 10,
  HP = 50,
}

-- Effects
Effects = { {
  Trigger = {
    Type = "PASSIVE",
    Event = "ATTACK_HIT",
  },
  UUID = ReservedUIDs[2],
  Execute = function(owner, target, fightInstance, meta)
    ---@diagnostic disable-next-line: undefined-global
    local addPercentage = utils.PercentOf(GetStat(owner, StatsConst.STAT_AD), 1)

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
        Caster = GetUUID(owner),
        ---@diagnostic disable-next-line: undefined-global
        Target = GetUUID(owner),
        Source = "SOURCE_ITEM",
        Meta = {
          ---@diagnostic disable-next-line: undefined-global
          Value = utils.PercentOf(GetStat(owner, StatsConst.STAT_HP) - GetCurrentHP(owner), 10 + addPercentage)
        }
      },
    })

    return nil
  end,
  CD = 10,
} }

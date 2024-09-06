ReservedUIDs = {
  "00000000-0000-0000-0000-000000000008",
  "00000000-0000-0001-0000-000000000008",
}

-- Meta
UUID = ReservedUIDs[1]
Name = "Dziedzictwo wojownika"
Description = "Zwiększa obrażenia w zależności od maks zdrowia."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  ATK = 20,
  HP = 50
}

-- Effects
Effects = { {
  Trigger = {
    Type = "PASSIVE",
    Event = "ATTACK_BEFORE",
  },
  UUID = ReservedUIDs[2],
  Execute = function(owner, target, fightInstance, meta)
    ---@diagnostic disable-next-line: undefined-global
    local dmgPercent = utils.PercentOf(GetStat(owner, StatsConst.STAT_HP_PLUS), 1)

    return {
      Effects = {
        {
          Value = dmgPercent,
          Type = 0,
          Percent = true,
        },
      },
    }
  end,
} }

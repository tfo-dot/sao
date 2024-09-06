ReservedUIDs = {
  "00000000-0000-0000-0000-00000000001B",
  "00000000-0000-0001-0000-00000000001B",
}

-- Meta
UUID = ReservedUIDs[1]
Name = "Zwiastun burzy"
Description = "Ataki zadają dodatkowe obrażenia w zależności od AP."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  ATK = 100,
  SPD = 5
}

-- Effects
Effects = { {
  Trigger = {
    Type = "PASSIVE",
    Event = "ATTACK_BEFORE"
  },
  UUID = ReservedUIDs[2],
  Execute = function(owner, target, fightInstance, meta)
    return {
      Effects = {
        {
          ---@diagnostic disable-next-line: undefined-global
          Value = utils.PercentOf(GetStat(owner, StatsConst.STAT_AP), 20),
          Type = "DMG_MAGICAL",
          Percent = false,
        },
      },
    }
  end,
} }

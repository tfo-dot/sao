ReservedUIDs = {
  "00000000-0000-0000-0000-000000000007",
  "00000000-0000-0001-0000-000000000007",
  "00000000-0000-0001-0001-000000000007",
}

-- Meta
UUID = ReservedUIDs[1]
Name = "Oblicze ataku"
Description = "Dostajesz HP w zależności od ATK."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  ATK = 20,
  DEF = 15,
  MR = 15,
}

-- Effects
Effects = { {
  Trigger = {
    Type = "PASSIVE",
    Event = "NONE",
  },
  UUID = ReservedUIDs[2],
  Events = {
    TRIGGER_UNLOCK = function(owner)
      ---@diagnostic disable-next-line: undefined-global
      AppendDerivedStat(owner, {
        Base = StatsConst.STAT_AD,
        Derived = StatsConst.STAT_HP,
        Percent = 10,
        Source = ReservedUIDs[3],
      })
    end
  },
} }

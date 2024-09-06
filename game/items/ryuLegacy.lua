ReservedUIDs = {
  "00000000-0000-0000-0000-00000000000B",
  "00000000-0000-0001-0000-00000000000B",
}

-- Meta
UUID = ReservedUIDs[1]
Name = "Dziedzictwo Ryu"
Description = "Zwiększa RES i DEF o 20%."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  HP = 150,
  DEF = 40,
  RES = 40,
}

-- Effects
Effects = { {
  Trigger = {
    Type = "PASSIVE",
    Event = "NONE",
  },
  UUID = ReservedUIDs[2],
  GetEvents = {
    TRIGGER_UNLOCK = function(owner)
      ---@diagnostic disable-next-line: undefined-global
      AppendDerivedStat(owner, {
        Base = StatsConst.STAT_DEF,
        Derived = StatsConst.STAT_DEF,
        Percent = 20,
        Source = ReservedUIDs[2],
      })

      ---@diagnostic disable-next-line: undefined-global
      AppendDerivedStat(owner, {
        Base = StatsConst.STAT_MR,
        Derived = StatsConst.STAT_MR,
        Percent = 20,
        Source = ReservedUIDs[2],
      })
    end
  },
} }

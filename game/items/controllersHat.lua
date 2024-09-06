ReservedUIDs = {
  "00000000-0000-0000-0000-000000000014",
  "00000000-0000-0001-0000-000000000014",
  "00000000-0000-0001-0001-000000000014",
}

-- Meta
UUID = ReservedUIDs[1]
Name = "Kapelusz kontrolera"
Description = "Daje siłę adaptacyjną w zależności od many."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  MANA = 5
}

-- Effects
Effects = { {
  Trigger = {
    Type = "PASSIVE",
    Event = "NONE"
  },
  UUID = ReservedUIDs[2],
  Events = {
    TRIGGER_UNLOCK = function(owner)
      ---@diagnostic disable-next-line: undefined-global
      AppendDerivedStat(owner, {
        Base = StatsConst.STAT_MANA_PLUS,
        Derived = StatsConst.STAT_ADAPTIVE,
        Percent = 100,
        Source = ReservedUIDs[3]
      })
    end
  },
} }

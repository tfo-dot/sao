ReservedUIDs = {
  "00000000-0000-0000-0000-000000000017",
  "00000000-0000-0001-0000-000000000017",
}

-- Meta
UUID = ReservedUIDs[1]
Name = "Mgliste wzmocenienie"
Description = "Otrzymujesz AP w zależności od siły leczenia i tarcz."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  HEAL_POWER = 15,
  AP = 30,
}

-- Effects
Effects = { {
  GetTrigger = {
    Type = "PASSIVE",
    Event = "NONE"
  },
  UUID = ReservedUIDs[2],
  Events = {
    TRIGGER_UNLOCK = function(owner)
      ---@diagnostic disable-next-line: undefined-global
      AppendDerivedStat(owner, {
        Base = StatsConst.STAT_HEAL_POWER,
        Derived = StatsConst.STAT_AP,
        Percent = 1000,
        Source = ReservedUIDs[2]
      })
    end
  },
} }

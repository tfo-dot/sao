ReservedUIDs = {
  "00000000-0000-0000-0000-00000000000E",
  "00000000-0000-0001-0000-00000000000E",
}

-- Meta
UUID = ReservedUIDs[1]
Name = "Płaszcz wzmacniający"
Description = "Zwiększa maksymalne zdrowie o 20%."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  HP = 300,
  DEF = 10,
  RES = 10,
}

Consts = {
  STATS_HP = 1
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
        Base = Consts.STATS_HP,
        Derived = Consts.STATS_HP,
        Percent = 20,
        Source = ReservedUIDs[2],
      })
    end,
  },
} }

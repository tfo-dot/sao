ReservedUIDs = {
  "00000000-0000-0000-0000-000000000006",
  "00000000-0000-0001-0000-000000000006",
  "00000000-0000-0001-0001-000000000006",
}

-- Meta
UUID = ReservedUIDs[1]
Name = "Oblicze obrony"
Description = "Dostajesz ATK w zależności od maks. HP."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  HP = 100,
  DEF = 15,
  MR = 15,
}

-- Effects
Effects = { {
  GetTrigger = {
    Type = "PASSIVE",
    Event = "NONE",
  },
  UUID = ReservedUIDs[2],
  GetEvents = function()
    return {
      TRIGGER_UNLOCK = function(owner)
        ---@diagnostic disable-next-line: undefined-global
        AppendDerivedStat(owner, {
          Base = StatsConst.STAT_HP,
          Derived = StatsConst.STAT_AD,
          Percent = 5,
          Source = ReservedUIDs[3],
        })
      end
    }
  end,
} }

ReservedUIDs = {
  "00000000-0000-0000-0000-00000000000C",
  "00000000-0000-0001-0000-00000000000C",
}

-- Meta
UUID = ReservedUIDs[1]
Name = "Ostrze obrońcy"
Description = "Zwiększa ataki o twój RES i DEF."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  HP = 150,
  DEF = 30,
  RES = 30,
  ATK = 20,
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
    local defStat = GetStat(owner, "STAT_DEF")
    ---@diagnostic disable-next-line: undefined-global
    local mrStat = GetStat(owner, "STAT_MR")

    return {
      Effects = {
        {
          Value = utils.PercentOf(defStat, 2) + utils.PercentOf(mrStat, 3),
          Type = 0,
          Percent = false
        },
      }
    }
  end,
} }

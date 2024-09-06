ReservedUIDs = {
  "00000000-0000-0000-0000-00000000001A",
  "00000000-0000-0001-0000-00000000001A",
}

-- Meta
UUID = ReservedUIDs[1]
Name = "PŁomień Shiki"
Description = "Obrażenia magiczne są zwiększone w zależności od zdrowia wroga"
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
    Event = "DAMAGE_BEFORE"
  },
  UUID = ReservedUIDs[2],
  Execute = function(owner, target, fightInstance, meta)
    ---@diagnostic disable-next-line: undefined-global
    local targetPercent = utils.percentOf(GetStat(target, StatsConst.STATS_HP), 5)

    return {
      Effects = {
        {
          ---@diagnostic disable-next-line: undefined-global
          Value = targetPercent + utils.percentOf(GetStat(owner, StatsConst.STAT_AP), 10),
          Type = 1,
          Percent = false,
        },
      },
    }
  end,
} }

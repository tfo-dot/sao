ReservedUIDs = {
  "00000000-0000-0000-0000-000000000002",
  "00000000-0000-0001-0000-000000000002",
}

-- Meta
UUID = ReservedUIDs[1]
Name = "Zabójca gigantów"
Description = "Zadaje dodatkowe obrażenia w zależności od pancerza przeciwnika."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  ATK = 25,
  LETHAL = 10,
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
    local damageValue = utils.PercentOf(GetStat(target, StatsConst.STAT_HP), 2)

    return {
      Effects = {
        { Value = damageValue, Type = 0, Percent = false }
      },
    }
  end,
} }

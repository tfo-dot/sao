ReservedUIDs = {
  "00000000-0000-0000-0000-000000000018",
  "00000000-0000-0001-0000-000000000018",
}

-- Meta
UUID = ReservedUIDs[1]
Name = "Wietrzne wzmocenienie"
Description = "Otrzymujesz SPD w zależności od siły leczenia i tarcz. Oraz leczysz przy ataku"
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  HEAL_POWER = 10,
  AD = 15,
}

-- Effects
Effects = { {
  Trigger = {
    Type = "PASSIVE",
    Event = "ATTACK_HIT"
  },
  UUID = ReservedUIDs[2],
  Execute = function(owner, target, fightInstance, meta)
    ---@diagnostic disable-next-line: undefined-global
    local validTargets = GetAlliesFor(fightInstance, GetUUID(owner))

    if #validTargets == 0 then
      return nil
    end

    local healTarget

    for idx = 1, #validTargets do
      if healTarget == nil then
        healTarget = validTargets[idx]
      end

      ---@diagnostic disable-next-line: undefined-global
      local healTargetPercent = GetCurrentHP(healTarget) / GetStat(healTarget, StatsConst.STAT_HP)
      ---@diagnostic disable-next-line: undefined-global
      local entityPercent = GetCurrentHP(validTargets[idx]) / GetStat(validTargets[idx], StatsConst.STAT_HP)

      if entityPercent < healTargetPercent then
        healTarget = validTargets[idx]
      end
    end

    ---@diagnostic disable-next-line: undefined-global
    local healValue = utils.PercentOf(GetStat(owner, StatsConst.STAT_AD), 10)

    ---@diagnostic disable-next-line: undefined-global
    HandleAction(fightInstance, {
      Event = "ACTION_EFFECT",
      ---@diagnostic disable-next-line: undefined-global
      Source = GetUUID(owner),
      ---@diagnostic disable-next-line: undefined-global
      Target = GetUUID(healTarget),
      Meta = {
        Effect = "EFFECT_HEAL",
        Value = 0,
        Duration = 1,
        Uuid = ReservedUIDs[2],
        Meta = {
          Value = healValue,
        },
        ---@diagnostic disable-next-line: undefined-global
        Caster = GetUUID(owner),
        ---@diagnostic disable-next-line: undefined-global
        Target = GetUUID(healTarget),
        Source = "SOURCE_ITEM",
      },
    })

    return nil
  end,
  Events = {
    TRIGGER_UNLOCK = function(owner)
      ---@diagnostic disable-next-line: undefined-global
      AppendDerivedStat(owner, {
        Base = StatsConst.STAT_HEAL_POWER,
        Derived = StatsConst.STAT_SPD,
        Percent = 100,
        Source = ReservedUIDs[2]
      })
    end
  },
} }

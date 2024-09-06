ReservedUIDs = {
  "00000000-0000-0000-0000-00000000000D",
  "00000000-0000-0001-0000-00000000000D",
}

-- Meta
UUID = ReservedUIDs[1]
Name = "Pancerz zwady"
Description = "Zadaje obrażenia wrogom, którzy cię uderzają i zmniejsza ich leczenie."
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
}

-- Effects
Effects = { {
  Trigger = {
    Type = "PASSIVE",
    Event = "ATTACK_GOT_HIT",
  },
  UUID = ReservedUIDs[2],
  Execute = function(owner, target, fightInstance, meta)
    ---@diagnostic disable-next-line: undefined-global
    HandleAction(fightInstance, {
      Event = "ACTION_DMG",
      Source = GetUUID(owner),
      Target = GetUUID(target),
      Meta = {
        Damage = {
          {
            ---@diagnostic disable-next-line: undefined-global
            Value = utils.PercentOf(GetStat(owner, StatsConst.STAT_DEF), 10),
            Type = 2,
            CanDodge = false,
          },
        },
      },
    })

    fightInstance:HandleAction({
      Event = "ACTION_EFFECT",
      Source = GetUUID(owner),
      Target = GetUUID(target),
      Meta = {
        Effect = "EFFECT_STAT_DEC",
        Value = -20,
        Duration = 1,
        Uuid = utils.GenerateUUID(),
        Meta = { Stat = StatsConst.STAT_HEAL_POWER, Value = -20, IsPercent = false },
        Caster = GetUUID(owner),
        Source = "SOURCE_ITEM",
      },
    })
    return nil
  end,
} }

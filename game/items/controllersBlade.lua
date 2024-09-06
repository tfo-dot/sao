ReservedUIDs = {
  "00000000-0000-0000-0000-000000000013",
  "00000000-0000-0001-0000-000000000013",
  "00000000-0000-0001-0001-000000000013",
}

-- Meta
UUID = ReservedUIDs[1]
Name = "Ostrze kontrolera"
Description = "Atakowanie zmniejsza prędkość wrogów."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  ATK = 15,
  SPD = 10
}

-- Effects
Effects = { {
  Trigger = {
    Type = "PASSIVE",
    Event = "ATTACK_HIT",
    Cooldown = {
      PassEvent = "ATTACK_HIT"
    }
  },
  UUID = ReservedUIDs[2],
  Execute = function(owner, target, fightInstance, meta)
    HandleAction(fightInstance, {
      Event = "ACTION_EFFECT",
      Source = GetUUID(owner),
      Target = GetUUID(target),
      Meta = {
        Effect = "EFFECT_STAT_DEC",
        Value = 0,
        Duration = 1,
        Uuid = ReservedUIDs[3],
        Meta = {
          Stat = StatsConst.STAT_SPD,
          Value = 10,
          IsPercent = false,
        },
        Caster = GetUUID(owner),
        Source = "SOURCE_ITEM",
      },
    })

    return nil
  end,
  CD = 3,
} }

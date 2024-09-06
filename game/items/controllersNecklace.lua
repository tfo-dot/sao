ReservedUIDs = {
  "00000000-0000-0000-0000-000000000012",
  "00000000-0000-0001-0000-000000000012",
  "00000000-0000-0001-0001-000000000012",
}

-- Meta
UUID = ReservedUIDs[1]
Name = "Naszyjnik kontrolera"
Description = "Nałożenie efektu CC zwiększa twoją prędkość."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  AP = 20,
  ATK = 10,
  SPD = 5
}

-- Effects
Effects = { {
  Trigger = {
    Type = "PASSIVE",
    Event = "APPLY_CROWD_CONTROL",
  },
  UUID = ReservedUIDs[2],
  Execute = function(owner, target, fightInstance, meta)
    ---@diagnostic disable-next-line: undefined-global
    HandleAction(fightInstance, {
      Event = "ACTION_EFFECT",
      Source = GetUUID(owner),
      Target = GetUUID(owner),
      Meta = {
        Effect = "EFFECT_STAT_INC",
        Value = 0,
        Duration = 1,
        Uuid = ReservedUIDs[3],
        Meta = {
          Stat = "STAT_SPD",
          Value = 10,
          IsPercent = false,
        },
        Caster = GetUUID(owner),
        Source = "SOURCE_ITEM",
      },
    })

    return nil
  end,
} }

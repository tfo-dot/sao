ReservedUIDs = {
  "00000000-0000-0000-0000-000000000004",
  "00000000-0000-0001-0000-000000000004",
}

-- Meta
UUID = ReservedUIDs[1]
Name = "Piaskowe ostrze"
Description = "Zadawanie obrażeń zmniejsza leczenie wroga."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  ATK = 30,
  SPD = 5,
}

-- Effects
Effects = { {
  Trigger = {
    Type = "PASSIVE",
    Event = "DAMAGE",
  },
  UUID = ReservedUIDs[2],
  Execute = function(owner, target, fightInstance, meta)
    ---@diagnostic disable-next-line: undefined-global
    HandleAction(fightInstance, {
      Event = "ACTION_EFFECT",
      Source = GetUUID(owner),
      Target = GetUUID(target),
      Meta = {
        Effect = "EFFECT_STAT_DEC",
        Value = -20,
        Duration = 1,
        Uuid = utils.GenerateUUID(),
        Meta = {
          Stat = StatsConst.STAT_HEAL_POWER,
          Value = -20,
          IsPercent = false,
        },
        Caster = GetUUID(owner),
        Source = "SOURCE_ITEM",
      },
    })

    return nil
  end,
} }

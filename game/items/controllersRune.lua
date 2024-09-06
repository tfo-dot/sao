ReservedUIDs = {
  "00000000-0000-0000-0000-000000000011",
  "00000000-0000-0001-0000-000000000011",
  "00000000-0000-0001-0001-000000000011",
}

-- Meta
UUID = ReservedUIDs[1]
Name = "Runa kontrolera"
Description = "Zabicie wroga objętego CC przywraca manę."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  AP = 20,
  ATK = 20,
}

-- Effects
Effects = { {
  Trigger = {
      Type = "PASSIVE",
      Event = "EXECUTE",
    },
  UUID = ReservedUIDs[2],
  Execute = function(owner, target, fightInstance, meta)
---@diagnostic disable-next-line: undefined-global
    HandleAction(fightInstance,{
      Event = "ACTION_EFFECT",
      Source = GetUUID(owner),
      Target = GetUUID(owner),
      Meta = {
        Effect = "EFFECT_MANA_RESTORE",
        Value = 1,
        Duration = 0,
        Uuid = ReservedUIDs[3],
        Meta = nil,
        Caster = GetUUID(owner),
        Source = "SOURCE_ITEM",
      },
    })

    return nil
  end,
} }

ReservedUIDs = {
  "00000000-0000-0000-0000-000000000101",
  "00000000-0000-0001-0000-000000000101",
}

UUID = ReservedUIDs[1]
Name = "Mała mikstura"
Description = "Leczy 25 punktów życia"

TakesSlot = false
Stacks = true
Consume = true
Count = 1
MaxCount = 5
Hidden = false

Stats = {}

Effects = { {
  Trigger = {
    Type = "ACTIVE",
    Event = "NONE"
  },
  UUID = ReservedUIDs[2],
  Execute = function(owner, target, fightInstance, meta)
    local healValue = 25

    ---@diagnostic disable-next-line: undefined-global
    HandleAction(fightInstance, {
      Event = "ACTION_EFFECT",
      ---@diagnostic disable-next-line: undefined-global
      Source = GetUUID(owner),
      ---@diagnostic disable-next-line: undefined-global
      Target = GetUUID(owner),
      Meta = {
        Effect = "EFFECT_HEAL",
        Value = 0,
        Duration = 0,
        ---@diagnostic disable-next-line: undefined-global
        Target = GetUUID(owner),
        ---@diagnostic disable-next-line: undefined-global
        Caster = GetUUID(owner),
        Source = "SOURCE_ITEM",
        Meta = {
          ---@diagnostic disable-next-line: undefined-global
          Value = healValue
        }
      },
    })

    return true
  end
} }

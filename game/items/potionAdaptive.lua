ReservedUIDs = {
  "00000000-0000-0000-0000-000000000103",
  "00000000-0000-0001-0000-000000000103",
}

UUID = ReservedUIDs[1]
Name = "Adaptacyjna mikstura"
Description = "Leczy 50+20% max HP"

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
    local healValue = 50 + utils.PercentOf(GetStat(owner, "STAT_HP"), 20)

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
        Uuid = utils.GenerateUUID(),
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

    return nil
  end
} }

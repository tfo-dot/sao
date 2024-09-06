ReservedUIDs = {
  "00000000-0000-0000-0000-000000000016",
  "00000000-0000-0001-0000-000000000016",
}

-- Meta
UUID = ReservedUIDs[1]
Name = "Syreni śpiew"
Description = "Leczenie i tarcze przeskakują na sojusznika"
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  HEAL_POWER = 10,
  AP = 40,
  HP = 50
}

-- Effects
Effects = { {
  Trigger = {
    Type = "PASSIVE",
    Event = "HEAL_OTHER"
  },
  UUID = ReservedUIDs[2],
  Execute = function(owner, target, fightInstance, meta)
    ---@diagnostic disable-next-line: undefined-global
    local validTargets = GetAlliesFor(fightInstance, GetUUID(owner))

    if #validTargets < 1 then
      return nil
    end

    local idx = -1

    for index = 1, #validTargets do
      ---@diagnostic disable-next-line: undefined-global
      if GetUUID(validTargets[idx]) == GetUUID(target) then
        idx = index
        break
      end
    end

    if idx ~= -1 then
      table.remove(validTargets, idx)
    end

    if #validTargets < 1 then
      return nil
    end

    local healValue = utils.PercentOf(meta.Value, 10)
    local healTarget = validTargets[math.random(#validTargets)]

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
        Duration = 0,
        ---@diagnostic disable-next-line: undefined-global
        Target = GetUUID(healTarget),
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
  end,
} }

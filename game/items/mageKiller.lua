ReservedUIDs = {
  "00000000-0000-0000-0000-000000000003",
  "00000000-0000-0001-0000-000000000003",
}

-- Meta
UUID = ReservedUIDs[1]
Name = "Zabójca magów"
Description = "Atakowanie celi osłoniętych tarczą zwiększa obrażenia twojego ataku."
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
    if GetEffectByType(target, "EFFECT_SHIELD") == nil then
      return nil
    end

    return {
      Effects = {
        { Value = 10, Type = 0, Percent = true }
      },
    }
  end,
} }

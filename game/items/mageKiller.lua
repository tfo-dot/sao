--ItemID
ReservedUIDs[0] = "00000000-0000-0000-0000-000000000003"
--SkillID
ReservedUIDs[2] = "00000000-0000-0001-0000-000000000003"

-- Meta
UUID = ReservedUIDs[0]
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
Effects[0] = {
  GetName = function() return "Zabójca magów" end,
  GetDescription = function() return "Atakowanie celi osłoniętych tarczą zwiększa obrażenia twojego ataku." end,
  GetTrigger = function()
    return {
      Type = "PASSIVE",
      Event = "ATTACK_BEFORE",
    }
  end,
  GetUUID = function() return ReservedUIDs[1] end,
  Execute = function(owner, target, fightInstance, meta)
    if target.GetEffectByType("EFFECT_SHIELD") == nil then
      return { Effect = {} }
    end

    return {
      Effects = {
        { Value = 10, Type = "DMG_PHYSICAL", Percent = true }
      },
    }
  end,
  GetEvents = function() return nil end,
  GetCD = function() return 0 end,
  GetCost = function() return 0 end
}

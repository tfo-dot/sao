--ItemID
ReservedUIDs[0] = "00000000-0000-0000-0000-000000000011"
--SkillID
ReservedUIDs[2] = "00000000-0000-0001-0000-000000000011"
--EffectID
ReservedUIDs[2] = "00000000-0000-0001-0001-000000000011"

-- Meta
UUID = ReservedUIDs[0]
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
Effects[0] = {
  GetName = function() return "Runa kontrolera" end,
  GetDescription = function() return "Zabicie wroga objętego CC przywraca manę." end,
  GetTrigger = function()
    return {
      Type = "PASSIVE",
      Event = "EXECUTE",
    }
  end,
  GetUUID = function() return ReservedUIDs[1] end,
  Execute = function(owner, target, fightInstance, meta)
    fightInstance:HandleAction({
      Event = "ACTION_EFFECT",
      Source = owner:GetUUID(),
      Target = owner:GetUUID(),
      Meta = {
        Effect = "EFFECT_MANA_RESTORE",
        Value = 1,
        Duration = 0,
        Uuid = ReservedUIDs[2],
        Meta = nil,
        Caster = owner:GetUUID(),
        Source = "SOURCE_ITEM",
      },
    })


    return nil
  end,
  GetEvents = function() return nil end,
  GetCD = function() return 0 end,
  GetCost = function() return 0 end
}

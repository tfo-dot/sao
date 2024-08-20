--ItemID
ReservedUIDs[0] = "00000000-0000-0000-0000-000000000004"
--SkillID
ReservedUIDs[2] = "00000000-0000-0001-0000-000000000004"

-- Meta
UUID = ReservedUIDs[0]
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
Effects[0] = {
  GetName = function() return "Piaskowe ostrze" end,
  GetDescription = function() return "Zadawanie obrażeń zmniejsza leczenie wroga." end,
  GetTrigger = function()
    return {
      Type = "PASSIVE",
      Event = "DAMAGE",
    }
  end,
  GetUUID = function() return ReservedUIDs[1] end,
  Execute = function(owner, target, fightInstance, meta)
    fightInstance.HandleAction({
      Event = "ACTION_EFFECT",
      Source = owner.GetUUID(),
      Target = target.GetUUID(),
      Meta = {
        Effect = "EFFECT_STAT_DEC",
        Value = -20,
        Duration = 1,
        Uuid = "New uuid",    --@TODO uuid.New(),
        Meta = {
          Stat = "STAT_HEAL_POWER",
          Value = -20,
          IsPercent = false,
        },
        Caster = owner.GetUUID(),
        Source = "SOURCE_ITEM",
      },
    })

    return nil
  end,
  GetEvents = function() return nil end,
  GetCD = function() return 0 end,
  GetCost = function() return 0 end
}

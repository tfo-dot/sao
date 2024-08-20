--ItemID
ReservedUIDs[0] = "00000000-0000-0000-0000-000000000012"
--SkillID
ReservedUIDs[2] = "00000000-0000-0001-0000-000000000012"
--EffectID
ReservedUIDs[2] = "00000000-0000-0001-0001-000000000012"

-- Meta
UUID = ReservedUIDs[0]
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
Effects[0] = {
  GetName = function() return "Naszyjnik kontrolera" end,
  GetDescription = function() return "Nałożenie efektu CC zwiększa twoją prędkość." end,
  GetTrigger = function()
    return {
      Type = "PASSIVE",
      Event = "APPLY_CROWD_CONTROL",
    }
  end,
  GetUUID = function() return ReservedUIDs[1] end,
  Execute = function(owner, target, fightInstance, meta)
    fightInstance:HandleAction({
      Event = "ACTION_EFFECT",
      Source = owner:GetUUID(),
      Target = owner:GetUUID(),
      Meta = {
        Effect = "EFFECT_STAT_INC",
        Value = 0,
        Duration = 1,
        Uuid = ReservedUIDs[2],
        Meta = {
          Stat = "STAT_SPD",
          Value = 10,
          IsPercent = false,
        },
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

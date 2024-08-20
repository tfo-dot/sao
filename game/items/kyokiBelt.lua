--ItemID
ReservedUIDs[0] = "00000000-0000-0000-0000-000000000019"
--SkillID
ReservedUIDs[2] = "00000000-0000-0001-0000-000000000019"

-- Meta
UUID = ReservedUIDs[0]
Name = "Pasek Kyoki"
Description = "Obrażenia magiczne są zwiększone przez losowy mnożnik (0.8-1.8)."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  ATK = 20,
  AP = 40
}

-- Effects
Effects[0] = {
  GetName = function() return "Pasek Kyoki" end,
  GetDescription = function() return "Obrażenia magiczne są zwiększone przez losowy mnożnik (0.8-1.8)." end,
  GetTrigger = function()
    return {
      Type = "PASSIVE",
      Event = "DAMAGE_BEFORE"
    }
  end,
  GetUUID = function() return ReservedUIDs[1] end,
  Execute = function(owner, target, fightInstance, meta)
    return {
      Effects = {
        {
          Value = utils.RandomNumber(0, 100) - 20,
          Type = "DMG_MAGICAL",
          Percent = true,
        },
      },
    }
  end,
  GetEvents = function() return nil end,
  GetCD = function() return 0 end,
  GetCost = function() return 0 end
}

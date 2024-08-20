--ItemID
ReservedUIDs[0] = "00000000-0000-0000-0000-00000000001B"
--SkillID
ReservedUIDs[2] = "00000000-0000-0001-0000-00000000001B"

-- Meta
UUID = ReservedUIDs[0]
Name = "Zwiastun burzy"
Description = "Ataki zadają dodatkowe obrażenia w zależności od AP."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  ATK = 100,
  SPD = 5
}

-- Effects
Effects[0] = {
  GetName = function() return "Zwiastun burzy" end,
  GetDescription = function() return "Ataki zadają dodatkowe obrażenia w zależności od AP." end,
  GetTrigger = function()
    return {
      Type = "PASSIVE",
      Event = "ATTACK_BEFORE"
    }
  end,
  GetUUID = function() return ReservedUIDs[1] end,
  Execute = function(owner, target, fightInstance, meta)
    return {
      Effects = {
        {
          Value = utils.PercentOf(owner:GetStat("STAT_AP"), 20),
          Type = "DMG_MAGICAL",
          Percent = false,
        },
      },
    }
  end,
  GetEvents = function() return nil end,
  GetCD = function() return 0 end,
  GetCost = function() return 0 end
}

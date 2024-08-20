--ItemID
ReservedUIDs[0] = "00000000-0000-0000-0000-00000000001A"
--SkillID
ReservedUIDs[2] = "00000000-0000-0001-0000-00000000001A"

-- Meta
UUID = ReservedUIDs[0]
Name = "PŁomień Shiki"
Description = "Obrażenia magiczne są zwiększone w zależności od zdrowia wroga"
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
  GetName = function() return "Pasek Kyoki" end,
  GetDescription = function() return "Obrażenia magiczne są zwiększone w zależności od zdrowia wroga" end,
  GetTrigger = function()
    return {
      Type = "PASSIVE",
      Event = "DAMAGE_BEFORE"
    }
  end,
  GetUUID = function() return ReservedUIDs[1] end,
  Execute = function(owner, target, fightInstance, meta)
    local targetPercent = utils.percentOf(target.GetStat("STAT_HP"), 5)

    return {
      Effects = {
        {
          Value = targetPercent + utils.percentOf(owner.GetStat("STAT_AP"), 10),
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

--ItemID
ReservedUIDs[0] = "00000000-0000-0000-0000-000000000001"
--SkillID
ReservedUIDs[2] = "00000000-0000-0001-0000-000000000001"

-- Meta
UUID = ReservedUIDs[0]
Name = "Pogromca gigantów"
Description = "Zadaje dodatkowe obrażenia w zależności od pancerza przeciwnika."
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
  GetName = function() return "Pogromca gigantów" end,
  GetDescription = function() return "Zadaje dodatkowe obrażenia w zależności od pancerza przeciwnika." end,
  GetTrigger = function()
    return {
      Type = "PASSIVE",
      Event = "ATTACK_BEFORE",
    }
  end,
  GetUUID = function() return ReservedUIDs[1] end,
  Execute = function(owner, target, fightInstance, meta)
    local damageValue = utils.PercentOf(target.GetStat("DEF"), 10)

    return {
      Effects = {
        { Value = damageValue, Type = "DMG_PHYSICAL" }
      },
    }
  end,
  GetEvents = function() return nil end,
  GetCD = function() return 0 end,
  GetCost = function() return 0 end
}

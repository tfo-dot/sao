--ItemID
ReservedUIDs[0] = "00000000-0000-0000-0000-000000000017"
--SkillID
ReservedUIDs[2] = "00000000-0000-0001-0000-000000000017"

-- Meta
UUID = ReservedUIDs[0]
Name = "Mgliste wzmocenienie"
Description = "Otrzymujesz AP w zależności od siły leczenia i tarcz."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  HEAL_POWER = 15,
  AP = 30,
}

-- Effects
Effects[0] = {
  GetName = function() return "Mgliste wzmocenienie" end,
  GetDescription = function() return "Otrzymujesz AP w zależności od siły leczenia i tarcz." end,
  GetTrigger = function()
    return {
      Type = "PASSIVE",
      Event = "NONE"
    }
  end,
  GetUUID = function() return ReservedUIDs[1] end,
  Execute = function(owner, target, fightInstance, meta) return nil end,
  GetEvents = function() return {
    TRIGGER_UNLOCK = function(owner)
      owner:AppendDerivedStat({
        Base = "STAT_HEAL_POWER",
        Derived = "STAT_AP",
        Percent = 1000,
        Source = ReservedUIDs[2]
      })
    end
  } end,
  GetCD = function() return 0 end,
  GetCost = function() return 0 end
}

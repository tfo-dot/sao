--ItemID
ReservedUIDs[0] = "00000000-0000-0000-0000-00000000000E"
--SkillID
ReservedUIDs[2] = "00000000-0000-0001-0000-00000000000E"

-- Meta
UUID = ReservedUIDs[0]
Name = "Płaszcz wzmacniający"
Description = "Zwiększa maksymalne zdrowie o 20%."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  HP = 300,
  DEF = 10,
  RES = 10,
}

-- Effects
Effects[0] = {
  GetName = function() return "Płaszcz wzmacniający" end,
  GetDescription = function() return "Zwiększa maksymalne zdrowie o 20%." end,
  GetTrigger = function()
    return {
      Type = "PASSIVE",
      Event = "NONE",
    }
  end,
  GetUUID = function() return ReservedUIDs[1] end,
  Execute = function(owner, target, fightInstance, meta) return nil end,
  GetEvents = function() return {
    TRIGGER_UNLOCK = function(owner)
      owner:AppendDerivedStat({
        Base = "STAT_HP",
        Derived = "STAT_HP",
        Percent = 20,
        Source = ReservedUIDs[1],
      })
    end
  } end,
  GetCD = function() return 0 end,
  GetCost = function() return 0 end
}

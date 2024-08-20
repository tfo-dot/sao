--ItemID
ReservedUIDs[0] = "00000000-0000-0000-0000-00000000000B"
--SkillID
ReservedUIDs[2] = "00000000-0000-0001-0000-00000000000B"

-- Meta
UUID = ReservedUIDs[0]
Name = "Dziedzictwo Ryu"
Description = "Zwiększa RES i DEF o 20%."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  HP = 150,
  DEF = 40,
  RES = 40,
}

-- Effects
Effects[0] = {
  GetName = function() return "Dziedzictwo Ryu" end,
  GetDescription = function() return "Zwiększa RES i DEF o 20%." end,
  GetTrigger = function()
    return {
      Type = "PASSIVE",
      Event = "NONE",
    }
  end,
  GetUUID = function() return ReservedUIDs[1] end,
  Execute = function(owner, target, fightInstance, meta) return nil end,
  GetEvents = function()
    return {
      TRIGGER_UNLOCK = function(owner)
        owner:AppendDerivedStat({
          Base = "STAT_DEF",
          Derived = "STAT_DEF",
          Percent = 20,
          Source = ReservedUIDs[1],
        })

        owner:AppendDerivedStat({
          Base = "STAT_MR",
          Derived = "STAT_MR",
          Percent = 20,
          Source = ReservedUIDs[1],
        })
      end
    }
  end,
  GetCD = function() return 0 end,
  GetCost = function() return 0 end
}

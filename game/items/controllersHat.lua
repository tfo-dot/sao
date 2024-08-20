--ItemID
ReservedUIDs[0] = "00000000-0000-0000-0000-000000000014"
--SkillID
ReservedUIDs[2] = "00000000-0000-0001-0000-000000000014"
--EffectID
ReservedUIDs[2] = "00000000-0000-0001-0001-000000000014"

-- Meta
UUID = ReservedUIDs[0]
Name = "Kapelusz kontrolera"
Description = "Daje siłę adaptacyjną w zależności od many."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  MANA = 5
}

-- Effects
Effects[0] = {
  GetName = function() return "Kapelusz kontrolera" end,
  GetDescription = function() return "Daje siłę adaptacyjną w zależności od many." end,
  GetTrigger = function()
    return {
      Type = "PASSIVE",
      Event = "NONE"
    }
  end,
  GetUUID = function() return ReservedUIDs[1] end,
  Execute = function(owner, target, fightInstance, meta) return nil end,
  GetEvents = function()
    return {
      TRIGGER_UNLOCK = function(owner)
        owner:AppendDerivedStat({
          Base = "STAT_MANA_PLUS",
          Derived = "STAT_ADAPTIVE",
          Percent = 100,
          Source = ReservedUIDs[2]
        })
      end
    }
  end,
  GetCD = function() return 0 end,
  GetCost = function() return 0 end
}

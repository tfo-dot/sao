--ItemID
ReservedUIDs[0] = "00000000-0000-0000-0000-000000000006"
--SkillID
ReservedUIDs[2] = "00000000-0000-0001-0000-000000000006"
--EffectID
ReservedUIDs[2] = "00000000-0000-0001-0001-000000000006"

-- Meta
UUID = ReservedUIDs[0]
Name = "Oblicze obrony"
Description = "Dostajesz ATK w zależności od maks. HP."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  HP = 100,
  DEF = 15,
  MR = 15,
}

-- Effects
Effects[0] = {
  GetName = function() return "Oblicze obrony" end,
  GetDescription = function() return "Dostajesz ATK w zależności od maks. HP." end,
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
        owner.AppendDerivedStat({
          Base = "STAT_HP",
          Derived = "STAT_AD",
          Percent = 5,
          Source = ReservedUIDs[2],
        })
      end
    }
  end,
  GetCD = function() return 0 end,
  GetCost = function() return 0 end
}

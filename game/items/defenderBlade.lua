--ItemID
ReservedUIDs[0] = "00000000-0000-0000-0000-00000000000C"
--SkillID
ReservedUIDs[2] = "00000000-0000-0001-0000-00000000000C"

-- Meta
UUID = ReservedUIDs[0]
Name = "Ostrze obrońcy"
Description = "Zwiększa ataki o twój RES i DEF."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  HP = 150,
  DEF = 30,
  RES = 30,
  ATK = 20,
}

-- Effects
Effects[0] = {
  GetName = function() return "Ostrze obrońcy" end,
  GetDescription = function() return "Zwiększa ataki o twój RES i DEF." end,
  GetTrigger = function()
    return {
      Type = "PASSIVE",
      Event = "ATTACK_BEFORE",
    }
  end,
  GetUUID = function() return ReservedUIDs[1] end,
  Execute = function(owner, target, fightInstance, meta)
    local defStat = owner:GetStat("STAT_DEF")
    local mrStat = owner:GetStat("STAT_MR")

    return {
      Effects = {
        {
          Value = utils.PercentOf(defStat, 2) + utils.PercentOf(mrStat, 3),
          Type = "DMG_PHYSICAL",
        },
      }
    }
  end,
  GetEvents = function() return nil end,
  GetCD = function() return 0 end,
  GetCost = function() return 0 end
}

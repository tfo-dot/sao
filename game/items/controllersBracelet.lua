--ItemID
ReservedUIDs[0] = "00000000-0000-0000-0000-00000000000F"
--SkillID
ReservedUIDs[2] = "00000000-0000-0001-0000-00000000000F"

-- Meta
UUID = ReservedUIDs[0]
Name = "Bransoleta kontrolera"
Description = "Nałożenie efektu CC leczy ciebie i sojusznika."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  AP = 20,
  ATK = 10,
  SPD = 5,
}

-- Effects
Effects[0] = {
  GetName = function() return "Bransoleta kontrolera" end,
  GetDescription = function() return "Nałożenie efektu CC leczy ciebie i sojusznika." end,
  GetTrigger = function()
    return {
      Type = "PASSIVE",
      Event = "APPLY_CROWD_CONTROL",
    }
  end,
  GetUUID = function() return ReservedUIDs[1] end,
  Execute = function(owner, target, fightInstance, meta)
    local healValue = utils.PercentOf(owner:GetStat("STAT_AD"), 15) + utils.PercentOf(owner:GetStat("STAT_AP"), 15)
    --@TODO Choose heal target (now it should heal ccd target)

    fightInstance:HandleAction({
      Event = "ACTION_EFFECT",
      Source = owner:GetUUID(),
      Target = target:GetUUID(),
      Meta = { Value = healValue },
    })

    fightInstance:HandleAction({
      Event = "ACTION_EFFECT",
      Source = owner:GetUUID(),
      Target = owner:GetUUID(),
      Meta = { Value = healValue },
    })

    return nil
  end,
  GetEvents = function() return nil end,
  GetCD = function() return 0 end,
  GetCost = function() return 0 end
}

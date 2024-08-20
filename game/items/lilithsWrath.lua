--ItemID
ReservedUIDs[0] = "00000000-0000-0000-0000-00000000000A"
--SkillID
ReservedUIDs[2] = "00000000-0000-0001-0000-00000000000A"

-- Meta
UUID = ReservedUIDs[0]
Name = "Gniew Lilith"
Description = "Co ture zadaje obrażenia w zależności od zdrowia użytkownika."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  HP = 200,
  DEF = 30,
}

-- Effects
Effects[0] = {
  GetName = function() return "Gniew Lilith" end,
  GetDescription = function() return "Co ture zadaje obrażenia w zależności od zdrowia użytkownika." end,
  GetTrigger = function()
    return {
      Type = "PASSIVE",
      Event = "TURN",
    }
  end,
  GetUUID = function() return ReservedUIDs[1] end,
  Execute = function(owner, target, fightInstance, meta)
    --@TODO For all enemies
    fightInstance.HandleAction({
      Event = "ACTION_DMG",
      Source = owner:GetUUID(),
      Target = target:GetUUID(),
      Meta = { {
        Value = utils.PercentOf(owner:GetStat("STAT_HP"), 5),
        Type = "DMG_PHYSICAL",
        CanDodge = false,
      } },
    })

    return nil
  end,
  GetEvents = function() return nil end,
  GetCD = function() return 0 end,
  GetCost = function() return 0 end
}

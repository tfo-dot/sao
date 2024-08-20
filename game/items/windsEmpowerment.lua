--ItemID
ReservedUIDs[0] = "00000000-0000-0000-0000-000000000018"
--SkillID
ReservedUIDs[2] = "00000000-0000-0001-0000-000000000018"

-- Meta
UUID = ReservedUIDs[0]
Name = "Wietrzne wzmocenienie"
Description = "Otrzymujesz SPD w zależności od siły leczenia i tarcz. Oraz leczysz przy ataku"
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  HEAL_POWER = 10,
  AD = 15,
}

-- Effects
Effects[0] = {
  GetName = function() return "Wietrzne wzmocenienie" end,
  GetDescription = function() return "Otrzymujesz SPD w zależności od siły leczenia i tarcz. Oraz leczysz przy ataku" end,
  GetTrigger = function()
    return {
      Type = "PASSIVE",
      Event = "ATTACK_HIT"
    }
  end,
  GetUUID = function() return ReservedUIDs[1] end,
  Execute = function(owner, target, fightInstance, meta)
    local validTargets = fightInstance:GetAlliesFor(owner:GetUUID())

    if validTargets.len == 0 then
      return nil
    end

    local healTarget

    for idx = 1, validTargets.len do
      if healTarget == nil then
        healTarget = validTargets[idx]
      end

      local healTargetPercent = healTarget.GetCurrentHP() / healTarget.GetStat("STAT_HP")
      local entityPercent = validTargets[idx].GetCurrentHP() / validTargets[idx].GetStat("STAT_HP")

      if entityPercent < healTargetPercent then
        healTarget = validTargets[idx]
      end
    end

    local healValue = utils.PercentOf(owner:GetStat("STAT_AD"), 10)

    fightInstance:HandleAction({
      Event = "ACTION_EFFECT",
      Source = owner:GetUUID(),
      Target = healTarget.GetUUID(),
      Meta = { Value = healValue },
    })

    return nil
  end,
  GetEvents = function()
    return {
      TRIGGER_UNLOCK = function(owner)
        owner:AppendDerivedStat({
          Base = "STAT_HEAL_POWER",
          Derived = "STAT_SPD",
          Percent = 100,
          Source = ReservedUIDs[2]
        })
      end
    }
  end,
  GetCD = function() return 0 end,
  GetCost = function() return 0 end
}

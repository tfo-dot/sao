--ItemID
ReservedUIDs[0] = "00000000-0000-0000-0000-000000000008"
--SkillID
ReservedUIDs[2] = "00000000-0000-0001-0000-000000000008"
--EffectID
ReservedUIDs[2] = "00000000-0000-0001-0001-000000000008"

-- Meta
UUID = ReservedUIDs[0]
Name = "Dziedzictwo wojownika"
Description = "Zwiększa obrażenia w zależności od maks zdrowia."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  ATK = 20,
  HP = 50
}

-- Effects
Effects[0] = {
  GetName = function() return "Dziedzictwo wojownika" end,
  GetDescription = function() return "Zwiększa obrażenia w zależności od maks zdrowia." end,
  GetTrigger = function()
    return {
      Type = "PASSIVE",
      Event = "ATTACK_BEFORE",
    }
  end,
  GetUUID = function() return ReservedUIDs[1] end,
  Execute = function(owner, target, fightInstance, meta)
    local dmgPercent = utils.PercentOf(owner.GetStat("STAT_HP_PLUS"), 1)

    return {
      Effects = {
        {
          Value = dmgPercent,
          Type = "DMG_PHYSICAL",
          Percent = true,
        },
      },
    }
  end,
  GetEvents = function() return nil end,
  GetCD = function() return 0 end,
  GetCost = function() return 0 end
}

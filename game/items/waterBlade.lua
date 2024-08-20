--ItemID
ReservedUIDs[0] = "00000000-0000-0000-0000-000000000005"
--SkillID
ReservedUIDs[2] = "00000000-0000-0001-0000-000000000005"

-- Meta
UUID = ReservedUIDs[0]
Name = "Wodne ostrze"
Description = "Zadawanie obrażeń leczy o brakujące zdrowie."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  ATK = 25,
  VAMP = 10,
  HP = 50,
}

-- Effects
Effects[0] = {
  GetName = function() return "Wodne ostrze" end,
  GetDescription = function() return "Zadawanie obrażeń leczy o brakujące zdrowie." end,
  GetTrigger = function()
    return {
      Type = "PASSIVE",
      Event = "ATTACK_HIT",
    }
  end,
  GetUUID = function() return ReservedUIDs[1] end,
  Execute = function(owner, target, fightInstance, meta)
    local addPercentage = utils.PercentOf(owner:GetStat("AD"), 1)

    fightInstance.HandleAction({
      Event = "ACTION_EFFECT",
      Source = owner:GetUUID(),
      Target = owner:GetUUID(),
      Meta = {
        Effect = "EFFECT_HEAL",
        Value = utils.PercentOf(owner.GetStat("HP") - owner.GetCurrentHP(), 10 + addPercentage),
        Duration = 0,
        Caster = owner.GetUUID(),
        Source = "SOURCE_ITEM",
      },
    })

    return nil
  end,
  GetEvents = function() return nil end,
  GetCD = function() return 10 end,
  GetCost = function() return 0 end
}

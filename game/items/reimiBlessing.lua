--ItemID
ReservedUIDs[0] = "00000000-0000-0000-0000-000000000000"
--SkillID
ReservedUIDs[1] = "00000000-0000-0001-0000-000000000001"
--EffectID
ReservedUIDs[2] = "00000000-0000-0001-0001-000000000001"

-- Meta
UUID = ReservedUIDs[0]
Name = "Błogosławieństwo Reimi"
Description = "Przeleczenie daje tarczę."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  ATK = 25,
  HP = 100,
  ATK_VAMP = 10,
}

-- Effects
Effects[0] = {
  GetName = function() return "Błogosławieństwo Reimi" end,
  GetDescription = function() return "Przeleczenie daje tarczę." end,
  GetTrigger = function()
    return {
      Type = "PASSIVE",
      Event = "HEAL_SELF",
    }
  end,
  GetUUID = function() return ReservedUIDs[1] end,
  Execute = function(owner, target, fightInstance, meta)
    local oldEffect = owner:GetEffectByUUID(ReservedUIDs[2])
    local maxShield = utils.percentOf(owner:GetStat("HP"), 25) + utils.percentOf(owner:GetStat("ATK"), 25)

    if oldEffect ~= nil then
      owner.RemoveEffect(oldEffect)
    else
      oldEffect = {
        Effect = "EFFECT_SHIELD",
        Value = 0,
        Duration = -1,
        Uuid = ReservedUIDs[2],
        Meta = nil,
        Caster = owner.GetUUID(),
        Source = "SOURCE_ITEM",
      }
    end

    if oldEffect.Value < 0 then
      oldEffect.Value = 0
    end

    oldEffect.Value = oldEffect.Value + meta.Value

    if oldEffect.Value > maxShield then
      oldEffect.Value = oldEffect.Value + maxShield
    end

    owner:ApplyEffect(oldEffect)

    return nil
  end,
  GetEvents = function()
    return {
      TRIGGER_UNLOCK = function(owner, target, fightInstance, meta)
        owner:ApplyEffect({
          Effect = "Shield",
          Value = 0,
          Duration = -1,
          Uuid = ReservedUIDs[2],
          Meta = nil,
          Caster = owner:GetUUID(),
          Source = "SOURCE_ITEM"
        })
      end,
    }
  end,
  GetCD = function() return 0 end,
  GetCost = function() return 0 end
}

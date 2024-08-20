--ItemID
ReservedUIDs[0] = "00000000-0000-0000-0000-000000000013"
--SkillID
ReservedUIDs[2] = "00000000-0000-0001-0000-000000000013"
--EffectID
ReservedUIDs[2] = "00000000-0000-0001-0001-000000000013"

-- Meta
UUID = ReservedUIDs[0]
Name = "Ostrze kontrolera"
Description = "Atakowanie zmniejsza prędkość wrogów."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  ATK = 15,
  SPD = 10
}

-- Effects
Effects[0] = {
  GetName = function() return "Ostrze kontrolera" end,
  GetDescription = function() return "Atakowanie zmniejsza prędkość wrogów." end,
  GetTrigger = function()
    return {
      Type = "PASSIVE",
      Event = "ATTACK_HIT",
      Cooldown = {
        PassEvent = "ATTACK_HIT"
      }
    }
  end,
  GetUUID = function() return ReservedUIDs[1] end,
  Execute = function(owner, target, fightInstance, meta)
    fightInstance:HandleAction({
      Event = "ACTION_EFFECT",
      Source = owner:GetUUID(),
      Target = owner:GetUUID(),
      Meta = {
        Effect = "EFFECT_STAT_DEC",
        Value = 0,
        Duration = 1,
        Uuid = ReservedUIDs[2],
        Meta = {
          Stat = "STAT_SPD",
          Value = -10,
          IsPercent = false,
        },
        Caster = owner:GetUUID(),
        Source = "SOURCE_ITEM",
      },
    })

    return nil
  end,
  GetEvents = function() return nil end,
  GetCD = function() return 3 end,
  GetCost = function() return 0 end
}

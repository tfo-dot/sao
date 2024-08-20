--ItemID
ReservedUIDs[0] = "00000000-0000-0000-0000-00000000000D"
--SkillID
ReservedUIDs[2] = "00000000-0000-0001-0000-00000000000D"

-- Meta
UUID = ReservedUIDs[0]
Name = "Pancerz zwady"
Description = "Zadaje obrażenia wrogom, którzy cię uderzają i zmniejsza ich leczenie."
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
}

-- Effects
Effects[0] = {
  GetName = function() return "Pancerz zwady" end,
  GetDescription = function() return "Zadaje obrażenia wrogom, którzy cię uderzają i zmniejsza ich leczenie." end,
  GetTrigger = function()
    return {
      Type = "PASSIVE",
      Event = "ATTACK_GOT_HIT",
    }
  end,
  GetUUID = function() return ReservedUIDs[1] end,
  Execute = function(owner, target, fightInstance, meta)
    fightInstance:HandleAction({
      Event = "ACTION_DMG",
      Source = owner:GetUUID(),
      Target = target:GetUUID(),
      Meta = {
        Damage = {
          {
            Value = utils.PercentOf(owner:GetStat("STAT_DEF"), 10),
            Type = "DMG_TRUE",
            CanDodge = false,
          },
        },
      },
    })

    fightInstance:HandleAction({
      Event = "ACTION_EFFECT",
      Source = owner:GetUUID(),
      Target = target:GetUUID(),
      Meta = {
        Effect = "EFFECT_STAT_DEC",
        Value = -20,
        Duration = 1,
        Uuid = "New uuid",    --@TODO uuid.New(),
        Meta = { Stat = "STAT_HEAL_POWER", Value = -20, IsPercent = false },
        Caster = owner:GetUUID(),
        Source = "SOURCE_ITEM",
      },
    })
    return nil
  end,
  GetEvents = function() return nil end,
  GetCD = function() return 0 end,
  GetCost = function() return 0 end
}

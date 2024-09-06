ReservedUIDs = {
  "00000000-0000-0000-0000-000000000000",
  "00000000-0000-0001-0000-000000000001",
  "00000000-0000-0001-0001-000000000001",
}

-- Meta
UUID = ReservedUIDs[1]
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
Effects = { {
  Trigger = {
    Type = "PASSIVE",
    Event = "HEAL_SELF",
  },
  UUID = ReservedUIDs[2],
  Execute = function(owner, target, fightInstance, meta)
    ---@diagnostic disable-next-line: undefined-global
    local oldEffect = GetEffectByUUID(owner, ReservedUIDs[3])
    ---@diagnostic disable-next-line: undefined-global
    local maxShield = utils.percentOf(GetStat(owner, StatsConst.STAT_HP), 25) +
        ---@diagnostic disable-next-line: undefined-global
        utils.percentOf(GetStat(owner, StatsConst.STAT_AD), 25)

    if oldEffect ~= nil then
      owner.RemoveEffect(ReservedUIDs[3])
    else
      oldEffect = {
        Effect = "EFFECT_SHIELD",
        Value = 0,
        Duration = -1,
        Uuid = ReservedUIDs[3],
        Meta = nil,
        Caster = GetUUID(owner),
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
          Effect = "EFFECT_SHIELD",
          Value = 0,
          Duration = -1,
          Uuid = ReservedUIDs[3],
          Meta = nil,
          Caster = GetUUID(owner),
          Source = "SOURCE_ITEM"
        })
      end,
    }
  end,
  GetCD = function() return 0 end,
  GetCost = function() return 0 end
} }

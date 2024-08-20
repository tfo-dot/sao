--ItemID
ReservedUIDs[0] = "00000000-0000-0000-0000-000000000010"
--SkillID
ReservedUIDs[2] = "00000000-0000-0001-0000-000000000010"

-- Meta
UUID = ReservedUIDs[0]
Name = "Przeklęty lód"
Description = "Efekty spowolnienia są mocniejsze"
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  AP = 20,
  ATK = 20,
}

-- Effects
Effects[0] = {
  GetName = function() return "Przeklęty lód" end,
  GetDescription = function() return "Efekty spowolnienia są mocniejsze" end,
  GetTrigger = function()
    return {
      Type = "PASSIVE",
      Event = "APPLY_CROWD_CONTROL",
    }
  end,
  GetUUID = function() return ReservedUIDs[1] end,
  Execute = function(owner, target, fightInstance, meta)
    if meta.Effect == "EFFECT_STAT_DEC" then
      if meta.Meta.Stat == "STAT_SPD" then
        return {
          Effects = {
            {
              Value = 20,
              Percent = true,
            },
          }
        }
      end
    end

    return nil
  end,
  GetEvents = function() return nil end,
  GetCD = function() return 0 end,
  GetCost = function() return 0 end
}

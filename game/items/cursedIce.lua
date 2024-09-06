ReservedUIDs = {
  "00000000-0000-0000-0000-000000000010",
  "00000000-0000-0001-0000-000000000010",
}

-- Meta
UUID = ReservedUIDs[1]
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
Effects = { {
  Trigger = {
    Type = "PASSIVE",
    Event = "APPLY_CROWD_CONTROL",
  },
  UUID = ReservedUIDs[2],
  Execute = function(owner, target, fightInstance, meta)
    if meta.Effect == "EFFECT_STAT_DEC" then
      if meta.Meta.Stat == StatsConst.STAT_SPD then
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
} }

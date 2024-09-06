ReservedUIDs = {
  "00000000-0000-0000-0000-00000000000A",
  "00000000-0000-0001-0000-00000000000A",
}

-- Meta
UUID = ReservedUIDs[1]
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
Effects = { {
  Trigger = {
    Type = "PASSIVE",
    Event = "TURN",
  },
  UUID = ReservedUIDs[2],
  Execute = function(owner, target, fightInstance, meta)
    ---@diagnostic disable-next-line: undefined-global
    local enemies = GetEnemies(fightInstance, owner)

    for idx = 1, #enemies do
      local enemy = enemies[idx]

      ---@diagnostic disable-next-line: undefined-global
      HandleAction(fightInstance, {
        Event = "ACTION_DMG",
        Source = GetUUID(owner),
        Target = GetUUID(enemy),
        Meta = { {
          ---@diagnostic disable-next-line: undefined-global
          Value = utils.PercentOf(GetStat(owner, StatsConst.STAT_HP), 5),
          Type = 0,
          CanDodge = false,
        } },
      })
    end

    return nil
  end,
} }

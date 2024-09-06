--Base
Id = "LV1_Toxic_Spider"
HP = 170
SPD = 45
ATK = 50
Name = "Jadowity Pająk"

Const = {
  ITEM = 0,
  EXP = 1,
  GOLD = 2
}

--Loot
Loot = {
  { Type = Const.EXP,  Count = 90 },
  { Type = Const.GOLD, Count = 120 }
}

Action = function(mob, fight)
  local turn = GetTurnFor(fight, GetUUID(mob))
  if turn % 3 == 0 and turn ~= 1 then
    local enemies = GetEnemiesFor(fight, GetUUID(mob))
    local target = enemies[math.random(#enemies)]

    local entityActions = DefaultAction(mob, fight)

    table.insert(entityActions, { {
      Event = "ACTION_EFFECT",
      ---@diagnostic disable-next-line: undefined-global
      Source = GetUUID(owner),
      ---@diagnostic disable-next-line: undefined-global
      Target = GetUUID(target),
      Meta = {
        Effect = "EFFECT_DOT",
        Value = 20,
        Duration = 3,
        Uuid = utils.GenerateUUID(),
        ---@diagnostic disable-next-line: undefined-global
        Caster = GetUUID(owner),
        ---@diagnostic disable-next-line: undefined-global
        Target = GetUUID(target),
        Source = "SOURCE_ND",
      },
      DefaultAction(mob, fight)
    }
    })

    return entityActions
  else
    return DefaultAction(mob, fight)
  end

  return nil
end

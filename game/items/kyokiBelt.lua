ReservedUIDs = {
  "00000000-0000-0000-0000-000000000019",
  "00000000-0000-0001-0000-000000000019",
}

-- Meta
UUID = ReservedUIDs[1]
Name = "Pasek Kyoki"
Description = "Obrażenia magiczne są zwiększone przez losowy mnożnik (0.8-1.8)."
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  ATK = 20,
  AP = 40
}

-- Effects
Effects = { {
  Trigger = {
    Type = "PASSIVE",
    Event = "DAMAGE_BEFORE"
  },
  UUID = ReservedUIDs[2],
  Execute = function(owner, target, fightInstance, meta)
    return {
      Effects = {
        {
          Value = math.random(0, 100) - 20,
          Type = 1,
          Percent = true,
        },
      },
    }
  end,
} }

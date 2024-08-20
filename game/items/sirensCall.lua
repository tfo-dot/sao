--ItemID
ReservedUIDs[0] = "00000000-0000-0000-0000-000000000016"
--SkillID
ReservedUIDs[2] = "00000000-0000-0001-0000-000000000016"
--EffectID
ReservedUIDs[2] = "00000000-0000-0001-0001-000000000016"

-- Meta
UUID = ReservedUIDs[0]
Name = "Syreni śpiew"
Description = "Leczenie i tarcze przeskakują na sojusznika"
TakesSlot = true
Stacks = false
Consume = false
Count = 1
MaxCount = 1
Hidden = false

-- Stats
Stats = {
  HEAL_POWER = 10,
  AP = 40,
  HP = 50
}

-- Effects
Effects[0] = {
  GetName = function() return "Syreni śpiew" end,
  GetDescription = function() return "Leczenie i tarcze przeskakują na sojusznika" end,
  GetTrigger = function()
    return {
      Type = "PASSIVE",
      Event = "HEAL_OTHER"
    }
  end,
  GetUUID = function() return ReservedUIDs[1] end,
  Execute = function(owner, target, fightInstance, meta)
    local validTargets = fightInstance:GetAlliesFor(owner:GetUUID())

    if validTargets.len <= 1 then
      return nil
    end

    local idx = -1

    for index = 1, validTargets.len do
      if validTargets[index]:GetUUID() == target:GetUUID() then
        idx = index
        break
      end
    end

    if idx > -1 then
      -- Remove the target from the list of valid targets
      --@TODO SOMEHOW PORT IT
      validTargets = utils.append(validTargets[idx], validTargets[idx + 1])
    end

    local healValue = utils.PercentOf(meta.Value, 10)
    local healTarget = utils.RandomElement(validTargets)

    fightInstance:HandleAction({
      Event = "ACTION_EFFECT",
      Source = owner:GetUUID(),
      Target = healTarget:GetUUID(),
      Meta = {
        Value = healValue,
      },
    })

    return nil
  end,
  GetEvents = function() return nil end,
  GetCD = function() return 0 end,
  GetCost = function() return 0 end
}

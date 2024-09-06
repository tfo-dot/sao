--Base
Id = "LV0_Error"
HP = 250
SPD = 45
ATK = 50

Name = "Błąd"

Const = {
  ITEM = 0,
  EXP = 1,
  GOLD = 2
}

--Loot
Loot = {}

-- Effects

OnDefeat = function(player)
  --TODO change it when release
  ---@diagnostic disable-next-line: undefined-global
  UnlockFloor(player, "beta-piętro-2")
end
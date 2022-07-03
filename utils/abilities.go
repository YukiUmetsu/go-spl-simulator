package utils

import (
	sim "simulator"
)

/** Abilities that summoner applies to friendly team at the start of the game */
func GetSummonerAbilities() []sim.Ability {
	return []sim.Ability{
		sim.ABILITY_AFFLICTION,
		sim.ABILITY_AMPLIFY,
		sim.ABILITY_BLAST,
		sim.ABILITY_DIVINE_SHIELD,
		sim.ABILITY_FLYING,
		sim.ABILITY_LAST_STAND,
		sim.ABILITY_MAGIC_REFLECT,
		sim.ABILITY_PIERCING,
		sim.ABILITY_RETURN_FIRE,
		sim.ABILITY_SNARE,
		sim.ABILITY_THORNS,
		sim.ABILITY_TRUE_STRIKE,
		sim.ABILITY_VOID,
		sim.ABILITY_VOID_ARMOR,
		sim.POISON,
	}
}

/** Buffs that summoner applies to friendly team at the start of the game */
func GetSummonerPreGameBuffAbilities() []sim.Ability {
	return []sim.Ability{
		sim.ABILITY_STRENGTHEN,
	}
}

/** Abilities that summoner applies to enemy team at the start of the game */
func GetSummonerPreGameDebuffAbilities() []sim.Ability {
	return []sim.Ability{
		sim.ABILITY_AFFLICTION,
		sim.ABILITY_BLIND,
	}
}

/** Abilities that monsters apply to friendly team at the start of the game */
func GetMonsterPreGameBuffAbilities() []sim.Ability {
	return []sim.Ability{
		sim.ABILITY_PROTECT,
		sim.ABILITY_STRENGTHEN,
		sim.ABILITY_SWIFTNESS,
		sim.ABILITY_INSPIRE,
	}
}

/** Abilities that monsters apply to enemy team at the start of the game */
func GetMonsterPreGameDebuffAbilities() []sim.Ability {
	return []sim.Ability{
		sim.ABILITY_AMPLIFY,
		sim.ABILITY_BLIND,
		sim.ABILITY_DEMORALIZE,
		sim.ABILITY_HEADWINDS,
		sim.ABILITY_RUST,
		sim.ABILITY_SLOW,
		sim.ABILITY_SNARE,
		sim.ABILITY_SILENCE,
		sim.ABILITY_WEAKEN,
	}
}

/**
 * Abilities that can't be cleansed. (These aren't actually debuffs but this app codes them as a debuff)
 */
func GetUncleansableDebuffs() []sim.Ability {
	return []sim.Ability{
		sim.ABILITY_AMPLIFY,
	}
}

/** Abilities that require a turn to do something */
func GetActionAbilities() []sim.Ability {
	return []sim.Ability{
		sim.ABILITY_REPAIR,
		sim.ABILITY_TANK_HEAL,
	}
}

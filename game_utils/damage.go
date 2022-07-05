package game_utils

import (
	"math"
	sim "simulator"
)

// TODO: Should this return the reduced damage or normal? Life steal against void?
/** Hits the monster with magic damage. Returns the remainder damage. */
func HitMonsterWithMagic(game sim.Game, target sim.GameMonster, magicDamage int) sim.BattleDamage {
	// consider forcefield
	if target.HasAbility(sim.ABILITY_FORCEFIELD) && magicDamage > sim.FORCEFIELD_MIN_DAMAGE {
		magicDamage = 1
	}

	// consider Divine shield
	if target.HasAbility(sim.ABILITY_DIVINE_SHIELD) {
		target.RemoveDivineShield()
		game.CreateAndAddBattleLog(sim.ABILITY_DIVINE_SHIELD, target, nil, 0)
		return sim.BattleDamage{
			attack:           1,
			damageDone:       0,
			remainder:        0,
			actualDamageDone: 0,
		}
	}

	if magicDamage < 1 {
		return sim.BattleDamage{}
	}

	// hit void
	if target.HasAbility(sim.ABILITY_VOID) {
		if magicDamage == 1 {
			return sim.BattleDamage{attack: 1, remainder: 1}
		}
		magicDamage = int(math.Floor(float64(magicDamage+1) / 2))
	}

	// hit void armor
	if target.HasAbility(sim.ABILITY_VOID_ARMOR) {
		if target.Armor > 0 {
			remainder := HitArmor(target, magicDamage)
			return sim.BattleDamage{
				attack:           1,
				damageDone:       0,
				remainder:        remainder,
				actualDamageDone: 0,
			}

		}

		// has void armor ability but no armor left
		remainder := HitHealth(target, magicDamage)
		return sim.BattleDamage{
			attack:           magicDamage,
			damageDone:       magicDamage,
			remainder:        remainder,
			actualDamageDone: magicDamage - remainder,
		}
	}

	// no void armor
	remainder := HitHealth(target, magicDamage)
	return sim.BattleDamage{
		attack:           magicDamage,
		damageDone:       magicDamage,
		remainder:        remainder,
		actualDamageDone: magicDamage - remainder,
	}
}

/**
Hits the monster with physical damage. Returns the remainder damage.
This doesn't consider Piercing, so return remainder and call this one more time after hitting the armor
*/
func HitMonsterWithPhysical(game *sim.Game, target sim.MonsterCard, damageAmount int) sim.BattleDamage {
	// consider forcefield
	if target.HasAbility(sim.ABILITY_FORCEFIELD) && damageAmount >= sim.FORCEFIELD_MIN_DAMAGE {
		damageAmount = 1
	}

	// For things like thorns, this returns 1 to show a successful attack.
	if target.HasAbility(sim.ABILITY_DIVINE_SHIELD) {
		target.RemoveDivineShield()
		game.CreateAndAddBattleLog(sim.ABILITY_DIVINE_SHIELD, target, nil, 0)
		return sim.BattleDamage{attack: 1}
	}

	if damageAmount < 1 {
		return sim.BattleDamage{}
	}

	// consider hitting shield
	if target.HasAbility(sim.ABILITY_SHIELD) {
		if damageAmount == 1 {
			return sim.BattleDamage{attack: 1, remainder: 1}
		}
		damageAmount = int(math.Floor(float64(damageAmount+1) / 2))
	}

	if target.Armor > 0 {
		remainder := HitArmor(target, damageAmount)
		return sim.BattleDamage{
			attack:           damageAmount,
			damageDone:       damageAmount,
			remainder:        remainder,
			actualDamageDone: damageAmount - remainder,
		}
	}

	// normal attack hitting health
	remainderDamage := HitHealth(target, damageAmount)
	return sim.BattleDamage{
		attack:           damageAmount,
		damageDone:       damageAmount,
		remainder:        remainderDamage,
		actualDamageDone: damageAmount - remainderDamage,
	}
}

/** Returns remainder damage after hitting armor. */
func HitArmor(target sim.MonsterCard, damageAmount int) int {
	remainderArmor := target.Armor - damageAmount
	if remainderArmor < 0 {
		// damage was bigger than the armor
		target.Armor = 0
		return remainderArmor * (-1)
	}

	target.Armor = remainderArmor
	return 0
}

/** Returns remainder damage after hitting health. */
func HitHealth(target sim.MonsterCard, damageAmount int) int {
	preHitHealth := target.Health
	target.AddHealth(-1 * damageAmount)
	if target.Health < 0 {
		return target.Health * (-1)
	}
	if target.Health == 0 {
		return damageAmount - preHitHealth
	}

	return 0
}

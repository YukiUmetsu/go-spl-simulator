package game_models

import (
	"math"
)

// TODO: Should this return the reduced damage or normal? Life steal against void?
/** Hits the monster with magic damage. Returns the remainder damage. */
func HitMonsterWithMagic(game Game, target *MonsterCard, magicDamage int) BattleDamage {
	// consider forcefield
	if target.HasAbility(ABILITY_FORCEFIELD) && magicDamage > FORCEFIELD_MIN_DAMAGE {
		magicDamage = 1
	}

	// consider Divine shield
	if target.HasAbility(ABILITY_DIVINE_SHIELD) {
		target.RemoveDivineShield()
		game.CreateAndAddBattleLog(BATTLE_ACTION_REMOVE_DIVINE_SHIELD, target, nil, 0)
		return BattleDamage{
			Attack:           1,
			DamageDone:       0,
			Remainder:        0,
			ActualDamageDone: 0,
		}
	}

	if magicDamage < 1 {
		return BattleDamage{}
	}

	// hit void
	if target.HasAbility(ABILITY_VOID) {
		if magicDamage == 1 {
			return BattleDamage{Attack: 1, Remainder: 1}
		}
		magicDamage = int(math.Floor(float64(magicDamage+1) / 2))
	}

	// hit void armor
	if target.HasAbility(ABILITY_VOID_ARMOR) {
		if target.Armor > 0 {
			remainder := HitArmor(*target, magicDamage)
			return BattleDamage{
				Attack:           1,
				DamageDone:       0,
				Remainder:        remainder,
				ActualDamageDone: 0,
			}

		}

		// has void armor ability but no armor left
		remainder := HitHealth(*target, magicDamage)
		return BattleDamage{
			Attack:           magicDamage,
			DamageDone:       magicDamage,
			Remainder:        remainder,
			ActualDamageDone: magicDamage - remainder,
		}
	}

	// no void armor
	remainder := HitHealth(*target, magicDamage)
	return BattleDamage{
		Attack:           magicDamage,
		DamageDone:       magicDamage,
		Remainder:        remainder,
		ActualDamageDone: magicDamage - remainder,
	}
}

/**
Hits the monster with physical damage. Returns the remainder damage.
This doesn't consider Piercing, so return remainder and call this one more time after hitting the armor
*/
func HitMonsterWithPhysical(game *Game, target MonsterCard, damageAmount int) BattleDamage {
	// consider forcefield
	if target.HasAbility(ABILITY_FORCEFIELD) && damageAmount >= FORCEFIELD_MIN_DAMAGE {
		damageAmount = 1
	}

	// For things like thorns, this returns 1 to show a successful attack.
	if target.HasAbility(ABILITY_DIVINE_SHIELD) {
		target.RemoveDivineShield()
		game.CreateAndAddBattleLog(BATTLE_ACTION_REMOVE_DIVINE_SHIELD, &target, nil, 0)
		return BattleDamage{Attack: 1}
	}

	if damageAmount < 1 {
		return BattleDamage{}
	}

	// consider hitting shield
	if target.HasAbility(ABILITY_SHIELD) {
		if damageAmount == 1 {
			return BattleDamage{Attack: 1, Remainder: 1}
		}
		damageAmount = int(math.Floor(float64(damageAmount+1) / 2))
	}

	if target.Armor > 0 {
		remainder := HitArmor(target, damageAmount)
		return BattleDamage{
			Attack:           damageAmount,
			DamageDone:       damageAmount,
			Remainder:        remainder,
			ActualDamageDone: damageAmount - remainder,
		}
	}

	// normal attack hitting health
	remainderDamage := HitHealth(target, damageAmount)
	return BattleDamage{
		Attack:           damageAmount,
		DamageDone:       damageAmount,
		Remainder:        remainderDamage,
		ActualDamageDone: damageAmount - remainderDamage,
	}
}

/** Returns remainder damage after hitting armor. */
func HitArmor(target MonsterCard, damageAmount int) int {
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
func HitHealth(target MonsterCard, damageAmount int) int {
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

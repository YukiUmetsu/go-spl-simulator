package game_utils

import (
	"math"
	"math/rand"
	sim "simulator"
)

func GetDidDodge(rulesets sim.Ruleset, attacker sim.MonsterCard, target sim.MonsterCard, attackType sim.AttackType) bool {
	// true strike
	if attacker.HasAbility(sim.ABILITY_TRUE_STRIKE) {
		return false
	}

	// magic always hits except phase
	if attackType == sim.AttackType.Magic && !target.HasAbility(sim.ABILITY_PHASE) {
		return false
	}

	// snare to flying
	if attacker.HasAbility(sim.ABILITY_SNARE) && target.HasAbility(sim.ABILITY_FLYING) {
		return false
	}

	// calculate dodge chance from speed difference
	speedDiff := target.GetPostAbilitySpeed() - attacker.GetPostAbilitySpeed()
	if StrArrContains(rulesets, sim.Ruleset.RULESET_REVERSE_SPEED) {
		speedDiff = -1 * speedDiff
	}
	var dodgeChance float64 = 0
	if speedDiff > 0 {
		dodgeChance = float64(speedDiff) * sim.SPEED_DODGE_CHANCE
	}

	// add dodge ability 25% chance to evade
	if target.HasAbility(sim.ABILITY_DODGE) {
		dodgeChance = dodgeChance + sim.DODGE_CHANCE
	}

	// add flying ability 25% chance to evade (if attacker doesn't have flying and snare)
	if target.HasAbility(sim.ABILITY_FLYING) && !attacker.HasAbility(sim.ABILITY_FLYING) && !target.HasDebuff(sim.ABILITY_SNARE) {
		dodgeChance = dodgeChance + sim.FLYING_DODGE_CHANCE
	}

	// +15% if attacker has blind
	if attacker.HasDebuff(sim.ABILITY_BLIND) {
		dodgeChance = dodgeChance + sim.BLIND_DODGE_CHANCE
	}
	return GetSuccessBelow(dodgeChance * 100)
}

func GetSuccessBelow(chance float64) bool {
	return math.Floor(float64(rand.Intn(101))) < chance
}

// Compare Attack Order
// https://support.splinterlands.com/hc/en-us/articles/4414334269460-Attack-Order
func NormalCompareAttackOrder(m1 sim.MonsterCard, m2 sim.MonsterCard) int {
	speedDiff := m1.GetPostAbilitySpeed() - m2.GetPostAbilitySpeed()
	if speedDiff != 0 {
		return speedDiff
	}

	if m1.Magic > 0 && m2.Magic == 0 {
		return 1
	}
	if m2.Magic > 0 && m1.Magic == 0 {
		return -1
	}
	if m1.Ranged > 0 && m2.Ranged == 0 {
		return 1
	}
	if m2.Ranged > 0 && m1.Ranged == 0 {
		return -1
	}

	if m1.GetRarity() != m2.GetRarity() {
		return m1.GetRarity() - m2.GetRarity()
	}
	return m1.GetLevel() - m2.GetLevel()
}

func ResolveFriendlyTies(m1 sim.MonsterCard, m2 sim.MonsterCard) int {
	m1Position := m1.GetCardPosition()
	m2Position := m2.GetCardPosition()
	if !m1.HasAction() && !m2.HasAction() {
		if m1Position < m2Position {
			return -1
		}
		return 1
	}

	return RandomTieBreaker()
}

func RandomTieBreaker() int {
	if rand.Intn(101) > 50 {
		return -1
	}
	return 1
}

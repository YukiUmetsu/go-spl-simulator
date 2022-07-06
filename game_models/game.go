package game_models

import (
	"math"
	"math/rand"
)

func GetDidDodge(rulesets Ruleset, attacker MonsterCard, target MonsterCard, attackType AttackType) bool {
	// true strike
	if attacker.HasAbility(ABILITY_TRUE_STRIKE) {
		return false
	}

	// magic always hits except phase
	if attackType == AttackType.Magic && !target.HasAbility(ABILITY_PHASE) {
		return false
	}

	// snare to flying
	if attacker.HasAbility(ABILITY_SNARE) && target.HasAbility(ABILITY_FLYING) {
		return false
	}

	// calculate dodge chance from speed difference
	speedDiff := target.GetPostAbilitySpeed() - attacker.GetPostAbilitySpeed()
	if StrArrContains(rulesets, Ruleset.RULESET_REVERSE_SPEED) {
		speedDiff = -1 * speedDiff
	}
	var dodgeChance float64 = 0
	if speedDiff > 0 {
		dodgeChance = float64(speedDiff) * SPEED_DODGE_CHANCE
	}

	// add dodge ability 25% chance to evade
	if target.HasAbility(ABILITY_DODGE) {
		dodgeChance = dodgeChance + DODGE_CHANCE
	}

	// add flying ability 25% chance to evade (if attacker doesn't have flying and snare)
	if target.HasAbility(ABILITY_FLYING) && !attacker.HasAbility(ABILITY_FLYING) && !target.HasDebuff(ABILITY_SNARE) {
		dodgeChance = dodgeChance + FLYING_DODGE_CHANCE
	}

	// +15% if attacker has blind
	if attacker.HasDebuff(ABILITY_BLIND) {
		dodgeChance = dodgeChance + BLIND_DODGE_CHANCE
	}
	return GetSuccessBelow(dodgeChance * 100)
}

func GetSuccessBelow(chance float64) bool {
	return math.Floor(float64(rand.Intn(101))) < chance
}

// Compare Attack Order
// https://support.splinterlands.com/hc/en-us/articles/4414334269460-Attack-Order
func NormalCompareAttackOrder(m1 MonsterCard, m2 MonsterCard) int {
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

func ResolveFriendlyTies(m1 MonsterCard, m2 MonsterCard) int {
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

func CardsArrIncludesMonster(cards []MonsterCard, m MonsterCard) bool {
	if len(cards) == 0 {
		return false
	}

	for _, c := range cards {
		if c.GameCard.Name == m.GameCard.Name {
			return true
		}
	}

	return false
}

// https://support.splinterlands.com/hc/en-us/articles/4414334269460-Attack-Order
func MonsterTurnComparator(m1 Monster, m2 Monster) bool {
	normalCompareDiff := NormalCompareAttackOrder(m1, m2)

	// Descending order
	if normalCompareDiff != 0 {
		return normalCompareDiff > 0
	}

	// resolve tie by order if the same team, else random
	if m1.GetTeamNumber() == m2.GetTeamNumber() {
		return ResolveFriendlyTies(m1, m2) > 0
	} else {
		return RandomTieBreaker() > 0
	}
}

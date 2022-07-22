package simulator_tests

import (
	"testing"

	. "github.com/YukiUmetsu/go-spl-simulator/game_models"
	"github.com/stretchr/testify/assert"
)

func TestSetCardPosition(t *testing.T) {
	m := GetDefaultFakeMonster(ATTACK_TYPE_MAGIC)
	m.SetCardPosition(3)
	assert.Equal(t, 3, m.GetCardPosition())
}

func TestIsAlive(t *testing.T) {
	// isAlive function returns true when health is above 0
	m := GetDefaultFakeMonster(ATTACK_TYPE_MAGIC)
	m.SetHealth(1)
	assert.True(t, m.IsAlive())

	// isAlive function returns false when heatlh is 0
	m.SetHealth(0)
	assert.False(t, m.IsAlive())
}

func TestSetHasTurnPassed(t *testing.T) {
	// passes turn correctly
	m := GetDefaultFakeMonster(ATTACK_TYPE_MAGIC)
	m.SetHasTurnPassed(true)
	assert.True(t, m.GetHasTurnPassed())
}

func TestSetIsOnlyMonster(t *testing.T) {
	m := GetDefaultFakeMonster(ATTACK_TYPE_MAGIC)
	m.SetIsOnlyMonster()
	assert.True(t, m.IsLastMonster())

	// last stand, x1.5 - round down
	m = GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MAGIC, []Ability{ABILITY_LAST_STAND})
	m.SetIsOnlyMonster()
	assert.Equal(t, 7, m.GetHealth())
}

func TestGetCleanCard(t *testing.T) {
	// GetCleanCard returns a card without any modifiers
	m := GetDefaultFakeMonster(ATTACK_TYPE_MAGIC)
	m.SetHealth(100)
	cleanCard := m.GetCleanCard()
	assert.Equal(t, 5, cleanCard.GetHealth())
}

func TestHasAction(t *testing.T) {
	// true for the first position
	m := GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.SetCardPosition(0)
	assert.True(t, m.HasAction())

	// returns false if not in first position with no abilities
	noAttackM := GetDefaultFakeMonster(ATTACK_TYPE_NO_ATTACK)
	noAttackM.SetCardPosition(2)
	assert.False(t, noAttackM.HasAction())
	noAttackM.SetCardPosition(0)
	assert.False(t, noAttackM.HasAction())

	// tank heal has action
	tankHealM := GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MELEE, []Ability{ABILITY_TANK_HEAL})
	tankHealM.SetCardPosition(1)
	assert.True(t, tankHealM.HasAction())

	// repair has action
	repairM := GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MELEE, []Ability{ABILITY_REPAIR})
	repairM.SetCardPosition(2)
	assert.True(t, repairM.HasAction())

	// cleanse has action
	cleanseM := GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MELEE, []Ability{ABILITY_CLEANSE})
	cleanseM.SetCardPosition(1)
	assert.True(t, cleanseM.HasAction())

	// triage has action
	triageM := GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MELEE, []Ability{ABILITY_TRIAGE})
	triageM.SetCardPosition(0)
	assert.True(t, triageM.HasAction())

	// returns true if not in first position with opportunity
	opM := GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MELEE, []Ability{ABILITY_OPPORTUNITY})
	opM.SetCardPosition(2)
	assert.True(t, opM.HasAction())

	// returns true if not in first position with sneak
	sneakM := GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MELEE, []Ability{ABILITY_SNEAK})
	sneakM.SetCardPosition(2)
	assert.True(t, sneakM.HasAction())

	// returns true if not in first position and is melee mayhem ruleset
	mmm := GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MELEE, []Ability{ABILITY_MELEE_MAYHEM})
	mmm.SetCardPosition(2)
	assert.True(t, mmm.HasAction())

	// returns false if in 3rd position with reach
	reachM := GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MELEE, []Ability{ABILITY_REACH})
	reachM.SetCardPosition(2)
	assert.False(t, reachM.HasAction())

	// returns true if in 2nd position with reach
	reachM.SetCardPosition(1)
	assert.True(t, reachM.HasAction())

	// magic can attack from anywhere
	magicM := GetDefaultFakeMonster(ATTACK_TYPE_MAGIC)
	magicM.SetCardPosition(3)
	assert.True(t, magicM.HasAction())

	// range can't attack from the first position
	rangeM := GetDefaultFakeMonster(ATTACK_TYPE_RANGED)
	rangeM.SetCardPosition(0)
	assert.False(t, rangeM.HasAction())

	// range can attack from the not first position
	rangeM.SetCardPosition(1)
	assert.True(t, rangeM.HasAction())

	// close range can attack from the first position
	cRangeM := GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_RANGED, []Ability{ABILITY_CLOSE_RANGE})
	cRangeM.SetCardPosition(0)
	assert.True(t, cRangeM.HasAction())

}

func TestAddDebuff(t *testing.T) {
	// doesn't add debuff when mosnter is dead
	m := GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.SetHealth(0)
	m.AddDebuff(ABILITY_BLIND)
	assert.False(t, m.HasDebuff(ABILITY_BLIND))

	// adds a debuff when applied
	m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.AddDebuff(ABILITY_BLIND)
	assert.Equal(t, 1, m.GetDebuffCount(ABILITY_BLIND))

	// Immunity
	m = GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MELEE, []Ability{ABILITY_IMMUNITY})
	m.AddDebuff(ABILITY_BLIND)
	assert.False(t, m.HasDebuff(ABILITY_BLIND))

	// Adds amplify debuff even if monster has immunity
	m = GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MELEE, []Ability{ABILITY_IMMUNITY})
	m.AddDebuff(ABILITY_AMPLIFY)
	assert.True(t, m.HasDebuff(ABILITY_AMPLIFY))

	// Cleanse removes all the GetDebuffs
	m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.AddDebuff(ABILITY_STUN)
	m.AddDebuff(ABILITY_BLIND)
	m.AddDebuff(ABILITY_POISON)
	assert.Equal(t, 3, len(m.GetDebuffs()))
	m.CleanseDebuffs()
	assert.Equal(t, 0, len(m.GetDebuffs()))

	// cleanse can only remove one cripple at a time
	m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.AddDebuff(ABILITY_CRIPPLE)
	m.AddDebuff(ABILITY_CRIPPLE)
	m.AddDebuff(ABILITY_CRIPPLE)
	m.SetHealth(2)
	assert.Equal(t, 3, m.GetDebuffs()[ABILITY_CRIPPLE])
	m.CleanseDebuffs()
	assert.Equal(t, 2, m.GetDebuffs()[ABILITY_CRIPPLE])

	// removing cripple adds a health back
	m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.AddDebuff(ABILITY_CRIPPLE)
	m.SetHealth(1)
	m.RemoveDebuff(ABILITY_CRIPPLE)
	assert.Equal(t, 2, m.GetHealth())

	// weaken reduce health
	m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.AddDebuff(ABILITY_WEAKEN)
	assert.Equal(t, TEST_DEFAULT_HEALTH-1, m.GetHealth())

	// removing weaken adds a health back
	m.RemoveDebuff(ABILITY_WEAKEN)
	assert.Equal(t, TEST_DEFAULT_HEALTH, m.GetHealth())

	// Rust removes 2 armor
	m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.AddDebuff(ABILITY_RUST)
	assert.Equal(t, TEST_DEFAULT_ARMOR-2, m.Armor)

	// adding rust sets armor to minimum of 0
	m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.SetArmor(1)
	m.AddDebuff(ABILITY_RUST)
	assert.Equal(t, 0, m.Armor)

	// removing rust adds armor back
	m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.AddDebuff(ABILITY_RUST)
	m.RemoveDebuff(ABILITY_RUST)
	assert.Equal(t, TEST_DEFAULT_ARMOR, m.GetArmor())
}

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

func TestRemoveBuff(t *testing.T) {
	// does nothing if the monster doesn't have the buff
	m := GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_RANGED, []Ability{ABILITY_WEAKEN})
	m.RemoveBuff(ABILITY_SCAVENGER)
	assert.Equal(t, TEST_DEFAULT_HEALTH, m.GetHealth())
	assert.Equal(t, 1, len(m.Abilities))

	// removes health if scavenger buff is removed
	m = GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_RANGED, []Ability{ABILITY_SCAVENGER})
	m.AddBuff(ABILITY_SCAVENGER)
	m.RemoveBuff(ABILITY_SCAVENGER)
	assert.Equal(t, TEST_DEFAULT_HEALTH, m.GetHealth())

	// doesn't remove health if scavenger buff is removed and health is already 1
	m = GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_RANGED, []Ability{ABILITY_SCAVENGER})
	m.AddBuff(ABILITY_SCAVENGER)
	m.SetHealth(1)
	m.RemoveBuff(ABILITY_SCAVENGER)
	assert.Equal(t, 1, m.GetHealth())

	// removes health if life leech buff is removed
	m = GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_RANGED, []Ability{ABILITY_LIFE_LEECH})
	m.AddBuff(ABILITY_LIFE_LEECH)
	m.SetHealth(2)
	m.RemoveBuff(ABILITY_LIFE_LEECH)
	assert.Equal(t, 1, m.GetHealth())

	// doesn't remove health if life leech buff is removed and health is already 1
	m = GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_RANGED, []Ability{ABILITY_LIFE_LEECH})
	m.AddBuff(ABILITY_LIFE_LEECH)
	m.SetHealth(1)
	m.RemoveBuff(ABILITY_LIFE_LEECH)
	assert.Equal(t, 1, m.GetHealth())

	// removes health if monster is at max post ability health when strengthen is removed
	m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.AddBuff(ABILITY_STRENGTHEN)
	assert.Equal(t, TEST_DEFAULT_HEALTH+1, m.GetHealth())
	m.RemoveBuff(ABILITY_STRENGTHEN)
	assert.Equal(t, TEST_DEFAULT_HEALTH, m.GetHealth())

	// sets armor down when removing protect and armor is higher than max
	m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.AddBuff(ABILITY_PROTECT)
	m.SetArmor(TEST_DEFAULT_ARMOR + 1)
	m.RemoveBuff(ABILITY_PROTECT)
	assert.Equal(t, TEST_DEFAULT_ARMOR, m.GetArmor())

	// doesn't change armor when removing protect and armor is lower than max
	m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.AddBuff(ABILITY_PROTECT)
	m.SetArmor(2)
	m.RemoveBuff(ABILITY_PROTECT)
	assert.Equal(t, 2, m.GetArmor())
}

func TestAddBuff(t *testing.T) {
	// doesn't add buff if monster is not alive
	m := GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.SetHealth(0)
	m.AddBuff(ABILITY_SCAVENGER)
	assert.Equal(t, 0, len(m.BuffMap))

	// adds buff normally if monster is aliv
	m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.AddBuff(ABILITY_SCAVENGER)
	assert.Equal(t, 1, len(m.BuffMap))

	// adds health if adding scavenger buff
	m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.AddBuff(ABILITY_SCAVENGER)
	assert.Equal(t, TEST_DEFAULT_HEALTH+1, m.GetHealth())

	// adds health if adding life leech buff
	m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.AddBuff(ABILITY_LIFE_LEECH)
	assert.Equal(t, TEST_DEFAULT_HEALTH+1, m.GetHealth())

	// adds health if adding strengthen buff
	m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.AddBuff(ABILITY_STRENGTHEN)
	assert.Equal(t, TEST_DEFAULT_HEALTH+1, m.GetHealth())

	// adds 2 armor if adding protect buff
	m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.AddBuff(ABILITY_PROTECT)
	assert.Equal(t, TEST_DEFAULT_ARMOR+2, m.GetArmor())
}

func TestPostAbilityMaxHealth(t *testing.T) {
	// Last stand x1.5 health
	m := GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_RANGED, []Ability{ABILITY_LAST_STAND})
	m.SetIsOnlyMonster()
	assert.Equal(t, 7, m.GetPostAbilityMaxHealth())

	// life leech increased health is also affected by summoner health
	m = GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MAGIC, []Ability{ABILITY_LAST_STAND})
	m.SetIsOnlyMonster()
	m.AddSummonerHealth(-1)
	assert.Equal(t, 6, m.GetPostAbilityMaxHealth())

	// life leech increased health is also affected by last stand
	m = GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MAGIC, []Ability{ABILITY_LAST_STAND})
	m.SetIsOnlyMonster()
	m.AddBuff(ABILITY_LIFE_LEECH)
	m.AddBuff(ABILITY_LIFE_LEECH)
	assert.Equal(t, 10, m.GetPostAbilityMaxHealth())

	// scavenger increased health is also affected by last stand
	m = GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MAGIC, []Ability{ABILITY_LAST_STAND})
	m.SetIsOnlyMonster()
	m.AddBuff(ABILITY_SCAVENGER)
	m.AddBuff(ABILITY_SCAVENGER)
	assert.Equal(t, 10, m.GetPostAbilityMaxHealth())

	// cripple decreased health is affected by last stand
	m = GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MAGIC, []Ability{ABILITY_LAST_STAND})
	m.SetIsOnlyMonster()
	m.AddDebuff(ABILITY_CRIPPLE)
	assert.Equal(t, 6, m.GetPostAbilityMaxHealth())

	// weaken is affected by last stand
	m = GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MAGIC, []Ability{ABILITY_LAST_STAND})
	m.SetIsOnlyMonster()
	m.AddDebuff(ABILITY_WEAKEN)
	assert.Equal(t, 6, m.GetPostAbilityMaxHealth())

	// strengthen is affected by last stand
	m = GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MAGIC, []Ability{ABILITY_LAST_STAND})
	m.SetIsOnlyMonster()
	m.AddBuff(ABILITY_STRENGTHEN)
	assert.Equal(t, 9, m.GetPostAbilityMaxHealth())

	// returns a minimum max health of 1
	m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.AddSummonerHealth(-10)
	assert.Equal(t, 1, m.GetHealth())
}

func TestHitHealth(t *testing.T) {
	// reduce health after hitting
	m := GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.HitHealth(2)
	assert.Equal(t, TEST_DEFAULT_HEALTH-2, m.GetHealth())

	// returns the remaining damage after hitting health down to 0
	m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	remainder := m.HitHealth(7)
	assert.Equal(t, 7-TEST_DEFAULT_HEALTH, remainder)

	// can't add health if monster is dead
	m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.HitHealth(7)
	m.AddHealth(1)
	assert.False(t, m.IsAlive())
	assert.Equal(t, 0, m.GetHealth())

	// starting health is changed when summoner increases it
	m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.AddSummonerHealth(2)
	assert.Equal(t, TEST_DEFAULT_HEALTH+2, m.GetStartingHealth())

	// starting health is changed when summoner decreases it
	m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.AddSummonerHealth(-2)
	assert.Equal(t, TEST_DEFAULT_HEALTH-2, m.GetStartingHealth())
}

func TestResurrect(t *testing.T) {
	// resurrects monster to 1 health
	m := GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.HitHealth(10)
	m.Resurrect()
	assert.Equal(t, 1, m.GetHealth())

	// resurrects the monster with divine shield if it had one
	m = GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MELEE, []Ability{ABILITY_DIVINE_SHIELD})
	m.RemoveDivineShield()
	assert.False(t, m.HasAbility(ABILITY_DIVINE_SHIELD))
	m.HitHealth(10)
	m.Resurrect()
	assert.True(t, m.HasAbility(ABILITY_DIVINE_SHIELD))

	// resurrects without divine shield if it didn't originally have it
	m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.HitHealth(10)
	m.Resurrect()
	assert.False(t, m.HasAbility(ABILITY_DIVINE_SHIELD))

	// resurrects with the max armor it can have
	m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m.SetArmor(0)
	m.HitHealth(10)
	m.Resurrect()
	assert.Equal(t, TEST_DEFAULT_ARMOR, m.GetArmor())
}

func TestGetPostAbilityMaxArmor(t *testing.T) {
	var m *MonsterCard

	type testCase struct {
		Name          string
		Buffs         []Ability
		Debuffs       []Ability
		SummonerArmor int
		ExpectedArmor int
		IsLastStand   bool
	}

	testCases := []testCase{
		testCase{
			"returns correct armor with last stand",
			[]Ability{},
			[]Ability{},
			0,
			7,
			true,
		},
		testCase{
			"returns correct armor with summoner armor modifiers",
			[]Ability{},
			[]Ability{},
			2,
			9,
			true,
		},
		testCase{
			"returns correct armor with protect",
			[]Ability{ABILITY_PROTECT, ABILITY_PROTECT},
			[]Ability{},
			0,
			11,
			true,
		},
		testCase{
			"returns correct armor with rust",
			[]Ability{},
			[]Ability{ABILITY_RUST, ABILITY_RUST},
			0,
			3,
			true,
		},
		testCase{
			"returns correct armor with all modifiers",
			[]Ability{ABILITY_PROTECT, ABILITY_PROTECT},
			[]Ability{ABILITY_RUST},
			-1,
			8,
			true,
		},
		testCase{
			"returns correct armor with no other modifiers",
			[]Ability{},
			[]Ability{},
			0,
			TEST_DEFAULT_ARMOR,
			false,
		},
		testCase{
			"returns correct armor with summoner armor modifiers",
			[]Ability{},
			[]Ability{},
			2,
			TEST_DEFAULT_ARMOR + 2,
			false,
		},
		testCase{
			"returns correct armor with protect",
			[]Ability{ABILITY_PROTECT},
			[]Ability{},
			0,
			TEST_DEFAULT_ARMOR + 2,
			false,
		},
		testCase{
			"returns correct armor with rust",
			[]Ability{},
			[]Ability{ABILITY_RUST},
			0,
			TEST_DEFAULT_ARMOR - 2,
			false,
		},
		testCase{
			"returns correct armor with all modifiers",
			[]Ability{ABILITY_PROTECT, ABILITY_PROTECT},
			[]Ability{ABILITY_RUST},
			-1,
			TEST_DEFAULT_ARMOR + 2 - 1,
			false,
		},
	}

	for _, tc := range testCases {
		if tc.IsLastStand {
			m = GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MELEE, []Ability{ABILITY_LAST_STAND})
			m.SetIsOnlyMonster()
		} else {
			m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
		}

		if len(tc.Buffs) > 0 {
			for _, buff := range tc.Buffs {
				m.AddBuff(buff)
			}
		}
		if len(tc.Debuffs) > 0 {
			for _, debuff := range tc.Debuffs {
				m.AddDebuff(debuff)
			}
		}
		if tc.SummonerArmor != 0 {
			m.AddSummonerArmor(tc.SummonerArmor)
		}
		assert.Equal(t, tc.ExpectedArmor, m.GetPostAbilityMaxArmor())
	}
}

func TestGetPostAbilitySpeed(t *testing.T) {
	var m *MonsterCard

	type testCase struct {
		Name           string
		Buffs          []Ability
		Debuffs        []Ability
		SummonerSpeed  int
		ExpectedSpeed  int
		IsLastStand    bool
		HasEnrage      bool
		ActivateEnrage bool
	}

	testCases := []testCase{
		testCase{
			"returns correct speed with no other modifiers",
			[]Ability{},
			[]Ability{},
			0,
			7,
			true,
			false,
			false,
		},
		testCase{
			"returns correct speed with summoner speed modifiers",
			[]Ability{},
			[]Ability{},
			2,
			10,
			true,
			false,
			false,
		},
		testCase{
			"returns correct speed with enrage ability and enrage active (enrage rounded up)",
			[]Ability{},
			[]Ability{},
			0,
			11,
			true,
			true,
			true,
		},
		testCase{
			"returns correct speed with enrage ability and enrage inactive",
			[]Ability{},
			[]Ability{},
			0,
			7,
			true,
			true,
			false,
		},
		testCase{
			"returns correct speed with swiftness",
			[]Ability{ABILITY_SWIFTNESS},
			[]Ability{},
			0,
			8,
			true,
			false,
			false,
		},
		testCase{
			"returns correct speed with slow",
			[]Ability{},
			[]Ability{ABILITY_SLOW},
			0,
			6,
			true,
			false,
			false,
		},
		testCase{
			"returns correct speed with summoner speed",
			[]Ability{},
			[]Ability{},
			2,
			10,
			true,
			false,
			false,
		},
		testCase{
			"returns correct speed with all modifiers",
			[]Ability{ABILITY_SWIFTNESS, ABILITY_SWIFTNESS},
			[]Ability{ABILITY_SLOW},
			2,
			16,
			true,
			true,
			true,
		},
		// Without last stand
		testCase{
			"returns correct speed with summoner speed modifierss",
			[]Ability{},
			[]Ability{},
			1,
			6,
			false,
			false,
			false,
		},
		testCase{
			"returns correct speed with enrage ability and enrage active",
			[]Ability{},
			[]Ability{},
			0,
			8,
			false,
			true,
			true,
		},
		testCase{
			"returns correct speed with enrage ability and enrage active",
			[]Ability{},
			[]Ability{},
			0,
			5,
			false,
			true,
			false,
		},
		testCase{
			"returns correct speed with swiftness",
			[]Ability{ABILITY_SWIFTNESS},
			[]Ability{},
			0,
			6,
			false,
			false,
			false,
		},
		testCase{
			"returns correct speed with slow",
			[]Ability{},
			[]Ability{ABILITY_SLOW},
			0,
			4,
			false,
			false,
			false,
		},
		testCase{
			"returns correct speed with all modifiers",
			[]Ability{ABILITY_SWIFTNESS, ABILITY_SWIFTNESS},
			[]Ability{ABILITY_SLOW},
			2,
			8,
			false,
			false,
			false,
		},
		testCase{
			"returns minimum of 1",
			[]Ability{},
			[]Ability{},
			-10,
			1,
			false,
			false,
			false,
		},
	}

	// returns correct speed with no other modifiers

	for _, tc := range testCases {
		if tc.IsLastStand {
			m = GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MELEE, []Ability{ABILITY_LAST_STAND})
			m.SetIsOnlyMonster()
		} else {
			m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
		}
		if tc.HasEnrage {
			// for enrage ability to work
			m.AddAbility(ABILITY_ENRAGE)
		}
		if tc.ActivateEnrage {
			m.SetHealth(TEST_DEFAULT_HEALTH - 1)
		}

		if len(tc.Buffs) > 0 {
			for _, buff := range tc.Buffs {
				m.AddBuff(buff)
			}
		}
		if len(tc.Debuffs) > 0 {
			for _, debuff := range tc.Debuffs {
				m.AddDebuff(debuff)
			}
		}
		if tc.SummonerSpeed != 0 {
			m.AddSummonerSpeed(tc.SummonerSpeed)
		}
		assert.Equal(t, tc.ExpectedSpeed, m.GetPostAbilitySpeed())
	}
}

func TestGetPostAbilityMagic(t *testing.T) {
	var m *MonsterCard

	type testCase struct {
		Name          string
		Debuffs       []Ability
		SummonerMagic int
		ExpectedMagic int
		IsLastStand   bool
	}

	testCases := []testCase{
		testCase{
			"returns 0 if there's no magic",
			[]Ability{},
			-10,
			1,
			false,
		},
		testCase{
			"returns correct magic with summoner magic",
			[]Ability{},
			1,
			6,
			false,
		},
		testCase{
			"returns correct magic with halving",
			[]Ability{ABILITY_HALVING},
			0,
			2,
			false,
		},
		testCase{
			"returns correct magic with last stand",
			[]Ability{},
			0,
			7,
			true,
		},
		testCase{
			"returns correct magic with silence",
			[]Ability{ABILITY_SILENCE, ABILITY_SILENCE},
			0,
			3,
			false,
		},
		testCase{
			"returns correct magic with all modifiers",
			[]Ability{ABILITY_HALVING, ABILITY_SILENCE},
			-1,
			1,
			true,
		},
	}

	for _, tc := range testCases {
		if tc.IsLastStand {
			m = GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MAGIC, []Ability{ABILITY_LAST_STAND})
			m.SetIsOnlyMonster()
		} else {
			m = GetDefaultFakeMonster(ATTACK_TYPE_MAGIC)
		}

		if len(tc.Debuffs) > 0 {
			for _, debuff := range tc.Debuffs {
				m.AddDebuff(debuff)
			}
		}
		if tc.SummonerMagic != 0 {
			m.AddSummonerMagic(tc.SummonerMagic)
		}
		assert.Equal(t, tc.ExpectedMagic, m.GetPostAbilityMagic())
	}

}

func TestGetPostAbilityRange(t *testing.T) {
	var m *MonsterCard

	type testCase struct {
		Name          string
		Debuffs       []Ability
		SummonerRange int
		ExpectedRange int
		IsLastStand   bool
	}

	testCases := []testCase{
		testCase{
			"returns correct ranged with summoner range'",
			[]Ability{},
			1,
			6,
			false,
		},
		testCase{
			"returns correct ranged with halving",
			[]Ability{ABILITY_HALVING},
			0,
			2,
			false,
		},
		testCase{
			"returns correct ranged with last stand",
			[]Ability{},
			0,
			7,
			true,
		},
		testCase{
			"returns correct ranged with headwind",
			[]Ability{ABILITY_HEADWINDS, ABILITY_HEADWINDS},
			0,
			3,
			false,
		},
		testCase{
			"returns correct ranged with all modifiers",
			[]Ability{ABILITY_HALVING, ABILITY_HEADWINDS},
			-1,
			1,
			true,
		},
	}

	for _, tc := range testCases {
		if tc.IsLastStand {
			m = GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_RANGED, []Ability{ABILITY_LAST_STAND})
			m.SetIsOnlyMonster()
		} else {
			m = GetDefaultFakeMonster(ATTACK_TYPE_RANGED)
		}

		if len(tc.Debuffs) > 0 {
			for _, debuff := range tc.Debuffs {
				m.AddDebuff(debuff)
			}
		}
		if tc.SummonerRange != 0 {
			m.AddSummonerRanged(tc.SummonerRange)
		}
		assert.Equal(t, tc.ExpectedRange, m.GetPostAbilityRange())
	}

}

func TestGetPostAbilityMelee(t *testing.T) {
	var m *MonsterCard

	type testCase struct {
		Name           string
		Buffs          []Ability
		Debuffs        []Ability
		SummonerMelee  int
		ExpectedMelee  int
		IsLastStand    bool
		HasEnraged     bool
		ActivateEnrage bool
	}

	testCases := []testCase{
		testCase{
			"returns correct melee with summoner range'",
			[]Ability{},
			[]Ability{},
			1,
			6,
			false,
			false,
			false,
		},
		testCase{
			"returns correct melee with halving",
			[]Ability{},
			[]Ability{ABILITY_HALVING},
			0,
			2,
			false,
			false,
			false,
		},
		testCase{
			"returns correct melee with last stand",
			[]Ability{},
			[]Ability{},
			0,
			7,
			true,
			false,
			false,
		},
		testCase{
			"returns correct melee with demoralize",
			[]Ability{},
			[]Ability{ABILITY_DEMORALIZE, ABILITY_DEMORALIZE},
			0,
			3,
			false,
			false,
			false,
		},
		testCase{
			"returns correct melee with inspire",
			[]Ability{ABILITY_INSPIRE},
			[]Ability{},
			0,
			6,
			false,
			false,
			false,
		},
		testCase{
			"returns correct melee with enrage",
			[]Ability{},
			[]Ability{},
			0,
			8,
			false,
			true,
			true,
		},
		testCase{
			"returns correct melee with all modifiers",
			[]Ability{ABILITY_INSPIRE},
			[]Ability{ABILITY_HALVING, ABILITY_DEMORALIZE},
			-1,
			2,
			true,
			false,
			false,
		},
	}

	for _, tc := range testCases {
		if tc.IsLastStand {
			m = GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MELEE, []Ability{ABILITY_LAST_STAND})
			m.SetIsOnlyMonster()
		} else {
			m = GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
		}

		if tc.HasEnraged {
			m.AddAbility(ABILITY_ENRAGE)
		}
		if tc.ActivateEnrage {
			m.AddHealth(-1)
		}
		if len(tc.Buffs) > 0 {
			for _, buff := range tc.Buffs {
				m.AddBuff(buff)
			}
		}

		if len(tc.Debuffs) > 0 {
			for _, debuff := range tc.Debuffs {
				m.AddDebuff(debuff)
			}
		}
		if tc.SummonerMelee != 0 {
			m.AddSummonerMelee(tc.SummonerMelee)
		}
		assert.Equal(t, tc.ExpectedMelee, m.GetPostAbilityMelee())
	}

}

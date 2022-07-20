package game_models

import (
	"fmt"
	"testing"

	utils "github.com/YukiUmetsu/go-spl-simulator/game_utils"

	"github.com/stretchr/testify/assert"
)

func TestGetDodgeChance(t *testing.T) {
	attacker := GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	target := GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	target.SetTeam(TEAM_NUM_TWO)

	attackerSpeed2 := GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	attackerSpeed2.Speed = 2
	targetPhaseSpeed8 := GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	targetPhaseSpeed8.AddAbility(ABILITY_PHASE)
	targetPhaseSpeed8.Speed = 8
	targetPhaseSpeed8.SetTeam(TEAM_NUM_TWO)

	attackerTrueStrike := GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	attackerTrueStrike.AddAbility(ABILITY_TRUE_STRIKE)

	attackerSnare := GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	attackerSnare.AddAbility(ABILITY_SNARE)
	targetFlying := GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	targetFlying.AddAbility(ABILITY_FLYING)

	targetSpeed3 := GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	targetSpeed3.Speed = 3

	targetSpeed1 := GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	targetSpeed1.Speed = 1

	attackerFlying := GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	attackerFlying.AddAbility(ABILITY_FLYING)

	attackerBlind := GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	attackerBlind.AddDebuff(ABILITY_BLIND)
	attackerBlind.AddDebuff(ABILITY_BLIND)

	targetDodge := GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	targetDodge.AddAbility(ABILITY_DODGE)

	type FuncInput struct {
		Rulesets   []Ruleset
		Attacker   *MonsterCard
		Target     *MonsterCard
		AttackType CardAttackType
	}

	type testCase struct {
		Name           string
		Input          FuncInput
		ExpectedOutput float64
	}

	testCases := []testCase{
		testCase{
			Name: "returns false if attack type is magic and target does not have phase",
			Input: FuncInput{
				Rulesets:   []Ruleset{RULESET_EARTHQUAKE},
				Attacker:   attacker,
				Target:     target,
				AttackType: ATTACK_TYPE_MAGIC,
			},
			ExpectedOutput: float64(0),
		},
		testCase{
			Name: "attack type is magic and target has phase",
			Input: FuncInput{
				Rulesets:   []Ruleset{RULESET_EARTHQUAKE},
				Attacker:   attackerSpeed2,
				Target:     targetPhaseSpeed8,
				AttackType: ATTACK_TYPE_MAGIC,
			},
			ExpectedOutput: float64(0.6),
		},
		testCase{
			Name: "returns false if attacking monster has true strike",
			Input: FuncInput{
				Rulesets:   []Ruleset{RULESET_EARTHQUAKE},
				Attacker:   attackerTrueStrike,
				Target:     targetPhaseSpeed8,
				AttackType: ATTACK_TYPE_RANGED,
			},
			ExpectedOutput: float64(0),
		},
		testCase{
			Name: "no dodging if attacking monster has snamre and attack target is flying",
			Input: FuncInput{
				Rulesets:   []Ruleset{RULESET_EARTHQUAKE},
				Attacker:   attackerSnare,
				Target:     targetFlying,
				AttackType: ATTACK_TYPE_RANGED,
			},
			ExpectedOutput: float64(0),
		},
		testCase{
			Name: "gives 10 percent chance per speed difference if attacking is slower",
			Input: FuncInput{
				Rulesets:   []Ruleset{RULESET_EARTHQUAKE},
				Attacker:   attackerSpeed2,
				Target:     targetSpeed3,
				AttackType: ATTACK_TYPE_RANGED,
			},
			ExpectedOutput: float64(0.1),
		},
		testCase{
			Name: "gives 10 percent chance per speed difference if attacking is faster and ruleset is reverse speed",
			Input: FuncInput{
				Rulesets:   []Ruleset{RULESET_REVERSE_SPEED},
				Attacker:   attackerSpeed2,
				Target:     targetSpeed1,
				AttackType: ATTACK_TYPE_RANGED,
			},
			ExpectedOutput: float64(0.1),
		},
		testCase{
			Name: "gives no chance to dodge if attacker is faster",
			Input: FuncInput{
				Rulesets:   []Ruleset{RULESET_STANDARD},
				Attacker:   attackerSpeed2,
				Target:     targetSpeed1,
				AttackType: ATTACK_TYPE_RANGED,
			},
			ExpectedOutput: float64(0),
		},
		testCase{
			Name: "gives 25 percent dodge chance for flying",
			Input: FuncInput{
				Rulesets:   []Ruleset{RULESET_STANDARD},
				Attacker:   attacker,
				Target:     targetFlying,
				AttackType: ATTACK_TYPE_RANGED,
			},
			ExpectedOutput: float64(0.25),
		},
		testCase{
			Name: "doesn't give dodge chance for flying if attacking is also flying",
			Input: FuncInput{
				Rulesets:   []Ruleset{RULESET_STANDARD},
				Attacker:   attackerFlying,
				Target:     targetFlying,
				AttackType: ATTACK_TYPE_RANGED,
			},
			ExpectedOutput: float64(0),
		},
		testCase{
			Name: "gives 15 percent chance for blind and it doesn't stack",
			Input: FuncInput{
				Rulesets:   []Ruleset{RULESET_STANDARD},
				Attacker:   attackerBlind,
				Target:     target,
				AttackType: ATTACK_TYPE_RANGED,
			},
			ExpectedOutput: float64(0.15),
		},
		testCase{
			Name: "gives 25 percent chance for dodge",
			Input: FuncInput{
				Rulesets:   []Ruleset{RULESET_STANDARD},
				Attacker:   attacker,
				Target:     targetDodge,
				AttackType: ATTACK_TYPE_RANGED,
			},
			ExpectedOutput: float64(0.25),
		},
	}

	for _, tc := range testCases {
		dodgeChance := GetDodgeChance(tc.Input.Rulesets, tc.Input.Attacker, tc.Input.Target, tc.Input.AttackType)
		isExpected := utils.AlmostEqualFloat(tc.ExpectedOutput, dodgeChance)
		assert.True(t, isExpected)
		if !isExpected {
			fmt.Println(tc.ExpectedOutput, " - ", dodgeChance)
		}
	}
}

func TestNormalCompareAttackOrder(t *testing.T) {
	mSpeed2 := GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	mSpeed2.Speed = 2
	mSpeed1 := GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	mSpeed1.Speed = 1

	magicM := GetDefaultFakeMonster(ATTACK_TYPE_MAGIC)
	meleeM := GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	rangedM := GetDefaultFakeMonster(ATTACK_TYPE_RANGED)

	legendaryLvl1 := CreateMonsterOfRarityAndLevel(4, 1, MONSTER, ATTACK_TYPE_MELEE)
	commonLvl1 := CreateMonsterOfRarityAndLevel(1, 1, MONSTER, ATTACK_TYPE_MELEE)
	commonLvl2 := CreateMonsterOfRarityAndLevel(1, 2, MONSTER, ATTACK_TYPE_MELEE)

	type InputMonsters struct {
		M1 *MonsterCard
		M2 *MonsterCard
	}
	type testCase struct {
		Name           string
		Input          InputMonsters
		ExpectedOutput int
	}

	testCases := []testCase{
		{
			Name: "returns the difference of the monster speed if speed is different",
			Input: InputMonsters{
				M1: mSpeed2,
				M2: mSpeed1,
			},
			ExpectedOutput: int(1),
		},
		{
			Name: "returns the difference of the monster speed if speed is different 2",
			Input: InputMonsters{
				M1: mSpeed1,
				M2: mSpeed2,
			},
			ExpectedOutput: int(-1),
		},
		// SPEED TIES
		{
			Name: "returns the magic user as faster if other one has no magic",
			Input: InputMonsters{
				M1: magicM,
				M2: meleeM,
			},
			ExpectedOutput: int(1),
		},
		{
			Name: "returns the ranged user as faster if other one has no ranged or magic",
			Input: InputMonsters{
				M1: rangedM,
				M2: meleeM,
			},
			ExpectedOutput: int(1),
		},
		{
			Name: "returns the rarer monster if the rest of stats are same",
			Input: InputMonsters{
				M1: legendaryLvl1,
				M2: commonLvl1,
			},
			ExpectedOutput: int(3),
		},
		{
			Name: "returns the higher level monster if the rest of stats are same",
			Input: InputMonsters{
				M1: commonLvl2,
				M2: commonLvl1,
			},
			ExpectedOutput: int(1),
		},
	}

	for _, tc := range testCases {
		out := NormalCompareAttackOrder(tc.Input.M1, tc.Input.M2)
		isExpected := assert.Equal(t, tc.ExpectedOutput, out)
		if !isExpected {
			fmt.Println(tc.ExpectedOutput, " - ", out)
		}
	}
}

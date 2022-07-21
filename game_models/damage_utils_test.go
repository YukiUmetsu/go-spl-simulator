package game_models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHitMonsterWithMagic(t *testing.T) {
	game := CreateFakeGame()
	mForceField := GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MAGIC, []Ability{ABILITY_FORCEFIELD})
	mForceField2 := GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MAGIC, []Ability{ABILITY_FORCEFIELD})
	mDivineShield := GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MAGIC, []Ability{ABILITY_DIVINE_SHIELD})
	mVoid := GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MAGIC, []Ability{ABILITY_VOID})
	mVoid2 := GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MAGIC, []Ability{ABILITY_VOID})
	mVoid3 := GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MAGIC, []Ability{ABILITY_VOID})
	mVoidArmor := GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MAGIC, []Ability{ABILITY_VOID_ARMOR})
	mVoidAndVoidArmor := GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MAGIC, []Ability{ABILITY_VOID_ARMOR, ABILITY_VOID})
	target := GetDefaultFakeMonster(ATTACK_TYPE_MAGIC)

	type testCase struct {
		Name                string
		Target              *MonsterCard
		MagicDamage         int
		OriginalHealth      int
		ExpectedHealthAfter int
		OriginalArmor       int
		ExpectedArmorAfter  int
	}

	testCases := []testCase{
		testCase{
			"only does 1 damage if attack target has force field and attack is >= 5",
			mForceField,
			5,
			mForceField.GetHealth(),
			mForceField.GetHealth() - 1,
			0,
			0,
		},
		testCase{
			"does full damage if attack target has force field and attack is < 5",
			mForceField2,
			2,
			mForceField2.GetHealth(),
			mForceField2.GetHealth() - 2,
			0,
			0,
		},
		testCase{
			"does no damage and removes divine shield if target has divine shield",
			mDivineShield,
			2,
			mDivineShield.GetHealth(),
			mDivineShield.GetHealth(),
			0,
			0,
		},
		testCase{
			"Void 4 -> 2",
			mVoid,
			4,
			mVoid.GetHealth(),
			mVoid.GetHealth() - 2,
			0,
			0,
		},
		testCase{
			"Void 1 -> 0",
			mVoid2,
			1,
			mVoid2.GetHealth(),
			mVoid2.GetHealth(),
			0,
			0,
		},
		testCase{
			"Void 3 -> 2",
			mVoid3,
			3,
			mVoid3.GetHealth(),
			mVoid3.GetHealth() - 2,
			0,
			0,
		},
		testCase{
			"Void Armor",
			mVoidArmor,
			1,
			mVoid3.GetHealth(),
			mVoid3.GetHealth(),
			mVoidArmor.Armor,
			mVoidArmor.Armor - 1,
		},
		testCase{
			"Void Armor Big Damage",
			mVoidArmor,
			100,
			mVoid3.GetHealth(),
			mVoid3.GetHealth(),
			mVoidArmor.Armor,
			0,
		},
		testCase{
			"Void & Void Armor",
			mVoidAndVoidArmor,
			3,
			mVoidAndVoidArmor.GetHealth(),
			mVoidAndVoidArmor.GetHealth(),
			mVoidAndVoidArmor.Armor,
			mVoidAndVoidArmor.Armor - 2,
		},
		testCase{
			"Hit correctly without any modifiers",
			target,
			3,
			target.GetHealth(),
			target.GetHealth() - 3,
			0,
			0,
		},
	}

	for _, tc := range testCases {
		HitMonsterWithMagic(game, tc.Target, tc.MagicDamage)
		assert.Equal(t, tc.ExpectedHealthAfter, tc.Target.GetHealth())

		if tc.OriginalArmor != 0 && tc.ExpectedArmorAfter != 0 {
			assert.Equal(t, tc.ExpectedArmorAfter, tc.Target.Armor)
		}

		if tc.Target.hadDivineShield {
			assert.False(t, tc.Target.HasAbility(ABILITY_DIVINE_SHIELD))
		}
	}
}

func TestHitMonsterWithPhysical(t *testing.T) {
	game := CreateFakeGame()
	target := GetDefaultFakeMonster(ATTACK_TYPE_MAGIC)
	mForceField := GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MAGIC, []Ability{ABILITY_FORCEFIELD})
	mForceField.Armor = 0
	mForceField2 := GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MAGIC, []Ability{ABILITY_FORCEFIELD})
	mForceField2.Armor = 0
	mDivineShield := GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MAGIC, []Ability{ABILITY_DIVINE_SHIELD})
	mShield := GetDefaultFakeMonsterWithAbility(ATTACK_TYPE_MAGIC, []Ability{ABILITY_SHIELD})
	mShield.Armor = 0

	type testCase struct {
		Name                string
		Target              *MonsterCard
		PhysicalDamage      int
		OriginalHealth      int
		ExpectedHealthAfter int
		OriginalArmor       int
		ExpectedArmorAfter  int
	}

	testCases := []testCase{
		testCase{
			"Hits armor if it has any",
			target,
			3,
			target.GetHealth(),
			target.GetHealth(),
			target.Armor,
			target.Armor - 3,
		},
		testCase{
			"Hits armor big and remove all armor",
			target,
			20,
			target.GetHealth(),
			target.GetHealth(),
			target.Armor,
			0,
		},
		testCase{
			"only does 1 damage if attack target has force field and attack is >= 5",
			mForceField,
			5,
			mForceField.GetHealth(),
			mForceField.GetHealth() - 1,
			0,
			0,
		},
		testCase{
			"does full damage if attack target has force field and attack is < 5",
			mForceField2,
			2,
			mForceField2.GetHealth(),
			mForceField2.GetHealth() - 2,
			0,
			0,
		},
		testCase{
			"does no damage and removes divine shield if target has divine shield",
			mDivineShield,
			2,
			mDivineShield.GetHealth(),
			mDivineShield.GetHealth(),
			0,
			0,
		},
		testCase{
			"1 attack does 0 damage to Shield",
			mShield,
			1,
			mShield.GetHealth(),
			mShield.GetHealth(),
			0,
			0,
		},
		testCase{
			"Shield 1/2 damage",
			mShield,
			3,
			mShield.GetHealth(),
			mShield.GetHealth() - 2,
			0,
			0,
		},
	}

	for _, tc := range testCases {
		HitMonsterWithPhysical(game, tc.Target, tc.PhysicalDamage)
		assert.Equal(t, tc.ExpectedHealthAfter, tc.Target.GetHealth())

		if tc.OriginalArmor != 0 && tc.ExpectedArmorAfter != 0 {
			assert.Equal(t, tc.ExpectedArmorAfter, tc.Target.Armor)
		}

		if tc.Target.hadDivineShield {
			assert.False(t, tc.Target.HasAbility(ABILITY_DIVINE_SHIELD))
		}
	}
}

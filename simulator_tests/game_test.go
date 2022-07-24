package simulator_tests

import (
	"testing"

	. "github.com/YukiUmetsu/go-spl-simulator/game_models"
	"github.com/stretchr/testify/assert"
)

func TestMaybeApplyThorns(t *testing.T) {
	var game *Game
	var t1 *GameTeam
	var t2 *GameTeam
	var attacker *MonsterCard
	var target *MonsterCard

	setUp := func() {
		game, t1, t2 = CreateFakeGameAndTeams()
		attacker = t1.GetFirstAliveMonster()
		target = t2.GetFirstAliveMonster()
		target.AddAbility(ABILITY_THORNS)
	}

	// doesn't do anything if attack type is not melee
	setUp()
	attacker = t1.GetMonstersList()[1]
	game.MaybeApplyThorns(attacker, target, ATTACK_TYPE_RANGED)
	assert.Equal(t, TEST_DEFAULT_HEALTH, attacker.GetHealth())
	assert.Equal(t, TEST_DEFAULT_ARMOR, attacker.GetArmor())

	attacker = t1.GetMonstersList()[2]
	game.MaybeApplyThorns(attacker, target, ATTACK_TYPE_MAGIC)
	assert.Equal(t, TEST_DEFAULT_HEALTH, attacker.GetHealth())
	assert.Equal(t, TEST_DEFAULT_ARMOR, attacker.GetArmor())

	// doesn't do anything if attack target has no thorns
	setUp()
	target.RemoveAbility(ABILITY_THORNS)
	game.MaybeApplyThorns(attacker, target, ATTACK_TYPE_MELEE)
	assert.Equal(t, TEST_DEFAULT_HEALTH, attacker.GetHealth())
	assert.Equal(t, TEST_DEFAULT_ARMOR, attacker.GetArmor())

	// doesn't do anything if attacker has reflection shield
	setUp()
	attacker.AddAbility(ABILITY_REFLECTION_SHIELD)
	game.MaybeApplyThorns(attacker, target, ATTACK_TYPE_MELEE)
	assert.Equal(t, TEST_DEFAULT_HEALTH, attacker.GetHealth())
	assert.Equal(t, TEST_DEFAULT_ARMOR, attacker.GetArmor())

	// attacker takes 2 damage if target has thorns
	setUp()
	game.MaybeApplyThorns(attacker, target, ATTACK_TYPE_MELEE)
	assert.Equal(t, TEST_DEFAULT_HEALTH, attacker.GetHealth())
	assert.Equal(t, TEST_DEFAULT_ARMOR-2, attacker.GetArmor())

	// attacker takes 3 damage if target has thorns and attacker has amplify
	setUp()
	attacker.AddDebuff(ABILITY_AMPLIFY)
	game.MaybeApplyThorns(attacker, target, ATTACK_TYPE_MELEE)
	assert.Equal(t, TEST_DEFAULT_HEALTH, attacker.GetHealth())
	assert.Equal(t, TEST_DEFAULT_ARMOR-3, attacker.GetArmor())

	// attacker takes 1 damage if target has thorns and attacking has shield
	setUp()
	attacker.AddAbility(ABILITY_SHIELD)
	game.MaybeApplyThorns(attacker, target, ATTACK_TYPE_MELEE)
	assert.Equal(t, TEST_DEFAULT_HEALTH, attacker.GetHealth())
	assert.Equal(t, TEST_DEFAULT_ARMOR-1, attacker.GetArmor())
}

func TestMaybeApplyMagicReflect(t *testing.T) {
	var game *Game
	var t1 *GameTeam
	var t2 *GameTeam
	var attacker *MonsterCard
	var target *MonsterCard

	setUp := func() {
		game, t1, t2 = CreateFakeGameAndTeams()
		attacker = t1.GetMonstersList()[1]
		target = t2.GetFirstAliveMonster()
		target.AddAbility(ABILITY_MAGIC_REFLECT)
	}

	// doesn't do anything if attack type is not magic
	setUp()
	attacker = t1.GetMonstersList()[0]
	game.MaybeApplyMagicReflect(attacker, target, ATTACK_TYPE_MELEE, 3)
	assert.Equal(t, TEST_DEFAULT_HEALTH, attacker.GetHealth())

	attacker = t1.GetMonstersList()[2]
	game.MaybeApplyMagicReflect(attacker, target, ATTACK_TYPE_RANGED, 3)
	assert.Equal(t, TEST_DEFAULT_HEALTH, attacker.GetHealth())

	// doesn't do anything if attack target has no magic reflect
	setUp()
	target.RemoveAbility(ABILITY_MAGIC_REFLECT)
	game.MaybeApplyMagicReflect(attacker, target, ATTACK_TYPE_MAGIC, 3)
	assert.Equal(t, TEST_DEFAULT_HEALTH, attacker.GetHealth())

	// doesn't do anything if attacker has reflection shield
	setUp()
	attacker.AddAbility(ABILITY_REFLECTION_SHIELD)
	game.MaybeApplyMagicReflect(attacker, target, ATTACK_TYPE_MAGIC, 3)
	assert.Equal(t, TEST_DEFAULT_HEALTH, attacker.GetHealth())

	// attacker takes 2 damage if target has magic reflect
	setUp()
	game.MaybeApplyMagicReflect(attacker, target, ATTACK_TYPE_MAGIC, 3)
	assert.Equal(t, TEST_DEFAULT_HEALTH-2, attacker.GetHealth())

	// attacker takes 3 damage if target has thorns and attacker has amplify
	setUp()
	attacker.AddDebuff(ABILITY_AMPLIFY)
	game.MaybeApplyMagicReflect(attacker, target, ATTACK_TYPE_MAGIC, 3)
	assert.Equal(t, TEST_DEFAULT_HEALTH-3, attacker.GetHealth())

	// attacker takes 1 damage if target has thorns and attacking has void
	setUp()
	attacker.AddAbility(ABILITY_VOID)
	game.MaybeApplyMagicReflect(attacker, target, ATTACK_TYPE_MAGIC, 3)
	assert.Equal(t, TEST_DEFAULT_HEALTH-1, attacker.GetHealth())

	// attacker takes 2 damage in void armor
	setUp()
	attacker.AddAbility(ABILITY_VOID_ARMOR)
	game.MaybeApplyMagicReflect(attacker, target, ATTACK_TYPE_MAGIC, 3)
	assert.Equal(t, TEST_DEFAULT_ARMOR-2, attacker.GetArmor())

	// attacker takes 1 damage in void armor if the attacker also has void
	setUp()
	attacker.AddAbility(ABILITY_VOID_ARMOR)
	attacker.AddAbility(ABILITY_VOID)
	game.MaybeApplyMagicReflect(attacker, target, ATTACK_TYPE_MAGIC, 3)
	assert.Equal(t, TEST_DEFAULT_ARMOR-1, attacker.GetArmor())
}

func TestMaybeApplyReturnFire(t *testing.T) {
	var game *Game
	var t1 *GameTeam
	var t2 *GameTeam
	var attacker *MonsterCard
	var target *MonsterCard

	setUp := func() {
		game, t1, t2 = CreateFakeGameAndTeams()
		attacker = t1.GetMonstersList()[2]
		target = t2.GetFirstAliveMonster()
		target.AddAbility(ABILITY_RETURN_FIRE)
	}

	// doesn't do anything if attack type is not range
	setUp()
	attacker = t1.GetMonstersList()[0]
	game.MaybeApplyReturnFire(attacker, target, ATTACK_TYPE_MELEE, 3)
	assert.Equal(t, TEST_DEFAULT_ARMOR, attacker.GetArmor())

	attacker = t1.GetMonstersList()[1]
	game.MaybeApplyReturnFire(attacker, target, ATTACK_TYPE_MAGIC, 3)
	assert.Equal(t, TEST_DEFAULT_ARMOR, attacker.GetArmor())

	// doesn't do anything if attack target has no return fire
	setUp()
	target.RemoveAbility(ABILITY_RETURN_FIRE)
	game.MaybeApplyReturnFire(attacker, target, ATTACK_TYPE_RANGED, 3)
	assert.Equal(t, TEST_DEFAULT_ARMOR, attacker.GetArmor())

	// doesn't do anything if attacker has reflection shield
	setUp()
	attacker.AddAbility(ABILITY_REFLECTION_SHIELD)
	game.MaybeApplyReturnFire(attacker, target, ATTACK_TYPE_RANGED, 3)
	assert.Equal(t, TEST_DEFAULT_ARMOR, attacker.GetArmor())

	// attacker takes 2 damage if target has return fire
	setUp()
	game.MaybeApplyReturnFire(attacker, target, ATTACK_TYPE_RANGED, 3)
	assert.Equal(t, TEST_DEFAULT_ARMOR-2, attacker.GetArmor())

	// attacker takes 2 damage in health if target has return fire
	setUp()
	attacker.SetArmor(0)
	game.MaybeApplyReturnFire(attacker, target, ATTACK_TYPE_RANGED, 3)
	assert.Equal(t, TEST_DEFAULT_HEALTH-2, attacker.GetHealth())

	// attacker takes 3 damage if target has return fire and attacker has amplify
	setUp()
	attacker.AddDebuff(ABILITY_AMPLIFY)
	game.MaybeApplyReturnFire(attacker, target, ATTACK_TYPE_RANGED, 3)
	assert.Equal(t, TEST_DEFAULT_ARMOR-3, attacker.GetArmor())

	// attacker takes 1 damage if target has return fire and attacking has shield
	setUp()
	attacker.AddAbility(ABILITY_SHIELD)
	game.MaybeApplyReturnFire(attacker, target, ATTACK_TYPE_RANGED, 3)
	assert.Equal(t, TEST_DEFAULT_ARMOR-1, attacker.GetArmor())
}

package simulator_tests

import (
	"log"
	"testing"

	. "github.com/YukiUmetsu/go-spl-simulator/game_models"
	"github.com/stretchr/testify/assert"
)

func TestDoRulesetPreGameBuff(t *testing.T) {

	createTeams := func() (*GameTeam, *GameTeam) {
		return CreateFakeGameTeam(), CreateFakeGameTeam()
	}

	// applies 2 armor in armored up ruleset
	t1, t2 := createTeams()
	DoRulesetPreGameBuff([]Ruleset{RULESET_ARMORED_UP}, t1, t2)
	t1Monsters := t1.GetAliveMonsters()
	assert.Equal(t, TEST_DEFAULT_ARMOR+2, t1Monsters[0].Armor)

	// removes all abilities in back to basics ruleset
	t1, t2 = createTeams()
	t1Monsters = t1.GetAliveMonsters()
	t2Monsters := t2.GetAliveMonsters()
	t1Monsters[0].AddAbility(ABILITY_FLYING)
	t2Monsters[0].AddAbility(ABILITY_BLAST)
	DoRulesetPreGameBuff([]Ruleset{RULESET_BACK_TO_BASICS}, t1, t2)
	t1Monsters = t1.GetAliveMonsters()
	t2Monsters = t2.GetAliveMonsters()
	assert.Equal(t, 0, len(t1Monsters[0].Abilities))
	assert.Equal(t, 0, len(t2Monsters[0].Abilities))

	// removes sneak/snipe ability of all monsters in fog of war ruleset
	t1, t2 = createTeams()
	t1Monsters = t1.GetAliveMonsters()
	t2Monsters = t2.GetAliveMonsters()
	t1Monsters[0].AddAbility(ABILITY_SNEAK)
	t2Monsters[0].AddAbility(ABILITY_SNIPE)
	DoRulesetPreGameBuff([]Ruleset{RULESET_FOG_OF_WAR}, t1, t2)
	t1Monsters = t1.GetAliveMonsters()
	t2Monsters = t2.GetAliveMonsters()
	assert.Equal(t, 0, len(t1Monsters[0].Abilities))
	assert.Equal(t, 0, len(t2Monsters[0].Abilities))

	// removes healing abilities of all monsters in healed out ruleset
	t1, t2 = createTeams()
	t1Monsters = t1.GetAliveMonsters()
	t1Monsters[0].AddAbility(ABILITY_HEAL)
	t1Monsters[1].AddAbility(ABILITY_TANK_HEAL)
	t1Monsters[2].AddAbility(ABILITY_TRIAGE)
	DoRulesetPreGameBuff([]Ruleset{RULESET_HEALED_OUT}, t1, t2)
	t1Monsters = t1.GetAliveMonsters()
	assert.Equal(t, 0, len(t1Monsters[0].Abilities))
	assert.Equal(t, 0, len(t1Monsters[1].Abilities))
	assert.Equal(t, 0, len(t1Monsters[2].Abilities))

	// gives summoner no stats in silenced summoner ruleset
	t1, t2 = createTeams()
	t1Summoner := t1.GetSummoner()
	t1Summoner.AddAbility(ABILITY_BLAST)
	DoRulesetPreGameBuff([]Ruleset{RULESET_SILENCED_SUMMONERS}, t1, t2)
	t1Summoner = t1.GetSummoner()
	assert.Equal(t, 0, len(t1Summoner.Abilities))

	// gives close range ability to all monsters in close range ruleset
	t1, t2 = createTeams()
	DoRulesetPreGameBuff([]Ruleset{RULESET_CLOSE_RANGE}, t1, t2)
	t1Monsters = t1.GetAliveMonsters()
	t2Monsters = t2.GetAliveMonsters()
	assert.Equal(t, Ability(ABILITY_CLOSE_RANGE), t1Monsters[0].Abilities[0])
	assert.Equal(t, Ability(ABILITY_CLOSE_RANGE), t2Monsters[0].Abilities[0])

	// gives opportunity ability to all monsters in equal opportunity ruleset
	t1, t2 = createTeams()
	DoRulesetPreGameBuff([]Ruleset{RULESET_EQUAL_OPPORTUNITY}, t1, t2)
	t1Monsters = t1.GetAliveMonsters()
	t2Monsters = t2.GetAliveMonsters()
	assert.Equal(t, Ability(ABILITY_OPPORTUNITY), t1Monsters[0].Abilities[0])
	assert.Equal(t, Ability(ABILITY_OPPORTUNITY), t2Monsters[0].Abilities[0])

	// gives blast ability to all monsters in explosive weaponry ruleset
	t1, t2 = createTeams()
	DoRulesetPreGameBuff([]Ruleset{RULESET_EXPLOSIVE_WEAPONRY}, t1, t2)
	t1Monsters = t1.GetAliveMonsters()
	t2Monsters = t2.GetAliveMonsters()
	assert.Equal(t, Ability(ABILITY_BLAST), t1Monsters[0].Abilities[0])
	assert.Equal(t, Ability(ABILITY_BLAST), t2Monsters[0].Abilities[0])

	// gives knock out ability to all monsters in heavy hitters ruleset
	t1, t2 = createTeams()
	DoRulesetPreGameBuff([]Ruleset{RULESET_HEAVY_HITTERS}, t1, t2)
	t1Monsters = t1.GetAliveMonsters()
	t2Monsters = t2.GetAliveMonsters()
	assert.Equal(t, Ability(ABILITY_KNOCK_OUT), t1Monsters[0].Abilities[0])
	assert.Equal(t, Ability(ABILITY_KNOCK_OUT), t2Monsters[0].Abilities[0])

	// gives custom melee mayhem ability to all monsters in melee mayhem ruleset
	t1, t2 = createTeams()
	DoRulesetPreGameBuff([]Ruleset{RULESET_MELEE_MAYHEM}, t1, t2)
	t1Monsters = t1.GetAliveMonsters()
	t2Monsters = t2.GetAliveMonsters()
	assert.Equal(t, Ability(ABILITY_MELEE_MAYHEM), t1Monsters[0].Abilities[0])
	assert.Equal(t, Ability(ABILITY_MELEE_MAYHEM), t2Monsters[0].Abilities[0])

	// gives enrage ability to all monsters in spreading fury ruleset
	t1, t2 = createTeams()
	DoRulesetPreGameBuff([]Ruleset{RULESET_SPREADING_FURY}, t1, t2)
	t1Monsters = t1.GetAliveMonsters()
	t2Monsters = t2.GetAliveMonsters()
	assert.Equal(t, Ability(ABILITY_ENRAGE), t1Monsters[0].Abilities[0])
	assert.Equal(t, Ability(ABILITY_ENRAGE), t2Monsters[0].Abilities[0])

	// gives sneak ability to all melee monsters in super sneak ruleset
	t1, t2 = createTeams()
	DoRulesetPreGameBuff([]Ruleset{RULESET_SUPER_SNEAK}, t1, t2)
	t1Monsters = t1.GetAliveMonsters()
	t2Monsters = t2.GetAliveMonsters()
	assert.Equal(t, Ability(ABILITY_SNEAK), t1Monsters[0].Abilities[0])
	assert.Equal(t, Ability(ABILITY_SNEAK), t2Monsters[0].Abilities[0])

	// gives void armor ability to all monsters in weak magic ruleset
	t1, t2 = createTeams()
	DoRulesetPreGameBuff([]Ruleset{RULESET_WEAK_MAGIC}, t1, t2)
	t1Monsters = t1.GetAliveMonsters()
	t2Monsters = t2.GetAliveMonsters()
	assert.Equal(t, Ability(ABILITY_VOID_ARMOR), t1Monsters[0].Abilities[0])
	assert.Equal(t, Ability(ABILITY_VOID_ARMOR), t2Monsters[0].Abilities[0])

	// gives snipe ability to all monsters that have ranged/magic in target practice ruleset
	t1, t2 = createTeams()
	DoRulesetPreGameBuff([]Ruleset{RULESET_TARGET_PRACTICE}, t1, t2)
	t1Monsters = t1.GetAliveMonsters()
	assert.Equal(t, Ability(ABILITY_SNIPE), t1Monsters[1].Abilities[0])
	assert.Equal(t, Ability(ABILITY_SNIPE), t1Monsters[2].Abilities[0])

	// sets health of all monsters to highest health monster in equalizer ruleset
	t1, t2 = createTeams()
	t1Monsters = t1.GetAliveMonsters()
	t2Monsters = t2.GetAliveMonsters()
	t1Monsters[0].SetHealth(15)
	DoRulesetPreGameBuff([]Ruleset{RULESET_EQUALIZER}, t1, t2)
	t2Monsters = t2.GetAliveMonsters()
	assert.Equal(t, 15, t2Monsters[0].Health)

	// give poison in noxious fumes ruleset
	t1, t2 = createTeams()
	DoRulesetPreGameBuff([]Ruleset{RULESET_NOXIOUS_FUMES}, t1, t2)
	t1Monsters = t1.GetAliveMonsters()
	t2Monsters = t2.GetAliveMonsters()
	t1mDebuffMap := t1Monsters[0].GetDebuffs()
	t2mDebuffMap := t2Monsters[0].GetDebuffs()
	if _, ok := t1mDebuffMap[Ability(ABILITY_POISON)]; !ok {
		log.Fatalln("poison was not given in noxious fumes ruleset")
	}
	if _, ok := t2mDebuffMap[Ability(ABILITY_POISON)]; !ok {
		log.Fatalln("poison was not given in noxious fumes ruleset")
	}
}

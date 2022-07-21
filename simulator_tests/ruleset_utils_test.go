package simulator_tests

import (
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
}

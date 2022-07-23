package simulator_tests

import (
	"testing"

	. "github.com/YukiUmetsu/go-spl-simulator/game_models"
	"github.com/stretchr/testify/assert"
)

func getFakeTeam() (*SummonerCard, []*MonsterCard, *GameTeam) {
	var fakeSummoner *SummonerCard
	fakeCards := make([]*MonsterCard, 0)
	var gameTeam GameTeam
	fakeSummoner = GetDefaultFakeSummoner()
	fakeCards = append(fakeCards, GetDefaultFakeMonster(ATTACK_TYPE_MELEE))
	fakeCards = append(fakeCards, GetDefaultFakeMonster(ATTACK_TYPE_MAGIC))
	fakeCards = append(fakeCards, GetDefaultFakeMonster(ATTACK_TYPE_RANGED))
	gameTeam.Create(fakeSummoner, fakeCards, "testuser")
	return fakeSummoner, fakeCards, &gameTeam
}

func TestSetTeamNumber(t *testing.T) {
	fakeSummoner, _, gameTeam := getFakeTeam()
	gameTeam.SetTeamNumber(TEAM_NUM_TWO)
	assert.Equal(t, TEAM_NUM_TWO, fakeSummoner.GetTeamNumber())
}

func TestGetMonsterPosition(t *testing.T) {
	_, fakeCards, gameTeam := getFakeTeam()
	assert.Equal(t, 1, gameTeam.GetMonsterPosition(fakeCards[1]))
	fakeCards[0].HitHealth(100)
	assert.Equal(t, 0, gameTeam.GetMonsterPosition(fakeCards[1]))
}

func TestGetSummoner(t *testing.T) {
	fakeSummoner, _, gameTeam := getFakeTeam()
	assert.Equal(t, fakeSummoner, gameTeam.GetSummoner())
}

func TestGetMonstersList(t *testing.T) {
	_, fakeCards, gameTeam := getFakeTeam()
	assert.Equal(t, len(fakeCards), len(gameTeam.GetMonstersList()))
}

func TestGetAliveMonsters(t *testing.T) {
	_, fakeCards, gameTeam := getFakeTeam()
	assert.Equal(t, len(fakeCards), len(gameTeam.GetAliveMonsters()))
	gameTeam.GetMonstersList()[0].HitHealth(100)
	assert.Equal(t, len(fakeCards)-1, len(gameTeam.GetAliveMonsters()))
}

func TestGetUnmovedMonsters(t *testing.T) {
	_, fakeCards, gameTeam := getFakeTeam()
	assert.Equal(t, len(fakeCards), len(gameTeam.GetUnmovedMonsters()))
	gameTeam.GetMonstersList()[0].SetHasTurnPassed(true)
	assert.Equal(t, len(fakeCards)-1, len(gameTeam.GetUnmovedMonsters()))
}

func TestGetFirstAliveMonster(t *testing.T) {
	_, fakeCards, gameTeam := getFakeTeam()
	assert.Equal(t, fakeCards[0].GetCardDetail().ID, gameTeam.GetFirstAliveMonster().GetCardDetail().ID)
	for i, m := range gameTeam.GetMonstersList() {
		if i != 2 {
			m.HitHealth(100)
		}
	}
	assert.Equal(t, fakeCards[2].GetCardDetail().ID, gameTeam.GetFirstAliveMonster().GetCardDetail().ID)
}

func TestMaybeSetLastStand(t *testing.T) {
	_, _, gameTeam := getFakeTeam()

	// kill the 3rd position monster
	gameTeam.GetMonstersList()[2].HitHealth(100)
	gameTeam.MaybeSetLastStand()
	for _, m := range gameTeam.GetMonstersList() {
		assert.False(t, m.IsLastMonster())
	}

	// kill the 2nd position monster
	gameTeam.GetMonstersList()[1].HitHealth(100)

	// check the possible last monster
	gameTeam.MaybeSetLastStand()
	for i, m := range gameTeam.GetMonstersList() {
		if i == 0 {
			assert.True(t, m.IsLastMonster())
		} else {
			assert.False(t, m.IsLastMonster())
		}
	}
}

func TestSetAllMonsterHealth(t *testing.T) {
	// sets monsters health to maximum
	_, _, gameTeam := getFakeTeam()
	gameTeam.GetMonstersList()[0].AddBuff(ABILITY_STRENGTHEN)
	gameTeam.GetMonstersList()[0].SetHealth(0)
	gameTeam.SetAllMonsterHealth()
	assert.Equal(t, TEST_DEFAULT_HEALTH+1, gameTeam.GetMonstersList()[0].GetHealth())
}

func TestGetScattershotTarget(t *testing.T) {
	_, _, gameTeam := getFakeTeam()
	targetMap := make(map[int]int)
	for i := 0; i < 100; i++ {
		target := gameTeam.GetScattershotTarget()
		targetID := target.GetCardDetail().ID
		if _, ok := targetMap[targetID]; !ok {
			targetMap[targetID] = targetID
		}
	}
	assert.True(t, len(targetMap) > 1)
}

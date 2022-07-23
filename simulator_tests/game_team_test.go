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

func TestResetTeam(t *testing.T) {
	// resets the card stats
	_, _, gameTeam := getFakeTeam()
	gameTeam.GetMonstersList()[0].SetHealth(0)
	gameTeam.GetMonstersList()[0].SetArmor(0)
	gameTeam.GetMonstersList()[1].AddAbility(ABILITY_LAST_STAND)
	gameTeam.GetMonstersList()[1].SetIsOnlyMonster()
	gameTeam.SetTeamNumber(TEAM_NUM_ONE)
	gameTeam.ResetTeam()
	assert.Equal(t, 3, len(gameTeam.GetAliveMonsters()))
	assert.Equal(t, 5, gameTeam.GetAliveMonsters()[0].GetHealth())

	// set only monster if there's only one monster
	fakeSummoner, _, gameTeam := getFakeTeam()
	fakeMonster := GetDefaultFakeMonster(ATTACK_TYPE_MAGIC)
	gameTeam.Create(fakeSummoner, []*MonsterCard{fakeMonster}, gameTeam.GetPlayerName())
	gameTeam.SetTeamNumber(TEAM_NUM_ONE)
	assert.False(t, fakeMonster.IsLastMonster())
	gameTeam.ResetTeam()
	assert.True(t, gameTeam.GetMonstersList()[0].IsLastMonster())
}

func TestGetSnipeTarget(t *testing.T) {
	// always returns taunt monster if there is one
	_, fakeCards, gameTeam := getFakeTeam()
	gameTeam.GetMonstersList()[2].AddAbility(ABILITY_TAUNT)
	assert.Equal(t, fakeCards[2].GetCardDetail().ID, gameTeam.GetSnipeTarget().GetCardDetail().ID)
	gameTeam.GetMonstersList()[2].RemoveAbility(ABILITY_TAUNT)
	gameTeam.GetMonstersList()[1].AddAbility(ABILITY_TAUNT)
	assert.Equal(t, fakeCards[1].GetCardDetail().ID, gameTeam.GetSnipeTarget().GetCardDetail().ID)

	// normal snipe target
	_, fakeCards, gameTeam = getFakeTeam()
	assert.Equal(t, fakeCards[1].GetCardDetail().ID, gameTeam.GetSnipeTarget().GetCardDetail().ID)

	// don't hit camoflage monster
	_, fakeCards, gameTeam = getFakeTeam()
	gameTeam.GetMonstersList()[1].AddAbility(ABILITY_CAMOUFLAGE)
	assert.Equal(t, fakeCards[2].GetCardDetail().ID, gameTeam.GetSnipeTarget().GetCardDetail().ID)

	// returns first monster if no snipe targets
	_, fakeCards, gameTeam = getFakeTeam()
	gameTeam.GetMonstersList()[1].AddAbility(ABILITY_CAMOUFLAGE)
	gameTeam.GetMonstersList()[2].AddAbility(ABILITY_CAMOUFLAGE)
	assert.Equal(t, fakeCards[0].GetCardDetail().ID, gameTeam.GetSnipeTarget().GetCardDetail().ID)
}

func TestGetOpportunityTarget(t *testing.T) {
	// always returns taunt monster if there is one
	_, fakeCards, gameTeam := getFakeTeam()
	gameTeam.GetMonstersList()[2].AddAbility(ABILITY_TAUNT)
	assert.Equal(t, fakeCards[2].GetCardDetail().ID, gameTeam.GetOpportunityTarget().GetCardDetail().ID)
	gameTeam.GetMonstersList()[2].RemoveAbility(ABILITY_TAUNT)
	gameTeam.GetMonstersList()[1].AddAbility(ABILITY_TAUNT)
	assert.Equal(t, fakeCards[1].GetCardDetail().ID, gameTeam.GetOpportunityTarget().GetCardDetail().ID)

	// returns the lowest health monster
	_, fakeCards, gameTeam = getFakeTeam()
	gameTeam.GetMonstersList()[1].SetHealth(1)
	assert.Equal(t, fakeCards[1].GetCardDetail().ID, gameTeam.GetOpportunityTarget().GetCardDetail().ID)

	// returns the correct target if original target has camoflauge
	_, fakeCards, gameTeam = getFakeTeam()
	gameTeam.GetMonstersList()[1].SetHealth(1)
	gameTeam.GetMonstersList()[2].SetHealth(2)
	gameTeam.GetMonstersList()[1].AddAbility(ABILITY_CAMOUFLAGE)
	assert.Equal(t, fakeCards[2].GetCardDetail().ID, gameTeam.GetOpportunityTarget().GetCardDetail().ID)

	// returns first monster if no opportunity targets
	_, fakeCards, gameTeam = getFakeTeam()
	gameTeam.GetMonstersList()[1].SetHealth(1)
	gameTeam.GetMonstersList()[2].SetHealth(2)
	gameTeam.GetMonstersList()[1].AddAbility(ABILITY_CAMOUFLAGE)
	gameTeam.GetMonstersList()[2].AddAbility(ABILITY_CAMOUFLAGE)
	assert.Equal(t, fakeCards[0].GetCardDetail().ID, gameTeam.GetOpportunityTarget().GetCardDetail().ID)
}

func TestGetSneakTarget(t *testing.T) {
	// always returns taunt monster if there is one
	_, fakeCards, gameTeam := getFakeTeam()
	gameTeam.GetMonstersList()[2].AddAbility(ABILITY_TAUNT)
	assert.Equal(t, fakeCards[2].GetCardDetail().ID, gameTeam.GetSneakTarget().GetCardDetail().ID)
	gameTeam.GetMonstersList()[2].RemoveAbility(ABILITY_TAUNT)
	gameTeam.GetMonstersList()[1].AddAbility(ABILITY_TAUNT)
	assert.Equal(t, fakeCards[1].GetCardDetail().ID, gameTeam.GetSneakTarget().GetCardDetail().ID)

	// get the last monster
	_, fakeCards, gameTeam = getFakeTeam()
	assert.Equal(t, fakeCards[2].GetCardDetail().ID, gameTeam.GetSneakTarget().GetCardDetail().ID)

	// returns the correct target if original target has camoflauge
	_, fakeCards, gameTeam = getFakeTeam()
	gameTeam.GetMonstersList()[2].AddAbility(ABILITY_CAMOUFLAGE)
	assert.Equal(t, fakeCards[1].GetCardDetail().ID, gameTeam.GetSneakTarget().GetCardDetail().ID)

	// returns first monster if no sneak targets
	_, fakeCards, gameTeam = getFakeTeam()
	gameTeam.GetMonstersList()[1].AddAbility(ABILITY_CAMOUFLAGE)
	gameTeam.GetMonstersList()[2].AddAbility(ABILITY_CAMOUFLAGE)
	assert.Equal(t, fakeCards[0].GetCardDetail().ID, gameTeam.GetSneakTarget().GetCardDetail().ID)
}

func TestGetRepairTarget(t *testing.T) {
	// returns null if no armor is lost
	var monsterNilPointer *MonsterCard
	_, fakeCards, gameTeam := getFakeTeam()
	assert.Equal(t, gameTeam.GetRepairTarget(), monsterNilPointer)

	// returns the monster that lost the most armor
	gameTeam.GetMonstersList()[1].SetArmor(3)
	assert.Equal(t, gameTeam.GetRepairTarget().GetCardDetail().ID, fakeCards[1].GetCardDetail().ID)
	gameTeam.GetMonstersList()[0].SetArmor(0)
	assert.Equal(t, gameTeam.GetRepairTarget().GetCardDetail().ID, fakeCards[0].GetCardDetail().ID)

	// doesn't repair the dead monster's armor
	_, fakeCards, gameTeam = getFakeTeam()
	gameTeam.GetMonstersList()[1].SetArmor(1)
	gameTeam.GetMonstersList()[1].SetHealth(0)
	gameTeam.GetMonstersList()[0].SetArmor(3)
	assert.Equal(t, gameTeam.GetRepairTarget().GetCardDetail().ID, fakeCards[0].GetCardDetail().ID)
}

func TestGetTriageHealTarget(t *testing.T) {
	// returns null if no hp is lost
	var monsterNilPointer *MonsterCard
	_, fakeCards, gameTeam := getFakeTeam()
	assert.Equal(t, gameTeam.GetTriageHealTarget(), monsterNilPointer)

	// returns the backline monster that lost the most health
	gameTeam.GetMonstersList()[1].SetHealth(3)
	assert.Equal(t, gameTeam.GetTriageHealTarget().GetCardDetail().ID, fakeCards[1].GetCardDetail().ID)
	gameTeam.GetMonstersList()[0].SetHealth(2)
	assert.Equal(t, gameTeam.GetTriageHealTarget().GetCardDetail().ID, fakeCards[1].GetCardDetail().ID)
	gameTeam.GetMonstersList()[2].SetHealth(1)
	assert.Equal(t, gameTeam.GetTriageHealTarget().GetCardDetail().ID, fakeCards[2].GetCardDetail().ID)

	// 	doesn't try to heal dead mosnter
	_, fakeCards, gameTeam = getFakeTeam()
	gameTeam.GetMonstersList()[1].SetHealth(0)
	gameTeam.GetMonstersList()[2].SetHealth(1)
	assert.Equal(t, gameTeam.GetTriageHealTarget().GetCardDetail().ID, fakeCards[2].GetCardDetail().ID)
}

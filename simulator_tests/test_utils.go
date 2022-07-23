package simulator_tests

import (
	"math/rand"
	"time"

	. "github.com/YukiUmetsu/go-spl-simulator/game_models"
)

const TEST_DEFAULT_MANA = 5
const TEST_DEFAULT_HEALTH = 5

const TEST_DEFAULT_SPEED = 5
const TEST_DEFAULT_ARMOR = 5
const TEST_DEFAULT_ATTACK = 5
const TEST_DEFAULT_RANGED = 5
const TEST_DEFAULT_MAGIC = 5

func GetDefaultFakeSummoner() *SummonerCard {
	abilities := []any{}
	stats := CardRawStats{
		Abilities: abilities,
		Mana:      float64(TEST_DEFAULT_MANA),
		Health:    float64(0),
		Speed:     float64(0),
		Armor:     float64(0),
		Attack:    float64(0),
		Ranged:    float64(0),
		Magic:     float64(0),
	}
	details := CreateFakeCardDetail(SUMMONER, stats)
	var fakeSummoner SummonerCard
	(&fakeSummoner).Setup(details, 4)
	return &fakeSummoner
}

func GetDefaultFakeMonster(attackType CardAttackType) *MonsterCard {
	var cardDetail CardDetail
	if attackType == ATTACK_TYPE_MELEE {
		cardDetail = GetDefaultFakeMeleeOnlyCardDetail()
	} else if attackType == ATTACK_TYPE_RANGED {
		cardDetail = GetDefaultFakeRangeOnlyCardDetail()
	} else if attackType == ATTACK_TYPE_MAGIC {
		cardDetail = GetDefaultFakeMagicOnlyCardDetail()
	} else {
		cardDetail = GetDefaultFakeNoAttackCardDetail()
	}

	var m MonsterCard
	(&m).Setup(cardDetail, 4)
	m.Team = TEAM_NUM_ONE
	return &m
}

func GetDefaultFakeMonsterWithAbility(attackType CardAttackType, abilities []Ability) *MonsterCard {
	m := GetDefaultFakeMonster(attackType)
	for _, ability := range abilities {
		m.AddAbility(ability)
	}
	return m
}

func CreateFakeCardDetail(cardType CardType, stats CardRawStats) CardDetail {
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 10000
	randomNumber := rand.Intn(max-min+1) + min

	return CardDetail{
		ID:        randomNumber,
		Color:     COLOR_BLACK,
		Type:      cardType,
		Rarity:    1,
		IsStarter: false,
		Editions:  "1",
		Stats:     stats,
	}
}

func GetDefaultFakeMagicOnlyCardDetail() CardDetail {
	stats := CardRawStats{
		Mana:   []any{TEST_DEFAULT_MANA, TEST_DEFAULT_MANA, TEST_DEFAULT_MANA, TEST_DEFAULT_MANA, TEST_DEFAULT_MANA, TEST_DEFAULT_MANA, TEST_DEFAULT_MANA, TEST_DEFAULT_MANA},
		Health: []any{TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH},
		Speed:  []any{TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED},
		Armor:  []any{TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR},
		Magic:  []any{TEST_DEFAULT_MAGIC, TEST_DEFAULT_MAGIC, TEST_DEFAULT_MAGIC, TEST_DEFAULT_MAGIC, TEST_DEFAULT_MAGIC, TEST_DEFAULT_MAGIC, TEST_DEFAULT_MAGIC, TEST_DEFAULT_MAGIC},
		Ranged: []any{0, 0, 0, 0, 0, 0, 0, 0},
		Attack: []any{0, 0, 0, 0, 0, 0, 0, 0},
	}
	return CreateFakeCardDetail(MONSTER, stats)
}

func GetDefaultFakeRangeOnlyCardDetail() CardDetail {
	stats := CardRawStats{
		Mana:   []any{TEST_DEFAULT_MANA, TEST_DEFAULT_MANA, TEST_DEFAULT_MANA, TEST_DEFAULT_MANA, TEST_DEFAULT_MANA, TEST_DEFAULT_MANA, TEST_DEFAULT_MANA, TEST_DEFAULT_MANA},
		Health: []any{TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH},
		Speed:  []any{TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED},
		Armor:  []any{TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR},
		Ranged: []any{TEST_DEFAULT_RANGED, TEST_DEFAULT_RANGED, TEST_DEFAULT_RANGED, TEST_DEFAULT_RANGED, TEST_DEFAULT_RANGED, TEST_DEFAULT_RANGED, TEST_DEFAULT_RANGED, TEST_DEFAULT_RANGED},
		Magic:  []any{0, 0, 0, 0, 0, 0, 0, 0},
		Attack: []any{0, 0, 0, 0, 0, 0, 0, 0},
	}
	return CreateFakeCardDetail(MONSTER, stats)
}

func GetDefaultFakeMeleeOnlyCardDetail() CardDetail {
	stats := CardRawStats{
		Mana:   []any{TEST_DEFAULT_MANA, TEST_DEFAULT_MANA, TEST_DEFAULT_MANA, TEST_DEFAULT_MANA, TEST_DEFAULT_MANA, TEST_DEFAULT_MANA, TEST_DEFAULT_MANA, TEST_DEFAULT_MANA},
		Health: []any{TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH},
		Speed:  []any{TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED},
		Armor:  []any{TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR},
		Attack: []any{TEST_DEFAULT_ATTACK, TEST_DEFAULT_ATTACK, TEST_DEFAULT_ATTACK, TEST_DEFAULT_ATTACK, TEST_DEFAULT_ATTACK, TEST_DEFAULT_ATTACK, TEST_DEFAULT_ATTACK, TEST_DEFAULT_ATTACK},
		Ranged: []any{0, 0, 0, 0, 0, 0, 0, 0},
		Magic:  []any{0, 0, 0, 0, 0, 0, 0, 0},
	}
	return CreateFakeCardDetail(MONSTER, stats)
}

func GetDefaultFakeNoAttackCardDetail() CardDetail {
	stats := CardRawStats{
		Mana:   []any{TEST_DEFAULT_MANA, TEST_DEFAULT_MANA, TEST_DEFAULT_MANA, TEST_DEFAULT_MANA, TEST_DEFAULT_MANA, TEST_DEFAULT_MANA, TEST_DEFAULT_MANA, TEST_DEFAULT_MANA},
		Health: []any{TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH, TEST_DEFAULT_HEALTH},
		Speed:  []any{TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED, TEST_DEFAULT_SPEED},
		Armor:  []any{TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR, TEST_DEFAULT_ARMOR},
		Attack: []any{0, 0, 0, 0, 0, 0, 0, 0},
		Magic:  []any{0, 0, 0, 0, 0, 0, 0, 0},
		Ranged: []any{0, 0, 0, 0, 0, 0, 0, 0},
	}
	return CreateFakeCardDetail(MONSTER, stats)
}

func CreateMonsterOfRarityAndLevel(rarity, level int, cardType CardType, attackType CardAttackType) *MonsterCard {
	var m *MonsterCard
	m = GetDefaultFakeMonster(attackType)
	cardDetail := m.GetCardDetail()
	cardDetail.Rarity = rarity
	m.SetCardDetail(cardDetail)
	m.CardLevel = level
	return m
}

func CreateSummonerOfRarityAndLevel(rarity, level int, cardType CardType) *SummonerCard {
	var s *SummonerCard
	s = GetDefaultFakeSummoner()
	cardDetail := s.GetCardDetail()
	cardDetail.Rarity = rarity
	s.SetCardDetail(cardDetail)
	s.CardLevel = level
	return s
}

func CreateFakeGameTeam() *GameTeam {
	var t GameTeam
	ml := make([]*MonsterCard, 0)
	s := GetDefaultFakeSummoner()
	m1 := GetDefaultFakeMonster(ATTACK_TYPE_MELEE)
	m2 := GetDefaultFakeMonster(ATTACK_TYPE_MAGIC)
	m3 := GetDefaultFakeMonster(ATTACK_TYPE_RANGED)
	ml = append(ml, m1, m2, m3)
	(&t).Create(s, ml, "test_player1")
	return &t
}

func CreateFakeGame() *Game {
	var game Game
	t1 := CreateFakeGameTeam()
	t2 := CreateFakeGameTeam()
	rulesets := make([]Ruleset, 0)
	rulesets = append(rulesets, RULESET_EQUAL_OPPORTUNITY)
	(&game).Create(t1, t2, rulesets, false)
	return &game
}

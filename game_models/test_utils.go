package game_models

func GetDefaultFakeSummoner() *SummonerCard {
	abilities := []any{}
	stats := CardRawStats{
		Abilities: abilities,
		Mana:      5,
		Health:    0,
		Speed:     0,
		Armor:     0,
		Attack:    0,
		Ranged:    0,
		Magic:     0,
	}
	details := CreateFakeCardDetail(SUMMONER, stats)
	var fakeSummoner *SummonerCard
	fakeSummoner.Setup(details, 4)
	return fakeSummoner
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

func CreateFakeCardDetail(cardType CardType, stats CardRawStats) CardDetail {
	return CardDetail{
		ID:        9999,
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
		Mana:   []any{5, 5, 5, 5, 5, 5, 5, 5},
		Health: []any{5, 5, 5, 5, 5, 5, 5, 5},
		Speed:  []any{5, 5, 5, 5, 5, 5, 5, 5},
		Armor:  []any{5, 5, 5, 5, 5, 5, 5, 5},
		Magic:  []any{5, 5, 5, 5, 5, 5, 5, 5},
		Ranged: []any{0, 0, 0, 0, 0, 0, 0, 0},
		Attack: []any{0, 0, 0, 0, 0, 0, 0, 0},
	}
	return CreateFakeCardDetail(MONSTER, stats)
}

func GetDefaultFakeRangeOnlyCardDetail() CardDetail {
	stats := CardRawStats{
		Mana:   []any{5, 5, 5, 5, 5, 5, 5, 5},
		Health: []any{5, 5, 5, 5, 5, 5, 5, 5},
		Speed:  []any{5, 5, 5, 5, 5, 5, 5, 5},
		Armor:  []any{5, 5, 5, 5, 5, 5, 5, 5},
		Ranged: []any{5, 5, 5, 5, 5, 5, 5, 5},
		Magic:  []any{0, 0, 0, 0, 0, 0, 0, 0},
		Attack: []any{0, 0, 0, 0, 0, 0, 0, 0},
	}
	return CreateFakeCardDetail(MONSTER, stats)
}

func GetDefaultFakeMeleeOnlyCardDetail() CardDetail {
	stats := CardRawStats{
		Mana:   []any{5, 5, 5, 5, 5, 5, 5, 5},
		Health: []any{5, 5, 5, 5, 5, 5, 5, 5},
		Speed:  []any{5, 5, 5, 5, 5, 5, 5, 5},
		Armor:  []any{5, 5, 5, 5, 5, 5, 5, 5},
		Attack: []any{5, 5, 5, 5, 5, 5, 5, 5},
		Ranged: []any{0, 0, 0, 0, 0, 0, 0, 0},
		Magic:  []any{0, 0, 0, 0, 0, 0, 0, 0},
	}
	return CreateFakeCardDetail(MONSTER, stats)
}

func GetDefaultFakeNoAttackCardDetail() CardDetail {
	stats := CardRawStats{
		Mana:   5,
		Health: 5,
		Speed:  5,
	}
	return CreateFakeCardDetail(MONSTER, stats)
}

func CreateMonsterOfRarityAndLevel(rarity, level int, cardType CardType, attackType CardAttackType) *MonsterCard {
	var m *MonsterCard
	m = GetDefaultFakeMonster(attackType)
	m.cardDetail.Rarity = rarity
	m.CardLevel = level
	return m
}

func CreateSummonerOfRarityAndLevel(rarity, level int, cardType CardType) *SummonerCard {
	var s *SummonerCard
	s = GetDefaultFakeSummoner()
	s.cardDetail.Rarity = rarity
	s.CardLevel = level
	return s
}

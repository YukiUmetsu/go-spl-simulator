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

	var m *MonsterCard
	m.Setup(cardDetail, 4)
	return m
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
		Mana:   5,
		Health: 5,
		Speed:  5,
		Armor:  5,
		Magic:  5,
	}
	return CreateFakeCardDetail(MONSTER, stats)
}

func GetDefaultFakeRangeOnlyCardDetail() CardDetail {
	stats := CardRawStats{
		Mana:   5,
		Health: 5,
		Speed:  5,
		Armor:  5,
		Ranged: 5,
	}
	return CreateFakeCardDetail(MONSTER, stats)
}

func GetDefaultFakeMeleeOnlyCardDetail() CardDetail {
	stats := CardRawStats{
		Mana:   5,
		Health: 5,
		Speed:  5,
		Armor:  5,
		Attack: 5,
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

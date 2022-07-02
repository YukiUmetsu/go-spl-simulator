package main

import (
	utilGeneral "utils/generals"
)

type SummonerCard struct {
	cardDetail     SummonerCardDetail
	cardLevel      int
	team           TeamNumber
	DebuffMap      map[Ability]int
	BuffMap        map[Ability]int
	Abilities      []Ability
	Speed          int
	StartingArmor  int
	Armor          int
	Health         int
	StartingHealth int
	Magic          int
	Melee          int
	Ranged         int
	Mana           int
}

func (c SummonerCard) Setup(cardDetail SummonerCardDetail, cardLevel int) {
	c.cardDetail = cardDetail
	c.cardLevel = cardLevel - 1
	c.SetStats(c.cardDetail.Stats)
}

func (c SummonerCard) SetTeam(teamNumber TeamNumber) {
	c.team = teamNumber
}

func (c SummonerCard) SetStats(stats FlatCardStats) {
	c.Speed = stats.Speed
	c.Armor = stats.Armor
	c.StartingArmor = stats.Armor
	c.Health = stats.Health
	c.StartingHealth = stats.Health
	c.Magic = stats.Magic
	c.Ranged = stats.Ranged
	c.Melee = stats.Attack
	c.Mana = stats.Mana
	c.AddAbilities(stats.Abilities)
}

func (c SummonerCard) GetStat(stat int) FlatCardStats {
	return c.cardDetail.Stats
}

func (c SummonerCard) AddAbilities(abilities []Ability) {
	for _, ability := range abilities {
		c.Abilities = append(c.Abilities, ability)
	}
}

func (c SummonerCard) GetCardDetail() SummonerCardDetail {
	return c.cardDetail
}

func (c SummonerCard) HasAbility(ability Ability) bool {
	return c.HasAbility(ability)
}

func (c SummonerCard) RemoveAbility(ability Ability) {
	c.Abilities = utilGeneral.Remove(c.Abilities, ability)
}

func (c SummonerCard) GetTeamNumber() TeamNumber {
	return c.team
}

func (c SummonerCard) GetRarity() int {
	return c.cardDetail.Rarity
}

func (c SummonerCard) GetName() string {
	return c.cardDetail.Name
}

func (c SummonerCard) GetLevel() int {
	return c.cardLevel
}

func (c SummonerCard) GetDebuffs() map[Ability]int {
	return c.DebuffMap
}

func (c SummonerCard) GetBuffs() map[Ability]int {
	return c.BuffMap
}

func (c SummonerCard) Clone() SummonerCard {
	clonedCard := SummonerCard{
		cardDetail:     c.cardDetail,
		cardLevel:      c.cardLevel,
		team:           c.team,
		DebuffMap:      c.DebuffMap,
		BuffMap:        c.BuffMap,
		Abilities:      c.Abilities,
		Speed:          c.Speed,
		StartingArmor:  c.StartingArmor,
		Armor:          c.Armor,
		StartingHealth: c.StartingHealth,
		Health:         c.Health,
		Magic:          c.Magic,
		Melee:          c.Melee,
		Ranged:         c.Ranged,
		Mana:           c.Mana,
	}
	clonedCard.SetTeam(c.GetTeamNumber())
	return clonedCard
}

func (c SummonerCard) AddAbilitiesWithArray(abilities []Ability) {
	for _, a := range abilities {
		c.Abilities = append(c.Abilities, a)
	}
}

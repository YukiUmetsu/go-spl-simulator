package main

import (
	utilGeneral "utils/generals"
)

type MonsterCard struct {
	cardDetail     MonsterCardDetail
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

func (c MonsterCard) Setup(cardDetail MonsterCardDetail, cardLevel int) {
	c.cardDetail = cardDetail
	c.cardLevel = cardLevel - 1
	c.SetStats(c.cardDetail.Stats)
}

func (c MonsterCard) SetTeam(teamNumber TeamNumber) {
	c.team = teamNumber
}

func (c MonsterCard) SetStats(stats CardStatsByLevel) {
	c.Speed = c.GetStat(stats.Speed)
	c.Armor = c.GetStat(stats.Armor)
	c.StartingArmor = c.GetStat(stats.Armor)
	c.Health = c.GetStat(stats.Health)
	c.StartingHealth = c.GetStat(stats.Health)
	c.Magic = c.GetStat(stats.Magic)
	c.Ranged = c.GetStat(stats.Ranged)
	c.Melee = c.GetStat(stats.Attack)
	c.Mana = c.GetStat(stats.Mana)
	c.AddAbilities(stats.Abilities)
}

func (c MonsterCard) GetStat(stats []int) int {
	return stats[c.cardLevel]
}

func (c MonsterCard) AddAbilities(abilitiesArray [][]Ability) {
	for i, abilities := range abilitiesArray {
		if i+1 <= c.cardLevel {
			for _, ability := range abilities {
				c.Abilities = append(c.Abilities, ability)
			}
		}
	}
}

func (c MonsterCard) GetCardDetail() MonsterCardDetail {
	return c.cardDetail
}

func (c MonsterCard) HasAbility(ability Ability) bool {
	return c.HasAbility(ability)
}

func (c MonsterCard) RemoveAbility(ability Ability) {
	c.Abilities = utilGeneral.Remove(c.Abilities, ability)
}

func (c MonsterCard) GetTeamNumber() TeamNumber {
	return c.team
}

func (c MonsterCard) GetRarity() int {
	return c.cardDetail.Rarity
}

func (c MonsterCard) GetName() string {
	return c.cardDetail.Name
}

func (c MonsterCard) GetLevel() int {
	return c.cardLevel
}

func (c MonsterCard) GetDebuffs() map[Ability]int {
	return c.DebuffMap
}

func (c MonsterCard) GetBuffs() map[Ability]int {
	return c.BuffMap
}

func (c MonsterCard) Clone() MonsterCard {
	clonedCard := MonsterCard{
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

func (c MonsterCard) AddAbilitiesWithArray(abilities []Ability) {
	for _, a := range abilities {
		c.Abilities = append(c.Abilities, a)
	}
}

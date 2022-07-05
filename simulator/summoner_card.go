package simulator

import (
	utils "game_utils"
)

type SummonerCard struct {
	GameCard
	cardDetail SummonerCardDetail
}

func (c *SummonerCard) Setup(cardDetail SummonerCardDetail, cardLevel int) {
	c.cardDetail = cardDetail
	c.CardLevel = cardLevel - 1
	c.SetStats(c.cardDetail.Stats)
}

func (c *SummonerCard) SetTeam(teamNumber TeamNumber) {
	c.Team = teamNumber
}

func (c *SummonerCard) SetStats(stats FlatCardStats) {
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

func (c *SummonerCard) GetStat(stat int) FlatCardStats {
	return c.cardDetail.Stats
}

func (c *SummonerCard) AddAbilities(abilities []Ability) {
	for _, ability := range abilities {
		c.Abilities = append(c.Abilities, ability)
	}
}

func (c *SummonerCard) GetCardDetail() SummonerCardDetail {
	return c.cardDetail
}

func (c *SummonerCard) HasAbility(ability Ability) bool {
	return c.HasAbility(ability)
}

func (c *SummonerCard) RemoveAbility(ability Ability) {
	c.Abilities = utils.Remove(c.Abilities, ability)
}

func (c *SummonerCard) RemoveAllAbilities() {
	c.Abilities = []Ability{}
}

func (c *SummonerCard) GetTeamNumber() TeamNumber {
	return c.Team
}

func (c *SummonerCard) GetRarity() int {
	return c.cardDetail.Rarity
}

func (c *SummonerCard) GetName() string {
	return c.cardDetail.Name
}

func (c *SummonerCard) GetLevel() int {
	return c.CardLevel
}

func (c *SummonerCard) GetDebuffs() map[Ability]int {
	return c.DebuffMap
}

func (c *SummonerCard) GetBuffs() map[Ability]int {
	return c.BuffMap
}

func (c *SummonerCard) Clone() SummonerCard {
	clonedCard := SummonerCard{
		cardDetail: c.cardDetail,
		GameCard: GameCard{
			CardLevel:      c.CardLevel,
			Team:           c.Team,
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
		},
	}
	clonedCard.SetTeam(c.GetTeamNumber())
	return clonedCard
}

func (c *SummonerCard) AddAbilitiesWithArray(abilities []Ability) {
	for _, a := range abilities {
		c.Abilities = append(c.Abilities, a)
	}
}

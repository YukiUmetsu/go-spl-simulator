package game_models

import (
	"fmt"

	utils "github.com/YukiUmetsu/go-spl-simulator/game_utils"
)

type SummonerCard struct {
	GameCard
	cardDetail CardDetail
}

func (c *SummonerCard) Setup(cardDetail CardDetail, cardLevel int) {
	c.cardDetail = cardDetail
	c.CardLevel = cardLevel - 1
	var summonerStats FlatCardStats

	var abilities []Ability = make([]Ability, 0)
	for _, a := range cardDetail.Stats.Abilities {
		if len(a.(string)) < 1 {
			continue
		}
		abilities = append(abilities, Ability(a.(string)))
	}
	summonerStats.Abilities = abilities
	summonerStats.Mana = int(cardDetail.Stats.Mana.(float64))
	summonerStats.Attack = int(cardDetail.Stats.Attack.(float64))
	summonerStats.Ranged = int(cardDetail.Stats.Ranged.(float64))
	summonerStats.Magic = int(cardDetail.Stats.Magic.(float64))
	summonerStats.Armor = int(cardDetail.Stats.Armor.(float64))
	summonerStats.Speed = int(cardDetail.Stats.Speed.(float64))
	summonerStats.Health = int(cardDetail.Stats.Health.(float64))
	c.SetStats(summonerStats)
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

func (c *SummonerCard) AddAbilities(abilities []Ability) {
	for _, ability := range abilities {
		c.Abilities = append(c.Abilities, ability)
	}
}

func (c *SummonerCard) AddAbility(ability Ability) {
	c.Abilities = append(c.Abilities, ability)
}

func (c *SummonerCard) GetCardDetail() CardDetail {
	return c.cardDetail
}

func (c *SummonerCard) SetCardDetail(cardDetail CardDetail) {
	c.cardDetail = cardDetail
}

func (c *SummonerCard) HasAbility(ability Ability) bool {
	return utils.Contains(c.Abilities, ability)
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

func (c *SummonerCard) Clone() GameCardInterface {
	var clonedCard GameCardInterface = &SummonerCard{
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

func (c *SummonerCard) String() string {
	return fmt.Sprintf(" S[ Name:%s(%v), Lvl: %v, Team: %v, Abilities: %v ]", c.cardDetail.Name, c.cardDetail.ID, c.CardLevel, c.GetTeamNumber(), c.Abilities)
}

/* Returns the card level (0 indexed) */
func (c *SummonerCard) GetCardLevel() int {
	return c.CardLevel
}

func (c *SummonerCard) GetCleanCard() *SummonerCard {
	var summoner *SummonerCard = &SummonerCard{}
	summoner.Setup(c.cardDetail, c.GetCardLevel()+1)
	return summoner
}

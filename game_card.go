package main

type GameCard interface {
	SetTeam(TeamNumber)
	HasAbility(ability Ability)
	RemoveAbility(ability Ability)
	GetTeamNumber() TeamNumber
	GetRarity() int
	GetName() string
	GetLevel() int
	GetDebuffs() map[Ability]int
	GetBuffs() map[Ability]int
	Clone() GameCard
}

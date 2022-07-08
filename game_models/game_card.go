package game_models

type GameCardInterface interface {
	SetTeam(TeamNumber)
	HasAbility(ability Ability) bool
	RemoveAbility(ability Ability)
	GetTeamNumber() TeamNumber
	GetRarity() int
	GetName() string
	GetLevel() int
	GetDebuffs() map[Ability]int
	GetBuffs() map[Ability]int
	Clone() GameCardInterface
	String() string
}

type GameCard struct {
	CardLevel      int
	Team           TeamNumber
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

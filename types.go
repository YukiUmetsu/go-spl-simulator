package main

type FlatCardStats struct {
	Abilities []Ability
	Mana      int
	Attack    int
	Ranged    int
	Magic     int
	Armor     int
	Health    int
	Speed     int
}

type CardStatsByLevel struct {
	Abilities [][]Ability
	Mana      []int
	Attack    []int
	Ranged    []int
	Magic     []int
	Armor     []int
	Health    []int
	Speed     []int
}

type CardDetailDistribution struct {
	CardDetailID int
	Gold         bool
	Edition      string
}

type CardDetail struct {
	ID              int
	Name            string
	Color           CardColor
	Card            CardType
	Rarity          int
	IsStarter       bool
	Editions        string
	DropRate        int
	SubType         string
	CreatedBlockNum int
	LastUpdateTx    string
	TotalPrinted    int
	IsPromo         bool
	Tier            string
	Distribution    []CardDetailDistribution
}

type SummonerCardDetail struct {
	CardDetail
	Stats FlatCardStats
}

type MonsterCardDetail struct {
	CardDetail
	Stats CardStatsByLevel
}

type Team struct {
	Summoner SummonerCardDetail
	Monsters []MonsterCardDetail
}

type BattleDamage struct {
	Attack           int // Some things care about if actually hit or not. Don't use this to check damage
	DamageDone       int // Actual damage done after modifiers, overkills. So 10 dmg to a 1 health is still 10 dmg.
	Remainder        int // Remainder damage after modifiers
	ActualDamageDone int // Actual damage done after modifiers, but does not overkill. 10 dmg to a 1 health is 1 dmg.
}

type BattleHistory struct {
	battle_queue_id_1       string
	battle_queue_id_2       string
	player_1_rating_initial int
	player_2_rating_initial int
	winner                  string
	player_1_rating_final   int
	player_2_rating_final   int
	player_1                string
	player_2                string
	created_date            string
	mana_cap                int
	ruleset                 string
	inactive                string
	settings                string
	details                 string
}

type CollectionCard struct {
	Player       string
	UUID         string
	CardDetailID int
	IsGoldFoil   bool
	Edition      int
	Level        int
}

type BattleTeam struct {
	Color    string
	Monsters []CollectionCard
	Summoner []CollectionCard
	Player   string
	Rating   int
}

type CustomCardStats struct {
	CardDetail
	// Custom
	IsSelected    bool
	Mana          int
	Edition       CardEdition
	AdjustedLevel int
}

type BattleAllowedCards struct {
	// Brawl and tournament?
	Foil string
	// Tournament only?
	Type     string
	Editions []int
}

type BattleLogAction interface{}

type BattleLog struct {
	/** The summoner or monster performing the action. This is a snapshot of the actor AFTER the action has been performed. */
	Actor GameCard
	/** The target of the action. This is a snapshot of the target AFTER the action has been performed. */
	Target GameCard
	/** The action */
	Action BattleLogAction
	/** The value, can be the amount of damage, or heal, etc. Based on the action */
	Value int
}

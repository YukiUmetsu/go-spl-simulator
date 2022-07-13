package game_models

type FlatCardStats struct {
	Mana      int       `json:"mana"`
	Attack    int       `json:"attack"`
	Ranged    int       `json:"ranged"`
	Magic     int       `json:"magic"`
	Armor     int       `json:"armor"`
	Health    int       `json:"health"`
	Speed     int       `json:"speed"`
	Abilities []Ability `json:"abilities,omitempty"`
}

type CardStatsByLevel struct {
	Mana      []int       `json:"mana"`
	Attack    []int       `json:"attack"`
	Ranged    []int       `json:"ranged"`
	Magic     []int       `json:"magic"`
	Armor     []int       `json:"armor"`
	Health    []int       `json:"health"`
	Speed     []int       `json:"speed"`
	Abilities [][]Ability `json:"abilities,omitempty"`
}

type CardDetailDistribution struct {
	CardDetailID int
	Gold         bool
	Edition      int
}

type CardDetail struct {
	ID              int                      `json:"id"`
	Name            string                   `json:"name"`
	Color           CardColor                `json:"color"`
	Type            CardType                 `json:"type"`
	Rarity          int                      `json:"rarity"`
	IsStarter       bool                     `json:"is_starter"`
	Editions        string                   `json:"editions"`
	DropRate        int                      `json:"drop_rate"`
	SubType         string                   `json:"sub_type"`
	CreatedBlockNum int                      `json:"created_block_num"`
	LastUpdateTx    string                   `json:"last_update_tx"`
	TotalPrinted    int                      `json:"total_printed"`
	IsPromo         bool                     `json:"is_promo"`
	Tier            int                      `json:"tier"`
	Distribution    []CardDetailDistribution `json:"distribution"`
	Stats           CardRawStats             `json:"stats"`
}

type CardRawStats struct {
	Mana      any   `json:"mana"`
	Attack    any   `json:"attack"`
	Ranged    any   `json:"ranged"`
	Magic     any   `json:"magic"`
	Armor     any   `json:"armor"`
	Health    any   `json:"health"`
	Speed     any   `json:"speed"`
	Abilities []any `json:"abilities,omitempty"`
}

type CardDetailMap map[int]CardDetail
type CardDetailMapPerName map[string]CardDetail

type Team struct {
	Summoner CardDetail
	Monsters []CardDetail
}

type BattleDamage struct {
	Attack           int // Some things care about if actually hit or not. Don't use this to check damage
	DamageDone       int // Actual damage done after modifiers, overkills. So 10 dmg to a 1 health is still 10 dmg.
	Remainder        int // Remainder damage after modifiers
	ActualDamageDone int // Actual damage done after modifiers, but does not overkill. 10 dmg to a 1 health is 1 dmg.
}

type BattleHistory struct {
	BattleQueueId1       string `json:"battle_queue_id_1"`
	BattleQueueId2       string `json:"battle_queue_id_2"`
	Player1RatingInitial int    `json:"player_1_rating_initial"`
	Player2RatingInitial int    `json:"player_2_rating_initial"`
	Winner               string `json:"winner"`
	Player1RatingFinal   int    `json:"player_1_rating_final"`
	Player2RatingFinal   int    `json:"player_2_rating_final"`
	Player1              string `json:"player_1"`
	Player2              string `json:"player_2"`
	CreatedDate          string `json:"created_date"`
	ManaCap              int    `json:"created_block_num"`
	Ruleset              string `json:"ruleset"`
	Inactive             string `json:"inactive"`
	Settings             string `json:"settings"`
	Details              string `json:"details"`
}

type BattleDetails struct {
	Loser  string     `json:"loser"`
	Winner string     `json:"winner"`
	Type   string     `json:"type"`
	Team1  BattleTeam `json:"team1"`
	Team2  BattleTeam `json:"team2"`
}

type CollectionCard struct {
	UID          string      `json:"uid"`
	XP           int         `json:"xp"`
	CardDetailID int         `json:"card_detail_id"`
	Gold         bool        `json:"gold"`
	Edition      int         `json:"edition"`
	Level        int         `json:"level"`
	State        BattleState `json:"state"`
}

type BattleState struct {
	Alive      bool  `json:"alive"`
	Stats      []int `json:"stats"`
	BaseHealth int   `json:"base_health"`
}

type BattleTeam struct {
	Color    string
	Monsters []CollectionCard
	Summoner CollectionCard
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

type Stringer interface {
	String() string
}

type BattleLog struct {
	/** The summoner or monster performing the action. This is a snapshot of the actor AFTER the action has been performed. */
	Actor GameCardInterface
	/** The target of the action. This is a snapshot of the target AFTER the action has been performed. */
	Target GameCardInterface
	/** The action */
	Action Stringer
	/** The value, can be the amount of damage, or heal, etc. Based on the action */
	Value int
}

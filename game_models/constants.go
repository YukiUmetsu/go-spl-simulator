package game_models

type CardType string

const (
	SUMMONER CardType = "Monster"
	MONSTER           = "Summoner"
)

type CardColor string

const (
	COLOR_BLACK CardColor = "Black"
	COLOR_BLUE            = "Blue"
	COLOR_GOLD            = "Gold"
	COLOR_GRAY            = "Gray"
	COLOR_GREEN           = "Green"
	COLOR_RED             = "Red"
	COLOR_WHITE           = "White"
)

type Ability string

const (
	ABILITY_AFFLICTION        Ability = "Affliction"
	ABILITY_AMPLIFY                   = "Amplify"
	ABILITY_BACKFIRE                  = "Backfire"
	ABILITY_BLAST                     = "Blast"
	ABILITY_BLIND                     = "Blind"
	ABILITY_BLOODLUST                 = "Bloodlust"
	ABILITY_CAMOUFLAGE                = "camouflage"
	ABILITY_CLEANSE                   = "Cleanse"
	ABILITY_CLOSE_RANGE               = "Close Range"
	ABILITY_CRIPPLE                   = "Cripple"
	ABILITY_DEATHBLOW                 = "Deathblow"
	ABILITY_DEMORALIZE                = "Demoralize"
	ABILITY_DISPEL                    = "Dispel"
	ABILITY_DIVINE_SHIELD             = "Divine Shield"
	ABILITY_DODGE                     = "Dodge"
	ABILITY_DOUBLE_STRIKE             = "Double Strike"
	ABILITY_ENRAGE                    = "Enrage"
	ABILITY_FLYING                    = "Flying"
	ABILITY_FORCEFIELD                = "Forcefield"
	ABILITY_GIANT_KILLER              = "Giant Killer"
	ABILITY_HALVING                   = "Halving"
	ABILITY_HEADWINDS                 = "Headwinds"
	ABILITY_HEAL                      = "Heal"
	ABILITY_IMMUNITY                  = "Immunity"
	ABILITY_INSPIRE                   = "Inspire"
	ABILITY_KNOCK_OUT                 = "Knock Out"
	ABILITY_LAST_STAND                = "Last Stand"
	ABILITY_LIFE_LEECH                = "Life Leech"
	ABILITY_MAGIC_REFLECT             = "Magic Reflect"
	ABILITY_OPPORTUNITY               = "Opportunity"
	ABILITY_OPPRESS                   = "Oppress"
	ABILITY_PHASE                     = "Phase"
	ABILITY_PIERCING                  = "Piercing"
	ABILITY_POISON                    = "Poison"
	ABILITY_PROTECT                   = "Protect"
	ABILITY_REACH                     = "Reach"
	ABILITY_RECHARGE                  = "Recharge"
	ABILITY_REDEMPTION                = "Redemption"
	ABILITY_REFLECTION_SHIELD         = "Reflection Shield"
	ABILITY_REPAIR                    = "Repair"
	ABILITY_RESURRECT                 = "Resurrect"
	ABILITY_RETALIATE                 = "Retaliate"
	ABILITY_RETURN_FIRE               = "Return Fire"
	ABILITY_RUST                      = "Rust"
	ABILITY_SCATTERSHOT               = "Scattershot"
	ABILITY_SCAVENGER                 = "Scavenger"
	ABILITY_SHATTER                   = "Shatter"
	ABILITY_SHIELD                    = "Shield"
	ABILITY_SILENCE                   = "Silence"
	ABILITY_SLOW                      = "Slow"
	ABILITY_SNARE                     = "Snare"
	ABILITY_SNEAK                     = "Sneak"
	ABILITY_SNIPE                     = "Snipe"
	ABILITY_STRENGTHEN                = "Strengthen"
	ABILITY_STUN                      = "Stun"
	ABILITY_SWIFTNESS                 = "Swiftness"
	ABILITY_TANK_HEAL                 = "Tank Heal"
	ABILITY_TAUNT                     = "Taunt"
	ABILITY_THORNS                    = "Thorns"
	ABILITY_TRAMPLE                   = "Trample"
	ABILITY_TRIAGE                    = "Triage"
	ABILITY_TRUE_STRIKE               = "True Strike"
	ABILITY_VOID                      = "Void"
	ABILITY_VOID_ARMOR                = "Void Armor"
	ABILITY_WEAKEN                    = "Weaken"

	/* RULESET ABILITY */
	ABILITY_MELEE_MAYHEM = "Melee Mayhem"
)

type CardEdition int

const (
	ALPHA CardEdition = iota
	BETA
	PROMO
	REWARD
	UNTAMED
	DICE
	GLADIUS
	CHAOS
)

type CardAttackType int

const (
	ATTACK_TYPE_MELEE CardAttackType = iota
	ATTACK_TYPE_RANGED
	ATTACK_TYPE_MAGIC
)

type FoilType int

const (
	STANDARD_FOIL FoilType = iota
	GOLD_FOIL
)

type TeamNumber int

const (
	TEAM_NUM_UNKNOWN TeamNumber = iota
	TEAM_NUM_ONE
	TEAM_NUM_TWO
	TEAM_NUM_TIE
)

type AdditionalBattleAction string

const (
	BATTLE_ACTION_DEATH         AdditionalBattleAction = "Death"
	BATTLE_ACTION_EARTHQUAKE    AdditionalBattleAction = "Earthquake"
	BATTLE_ACTION_FATIGUE       AdditionalBattleAction = "Fatigue"
	BATTLE_ACTION_MELEE         AdditionalBattleAction = "Melee attack"
	BATTLE_ACTION_RANGED        AdditionalBattleAction = "Ranged attack"
	BATTLE_ACTION_MAGIC         AdditionalBattleAction = "Magic attack"
	BATTLE_ACTION_ROUND_START   AdditionalBattleAction = "Round start"
	BATTLE_ACTION_ATTACK        AdditionalBattleAction = "Attack"
	BATTLE_ACTION_ATTACK_DODGED AdditionalBattleAction = "Dodged"

	// abilities
	BATTLE_ACTION_AFFLICTION           AdditionalBattleAction = "Affliction"
	BATTLE_ACTION_BACKFIRE             AdditionalBattleAction = "Backfire"
	BATTLE_ACTION_BLAST                AdditionalBattleAction = "Blast"
	BATTLE_ACTION_BLOODLUST            AdditionalBattleAction = "Bloodlust"
	BATTLE_ACTION_CLEANSE              AdditionalBattleAction = "Cleanse"
	BATTLE_ACTION_CRIPPLE              AdditionalBattleAction = "Cripple"
	BATTLE_ACTION_DIVINE_SHIELD        AdditionalBattleAction = "Divine shield"
	BATTLE_ACTION_HALVING              AdditionalBattleAction = "Halving"
	BATTLE_ACTION_HEAL                 AdditionalBattleAction = "Heal"
	BATTLE_ACTION_LIFE_LEECH           AdditionalBattleAction = "Life leech"
	BATTLE_ACTION_MAGIC_REFLECT        AdditionalBattleAction = "Magic reflect"
	BATTLE_ACTION_POISON               AdditionalBattleAction = "Poison"
	BATTLE_ACTION_REMOVE_DIVINE_SHIELD AdditionalBattleAction = "Remove divine shield"
	BATTLE_ACTION_REPAIR               AdditionalBattleAction = "Repair"
	BATTLE_ACTION_RESURRECT            AdditionalBattleAction = "Resurrect"
	BATTLE_ACTION_RETALIATE            AdditionalBattleAction = "Retaliate"
	BATTLE_ACTION_RETURN_FIRE          AdditionalBattleAction = "Return fire"
	BATTLE_ACTION_STUN                 AdditionalBattleAction = "Stun"
	BATTLE_ACTION_SCAVENGER            AdditionalBattleAction = "Scavenger"
	BATTLE_ACTION_SNARE                AdditionalBattleAction = "Snare"
	BATTLE_ACTION_TANK_HEAL            AdditionalBattleAction = "Tank heal"
	BATTLE_ACTION_THORNS               AdditionalBattleAction = "Thorns"
	BATTLE_ACTION_TRIAGE               AdditionalBattleAction = "Triage"
)

func (a AdditionalBattleAction) String() string {
	return string(a)
}

type Ruleset string

const (
	RULESET_AIM_TRUE              Ruleset = "Aim True"
	RULESET_ARMORED_UP                    = "Armored Up"
	RULESET_BACK_TO_BASICS                = "Back to Basics"
	RULESET_BROKEN_ARROWS                 = "Broken Arrows"
	RULESET_CLOSE_RANGE                   = "Close Range"
	RULESET_EARTHQUAKE                    = "Earthquake"
	RULESET_EQUAL_OPPORTUNITY             = "Equal Opportunity"
	RULESET_EQUALIZER                     = "Equalizer"
	RULESET_EVEN_STEVENS                  = "Even Stevens"
	RULESET_EXPLOSIVE_WEAPONRY            = "Explosive Weaponry"
	RULESET_FOG_OF_WAR                    = "Fog of War"
	RULESET_HEALED_OUT                    = "Healed Out"
	RULESET_HEAVY_HITTERS                 = "Heavy Hitters"
	RULESET_HOLY_PROTECTION               = "Holy Protection"
	RULESET_KEEP_YOUR_DISTANCE            = "Keep Your Distance"
	RULESET_LITTLE_LEAGUE                 = "Little League"
	RULESET_LOST_LEGENDARIES              = "Lost Legendaries"
	RULESET_LOST_MAGIC                    = "Lost Magic"
	RULESET_MELEE_MAYHEM                  = "Melee Mayhem"
	RULESET_NOXIOUS_FUMES                 = "Noxious Fumes"
	RULESET_ODD_ONES_OUT                  = "Odd Ones Out"
	RULESET_REVERSE_SPEED                 = "Reverse Speed"
	RULESET_RISE_OF_THE_COMMONS           = "Rise of the Commons"
	RULESET_SILENCED_SUMMONERS            = "Silenced Summoners"
	RULESET_SPREADING_FURY                = "Spreading Fury"
	RULESET_STAMPEDE                      = "Stampede"
	RULESET_STANDARD                      = "Standard"
	RULESET_SUPER_SNEAK                   = "Super Sneak"
	RULESET_TAKING_SIDES                  = "Taking Sides"
	RULESET_TARGET_PRACTICE               = "Target Practice"
	RULESET_UNPROTECTED                   = "Unprotected"
	RULESET_UP_CLOSE_AND_PERSONAL         = "Up Close & Personal"
	RULESET_WEAK_MAGIC                    = "Weak Magic"
)

type Stat int

const (
	STAT_MANA Stat = iota
	STAT_ATTACK
	STAT_MAGIC
	STAT_RANGED
	STAT_ARMOR
	STAT_SPEED
	STAT_HEALTH
)

const (
	// Abilities multiplier
	EARTHQUAKE_DAMAGE      = 2
	RUST_AMOUNT            = 2
	REDEMPTION_DAMAGE      = 1
	BACKFIRE_DAMAGE        = 2
	THORNS_DAMAGE          = 2
	PROTECT_AMOUNT         = 2
	CRIPPLE_AMOUNT         = 1
	WEAKEN_AMOUNT          = 1
	STRENGTHEN_AMOUNT      = 1
	LIFE_LEECH_AMOUNT      = 1
	SCAVENGER_AMOUNT       = 1
	REPAIR_AMOUNT          = 2
	TANK_HEAL_MULTIPLIER   = 1 / 3
	TRIAGE_HEAL_MULTIPLIER = 1 / 3
	MINIMUM_TRIAGE_HEAL    = 2
	MINIMUM_SELF_HEAL      = 2
	BLAST_MULTIPLIER       = 1 / 2
	FORCEFIELD_MIN_DAMAGE  = 5
	POISON_DAMAGE          = 2

	// Ability hit chance
	AFFLICTION_CHANCE       = 1 / 2
	RETALIATE_CHANCE        = 1 / 2
	STUN_CHANCE             = 1 / 2
	POISON_CHANCE           = 1 / 2
	DODGE_CHANCE            = 0.25
	FLYING_DODGE_CHANCE     = 0.25
	BLIND_DODGE_CHANCE      = 0.15
	SPEED_DIFF_DODGE_CHANCE = 1 / 10

	// Multiplier
	LAST_STAND_MULTIPLIER = 1.5
	ENRAGE_MULTIPLIER     = 1.5
)

const (
	FATIGUE_ROUND_NUMBER = 20
)

package simulator

type CardType int

const (
	SUMMONER CardType = iota
	MONSTER
)

type CardColor int

const (
	COLOR_BLACK CardColor = iota
	COLOR_BLUE
	COLOR_GOLD
	COLOR_GRAY
	COLOR_GREEN
	COLOR_RED
	COLOR_WHITE
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

type AdditionalBattleAction int

const (
	BATTLE_ACTION_DEATH AdditionalBattleAction = iota
	BATTLE_ACTION_EARTHQUAKE
	BATTLE_ACTION_FATIGUE
	BATTLE_ACTION_MELEE
	BATTLE_ACTION_RANGED
	BATTLE_ACTION_MAGIC
	BATTLE_ACTION_ROUND_START
	BATTLE_ACTION_SCAVENGER
)

type Ruleset int

const (
	RULESET_AIM_TRUE Ruleset = iota
	RULESET_ARMORED_UP
	RULESET_BACK_TO_BASICS
	RULESET_BROKEN_ARROWS
	RULESET_CLOSE_RANGE
	RULESET_EARTHQUAKE
	RULESET_EQUAL_OPPORTUNITY
	RULESET_EQUALIZER
	RULESET_EVEN_STEVENS
	RULESET_EXPLOSIVE_WEAPONRY
	RULESET_FOG_OF_WAR
	RULESET_HEALED_OUT
	RULESET_HEAVY_HITTERS
	RULESET_HOLY_PROTECTION
	RULESET_KEEP_YOUR_DISTANCE
	RULESET_LITTLE_LEAGUE
	RULESET_LOST_LEGENDARIES
	RULESET_LOST_MAGIC
	RULESET_MELEE_MAYHEM
	RULESET_NOXIOUS_FUMES
	RULESET_ODD_ONES_OUT
	RULESET_REVERSE_SPEED
	RULESET_RISE_OF_THE_COMMONS
	RULESET_SILENCED_SUMMONERS
	RULESET_SPREADING_FURY
	RULESET_STAMPEDE
	RULESET_STANDARD
	RULESET_SUPER_SNEAK
	RULESET_TAKING_SIDES
	RULESET_TARGET_PRACTICE
	RULESET_UNPROTECTED
	RULESET_UP_CLOSE_AND_PERSONAL
	RULESET_WEAK_MAGIC
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
	AFFLICTION_CHANCE   = 1 / 2
	RETALIATE_CHANCE    = 1 / 2
	STUN_CHANCE         = 1 / 2
	POISON_CHANCE       = 1 / 2
	DODGE_CHANCE        = 0.25
	FLYING_DODGE_CHANCE = 0.25
	BLIND_DODGE_CHANCE  = 0.15

	// Multiplier
	LAST_STAND_MULTIPLIER = 1.5
	ENRAGE_MULTIPLIER     = 1.5
)

const (
	FATIGUE_ROUND_NUMBER = 20
)

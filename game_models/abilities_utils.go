package game_models

import (
	"math"

	utils "github.com/YukiUmetsu/go-spl-simulator/game_utils"
)

/** Abilities that summoner applies to friendly team at the start of the game */
func GetSummonerAbilityAbilities() []Ability {
	return []Ability{
		ABILITY_BLAST,
		ABILITY_DIVINE_SHIELD,
		ABILITY_FLYING,
		ABILITY_LAST_STAND,
		ABILITY_MAGIC_REFLECT,
		ABILITY_PIERCING,
		ABILITY_RETURN_FIRE,
		ABILITY_SNARE,
		ABILITY_THORNS,
		ABILITY_TRUE_STRIKE,
		ABILITY_VOID,
		ABILITY_VOID_ARMOR,
		ABILITY_POISON,
	}
}

/** Buffs that summoner applies to friendly team at the start of the game */
func GetSummonerPreGameBuffAbilities() []Ability {
	return []Ability{
		ABILITY_STRENGTHEN,
	}
}

/** Abilities that summoner applies to enemy team at the start of the game */
func GetSummonerPreGameDebuffAbilities() []Ability {
	return []Ability{
		ABILITY_AFFLICTION,
		ABILITY_BLIND,
	}
}

/** Abilities that monsters apply to friendly team at the start of the game */
func GetMonsterPreGameBuffAbilities() []Ability {
	return []Ability{
		ABILITY_PROTECT,
		ABILITY_STRENGTHEN,
		ABILITY_SWIFTNESS,
		ABILITY_INSPIRE,
	}
}

/** Abilities that monsters apply to enemy team at the start of the game */
func GetMonsterPreGameDebuffAbilities() []Ability {
	return []Ability{
		ABILITY_AMPLIFY,
		ABILITY_BLIND,
		ABILITY_DEMORALIZE,
		ABILITY_HEADWINDS,
		ABILITY_RUST,
		ABILITY_SLOW,
		ABILITY_SNARE,
		ABILITY_SILENCE,
		ABILITY_WEAKEN,
	}
}

/**
 * Abilities that can't be cleansed. (These aren't actually debuffs but this app codes them as a debuff)
 */
func GetUncleansableDebuffs() []Ability {
	return []Ability{
		ABILITY_AMPLIFY,
	}
}

/** Abilities that require a turn to do something */
func GetActionAbilities() []Ability {
	return []Ability{
		ABILITY_REPAIR,
		ABILITY_TANK_HEAL,
	}
}

/* Check if the monster has pre-game debuff abilities. If so, return those debuffs, otherwise empty array */
func MonsterHasDebuffAbilities(m *MonsterCard) []Ability {
	if m == nil {
		return []Ability{}
	}
	monsterDebuffs := []Ability{}
	debuffs := GetMonsterPreGameDebuffAbilities()
	for _, ability := range m.Abilities {
		if utils.Contains(debuffs, ability) {
			monsterDebuffs = append(monsterDebuffs, ability)
		}
	}
	return monsterDebuffs
}

func MonsterHasBuffsAbilities(m *MonsterCard) []Ability {
	monsterBuffs := []Ability{}
	if m == nil {
		return monsterBuffs
	}
	buffs := GetMonsterPreGameBuffAbilities()
	for _, ability := range m.Abilities {
		if utils.Contains(buffs, ability) {
			monsterBuffs = append(monsterBuffs, ability)
		}
	}
	return monsterBuffs
}

func RepairMonsterArmor(m *MonsterCard) int {
	if m == nil {
		return 0
	}
	previousArmor := m.Armor
	maxArmor := m.GetPostAbilityMaxArmor()
	newArmorAmount := utils.GetSmaller(maxArmor, (m.Armor + REPAIR_AMOUNT))
	m.Armor = newArmorAmount
	return newArmorAmount - previousArmor
}

func TankHealMonster(m *MonsterCard) int {
	if m == nil {
		return 0
	}
	if m.HasDebuff(ABILITY_AFFLICTION) {
		return 0
	}
	previousHealth := m.Health
	maxHealth := m.GetPostAbilityMaxHealth()
	healAmount := int(math.Floor(float64(maxHealth) * TANK_HEAL_MULTIPLIER))
	healAmount = utils.GetBigger(healAmount, 2)
	m.AddHealth(healAmount)
	return m.Health - previousHealth
}

func TriageHealMonster(m *MonsterCard) int {
	if m == nil || m.HasDebuff(ABILITY_AFFLICTION) {
		return 0
	}

	previousHealth := m.Health
	maxHealth := m.GetPostAbilityMaxHealth()
	healAmount := int(math.Floor(float64(maxHealth) * TRIAGE_HEAL_MULTIPLIER))
	healAmount = utils.GetBigger(healAmount, MINIMUM_TRIAGE_HEAL)
	m.AddHealth(healAmount)
	return m.Health - previousHealth
}

func RustMonster(m *MonsterCard) {
	if m == nil {
		return
	}
	m.Armor = utils.GetBigger(0, m.Armor-RUST_AMOUNT)
}

func ScavengerMonster(m *MonsterCard) {
	if m == nil {
		return
	}
	m.AddHealth(1)
}

func LifeLeechMonster(m *MonsterCard) {
	if m == nil {
		return
	}
	m.AddHealth(1)
}

func StrengthenMonster(m *MonsterCard) {
	if m == nil {
		return
	}
	m.AddHealth(1)
}

func ProtectMonster(m *MonsterCard) {
	if m == nil {
		return
	}
	m.Armor = m.Armor + PROTECT_AMOUNT
}

func WeakenMonster(m *MonsterCard) {
	m.Health = utils.GetBigger(m.Health-1, 1)
}

func SelfHealMonster(m *MonsterCard) int {
	if m.HasDebuff(ABILITY_AFFLICTION) {
		return 0
	}

	previousHealth := m.Health
	maxHealth := float64(m.GetPostAbilityMaxHealth())
	healAmount := utils.GetBigger(int(math.Floor(maxHealth/3)), MINIMUM_SELF_HEAL)
	m.AddHealth(healAmount)
	return m.Health - previousHealth
}

func DispelBuffs(m *MonsterCard) {
	buffMap := m.GetAllBuffs()
	for buff, _ := range buffMap {
		if buff == ABILITY_SCAVENGER || buff == ABILITY_LIFE_LEECH {
			m.RemoveBuff(buff)
		}
		m.RemoveAllBuff(buff)
	}
}

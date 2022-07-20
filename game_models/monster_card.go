package game_models

import (
	"fmt"
	"math"

	utils "github.com/YukiUmetsu/go-spl-simulator/game_utils"
)

type MonsterCard struct {
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

	cardDetail CardDetail

	// only monsters
	cardPosition    int
	isOnlyMonster   bool
	hasTurnPassed   bool
	summonerSpeed   int
	summonerArmor   int
	summonerMelee   int
	summonerRanged  int
	summonerMagic   int
	hadDivineShield bool
}

func (c *MonsterCard) Setup(cardDetail CardDetail, cardLevel int) {
	c.cardDetail = cardDetail
	c.CardLevel = cardLevel
	var cardStatsByLevel CardStatsByLevel

	// convert interface to ability
	abilityByLevel := make([][]Ability, 0)
	for _, abilityArr := range cardDetail.Stats.Abilities {
		abilitiesInLevel := []Ability{}
		for _, ability := range abilityArr.([]any) {
			abilitiesInLevel = append(abilitiesInLevel, Ability(ability.(string)))
		}
		abilityByLevel = append(abilityByLevel, abilitiesInLevel)
	}
	cardStatsByLevel.Abilities = abilityByLevel

	// convert interface to []int
	manaByLevel := make([]int, 0)
	for _, m := range cardDetail.Stats.Mana.([]any) {
		switch m.(type) {
		case float64:
			manaByLevel = append(manaByLevel, int(m.(float64)))
		default:
			manaByLevel = append(manaByLevel, m.(int))
		}
	}
	cardStatsByLevel.Mana = manaByLevel

	atkByLevel := make([]int, 0)
	for _, m := range cardDetail.Stats.Attack.([]any) {
		switch m.(type) {
		case float64:
			atkByLevel = append(atkByLevel, int(m.(float64)))
		default:
			atkByLevel = append(atkByLevel, m.(int))
		}
	}
	cardStatsByLevel.Attack = atkByLevel

	rangeByLevel := make([]int, 0)
	for _, m := range cardDetail.Stats.Ranged.([]any) {
		switch m.(type) {
		case float64:
			rangeByLevel = append(rangeByLevel, int(m.(float64)))
		default:
			rangeByLevel = append(rangeByLevel, m.(int))
		}
	}
	cardStatsByLevel.Ranged = rangeByLevel

	magicByLevel := make([]int, 0)
	for _, m := range cardDetail.Stats.Magic.([]any) {
		switch m.(type) {
		case float64:
			magicByLevel = append(magicByLevel, int(m.(float64)))
		default:
			magicByLevel = append(magicByLevel, m.(int))
		}
	}
	cardStatsByLevel.Magic = magicByLevel

	armorByLevel := make([]int, 0)
	for _, m := range cardDetail.Stats.Armor.([]any) {
		switch m.(type) {
		case float64:
			armorByLevel = append(armorByLevel, int(m.(float64)))
		default:
			armorByLevel = append(armorByLevel, m.(int))
		}
	}
	cardStatsByLevel.Armor = armorByLevel

	speedByLevel := make([]int, 0)
	for _, m := range cardDetail.Stats.Speed.([]any) {
		switch m.(type) {
		case float64:
			speedByLevel = append(speedByLevel, int(m.(float64)))
		default:
			speedByLevel = append(speedByLevel, m.(int))
		}
	}
	cardStatsByLevel.Speed = speedByLevel

	hpByLevel := make([]int, 0)
	for _, m := range cardDetail.Stats.Health.([]any) {
		switch m.(type) {
		case float64:
			hpByLevel = append(hpByLevel, int(m.(float64)))
		default:
			hpByLevel = append(hpByLevel, m.(int))
		}
	}
	cardStatsByLevel.Health = hpByLevel

	c.SetStats(cardStatsByLevel)
}

func (c *MonsterCard) SetTeam(teamNumber TeamNumber) {
	c.Team = teamNumber
}

func (c *MonsterCard) GetCardLevel() int {
	return c.CardLevel
}

func (c *MonsterCard) GetCleanCard() *MonsterCard {
	var monster *MonsterCard = &MonsterCard{}
	monster.Setup(c.cardDetail, c.GetCardLevel())
	return monster
}

func (c *MonsterCard) SetStats(stats CardStatsByLevel) {
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

func (c *MonsterCard) GetStat(stats []int) int {
	return stats[c.CardLevel-1]
}

func (c *MonsterCard) AddAbilities(abilitiesArray [][]Ability) {
	for i, abilities := range abilitiesArray {
		if i+1 <= c.CardLevel {
			for _, ability := range abilities {
				c.Abilities = append(c.Abilities, ability)
			}
		}
	}
}

func (c *MonsterCard) AddAbility(ability Ability) {
	c.Abilities = append(c.Abilities, ability)
}

func (c *MonsterCard) GetCardDetail() CardDetail {
	return c.cardDetail
}

func (c *MonsterCard) HasAbility(ability Ability) bool {
	return utils.Contains(c.Abilities, ability)
}

func (c *MonsterCard) RemoveAbility(ability Ability) {
	c.Abilities = utils.Remove(c.Abilities, ability)
}

func (c *MonsterCard) RemoveAllAbilities() {
	c.Abilities = []Ability{}
}

func (c *MonsterCard) GetTeamNumber() TeamNumber {
	return c.Team
}

func (c *MonsterCard) GetRarity() int {
	return c.cardDetail.Rarity
}

func (c *MonsterCard) GetName() string {
	return c.cardDetail.Name
}

func (c *MonsterCard) GetLevel() int {
	return c.CardLevel
}

func (c *MonsterCard) GetDebuffs() map[Ability]int {
	return c.DebuffMap
}

func (c *MonsterCard) GetBuffs() map[Ability]int {
	return c.BuffMap
}

func (c *MonsterCard) AddBuff(buff Ability) {
	if !c.IsAlive() {
		return
	}

	var buffsAmount int
	if value, ok := c.BuffMap[buff]; ok {
		buffsAmount = value + 1
	} else {
		buffsAmount = 1
	}
	if c.BuffMap == nil {
		c.BuffMap = make(map[Ability]int)
	}
	c.BuffMap[buff] = buffsAmount

	if buff == ABILITY_SCAVENGER {
		ScavengerMonster(c)
	} else if buff == ABILITY_LIFE_LEECH {
		LifeLeechMonster(c)
	} else if buff == ABILITY_STRENGTHEN {
		StrengthenMonster(c)
	} else if buff == ABILITY_PROTECT {
		ProtectMonster(c)
	}
}

func (c *MonsterCard) AddDebuff(debuff Ability) {
	if !c.IsAlive() {
		return
	}

	// the card has immunity and it's not an uncleansable debuff => ignore
	uncleansableBuffs := GetUncleansableDebuffs()
	if utils.Contains(c.Abilities, ABILITY_IMMUNITY) && !utils.Contains(uncleansableBuffs, debuff) {
		return
	}

	// debuff is snare and snare is already applied => ignore
	if debuff == ABILITY_SNARE && c.HasDebuff(ABILITY_SNARE) {
		return
	}

	if debuff == ABILITY_WEAKEN {
		WeakenMonster(c)
	} else if debuff == ABILITY_RUST {
		RustMonster(c)
	}

	var debuffAmount int
	if value, ok := c.DebuffMap[debuff]; ok {
		debuffAmount = value + 1
	} else {
		debuffAmount = 1
	}

	if c.DebuffMap == nil {
		c.DebuffMap = make(map[Ability]int)
	}
	c.DebuffMap[debuff] = debuffAmount
}

func (c *MonsterCard) GetHasTurnPassed() bool {
	return c.hasTurnPassed
}

func (c *MonsterCard) SetHasTurnPassed(hasTurnPassed bool) {
	c.hasTurnPassed = hasTurnPassed
}

func (c *MonsterCard) Clone() GameCardInterface {
	var clonedCard GameCardInterface = &MonsterCard{
		cardDetail:     c.cardDetail,
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
	}
	clonedCard.SetTeam(c.GetTeamNumber())
	return clonedCard
}

func (c *MonsterCard) AddAbilitiesWithArray(abilities []Ability) {
	for _, a := range abilities {
		c.Abilities = append(c.Abilities, a)
	}
}

/* Returns remaining damage after hitting health */
func (c *MonsterCard) HitHealth(damage int) int {
	preHitHealth := c.Health
	c.AddHealth(-1 * damage)
	if c.Health < 0 {
		return c.Health * -1
	}

	if c.Health == 0 {
		return damage - preHitHealth
	}

	return 0
}

// monster only
func (c *MonsterCard) SetCardPosition(position int) {
	c.cardPosition = position
}

func (c *MonsterCard) GetCardPosition() int {
	return c.cardPosition
}

func (c *MonsterCard) IsAlive() bool {
	return c.Health > 0
}

func (c *MonsterCard) SetHealth(health int) {
	c.Health = health
}

func (c *MonsterCard) GetHealth() int {
	return c.Health
}

func (c *MonsterCard) AddHealth(amount int) {
	if !c.IsAlive() || amount == 0 {
		return
	}

	finalHealth := utils.GetSmaller(c.Health+amount, c.GetPostAbilityMaxHealth())
	finalHealth = utils.GetBigger(finalHealth, 0)
	c.SetHealth(finalHealth)
}

func (c *MonsterCard) HasBuff(buff Ability) bool {
	_, ok := c.BuffMap[buff]
	return ok
}

func (c *MonsterCard) GetBuffCount(buff Ability) int {
	if val, ok := c.BuffMap[buff]; ok {
		return val
	} else {
		return 0
	}
}

func (c *MonsterCard) HasDebuff(debuff Ability) bool {
	_, ok := c.DebuffMap[debuff]
	return ok
}

func (c *MonsterCard) GetDebuffCount(debuff Ability) int {
	if val, ok := c.DebuffMap[debuff]; ok {
		return val
	} else {
		return 0
	}
}

func (c *MonsterCard) GetIsLastStand() bool {
	return c.isOnlyMonster && c.HasAbility(ABILITY_LAST_STAND)
}

func (c *MonsterCard) IsLastMonster() bool {
	return c.isOnlyMonster
}

func (c *MonsterCard) SetIsOnlyMonster() {
	if c.HasAbility(ABILITY_LAST_STAND) {
		prevMaxHealth := c.GetPostAbilityMaxHealth()
		dmgTaken := prevMaxHealth - c.Health
		c.isOnlyMonster = true
		c.SetHealth(c.GetPostAbilityMaxHealth() - dmgTaken)
	}
	c.isOnlyMonster = true
}

func (c *MonsterCard) HasAttack() bool {
	return c.Melee > 0 || c.Ranged > 0 || c.Magic > 0
}

func (c *MonsterCard) GetPostAbilityMaxHealth() int {
	maxHealth := 1
	if c.StartingHealth > maxHealth {
		maxHealth = c.StartingHealth
	}

	// Life leech and scavenger are affected by last stand multiplier
	if c.HasBuff(ABILITY_LIFE_LEECH) {
		lifeLeechAmount := c.GetBuffCount(ABILITY_LIFE_LEECH)
		maxHealth = maxHealth + LIFE_LEECH_AMOUNT*lifeLeechAmount
	}
	if c.HasBuff(ABILITY_SCAVENGER) {
		scavengerAmout := c.GetBuffCount(ABILITY_SCAVENGER)
		maxHealth = maxHealth + SCAVENGER_AMOUNT*scavengerAmout
	}

	if c.HasBuff(ABILITY_CRIPPLE) {
		maxHealth = maxHealth - CRIPPLE_AMOUNT*c.GetDebuffCount(ABILITY_CRIPPLE)
	}
	if c.GetIsLastStand() {
		maxHealth = int(math.Floor(float64(maxHealth) * LAST_STAND_MULTIPLIER))
	}
	if c.HasDebuff(ABILITY_WEAKEN) {
		maxHealth = maxHealth - WEAKEN_AMOUNT*c.GetDebuffCount(ABILITY_WEAKEN)
	}
	if c.HasBuff(ABILITY_STRENGTHEN) {
		maxHealth = maxHealth + STRENGTHEN_AMOUNT*c.GetBuffCount(ABILITY_STRENGTHEN)
	}

	// The summoner skill made this starting health 0 or negative
	if c.StartingHealth < 1 {
		maxHealth = maxHealth + c.StartingHealth - 1
	}

	return utils.GetBigger(maxHealth, 1)
}

func (c *MonsterCard) GetPostAbilityAttackOfType(attackType CardAttackType) int {
	if attackType == ATTACK_TYPE_MAGIC {
		return c.GetPostAbilityMagic()
	}
	if attackType == ATTACK_TYPE_RANGED {
		return c.GetPostAbilityRange()
	}
	if attackType == ATTACK_TYPE_MELEE {
		return c.GetPostAbilityMelee()
	}
	return 0
}

/**  How much magic damage this will do */
func (c *MonsterCard) GetPostAbilityMagic() int {
	if c.Magic == 0 {
		return 0
	}
	postMagic := c.Magic
	if c.HasDebuff(ABILITY_HALVING) {
		postMagic = int(math.Floor((float64(postMagic) + 1) / 2))
		postMagic = utils.GetBigger(postMagic, 1)
	}
	if c.GetIsLastStand() {
		postMagic = int(math.Floor((float64(postMagic) * LAST_STAND_MULTIPLIER)))
	}

	// calculate the magic modifier
	magicModifier := 0
	for i := 0; i < c.GetDebuffCount(ABILITY_SILENCE); i++ {
		magicModifier = magicModifier - 1
	}
	if c.HasDebuff(ABILITY_HALVING) {
		magicModifier = magicModifier + int(math.Floor((float64(c.summonerMagic) / 2)))
	} else {
		magicModifier = magicModifier + c.summonerMagic
	}

	return utils.GetBigger(postMagic+magicModifier, 1)
}

/**  How much range damage this will do */
func (c *MonsterCard) GetPostAbilityRange() int {
	if c.Ranged == 0 {
		return 0
	}
	postRange := c.Ranged
	if c.HasDebuff(ABILITY_HALVING) {
		postRange = int(math.Floor((float64(postRange) + 1) / 2))
		postRange = utils.GetBigger(postRange, 1)
	}
	if c.GetIsLastStand() {
		postRange = int(math.Floor((float64(postRange) * LAST_STAND_MULTIPLIER)))
	}

	// calculate the magic modifier
	rangeModifier := 0
	for i := 0; i < c.GetDebuffCount(ABILITY_HEADWINDS); i++ {
		rangeModifier = rangeModifier - 1
	}
	if c.HasDebuff(ABILITY_HALVING) {
		rangeModifier = rangeModifier + int(math.Floor((float64(c.summonerRanged) / 2)))
	} else {
		rangeModifier = rangeModifier + c.summonerRanged
	}

	return utils.GetBigger(postRange+rangeModifier, 1)
}

/**  How much melee damage this will do */
func (c *MonsterCard) GetPostAbilityMelee() int {
	if c.Melee == 0 {
		return 0
	}
	postMelee := c.Melee
	if c.HasDebuff(ABILITY_HALVING) {
		postMelee = int(math.Floor((float64(postMelee) + 1) / 2))
		postMelee = utils.GetBigger(postMelee, 1)
	}
	if c.GetIsLastStand() {
		postMelee = int(math.Floor((float64(postMelee) * LAST_STAND_MULTIPLIER)))
	}

	// calculate the magic modifier
	meleeModifier := 0
	for i := 0; i < c.GetDebuffCount(ABILITY_DEMORALIZE); i++ {
		meleeModifier = meleeModifier - 1
	}
	for i := 0; i < c.GetBuffCount(ABILITY_INSPIRE); i++ {
		meleeModifier = meleeModifier + 1
	}
	if c.HasDebuff(ABILITY_HALVING) {
		meleeModifier = meleeModifier + int(math.Floor((float64(c.summonerMelee) / 2)))
	} else {
		meleeModifier = meleeModifier + c.summonerMelee
	}

	currentMelee := utils.GetBigger(postMelee+meleeModifier, 1)
	if c.IsEnraged() {
		return int(math.Ceil(float64(currentMelee) * ENRAGE_MULTIPLIER))
	}

	return currentMelee
}

func (c *MonsterCard) GetPostAbilitySpeed() int {
	speedModifier := 0
	speed := c.Speed + c.summonerSpeed
	if c.GetIsLastStand() {
		speed = int(math.Floor(float64(speed) * LAST_STAND_MULTIPLIER))
	}
	if c.IsEnraged() {
		speed = int(math.Ceil(float64(speed) * ENRAGE_MULTIPLIER))
	}
	if c.HasBuff(ABILITY_SWIFTNESS) {
		speedModifier = speedModifier + c.GetBuffCount(ABILITY_SWIFTNESS)
	}
	if c.HasDebuff(ABILITY_SLOW) {
		speedModifier = speedModifier - c.GetDebuffCount(ABILITY_SLOW)
	}
	return utils.GetBigger(speed+speedModifier, 1)
}

func (c *MonsterCard) GetPostAbilityMaxArmor() int {
	postArmor := c.StartingArmor
	if c.GetIsLastStand() {
		postArmor = int(math.Floor(float64(postArmor) * LAST_STAND_MULTIPLIER))
	}

	armorModifier := 0
	for i := 0; i < c.GetBuffCount(ABILITY_PROTECT); i++ {
		armorModifier = armorModifier + PROTECT_AMOUNT
	}
	for i := 0; i < c.GetDebuffCount(ABILITY_RUST); i++ {
		armorModifier = armorModifier - RUST_AMOUNT
	}

	return utils.GetBigger(postArmor+armorModifier, 0)
}

func (c *MonsterCard) RemoveBuff(buff Ability) {
	if !c.HasBuff(buff) {
		return
	}

	newBufCount := c.GetBuffCount(buff) - 1
	if newBufCount == 0 {
		// remove buff key from the map
		delete(c.BuffMap, buff)
	} else {
		c.BuffMap[buff] = newBufCount
	}

	// TODO: comeback and check the validity of this code because of the rule
	if buff == ABILITY_SCAVENGER {
		c.RemoveBuffHealth(1)
	} else if buff == ABILITY_LIFE_LEECH {
		c.RemoveBuffHealth(1)
	} else if buff == ABILITY_STRENGTHEN {
		c.RemoveStrengthenHealth(1)
	} else if buff == ABILITY_PROTECT {
		c.Armor = utils.GetSmaller(c.GetPostAbilityMaxArmor(), c.Armor)
	}
}

// TODO: comeback and check the validity of this code because of the rule
func (c *MonsterCard) RemoveStrengthenHealth(healthAmount int) {
	if c.Health < c.GetPostAbilityMaxHealth() {
		c.RemoveBuffHealth(healthAmount)
	}
}

// TODO: comeback and check the validity of this code because of the rule
func (c *MonsterCard) RemoveBuffHealth(healthAmount int) {
	if c.Health < 1 {
		return
	}
	c.Health = utils.GetBigger(c.Health-healthAmount, 1)
}

func (c *MonsterCard) RemoveDebuff(debuff Ability) {
	debuffAmount := c.GetDebuffCount(debuff) - 1
	if debuffAmount < 1 {
		if _, ok := c.DebuffMap[debuff]; ok {
			delete(c.DebuffMap, debuff)
		}
	} else {
		c.DebuffMap[debuff] = debuffAmount
	}

	if debuff == ABILITY_WEAKEN {
		c.AddHealth(1)
	} else if debuff == ABILITY_CRIPPLE {
		c.AddHealth(1)
	} else if debuff == ABILITY_RUST {
		c.Armor = utils.GetSmaller(c.Armor+2, c.GetPostAbilityMaxArmor())
	}
}

func (c *MonsterCard) RemoveAllDebuff(debuff Ability) {
	debuffAmount := c.GetDebuffCount(debuff)
	for i := 0; i < debuffAmount; i++ {
		c.RemoveDebuff(debuff)
	}
}

func (c *MonsterCard) CleanseDebuffs() {
	// Special case, cleanse only removes 1 cripple
	crippleAmount := c.GetDebuffCount(ABILITY_CRIPPLE)
	for ability, _ := range c.DebuffMap {
		if !utils.Contains(GetUncleansableDebuffs(), ability) {
			c.RemoveAllDebuff(ability)
		}
	}

	if crippleAmount > 1 {
		c.DebuffMap[ABILITY_CRIPPLE] = crippleAmount - 1
	}
	if crippleAmount > 0 && c.Health > 0 {
		c.AddHealth(1)
	}
}

func (c *MonsterCard) CleanseDebuffsAfterResurrect() {
	for _, aDebuff := range GetActionDebuffs() {
		c.RemoveDebuff(aDebuff)
	}
}

func (c *MonsterCard) Resurrect() {
	if c.hadDivineShield {
		c.AddAbilitiesWithArray([]Ability{ABILITY_DIVINE_SHIELD})
	}
	c.Armor = c.GetPostAbilityMaxArmor()
	c.CleanseDebuffsAfterResurrect()
	c.Health = 1
}

func (c *MonsterCard) IsEnraged() bool {
	return c.HasAbility(ABILITY_ENRAGE) && c.Health < c.GetPostAbilityMaxHealth()
}

func (c *MonsterCard) IsPureMelee() bool {
	return c.Melee > 0 && c.Ranged == 0 && c.Magic == 0
}

func (c *MonsterCard) CanMeleeAttack() bool {
	if c.HasAbility(ABILITY_OPPORTUNITY) ||
		c.HasAbility(ABILITY_SNEAK) ||
		c.cardPosition == 0 ||
		c.HasAbility(ABILITY_MELEE_MAYHEM) ||
		(c.HasAbility(ABILITY_REACH) && c.cardPosition == 1) {
		return true
	}
	return false
}

/* Summoner Related stuff */
func (c *MonsterCard) GetSummonerArmor() int {
	return c.summonerArmor
}

func (c *MonsterCard) AddSummonerArmor(stat int) {
	c.summonerArmor = c.summonerArmor + stat
	c.Armor = utils.GetBigger(c.Armor+stat, 1)
}

func (c *MonsterCard) AddSummonerHealth(stat int) {
	c.StartingHealth = c.StartingHealth + stat
	c.Health = utils.GetBigger(c.Health+stat, 1)
}

func (c *MonsterCard) AddSummonerSpeed(stat int) {
	c.summonerSpeed = c.summonerSpeed + stat
}

func (c *MonsterCard) AddSummonerMelee(stat int) {
	c.summonerMelee = c.summonerMelee + stat
}

func (c *MonsterCard) AddSummonerRanged(stat int) {
	c.summonerRanged = c.summonerRanged + stat
}

func (c *MonsterCard) AddSummonerMagic(stat int) {
	c.summonerMagic = c.summonerMagic + stat
}

func (c *MonsterCard) RemoveDivineShield() {
	c.hadDivineShield = true
	c.RemoveAbility(ABILITY_DIVINE_SHIELD)
}

func (c *MonsterCard) SetHasTurnPasses(hasPassed bool) {
	c.hasTurnPassed = hasPassed
}

func (c *MonsterCard) GetAllBuffs() map[Ability]int {
	return c.BuffMap
}

func (c *MonsterCard) GetBuffAmount(buff Ability) int {
	if val, ok := c.BuffMap[buff]; ok {
		return val
	}
	return 0
}

func (c *MonsterCard) RemoveAllBuff(buff Ability) {
	buffAmount := c.GetBuffAmount(buff)
	for i := 0; i < buffAmount; i++ {
		c.RemoveBuff(buff)
	}
}

func (c *MonsterCard) String() string {
	return fmt.Sprintf(
		"M[ Name: %s(%v), Lvl:%v, Team:%v, HP:%v, Speed:%v, Armor:%v, buffs:%+v, debuffs:%+v, abilities:%+v, (ML:%v, RG:%v, MG:%v) ] ",
		c.cardDetail.Name,
		c.cardDetail.ID,
		c.CardLevel,
		c.GetTeamNumber(),
		c.Health,
		c.GetPostAbilitySpeed(),
		c.Armor,
		c.BuffMap,
		c.DebuffMap,
		c.Abilities,
		c.GetPostAbilityMelee(),
		c.GetPostAbilityRange(),
		c.GetPostAbilityMagic(),
	)
}

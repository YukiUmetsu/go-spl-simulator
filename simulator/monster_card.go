package simulator

import (
	"math"
	utils "utils"
)

type MonsterCard struct {
	GameCard
	cardDetail MonsterCardDetail

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

func (c MonsterCard) Setup(cardDetail MonsterCardDetail, cardLevel int) {
	c.cardDetail = cardDetail
	c.CardLevel = cardLevel - 1
	c.SetStats(c.cardDetail.Stats)
}

func (c MonsterCard) SetTeam(teamNumber TeamNumber) {
	c.Team = teamNumber
}

func (c MonsterCard) SetStats(stats CardStatsByLevel) {
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

func (c MonsterCard) GetStat(stats []int) int {
	return stats[c.CardLevel]
}

func (c MonsterCard) AddAbilities(abilitiesArray [][]Ability) {
	for i, abilities := range abilitiesArray {
		if i+1 <= c.CardLevel {
			for _, ability := range abilities {
				c.Abilities = append(c.Abilities, ability)
			}
		}
	}
}

func (c MonsterCard) GetCardDetail() MonsterCardDetail {
	return c.cardDetail
}

func (c MonsterCard) HasAbility(ability Ability) bool {
	return c.HasAbility(ability)
}

func (c MonsterCard) RemoveAbility(ability Ability) {
	c.Abilities = utils.Remove(c.Abilities, ability)
}

func (c MonsterCard) GetTeamNumber() TeamNumber {
	return c.Team
}

func (c MonsterCard) GetRarity() int {
	return c.cardDetail.Rarity
}

func (c MonsterCard) GetName() string {
	return c.cardDetail.Name
}

func (c MonsterCard) GetLevel() int {
	return c.CardLevel
}

func (c MonsterCard) GetDebuffs() map[Ability]int {
	return c.DebuffMap
}

func (c MonsterCard) GetBuffs() map[Ability]int {
	return c.BuffMap
}

func (c MonsterCard) Clone() MonsterCard {
	clonedCard := MonsterCard{
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

func (c MonsterCard) AddAbilitiesWithArray(abilities []Ability) {
	for _, a := range abilities {
		c.Abilities = append(c.Abilities, a)
	}
}

// monster only
func (c MonsterCard) SetCardPosition(position int) {
	c.cardPosition = position
}

func (c MonsterCard) IsAlive() bool {
	return c.GameCard.Health >= 0
}

func (c MonsterCard) SetHealth(health int) {
	c.GameCard.Health = health
}

func (c MonsterCard) AddHealth(amount int) {
	if !c.IsAlive() {
		return
	}

	var finalHealth int
	postAbilityMaxHealth := c.GetPostAbilityMaxHealth()
	if c.GameCard.Health+amount > postAbilityMaxHealth {
		finalHealth = c.GameCard.Health + amount
	} else {
		finalHealth = postAbilityMaxHealth
	}

	if finalHealth > 0 {
		c.SetHealth(finalHealth)
	} else {
		c.SetHealth(0)
	}
}

func (c MonsterCard) HasBuff(buff Ability) bool {
	return utils.StrArrContains(c.GameCard.BuffMap, buff)
}

func (c MonsterCard) GetBuffCount(buff Ability) int {
	if val, ok := c.GameCard.BuffMap[buff]; ok {
		return val
	} else {
		return 0
	}
}

func (c MonsterCard) HasDebuff(buff Ability) bool {
	return utils.StrArrContains(c.GameCard.DebuffMap, buff)
}

func (c MonsterCard) GetDebuffCount(debuff Ability) int {
	if val, ok := c.GameCard.DebuffMap[debuff]; ok {
		return val
	} else {
		return 0
	}
}

func (c MonsterCard) GetIsLastStand() bool {
	return c.isOnlyMonster && c.HasAbility(ABILITY_LAST_STAND)
}

func (c MonsterCard) IsLastMonster() bool {
	return c.isOnlyMonster
}

func (c MonsterCard) SetIsOnlyMonster() {
	if c.HasAbility(ABILITY_LAST_STAND) {
		prevMaxHealth := c.GetPostAbilityMaxHealth()
		dmgTaken := prevMaxHealth - c.Health
		c.isOnlyMonster = true
		c.SetHealth(c.GetPostAbilityMaxHealth() - dmgTaken)
	}
	c.isOnlyMonster = true
}

func (c MonsterCard) HasAttack() bool {
	return c.Melee > 0 || c.Ranged > 0 || c.Magic > 0
}

func (c MonsterCard) GetPostAbilityMaxHealth() int {
	maxHealth := 1
	if c.GameCard.StartingHealth > maxHealth {
		maxHealth = c.GameCard.StartingHealth
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
		maxHealth = maxHealth - STRENGTHEN_AMOUNT*c.GetBuffCount(ABILITY_STRENGTHEN)
	}

	// The summoner skill made this starting health 0 or negative
	if c.StartingHealth < 1 {
		maxHealth = maxHealth + c.StartingHealth - 1
	}

	return utils.GetBigger(maxHealth, 1)
}

func (c MonsterCard) GetPostAbilityAttackOfType(attackType CardAttackType) int {
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
func (c MonsterCard) GetPostAbilityMagic() int {
	if c.Magic == 0 {
		return 0
	}
	postMagic := c.Magic
	if c.HasDebuff(ABILITY_HALVING) {
		postMagic = int(math.Floor((float64(postMagic) + 1) / 2))
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
func (c MonsterCard) GetPostAbilityRange() int {
	if c.Ranged == 0 {
		return 0
	}
	postRange := c.Ranged
	if c.HasDebuff(ABILITY_HALVING) {
		postRange = int(math.Floor((float64(postRange) + 1) / 2))
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
func (c MonsterCard) GetPostAbilityMelee() int {
	if c.Melee == 0 {
		return 0
	}
	postMelee := c.Melee
	if c.HasDebuff(ABILITY_HALVING) {
		postMelee = int(math.Floor((float64(postMelee) + 1) / 2))
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

	return utils.GetBigger(postMelee+meleeModifier, 1)
}

func (c MonsterCard) GetPostAbilitySpeed() int {
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
	if c.HasBuff(ABILITY_SLOW) {
		speedModifier = speedModifier - c.GetBuffCount(ABILITY_SLOW)
	}
	return utils.GetBigger(speed+speedModifier, 1)
}

func (c MonsterCard) GetPostAbilityMaxArmor() int {
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

func (c MonsterCard) RemoveDebuff(debuff Ability) {
	debuffAmount := c.GetDebuffCount(debuff) - 1
	if debuffAmount == 0 {
		delete(c.DebuffMap, debuff)
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

func (c MonsterCard) RemoveAllDebuff(debuff Ability) {
	debuffAmount := c.GetDebuffCount(debuff)
	for i := 0; i < debuffAmount; i++ {
		c.RemoveDebuff(debuff)
	}
}

func (c MonsterCard) CleanseDebuffs() {
	// Special case, cleanse only removes 1 cripple
	crippleAmount := c.GetDebuffCount(ABILITY_CRIPPLE)
	for ability, _ := range c.DebuffMap {
		if !utils.StrArrContains(utils.GetUncleansableDebuffs(), ability) {
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

func (c MonsterCard) Resurrect() {
	c.Health = 1
	if c.hadDivineShield {
		c.AddAbilitiesWithArray([]Ability{ABILITY_DIVINE_SHIELD})
	}
	c.Armor = c.GetPostAbilityMaxArmor()
	c.CleanseDebuffs()
}

func (c MonsterCard) IsEnraged() bool {
	return c.HasAbility(ABILITY_ENRAGE) && c.Health < c.GetPostAbilityMaxHealth()
}

func (c MonsterCard) IsPureMelee() bool {
	return c.Melee > 0 && c.Ranged == 0 && c.Magic == 0
}

func (c MonsterCard) CanMeleeAttack() bool {
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
func (c MonsterCard) GetSummonerArmor() int {
	return c.summonerArmor
}

func (c MonsterCard) AddSummonerArmor(stat int) {
	c.summonerArmor = c.summonerArmor + stat
	c.Armor = utils.GetBigger(c.Armor+stat, 1)
}

func (c MonsterCard) AddSummonerHealth(stat int) {
	c.StartingHealth = c.StartingHealth + stat
	c.Health = utils.GetBigger(c.Health+stat, 1)
}

func (c MonsterCard) AddSummonerSpeed(stat int) {
	c.summonerSpeed = c.summonerSpeed + stat
}

func (c MonsterCard) AddSummonerMelee(stat int) {
	c.summonerMelee = c.summonerMelee + stat
}

func (c MonsterCard) AddSummonerRanged(stat int) {
	c.summonerRanged = c.summonerRanged + stat
}

func (c MonsterCard) AddSummonerMagic(stat int) {
	c.summonerMagic = c.summonerMagic + stat
}

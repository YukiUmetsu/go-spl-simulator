package game_models

import utils "github.com/YukiUmetsu/go-spl-simulator/game_utils"

func DoRulesetPreGameBuff(rulesets []Ruleset, team1, team2 *GameTeam) {
	if utils.Contains(rulesets, RULESET_ARMORED_UP) {
		ApplyToBothTeamMonsters(team1, team2, ApplyArmorUpRuleset)
	}
	if utils.Contains(rulesets, RULESET_BACK_TO_BASICS) {
		ApplyToBothTeamMonsters(team1, team2, ApplyBackToBasicsRuleset)
	}
	if utils.Contains(rulesets, RULESET_CLOSE_RANGE) {
		ApplyToBothTeamMonsters(team1, team2, ApplyCloseRangeRuleset)
	}
	if utils.Contains(rulesets, RULESET_EQUAL_OPPORTUNITY) {
		ApplyToBothTeamMonsters(team1, team2, ApplyEqualOpportunityRuleset)
	}
	if utils.Contains(rulesets, RULESET_EQUALIZER) {
		ApplyEqualizer(team1, team2)
	}
	if utils.Contains(rulesets, RULESET_EXPLOSIVE_WEAPONRY) {
		ApplyToBothTeamMonsters(team1, team2, ApplyExplosiveWeaponRuleset)
	}
	if utils.Contains(rulesets, RULESET_FOG_OF_WAR) {
		ApplyToBothTeamMonsters(team1, team2, ApplyFogOfWarRuleset)
	}
	if utils.Contains(rulesets, RULESET_HEALED_OUT) {
		ApplyToBothTeamMonsters(team1, team2, ApplyHealedOutRuleset)
	}
	if utils.Contains(rulesets, RULESET_HEAVY_HITTERS) {
		ApplyToBothTeamMonsters(team1, team2, ApplyHeavyHittersRuleset)
	}
	if utils.Contains(rulesets, RULESET_HOLY_PROTECTION) {
		ApplyToBothTeamMonsters(team1, team2, ApplyHolyProtectionRuleset)
	}
	if utils.Contains(rulesets, RULESET_MELEE_MAYHEM) {
		ApplyToBothTeamMonsters(team1, team2, ApplyMeleeMayhemRuleset)
	}
	if utils.Contains(rulesets, RULESET_NOXIOUS_FUMES) {
		ApplyToBothTeamMonsters(team1, team2, ApplyNoxiousFumesRuleset)
	}
	if utils.Contains(rulesets, RULESET_SILENCED_SUMMONERS) {
		ApplySilencedSummonersRuleset(team1, team2)
	}
	if utils.Contains(rulesets, RULESET_SPREADING_FURY) {
		ApplyToBothTeamMonsters(team1, team2, ApplySpreadingFuryRuleset)
	}
	if utils.Contains(rulesets, RULESET_SUPER_SNEAK) {
		ApplyToBothTeamMonsters(team1, team2, ApplySuperSneakRuleset)
	}
	if utils.Contains(rulesets, RULESET_WEAK_MAGIC) {
		ApplyToBothTeamMonsters(team1, team2, ApplyWeakMagicRuleset)
	}
	if utils.Contains(rulesets, RULESET_TARGET_PRACTICE) {
		ApplyToBothTeamMonsters(team1, team2, ApplyTargetPracticeRuleset)
	}
}

func DoRulesetPreGamePostBuff(rulesets []Ruleset, team1, team2 *GameTeam) {
	if utils.Contains(rulesets, RULESET_UNPROTECTED) {
		ApplyToBothTeamMonsters(team1, team2, ApplyUnprotectedRuleset)
	}
}

func ApplyToBothTeamMonsters(team1, team2 *GameTeam, fn func(*MonsterCard)) {
	for _, m := range team1.GetMonstersList() {
		fn(m)
	}
	for _, m := range team2.GetMonstersList() {
		fn(m)
	}
}

/* All monsters gain 2 armors */
func ApplyArmorUpRuleset(m *MonsterCard) {
	m.AddSummonerArmor(2)
}

/* All monsters lose all of their abilities */
func ApplyBackToBasicsRuleset(m *MonsterCard) {
	m.RemoveAllAbilities()
}

/* All monsters have Close Range */
func ApplyCloseRangeRuleset(m *MonsterCard) {
	m.AddAbilitiesWithArray([]Ability{ABILITY_CLOSE_RANGE})
}

/* All monsters have Opportunity */
func ApplyEqualOpportunityRuleset(m *MonsterCard) {
	if !m.HasAbility(ABILITY_SNEAK) && !m.HasAbility(ABILITY_SNIPE) {
		m.AddAbilitiesWithArray([]Ability{ABILITY_OPPORTUNITY})
	}
}

func ApplyEqualizer(team1, team2 *GameTeam) {
	allMonsters := make([]*MonsterCard, 0)
	allMonsters = append(team1.GetMonstersList(), team2.GetMonstersList()...)
	highestHp := 0
	for _, m := range allMonsters {
		highestHp = utils.GetBigger(m.Health, highestHp)
	}
	for _, m := range allMonsters {
		m.Health = highestHp
		m.StartingHealth = highestHp
	}
}

/* All monsters have blast */
func ApplyExplosiveWeaponRuleset(m *MonsterCard) {
	m.AddAbilitiesWithArray([]Ability{ABILITY_BLAST})
}

/* No Sneak or Snipe */
func ApplyFogOfWarRuleset(m *MonsterCard) {
	m.RemoveAbility(ABILITY_SNEAK)
	m.RemoveAbility(ABILITY_SNIPE)
}

/* No healing abilities */
func ApplyHealedOutRuleset(m *MonsterCard) {
	m.RemoveAbility(ABILITY_TANK_HEAL)
	m.RemoveAbility(ABILITY_HEAL)
	m.RemoveAbility(ABILITY_TRIAGE)
}

/* All monsters have holy protection */
func ApplyHolyProtectionRuleset(m *MonsterCard) {
	m.AddAbilitiesWithArray([]Ability{ABILITY_DIVINE_SHIELD})
}

/* Monsters can attack from any position */
func ApplyMeleeMayhemRuleset(m *MonsterCard) {
	m.AddAbilitiesWithArray([]Ability{ABILITY_MELEE_MAYHEM})
}

/* All monsters poisoned */
func ApplyNoxiousFumesRuleset(m *MonsterCard) {
	m.AddDebuff(ABILITY_POISON)
}

/* Summoners don't do anything */
func ApplySilencedSummonersRuleset(team1, team2 *GameTeam) {
	silenceSummoner(team1.GetSummoner())
	silenceSummoner(team2.GetSummoner())
}

func silenceSummoner(summoner *SummonerCard) {
	summoner.RemoveAllAbilities()
	summoner.Health = 0
	summoner.Armor = 0
	summoner.Speed = 0
	summoner.Melee = 0
	summoner.Ranged = 0
	summoner.Magic = 0
}

/* All monsters has enrage */
func ApplySpreadingFuryRuleset(m *MonsterCard) {
	m.AddAbilitiesWithArray([]Ability{ABILITY_ENRAGE})
}

/* All Melee monsters have sneak */
func ApplySuperSneakRuleset(m *MonsterCard) {
	if m.Melee > 0 {
		m.AddAbilitiesWithArray([]Ability{ABILITY_SNEAK})
	}
}

/* All monsters have void armor */
func ApplyWeakMagicRuleset(m *MonsterCard) {
	m.AddAbilitiesWithArray([]Ability{ABILITY_VOID_ARMOR})
}

/* All ranged and magic have snipe */
func ApplyTargetPracticeRuleset(m *MonsterCard) {
	if m.Ranged > 0 || m.Magic > 0 {
		m.AddAbilitiesWithArray([]Ability{ABILITY_SNIPE})
	}
}

/* Monsters don't have armor */
func ApplyUnprotectedRuleset(m *MonsterCard) {
	m.Armor = 0
	m.StartingArmor = -99
}

func ApplyHeavyHittersRuleset(m *MonsterCard) {
	m.AddAbilitiesWithArray([]Ability{ABILITY_KNOCK_OUT})
}

func RulesetsContains(rulesets []Ruleset, ruleset Ruleset) bool {
	for _, r := range rulesets {
		if r == ruleset {
			return true
		}
	}
	return false
}

func ApplyEarthquake(g *Game, m *MonsterCard) BattleDamage {
	if g == nil || m == nil {
		return BattleDamage{}
	}

	if !m.HasAbility(ABILITY_FLYING) || m.HasDebuff(ABILITY_SNARE) {
		return HitMonsterWithPhysical(g, m, EARTHQUAKE_DAMAGE)
	}

	return BattleDamage{}
}

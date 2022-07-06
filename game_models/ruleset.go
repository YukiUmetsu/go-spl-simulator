package game_models

func DoRulesetPreGamePostBuff(rulesets []Ruleset, team1 GameTeam, team2 GameTeam) {
	if StrArrContains(rulesets, RULESET_ARMORED_UP) {
		ApplyToBothTeamMonsters(team1, team2, ApplyArmorUpRuleset)
	}
	if StrArrContains(rulesets, RULESET_BACK_TO_BASICS) {
		ApplyToBothTeamMonsters(team1, team2, ApplyBackToBasicsRuleset)
	}
	if StrArrContains(rulesets, RULESET_CLOSE_RANGE) {
		ApplyToBothTeamMonsters(team1, team2, ApplyCloseRangeRuleset)
	}
	if StrArrContains(rulesets, RULESET_EQUAL_OPPORTUNITY) {
		ApplyToBothTeamMonsters(team1, team2, ApplyEqualOpportunityRuleset)
	}
	if StrArrContains(rulesets, RULESET_EQUALIZER) {
		ApplyEqualizer(team1, team2)
	}
	if StrArrContains(rulesets, RULESET_EXPLOSIVE_WEAPONRY) {
		ApplyToBothTeamMonsters(team1, team2, ApplyExplosiveWeaponRuleset)
	}
	if StrArrContains(rulesets, RULESET_FOG_OF_WAR) {
		ApplyToBothTeamMonsters(team1, team2, ApplyFogOfWarRuleset)
	}
	if StrArrContains(rulesets, RULESET_HEALED_OUT) {
		ApplyToBothTeamMonsters(team1, team2, ApplyHealedOutRuleset)
	}
	if StrArrContains(rulesets, RULESET_HEAVY_HITTERS) {
		ApplyToBothTeamMonsters(team1, team2, ApplyHeavyHittersRuleset)
	}
	if StrArrContains(rulesets, RULESET_HOLY_PROTECTION) {
		ApplyToBothTeamMonsters(team1, team2, ApplyHolyProtectionRuleset)
	}
	if StrArrContains(rulesets, RULESET_MELEE_MAYHEM) {
		ApplyToBothTeamMonsters(team1, team2, ApplyMeleeMayhemRuleset)
	}
	if StrArrContains(rulesets, RULESET_NOXIOUS_FUMES) {
		ApplyToBothTeamMonsters(team1, team2, ApplyNoxiousFumesRuleset)
	}
	if StrArrContains(rulesets, RULESET_SILENCED_SUMMONERS) {
		ApplySilencedSummonersRuleset(team1, team2)
	}
	if StrArrContains(rulesets, RULESET_SPREADING_FURY) {
		ApplyToBothTeamMonsters(team1, team2, ApplySpreadingFuryRuleset)
	}
	if StrArrContains(rulesets, RULESET_SUPER_SNEAK) {
		ApplyToBothTeamMonsters(team1, team2, ApplySuperSneakRuleset)
	}
	if StrArrContains(rulesets, RULESET_WEAK_MAGIC) {
		ApplyToBothTeamMonsters(team1, team2, ApplyWeakMagicRuleset)
	}
	if StrArrContains(rulesets, RULESET_TARGET_PRACTICE) {
		ApplyToBothTeamMonsters(team1, team2, ApplyTargetPracticeRuleset)
	}
}

func ApplyToBothTeamMonsters(team1 GameTeam, team2 GameTeam, fn func(MonsterCard)) {
	for _, m := range team1.GetMonstersList() {
		fn(m)
	}
	for _, m := range team2.GetMonstersList() {
		fn(m)
	}
}

/* All monsters gain 2 armors */
func ApplyArmorUpRuleset(m MonsterCard) {
	m.AddSummonerArmor(2)
}

/* All monsters lose all of their abilities */
func ApplyBackToBasicsRuleset(m MonsterCard) {
	m.RemoveAllAbilities()
}

/* All monsters have Close Range */
func ApplyCloseRangeRuleset(m MonsterCard) {
	m.AddAbility(ABILITY_CLOSE_RANGE)
}

/* All monsters have Opportunity */
func ApplyEqualOpportunityRuleset(m MonsterCard) {
	if !m.HasAbility(ABILITY_SNEAK) && !m.HasAbility(ABILITY_SNIPE) {
		m.AddAbility(ABILITY_OPPORTUNITY)
	}
}

func ApplyEqualizer(team1 MonsterCard, team2 MonsterCard) {
	allMonsters := make([]MonsterCard, 0)
	allMonsters = append(allMonsters, team1.GetMonstersList(), team2.GetMonstersList())
	highestHp := 0
	for _, m := range allMonsters {
		highestHp = GetBigger(m.Health, highestHp)
	}
	for _, m := range allMonsters {
		m.Health = highestHp
		m.StartingHealth = highestHp
	}
}

/* All monsters have blast */
func ApplyExplosiveWeaponRuleset(m MonsterCard) {
	m.AddAbility(ABILITY_BLAST)
}

/* No Sneak or Snipe */
func ApplyFogOfWarRuleset(m MonsterCard) {
	m.RemoveAbility(ABILITY_SNEAK)
	m.RemoveAbility(ABILITY_SNIPE)
}

/* No healing abilities */
func ApplyHealedOutRuleset(m MonsterCard) {
	m.RemoveAbility(ABILITY_TANK_HEAL)
	m.RemoveAbility(ABILITY_HEAL)
	m.RemoveAbility(ABILITY_TRIAGE)
}

/* All monsters have holy protection */
func ApplyHolyProtectionRuleset(m MonsterCard) {
	m.AddAbility(ABILITY_DIVINE_SHIELD)
}

/* Monsters can attack from any position */
func ApplyMeleeMayhemRuleset(m MonsterCard) {
	m.AddAbility(ABILITY_MELEE_MAYHEM)
}

/* All monsters poisoned */
func ApplyNoxiousFumesRuleset(m MonsterCard) {
	m.AddDebuff(ABILITY_POISON)
}

/* Summoners don't do anything */
func ApplySilencedSummonersRuleset(team1 MonsterCard, team2 MonsterCard) {
	silenceSummoner(team1.GetSummoner())
	silenceSummoner(team2.GetSummoner())
}

func silenceSummoner(summoner SummonerCard) {
	summoner.RemoveAllAbilities()
	summoner.Health = 0
	summoner.Armor = 0
	summoner.Speed = 0
	summoner.Melee = 0
	summoner.Ranged = 0
	summoner.Magic = 0
}

/* All monsters has enrage */
func ApplySpreadingFuryRuleset(m MonsterCard) {
	m.AddAbility(ABILITY_ENRAGE)
}

/* All Melee monsters have sneak */
func ApplySuperSneakRuleset(m MonsterCard) {
	if m.Melee > 0 {
		m.AddAbility(ABILITY_SNEAK)
	}
}

/* All monsters have void armor */
func ApplyWeakMagicRuleset(m MonsterCard) {
	m.AddAbility(ABILITY_VOID_ARMOR)
}

/* All ranged and magic have snipe */
func ApplyTargetPracticeRuleset(m MonsterCard) {
	if m.Ranged > 0 || m.Magic > 0 {
		m.AddAbility(ABILITY_SNIPE)
	}
}

/* Monsters don't have armor */
func ApplyUnprotectedRuleset(m MonsterCard) {
	m.Armor = 0
	m.StartingArmor = -99
}

func ApplyHeavyHittersRuleset(m MonsterCard) {
	m.AddAbility(ABILITY_KNOCK_OUT)
}

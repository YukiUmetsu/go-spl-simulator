package game_utils

import (
	sim "simulator"
)

func DoRulesetPreGamePostBuff(rulesets []sim.Ruleset, team1 sim.GameTeam, team2 sim.GameTeam) {
	if StrArrContains(rulesets, sim.RULESET_ARMORED_UP) {
		ApplyToBothTeamMonsters(team1, team2, ApplyArmorUpRuleset)
	}
	if StrArrContains(rulesets, sim.RULESET_BACK_TO_BASICS) {
		ApplyToBothTeamMonsters(team1, team2, ApplyBackToBasicsRuleset)
	}
	if StrArrContains(rulesets, sim.RULESET_CLOSE_RANGE) {
		ApplyToBothTeamMonsters(team1, team2, ApplyCloseRangeRuleset)
	}
	if StrArrContains(rulesets, sim.RULESET_EQUAL_OPPORTUNITY) {
		ApplyToBothTeamMonsters(team1, team2, ApplyEqualOpportunityRuleset)
	}
	if StrArrContains(rulesets, sim.RULESET_EQUALIZER) {
		ApplyEqualizer(team1, team2)
	}
	if StrArrContains(rulesets, sim.RULESET_EXPLOSIVE_WEAPONRY) {
		ApplyToBothTeamMonsters(team1, team2, ApplyExplosiveWeaponRuleset)
	}
	if StrArrContains(rulesets, sim.RULESET_FOG_OF_WAR) {
		ApplyToBothTeamMonsters(team1, team2, ApplyFogOfWarRuleset)
	}
	if StrArrContains(rulesets, sim.RULESET_HEALED_OUT) {
		ApplyToBothTeamMonsters(team1, team2, ApplyHealedOutRuleset)
	}
	if StrArrContains(rulesets, sim.RULESET_HEAVY_HITTERS) {
		ApplyToBothTeamMonsters(team1, team2, ApplyHeavyHittersRuleset)
	}
	if StrArrContains(rulesets, sim.RULESET_HOLY_PROTECTION) {
		ApplyToBothTeamMonsters(team1, team2, ApplyHolyProtectionRuleset)
	}
	if StrArrContains(rulesets, sim.RULESET_MELEE_MAYHEM) {
		ApplyToBothTeamMonsters(team1, team2, ApplyMeleeMayhemRuleset)
	}
	if StrArrContains(rulesets, sim.RULESET_NOXIOUS_FUMES) {
		ApplyToBothTeamMonsters(team1, team2, ApplyNoxiousFumesRuleset)
	}
	if StrArrContains(rulesets, sim.RULESET_SILENCED_SUMMONERS) {
		ApplySilencedSummonersRuleset(team1, team2)
	}
	if StrArrContains(rulesets, sim.RULESET_SPREADING_FURY) {
		ApplyToBothTeamMonsters(team1, team2, ApplySpreadingFuryRuleset)
	}
	if StrArrContains(rulesets, sim.RULESET_SUPER_SNEAK) {
		ApplyToBothTeamMonsters(team1, team2, ApplySuperSneakRuleset)
	}
	if StrArrContains(rulesets, sim.RULESET_WEAK_MAGIC) {
		ApplyToBothTeamMonsters(team1, team2, ApplyWeakMagicRuleset)
	}
	if StrArrContains(rulesets, sim.RULESET_TARGET_PRACTICE) {
		ApplyToBothTeamMonsters(team1, team2, ApplyTargetPracticeRuleset)
	}
}

func ApplyToBothTeamMonsters(team1 sim.GameTeam, team2 sim.GameTeam, fn func(sim.MonsterCard)) {
	for _, m := range team1.GetMonstersList() {
		fn(m)
	}
	for _, m := range team2.GetMonstersList() {
		fn(m)
	}
}

/* All monsters gain 2 armors */
func ApplyArmorUpRuleset(m sim.MonsterCard) {
	m.AddSummonerArmor(2)
}

/* All monsters lose all of their abilities */
func ApplyBackToBasicsRuleset(m sim.MonsterCard) {
	m.RemoveAllAbilities()
}

/* All monsters have Close Range */
func ApplyCloseRangeRuleset(m sim.MonsterCard) {
	m.AddAbility(sim.ABILITY_CLOSE_RANGE)
}

/* All monsters have Opportunity */
func ApplyEqualOpportunityRuleset(m sim.MonsterCard) {
	if !m.HasAbility(sim.ABILITY_SNEAK) && !m.HasAbility(sim.ABILITY_SNIPE) {
		m.AddAbility(sim.ABILITY_OPPORTUNITY)
	}
}

func ApplyEqualizer(team1 sim.MonsterCard, team2 sim.MonsterCard) {
	allMonsters := make([]sim.MonsterCard, 0)
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
func ApplyExplosiveWeaponRuleset(m sim.MonsterCard) {
	m.AddAbility(sim.ABILITY_BLAST)
}

/* No Sneak or Snipe */
func ApplyFogOfWarRuleset(m sim.MonsterCard) {
	m.RemoveAbility(sim.ABILITY_SNEAK)
	m.RemoveAbility(sim.ABILITY_SNIPE)
}

/* No healing abilities */
func ApplyHealedOutRuleset(m sim.MonsterCard) {
	m.RemoveAbility(sim.ABILITY_TANK_HEAL)
	m.RemoveAbility(sim.ABILITY_HEAL)
	m.RemoveAbility(sim.ABILITY_TRIAGE)
}

/* All monsters have holy protection */
func ApplyHolyProtectionRuleset(m sim.MonsterCard) {
	m.AddAbility(sim.ABILITY_DIVINE_SHIELD)
}

/* Monsters can attack from any position */
func ApplyMeleeMayhemRuleset(m sim.MonsterCard) {
	m.AddAbility(sim.ABILITY_MELEE_MAYHEM)
}

/* All monsters poisoned */
func ApplyNoxiousFumesRuleset(m sim.MonsterCard) {
	m.AddDebuff(sim.ABILITY_POISON)
}

/* Summoners don't do anything */
func ApplySilencedSummonersRuleset(team1 sim.MonsterCard, team2 sim.MonsterCard) {
	silenceSummoner(team1.GetSummoner())
	silenceSummoner(team2.GetSummoner())
}

func silenceSummoner(summoner sim.SummonerCard) {
	summoner.RemoveAllAbilities()
	summoner.Health = 0
	summoner.Armor = 0
	summoner.Speed = 0
	summoner.Melee = 0
	summoner.Ranged = 0
	summoner.Magic = 0
}

/* All monsters has enrage */
func ApplySpreadingFuryRuleset(m sim.MonsterCard) {
	m.AddAbility(sim.ABILITY_ENRAGE)
}

/* All Melee monsters have sneak */
func ApplySuperSneakRuleset(m sim.MonsterCard) {
	if m.Melee > 0 {
		m.AddAbility(sim.ABILITY_SNEAK)
	}
}

/* All monsters have void armor */
func ApplyWeakMagicRuleset(m sim.MonsterCard) {
	m.AddAbility(sim.ABILITY_VOID_ARMOR)
}

/* All ranged and magic have snipe */
func ApplyTargetPracticeRuleset(m sim.MonsterCard) {
	if m.Ranged > 0 || m.Magic > 0 {
		m.AddAbility(sim.ABILITY_SNIPE)
	}
}

/* Monsters don't have armor */
func ApplyUnprotectedRuleset(m sim.MonsterCard) {
	m.Armor = 0
	m.StartingArmor = -99
}

func ApplyHeavyHittersRuleset(m sim.MonsterCard) {
	m.AddAbility(sim.ABILITY_KNOCK_OUT)
}

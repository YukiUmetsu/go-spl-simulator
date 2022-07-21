package game_models

import (
	"errors"
	"fmt"
	"math"
	"sort"

	utils "github.com/YukiUmetsu/go-spl-simulator/game_utils"
)

type Game struct {
	team1      *GameTeam
	team2      *GameTeam
	rulesets   []Ruleset
	battleLogs []BattleLog
	shouldLog  bool
	/* 0: unknown, 1: team1, 2: team2, 3: Tie */
	winner       TeamNumber
	deadMonsters []*MonsterCard
	roundNumber  int
	stunData     map[string][]*MonsterCard // key: "[team number]-[monster name]" e.g. "1-Magnor"
}

func (g *Game) Create(team1, team2 *GameTeam, rulesets []Ruleset, shouldLog bool) {
	g.team1 = team1
	g.team2 = team2
	g.rulesets = rulesets
	g.shouldLog = shouldLog
	g.team1.SetTeamNumber(TEAM_NUM_ONE)
	g.team2.SetTeamNumber(TEAM_NUM_TWO)
	g.stunData = make(map[string][]*MonsterCard, 0)
}

func (g *Game) Reset() {
	g.roundNumber = 0
	g.winner = TEAM_NUM_UNKNOWN
	g.deadMonsters = make([]*MonsterCard, 0)
	g.team1.ResetTeam()
	g.team2.ResetTeam()
	g.stunData = make(map[string][]*MonsterCard, 0)
	g.battleLogs = []BattleLog{}
}

func (g *Game) GetWinner() TeamNumber {
	return g.winner
}

func (g *Game) RemoveStunsThatThisMonsterApplied(m *MonsterCard) {
	stunDataKey := g.GetStunDataKey(g.roundNumber, m)
	if _, ok := g.stunData[stunDataKey]; ok {
		// the monster stunned a monster previously
		hadStunMonsters := g.stunData[stunDataKey]
		if hadStunMonsters != nil && len(hadStunMonsters) > 0 {
			for _, stunnedMonster := range hadStunMonsters {
				stunnedMonster.RemoveAllDebuff(ABILITY_STUN)
				delete(g.stunData, stunDataKey)
				g.CreateAndAddBattleLog(BATTLE_ACTION_STUN_REMOVED, m, stunnedMonster, 0)
			}
		}
	}
}

func (g *Game) GetStunDataKey(roundNumber int, stunApplier GameCardInterface) string {
	return fmt.Sprintf("%d-%v", roundNumber, stunApplier.GetName())
}

func (g *Game) GetBattleLogs() ([]BattleLog, error) {
	if !g.shouldLog {
		return []BattleLog{}, errors.New("you must instantiate the game with enableLogs as true")
	}
	return g.battleLogs, nil
}

func (g *Game) PlayGame() {
	g.Reset()
	team1Summoner := g.team1.GetSummoner()
	team1Monsters := g.team1.GetMonstersList()
	team2Summoner := g.team2.GetSummoner()
	team2Monsters := g.team2.GetMonstersList()

	// pre game rulesets
	DoRulesetPreGameBuff(g.rulesets, g.team1, g.team2)

	// Summoner pre-game buffs
	g.DoSummonerPreGameBuff(team1Summoner, team1Monsters)
	g.DoSummonerPreGameBuff(team2Summoner, team2Monsters)

	// Summoner pre-game debuffs
	g.DoSummonerPreGameDebuff(team1Summoner, team2Monsters)
	g.DoSummonerPreGameDebuff(team2Summoner, team1Monsters)

	// Monsters pre-game buffs
	g.DoMonsterPreGameBuff(team1Monsters)
	g.DoMonsterPreGameBuff(team2Monsters)

	// Monsters pre-game debuffs
	g.DoMonsterPreGameDebuff(team1Monsters, team2Monsters)

	// Apply ruleset rules that apply post buff phase
	DoRulesetPreGamePostBuff(g.rulesets, g.team1, g.team2)

	g.team1.SetAllMonsterHealth()
	g.team2.SetAllMonsterHealth()

	g.PlayRoundsUntilGameEnd(0)
}

func (g *Game) DoSummonerPreGameBuff(summoner *SummonerCard, friendlyMonsters []*MonsterCard) {
	// add summoner abilities that increase stats (aka, strengthen)
	for _, ability := range GetSummonerPreGameBuffAbilities() {
		if summoner.HasAbility(ability) {
			g.ApplyBuffToMonsters(friendlyMonsters, ability)
		}
	}

	// add summoner abilities
	for _, ability := range GetSummonerAbilityAbilities() {
		if summoner.HasAbility(ability) {
			g.ApplyAbilityToMonsters(friendlyMonsters, ability)
			logAction := "Summoner Pre Game Ability Buff " + string(ability)
			for _, m := range friendlyMonsters {
				g.CreateAndAddBattleLog(AdditionalBattleAction(logAction), summoner, m, 0)
			}
		}
	}

	// add summoner stats (e.g. +1 melee, +1 archery, +1 magic etc...)
	for _, m := range friendlyMonsters {
		if summoner.Armor > 0 {
			g.CreateAndAddBattleLog(AdditionalBattleAction("Summoner Armor Buff"), summoner, m, summoner.Armor)
			m.AddSummonerArmor(summoner.Armor)
		}
		if summoner.Health > 0 {
			g.CreateAndAddBattleLog(AdditionalBattleAction("Summoner Health Buff"), summoner, m, summoner.Health)
			m.AddSummonerHealth(summoner.Health)
		}
		if summoner.Speed > 0 {
			g.CreateAndAddBattleLog(AdditionalBattleAction("Summoner Speed Buff"), summoner, m, summoner.Speed)
			m.AddSummonerSpeed(summoner.Speed)
		}
		if summoner.Melee > 0 {
			g.CreateAndAddBattleLog(AdditionalBattleAction("Summoner Melee Buff"), summoner, m, summoner.Melee)
			m.AddSummonerMelee(summoner.Melee)
		}
		if summoner.Ranged > 0 {
			g.CreateAndAddBattleLog(AdditionalBattleAction("Summoner Ranged Buff"), summoner, m, summoner.Ranged)
			m.AddSummonerRanged(summoner.Ranged)
		}
		if summoner.Magic > 0 {
			g.CreateAndAddBattleLog(AdditionalBattleAction("Summoner Magic Buff"), summoner, m, summoner.Magic)
			m.AddSummonerMagic(summoner.Magic)
		}
	}
}

// Add all summoner abilities to all enemy monsters which are in SUMMONER_DEBUFF_ABILITIES
func (g *Game) DoSummonerPreGameDebuff(summoner *SummonerCard, targetMonsters []*MonsterCard) {
	// add summoner debuffs (aka, affliciton, blind)
	for _, debuff := range GetSummonerPreGameDebuffAbilities() {
		if summoner.HasAbility(debuff) {
			g.ApplyDebuffToMonsters(targetMonsters, debuff)
		}
	}

	// add summoner stats (e.g. -1 melee, -1 archery, -1 magic etc...)
	for _, m := range targetMonsters {
		if summoner.Armor < 0 {
			m.AddSummonerArmor(summoner.Armor)
		}
		if summoner.Health < 0 {
			m.AddSummonerHealth(summoner.Health)
		}
		if summoner.Speed < 0 {
			m.AddSummonerSpeed(summoner.Speed)
		}
		if summoner.Melee < 0 {
			m.AddSummonerMelee(summoner.Melee)
		}
		if summoner.Ranged < 0 {
			m.AddSummonerRanged(summoner.Ranged)
		}
		if summoner.Magic < 0 {
			m.AddSummonerMagic(summoner.Magic)
		}
	}
}

// Add all monster abilities to all friendly monsters which are in MONSTER_BUFF_ABILITIES
func (g *Game) DoMonsterPreGameBuff(friendlyMonsters []*MonsterCard) {
	for _, m := range friendlyMonsters {
		for _, buff := range GetMonsterPreGameBuffAbilities() {
			if !m.HasAbility(buff) {
				continue
			}

			g.ApplyBuffToMonsters(friendlyMonsters, buff)
			for _, fm := range friendlyMonsters {
				g.CreateAndAddBattleLog(AdditionalBattleAction("Monster Pre-Game Buff "+string(buff)), m, fm, 0)
			}
		}
	}
}

func (g *Game) DoMonsterPreGameDebuff(team1Monsters []*MonsterCard, team2Monsters []*MonsterCard) {
	for _, debuff := range GetMonsterPreGameDebuffAbilities() {
		// team1 debuff team2
		for _, m := range team1Monsters {
			if !m.HasAbility(debuff) {
				continue
			}
			g.ApplyDebuffToMonsters(team2Monsters, debuff)

			// battle log
			for _, m2 := range team2Monsters {
				g.CreateAndAddBattleLog(AdditionalBattleAction("Monster Pre-Game Debuff "+string(debuff)), m, m2, 0)
			}
		}

		// team2 debuff team1
		for _, m := range team2Monsters {
			if !m.HasAbility(debuff) {
				continue
			}
			g.ApplyDebuffToMonsters(team1Monsters, debuff)

			// battle log
			for _, m1 := range team1Monsters {
				g.CreateAndAddBattleLog(AdditionalBattleAction("Monster Pre-Game Debuff "+string(debuff)), m, m1, 0)
			}
		}
	}
}

func (g *Game) ApplyBuffToMonsters(monsters []*MonsterCard, buff Ability) {
	// add stats buff
	for _, m := range monsters {
		m.AddBuff(buff)
	}
}

func (g *Game) ApplyAbilityToMonsters(monsters []*MonsterCard, ability Ability) {
	for _, m := range monsters {
		m.AddAbilitiesWithArray([]Ability{ability})
	}
}

func (g *Game) ApplyDebuffToMonsters(monsters []*MonsterCard, debuff Ability) {
	for _, m := range monsters {
		m.AddDebuff(debuff)
	}
}

func (g *Game) GetAllAliveMonsters() []*MonsterCard {
	aliveMonsters := make([]*MonsterCard, 0)
	t1 := g.team1.GetAliveMonsters()
	t2 := g.team2.GetAliveMonsters()
	aliveMonsters = append(t1, t2...)
	return aliveMonsters
}

// Plays the game rounds until the game is over
func (g *Game) PlayRoundsUntilGameEnd(roundNumber int) {
	g.roundNumber = roundNumber
	if g.winner != TEAM_NUM_UNKNOWN {
		g.LogGameOver()
		return
	}

	// if round >= 50, game is tie
	if roundNumber > 50 {
		g.winner = TEAM_NUM_TIE
		g.LogGameOver()
	}

	// Fatigue
	if roundNumber >= FATIGUE_ROUND_NUMBER {
		g.FatigueMonsters(roundNumber)
		g.CheckAndSetGameWinner()
		if g.winner != TEAM_NUM_UNKNOWN {
			g.LogGameOver()
			return
		}
	}

	// Play one round
	g.PlaySingleRound()
	g.CheckAndSetGameWinner()
	if g.winner != TEAM_NUM_UNKNOWN {
		g.LogGameOver()
		return
	}

	// Post round including earthquake
	g.DoPostRound()
	g.CheckAndSetGameWinner()

	g.PlayRoundsUntilGameEnd(roundNumber + 1)
}

func (g *Game) LogGameOver() {
	if g.winner == TEAM_NUM_UNKNOWN {
		return
	}

	winnerTeamLabel := ""
	if g.winner == TEAM_NUM_ONE {
		winnerTeamLabel = "1 - " + g.team1.GetPlayerName()
	} else if g.winner == TEAM_NUM_TWO {
		winnerTeamLabel = "2 - " + g.team2.GetPlayerName()
	} else {
		winnerTeamLabel = "TIE"
	}
	g.CreateAndAddBattleLog(AdditionalBattleAction(string(BATTLE_ACTION_GAME_OVER)+" Winner: "+winnerTeamLabel), nil, nil, int(g.winner))
}

func (g *Game) FatigueMonsters(roundNumber int) {
	fatigueDamage := roundNumber - FATIGUE_ROUND_NUMBER + 1
	allAliveMonsters := g.GetAllAliveMonsters()

	for _, m := range allAliveMonsters {
		g.CreateAndAddBattleLog(BATTLE_ACTION_FATIGUE, m, nil, fatigueDamage)
		m.HitHealth(fatigueDamage)
		g.ProcessIfDead(m)
	}

	g.CheckAndSetGameWinner()
	if g.winner != TEAM_NUM_UNKNOWN {
		g.LogGameOver()
		return
	}
}

func (g *Game) CreateAndAddBattleLog(action Stringer, cardOne GameCardInterface, cardTwo GameCardInterface, value int) {
	if !g.shouldLog {
		return
	}

	var actor GameCardInterface
	var target GameCardInterface
	if cardOne != nil {
		actor = cardOne.Clone()
	}
	if cardTwo != nil {
		target = cardTwo.Clone()
	}

	log := BattleLog{
		Actor:  actor,
		Action: action,
		Target: target,
		Value:  value,
	}
	g.battleLogs = append(g.battleLogs, log)
}

func (g *Game) CheckAndSetGameWinner() {
	team1AliveMonstersCount := len(g.team1.GetAliveMonsters())
	team2AliveMonstersCount := len(g.team2.GetAliveMonsters())

	if team1AliveMonstersCount == 1 && team2AliveMonstersCount == 1 {
		g.winner = TEAM_NUM_UNKNOWN
	} else if team2AliveMonstersCount == 0 {
		g.winner = TEAM_NUM_ONE
	} else if team1AliveMonstersCount == 0 {
		g.winner = TEAM_NUM_TWO
	}
}

func (g *Game) ProcessIfDead(m *MonsterCard) {
	if m == nil || m.IsAlive() {
		return
	}

	// TODO: [ProcessIfDead] this is not right. even if stun applier dies, stun won't go away
	g.RemoveStunsThatThisMonsterApplied(m)

	// monster is dead
	g.CreateAndAddBattleLog(BATTLE_ACTION_DEATH, m, nil, 0)
	g.deadMonsters = append(g.deadMonsters, m)
	m.SetHasTurnPassed(true)

	friendlyTeam := g.GetTeamOfMonster(m)
	aliveFriendlyMonsters := friendlyTeam.GetAliveMonsters()
	enemyTeam := g.GetEnemyTeamOfMonster(m)
	aliveEnemyMonsters := enemyTeam.GetAliveMonsters()

	// Redemption
	if m.HasAbility(ABILITY_REDEMPTION) {
		for _, e := range aliveEnemyMonsters {
			HitMonsterWithPhysical(
				g,
				e,
				REDEMPTION_DAMAGE,
			)

			g.ProcessIfDead(e)
		}
	}

	// Ressurect
	friendlySummoner := friendlyTeam.GetSummoner()
	// summoner resurrect
	wasResurrected := g.ProcessIfResurrect(friendlySummoner, m)
	// friendly monster resurrect
	for _, friendlyMonster := range aliveFriendlyMonsters {
		if wasResurrected {
			break
		}
		wasResurrected = g.ProcessIfResurrect(friendlyMonster, m)
	}

	// remove debuffs and buffs if not resurrected
	if !wasResurrected {
		// remove debuffs from the enemy team
		monsterDebuffs := MonsterHasDebuffAbilities(m)
		g.RemoveMonsterDebuff(m, monsterDebuffs, aliveEnemyMonsters)

		// remove buffs from the friendly monsters
		monsterBuffs := MonsterHasBuffsAbilities(m)
		g.RemoveMonsterBuff(m, monsterBuffs, aliveFriendlyMonsters)
	}

	// handle scavenger & battle log
	for _, fm := range aliveFriendlyMonsters {
		g.OnMonsterDeath(fm, m)
	}
	for _, em := range aliveEnemyMonsters {
		g.OnMonsterDeath(em, m)
	}

	if !wasResurrected {
		friendlyTeam.MaybeSetLastStand()
	}
}

/**
* Plays a single round
* 1. Summoners do their pre round abilities.
* 2. Get turn order of monsters
* 3. For each monster, check if alive then
* 3a. Do onPreTurn
* 3b. Get target, continue to next monster if no attack
* 3c. Attack target
* 3d. Resolve damage, check if dead monsters
* 3e. (If dead) Trigger onDeath on all alive monsters and summoners
 */
func (g *Game) PlaySingleRound() {
	// pre round buffs etc
	g.CreateAndAddBattleLog(BATTLE_ACTION_ROUND_START, nil, nil, g.roundNumber+1)
	g.deadMonsters = []*MonsterCard{}
	g.DoGamePreRound()
	g.DoSummonerPreRound(g.team1)
	g.DoSummonerPreRound(g.team2)

	// loop through each monster's turn
	currentMonster := g.GetNextMonsterTurn()
	for currentMonster != nil {
		if !currentMonster.IsAlive() {
			continue
		}

		// remove stun state
		g.RemoveStunsThatThisMonsterApplied(currentMonster)

		// check stun
		if currentMonster.HasDebuff(ABILITY_STUN) {
			currentMonster.SetHasTurnPassed(true)
			currentMonster = g.GetNextMonsterTurn()
			continue
		}

		// handle monster attack
		g.DoMonsterPreTurn(currentMonster)
		g.ResolveAttackForMonster(currentMonster)
		if currentMonster.HasAbility(ABILITY_DOUBLE_STRIKE) {
			g.ResolveAttackForMonster(currentMonster)
		}

		currentMonster = g.GetNextMonsterTurn()
	}
}

func (g *Game) DoGamePreRound() {
	// Add pre round actions when necessary
}

// Handle Summoner's pre turn actions (e.g. cleanse, tank heal, repair, triage)
func (g *Game) DoSummonerPreRound(t *GameTeam) {
	summoner := t.GetSummoner()

	// Cleanse
	if summoner.HasAbility(ABILITY_CLEANSE) {
		firstMonster := t.GetFirstAliveMonster()
		firstMonster.CleanseDebuffs()
		g.CreateAndAddBattleLog(BATTLE_ACTION_CLEANSE, summoner, firstMonster, 0)
	}

	// Repair
	if summoner.HasAbility(ABILITY_REPAIR) {
		repairTarget := t.GetRepairTarget()
		if repairTarget != nil {
			RepairMonsterArmor(repairTarget)
			g.CreateAndAddBattleLog(BATTLE_ACTION_REPAIR, summoner, repairTarget, 0)
		}
	}

	// Tank heal
	if summoner.HasAbility(ABILITY_TANK_HEAL) {
		firstMonster := t.GetFirstAliveMonster()
		TankHealMonster(firstMonster)
		g.CreateAndAddBattleLog(BATTLE_ACTION_TANK_HEAL, summoner, firstMonster, 0)
	}

	// Triage
	if summoner.HasAbility(ABILITY_TRIAGE) {
		healTarget := t.GetTriageHealTarget()
		if healTarget != nil {
			healAmount := TriageHealMonster(healTarget)
			g.CreateAndAddBattleLog(BATTLE_ACTION_TRIAGE, summoner, healTarget, healAmount)
		}
	}
}

// TODO Test  Get the next monster that should attack
func (g *Game) GetNextMonsterTurn() *MonsterCard {
	allUnmovedMonsters := append(g.team1.GetUnmovedMonsters(), g.team2.GetUnmovedMonsters()...)
	if len(allUnmovedMonsters) == 0 {
		return nil
	}

	// sort unmoved monsters
	sort.SliceStable(allUnmovedMonsters, func(i, j int) bool {
		m1 := allUnmovedMonsters[i]
		m2 := allUnmovedMonsters[j]
		normalCompareDiff := NormalCompareAttackOrder(m1, m2)

		// Descending order
		if normalCompareDiff != 0 {
			return normalCompareDiff < 0
		}

		// resolve tie by order if the same team, else random
		if m1.GetTeamNumber() == m2.GetTeamNumber() {
			return ResolveFriendlyTies(m1, m2) > 0
		} else {
			return RandomTieBreaker() > 0
		}
	})

	if utils.Contains(g.rulesets, RULESET_REVERSE_SPEED) {
		return allUnmovedMonsters[0]
	}

	return allUnmovedMonsters[len(allUnmovedMonsters)-1]
}

func (g *Game) DoMonsterPreTurn(m *MonsterCard) {
	m.SetHasTurnPasses(true)
	friendlyTeam := g.GetTeamOfMonster(m)

	// Cleanse
	if m.HasAbility(ABILITY_CLEANSE) {
		cleanseTarget := friendlyTeam.GetFirstAliveMonster()
		cleanseTarget.CleanseDebuffs()
		g.CreateAndAddBattleLog(BATTLE_ACTION_CLEANSE, m, cleanseTarget, 0)
	}

	// Tank heal
	if m.HasAbility(ABILITY_TANK_HEAL) {
		tankHealTarget := friendlyTeam.GetFirstAliveMonster()
		healAmount := TankHealMonster(tankHealTarget)
		g.CreateAndAddBattleLog(BATTLE_ACTION_TANK_HEAL, m, tankHealTarget, healAmount)
	}

	// Repair
	if m.HasAbility(ABILITY_REPAIR) {
		repairTarget := friendlyTeam.GetRepairTarget()
		if repairTarget != nil {
			repairAmount := RepairMonsterArmor(repairTarget)
			g.CreateAndAddBattleLog(BATTLE_ACTION_REPAIR, m, repairTarget, repairAmount)
		}
	}

	// Triage
	if m.HasAbility(ABILITY_TRIAGE) {
		triageTarget := friendlyTeam.GetTriageHealTarget()
		if triageTarget != nil {
			triageAmount := TriageHealMonster(triageTarget)
			g.CreateAndAddBattleLog(BATTLE_ACTION_TRIAGE, m, triageTarget, triageAmount)
		}
	}

	// Self heal
	if m.HasAbility(ABILITY_HEAL) {
		healAmount := SelfHealMonster(m)
		g.CreateAndAddBattleLog(BATTLE_ACTION_HEAL, m, m, healAmount)
	}
}

func (g *Game) ResolveAttackForMonster(attacker *MonsterCard) {
	if !attacker.HasAttack() {
		return
	}

	currentDeadMonstersCount := len(g.deadMonsters)

	// Magic attack
	if attacker.Magic > 0 {
		target := g.GetTargetForAttackType(attacker, ATTACK_TYPE_MAGIC)
		if target != nil {
			g.AttackMonsterPhase(attacker, target, ATTACK_TYPE_MAGIC)
		}
	}

	// Range attack
	if attacker.Ranged > 0 {
		target := g.GetTargetForAttackType(attacker, ATTACK_TYPE_RANGED)
		if target != nil {
			g.AttackMonsterPhase(attacker, target, ATTACK_TYPE_RANGED)
		}
	}

	// Melee attack
	if attacker.Melee > 0 {
		target := g.GetTargetForAttackType(attacker, ATTACK_TYPE_MELEE)
		if target != nil {
			g.ResolveMeleeAttackForMonster(attacker, target, ATTACK_TYPE_MELEE, false)
		}
	}

	// Check Bloodlust
	deadMonstersCount := len(g.deadMonsters) - currentDeadMonstersCount
	if attacker.HasAbility(ABILITY_BLOODLUST) && deadMonstersCount > 0 {
		isReverseSpeed := utils.Contains(g.rulesets, RULESET_REVERSE_SPEED)
		if deadMonstersCount > 1 {
			// two monsters might die from one attack with blast
			for i := 0; i < deadMonstersCount; i++ {
				g.MaybeApplyBloodlust(attacker, isReverseSpeed)
			}

		} else {
			g.MaybeApplyBloodlust(attacker, isReverseSpeed)
		}
	}
}

// Who this monster will target, if any. Null if none
func (g *Game) GetTargetForAttackType(m *MonsterCard, attackType CardAttackType) *MonsterCard {
	if !m.IsAlive() {
		return nil
	}
	enemyMonsters := g.GetEnemyTeamOfMonster(m)
	if len(enemyMonsters.GetAliveMonsters()) == 0 {
		return nil
	}

	if attackType == ATTACK_TYPE_MAGIC {
		return g.GetTargetForMagicAttack(m)
	} else if attackType == ATTACK_TYPE_RANGED {
		return g.GetTargetForRangedAttack(m)
	} else if attackType == ATTACK_TYPE_MELEE {
		return g.GetTargetForMeleeAttack(m)
	}
	return nil
}

func (g *Game) GetTargetForMagicAttack(m *MonsterCard) *MonsterCard {
	// check if the monster is in the first position
	friendlyTeam := g.GetTeamOfMonster(m)
	mPosition := friendlyTeam.GetMonsterPosition(m)
	if mPosition == 0 {
		enemyTeam := g.GetEnemyTeamOfMonster(m)
		return enemyTeam.GetFirstAliveMonster()
	}

	return g.GetTargetForNonMelee(m)
}

func (g *Game) GetTargetForNonMelee(m *MonsterCard) *MonsterCard {
	if m == nil {
		return nil
	}

	enemyTeam := g.GetEnemyTeamOfMonster(m)
	// Scattershot target
	if m.HasAbility(ABILITY_SCATTERSHOT) {
		return enemyTeam.GetScattershotTarget()
	}

	// Taunt
	tauntMonster := enemyTeam.GetTauntMonster()
	if tauntMonster != nil {
		return tauntMonster
	}

	// Sneak target
	if m.HasAbility(ABILITY_SNEAK) {
		return enemyTeam.GetSneakTarget()
	}

	// Snipe target
	if m.HasAbility(ABILITY_SNIPE) {
		return enemyTeam.GetSnipeTarget()
	}

	// Opportunity
	if m.HasAbility(ABILITY_OPPORTUNITY) {
		return enemyTeam.GetOpportunityTarget()
	}

	return enemyTeam.GetFirstAliveMonster()
}

func (g *Game) GetTargetForRangedAttack(m *MonsterCard) *MonsterCard {
	if m == nil {
		return nil
	}
	hasCloseRange := m.HasAbility(ABILITY_CLOSE_RANGE)
	friendlyTeam := g.GetTeamOfMonster(m)
	mPosition := friendlyTeam.GetMonsterPosition(m)

	// close range first position
	if hasCloseRange && mPosition == 0 {
		enemyTeam := g.GetEnemyTeamOfMonster(m)
		return enemyTeam.GetFirstAliveMonster()
	}
	// can't attack in the first position
	if mPosition == 0 {
		return nil
	}
	return g.GetTargetForNonMelee(m)
}

func (g *Game) GetTargetForMeleeAttack(m *MonsterCard) *MonsterCard {
	if m == nil {
		return nil
	}

	friendlyTeam := g.GetTeamOfMonster(m)
	enemyTeam := g.GetEnemyTeamOfMonster(m)
	mPosition := friendlyTeam.GetMonsterPosition(m)

	if mPosition == 0 {
		return enemyTeam.GetFirstAliveMonster()
	}

	// Sneak target
	if m.HasAbility(ABILITY_SNEAK) {
		return enemyTeam.GetSneakTarget()
	}

	// Opportunity
	if m.HasAbility(ABILITY_OPPORTUNITY) {
		return enemyTeam.GetOpportunityTarget()
	}

	// Melee mayhem
	if m.HasAbility(ABILITY_MELEE_MAYHEM) {
		return enemyTeam.GetFirstAliveMonster()
	}

	// Reach
	if mPosition == 1 && m.HasAbility(ABILITY_REACH) {
		return enemyTeam.GetFirstAliveMonster()
	}
	return nil
}

func (g *Game) AttackMonsterPhase(attacker, target *MonsterCard, attackType CardAttackType) {
	if attacker == nil || target == nil {
		return
	}

	// Recharge attack in odd index turns (starting 0)
	if attacker.HasAbility(ABILITY_RECHARGE) && g.roundNumber%2 == 0 {
		return
	}
	isAimTrue := utils.Contains(g.rulesets, RULESET_AIM_TRUE)
	wasAttackDoged := GetDidDodge(g.rulesets, attacker, target, attackType)
	if !isAimTrue && wasAttackDoged {
		g.CreateAndAddBattleLog(BATTLE_ACTION_ATTACK_DODGED, attacker, target, 0)
		g.MaybeApplyBackFire(attacker, target, attackType)
		// no more calculation since attack was dodged
		return
	}

	// Prepare variables
	damageMultiplier, _ := g.GetDamageMultiplier(attacker, target)
	baseDamage := attacker.GetPostAbilityAttackOfType(attackType)
	damageAmount := baseDamage * damageMultiplier
	attackedTeam := g.GetTeamOfMonster(target)
	attackedPosition := attackedTeam.GetMonsterPosition(target)
	attackedTeamAliveMonsters := attackedTeam.GetAliveMonsters()

	var prevMonster *MonsterCard
	if attackedPosition-1 >= 0 {
		prevMonster = attackedTeamAliveMonsters[attackedPosition-1]
	}
	var nextMonster *MonsterCard
	if attackedPosition+1 < len(attackedTeamAliveMonsters) {
		nextMonster = attackedTeamAliveMonsters[attackedPosition+1]
	}

	// Snare the target monster
	if target.HasAbility(ABILITY_FLYING) && attacker.HasAbility(ABILITY_SNARE) && !target.HasDebuff(ABILITY_SNARE) {
		g.AddMonsterDebuffToAMonster(attacker, target, ABILITY_SNARE, BATTLE_ACTION_SNARE)
	}

	// Divine shield
	if target.HasAbility(ABILITY_DIVINE_SHIELD) {
		g.HandleDivineShield(attacker, target, attackType, baseDamage, prevMonster, nextMonster, damageAmount)
		return
	}

	battleDamage := g.ActuallyHitMonster(attacker, target, attackType)
	g.CreateAndAddBattleLog(BATTLE_ACTION_ATTACK, attacker, target, battleDamage.DamageDone)

	// Pierce
	if attacker.HasAbility(ABILITY_PIERCING) && battleDamage.Remainder > 0 {
		// remainder already halved by shield or void. it just needs to hit health
		remainderDamage := HitHealth(target, battleDamage.Remainder)
		g.CreateAndAddBattleLog(BATTLE_ACTION_PIERCING_REMAINDER, attacker, target, remainderDamage)
	}

	// TODO: this doesn't account for the pierce
	g.MaybeApplyLifeLeech(attacker, battleDamage.ActualDamageDone)
	g.MaybeApplyThorns(attacker, target, attackType)
	g.MaybeApplyMagicReflect(attacker, target, attackType, attacker.GetPostAbilityAttackOfType(attackType))
	g.MaybeApplyReturnFire(attacker, target, attackType, attacker.GetPostAbilityAttackOfType(attackType))
	g.MaybeApplyRetaliate(attacker, target, attackType)
	g.MaybeApplyHalving(attacker, target)

	// check if dead
	g.ProcessIfDead(attacker)
	g.ProcessIfDead(target)

	// check blast
	g.MaybeApplyBlast(attacker, prevMonster, attackType, damageAmount)
	g.MaybeApplyBlast(attacker, nextMonster, attackType, damageAmount)

	// Shatter
	if attacker.HasAbility(ABILITY_SHATTER) {
		target.Armor = 0
	}

	// Stun
	g.MaybeApplyStun(attacker, target)

	// Poison
	g.MaybeApplyPoison(attacker, target)

	// Cripple
	g.MaybeApplyCripple(attacker, target)

	// Affliction
	if attacker.HasAbility(ABILITY_AFFLICTION) && !target.HasDebuff(ABILITY_AFFLICTION) && GetSuccessBelow(AFFLICTION_CHANCE*100) {
		g.CreateAndAddBattleLog(BATTLE_ACTION_AFFLICTION, attacker, target, battleDamage.DamageDone)
		g.AddMonsterDebuffToAMonster(attacker, target, ABILITY_AFFLICTION, BATTLE_ACTION_AFFLICTION)
	}

	// Dispel
	if attacker.HasAbility(ABILITY_DISPEL) {
		DispelBuffs(target)
	}
}

// Handle dodged log and backfire
// return true if backfire applied
func (g *Game) MaybeApplyBackFire(attacker, target *MonsterCard, attackType CardAttackType) bool {
	if !target.HasAbility(ABILITY_BACKFIRE) {
		return false
	}

	backfireBattleDamage := HitMonsterWithPhysical(g, target, BACKFIRE_DAMAGE)
	g.CreateAndAddBattleLog(BATTLE_ACTION_BACKFIRE, attacker, target, backfireBattleDamage.ActualDamageDone)
	// attacker gets damage from backfire and might die from it
	g.ProcessIfDead(attacker)
	return true
}

// Return damage multiplier and applied abilities
func (g *Game) GetDamageMultiplier(attacker, target *MonsterCard) (int, []Ability) {
	multiplier := 1
	appliedAbilities := []Ability{}
	// Recharge x3
	if attacker.HasAbility(ABILITY_RECHARGE) {
		multiplier *= 3
		appliedAbilities = append(appliedAbilities, ABILITY_RECHARGE)
	}

	// Giant killer x2
	if attacker.HasAbility(ABILITY_GIANT_KILLER) && target.Mana >= 10 {
		multiplier *= 2
		appliedAbilities = append(appliedAbilities, ABILITY_GIANT_KILLER)
	}

	// Deathblow x2
	if attacker.HasAbility(ABILITY_DEATHBLOW) && target.IsLastMonster() {
		multiplier *= 2
		appliedAbilities = append(appliedAbilities, ABILITY_DEATHBLOW)
	}

	// Knock out x2 if stunned
	if attacker.HasAbility(ABILITY_KNOCK_OUT) && target.HasDebuff(ABILITY_STUN) {
		multiplier *= 2
		appliedAbilities = append(appliedAbilities, ABILITY_KNOCK_OUT)
	}

	// Opress x2
	if attacker.HasAbility(ABILITY_OPPRESS) && !target.HasAttack() {
		multiplier *= 2
		appliedAbilities = append(appliedAbilities, ABILITY_OPPRESS)
	}
	return multiplier, appliedAbilities
}

func (g *Game) AddMonsterDebuffToAMonster(caster *MonsterCard, target *MonsterCard, debuff Ability, battleAction AdditionalBattleAction) {
	target.AddDebuff(debuff)
	g.CreateAndAddBattleLog(battleAction, caster, target, 0)
}

func (g *Game) HandleDivineShield(
	attacker, target *MonsterCard,
	attackType CardAttackType,
	AttackDamageForReflections int,
	prevMonster *MonsterCard,
	nextMonster *MonsterCard,
	damageAmount int,
) {
	target.RemoveDivineShield()
	g.CreateAndAddBattleLog(BATTLE_ACTION_REMOVE_DIVINE_SHIELD, attacker, target, 0)

	// Handle Reflective Damage
	if attackType == ATTACK_TYPE_MAGIC {
		// magic reflect
		g.MaybeApplyMagicReflect(attacker, target, attackType, AttackDamageForReflections)
	} else if attackType == ATTACK_TYPE_RANGED {
		g.MaybeApplyReturnFire(attacker, target, attackType, AttackDamageForReflections)
	} else {
		g.MaybeApplyThorns(attacker, target, attackType)
		g.MaybeApplyRetaliate(attacker, target, attackType)
	}

	// Shatter
	if attacker.HasAbility(ABILITY_SHATTER) {
		target.Armor = 0
	}

	// check if dead
	g.ProcessIfDead(attacker)
	g.ProcessIfDead(target)

	// Do buffs
	g.MaybeApplyStun(attacker, target)
	g.MaybeApplyPoison(attacker, target)
	g.MaybeApplyCripple(attacker, target)
	g.MaybeApplyHalving(attacker, target)

	// Debuffs
	// Affliction
	if attacker.HasAbility(ABILITY_AFFLICTION) && !target.HasDebuff(ABILITY_AFFLICTION) && GetSuccessBelow(AFFLICTION_CHANCE*100) {
		g.CreateAndAddBattleLog(BATTLE_ACTION_AFFLICTION, attacker, target, 0)
		g.AddMonsterDebuffToAMonster(attacker, target, ABILITY_AFFLICTION, BATTLE_ACTION_AFFLICTION)
	}

	// Dispel
	if attacker.HasAbility(ABILITY_DISPEL) {
		DispelBuffs(target)
	}

	g.MaybeApplyBlast(attacker, prevMonster, attackType, damageAmount)
	g.MaybeApplyBlast(attacker, nextMonster, attackType, damageAmount)
}

// Resolve attack involves Trample
func (g *Game) ResolveMeleeAttackForMonster(attacker, target *MonsterCard, attackType CardAttackType, hasTrampled bool) {
	beforeAttackDeadMonsterCount := len(g.deadMonsters)
	enemyTeam := g.GetTeamOfMonster(target)
	mPosition := enemyTeam.GetMonsterPosition(target)
	aliveEnemyMonsters := enemyTeam.GetAliveMonsters()
	var nextMonster *MonsterCard
	if mPosition+1 < len(aliveEnemyMonsters) {
		nextMonster = aliveEnemyMonsters[mPosition+1]
	}
	g.AttackMonsterPhase(attacker, target, ATTACK_TYPE_MELEE)

	// Check Trample
	if !hasTrampled && nextMonster != nil && attacker.HasAbility(ABILITY_TRAMPLE) && len(g.deadMonsters) > beforeAttackDeadMonsterCount {
		g.ResolveMeleeAttackForMonster(attacker, nextMonster, attackType, utils.Contains(g.rulesets, RULESET_STAMPEDE))
	}
}

func (g *Game) MaybeApplyBloodlust(attacker *MonsterCard, isReverseSpeed bool) {
	if !attacker.HasAbility(ABILITY_BLOODLUST) || attacker.Health <= 0 {
		return
	}

	// add attack if have attack
	if attacker.Magic > 0 {
		attacker.Magic += 1
	}
	if attacker.Ranged > 0 {
		attacker.Ranged += 1
	}
	if attacker.Melee > 0 {
		attacker.Melee += 1
	}
	if attacker.GetPostAbilityMaxArmor() > 0 {
		attacker.Armor += 1
	}

	if isReverseSpeed && attacker.Speed > 0 {
		attacker.Speed -= 1
	} else {
		attacker.Speed += 1
	}
	attacker.Health += 1
	g.CreateAndAddBattleLog(BATTLE_ACTION_BLOODLUST, attacker, nil, 0)
}

func (g *Game) MaybeApplyMagicReflect(attacker *MonsterCard, target *MonsterCard, attackType CardAttackType, attackDamageForReflections int) {
	if !target.HasAbility(ABILITY_MAGIC_REFLECT) || attackType != ATTACK_TYPE_MAGIC {
		return
	}

	attackDamage := 0
	if attackDamageForReflections != 0 {
		attackDamage = attackDamageForReflections
	}

	// half of the magic attack damage (round up)
	reflectDamage := int(math.Ceil(float64(attackDamage) / 2))

	// Amplify
	if attacker.HasDebuff(ABILITY_AMPLIFY) {
		reflectDamage += 1
	}

	// Reflection shield
	if attacker.HasAbility(ABILITY_REFLECTION_SHIELD) {
		reflectDamage = 0
	}

	battleDamage := HitMonsterWithMagic(g, attacker, reflectDamage)
	g.CreateAndAddBattleLog(BATTLE_ACTION_MAGIC_REFLECT, attacker, target, battleDamage.DamageDone)
}

func (g *Game) MaybeApplyReturnFire(attacker *MonsterCard, target *MonsterCard, attackType CardAttackType, attackDamageForReflections int) {
	if !target.HasAbility(ABILITY_RETURN_FIRE) || attackType != ATTACK_TYPE_RANGED {
		return
	}

	attackDamage := 0
	if attackDamageForReflections != 0 {
		attackDamage = attackDamageForReflections
	}

	// half of the ranged attack damage (round up)
	reflectDamage := int(math.Ceil(float64(attackDamage) / 2))

	// Amplify
	if attacker.HasDebuff(ABILITY_AMPLIFY) {
		reflectDamage += 1
	}

	// Reflection shield
	if attacker.HasAbility(ABILITY_REFLECTION_SHIELD) {
		reflectDamage = 0
	}

	battleDamage := HitMonsterWithPhysical(g, attacker, reflectDamage)
	g.CreateAndAddBattleLog(BATTLE_ACTION_RETURN_FIRE, attacker, target, battleDamage.DamageDone)
}

func (g *Game) MaybeApplyThorns(attacker *MonsterCard, target *MonsterCard, attackType CardAttackType) {
	if !target.HasAbility(ABILITY_THORNS) || attackType != ATTACK_TYPE_MELEE {
		return
	}

	// Thorns always do the same damage amount
	reflectDamage := THORNS_DAMAGE

	// Amplify
	if attacker.HasDebuff(ABILITY_AMPLIFY) {
		reflectDamage += 1
	}

	// Reflection shield
	if attacker.HasAbility(ABILITY_REFLECTION_SHIELD) {
		reflectDamage = 0
	}

	battleDamage := HitMonsterWithPhysical(g, attacker, reflectDamage)
	g.CreateAndAddBattleLog(BATTLE_ACTION_THORNS, attacker, target, battleDamage.DamageDone)
}

func (g *Game) MaybeApplyRetaliate(attacker *MonsterCard, target *MonsterCard, attackType CardAttackType) {
	if attacker == nil || target == nil || !target.HasAbility(ABILITY_RETALIATE) || attackType != ATTACK_TYPE_MELEE {
		return
	}

	// Retaliate chance
	doesRetaliate := GetSuccessBelow(RETALIATE_CHANCE * 100)
	if !doesRetaliate {
		return
	}

	g.CreateAndAddBattleLog(BATTLE_ACTION_RETALIATE, target, attacker, attacker.GetPostAbilityMelee())
	g.AttackMonsterPhase(target, attacker, ATTACK_TYPE_MELEE)
}

func (g *Game) MaybeApplyStun(attacker, target *MonsterCard) {
	if attacker.HasAbility(ABILITY_STUN) && GetSuccessBelow(STUN_CHANCE*100) {
		stunDataKey := g.GetStunDataKey(g.roundNumber, attacker)
		prevStunnedMonsters := g.stunData[stunDataKey]
		if prevStunnedMonsters == nil {
			prevStunnedMonsters = []*MonsterCard{}
		}
		prevStunnedMonsters = append(prevStunnedMonsters, target)
		g.stunData[stunDataKey] = prevStunnedMonsters
		g.AddMonsterDebuffToAMonster(attacker, target, ABILITY_STUN, BATTLE_ACTION_STUN)
	}
}

func (g *Game) MaybeApplyPoison(attacker, target *MonsterCard) {
	if attacker.HasAbility(ABILITY_POISON) && GetSuccessBelow(POISON_CHANCE*100) && !target.HasDebuff(ABILITY_POISON) && target.IsAlive() {
		g.AddMonsterDebuffToAMonster(attacker, target, ABILITY_POISON, BATTLE_ACTION_POISON)
	}
}

func (g *Game) MaybeApplyCripple(attacker, target *MonsterCard) {
	if attacker.HasAbility(ABILITY_CRIPPLE) && target.IsAlive() {
		g.AddMonsterDebuffToAMonster(attacker, target, ABILITY_CRIPPLE, BATTLE_ACTION_CRIPPLE)
	}
}

func (g *Game) MaybeApplyHalving(attacker, target *MonsterCard) {
	if attacker.HasAbility(ABILITY_HALVING) && target.IsAlive() {
		g.AddMonsterDebuffToAMonster(attacker, target, ABILITY_HALVING, BATTLE_ACTION_HALVING)
	}
}

func (g *Game) MaybeApplyLifeLeech(attacker *MonsterCard, damage int) {
	if !attacker.IsAlive() || !attacker.HasAbility(ABILITY_LIFE_LEECH) {
		return
	}

	lifeLeechAmount := int(math.Ceil(float64(damage) / 2))
	if lifeLeechAmount > 0 {
		for i := 0; i < lifeLeechAmount; i++ {
			attacker.AddBuff(ABILITY_LIFE_LEECH)
		}
		g.CreateAndAddBattleLog(BATTLE_ACTION_LIFE_LEECH, attacker, nil, lifeLeechAmount)
	}
}

// Blast has a ton of edge cases... https://support.splinterlands.com/hc/en-us/articles/4414966685332-Abilities-Status-Effects
func (g *Game) MaybeApplyBlast(attacker, blastTarget *MonsterCard, attackType CardAttackType, damageAmount int) {
	if !attacker.HasAbility(ABILITY_BLAST) || blastTarget == nil {
		return
	}

	baseBlastDamage := int(math.Ceil(float64(damageAmount) / 2))
	damageMultiplier := g.GetPostBlastDamageMultiplier(attacker, blastTarget)
	blastDamage := baseBlastDamage * damageMultiplier

	// Forcefield
	if blastTarget.HasAbility(ABILITY_FORCEFIELD) && blastDamage >= FORCEFIELD_MIN_DAMAGE {
		blastDamage = 1
	}

	// Snare
	if blastTarget.HasAbility(ABILITY_FLYING) && attacker.HasAbility(ABILITY_SNARE) && !blastTarget.HasDebuff(ABILITY_SNARE) {
		g.AddMonsterDebuffToAMonster(attacker, blastTarget, ABILITY_SNARE, BATTLE_ACTION_SNARE)
	}

	// Reflection shield
	if blastTarget.HasAbility(ABILITY_REFLECTION_SHIELD) {
		blastDamage = 0
	}

	// Magic blast damage
	if attackType == ATTACK_TYPE_MAGIC {
		battleDamage := HitMonsterWithMagic(g, blastTarget, blastDamage)
		g.CreateAndAddBattleLog(BATTLE_ACTION_BLAST, attacker, blastTarget, battleDamage.DamageDone)
		g.MaybeApplyMagicReflect(attacker, blastTarget, attackType, battleDamage.Attack)
		g.MaybeApplyLifeLeech(attacker, (blastDamage - battleDamage.Remainder))
	} else {
		// melee or range attack
		battleDamage := HitMonsterWithPhysical(g, blastTarget, blastDamage)
		g.CreateAndAddBattleLog(BATTLE_ACTION_BLAST, attacker, blastTarget, battleDamage.DamageDone)
		g.MaybeApplyReturnFire(attacker, blastTarget, attackType, battleDamage.Attack)
	}

	// check dead monster
	g.ProcessIfDead(blastTarget)
	g.ProcessIfDead(attacker)
}

func (g *Game) GetPostBlastDamageMultiplier(attacker, blastTarget *MonsterCard) int {
	multiplier := 1

	// Recharge
	if attacker.HasAbility(ABILITY_RECHARGE) {
		multiplier *= 3
	}

	// Giant killer
	if attacker.HasAbility(ABILITY_GIANT_KILLER) && blastTarget.Mana >= 10 {
		multiplier *= 2
	}

	// Knock out
	if blastTarget.HasDebuff(ABILITY_STUN) && attacker.HasAbility(ABILITY_KNOCK_OUT) {
		multiplier *= 2
	}

	// Oppress
	if !blastTarget.HasAttack() && attacker.HasAbility(ABILITY_OPPRESS) {
		multiplier *= 2
	}

	return multiplier
}

func (g *Game) DoPostRound() {
	aliveTeam1 := g.team1.GetAliveMonsters()
	aliveTeam2 := g.team2.GetAliveMonsters()

	// Earthquake
	if utils.Contains(g.rulesets, RULESET_EARTHQUAKE) {
		g.DoPostRoundEarthquake(aliveTeam1)
		g.DoPostRoundEarthquake(aliveTeam2)
		aliveTeam1 = g.team1.GetAliveMonsters()
		aliveTeam2 = g.team2.GetAliveMonsters()
	}

	// Poison
	g.DoPostRoundPoison(aliveTeam1)
	g.DoPostRoundPoison(aliveTeam2)
}

// handle earthquake
func (g *Game) DoPostRoundEarthquake(monsters []*MonsterCard) {
	for _, m := range monsters {
		battleDamage := ApplyEarthquake(g, m)
		g.CreateAndAddBattleLog(BATTLE_ACTION_EARTHQUAKE, m, nil, battleDamage.DamageDone)
		g.ProcessIfDead(m)
		g.CheckAndSetGameWinner()
		if g.winner != TEAM_NUM_UNKNOWN {
			g.LogGameOver()
			return
		}
	}
}

// Handle poison (another name monstersOnPostRound)
func (g *Game) DoPostRoundPoison(monsters []*MonsterCard) {
	for _, m := range monsters {
		if m.HasDebuff(ABILITY_POISON) {
			m.Health -= POISON_DAMAGE
			g.ProcessIfDead(m)
			g.CreateAndAddBattleLog(BATTLE_ACTION_POISON, m, nil, POISON_DAMAGE)
		}
		m.SetHasTurnPassed(false)
	}
}

func (g *Game) GetTeamOfMonster(m *MonsterCard) *GameTeam {
	if m.GetTeamNumber() == 1 {
		return g.team1
	}
	return g.team2
}

func (g *Game) GetEnemyTeamOfMonster(m *MonsterCard) *GameTeam {
	if m.GetTeamNumber() == 1 {
		return g.team2
	}
	return g.team1
}

// Returns true if resurrected, false otherwise
func (g *Game) ProcessIfResurrect(caster GameCardInterface, deadMonster *MonsterCard) bool {
	if caster.HasAbility(ABILITY_RESURRECT) && !deadMonster.IsAlive() {
		caster.RemoveAbility(ABILITY_RESURRECT)
		deadMonster.Resurrect()

		// remove it from the dead monsters list
		deadMonsterList := make([]*MonsterCard, 0)
		for _, m := range g.deadMonsters {
			if m.cardDetail.Name == deadMonster.cardDetail.Name {
				continue
			}
			deadMonsterList = append(deadMonsterList, m)
		}
		g.deadMonsters = deadMonsterList
		g.CreateAndAddBattleLog(BATTLE_ACTION_RESURRECT, caster, deadMonster, 0)
		return true
	}

	return false
}

func (g *Game) RemoveMonsterDebuff(m *MonsterCard, debuffs []Ability, enemyMonsters []*MonsterCard) {
	if len(debuffs) == 0 {
		return
	}

	for _, debuff := range debuffs {
		for _, enemy := range enemyMonsters {
			enemy.RemoveDebuff(debuff)
		}
	}
}

func (g *Game) RemoveMonsterBuff(m *MonsterCard, buffs []Ability, friendlyMonsters []*MonsterCard) {
	if len(buffs) == 0 {
		return
	}

	for _, buff := range buffs {
		for _, friendlyMonster := range friendlyMonsters {
			friendlyMonster.RemoveBuff(buff)
		}
	}
}

// Handle scavenger and battle log
func (g *Game) OnMonsterDeath(m *MonsterCard, deadMonster *MonsterCard) {
	// Scavenger
	if m.HasAbility(ABILITY_SCAVENGER) {
		m.AddBuff(ABILITY_SCAVENGER)
		g.CreateAndAddBattleLog(BATTLE_ACTION_SCAVENGER, m, deadMonster, 1)
	}
}

func (g *Game) ActuallyHitMonster(attacker, target *MonsterCard, attackType CardAttackType) BattleDamage {
	damageMultiplier, _ := g.GetDamageMultiplier(attacker, target)
	damageAmount := attacker.GetPostAbilityAttackOfType(attackType) * damageMultiplier
	battleDamage := BattleDamage{Attack: 0, DamageDone: 0, Remainder: 0, ActualDamageDone: 0}

	if attackType == ATTACK_TYPE_MAGIC {
		battleDamage = HitMonsterWithMagic(g, target, damageAmount)
	} else if attackType == ATTACK_TYPE_RANGED {
		battleDamage = HitMonsterWithPhysical(g, target, damageAmount)
	} else if attackType == ATTACK_TYPE_MELEE {
		battleDamage = HitMonsterWithPhysical(g, target, damageAmount)
	}

	return battleDamage
}

package simulator

import (
	"errors"
	. "game_models"
	utils "game_utils"
	"sort"
)

type Game struct {
	team1      GameTeam
	team2      GameTeam
	rulesets   []Ruleset
	battleLogs []BattleLog
	shouldLog  bool
	/* 1: team1, 2: team2, 3: Tie */
	winner       *TeamNumber
	deadMonsters []MonsterCard
	roundNumber  int
}

func (g *Game) Create(team1 GameTeam, team2 GameTeam, rulesets []Ruleset, shouldLog bool) {
	g.team1 = team1
	g.team2 = team2
	g.rulesets = rulesets
	g.shouldLog = shouldLog
	g.team1.SetTeamNumber(TEAM_NUM_ONE)
	g.team2.SetTeamNumber(TEAM_NUM_TWO)
}

func (g *Game) GetWinner() *TeamNumber {
	return g.winner
}

func (g *Game) GetBattleLogs() ([]BattleLog, error) {
	if !g.shouldLog {
		return []BattleLog{}, errors.New("you must instantiate the game with enableLogs as true")
	}
	return g.battleLogs, nil
}

func (g *Game) PlayGame() {
	g.roundNumber = 0
	team1Summoner := g.team1.GetSummoner()
	team1Monsters := g.team1.GetMonstersList()
	team2Summoner := g.team2.GetSummoner()
	team2Monsters := g.team2.GetMonstersList()

	// pre game rulesets
	utils.DoRulesetPreGameBuff(g.rulesets, g.team1, g.team2)

	// Summoner pre-game buffs
	g.DoSummonerPreGameBuff(team1Summoner, team1Monsters)
	g.DoSummonerPreGameBuff(team2Summoner, team2Monsters)

	// Monsters pre-game buffs
	g.DoMonsterPreGameBuff(team1Monsters)
	g.DoMonsterPreGameBuff(team2Monsters)

	// Summoner pre-game debuffs
	g.DoSummonerPreGameDebuff(team1Summoner, team2Monsters)
	g.DoSummonerPreGameDebuff(team2Summoner, team1Monsters)

	// Monsters pre-game debuffs
	g.DoMonsterPreGameDebuff(team1Monsters, team2Monsters)

	// Apply ruleset rules that apply post buff phase
	utils.DoRulesetPreGamePostBuff(g.rulesets, g.team1, g.team2)

	g.team1.SetAllMonsterHealth()
	g.team2.SetAllMonsterHealth()
	g.PlayRoundsUntilGameEnd(0)
}

func (g *Game) DoSummonerPreGameBuff(summoner SummonerCard, friendlyMonsters []MonsterCard) {
	// add summoner abilities that increase stats (aka, strengthen)
	for _, ability := range utils.GetSummonerPreGameBuffAbilities() {
		if summoner.HasAbility(ability) {
			g.ApplyBuffToMonsters(friendlyMonsters, ability)
		}
	}

	// add summoner abilities
	for _, ability := range utils.GetSummonerAbilityAbilities() {
		if summoner.HasAbility(ability) {
			g.ApplyAbilityToMonsters(friendlyMonsters, ability)
		}
	}

	// add summoner stats (e.g. +1 melee, +1 archery, +1 magic etc...)
	for _, m := range friendlyMonsters {
		if summoner.Armor > 0 {
			m.AddSummonerArmor(summoner.Armor)
		}
		if summoner.Health > 0 {
			m.AddSummonerHealth(summoner.Health)
		}
		if summoner.Speed > 0 {
			m.AddSummonerSpeed(summoner.Speed)
		}
		if summoner.Melee > 0 {
			m.AddSummonerMelee(summoner.Melee)
		}
		if summoner.Ranged > 0 {
			m.AddSummonerRanged(summoner.Ranged)
		}
		if summoner.Magic > 0 {
			m.AddSummonerMagic(summoner.Magic)
		}
	}
}

// Add all summoner abilities to all enemy monsters which are in SUMMONER_DEBUFF_ABILITIES
func (g *Game) DoSummonerPreGameDebuff(summoner SummonerCard, targetMonsters []MonsterCard) {
	// add summoner debuffs (aka, affliciton, blind)
	for _, ability := range utils.GetSummonerPreGameDebuffAbilities() {
		if summoner.HasAbility(ability) {
			g.ApplyDebuffToMonsters(targetMonsters, ability)
		}
	}

	// add summoner abilities
	for _, ability := range utils.GetSummonerAbilityAbilities() {
		if summoner.HasAbility(ability) {
			g.ApplyAbilityToMonsters(targetMonsters, ability)
		}
	}

	// add summoner stats (e.g. +1 melee, +1 archery, +1 magic etc...)
	for _, m := range targetMonsters {
		if summoner.Armor > 0 {
			m.AddSummonerArmor(summoner.Armor)
		}
		if summoner.Health > 0 {
			m.AddSummonerHealth(summoner.Health)
		}
		if summoner.Speed > 0 {
			m.AddSummonerSpeed(summoner.Speed)
		}
		if summoner.Melee > 0 {
			m.AddSummonerMelee(summoner.Melee)
		}
		if summoner.Ranged > 0 {
			m.AddSummonerRanged(summoner.Ranged)
		}
		if summoner.Magic > 0 {
			m.AddSummonerMagic(summoner.Magic)
		}
	}
}

// Add all monster abilities to all friendly monsters which are in MONSTER_BUFF_ABILITIES
func (g *Game) DoMonsterPreGameBuff(friendlyMonsters []MonsterCard) {
	for _, m := range friendlyMonsters {
		for _, buff := range utils.GetMonsterPreGameBuffAbilities() {
			if !m.HasAbility(buff) {
				continue
			}

			g.ApplyBuffToMonsters(friendlyMonsters, buff)
		}
	}
}

func (g *Game) DoMonsterPreGameDebuff(team1Monsters []MonsterCard, team2Monsters []MonsterCard) {
	for _, debuff := range utils.GetMonsterPreGameDebuffAbilities() {
		// team1 debuff team2
		for _, m := range team1Monsters {
			if !m.HasAbility(debuff) {
				continue
			}
			g.ApplyDebuffToMonsters(team2Monsters, debuff)
		}

		// team2 debuff team1
		for _, m := range team2Monsters {
			if !m.HasAbility(debuff) {
				continue
			}
			g.ApplyDebuffToMonsters(team1Monsters, debuff)
		}
	}
}

func (g *Game) ApplyBuffToMonsters(monsters []MonsterCard, buff Ability) {
	// add stats buff
	for _, m := range monsters {
		m.AddBuff(buff)
	}
}

func (g *Game) ApplyAbilityToMonsters(monsters []MonsterCard, ability Ability) {
	for _, m := range monsters {
		m.AddAbilitiesWithArray([]Ability{ability})
	}
}

func (g *Game) ApplyDebuffToMonsters(monsters []MonsterCard, debuff Ability) {
	for _, m := range monsters {
		m.AddDebuff(debuff)
	}
}

func (g *Game) GetAllAliveMonsters() []MonsterCard {
	aliveMonsters := make([]MonsterCard, 0)
	t1 := g.team1.GetAliveMonsters()
	t2 := g.team2.GetAliveMonsters()
	aliveMonsters = append(t1, t2...)
	return aliveMonsters
}

// Plays the game rounds until the game is over
func (g *Game) PlayRoundsUntilGameEnd(roundNumber int) {
	g.roundNumber = roundNumber
	if g.winner != nil {
		return
	}

	// if round >= 50, game is tie
	if roundNumber > 50 {
		*g.winner = TEAM_NUM_UNKNOWN
	}

	// Fatigue
	if roundNumber >= FATIGUE_ROUND_NUMBER {
		g.FatigueMonsters(roundNumber)
		g.CheckAndSetGameWinner()
		if g.winner != nil {
			return
		}
	}

	// Play one round
	g.PlaySingleRound()
	g.CheckAndSetGameWinner()
	if g.winner != nil {
		return
	}

	// Post round including earthquake
	g.DoPostRound()
	g.CheckAndSetGameWinner()

	g.PlayRoundsUntilGameEnd(roundNumber + 1)
}

func (g *Game) FatigueMonsters(roundNumber int) {
	fatigueDamage := roundNumber - FATIGUE_ROUND_NUMBER + 1
	allAliveMonsters := g.GetAllAliveMonsters()

	for _, m := range allAliveMonsters {
		g.CreateAndAddBattleLog(BATTLE_ACTION_FATIGUE, &m, nil, fatigueDamage)
		m.HitHealth(fatigueDamage)
		g.ProcessIfDead(m)
	}

	g.CheckAndSetGameWinner()
	if g.winner != nil {
		return
	}
}

func (g *Game) CreateAndAddBattleLog(action AdditionalBattleAction, cardOne GameCardInterface, cardTwo GameCardInterface, value int) {
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
		*g.winner = TEAM_NUM_UNKNOWN
	} else if team2AliveMonstersCount == 0 {
		*g.winner = TEAM_NUM_ONE
	} else if team1AliveMonstersCount == 0 {
		*g.winner = TEAM_NUM_TWO
	}
}

func (g *Game) ProcessIfDead(m MonsterCard) {
	if m.IsAlive() || !utils.CardsArrIncludesMonster(g.deadMonsters, m) {
		return
	}

	// monster is dead
	g.CreateAndAddBattleLog(BATTLE_ACTION_DEATH, &m, nil, 0)
	g.deadMonsters = append(g.deadMonsters, m)
	m.SetHasTurnPassed(true)

	friendlyTeam := g.GetTeamOfMonster(m)
	aliveFriendlyMonsters := friendlyTeam.GetAliveMonsters()
	enemyTeam := g.GetEnemyTeamOfMonster()
	aliveEnemyMonsters := enemyTeam.GetAliveMonsters()

	// Redemption
	if m.HasAbility(ABILITY_REDEMPTION) {
		for _, e := range aliveEnemyMonsters {
			utils.HitMonsterWithPhysical(
				g,
				e,
				utils.REDEMPTION_DAMAGE,
			)

			g.ProcessIfDead(e)
		}
	}

	// Ressurect
	friendlySummoner := friendlyTeam.GetSummoner()
	// summoner resurrect
	wasResurrected := g.ProcessIfResurrect(&friendlySummoner, m)
	// friendly monster resurrect
	for _, friendlyMonster := range aliveFriendlyMonsters {
		if wasResurrected {
			break
		}
		wasResurrected = g.ProcessIfResurrect(&friendlyMonster, m)
	}

	// remove debuffs and buffs if not resurrected
	if !wasResurrected {
		// remove debuffs from the enemy team
		monsterDebuffs := utils.MonsterHasDebuffAbilities(m)
		g.RemoveMonsterDebuff(m, monsterDebuffs, aliveEnemyMonsters)

		// remove buffs from the friendly monsters
		monsterBuffs := utils.MonsterHasBuffsAbilities(m)
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
	g.deadMonsters = []MonsterCard{}
	g.DoGamePreRound()
	g.DoSummonerPreRound(g.team1)
	g.DoSummonerPreRound(g.team2)

	// loop through each monster's turn
	stunnedMonsters := []MonsterCard{}
	currentMonster := g.GetNextMonsterTurn()
	for currentMonster != nil {
		if !currentMonster.IsAlive() {
			continue
		}

		// check stun
		if currentMonster.HasDebuff(ABILITY_STUN) {
			stunnedMonsters = append(stunnedMonsters, *currentMonster)
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

	// remove stun state
	for _, sm := range stunnedMonsters {
		sm.RemoveAllDebuff(ABILITY_STUN)
	}
}

func (g *Game) DoGamePreRound() {
	// Add pre round actions when necessary
}

// Handle Summoner's pre turn actions (e.g. cleanse, tank heal, repair, triage)
func (g *Game) DoSummonerPreRound(t GameTeam) {
	summoner := t.GetSummoner()

	// Cleanse
	if summoner.HasAbility(ABILITY_CLEANSE) {
		firstMonster := t.GetFirstAliveMonster()
		firstMonster.CleanseDebuffs()
		g.CreateAndAddBattleLog(BATTLE_ACTION_CLEANSE, &summoner, firstMonster, 0)
	}

	// Repair
	if summoner.HasAbility(ABILITY_REPAIR) {
		repairTarget := t.GetRepairTarget()
		if repairTarget != nil {
			utils.RepairMonsterArmor(repairTarget)
			g.CreateAndAddBattleLog(BATTLE_ACTION_REPAIR, &summoner, repairTarget, 0)
		}
	}

	// Tank heal
	if summoner.HasAbility(ABILITY_TANK_HEAL) {
		firstMonster := t.GetFirstAliveMonster()
		utils.TankHealMonster(firstMonster)
		g.CreateAndAddBattleLog(BATTLE_ACTION_TANK_HEAL, &summoner, firstMonster, 0)
	}

	// Triage
	if summoner.HasAbility(ABILITY_TRIAGE) {
		healTarget := t.GetTriageHealTarget()
		if healTarget != nil {
			healAmount := utils.TriageHealMonster(healTarget)
			g.CreateAndAddBattleLog(BATTLE_ACTION_TANK_HEAL, &summoner, healTarget, healAmount)
		}
	}
}

// TODO Test
func (g *Game) GetNextMonsterTurn() *MonsterCard {
	allUnmovedMonsters := append(g.team1.GetUnmovedMonsters(), g.team2.GetUnmovedMonsters()...)
	if len(allUnmovedMonsters) == 0 {
		return nil
	}

	// sort unmoved monsters
	sort.SliceStable(allUnmovedMonsters, utils.MonsterTurnComparator)

	if utils.Contains(g.rulesets, RULESET_REVERSE_SPEED) {
		return allUnmovedMonsters[0]
	}

	return allUnmovedMonsters[len(allUnmovedMonsters)-1]
}

// TODO
func (g *Game) DoMonsterPreTurn(m *MonsterCard) {
	m.SetHasTurnPasses(true)
	friendlyTeam := g.GetTeamOfMonster(m)

	// Cleanse
	if m.HasAbility(ABILITY_CLEANSE) {
		cleanseTarget := friendlyTeam.GetFirstAliveMonster()
		cleanseTarget.CleanseDebuffs()
		g.CreateAndAddBattleLog(ABILITY_CLEANSE, &m, cleanseTarget, 0)
	}

	// Tank heal
	if m.HasAbility(ABILITY_TANK_HEAL) {
		tankHealTarget := friendlyTeam.GetFirstAliveMonster()
		healAmount := utils.TankHealMonster(tankHealTarget)
		g.CreateAndAddBattleLog(BATTLE_ACTION_TANK_HEAL, m, tankHealTarget, healAmount)
	}

	// Repair
	if m.HasAbility(ABILITY_REPAIR) {
		repairTarget := friendlyTeam.GetRepairTarget()
		if repairTarget {
			repairAmount = utils.RepairMonsterArmor(repairTarget)
		}
	}
}

// TODO
func (g *Game) ResolveAttackForMonster(m *MonsterCard) {

}

// TODO
func (g *Game) DoPostRound() {

}

// TODO
func (g *Game) GetTeamOfMonster(m MonsterCard) GameTeam {

	return GameTeam{}
}

// TODO
func (g *Game) GetEnemyTeamOfMonster(m MonsterCard) GameTeam {
	return GameTeam{}
}

// TODO Returns true if resurrected, false otherwise
func (g *Game) ProcessIfResurrect(caster GameCardInterface, deadMonster MonsterCard) bool {
	return false
}

func (g *Game) RemoveMonsterDebuff(m MonsterCard, debuffs []Ability, enemyMonsters []MonsterCard) {
	if len(debuffs) == 0 {
		return
	}

	for _, debuff := range debuffs {
		for _, enemy := range enemyMonsters {
			enemy.RemoveDebuff(debuff)
		}
	}
}

func (g *Game) RemoveMonsterBuff(m MonsterCard, buffs []Ability, friendlyMonsters []MonsterCard) {
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
func (g *Game) OnMonsterDeath(m MonsterCard, deadMonster MonsterCard) {
	// Scavenger
	if m.HasAbility(ABILITY_SCAVENGER) {
		m.AddBuff(ABILITY_SCAVENGER)
		g.CreateAndAddBattleLog(BATTLE_ACTION_SCAVENGER, &m, &deadMonster, 1)
	}
}

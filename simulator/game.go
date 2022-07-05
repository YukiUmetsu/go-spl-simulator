package simulator

import (
	"errors"
	utils "game_utils"
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

func (g *Game) CreateAndAddBattleLog(action AdditionalBattleAction, cardOne *MonsterCard, cardTwo *MonsterCard, value int) {
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

// TODO
func (g *Game) PlaySingleRound() {

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

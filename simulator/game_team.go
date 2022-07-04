package simulator

import (
	"math/rand"
	"time"
)

type GameTeam struct {
	summoner    SummonerCard
	monsterList []MonsterCard
}

func (t *GameTeam) Create(summoner SummonerCard, monsterList []MonsterCard) {
	t.summoner = summoner
	t.monsterList = monsterList
	if len(t.monsterList) > 1 {
		t.monsterList[0].SetIsOnlyMonster()
	}
	t.SetMonsterPositions()
}

func (t *GameTeam) SetMonsterPositions() {
	for i := range t.monsterList {
		t.monsterList[i].SetCardPosition(i)
	}
}

func (t *GameTeam) SetTeamNumber(teamNumber TeamNumber) {
	t.summoner.SetTeam(teamNumber)
	for _, monster := range t.monsterList {
		monster.SetTeam(teamNumber)
	}
}

/** Position of the alive monsters */
func (t *GameTeam) GetMonsterPosition(monster MonsterCard) int {
	aliveMonsters := t.GetAliveMonsters()
	for i, m := range aliveMonsters {
		if m.cardDetail.ID == monster.cardDetail.ID {
			return i
		}
	}
	return -1
}

func (t *GameTeam) GetSummoner() SummonerCard {
	return t.summoner
}

func (t *GameTeam) GetMonstersList() []MonsterCard {
	return t.monsterList
}

func (t *GameTeam) GetFirstAliveMonster() *MonsterCard {
	aliveMonsters := t.GetAliveMonsters()
	if len(aliveMonsters) == 0 {
		return nil
	}
	return &aliveMonsters[0]
}

func (t *GameTeam) GetUnmovedMonsters() []MonsterCard {
	unmovedMonsters := make([]MonsterCard, 0)
	for _, m := range t.GetAliveMonsters() {
		if !m.GetHasTurnPassed() {
			unmovedMonsters = append(unmovedMonsters, m)
		}
	}
	return unmovedMonsters
}

func (t *GameTeam) GetAliveMonsters() []MonsterCard {
	aliveMonsters := make([]MonsterCard, 0)
	for _, monster := range t.monsterList {
		if monster.IsAlive() {
			aliveMonsters = append(aliveMonsters, monster)
		}
	}
	return aliveMonsters
}

func (t *GameTeam) MaybeSetLastStand() {
	aliveMonsters := t.GetAliveMonsters()
	if len(aliveMonsters) == 1 {
		aliveMonsters[0].SetIsOnlyMonster()
	}
}

func (t *GameTeam) SetAllMonsterHealth() {
	for _, m := range t.monsterList {
		m.Health = m.GetPostAbilityMaxHealth()
	}
}

func (t *GameTeam) GetScattershotTarget() MonsterCard {
	aliveMonsters := t.GetAliveMonsters()
	rand.Seed(time.Now().Unix())
	randomMonsterNum := rand.Intn(len(aliveMonsters))
	return aliveMonsters[randomMonsterNum]
}

func (t *GameTeam) GetSnipeTarget() *MonsterCard {
	tauntMonster := t.GetTauntMonster()
	if tauntMonster != nil {
		return tauntMonster
	}
	// no taunt monster
	backlineAliveMonsters := t.GetBacklineAliveMonsters()
	for _, m := range backlineAliveMonsters {
		canBeSniped := !m.HasAbility(ABILITY_CAMOUFLAGE) && (!m.HasAttack() || m.Ranged > 0 || m.Magic > 0)
		if canBeSniped {
			return &m
		}
	}

	// no backline snipe target
	firstAliveMonster := t.GetFirstAliveMonster()
	return firstAliveMonster
}

func (t *GameTeam) GetOpportunityTarget() *MonsterCard {
	tauntMonster := t.GetTauntMonster()
	if tauntMonster != nil {
		return tauntMonster
	}

	aliveMonsters := t.GetAliveMonsters()
	if len(aliveMonsters) == 0 {
		return nil
	}

	target := aliveMonsters[0]
	lowestHealth := target.Health
	for _, m := range aliveMonsters {
		if m.HasAbility(ABILITY_CAMOUFLAGE) {
			// first monster is already set as lowest health, so can ignore the first position camoflage monster
			continue
		}
		if m.Health < lowestHealth {
			target = m
			lowestHealth = m.Health
		}
	}
	return &target
}

func (t *GameTeam) GetSneakTarget() *MonsterCard {
	tauntMonster := t.GetTauntMonster()
	if tauntMonster != nil {
		return tauntMonster
	}

	aliveMonsters := t.GetAliveMonsters()
	numOfAliveMonsters := len(aliveMonsters)
	if numOfAliveMonsters == 0 {
		return nil
	}

	for i := len(aliveMonsters) - 1; i > 0; i-- {
		m := aliveMonsters[i]
		if !m.HasAbility(ABILITY_CAMOUFLAGE) {
			return &m
		}
	}

	return &aliveMonsters[0]
}

/** Which monster to repair, returns NULL if none. Repair target is the one that lost the most armor. */
func (t *GameTeam) GetRepairTarget() *MonsterCard {
	largestArmorDiff := 0
	monsterToRepair := nil

	for _, m := range t.GetAliveMonsters() {
		armorDiff := m.GetPostAbilityMaxArmor() - m.Armor
		if armorDiff > largestArmorDiff {
			largestArmorDiff = armorDiff
			monsterToRepair = m
		}
	}

	return monsterToRepair
}

/** Which backline monster to triage (monster who got the most damage), returns NULL if none */
func (t *GameTeam) GetTriageHealTarget() *MonsterCard {
	largestHealthDiff := 0
	var monsterToTriage *MonsterCard

	for _, m := range t.GetAliveMonsters() {
		healthDiff := m.GetPostAbilityMaxHealth() - m.Health
		if healthDiff > largestHealthDiff {
			largestHealthDiff = healthDiff
			monsterToTriage = &m
		}
	}

	return monsterToTriage
}

func (t *GameTeam) GetBacklineAliveMonsters() []MonsterCard {
	aliveMonsters := t.GetAliveMonsters()
	if len(aliveMonsters) <= 1 {
		return []MonsterCard{}
	}
	return aliveMonsters[1:len(aliveMonsters)]
}

func (t *GameTeam) GetTauntMonster() *MonsterCard {
	aliveMonsters := t.GetAliveMonsters()
	for _, monster := range aliveMonsters {
		if monster.HasAbility(ABILITY_TAUNT) {
			return &monster
		}
	}
	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"

	. "github.com/YukiUmetsu/go-spl-simulator/game_models"
)

const SPL_API_URL = "https://api2.splinterlands.com/"
const GET_ALL_CARDS_ENDPOIONT = "cards/get_details"
const BATTLE_HISTORY_ENDPOINT = "battle/result?id="

/*  Creates the game using the card id. Returns the battle logs. */
func SimulateBattle(battleId string, shouldLog bool) ([]BattleLog, error) {
	cardDetailMap := GetAllCardDetail()
	historicBattle := GetHistoricBattle(battleId)

	var battleDetails BattleDetails
	err := json.Unmarshal([]byte(historicBattle.Details), &battleDetails)
	if err != nil {
		log.Fatalln(err)
	}

	rulesetStrArr := strings.Split(historicBattle.Ruleset, "|")
	rulesets := make([]Ruleset, 0)
	for _, rulesetStr := range rulesetStrArr {
		rulesets = append(rulesets, Ruleset(rulesetStr))
	}
	game := CreateGame(cardDetailMap, battleDetails, rulesets, shouldLog)
	game.PlayGame()
	return game.GetBattleLogs()
}

func GetWinrateOfBattle(battleId string, playerNum int) float64 {
	winCount := 0
	cardDetailMap := GetAllCardDetail()
	historicBattle := GetHistoricBattle(battleId)

	var battleDetails BattleDetails
	err := json.Unmarshal([]byte(historicBattle.Details), &battleDetails)
	if err != nil {
		log.Fatalln(err)
	}

	rulesetStrArr := strings.Split(historicBattle.Ruleset, "|")
	rulesets := make([]Ruleset, 0)
	for _, rulesetStr := range rulesetStrArr {
		rulesets = append(rulesets, Ruleset(rulesetStr))
	}

	for i := 0; i < 100; i++ {
		game := CreateGame(cardDetailMap, battleDetails, rulesets, false)
		game.PlayGame()
		winner := game.GetWinner()
		if int(winner) == playerNum {
			winCount += 1
		}
	}
	return math.Round(float64(winCount) * 100 / 100)
}

func GetAllCardDetail() CardDetailMap {
	resp, err := http.Get(SPL_API_URL + GET_ALL_CARDS_ENDPOIONT)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	var cardDetails []CardDetail
	err = json.NewDecoder(resp.Body).Decode(&cardDetails)
	if err != nil {
		log.Fatal(err)
	}

	cardDetailMap := make(CardDetailMap)
	for _, cd := range cardDetails {
		cardDetailMap[cd.ID] = cd
	}
	return cardDetailMap
}

func GetHistoricBattle(battleID string) BattleHistory {
	resp, err := http.Get(SPL_API_URL + BATTLE_HISTORY_ENDPOINT + battleID)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	var bh BattleHistory
	err = json.NewDecoder(resp.Body).Decode(&bh)
	if err != nil {
		log.Fatal(err)
	}

	return bh
}

func CreateGameTeam(cardDetailMap CardDetailMap, battleTeam BattleTeam) *GameTeam {
	var summoner SummonerCard
	summonerDetail := cardDetailMap[battleTeam.Summoner.CardDetailID]
	summoner.Setup(summonerDetail, battleTeam.Summoner.Level)

	monsterList := make([]*MonsterCard, 0)
	for _, m := range battleTeam.Monsters {
		mDetail := cardDetailMap[m.CardDetailID]
		monster := MonsterCard{}
		monster.Setup(mDetail, m.Level)
		monsterList = append(monsterList, &monster)
	}
	var team GameTeam
	team.Create(&summoner, monsterList, battleTeam.Player)
	return &team
}

func CreateGame(cardDetailMap CardDetailMap, battleDetails BattleDetails, rulesets []Ruleset, shouldLog bool) Game {
	gameTeam1 := CreateGameTeam(cardDetailMap, battleDetails.Team1)
	gameTeam2 := CreateGameTeam(cardDetailMap, battleDetails.Team2)
	var game Game
	game.Create(gameTeam1, gameTeam2, rulesets, shouldLog)
	return game
}

func PrintStruct(value any) {
	jsonData, err := json.Marshal(&value)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(jsonData)
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	. "github.com/YukiUmetsu/go-spl-simulator/game_models"
)

const SPL_API_URL = "https://api2.splinterlands.com/"
const GET_ALL_CARDS_ENDPOIONT = "cards/get_details"
const BATTLE_HISTORY_ENDPOINT = "battle/result?id="

/*  Creates the game using the card id. Returns the battle logs. */
func ExampleHistoricBattle() ([]BattleLog, error) {
	cardDetailMap := GetAllCardDetail()
	historicBattle := GetHistoricBattle("sl_04e5faecf4e5996350fa97e1b94c9658")

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
	game := CreateGame(cardDetailMap, battleDetails, rulesets)
	game.PlayGame()
	return game.GetBattleLogs()
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

func CreateGameTeam(cardDetailMap CardDetailMap, battleTeam BattleTeam) GameTeam {
	summoner := SummonerCard{}
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
	team.Create(&summoner, monsterList)
	return team
}

func CreateGame(cardDetailMap CardDetailMap, battleDetails BattleDetails, rulesets []Ruleset) Game {
	gameTeam1 := CreateGameTeam(cardDetailMap, battleDetails.Team1)
	gameTeam2 := CreateGameTeam(cardDetailMap, battleDetails.Team2)
	var game Game
	game.Create(gameTeam1, gameTeam2, rulesets, true)
	return game
}

func PrintStruct(value any) {
	jsonData, err := json.Marshal(&value)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(jsonData)
}

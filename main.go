package simulator

import (
	"fmt"
	"log"
	"strings"
)

func main() {
	simulateMainBattle("sl_b0a5f35aa314b8d793b669a38f769964", false, 2)
}

func simulateMainBattle(battleId string, isforWinrate bool, playerNum int) {
	if !isforWinrate {
		logs, err := SimulateBattle(battleId, true)

		if err != nil {
			log.Fatalln("final error: ", err)
		}
		turnCount := 0
		actionCount := 0
		for i := 0; i < len(logs); i++ {
			l := logs[i]
			if turnCount > 100 {
				break
			}
			if strings.Contains(l.Action.String(), "Round") {
				turnCount += 1
				actionCount = 0
				fmt.Printf("\n--------------------------------------\n")
			}
			actionCount += 1
			fmt.Printf("\n%v Action: %v, Actor: %+v, Target: %+v, Value: %d\n", actionCount, l.Action, l.Actor, l.Target, l.Value)
		}
		return
	}

	winRate, playerName := GetWinrateOfBattle(battleId, playerNum)
	fmt.Printf("player: %s, winrate: %f\n", playerName, winRate)
}

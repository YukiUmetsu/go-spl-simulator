package main

import (
	"fmt"
	"log"
	"strings"
)

func main() {
	logs, err := SimulateBattle("sl_001312439e6630c7755cb58151db8eed", true)

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
}

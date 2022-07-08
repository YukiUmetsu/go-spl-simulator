package main

import (
	"fmt"
	"log"
	"strings"
)

func main() {
	logs, err := ExampleHistoricBattle()
	if err != nil {
		log.Fatalln("final error: ", err)
	}
	count := 0
	for i := 0; i < len(logs); i++ {
		l := logs[i]
		if count > 5 {
			break
		}
		if strings.Contains(string(l.Action), "Round") {
			count += 1
			fmt.Printf("\n--------------------------------------\n")
		}
		fmt.Printf("\nAction: %v, Actor: %+v, Target: %+v, Value: %d\n", l.Action, l.Actor, l.Target, l.Value)
	}
}

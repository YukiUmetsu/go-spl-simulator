package game_utils

import (
    "encoding/json"
    "golang.org/x/exp/constraints"
    "log"
    "fmt"
)

func Remove[T comparable](l []T, item T) []T {
    for i, other := range l {
        if other == item {
            return append(l[:i], l[i+1:]...)
        }
    }
    return l
}

func Contains[T comparable](s []T, e T) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func GetBigger[T constraints.Ordered](a T, b T) T {
    if (a > b) {
        return a
    }
    return b
}

func GetSmaller[T constraints.Ordered](a T, b T) T {
    if (a > b) {
        return b
    }
    return a
}

func StrArrContains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func PrintStructnWithLabel(label string, value any) {
	jsonData, err := json.Marshal(&value)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n%s %+v\n", label, string(jsonData))
}

func PrintStruct(value any) {
	jsonData, err := json.Marshal(&value)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))
}
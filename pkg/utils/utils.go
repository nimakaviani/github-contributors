package utils

import "fmt"

var Debug bool

func Log(msg ...string) {
	if Debug {
		for _, m := range msg {
			fmt.Printf("%s ", m)
		}
		println()
	}
}

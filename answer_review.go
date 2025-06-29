package main

import (
	"fmt"
)

func review(wrong []string) {
	for i := range wrong {
		fmt.Println(wrong[i])
	}
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// create a basic interactable visual
// however i do not know how to do this at this time
// sad!

func mc_draw(options [Choices]string, answer string) bool {
	hi_lited := options[0]
	for _, opt := range options {
		if opt == hi_lited {
			fmt.Printf("\x1b[7m>> %s <<\x1b[0m\n", opt)
		} else {
			fmt.Println(opt)
		}
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		in, err := reader.ReadByte()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", in)
		if in != 0 {
			break
		}
	}
	return false
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

const Choices int = 4

func quiz(ans string) bool {
	has_int := strings.ContainsAny(ans, "0123456789")
	if has_int {
		return get_mult_choic(ans)
	}
	scanr := bufio.NewScanner(os.Stdin)
	fmt.Print(">> ")
	scanr.Scan()
	response := scanr.Text()
	return response == ans
}

func main() {
	scanr := bufio.NewScanner(os.Stdin)
	fmt.Print("file path>> ")
	scanr.Scan()
	fpath := scanr.Text()
	file, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	line_count := 0
	fscanr := bufio.NewScanner(file)
	var line_slc = []string{}
	for fscanr.Scan() {
		line := fscanr.Text()
		line_slc = append(line_slc, line)
		line_count++
	}
	a_set := []string{}
	for i := range len(line_slc) {
		a := strings.Split(line_slc[i], ":")[1]
		a_set = append(a_set, a)
	}
	fmt.Print("# Questions to ask>> ")
	scanr.Scan()
	qstr := scanr.Text()
	qquant, qerr := strconv.Atoi(qstr)
	if qerr != nil {
		qquant = 0
	}
	counter := 0
	correct := false
	quit := false
	wrong_set := []string{}
	seed := rand.New(rand.NewSource(99))
	for !quit {
		if qquant == counter {
			quit = true
			continue
		}
		rand_q_i := seed.Intn(line_count)
		r_line := line_slc[rand_q_i]
		question := strings.Split(r_line, ":")[0]
		answer := strings.Split(r_line, ":")[1]
		fmt.Println(question)
		correct = quiz(answer)
		if !correct {
			wrong_set = append(wrong_set, r_line)
		}
		counter++
	}
}

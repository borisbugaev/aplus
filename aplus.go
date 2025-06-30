package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const Choices int = 4
const Choices_ext int = 6

func quiz(ans string, acro string) bool {
	is_cs := strings.Contains(ans, ",")
	if is_cs {
		return get_multi_answrs(ans, acro)
	}
	has_int := strings.ContainsAny(ans, "0123456789")
	has_acro := false
	acwrds := strings.Fields(acro)
	for i := range len(acwrds) {
		if strings.Contains(ans, acwrds[i]) {
			has_acro = true
		}
	}
	if has_int || has_acro {
		return get_mult_choic(ans, acro)
	}
	scanr := bufio.NewScanner(os.Stdin)
	fmt.Print(">> ")
	scanr.Scan()
	response := scanr.Text()
	return response == ans
}

func main() {
	file, err := os.Open("QUESTIONS.TXT")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	afile, err := os.Open("ACRONYMS.TXT")
	if err != nil {
		log.Fatal(err)
	}
	defer afile.Close()

	line_count := 0
	fscanr := bufio.NewScanner(file)
	var line_slc = []string{}
	for fscanr.Scan() {
		line := fscanr.Text()
		line_slc = append(line_slc, line)
		line_count++
	}
	fscanr = bufio.NewScanner(afile)
	fscanr.Scan()
	aln := fscanr.Text()
	scanr := bufio.NewScanner(os.Stdin)
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
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	prev_val := map[int]bool{}
	for !quit {
		if qquant == counter {
			quit = true
			continue
		}
		rand_q_i := seed.Intn(line_count)
		for prev_val[rand_q_i] {
			rand_q_i = seed.Intn(line_count)
		}
		prev_val[rand_q_i] = true
		r_line := line_slc[rand_q_i]
		question := strings.Split(r_line, ":")[0]
		answer := strings.Split(r_line, ":")[1]
		fmt.Println(question)
		correct = quiz(answer, aln)
		if !correct {
			wrong_set = append(wrong_set, r_line)
		}
		counter++
	}
	review(wrong_set)
}

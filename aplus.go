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

func quiz(ans string, typesmap map[string]string) bool {
	type_names := strings.Split(typesmap["DEFAULT"], ",")
	var my_type string = "NOTYPE"
	for _, name := range type_names {
		twords := strings.Split(typesmap[name], ",")
		for _, word := range twords {
			if strings.Contains(ans, word) {
				my_type = name
				break
			}
		}
	}
	is_cs := strings.Contains(ans, ",")
	if is_cs {
		return get_multi_answrs(ans, typesmap[my_type])
	}
	if my_type == "NOTYPE" {
		scanr := bufio.NewScanner(os.Stdin)
		fmt.Print(">> ")
		scanr.Scan()
		response := scanr.Text()
		return response == ans
	}

	return get_mult_choic(ans, typesmap[my_type])
}

func main() {
	q_file, err := os.Open("QUESTIONS.TXT")
	if err != nil {
		log.Fatal(err)
	}
	defer q_file.Close()
	use_file, err := os.Open("USE.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer use_file.Close()

	line_count := 0
	fscanr := bufio.NewScanner(q_file)
	var line_slc = []string{}
	for fscanr.Scan() {
		line := fscanr.Text()
		line_slc = append(line_slc, line)
		line_count++
	}
	fscanr = bufio.NewScanner(use_file)
	fscanr.Scan()
	types_csv := fscanr.Text()
	type_names := strings.Split(types_csv, ",")
	type_map := map[string]string{}
	type_map["DEFAULT"] = types_csv
	for _, typename := range type_names {
		t_f_name := fmt.Sprintf("USING/%s.csv", typename)
		t_file, err := os.Open(t_f_name)
		if err != nil {
			log.Fatal(err)
		}
		defer t_file.Close()
		fscanr = bufio.NewScanner(t_file)
		fscanr.Scan()
		type_map[typename] = fscanr.Text()
	}
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
		correct = quiz(answer, type_map)
		if !correct {
			wrong_set = append(wrong_set, r_line)
		}
		counter++
	}
	review(wrong_set)
}

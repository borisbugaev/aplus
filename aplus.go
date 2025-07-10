package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	printutils "github.com/borisbugaev/go_print_utils/printutils"
)

func print_quant(content string) int {
	return printutils.Print_Lines(content)
}

func clear_lines(count int) {
	printutils.Clear_Lines(count)
}

const Choices int = 4
const Choices_ext int = 6

func quiz(ans string, typesmap map[string]string) bool {
	var my_type string = "NOTYPE"
	type_names := strings.SplitSeq(typesmap["DEFAULT"], ",")
	for name := range type_names {
		if name == "DEFAULT" {
			continue
		}
		twords := strings.SplitSeq(typesmap[name], ",")
		for word := range twords {
			if strings.Contains(ans, word) {
				my_type = name
				break
			}
		}
	}
	if my_type == "NOTYPE" {
		scanr := bufio.NewScanner(os.Stdin)
		line_count := print_quant(">> ")
		scanr.Scan()
		response := scanr.Text()
		response = strings.Trim(response, " ")
		response = strings.ToLower(response)
		ans = strings.ToLower(ans)
		clear_lines(line_count)
		return response == ans
	}
	is_cs := strings.Contains(ans, ",")
	if is_cs {
		return get_multi_answrs(ans, typesmap[my_type])
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

	file_line_count := 0
	fscanr := bufio.NewScanner(q_file)
	var line_slc = []string{}
	for fscanr.Scan() {
		line := fscanr.Text()
		line_slc = append(line_slc, line)
		file_line_count++
	}
	q_file.Close()
	fscanr = bufio.NewScanner(use_file)
	fscanr.Scan()
	types_csv := fscanr.Text()
	type_names := strings.Split(types_csv, ",")
	type_map := map[string]string{}
	use_file.Close()
	type_map["DEFAULT"] = types_csv
	for _, typename := range type_names {
		if typename == "" || typename == "DEFAULT" {
			continue
		} else {
			t_f_name := fmt.Sprintf("USING/%s.csv", typename)
			t_file, err := os.Open(t_f_name)
			if err != nil {
				log.Fatal(err)
			}
			defer t_file.Close()
			fscanr = bufio.NewScanner(t_file)
			fscanr.Scan()
			type_map[typename] = fscanr.Text()
			t_file.Close()
		}
	}
	mqs_ptr := flag.Bool("mqs", false, "if true, question count is runtime defined")
	force_ptr := flag.Int("force", -1, "force first question")
	pedant_ptr := flag.Bool("pedant", false, "provide answer checks immediately")
	flag.Parse()
	qquant := 30 // else question count defaults to 30
	if *mqs_ptr {
		scanr := bufio.NewScanner(os.Stdin)
		fmt.Print("# Questions to ask>> ")
		scanr.Scan()
		qstr := scanr.Text()
		iqquant, qerr := strconv.Atoi(qstr)
		if qerr != nil {
			qquant = 0
		}
		qquant = iqquant
	}
	counter := 0
	correct := false
	quit := false
	wrong_set := []string{}
	seed := rand.New(rand.NewSource(time.Now().UnixMilli()))
	prev_val := map[int]bool{}
	for !quit {
		if qquant == counter {
			quit = true
			continue
		}
		if counter == 0 && counter < *force_ptr {
			line := line_slc[*force_ptr]
			prev_val[*force_ptr] = true
			q := strings.Split(line, ":")[0]
			fmt.Println(q)
			a := strings.Split(line, ":")[1]
			correct = quiz(a, type_map)
			if !correct {
				wrong_set = append(wrong_set, line)
				if *pedant_ptr {
					fmt.Print("\x1b[1A\x1b[93;41m>>\x1b[G\n\x1b[0m")
					fmt.Printf("\x1b[3m>> %s\x1b[0m\n", a)
				}
			}
		}
		rand_q_i := seed.Intn(file_line_count)
		for prev_val[rand_q_i] {
			rand_q_i = seed.Intn(file_line_count)
		}
		prev_val[rand_q_i] = true
		r_line := line_slc[rand_q_i]
		question := strings.Split(r_line, ":")[0]
		answer := strings.Split(r_line, ":")[1]
		question = strings.Trim(question, " ")
		answer = strings.Trim(answer, " ")
		question = fmt.Sprintln(question)
		line_count := print_quant(question)
		correct = quiz(answer, type_map)
		if !correct {
			wrong_set = append(wrong_set, r_line)
			if *pedant_ptr {
				line_count += print_quant("\x1b[1A\x1b[93;41m>>\x1b[G\n\x1b[0m")
				line_count += print_quant(fmt.Sprintf("\x1b[3m>> %s\x1b[0m\n", answer))
			}
		}
		clear_lines(line_count)
		counter++
	}
	review(wrong_set)
}

package main

import (
	"bufio"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	printutils "github.com/borisbugaev/go_print_utils/printutils"
)

func true_false(answer bool) bool {
	line_count := print_quant("[T]RUE\t\t\t[F]ALSE\n>> ")
	scnr := bufio.NewScanner(os.Stdin)
	scnr.Scan()
	response := strings.ToLower(scnr.Text())
	response_bool := strings.Contains(response, "t")
	correct := response_bool == answer
	clear_lines(line_count)
	return correct
}

func q_concat(a [Choices]string, b [Choices]string) [Choices]string {
	for i := range Choices {
		if strings.Contains(b[i], "\a") {
			continue
		} else if a[i] == b[i] {
			continue
		} else if i%2 == 0 {
			a[i] = b[i]
			continue
		} else {
			continue
		}
	}
	return a
}

func get_multi_answrs(answer string, of_type string) bool {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	answrs := strings.Split(answer, ",")
	optns := []string{}
	set := map[string]bool{}
	for i := range len(answrs) {
		optns = append(optns, answrs[i])
		set[answrs[i]] = true
	}
	answer_set := set
	words_of_type := strings.Split(of_type, ",")
	for range Choices_ext {
		to_insert := words_of_type[seed.Intn(len(words_of_type))]
		_, includes := set[to_insert]
		for includes {
			to_insert = words_of_type[seed.Intn(len(words_of_type))]
			_, includes = set[to_insert]
		}
		set[to_insert] = true
		optns = append(optns, to_insert)
		if len(optns) == Choices_ext {
			break
		}
	}
	seed.Shuffle(len(optns), func(i, j int) {
		optns[i], optns[j] = optns[j], optns[i]
	})
	response_cs := printutils.Line_Select_MC(optns)
	r_seq := strings.SplitSeq(response_cs, ",")
	counter := 0
	correct := true
	for response := range r_seq {
		if response == "\a" {
			correct = false
			continue
		}
		_, includes := answer_set[response]
		if !includes {
			correct = false
			continue
		}
		counter++
	}
	if counter != len(answrs) {
		correct = false
	}
	return correct
}

func mlt_chc_i_rndmz(txt []string, answer_value int) [Choices]string {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	set := map[int]bool{}
	order := [Choices]int{}
	for i := range Choices {
		ith_value := seed.Intn(Choices)
		_, includes := set[ith_value]
		for includes {
			ith_value = seed.Intn(Choices)
			_, includes = set[ith_value]
		}
		set[ith_value] = true
		order[i] = ith_value
	}
	set = map[int]bool{} //clear set
	values := [Choices]int{}
	values[0] = answer_value
	rndm_range := min(12, answer_value)
	j := 1
	for j < Choices {
		diff := seed.Intn(rndm_range)
		if seed.Intn(2) > 0 {
			diff *= -1
		}
		_, includes := set[diff]
		for includes {
			diff := seed.Intn(rndm_range)
			if seed.Intn(2) > 0 {
				diff *= -1
			}
			_, includes = set[diff]
		}
		set[diff] = true
		values[j] = answer_value + diff
		j++
	}
	options := [Choices]string{}
	for i := range Choices {
		options[order[i]] = txt[0] + strconv.Itoa(values[i]) + txt[1]
	}
	return options
}

func mlt_chc_acr_r(txt []string, answer string, of_type string) [Choices]string {
	words_of_type := strings.Split(of_type, ",")
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	set := map[int]bool{}
	order := [Choices]int{}
	for i := range Choices {
		ith_val := seed.Intn(Choices)
		_, includes := set[ith_val]
		for includes {
			ith_val = seed.Intn(Choices)
			_, includes = set[ith_val]
		}
		set[ith_val] = true
		order[i] = ith_val
	}
	set = map[int]bool{} //clear set
	vals := [Choices]int{}
	vals[0] = slices.Index(words_of_type, answer)
	set[vals[0]] = true
	j := 1
	for j < min(Choices, len(words_of_type)) {
		index := seed.Intn(len(words_of_type))
		_, includes := set[index]
		for includes {
			index = seed.Intn(len(words_of_type))
			_, includes = set[index]
		}
		set[index] = true
		vals[j] = index
		j++
	}
	out_acrw := [Choices]string{}
	for i := range Choices {
		if !set[vals[i]] {
			out_acrw[i] = "\a"
		} else {
			out_acrw[i] = words_of_type[vals[i]]
		}
	}
	optns := [Choices]string{}
	for i := range Choices {
		optns[order[i]] = txt[0] + out_acrw[i] + txt[1]
	}
	return optns
}

func mc_caller(answer string, options []string, my_type string) bool {
	if my_type == "headless" {
		// pass options somewhere that gets a selection
	} else {
		selection := printutils.Line_Select_MC(options)
		if answer == selection {
			return true
		}
	}
	return false
}

func get_mult_choic(answer string, of_type string, strict bool) bool {
	words := strings.Fields(answer)
	words_of_type := strings.Split(of_type, ",")
	if len(words_of_type) == 2 && words_of_type[0] == "True" {
		correct_answer := words[0] == words_of_type[0]
		return true_false(correct_answer)
	}
	intoptns := [Choices]string{"DEFAULT"}
	for i := range len(words) {
		if strict {
			break
		}
		w_num, err := strconv.Atoi(words[i])
		if err == nil {
			if len(words) == 1 {
				txt := []string{"", ""}
				intoptns = mlt_chc_i_rndmz(txt, w_num)
				break
			} else {
				txt := []string{"", ""}
				var cut_success bool
				txt[0], txt[1], cut_success = strings.Cut(answer, words[i])
				if cut_success {
					intoptns = mlt_chc_i_rndmz(txt, w_num)
					break
				}
			}
		}
	}
	acroptns := [Choices]string{"DEFAULT"}
	for j := range len(words_of_type) {
		if strings.Contains(answer, words_of_type[j]) {
			if len(words) == 1 && answer == words_of_type[j] {
				txt := []string{"", ""}
				acroptns = mlt_chc_acr_r(txt, words_of_type[j], of_type)
				break
			} else {
				txt := []string{"", ""}
				var cut_success bool
				txt[0], txt[1], cut_success = strings.Cut(answer, words_of_type[j])
				if cut_success {
					acroptns = mlt_chc_acr_r(txt, words_of_type[j], of_type)
					break
				}
			}
		}
	}
	optns := [Choices]string{"DEFAULT"}
	if intoptns[0] != "DEFAULT" && acroptns[0] != "DEFAULT" {
		// concat
		optns = q_concat(intoptns, acroptns)
	} else if intoptns[0] != "DEFAULT" {
		optns = intoptns
	} else if acroptns[0] != "DEFAULT" {
		optns = acroptns
	} else {
		// default case, should not occur
	}
	// print options and get answer
	return mc_caller(answer, optns[:], "")
}

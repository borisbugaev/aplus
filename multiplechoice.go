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

func true_false(a bool) bool {
	line_count := print_quant("TRUE\t\t\tFALSE\n>> ")
	scnr := bufio.NewScanner(os.Stdin)
	scnr.Scan()
	answr := scnr.Text()
	answr = strings.ToLower(answr)
	if answr == "t" {
		answr = "true"
	}
	if answr == "f" {
		answr = "false"
	}
	correct := false
	if answr == "true" && a {
		correct = true
	} else if answr == "false" && !a {
		correct = true
	}
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

func get_multi_answrs(ans string, acro string) bool {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	answrs := strings.Split(ans, ",")
	optns := []string{}
	set := map[string]bool{}
	for i := range len(answrs) {
		optns = append(optns, answrs[i])
		set[answrs[i]] = true
	}
	a_set := set
	acrw := strings.Split(acro, ",")
	for range Choices_ext {
		to_insert := acrw[seed.Intn(len(acrw))]
		_, includes := set[to_insert]
		for includes {
			to_insert = acrw[seed.Intn(len(acrw))]
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
		_, includes := a_set[response]
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

func mlt_chc_i_rndmz(txt []string, val int) [Choices]string {
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
	vals[0] = val
	rndm_range := min(12, val)
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
		vals[j] = val + diff
		j++
	}
	opts := [Choices]string{}
	for i := range Choices {
		opts[order[i]] = txt[0] + strconv.Itoa(vals[i]) + txt[1]
	}
	return opts
}

func mlt_chc_acr_r(txt []string, ans string, acro string) [Choices]string {
	acrw := strings.Split(acro, ",")
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
	vals[0] = slices.Index(acrw, ans)
	set[vals[0]] = true
	j := 1
	for j < min(Choices, len(acrw)) {
		index := seed.Intn(len(acrw))
		_, includes := set[index]
		for includes {
			index = seed.Intn(len(acrw))
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
			out_acrw[i] = acrw[vals[i]]
		}
	}
	optns := [Choices]string{}
	for i := range Choices {
		optns[order[i]] = txt[0] + out_acrw[i] + txt[1]
	}
	return optns
}

func get_mult_choic(ans string, acro string) bool {
	words := strings.Fields(ans)
	acrw := strings.Split(acro, ",")
	if len(acrw) == 2 && acrw[0] == "True" {
		correct_answer := words[0] == acrw[0]
		return true_false(correct_answer)
	}
	intoptns := [Choices]string{"DEFAULT"}
	for i := range len(words) {
		w_num, err := strconv.Atoi(words[i])
		if err == nil {
			if len(words) == 1 {
				txt := []string{"", ""}
				intoptns = mlt_chc_i_rndmz(txt, w_num)
				break
			} else {
				txt := []string{"", ""}
				var cut_success bool
				txt[0], txt[1], cut_success = strings.Cut(ans, words[i])
				if cut_success {
					intoptns = mlt_chc_i_rndmz(txt, w_num)
					break
				}
			}
		}
	}
	acroptns := [Choices]string{"DEFAULT"}
	for j := range len(acrw) {
		if strings.Contains(ans, acrw[j]) {
			if len(words) == 1 && ans == acrw[j] {
				txt := []string{"", ""}
				acroptns = mlt_chc_acr_r(txt, acrw[j], acro)
				break
			} else {
				txt := []string{"", ""}
				var cut_success bool
				txt[0], txt[1], cut_success = strings.Cut(ans, acrw[j])
				if cut_success {
					acroptns = mlt_chc_acr_r(txt, acrw[j], acro)
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
	response := printutils.Line_Select_MC(optns[:])
	correct := false
	if response == ans {
		correct = true
	}
	return correct
}

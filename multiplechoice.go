package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

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
	acrw := strings.Fields(acro)
	for range Choices_ext {
		to_insert := acrw[seed.Intn(len(acrw))]
		_, includes := set[to_insert]
		for includes {
			to_insert := acrw[seed.Intn(len(answrs))]
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
	out_optns := [Choices_ext]string{}
	mp_optns := map[string]string{}
	for i := range Choices_ext {
		lttr := fmt.Sprintf("%c", 'A'+i)
		out_optns[i] = fmt.Sprintf("%s: %s", lttr, optns[i])
		mp_optns[lttr] = optns[i]
		llttr := fmt.Sprintf("%c", 'a'+i)
		mp_optns[llttr] = optns[i]
	}
	for i := range Choices_ext {
		fmt.Printf("%s\n", out_optns[i])
	}
	scnnr := bufio.NewScanner(os.Stdin)
	fmt.Print(">> ")
	scnnr.Scan()
	my_nswr := scnnr.Text()
	sl_my := strings.Fields(my_nswr)
	correct := true
	for i := range len(sl_my) {
		if len(sl_my) != len(answrs) {
			correct = false
			break
		}
		cstr, includes := mp_optns[sl_my[i]]
		if !includes {
			correct = false
			break
		}
		_, includes = a_set[cstr]
		if !includes {
			correct = false
			break
		}
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
	}
	opts := [Choices]string{}
	for i := range Choices {
		opts[order[i]] = txt[0] + strconv.Itoa(vals[i]) + txt[1]
	}
	return opts
}

func mlt_chc_acr_r(txt []string, ans string, acro string) [Choices]string {
	acrw := strings.Fields(acro)
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
	for j < Choices {
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
	optns := [Choices]string{}
	for i := range Choices {
		optns[order[i]] = txt[0] + acrw[vals[i]] + txt[1]
	}
	return optns
}

func get_mult_choic(ans string, acro string) bool {
	optns := [Choices]string{}
	acrw := strings.Fields(acro)
	words := strings.Fields(ans)
	for i := range len(words) {
		w_num, err := strconv.Atoi(words[i])
		if err == nil {
			if len(words) == 1 {
				txt := []string{"", ""}
				optns = mlt_chc_i_rndmz(txt, w_num)
			}
			splt := strings.Split(ans, words[i])
			if len(splt) == 2 {
				optns = mlt_chc_i_rndmz(splt, w_num)
			}
		}
	}
	for j := range len(acrw) {
		if strings.Contains(ans, acrw[j]) {
			if len(ans) == 1 {
				txt := []string{"", ""}
				optns = mlt_chc_acr_r(txt, acrw[j], acro)
			}
			splt := strings.Split(ans, acrw[j])
			if len(splt) == 2 {
				optns = mlt_chc_acr_r(splt, acrw[j], acro)
			}
		}
	}

	// print options and get answer
	out_optns := [Choices]string{}
	mp_optns := map[string]string{}
	for i := range Choices {
		lttr := fmt.Sprintf("%c", 'A'+i)
		out_optns[i] = fmt.Sprintf("%s: %s", lttr, optns[i])
		mp_optns[lttr] = optns[i]
		llttr := fmt.Sprintf("%c", 'a'+i)
		mp_optns[llttr] = optns[i]
	}
	for i := range Choices {
		fmt.Printf("%s\n", out_optns[i])
	}
	scnnr := bufio.NewScanner(os.Stdin)
	fmt.Print(">> ")
	scnnr.Scan()
	my_nswr := scnnr.Text()
	cstr := mp_optns[my_nswr]
	correct := false
	if cstr == ans {
		correct = true
	}
	return correct
}

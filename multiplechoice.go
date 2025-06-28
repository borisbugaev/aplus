package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func mlt_chc_i_rndmz(txt []string, val int) [Choices]string {
	seed := rand.New(rand.NewSource(98))
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

func get_mult_choic(ans string) bool {
	optns := [Choices]string{}
	words := strings.Fields(ans)
	for i := range len(words) {
		w_num, err := strconv.Atoi(words[i])
		if err == nil {
			splt := strings.Split(ans, words[i])
			if len(splt) == 2 {
				optns = mlt_chc_i_rndmz(splt, w_num)
			}
		}
	}

	// print options and get answer
	out_optns := [Choices]string{}
	mp_optns := map[string]string{}
	for i := range Choices {
		lttr := fmt.Sprintf("%v", 'A'+i)
		out_optns[i] = fmt.Sprintf("%s: %s", lttr, optns[i])
		mp_optns[lttr] = optns[i]
		llttr := fmt.Sprintf("%v", 'a'+i)
		mp_optns[llttr] = optns[i]
	}
	if len(ans) < 12 {
		fmt.Printf("%s\t%s\n", out_optns[0], out_optns[2])
		fmt.Printf("%s\t%s\n", out_optns[1], out_optns[3])
	} else {
		for i := range Choices {
			fmt.Printf("%s\n", out_optns[i])
		}
	}
	scnnr := bufio.NewScanner(os.Stdin)
	fmt.Print(">> ")
	scnnr.Scan()
	my_nswr := scnnr.Text()
	cstr, err := mp_optns[my_nswr]
	correct := false
	for err {
		fmt.Print(">> ")
		scnnr.Scan()
		my_nswr := scnnr.Text()
		cstr, err = mp_optns[my_nswr]
	}
	if cstr == ans {
		correct = true
	}
	return correct
}

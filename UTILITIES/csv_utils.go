package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	printutils "github.com/borisbugaev/go_print_utils/printutils"
)

func clear_lines(count int) {
	printutils.Clear_Lines(count)
}

func print_quant(lines string) int {
	return printutils.Print_Lines(lines)
}

func prune_from_file(dir string, file string, content string) {

}

func pruner(dir string, files string, content string) {
	line_count, index := 1, 0
	lines := fmt.Sprintf("value %s in files...\n", content)
	options := fmt.Sprintf("%s,cancel", files)
	opt_seq := strings.SplitSeq(options, ",")
	opt_at := map[string]string{}
	for opt := range opt_seq {
		lttr := fmt.Sprintf("%c", 'A'+index)
		lines = fmt.Sprintf("%s%s> %s\n", lines, lttr, opt)
		opt_at[lttr] = opt
		index++
	}
	lines = fmt.Sprintf("%sdelete %s from...\n>> ", lines, content)
	line_count += print_quant(lines)
	scnr := bufio.NewScanner(os.Stdin)
	scnr.Scan()
	opt_selected, ok := opt_at[strings.ToUpper(scnr.Text())]
	if !ok {
		opt_selected = "cancel"
	}
	clear_lines(line_count)
	if opt_selected == "cancel" {
		return
	}
	prune_from_file(dir, opt_selected, content)
}

func overlap(my_dir string, prune bool) {
	to_console := fmt.Sprintf("Compare CSV in dir %s:", my_dir)
	fmt.Printf("%s", to_console)
	if strings.Contains(to_console, "\a") {
		scnr := bufio.NewScanner(os.Stdin)
		scnr.Scan()
		my_dir = scnr.Text()
	}
	files, err := os.ReadDir(my_dir)
	if err != nil {
		log.Fatal(err)
	}
	contents := map[string]string{}
	names := []string{}
	for _, file := range files {
		if file.Name() == "ALL.csv" {
			continue
		}
		location := fmt.Sprintf("%s\\%s", my_dir, file.Name())
		my_file, err := os.Open(location)
		if err != nil {
			log.Fatal(err)
		}
		defer my_file.Close()
		scnr := bufio.NewScanner(my_file)
		scnr.Scan()
		contents[file.Name()] = scnr.Text()
		my_file.Close()
		names = append(names, file.Name())
	}
	content_set := map[string]string{}
	entry_set_str := "\a"
	for _, name := range names {
		name_cont_seq := strings.SplitSeq(contents[name], ",")
		for item := range name_cont_seq {
			c_str, includes := content_set[item]
			if includes {
				c_str = fmt.Sprintf("%s,%s", c_str, name)
			} else {
				c_str = name
				entry_set_str = fmt.Sprintf("%s,%s", entry_set_str, item)
			}
			content_set[item] = c_str
		}
	}
	entry_seq := strings.SplitSeq(entry_set_str, ",")
	for entry := range entry_seq {
		if strings.Contains(entry, "\a") {
			continue
		} else if strings.Contains(content_set[entry], ",") {
			if prune {
				pruner(my_dir, content_set[entry], entry)
				return
			}
			out_str := fmt.Sprintf("%s in files %s\n", entry, content_set[entry])
			fmt.Print(out_str)
		}
	}
}

func initialize(repeats bool) {
	q_file, err := os.Open("../QUESTIONS.TXT")
	if err != nil {
		log.Fatal(err)
	}
	defer q_file.Close()
	fscanr := bufio.NewScanner(q_file)
	var line_slc = []string{}
	for fscanr.Scan() {
		line := fscanr.Text()
		line_slc = append(line_slc, line)
	}
	q_file.Close()
	all_file, err := os.Create("../USING/ALL.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer all_file.Close()
	all_str := ""
	for _, line := range line_slc {
		answer := strings.Split(line, ":")[1]
		answer = strings.Trim(answer, " ")
		all_str = fmt.Sprintf("%s%s,", all_str, answer)
	}
	if !repeats {
		ans := strings.SplitSeq(all_str, ",")
		an_set := map[string]bool{}
		no_repeat := ""
		for a := range ans {
			_, includes := an_set[a]
			if !includes {
				no_repeat = fmt.Sprintf("%s%s,", no_repeat, a)
			}
			an_set[a] = true
		}
		all_str = no_repeat
	}
	all_str = strings.Trim(all_str, ",")
	all_file.WriteString(all_str)
	all_file.Sync()
	all_file.Close()
}

func new(dir string, contents map[string]string) string {
	line_count := 0
	lines := "name of new .csv\n>> "
	line_count += print_quant(lines)
	scnr := bufio.NewScanner(os.Stdin)
	scnr.Scan()
	my_text := fmt.Sprintf("%s/%s.csv", dir, strings.ToUpper(scnr.Text()))
	_, incl := contents[my_text]
	for incl {
		fmt.Print(">> ")
		line_count++
		scnr.Scan()
		my_text = fmt.Sprintf("%s/%s.csv", dir, strings.ToUpper(scnr.Text()))
		_, incl = contents[my_text]
	}
	clear_lines(line_count)
	return my_text
}

func sort(dir string, to_sort string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	names := []string{}
	contents := map[string]string{}
	for _, file := range files {
		location := fmt.Sprintf("%s/%s", dir, file.Name())
		my_file, err := os.Open(location)
		if err != nil {
			log.Fatal(err)
		}
		defer my_file.Close()
		scnr := bufio.NewScanner(my_file)
		scnr.Scan()
		contents[file.Name()] = scnr.Text()
		my_file.Close()
		names = append(names, file.Name())
	}
	else_subset := map[string]bool{}
	for _, name := range names {
		if name == to_sort {
			continue
		}
		name_cont_seq := strings.SplitSeq(contents[name], ",")
		for item := range name_cont_seq {
			_, includes := else_subset[item]
			if !includes {
				else_subset[item] = true
			}
		}
	}
	all_seq := strings.SplitSeq(contents[to_sort], ",")
	to_write := map[string]bool{}
	for entry := range all_seq {
		_, includes := else_subset[entry]
		if includes {
			continue
		}
		line_count := 1
		lines := fmt.Sprintf("in %s--Entry %s: place into...\n", to_sort, entry)
		line_count += print_quant(lines)
		to_poll := []string{}
		default_options := []string{"new...", "skip", "exit"}
		for _, name := range names {
			if name == to_sort || name == "ALL.csv" {
				continue
			}
			to_poll = append(to_poll, name)
		}
		for _, option := range default_options {
			to_poll = append(to_poll, option)
		}
		scnr := bufio.NewScanner(os.Stdin)
		resp := printutils.Line_Select_MC(to_poll)
		for resp == "exit" {
			line_count += print_quant("save? y/n\n>>")
			scnr.Scan()
			save := strings.ToLower(scnr.Text())
			if save == "y" {
				break
			} else if save == "n" {
				clear_lines(line_count)
				return
			}
		}
		if resp == "exit" {
			clear_lines(line_count)
			break
		}
		if resp == "skip" {
			clear_lines(line_count)
			continue
		}
		if resp == "new..." {
			new_name := new(dir, contents)
			nfile, err := os.Create(new_name)
			if err != nil {
				log.Fatal(err)
			}
			defer nfile.Close()
			nfile.WriteString(entry)
			names = append(names, new_name)
			contents[new_name] = entry
			else_subset[entry] = true
			clear_lines(line_count)
			continue
		}
		contents[resp] = fmt.Sprintf("%s,%s", contents[resp], entry)
		_, exists := to_write[resp]
		if !exists {
			to_write[resp] = true
		}
		clear_lines(line_count)
	}
	for target_file := range to_write {
		location := fmt.Sprintf("%s/%s", dir, target_file)
		// currently this rewrites the entire file modified. should be changed to append.
		file, err := os.Create(location)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		file.WriteString(contents[target_file])
		file.Close()
	}
}

func main() {
	dir_ptr := flag.String("dir", "../USING", "directory of CSV files")
	o_ptr := flag.Bool("o", false, "should utilities run overlap?")
	prune_ptr := flag.Bool("prune", false, "should overlap prune files")
	init_ptr := flag.Bool("init", false, "generate initial csv of answers")
	re_ptr := flag.Bool("re", false, "include repeats")
	sort_ptr := flag.String("sort", "", "run sort function on designated csv")
	flag.Parse()
	if *o_ptr {
		overlap(*dir_ptr, *prune_ptr)
	}
	if *init_ptr {
		initialize(*re_ptr)
	}
	if *sort_ptr != "" {
		sort_f := fmt.Sprintf("%s.csv", *sort_ptr)
		sort(*dir_ptr, sort_f)
	}
}

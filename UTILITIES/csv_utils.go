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

const IGNORED_FILES string = "QUESTIONS.TXT,USE.CSV"

func clear_lines(count int) {
	printutils.Clear_Lines(count)
}

func print_quant(lines string) int {
	return printutils.Print_Lines(lines)
}

func is_any(line string, cs_set string) bool {
	cs_seq := strings.SplitSeq(cs_set, ",")
	for str := range cs_seq {
		if str == line {
			return true
		}
	}
	return false
}

func prune_from_file(dir string, filename string, content string) {
	location := fmt.Sprintf("%s/%s", dir, filename)
	file, err := os.Open(location)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scnr := bufio.NewScanner(file)
	scnr.Scan()
	pruned_str := ""
	entry_seq := strings.SplitSeq(scnr.Text(), ",")
	file.Close()
	for entry := range entry_seq {
		if entry == content {
			continue
		}
		pruned_str = fmt.Sprintf("%s,%s", pruned_str, entry)
	}
	pruned_str = strings.Trim(pruned_str, ",")
	pfile, err := os.Create(location)
	if err != nil {
		log.Fatal(err)
	}
	defer pfile.Close()
	pfile.WriteString(pruned_str)
	pfile.Close()
}

func pruner(dir string, files string, content string) {
	line_count := 1
	lines := fmt.Sprintf("value %s in files...\n", content)
	options := fmt.Sprintf("%s,cancel", files)
	opt_seq := strings.SplitSeq(options, ",")
	opts := []string{}
	for opt := range opt_seq {
		opts = append(opts, opt)
	}
	lines = fmt.Sprintf("%sdelete %s from...\n", lines, content)
	line_count += print_quant(lines)
	opt_selected := printutils.Line_Select_MC(opts)
	clear_lines(line_count)
	if opt_selected == "cancel" || opt_selected == "\a" {
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
	to_ignore := fmt.Sprintf("%s,ALL.csv", IGNORED_FILES)
	for _, file := range files {
		if is_any(file.Name(), to_ignore) {
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
				continue
			}
			out_str := fmt.Sprintf("%s in files %s\n", entry, content_set[entry])
			fmt.Print(out_str)
		}
	}
}

func initialize(dir string, repeats bool) {
	question_location := fmt.Sprintf("%s/QUESTIONS.TXT", dir)
	q_file, err := os.Open(question_location)
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
	all_location := fmt.Sprintf("%s/ALL.csv", dir)
	all_file, err := os.Create(all_location)
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
		answer_seq := strings.SplitSeq(all_str, ",")
		answer_set := map[string]bool{}
		no_repeat := ""
		for a := range answer_seq {
			_, includes := answer_set[a]
			if !includes {
				no_repeat = fmt.Sprintf("%s%s,", no_repeat, a)
			}
			answer_set[a] = true
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
		line_count += print_quant(">> ")
		scnr.Scan()
		my_text = fmt.Sprintf("%s/%s.csv", dir, strings.ToUpper(scnr.Text()))
		_, incl = contents[my_text]
	}
	clear_lines(line_count)
	return my_text
}

func sort(dir string, to_sort string, move bool) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	names := []string{}
	contents := map[string]string{}
	for _, file := range files {
		if is_any(file.Name(), IGNORED_FILES) {
			continue
		}
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
	to_ignore := fmt.Sprintf("%s,ALL.csv,%s", IGNORED_FILES, to_sort)
	for _, name := range names {
		if is_any(name, to_ignore) {
			continue
		}
		name_content_seq := strings.SplitSeq(contents[name], ",")
		for item := range name_content_seq {
			_, includes := else_subset[item]
			if !includes {
				else_subset[item] = true
			}
		}
	}
	all_seq := strings.SplitSeq(contents[to_sort], ",")
	to_sort_contents_copy := contents[to_sort]
	to_write := map[string]bool{}
	if move {
		to_write[to_sort] = true
	}
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
			if is_any(name, to_ignore) {
				continue
			}
			to_poll = append(to_poll, name)
		}
		to_poll = append(to_poll, default_options...)
		scnr := bufio.NewScanner(os.Stdin)
		response := printutils.Line_Select_MC(to_poll)
		for response == "exit" {
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
		if response == "exit" {
			clear_lines(line_count)
			break
		}
		if response == "skip" {
			clear_lines(line_count)
			continue
		}
		if response == "new..." {
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
		contents[response] = fmt.Sprintf("%s,%s", contents[response], entry)
		if move {
			splits := strings.Split(to_sort_contents_copy, entry)
			for i, split := range splits {
				split = strings.Trim(split, ",")
				if i == 0 {
					to_sort_contents_copy = split
					continue
				} else {
					to_sort_contents_copy = fmt.Sprintf("%s,%s", to_sort_contents_copy, split)
				}
			}
		}
		_, exists := to_write[response]
		if !exists {
			to_write[response] = true
		}
		clear_lines(line_count)
	}
	if move {
		contents[to_sort] = to_sort_contents_copy
	}
	for target_file := range to_write {
		location := fmt.Sprintf("%s/%s", dir, target_file)
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
	mov_ptr := flag.Bool("mov", false, "remove sorted file from original location")
	flag.Parse()
	if *o_ptr {
		overlap(*dir_ptr, *prune_ptr)
	}
	if *init_ptr {
		initialize(*dir_ptr, *re_ptr)
	}
	if *sort_ptr != "" {
		sort_f := fmt.Sprintf("%s.csv", *sort_ptr)
		sort(*dir_ptr, sort_f, *mov_ptr)
	}
}

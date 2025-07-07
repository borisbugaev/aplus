package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func overlap(my_dir string) {
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

func main() {
	dir_ptr := flag.String("dir", "../USING", "directory of CSV files")
	o_ptr := flag.Bool("o", false, "should utilities run overlap?")
	init_ptr := flag.Bool("init", false, "generate initial csv of answers")
	re_ptr := flag.Bool("re", false, "include repeats")
	flag.Parse()
	my_dir := *dir_ptr
	if *o_ptr {
		overlap(my_dir)
	}
	if *init_ptr {
		initialize(*re_ptr)
	}
}

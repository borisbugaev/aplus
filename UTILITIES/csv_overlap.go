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

func main() {
	dir_ptr := flag.String("dir", "\a", "directory of CSV files")
	o_ptr := flag.Bool("o", false, "should utilities run overlap?")
	flag.Parse()
	my_dir := *dir_ptr
	if *o_ptr {
		overlap(my_dir)
	}
}

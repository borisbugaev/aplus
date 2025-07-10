# UTILITIES for "APLUS"

### tools for managing .csv input used by the APLUS quiz program

## USAGE

The primary way to interact with this program is with runtime flags in the command line.

### Flags:

**directory** *-dir* (string)

- manually set directory of .csv files (*-dir=$directoryname*)
- if not invoked, default directory is ../USING

**initialize** *-init* (bool) *-re* (bool)

- creates file "ALL.csv"
- "ALL.csv" is a list of all answers from QUESTIONS.txt
- if *-re* is invoked, answers will be included as many times as they appear in QUESTIONS
- otherwise, each answer appears only once.

**sort** *-sort* (string)

- iterates though each entry in the .csv file indicated by the flag (*-sort=$FILENAME*)
- prompts user to select a new location for each entry
- user can select "new" to create a new .csv file to place the entry into
- user can skip an entry, which will not move it
- user can exit at any time with or without saving
- leaves a copy of the entry in the original file

**overlap** *-o* && **prune** *-prune* (bool, bool)

- scans all files in the directory and identifies all entries which appear in multiple files
- if *-prune* is invoked, prompts user to delete the entry from one of the files
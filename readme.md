# APLUS
### a program to help me study for qualification exams and brush up on go

"APLUS" is a straightforward go program to quiz oneself in the command line.

At minimum, the program requires a `QUESTIONS.TXT` file to exist in the `_location_/USING` directory.

QUESTIONS should be formatted as follows:

`_question_text_ : _answer_text_\\n`

Where a colon serves as delimiter between question and answer, and newline is delimiter between q:a entries.

It asks questions from the file in random[^1] order, and checks the given answer for correctness.

The program works best with a few `_TOPIC_.csv` files. These files contain plausible answer options for a given question topic.

Use of the `_TOPIC_.csv` files allows the program to ask multiple-choice questions, by randomly[^2] populating each option from a relevant list. All `_TOPIC_` files are also stored in /USING.

Also included is a separate [utility program] for managing the .csv files. For more information, see [util_README]

Makes use of [go_print_utils](github.com/borisbugaev/go_print_utils), an extremely minimal package for managing terminal outputs. 

At some point I will likely use this as a framework for a more sophisticated, possibly GUI, application.

[^1]: no, it's not semantically random

[^2]: "iTunes Random" is how I'd define this if pushed

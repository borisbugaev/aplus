>> a program to help me study for qualification exams
>>> an opportunity to learn go

"APLUS" is a straightforward go program to quiz oneself in the command line. It asks questions from a user-provided file in quasi-random order, and checks the given answer for correctness. When possible, questions are presented in multiple-choice form. The program generates multiple-choice options at (quasi)random. To do this, questions are sorted into types, and the answers to all questions of a given type are stored in a .csv file. This allows the program to suggest plausible, but incorrect, answers for multiple-choice questions. Question type .csv are stored in the "USING" directory, and the types to be used by the program are defined in USE.csv in the program's root directory.


Also included is a separate utility program for managing the .csv files.


At some point I will likely use this as a framework for a more sophisticated, possibly GUI, application.
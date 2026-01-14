package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)


func main() {
	if length := len(os.Args); length != 2 {
		log.Fatalf("Expected 1 argument, got %d\n", length - 1)
	}

	fileContent, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	lines := strings.Split(string(fileContent), "\n")
	// discrad last line --> is empty
	lines = lines[:len(lines) - 1]
	problemCount := len(strings.Fields(lines[0]))

	// operations
	ops := strings.Fields(lines[len(lines) - 1])
	problems := make([]int64, problemCount)

	// loop through all but the last line.
	for lineNumber, line := range lines[:len(lines) - 1] {
		lineProblems := strings.Fields(line)
		if length := len(lineProblems); length != problemCount {
			log.Fatalf("Lines don't contain the same number of problems. Expected %d, got %d on line %d", problemCount, length, lineNumber + 1)
		}
		for index, problem := range lineProblems {
			parsedNumber, err := strconv.ParseInt(problem, 10, 64)
			if err != nil {
				log.Fatalf("Error parsing number: %v", err)
			}
			if lineNumber == 0 {
				problems[index] = parsedNumber
				continue
			}

			switch ops[index] {
				case "+":
					problems[index] += parsedNumber
				case "*":
					problems[index] *= parsedNumber
				default:
					log.Fatalf("Unexpected Operation in Problem %d", index + 1)
			}
		}
	}

	output := 0
	for _, problem := range problems {
		output += int(problem)
	}

	fmt.Printf("Total is %d\n", output)
}

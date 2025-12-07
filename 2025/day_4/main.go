package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type horizontalLine struct {
	lines [3]string
}

func (h *horizontalLine) setLines(lines [3]string) error {
	if len(lines[0]) != len(lines[1]) || len(lines[0]) != len(lines[2]) {
		return fmt.Errorf("Lines don't have the same length.")
	}

	h.lines = lines

	return nil
}

func (h *horizontalLine) insertNewLine(newLine string) error {
	// shift all lines up by 1 and add the new one to the end.
	if len(newLine) != len(h.lines[0]) {
		return fmt.Errorf("newLine does not match the other lines' length.")
	}

	h.lines[0] = h.lines[1]
	h.lines[1] = h.lines[2]
	h.lines[2] = newLine

	return nil
}

func (h horizontalLine) getNeighbors(index int) ([8]rune, error) {
	/*
		neighbor numbering:
		0 1 2
		3 * 4
		5 6 7
	*/

	output := [8]rune{}

	if index > len(h.lines[0])-1 || index < 0 {
		return output, fmt.Errorf("Index is out of range")
	}

	if index == 0 {
		output[0] = '.'
		output[3] = '.'
		output[5] = '.'
	} else {
		output[0] = rune(h.lines[0][index-1])
		output[3] = rune(h.lines[1][index-1])
		output[5] = rune(h.lines[2][index-1])
	}

	output[1] = rune(h.lines[0][index])
	output[6] = rune(h.lines[2][index])

	if index == len(h.lines[0])-1 {
		output[2] = '.'
		output[4] = '.'
		output[7] = '.'
	} else {
		output[2] = rune(h.lines[0][index+1])
		output[4] = rune(h.lines[1][index+1])
		output[7] = rune(h.lines[2][index+1])
	}

	return output, nil
}

func (h horizontalLine) getNumberOfNeighboringSymbols(index int, symbol rune) (int, error) {
	neighbors, err := h.getNeighbors(index)
	if err != nil {
		return 0, fmt.Errorf("Error getting neighbors: %w", err)
	}

	log.Printf("Neighbors of roll %d:\n%c %c %c \n%c * %c \n%c %c %c \n", index, neighbors[0], neighbors[1], neighbors[2], neighbors[3], neighbors[4], neighbors[5], neighbors[6], neighbors[7])

	total := 0
	for _, n := range neighbors {
		if n == symbol {
			total++
		}
	}

	log.Printf("Found %d occurences of %c\n", total, symbol)

	return total, nil
}

func (h horizontalLine) getReachableRolls() (int, error) {
	log.Printf("Getting number of reachable rolls from lines:\n%s\n%s\n%s", h.lines[0], h.lines[1], h.lines[2])
	total := 0
	for index := range h.lines[1] {
		if h.lines[1][index] != '@' {
			continue
		}
		neighboringRolls, err := h.getNumberOfNeighboringSymbols(index, '@')
		if err != nil {
			return 0, fmt.Errorf("Error getting number of neighboring rolls: %w", err)
		}
		if neighboringRolls < 4 {
			log.Printf("Adding roll %d to the total\n", index)
			total++
		}
	}

	log.Printf("Number of reachable rolls in this line: %d\n", total)

	return total, nil
}

func getReachableRolls(s *bufio.Scanner) (int, error) {
	lines := [3]string{}
	for i := range 2 {
		if !s.Scan() {
			return 0, fmt.Errorf("Error getting line %d from scanner: %w", i, s.Err())
		}
		lines[i + 1] = s.Text()
	}

	emptyLine := strings.Repeat(".", len(lines[1]))
	lines[0] = emptyLine

	l := horizontalLine{}
	err := l.setLines(lines)
	if err != nil {
		return 0, fmt.Errorf("Error setting lines: %w", err)
	}

	total, err := l.getReachableRolls()
	if err != nil {
		return  0, fmt.Errorf("Error getting number of reachable rolls from first line: %w", err)
	}

	for s.Scan() {
		if err := s.Err(); err != nil {
			return 0, fmt.Errorf("Error scanning next line: %w", err)
		}

		newLine := s.Text()
		l.insertNewLine(newLine)

		lineTotal, err := l.getReachableRolls()
		if err != nil {
			return  0, fmt.Errorf("Error getting number of reachable rolls from first line: %w", err)
		}

		total += lineTotal
	}

	l.insertNewLine(emptyLine)

	lineTotal, err := l.getReachableRolls()
	if err != nil {
		return  0, fmt.Errorf("Error getting number of reachable rolls from last line: %w", err)
	}

	total += lineTotal

	return total, nil
}

func test() {
	input := `..@@.@@@@.
@@@.@.@.@@
@@@@@.@.@@
@.@@@@..@.
@@.@@@@.@@
.@@@@@@@.@
.@.@.@.@@@
@.@@@.@@@@
.@@@@@@@@.
@.@.@@@.@.`

	r := strings.NewReader(input)
	s := bufio.NewScanner(r)

	total, err := getReachableRolls(s)
	if err != nil {
		log.Fatalf("Error getting reachable rolls: %v", err)
	}

	fmt.Printf("Total reachable rolls: %d\n", total)
}

func testPaperGrid() {
	input := `..@@.@@@@.
@@@.@.@.@@
@@@@@.@.@@
@.@@@@..@.
@@.@@@@.@@
.@@@@@@@.@
.@.@.@.@@@
@.@@@.@@@@
.@@@@@@@@.
@.@.@@@.@.`

	r := strings.NewReader(input)
	s := bufio.NewScanner(r)

	p, err := NewGrid(s, '.', '@')
	if err != nil {
		log.Fatalf("Error creating paperGrid: %v", err)
	}

	fmt.Printf("Total reachable rolls: %d\n", p.getNumberOfReachableRolls())

	fmt.Printf("%s\n", p)

	rollsRemoved, err := p.removeAllReachableRolls()
	if err != nil {
		log.Fatalf("Error removing rolls: %v", err)
	}

	fmt.Printf("%s\n", p)

	fmt.Printf("Total rolls removed: %d\n", rollsRemoved)
}

func main() {
	if length := len(os.Args); length != 2 {
		log.Fatalf("Expected 1 argument, got %d\n", length - 1)
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	s := bufio.NewScanner(file)
	p, err := NewGrid(s, '.', '@')
	if err != nil {
		log.Fatalf("Error creating paperGrid: %v", err)
	}

	total := p.getNumberOfReachableRolls()
	fmt.Printf("Total reachable rolls: %d\n", total)

	total, err = p.removeAllReachableRolls()
	if err != nil {
		log.Fatalf("Error removing rolls: %v", err)
	}

	fmt.Printf("Total rolls removed: %d\n", total)
}

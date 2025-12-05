package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type Lock struct {
	currentPosition  int
	timesReachedZero int
}

func (l *Lock) rotate(ticks int) {
	log.Printf("Starting Rotation from %d\n", l.currentPosition)
	log.Printf("Rotating %d\n", ticks)

	l.currentPosition += ticks
	l.currentPosition %= 100

	if l.currentPosition < 0 {
		l.currentPosition += 100
	}

	log.Printf("New position is %d\n", l.currentPosition)

	if l.currentPosition == 0 {
		l.timesReachedZero++
	}
}

func parseTicks(line string) (int, error) {
	ticks, err := strconv.ParseInt(line[1:], 10, 0)
	if err != nil {
		return 0, fmt.Errorf("Error parsing number of ticks: %w", err)
	}

	switch line[0] {
	case 'R':
		return int(ticks), nil
	case 'L':
		return -int(ticks), nil
	}

	return 0, fmt.Errorf("Error parsing instruction. Expected 'L' or 'R', got %c", line[0])
}

func (l *Lock) rotateFromInstruction(line string) error {
	ticks, err := parseTicks(line)
	if err != nil {
		return err
	}

	l.rotate(ticks)
	return nil
}

func passwordFromInstructionFile(pathToInstructions string) error {
	file, err := os.Open(pathToInstructions)
	if err != nil {
		return fmt.Errorf("Error opening file: %w", err)
	}
	defer file.Close()

	l := Lock{currentPosition: 50}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("Error scanning file: %w", err)
		}

		line := scanner.Text()
		l.rotateFromInstruction(line)
	}

	fmt.Printf("The password is: %d\n", l.timesReachedZero)
	return nil
}

func testInput() {
	l := Lock{currentPosition: 50}

	l.rotateFromInstruction("L68")
	l.rotateFromInstruction("L30")
	l.rotateFromInstruction("R48")
	l.rotateFromInstruction("L5")
	l.rotateFromInstruction("R60")
	l.rotateFromInstruction("L55")
	l.rotateFromInstruction("L1")
	l.rotateFromInstruction("L99")
	l.rotateFromInstruction("R14")
	l.rotateFromInstruction("L82")

	fmt.Printf("Password is: %d\n", l.timesReachedZero)
}

func main() {
	log.SetOutput(io.Discard)

	if length := len(os.Args); length != 2 {
		log.Fatalf("Expected 1 argument, got %d", length)
	}

	err := passwordFromInstructionFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}

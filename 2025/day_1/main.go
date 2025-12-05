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

func (l *Lock) rotateRight(ticks int) {
	log.Printf("[Start] Rotating right from %d + %d\n", l.currentPosition, ticks)

	l.currentPosition += ticks

	if newZeros := l.currentPosition / 100; newZeros > 0 {
		log.Printf("Adding %d 0s\n", newZeros)
		l.timesReachedZero += newZeros
	}

	l.currentPosition %= 100

	log.Printf("[End] New Position is %d\n", l.currentPosition)
}

func (l *Lock) rotateLeft(ticks int) {
	log.Printf("[Start] Rotating left from %d - %d\n", l.currentPosition, ticks)

	if l.currentPosition == 0 {
		l.currentPosition = 100
	}

	l.currentPosition -= ticks

	if newZeros := l.currentPosition / -100; newZeros > 0 {
		log.Printf("Adding %d 0s\n", newZeros)
		l.timesReachedZero += newZeros
	}

	l.currentPosition %= 100
	if l.currentPosition < 0 {
		log.Printf("Adding another 0 to the counter, because currentPosition is still < 0\n")
		l.timesReachedZero++
		l.currentPosition += 100
	}
	if l.currentPosition == 0 {
		log.Printf("New position is 0! Adding to counter.\n")
		l.timesReachedZero++
	}

	log.Printf("[End] New Position is %d\n", l.currentPosition)
}

func (l *Lock) rotateFromInstruction(line string) error {
	ticks, err := strconv.ParseInt(line[1:], 10, 0)
	if err != nil {
		return fmt.Errorf("Error parsing number of ticks: %w", err)
	}

	switch line[0] {
	case 'R':
		l.rotateRight(int(ticks))
		return nil
	case 'L':
		l.rotateLeft(int(ticks))
		return nil
	}

	return fmt.Errorf("Error parsing instruction. Expected 'L' or 'R', got %c", line[0])
}

func passwordFromInstructionFile(pathToInstructions string) error {
	log.SetOutput(io.Discard)

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

func testBigRotation() {
	l := Lock{currentPosition: 50}
	l.rotateFromInstruction("R1000")
	l.rotateFromInstruction("R550")
	fmt.Printf("password would be: %d\n", l.timesReachedZero)

	l = Lock{currentPosition: 50}
	l.rotateFromInstruction("L1000")
	l.rotateFromInstruction("L550")
	fmt.Printf("password would be: %d\n", l.timesReachedZero)
}

func main() {
	if length := len(os.Args); length != 2 {
		log.Fatalf("Expected 1 argument, got %d", length)
	}

	err := passwordFromInstructionFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}

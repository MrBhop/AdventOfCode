package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func getJoltage(batteries []byte) (int, error) {
	frontPart := batteries[:len(batteries)-1]
	log.Printf("len(batteris) = %d, len(frontPart) = %d\n", len(batteries), len(frontPart))

	maxIndexFront, err := getMaxIndex(frontPart)
	if err != nil {
		return 0, fmt.Errorf("Error getting max index of front part: %w", err)
	}

	backPart := batteries[maxIndexFront+1:]
	log.Printf("maxIndexFront = %d, len(backPart) = %d\n", maxIndexFront, len(backPart))
	maxIndexBack, err := getMaxIndex(backPart)
	if err != nil {
		return 0, fmt.Errorf("Error getting max index of back part: %w", err)
	}

	joltageString := fmt.Sprintf("%c%c", frontPart[maxIndexFront], backPart[maxIndexBack])
	log.Printf("joltageString: %s\n", joltageString)
	joltage, err := strconv.ParseInt(joltageString, 10, 0)
	return int(joltage), err
}

func getMaxIndex(bytes []byte) (int, error) {
	log.Printf("Getting max index of batteries: %s\n", bytes)
	maxIndex := 0
	for index, b := range bytes {
		if b < '0' || b > '9' {
			return 0, fmt.Errorf("non numeric byte encountered")
		}
		if b == '9' {
			return index, nil
		}
		if b > bytes[maxIndex] {
			maxIndex = index
		}
	}

	return maxIndex, nil
}

func testJoltage() {
	input := `987654321111111
811111111111119
234234234234278
818181911112111`

	banks := strings.Split(input, "\n")
	log.Printf("banks: %d, %+v\n", len(banks), banks)
	totalJoltage := 0
	for _, bank := range banks {
		log.Printf("bank: %s\n", bank)
		joltage, err := getJoltage([]byte(bank))
		if err != nil {
			log.Fatalf("Error getting joltage of batteries: %v", err)
		}

		totalJoltage += joltage
	}

	log.Printf("Total joltage: %d\n", totalJoltage)
}

func totalJoltageFromFile(inputFilePath string) error {
	log.SetOutput(io.Discard)

	file, err := os.Open(inputFilePath)
	if err != nil {
		return fmt.Errorf("Error opening file: %w", err)
	}
	defer file.Close()

	totalJoltage := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("Error scanning file: %w", err)
		}

		line := scanner.Text()
		joltage, err := getJoltage([]byte(line))
		if err != nil {
			return fmt.Errorf("Error getting joltage from line: %w", err)
		}

		totalJoltage += joltage
	}

	fmt.Printf("Total joltage: %d\n", totalJoltage)

	return nil
}

func main() {
	if length := len(os.Args); length != 2 {
		log.Fatalf("Expected 1 argumetn, got %d", length-1)
	}

	err := totalJoltageFromFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}

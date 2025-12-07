package main

import (
	"bufio"
	"fmt"
	"log"
)

type paperGrid struct {
	grid [][]byte
	emptySymbol byte
	occupiedSymbol byte
}

func (p paperGrid) String() string {
	output := []byte{}

	for _, line := range p.grid {
		output = append(output, line...)
		output = append(output, '\n')
	}

	return fmt.Sprintf("%s", output)
}

func (p paperGrid) height() int {
	return len(p.grid)
}

func (p paperGrid) width() int {
	if p.height() == 0 {
		return 0
	}

	return len(p.grid[0])
}

func (p paperGrid) getSymbol(x, y int) (symbol byte, isInbounds bool) {
	if min(x, y) < 0 || x > p.width() - 1 || y > p.height() - 1 {
		return p.emptySymbol, false
	}

	return p.grid[y][x], true
}

func (p paperGrid) occurencesInNeighbors(x, y int, symbol byte) int {
	offsets := [3]int{-1, 0, 1}

	total := 0
	for _, ox := range offsets {
		for _, oy := range offsets {
			if ox == 0 && oy == 0 {
				// 0 offest is not a neighbor.
				continue
			}

			neighbor, _ := p.getSymbol(x + ox, y + oy)
			if neighbor == symbol {
				total++
			}
		}
	}

	return total
}

func (p paperGrid) rollIsReachable(x, y int) bool {
	result :=  p.occurencesInNeighbors(x, y, p.occupiedSymbol) < 4
	return result

}

func (p paperGrid) getNumberOfReachableRolls() int {
	total := 0
	for y, line := range p.grid {
		for x := range line {
			if symbol, _ := p.getSymbol(x, y); symbol != p.occupiedSymbol {
				continue
			}
			if p.rollIsReachable(x, y) {
				total++
			}
		}
	}

	return total
}

func (p *paperGrid) updateSymbol(x, y int, newSymbol byte) error {
	if _, isInbounds := p.getSymbol(x, y); !isInbounds {
		return fmt.Errorf("%d, %d is not in bounds.", x, y)
	}

	p.grid[y][x] = newSymbol
	return nil
}

func (p *paperGrid) removeReachableRollsOnce() (int, error) {
	total := 0
	for y, line := range p.grid {
		for x := range line {
			log.Printf("Currently on spot %d, %d\n", x, y)
			if symbol, _ := p.getSymbol(x, y); symbol != p.occupiedSymbol {
				log.Printf("Spot is not a roll. Continuing.\n")
				continue
			}
			log.Printf("Spot is a roll.\n")
			if p.rollIsReachable(x, y) {
				log.Printf("Roll is reachable.\n")
				total++

				err := p.updateSymbol(x, y, p.emptySymbol)
				if err != nil {
					return 0, fmt.Errorf("Error updating symbol: %w", err)
				}
			}
		}
	}

	log.Printf("Removed %d rolls\n", total)
	return total, nil
}

func (p *paperGrid) removeAllReachableRolls() (int, error) {
	total := 0

	for {
		log.Printf("\n%s", p)
		removedRolls, err := p.removeReachableRollsOnce()
		if err != nil {
			return 0, fmt.Errorf("Error removing rolls: %w", err)
		}
		if removedRolls == 0 {
			break
		}

		total += removedRolls
	}

	return total, nil
}

func NewGrid(s *bufio.Scanner, emptySymbol, occupiedSymbol byte) (*paperGrid, error) {
	grid := [][]byte{}
	for s.Scan() {
		if err := s.Err(); err != nil {
			return nil, fmt.Errorf("Error while scanning: %w", err)
		}

		grid = append(grid, []byte(s.Text()))
		if len(grid[0]) != len(grid[len(grid) - 1]) {
			return nil, fmt.Errorf("Lines don't have the same length.")
		}
	}

	return &paperGrid{grid: grid, emptySymbol: emptySymbol, occupiedSymbol: occupiedSymbol}, nil
}

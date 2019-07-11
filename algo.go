package main

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten"
)

func (env *Env) botPlayer() {
	// wait for graphics
	fmt.Println("Press SPACE to start bot...")
	for {
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			break
		}
		time.Sleep(DELAY)
	}
	// start algo
	env.algo()
}

func (env *Env) algo() {
	env.buildFinished()
	env.aStar()
}

func (env *Env) aStar() {
	var closedList []*grid
	var openList []*grid
	// Append start node to open list
	env.grid.cost = 0
	env.grid.heuristic = env.globalHeuristic(env.grid)
	openList = append(openList, env.grid)
	for len(openList) != 0 {
		// Unstack first cell of open list
		currGrid := openList[0]
		env.grid = copyGrid(currGrid)
		// Check end
		if env.isFinished() {
			closedList = append(closedList, currGrid)
			fmt.Println("aStar done")
			return
		}
		//for each possible move
		movesList := env.getMoves(currGrid)
		for _, newGrid := range movesList {
			openList = append(openList, newGrid)
		}
		//append currGrid to closedList
		closedList = append(closedList, currGrid)
		// pop currGrid from openList
		openList = openList[1:]
		// sort openList
		sort.Slice(openList, func(i, j int) bool {
			return openList[i].heuristic < openList[j].heuristic
		})
	}
	fmt.Println("aStar returned no solution")
}

func isPresentID(currGrid *grid, gridList []*grid) int {
	for i, grid := range gridList {
		for cell := 0; cell < len(currGrid.mapping); cell++ {
			if grid.mapping[cell] == currGrid.mapping[cell] {
				return i
			}
		}
	}
	return -1
}

func (env *Env) getMoves(currGrid *grid) []*grid {
	var gridList []*grid
	for direction := 1; direction < 5; direction++ {
		i := env.checkMove(currGrid, direction)
		if i >= 0 {
			newGrid := env.virtualMove(currGrid, direction, i)
			gridList = append(gridList, newGrid)
		}
	}
	return gridList
}

func (env *Env) virtualMove(currGrid *grid, direction int, i int) *grid {
	newGrid := copyGrid(currGrid)
	if i >= 0 {
		switch {
		case direction == UP:
			newGrid.mapping[0].Y++
			newGrid.mapping[i].Y--
		case direction == DOWN:
			newGrid.mapping[0].Y--
			newGrid.mapping[i].Y++
		case direction == LEFT:
			newGrid.mapping[0].X++
			newGrid.mapping[i].X--
		case direction == RIGHT:
			newGrid.mapping[0].X--
			newGrid.mapping[i].X++
		}
		newGrid.cost = newGrid.cost + 1
		newGrid.heuristic = newGrid.cost + env.globalHeuristic(newGrid)
	}
	return newGrid
}

func copyGrid(srcGrid *grid) *grid {
	newGrid := &grid{}
	newGrid.mapping = make([]*cell, len(srcGrid.mapping))
	for i := 0; i < len(srcGrid.mapping); i++ {
		newCell := cell{}
		newCell.X = srcGrid.mapping[i].X
		newCell.Y = srcGrid.mapping[i].Y
		newCell.cellImg = srcGrid.mapping[i].cellImg
		newGrid.mapping[i] = &newCell
	}
	newGrid.cost = srcGrid.cost
	newGrid.heuristic = srcGrid.heuristic
	return newGrid
}

func (env *Env) globalHeuristic(currGrid *grid) int {
	gManDist := 0
	for id := 0; id < len(currGrid.mapping); id++ {
		switch {
		case env.heuristic == "" || env.heuristic == "md":
			gManDist += manhattanDistance(currGrid.mapping[id], env.finishedMap.mapping[id])
		case env.heuristic == "c":
			gManDist += countLeft(currGrid.mapping[id], env.finishedMap.mapping[id])
		}
	}
	return gManDist
}

// Heuristics :

func manhattanDistance(a, b *cell) int {
	return int(math.Abs(float64(a.X)-float64(b.X)) +
		math.Abs(float64(a.Y)-float64(b.Y)))
}

func countLeft(a, b *cell) int {
	if a.X != b.X || a.Y != b.Y {
		return 1
	}
	return 0
}

package main

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten"
)

func (env *Env) botPlayer() {
	env.buildFinished()
	// Wait for graphics
	if env.graph {
		fmt.Println("Press SPACE to start bot...")
		for {
			if ebiten.IsKeyPressed(ebiten.KeySpace) {
				break
			}
			time.Sleep(DELAY)
		}
	}
	fmt.Println("Start...")
	// Start algo
	env.aStar()
}

func (env *Env) aStar() {
	var closedList []*Grid
	var openList []*Grid
	// Append start node to open list
	openList = append(openList, env.grid)
	env.grid.cost = 0
	env.grid.heuristic = env.globalHeuristic(env.grid)
	lastGrid := env.grid
	for len(openList) != 0 {
		// Unstack first cell of open list
		currGrid := openList[0]
		// Clear openList
		openList = openList[:0]
		// Update state
		env.grid = currGrid
		// Check end
		if env.isFinished() {
			closedList = append(closedList, currGrid)
			fmt.Println("Astar done in ", len(closedList)-1, "turns")
			return
		}
		// For each possible move
		movesList := env.getMoves(currGrid)
		for _, newGrid := range movesList {
			if !equal(newGrid.mapping, lastGrid.mapping) {
				// Append newGrid to openList
				openList = append(openList, newGrid)
			}
		}
		// Append currGrid to closedList
		closedList = append(closedList, currGrid)
		lastGrid = currGrid
		// Sort openList
		sort.Slice(openList, func(i, j int) bool {
			return openList[i].heuristic < openList[j].heuristic
		})
	}
	fmt.Println("Astar returned no solution")
}

func (env *Env) getMoves(currGrid *Grid) []*Grid {
	var gridList []*Grid
	for _, direction := range env.seed.Perm(4) {
		direction++
		i := env.checkMove(currGrid, direction)
		if i >= 0 {
			newGrid := env.virtualMove(currGrid, direction, i)
			gridList = append(gridList, newGrid)
		}
	}
	return gridList
}

func (env *Env) virtualMove(currGrid *Grid, direction int, i int) *Grid {
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

func (env *Env) globalHeuristic(currGrid *Grid) int {
	gHeur := 0

	for id := 1; id < len(currGrid.mapping); id++ {
		switch {
		case env.heuristic == "" || env.heuristic == "md":
			gHeur += manhattanDistance(currGrid.mapping[id], env.finishedMap.mapping[id])
		case env.heuristic == "hd":
			gHeur += hammingDistance(currGrid.mapping[id], env.finishedMap.mapping[id])
		case env.heuristic == "ed":
			gHeur += euclidianDistance(currGrid.mapping[id], env.finishedMap.mapping[id])
		case env.heuristic == "lc":
			gHeur += env.linearConflicts(currGrid, currGrid.mapping[id], env.finishedMap.mapping[id], id)
		}

	}
	return gHeur
}

// Heuristics :

func manhattanDistance(a, b *cell) int {
	return int(math.Abs(float64(a.X)-float64(b.X)) +
		math.Abs(float64(a.Y)-float64(b.Y)))
}

func euclidianDistance(a, b *cell) int {
	return int(math.Sqrt(math.Pow(float64(a.X)-float64(b.X), 2) +
		math.Pow(float64(a.Y)-float64(b.Y), 2)))
}

func hammingDistance(a, b *cell) int {
	if a.X != b.X || a.Y != b.Y {
		return 1
	}
	return 0
}

func (env *Env) linearConflicts(currGrid *Grid, a, b *cell, id int) int {
	md := manhattanDistance(a, b)
	lc := 0
	if a.X == b.X && a.Y != b.Y {
		for i := 0; i < env.size; i++ {
			if i != a.Y {
				idx := idxByXY(currGrid, a.X, i)
				if idx != 0 {
					if currGrid.mapping[idx].X == env.finishedMap.mapping[idx].X &&
						currGrid.mapping[idx].Y != env.finishedMap.mapping[idx].Y {
						lc++
					}
				}
			}
		}
	} else if a.X != b.X && a.Y == b.Y {
		for i := 0; i < env.size; i++ {
			if i != a.X {
				idx := idxByXY(currGrid, i, a.Y)
				if idx != 0 {
					if currGrid.mapping[idx].X != env.finishedMap.mapping[idx].X &&
						currGrid.mapping[idx].Y == env.finishedMap.mapping[idx].Y {
						lc++
					}
				}
			}
		}
	}
	return (lc / 2) + md
}

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
	// wait for graphics
	fmt.Println("Press SPACE to start bot...")
	for {
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			break
		}
		time.Sleep(DELAY)
	}
	// start algo
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
		//for each possible move
		movesList := env.getMoves(currGrid)
		for _, newGrid := range movesList {
			if !equal(newGrid.mapping, lastGrid.mapping) {
				openList = append(openList, newGrid)
			}
		}
		//append currGrid to closedList
		closedList = append(closedList, currGrid)
		lastGrid = currGrid
		// sort openList
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

	for id := 0; id < len(currGrid.mapping); id++ {
		switch {
		case env.heuristic == "" || env.heuristic == "md":
			gHeur += manhattanDistance(currGrid.mapping[id], env.finishedMap.mapping[id], id)
		case env.heuristic == "hd":
			gHeur += hammingDistance(currGrid.mapping[id], env.finishedMap.mapping[id], id)
		case env.heuristic == "lc":
			gHeur += linearConflicts(currGrid, currGrid.mapping[id], env.finishedMap.mapping[id], id)
		}

	}
	return gHeur
}

// Heuristics :

func manhattanDistance(a, b *cell, id int) int {
	if id == 0 {
		return 0
	}
	return int(math.Abs(float64(a.X)-float64(b.X)) +
		math.Abs(float64(a.Y)-float64(b.Y)))
}

func hammingDistance(a, b *cell, id int) int {
	if id != 0 && (a.X != b.X || a.Y != b.Y) {
		return 1
	}
	return 0
}

func linearConflicts(currGrid *Grid, a, b *cell, id int) int {
	if id == 0 {
		return 0
	}
	md := manhattanDistance(a, b, id)
	lc := 0
	if (a.X == b.X && a.Y != b.Y) || (a.X != b.X && a.Y == b.Y) {
		lc++
	}
	return (2 * lc) + md
}

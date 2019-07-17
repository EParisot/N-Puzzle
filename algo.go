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
	//env.aStar()
	env.idAstar()
}

func (env *Env) reconstructPath(closedList []*Grid, endGrid *Grid) {
	var finalList []*Grid
	var parentID int

	finalList = append(finalList, endGrid)
	parentID = endGrid.parentID
	for i := 0; i < len(closedList); i++ {
		if closedList[i].id == parentID {
			finalList = append(finalList, closedList[i])
			parentID = closedList[i].parentID
			if parentID == 0 {
				//End
				//Print solution in reverse order
				fmt.Println("Ordered sequence of states that make up the solution : ")
				for j := len(finalList) - 1; j != -1; j-- {
					env.printGrid(finalList[j])
				}
				fmt.Println("Number of moves required : ", len(finalList)-1)
				fmt.Println("Complexity in size : ", len(closedList))
				return
			}
			i = -1
		}
	}

}

/////// IDASTAR TEST ///////

func (env *Env) idAstar() {
	threshold := env.globalHeuristic(env.grid)
	env.grid.id = env.getID()
	env.grid.parentID = 0
	var closedList []*Grid
	closedList = append(closedList, env.grid)

	for {
		tmpThres, closedList := env.search(threshold, closedList)
		if tmpThres == -1 {
			fmt.Println("IDAstar Done")
			// test
			for _, grid := range closedList {
				env.grid = grid
			}
			// end test
			return
		} else if tmpThres >= 100000 {
			fmt.Println("IDAstar returned no solution")
			return
		}
		threshold = tmpThres
	}
}

func (env *Env) search(threshold int, closedList []*Grid) (int, []*Grid) {
	currGrid := closedList[len(closedList)-1]
	if currGrid.heuristic > threshold {
		return currGrid.heuristic, closedList
	}
	if env.isFinished(currGrid) {
		return -1, closedList
	}
	min := 100000
	childsList := env.getMoves(currGrid)
	for _, child := range childsList {
		if !existInClosedList(child, closedList) {
			closedList = append(closedList, child)
			tmp, closedList := env.search(threshold, closedList)
			if tmp == -1 {
				return -1, closedList
			}
			if tmp < min {
				min = tmp
			}
			closedList = closedList[:len(closedList)-1]
		}
	}
	return min, closedList
}

/////// IDASTAR TEST ///////

func (env *Env) aStar() {
	var closedList []*Grid
	var openList []*Grid
	var currGrid *Grid

	// Append start node to open list
	openList = append(openList, env.grid)
	env.grid.cost = 0
	env.grid.heuristic = env.globalHeuristic(env.grid)
	env.grid.id = env.getID()
	env.grid.parentID = 0

	for len(openList) != 0 {
		// Unstack first cell of open list
		currGrid, openList = openList[0], openList[1:]
		// Update state
		env.grid = currGrid
		// Check end
		if env.isFinished(nil) {
			env.reconstructPath(closedList, currGrid)
			return
		}
		// For each possible move
		movesList := env.getMoves(currGrid)
		for _, newGrid := range movesList {
			if existInClosedList(newGrid, closedList) || existInOpenListWithInferiorCost(newGrid, openList) {
			} else {
				openList = append(openList, newGrid)
				sort.Slice(openList, func(i, j int) bool {
					return openList[i].heuristic < openList[j].heuristic
				})
			}
		}
		// Append currGrid to closedList
		closedList = append(closedList, currGrid)
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
		newGrid.id = env.getID()
		newGrid.parentID = env.grid.id
		newGrid.cost = newGrid.cost + 1
		newGrid.heuristic = newGrid.cost + env.globalHeuristic(newGrid)*5
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

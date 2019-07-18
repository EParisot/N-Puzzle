package main

import (
	"fmt"
	"math"
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
	fmt.Println("Working...")
	// Start algo
	env.idAstar()
}

func (env *Env) reconstructPathIDA(closedList []*Grid, endGrid *Grid) {
	fmt.Println("Ordered sequence of states that make up the solution : ")
	for _, step := range closedList {
		env.grid = step
		env.printGrid(step)
	}
	fmt.Println("Number of moves required : ", len(closedList))
	fmt.Println("Number of moves checked : ", env.timeComp)
	fmt.Println("Max Number of moves represented in memory : ", env.sizeComp)
}

func (env *Env) idAstar() {
	threshold := env.globalHeuristic(env.grid)
	var closedList []*Grid
	closedList = append(closedList, env.grid)
	env.timeComp++
	sizeComp := env.sizeComp
	env.startTime = time.Now()
	for {
		tmpThres, closedList, _ := env.search(threshold, &closedList, sizeComp)
		if tmpThres == -1 {
			fmt.Println("IDAstar Done")
			env.reconstructPathIDA(*closedList, (*closedList)[len(*closedList)-1])
			return
		} else if tmpThres >= 10000 {
			fmt.Println("IDAstar returned no solution")
			return
		}
		threshold = tmpThres
	}
}

func (env *Env) search(threshold int, closedList *[]*Grid, sizeComp int) (int, *[]*Grid, int) {
	if time.Since(env.startTime) >= 10000000000 {
		var closedList []*Grid
		closedList = append(closedList, env.grid)
		fmt.Println("Incrementing W to ", env.w+1)
		env.w++
		env.startTime = time.Now()
		env.timeComp = 0
		env.sizeComp = 0
		return env.globalHeuristic(env.grid), &closedList, sizeComp
	}
	sizeComp++
	if sizeComp > env.sizeComp {
		env.sizeComp = sizeComp
	}
	currGrid := (*closedList)[len(*closedList)-1]
	if currGrid.heuristic > threshold {
		return currGrid.heuristic, closedList, sizeComp
	}
	if env.isFinished(currGrid) {
		return -1, closedList, sizeComp
	}
	min := 100000
	childsList := env.getMoves(currGrid)
	for _, child := range childsList {
		if !existInClosedList(child, *closedList) {
			*closedList = append(*closedList, child)
			env.timeComp++
			tmp, closedList, sizeComp := env.search(threshold, closedList, sizeComp)
			if tmp == -1 {
				sizeComp--
				return -1, closedList, sizeComp
			}
			if tmp < min {
				min = tmp
			}
			if len(*closedList) > 1 {
				*closedList = (*closedList)[:len(*closedList)-1]
			} else {
				sizeComp--
				return tmp, closedList, sizeComp
			}
		}
	}
	sizeComp--
	return min, closedList, sizeComp
}

func (env *Env) getMoves(currGrid *Grid) []*Grid {
	var gridList []*Grid
	for direction := 1; direction < 5; direction++ {
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
		if env.greedySearch == false {
			newGrid.cost = newGrid.cost + 1
		}
		newGrid.heuristic = newGrid.cost + env.globalHeuristic(newGrid)*env.w
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

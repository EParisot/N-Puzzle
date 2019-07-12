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
	//env.aStar()
	env.idAstar()
}

/////// IDASTAR TEST ///////
func (env *Env) idAstar() {
	threshold := env.globalHeuristic(env.grid)
	best := env.grid
	for {
		tmpThres, best := env.search(best, threshold)
		if tmpThres == -1 {
			env.grid = best
			fmt.Println("IDAstar Done")
			return
		} else if tmpThres >= 10000 {
			fmt.Println("IDAstar returned no solution")
			return
		}
		//env.grid = best
		threshold = tmpThres
	}
}

func (env *Env) search(currGrid *Grid, threshold int) (int, *Grid) {
	if currGrid.heuristic > threshold {
		return currGrid.heuristic, currGrid
	}
	if env.isFinished(currGrid) {
		return -1, currGrid
	}
	min := 10000
	childsList := env.getMoves(currGrid)
	var best *Grid
	for _, child := range childsList {
		tmp, child := env.search(child, threshold)
		if tmp == -1 {
			return -1, child
		}
		if tmp <= min {
			min = tmp
			best = child
		}
	}
	return min, best
}

/////// IDASTAR TEST ///////

func (env *Env) aStar() {
	var closedList []*Grid
	var openList []*Grid
	var antiBoucle []*Grid
	var i int
	i = 0
	// Append start node to open list
	env.grid.cost = 0
	env.grid.heuristic = env.globalHeuristic(env.grid)
	openList = append(openList, env.grid)
	for len(openList) != 0 {
		// Unstack first cell of open list
		currGrid := openList[0]
		env.grid = CopyGrid(currGrid)
		// Check end
		if env.isFinished(nil) {
			closedList = append(closedList, currGrid)
			fmt.Println(len(closedList))
			return
		}
		//for each possible move
		movesList := env.getMoves(currGrid)
		var j int
		j = 0
		for _, newGrid := range movesList {
			if !havedouble(newGrid, antiBoucle) {
				openList = append(openList, newGrid)
				antiBoucle = append(antiBoucle, newGrid)
				j++
			}
		}
		if j != 0 {
			//append currGrid to closedList
			closedList = append(closedList, currGrid)
		}
		// pop currGrid from openList
		openList = openList[1:]
		// sort openList
		sort.Slice(openList, func(i, j int) bool {
			return openList[i].heuristic < openList[j].heuristic
		})
		i++
	}
	fmt.Println("aStar returned no solution")
}

func equal(a, b []*cell) bool {
	for i := 0; i < len(a); i++ {
		if a[i].X != b[i].X || a[i].Y != b[i].Y {
			return false
		}
	}
	return true
}

func havedouble(gridToCheck *Grid, openList []*Grid) bool {
	for i := 0; i < len(openList); i++ {
		if equal(gridToCheck.mapping, openList[i].mapping) {
			return true
		}
	}
	return false
}

func isPresentID(currGrid *Grid, gridList []*Grid) int {
	for i, grid := range gridList {
		for cell := 0; cell < len(currGrid.mapping); cell++ {
			if grid.mapping[cell] == currGrid.mapping[cell] {
				return i
			}
		}
	}
	return -1
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
	newGrid := CopyGrid(currGrid)
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

// CopyGrid copy a grid
func CopyGrid(srcGrid *Grid) *Grid {
	newGrid := &Grid{}
	newGrid.mapping = make([]*cell, len(srcGrid.mapping))
	for i := 0; i < len(srcGrid.mapping); i++ {
		newCell := cell{}
		newCell.X = srcGrid.mapping[i].X
		newCell.Y = srcGrid.mapping[i].Y
		newCell.cellImg = srcGrid.mapping[i].cellImg
		newCell.digitImg = srcGrid.mapping[i].digitImg
		newGrid.mapping[i] = &newCell
	}
	newGrid.cost = srcGrid.cost
	newGrid.heuristic = srcGrid.heuristic
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
	/*if (((a.X == b.X+1 && IdxByXY(currGrid, a.X+1, a.Y) != 0) ||
	(a.X == b.X-1 && IdxByXY(currGrid, a.X-1, a.Y) != 0)) && a.Y == b.Y) ||
	(((a.Y == b.Y+1 && IdxByXY(currGrid, a.X, a.Y+1) != 0) ||
		(a.Y == b.Y-1 && IdxByXY(currGrid, a.X, a.Y-1) != 0)) && a.X == b.X) {*/
	if (a.X == b.X && a.Y != b.Y) || (a.X != b.X && a.Y == b.Y) {
		lc++
	}
	return (2 * lc) + md
}

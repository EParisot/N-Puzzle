package main

import (
	"fmt"
)

func (env *Env) checkMove(currGrid *Grid, move int) int {
	if move == UP && currGrid.mapping[0].Y != env.size-1 {
		for i := range currGrid.mapping {
			if currGrid.mapping[i].X == currGrid.mapping[0].X && currGrid.mapping[i].Y == currGrid.mapping[0].Y+1 {
				return i
			}
		}
	} else if move == DOWN && currGrid.mapping[0].Y != 0 {
		for i := range currGrid.mapping {
			if currGrid.mapping[i].X == currGrid.mapping[0].X && currGrid.mapping[i].Y == currGrid.mapping[0].Y-1 {
				return i
			}
		}
	} else if move == LEFT && currGrid.mapping[0].X != env.size-1 {
		for i := range currGrid.mapping {
			if currGrid.mapping[i].X == currGrid.mapping[0].X+1 && currGrid.mapping[i].Y == currGrid.mapping[0].Y {
				return i
			}
		}
	} else if move == RIGHT && currGrid.mapping[0].X != 0 {
		for i := range currGrid.mapping {
			if currGrid.mapping[i].X == currGrid.mapping[0].X-1 && currGrid.mapping[i].Y == currGrid.mapping[0].Y {
				return i
			}
		}
	}
	return -1
}

func (env *Env) moveCell(currGrid *Grid, direction int) {
	i := env.checkMove(currGrid, direction)
	if i >= 0 {
		switch {
		case direction == UP:
			currGrid.mapping[0].Y++
			currGrid.mapping[i].Y--
		case direction == DOWN:
			currGrid.mapping[0].Y--
			currGrid.mapping[i].Y++
		case direction == LEFT:
			currGrid.mapping[0].X++
			currGrid.mapping[i].X--
		case direction == RIGHT:
			currGrid.mapping[0].X--
			currGrid.mapping[i].X++
		}
	}
}

func (env *Env) buildFinished() {
	x := 0
	y := 0
	countSide := 0
	countCell := 0
	offset := 0
	way := 0
	finished := make([]*cell, len(env.grid.mapping))
	for id := 1; id < len(env.grid.mapping); id++ {
		finished[id] = &cell{}
		finished[id].X = x
		finished[id].Y = y
		if countCell+offset == env.size-1 {
			countCell = 0
			if countSide%2 == 0 {
				offset++
			}
			countSide++
			if way == 3 {
				way = 0
			} else {
				way++
			}
		} else {
			countCell++
		}

		switch {
		case way == 0:
			x++
		case way == 1:
			y++
		case way == 2:
			x--
		case way == 3:
			y--
		}
	}
	finished[0] = &cell{}
	finished[0].X = x
	finished[0].Y = y
	env.finishedMap.mapping = finished
}

func (env *Env) isFinished() bool {
	if len(env.finishedMap.mapping) == 0 {
		env.buildFinished()
	}
	for id := 1; id < env.size*env.size; id++ {
		if env.grid.mapping[id].X != env.finishedMap.mapping[id].X || env.grid.mapping[id].Y != env.finishedMap.mapping[id].Y {
			return false
		}
	}
	return true
}

func (env *Env) printGrid(grid *Grid) {
	for y := 0; y < env.size; y++ {
		for x := 0; x < env.size; x++ {
			for i := 0; i < len(grid.mapping); i++ {
				if grid.mapping[i].X == x && grid.mapping[i].Y == y {
					fmt.Print(i)
					break
				}
			}
			fmt.Print(" ")
		}
		fmt.Print("\n")
	}
	fmt.Println()
}

func equal(a, b []*cell) bool {
	for i := 0; i < len(a); i++ {
		if a[i].X != b[i].X || a[i].Y != b[i].Y {
			return false
		}
	}
	return true
}

func existInClosedList(newGrid *Grid, closedList []*Grid) bool {
	for i := 0; i < len(closedList); i++ {
		if equal(closedList[i].mapping, newGrid.mapping) {
			return true
		}
	}
	return false
}

func existInOpenListWithInferiorCost(newGrid *Grid, openList []*Grid) bool {
	for i := 0; i < len(openList); i++ {
		if equal(newGrid.mapping, openList[i].mapping) {
			if newGrid.cost < openList[i].cost {
				return true
			}
		}
	}
	return false
}

func (env *Env) isPresent(idToTest int) bool {
	for id, cell := range env.grid.mapping {
		if id == idToTest && cell != nil {
			return true
		}
	}
	return false
}

func copyGrid(srcGrid *Grid) *Grid {
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

func idxByVAL(list []int, val int) int {
	i := 0
	for i = range list {
		if list[i] == val {
			break
		}
	}
	return i
}

func idxByXY(grid *Grid, x, y int) int {
	i := 0
	for i = range grid.mapping {
		if grid.mapping[i].X == x && grid.mapping[i].Y == y {
			break
		}
	}
	return i
}

func (env *Env) getID() int {
	env.incrementID++
	return env.incrementID
}

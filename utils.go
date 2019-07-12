package main

import (
	"fmt"
)

func (env *Env) isPresent(idToTest int) bool {
	for id, cell := range env.grid.mapping {
		if id == idToTest && cell != nil {
			return true
		}
	}
	return false
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

// if grid is nil, eval the env.grid
func (env *Env) isFinished(grid *Grid) bool {
	if grid == nil {
		grid = env.grid
	}
	if len(env.finishedMap.mapping) == 0 {
		env.buildFinished()
	}
	for id := 1; id < env.size*env.size; id++ {
		if grid.mapping[id].X != env.finishedMap.mapping[id].X ||
			grid.mapping[id].Y != env.finishedMap.mapping[id].Y {
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

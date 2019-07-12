package main

import (
	"fmt"
	"math/rand"
)

func (env *Env) buildMap() {
	newMap := make([]*cell, env.size*env.size)
	env.grid.mapping = newMap
	env.buildFinished()
	fmt.Println("Building map...")
	for env.checkSolvability() == false {
		env.grid = CopyGrid(env.finishedMap)
		for i := 0; i < 10000; i++ {
			env.shuffle(env.seed)
		}
	}
}

func (env *Env) shuffle(r *rand.Rand) {
	move := r.Intn(4) + 1
	env.moveCell(env.grid, move)
}

func (env *Env) checkSolvability() bool {
	if env.grid.mapping[0] == nil {
		return false
	}
	inversions := env.countInversions()
	// if size is ODD
	if env.size%2 != 0 {
		//inversions must be EVEN
		if inversions%2 != 0 {
			return false
		}
	} else {
		// if blank row is ODD
		if env.grid.mapping[0].Y%2 != 0 {
			// inversions shoul be EVEN
			if inversions%2 != 0 {
				return false
			}
			// if blank row is EVEN
		} else {
			// inversions shoulb be ODD
			if inversions%2 == 0 {
				return false
			}
		}
	}
	return true
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

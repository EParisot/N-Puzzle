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
	for env.checkSolvability(env.grid) == false {
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

func (env *Env) checkSolvability(grid *Grid) bool {
	if grid.mapping[0] == nil {
		return false
	}
	inversions := env.countInversions()
	// if size is ODD
	if env.size%2 != 0 {
		//inversions must be EVEN
		if inversions%2 != 0 {
			return false
		}
		// if size is EVEN
	} else {
		// if dist of blank row from final blank row (idx start at 1) is EVEN
		if (env.finishedMap.mapping[0].Y-grid.mapping[0].Y)%2 == 0 {
			// inversions should be EVEN
			if inversions%2 != 0 {
				return false
			}
			// if blank row is EVEN
		} else {
			// inversions should be ODD
			if inversions%2 == 0 {
				return false
			}
		}
	}
	return true
}

func (env *Env) countInversions() int {
	var currList []int
	var finishedList []int
	finishedMap := env.finishedMap
	for y := 0; y < env.size; y++ {
		for x := 0; x < env.size; x++ {
			currList = append(currList, IdxByXY(env.grid, x, y))
			finishedList = append(finishedList, IdxByXY(finishedMap, x, y))
		}
	}
	// iter on ids to count inversions
	inversions := 0
	for pivot := range currList {
		if currList[pivot] != 0 {
			// find pivot in result
			k := idxByVAL(finishedList, currList[pivot])
			// for each next id in curr
			for i := range currList[pivot+1:] {
				if currList[pivot+1+i] != 0 {
					// check if next val in curr < pos pivot in res
					j := idxByVAL(finishedList, currList[pivot+1+i])
					if j < k {
						inversions++
					}
				}
			}
		}
	}
	return inversions
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

func IdxByXY(grid *Grid, x, y int) int {
	i := 0
	for i = range grid.mapping {
		if grid.mapping[i].X == x && grid.mapping[i].Y == y {
			break
		}
	}
	return i
}

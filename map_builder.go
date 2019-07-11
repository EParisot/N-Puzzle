package main

import (
	"fmt"
	"math/rand"
)

func (env *Env) buildMap() {
	switch {
	case env.difficulty == "E" || env.difficulty == "Easy":
		env.size = 3
	case env.difficulty == "M" || env.difficulty == "Medium":
		env.size = 8
	case env.difficulty == "H" || env.difficulty == "Hard":
		env.size = 16
	}
	newMap := make([]*cell, env.size*env.size)
	env.grid.mapping = newMap
	r := rand.New(rand.NewSource(42)) //time.Now().UnixNano()))
	newID := r.Intn(env.size * env.size)
	fmt.Println("Building map...")
	for env.checkSolvability() == false {
		for y := 0; y < env.size; y++ {
			for x := 0; x < env.size; x++ {
				for env.isPresent(newID) {
					newID = r.Intn(env.size * env.size)
				}
				env.grid.mapping[newID] = &cell{
					X: x,
					Y: y,
				}
			}
		}
	}
}

func (env *Env) checkSolvability() bool {
	if env.grid.mapping[0] == nil {
		return false
	}
	var currList []int
	var finishedList []int
	env.buildFinished()
	finishedMap := env.finishedMap
	inversions := 0
	for y := 0; y < env.size; y++ {
		for x := 0; x < env.size; x++ {
			currList = append(currList, idxByXY(env.grid, x, y))
			finishedList = append(finishedList, idxByXY(finishedMap, x, y))
		}
	}
	// iter on ids
	for pivot := range currList {
		if currList[pivot] != 0 {
			// find pivot in result
			k := idxByVAL(finishedList, currList[pivot])
			// for each next id in curr
			for i := range currList[pivot+1:] {
				// check if next val in curr < pos pivot in res
				j := idxByVAL(finishedList, currList[pivot+i])
				if j < k {
					inversions++
				}
			}
		}
	}
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

func idxByXY(grid *grid, x, y int) int {
	i := 0
	for i = range grid.mapping {
		if grid.mapping[i].X == x && grid.mapping[i].Y == y {
			break
		}
	}
	return i
}

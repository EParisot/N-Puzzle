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
	//print(env.globalManDist())
}

func (env *Env) aStar() {
	var closedList []*cell
	var openList []*cell
	// Append start node to open list
	env.grid[0].cost = 0
	env.grid[0].heuristic = env.globalManDist()
	openList = append(openList, env.grid[0])
	for len(openList) != 0 {
		// Sort open list by heuristic
		sort.Slice(openList, func(i, j int) bool {
			return openList[i].heuristic < openList[j].heuristic
		})
		// Unstack first cell of open list
		currPos := openList[0]
		openList[0] = nil
		openList = openList[1:]
		// Check end
		if env.isFinished() {
			closedList = append(closedList, currPos)
			return
		}
		//for each possible move
		//	if already present in closedList
		//		then continue
		//	elif already present in openList with lower cost
		//		then continue
		//
		//	append move to openList
		//append move to closedList
	}
}

func manhattanDistance(a, b *cell) int {
	return int(math.Abs(float64(a.X)-float64(b.X)) +
		math.Abs(float64(a.Y)-float64(b.Y)))
}

func (env *Env) globalManDist() int {
	gManDist := 0
	for id := 0; id < len(env.grid); id++ {
		gManDist += manhattanDistance(env.grid[id], env.finishedMap[id])
	}
	return gManDist
}

package main

import (
	"fmt"
	"math/rand"
	"time"
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
	env.grid = make([]*cell, env.size*env.size)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	newID := r.Intn(env.size * env.size)
	fmt.Println("Building map...")
	for y := 0; y < env.size; y++ {
		for x := 0; x < env.size; x++ {
			for env.isPresent(newID) {
				newID = r.Intn(env.size * env.size)
			}
			env.grid[newID] = &cell{
				X: x,
				Y: y,
			}
		}
	}
}

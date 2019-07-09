package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten"
)

// Env main game struct
type Env struct {
	mapFile    string
	imgFile    string
	difficulty string
	grid       []*cell
	size       int
	autoMode   bool
	heuristic  string
}

type cell struct {
	X   int
	Y   int
	img *ebiten.Image
}

func main() {
	env := Env{autoMode: false}
	err := env.parseFile()
	if err != nil {
		log.Fatal(err)
	}
	//DEBUG
	fmt.Println("size : ", env.size)
	for i := range env.grid {
		fmt.Println(i, env.grid[i])
	}
	env.cropImage(env.imgFile)
	go env.getKey()
	if err := ebiten.Run(env.update, 300, 300, 2, "N-Puzzle"); err != nil {
		log.Fatal(err)
	}
	fmt.Println(env.isFinished())
	//TODO start Algo
}

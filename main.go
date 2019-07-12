package main

import (
	"image"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
)

// Env main game struct
type Env struct {
	mapFile     string
	imgFile     string
	grid        *Grid
	finishedMap *Grid
	size        int
	sizeWindows int
	heuristic   string
	digit       bool
	seed        *rand.Rand
}

// Grid hold the map
type Grid struct {
	mapping   []*cell
	cost      int
	heuristic int
}

type cell struct {
	X        int
	Y        int
	cellImg  image.Image
	digitImg image.Image
}

func main() {
	env := Env{
		grid:        &Grid{},
		finishedMap: &Grid{},
		heuristic:   "md",
		seed:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	err := env.parseFile()
	if err != nil {
		log.Fatal(err)
	}
	//Find the perfect size for the windows
	if env.size == 0 {
		log.Fatal("error missing size value")
	}
	env.sizeWindows = 300 + (300 % env.size)
	// Default images
	if env.imgFile == "" {
		env.imgFile = "images/default.png"
	}
	env.cropImage(env.imgFile)
	//start Algo
	go env.botPlayer()
	go env.getKey()
	if err := ebiten.Run(env.update, env.sizeWindows, env.sizeWindows, 2, "N-Puzzle"); err != nil {
		log.Fatal(err)
	}
}

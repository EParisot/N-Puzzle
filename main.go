package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten"
)

// Env main game struct
type Env struct {
	mapFile     string
	imgFile     string
	difficulty  string
	grid        *grid
	finishedMap *grid
	size        int
	sizeWindows int
	autoMode    bool
	heuristic   string
	digit       bool
}

type grid struct {
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
		grid:        &grid{},
		finishedMap: &grid{},
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

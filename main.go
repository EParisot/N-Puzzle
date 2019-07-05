package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten"
)

// Env main game struct
type Env struct {
	grid []*cell
	size int
	mode string
}

type cell struct {
	id  int
	X   int
	Y   int
	img *ebiten.Image
}

func main() {
	env := Env{}
	err := env.parseFile()
	if err != nil {
		log.Fatal(err)
	}
	//DEBUG
	fmt.Println("size : ", env.size)
	for i := range env.grid {
		fmt.Println(env.grid[i])
	}
}

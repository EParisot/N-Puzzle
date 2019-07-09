package main

import (
	"fmt"
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

}

func (env *Env) astar() {

}

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

const (
	UP    = 1
	DOWN  = 2
	LEFT  = 3
	RIGHT = 4

	DELAY = time.Second / 4
)

var square *ebiten.Image

func (env *Env) update(screen *ebiten.Image) error {

	//Fill the screen with background color
	screen.Fill(color.NRGBA{0xff, 0x00, 0x00, 0xff})
	//Handle controls
	if !env.isFinished() {
		env.getKey()
	}
	for i := range env.grid.mapping {
		//Add cells
		env.addSquare(float64(env.grid.mapping[i].X*(env.sizeWindows/env.size)),
			float64(env.grid.mapping[i].Y*(env.sizeWindows/env.size)),
			square,
			screen,
			i,
		)
	}

	return nil
}

func (env *Env) getKey() {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		env.moveCell(env.grid, UP)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		env.moveCell(env.grid, DOWN)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		env.moveCell(env.grid, LEFT)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		env.moveCell(env.grid, RIGHT)
	}
}

func (env *Env) checkMove(currGrid *grid, move int) int {
	if move == UP && currGrid.mapping[0].Y != env.size-1 {
		for i := range currGrid.mapping {
			if currGrid.mapping[i].X == currGrid.mapping[0].X && currGrid.mapping[i].Y == currGrid.mapping[0].Y+1 {
				return i
			}
		}
	} else if move == DOWN && currGrid.mapping[0].Y != 0 {
		for i := range currGrid.mapping {
			if currGrid.mapping[i].X == currGrid.mapping[0].X && currGrid.mapping[i].Y == currGrid.mapping[0].Y-1 {
				return i
			}
		}
	} else if move == LEFT && currGrid.mapping[0].X != env.size-1 {
		for i := range currGrid.mapping {
			if currGrid.mapping[i].X == currGrid.mapping[0].X+1 && currGrid.mapping[i].Y == currGrid.mapping[0].Y {
				return i
			}
		}
	} else if move == RIGHT && currGrid.mapping[0].X != 0 {
		for i := range currGrid.mapping {
			if currGrid.mapping[i].X == currGrid.mapping[0].X-1 && currGrid.mapping[i].Y == currGrid.mapping[0].Y {
				return i
			}
		}
	}
	return -1
}

func (env *Env) moveCell(currGrid *grid, direction int) {
	i := env.checkMove(currGrid, direction)
	if i >= 0 {
		switch {
		case direction == UP:
			currGrid.mapping[0].Y++
			currGrid.mapping[i].Y--
		case direction == DOWN:
			currGrid.mapping[0].Y--
			currGrid.mapping[i].Y++
		case direction == LEFT:
			currGrid.mapping[0].X++
			currGrid.mapping[i].X--
		case direction == RIGHT:
			currGrid.mapping[0].X--
			currGrid.mapping[i].X++
		}
	}
}

func (env *Env) addSquare(x float64, y float64, square *ebiten.Image, screen *ebiten.Image, i int) {

	var err error

	if i != 0 {
		square, err = ebiten.NewImageFromImage(env.grid.mapping[i].cellImg, ebiten.FilterDefault)
		if err != nil {
			log.Fatal("Error new images", err)
		}
	} else {
		square, _ = ebiten.NewImage((env.sizeWindows / env.size), (env.sizeWindows / env.size), ebiten.FilterNearest)
		square.Fill(color.Black)
	}

	opts := &ebiten.DrawImageOptions{}

	// Add the Translate effect to the option struct.
	opts.GeoM.Translate(x, y)
	screen.DrawImage(square, opts)

}

func (env *Env) cropImage(images string) {

	f, err := os.Open(images)
	if err != nil {
		log.Fatal("Cannot open file", err)
	}
	// Accept for now only png
	img, err := png.Decode(f)
	if err != nil {
		log.Fatal("Cannot decode image:", err)
	}
	//Resize the picture to sizeWindows
	newImage := resize.Resize(uint(env.sizeWindows), uint(env.sizeWindows), img, resize.Lanczos3)

	position_x := 0
	position_y := 0
	countSide := 0
	countCell := 0
	offset := 0
	way := 0
	for i := range env.grid.mapping {
		if i == 0 {
			continue
		}
		// Crop the image to multiple square
		cImg, err := cutter.Crop(newImage, cutter.Config{
			Height:  (env.sizeWindows / env.size),                                                          // height in pixel or Y ratio(see Ratio Option below)
			Width:   (env.sizeWindows / env.size),                                                          // width in pixel or X ratio
			Mode:    cutter.TopLeft,                                                                        // Accepted Mode: TopLeft, Centered
			Anchor:  image.Point{position_x * (env.sizeWindows / env.size), position_y * (300 / env.size)}, // Position of the top left point
			Options: 0,                                                                                     // Accepted Option: Ratio
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Each cell fill with a square of the image
		env.grid.mapping[i].cellImg = cImg

		if countCell+offset == env.size-1 {
			countCell = 0
			if countSide%2 == 0 {
				offset++
			}
			countSide++
			if way == 3 {
				way = 0
			} else {
				way++
			}
		} else {
			countCell++
		}

		switch {
		case way == 0:
			position_x++
		case way == 1:
			position_y++
		case way == 2:
			position_x--
		case way == 3:
			position_y--
		}
	}
	f.Close()
}

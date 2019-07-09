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

	env.getKey()

	for i := range env.grid {
		if i == 0 {
			continue
		}
		//Add cells
		env.addSquare(float64(env.grid[i].X*(300/env.size)),
			float64(env.grid[i].Y*(300/env.size)),
			square,
			screen,
			i,
		)
	}

	return nil
}

func (env *Env) getKey() {
	//for {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		env.moveCell(UP)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		env.moveCell(DOWN)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		env.moveCell(LEFT)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		env.moveCell(RIGHT)
	}
	if env.isFinished() {
		//TODO GAME OVER
		return
	}
	//		time.Sleep(DELAY)
}

//}

func (env *Env) moveCell(direction int) {

	if direction == UP {
		if env.grid[0].Y == env.size-1 {
			fmt.Println("You can't do that")
		} else {
			for i := range env.grid {
				if env.grid[i].X == env.grid[0].X && env.grid[i].Y == env.grid[0].Y+1 {
					env.grid[0].Y++
					env.grid[i].Y--
					break
				}
			}
		}
	} else if direction == DOWN {
		if env.grid[0].Y == 0 {
			fmt.Println("You can't do that")
		} else {
			for i := range env.grid {
				if env.grid[i].X == env.grid[0].X && env.grid[i].Y == env.grid[0].Y-1 {
					env.grid[0].Y--
					env.grid[i].Y++
					break
				}
			}
		}
	} else if direction == LEFT {
		if env.grid[0].X == env.size-1 {
			fmt.Println("You can't do that")
		} else {
			for i := range env.grid {
				if env.grid[i].X == env.grid[0].X+1 && env.grid[i].Y == env.grid[0].Y {
					env.grid[0].X++
					env.grid[i].X--
					break
				}
			}
		}
	} else if direction == RIGHT {
		if env.grid[0].X == 0 {
			fmt.Println("You can't do that")
		} else {
			for i := range env.grid {
				if env.grid[i].X == env.grid[0].X-1 && env.grid[i].Y == env.grid[0].Y {
					env.grid[0].X--
					env.grid[i].X++
					break
				}
			}
		}
	}
}

func (env *Env) addSquare(x float64, y float64, square *ebiten.Image, screen *ebiten.Image, i int) {

	var err error

	square, err = ebiten.NewImageFromImage(env.grid[i].cellImg, ebiten.FilterDefault)
	if err != nil {
		log.Fatal("Error new images", err)
	}

	opts := &ebiten.DrawImageOptions{}

	// Add the Translate effect to the option struct.
	opts.GeoM.Translate(x, y)
	screen.DrawImage(square, opts)

}

func (env *Env) cropImage(images string) {

	fmt.Println(images)
	f, err := os.Open(images)
	if err != nil {
		log.Fatal("Cannot open file", err)
	}
	// Accept for now only png
	img, err := png.Decode(f)
	if err != nil {
		log.Fatal("Cannot decode image:", err)
	}

	//Resize the picture to 300 * 300
	newImage := resize.Resize(300, 300, img, resize.Lanczos3)

	position_x := 0
	position_y := 0
	countSide := 0
	countCell := 0
	offset := 0
	way := 0
	for i := range env.grid {
		if i == 0 {
			continue
		}
		// Crop the image to multiple square
		cImg, err := cutter.Crop(newImage, cutter.Config{
			Height:  (300 / env.size),                                                          // height in pixel or Y ratio(see Ratio Option below)
			Width:   (300 / env.size),                                                          // width in pixel or X ratio
			Mode:    cutter.TopLeft,                                                            // Accepted Mode: TopLeft, Centered
			Anchor:  image.Point{position_x * (300 / env.size), position_y * (300 / env.size)}, // Position of the top left point
			Options: 0,                                                                         // Accepted Option: Ratio
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Each cell fill with a square of the image
		env.grid[i].cellImg = cImg

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

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strconv"
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

	if env.digit {
		var err error

		if i != 0 {
			square, err = ebiten.NewImageFromImage(env.grid.mapping[i].digitImg, ebiten.FilterDefault)
			if err != nil {
				log.Fatal("Error new images", err)
			}
		}
		opts := &ebiten.DrawImageOptions{}

		// Add the Translate effect to the option struct.
		opts.GeoM.Translate(x, y)
		screen.DrawImage(square, opts)

	}
}

func (env *Env) imgDigit(digit int) image.Image {
	f, err := os.Open("digits/" + strconv.Itoa(digit) + ".png")
	if err != nil {
		log.Fatal("Cannot open Digit file", err)
	}
	// Accept for now only png
	img, err := png.Decode(f)
	if err != nil {
		log.Fatal("Cannot decode image:", err)
	}
	return img
}

func (env *Env) mergeTwoImages(img1, img2 image.Image) image.Image {
	//starting position of the second image (bottom left)
	sp2 := image.Point{img1.Bounds().Dx(), 0}
	//new rectangle for the second image
	r2 := image.Rectangle{sp2, sp2.Add(img2.Bounds().Size())}
	//rectangle for the big image
	r := image.Rectangle{image.Point{0, 0}, r2.Max}
	rgba := image.NewRGBA(r)
	draw.Draw(rgba, img1.Bounds(), img1, image.Point{0, 0}, draw.Src)
	draw.Draw(rgba, r2, img2, image.Point{0, 0}, draw.Src)
	return rgba
}

func (env *Env) getDigit(digit int) image.Image {

	var rgba image.Image

	if digit/100 > 0 {
		rgba = env.mergeTwoImages(env.imgDigit(digit/100), env.imgDigit((digit%100)/10))
		rgba = env.mergeTwoImages(rgba, env.imgDigit(digit%10))
	} else if digit/10 > 0 {
		rgba = env.mergeTwoImages(env.imgDigit(digit/10), env.imgDigit(digit%10))
	} else {
		rgba = env.imgDigit(digit)
	}

	newImage := resize.Resize(uint((env.sizeWindows/env.size)/3), uint((env.sizeWindows/env.size)/3), rgba, resize.Lanczos3)
	return newImage
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

		if env.digit {
			env.grid.mapping[i].digitImg = env.getDigit(i)
		}
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

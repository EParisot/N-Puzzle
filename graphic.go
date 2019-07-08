package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

var square *ebiten.Image

func (env *Env) update(screen *ebiten.Image) error {

	// Fill the screen with #FF0000 color
	screen.Fill(color.NRGBA{0xff, 0x00, 0x00, 0xff})

	// Display the text though the debug function
	ebitenutil.DebugPrint(screen, "Our first game in Ebiten!")

	for i := range env.grid {
		if i == 0 {
			continue
		}
		env.Addsquare(float64(env.grid[i].X*(300/env.size)),
			float64(env.grid[i].Y*(300/env.size)),
			square,
			screen,
			i,
		)
	}
	// Draw the square image to the screen with an empty option

	return nil
}

func (env *Env) Addsquare(x float64, y float64, square *ebiten.Image, screen *ebiten.Image, i int) {

	var err error

	if square == nil {
		square, _, err = ebitenutil.NewImageFromFile(".tmp/"+strconv.Itoa(i)+".png", ebiten.FilterDefault)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Fill the square with the white color
	//square.Fill(color.RGBA{0xff, uint8(colors), 0, 0xff})

	// The previous empty option struct
	opts := &ebiten.DrawImageOptions{}

	// Add the Translate effect to the option struct.
	opts.GeoM.Translate(x, y)
	screen.DrawImage(square, opts)
	// Draw the square image to the screen with an empty option

}

func (env *Env) CropImage(images string) {

	f, err := os.Open(images)
	if err != nil {
		log.Fatal("Cannot open file", err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal("Cannot decode image:", err)
	}

	//Size the image

	newImage := resize.Resize(300, 300, img, resize.Lanczos3)

	//fmt.Println(newImage)

	//Clean the .tmp directory
	err = RemoveContents(".tmp")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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
		fmt.Println("position_x = ", position_x)
		fmt.Println("position_y = ", position_y)
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

		out, err := os.Create(".tmp/" + strconv.Itoa(i) + ".png")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = png.Encode(out, cImg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		out.Close()
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

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (env *Env) parseFile() error {
	// Parse Args
	err := env.parseArgs()
	if err != nil {
		return err
	}
	if env.mapFile != "" {
		// Read file
		file, err := os.Open(env.mapFile)
		if err != nil {
			return (err)
		}
		defer file.Close()
		reader := bufio.NewReader(file)
		// Read size
		err = env.readSize(reader)
		if err != nil {
			return err
		}
		// Read map
		err = env.readMap(reader)
		if err != nil {
			return err
		}
	}
	return nil
}

func printUsage() {
	fmt.Println(`Usage : N-Puzzle map_file [-m map] [-i image] [-d difficulty] [-a heuristic]
			-m map        = 'map_file.map'
			-i image      = 'image_file.png'
			-s size       = map size (int)
			-h heuristic  = 'heuristic' ('md' (default), 'hd', 'i')
			-dg (Add numbers to the picture)
			`)
}

func (env *Env) parseArgs() error {
	if len(os.Args[1:]) < 1 {
		printUsage()
		return errors.New("")
	}
	for i, arg := range os.Args {
		if arg == "-m" && i+1 < len(os.Args) &&
			strings.HasSuffix(os.Args[i+1], ".map") {
			env.mapFile = os.Args[i+1]
		} else if arg == "-i" && i+1 < len(os.Args) &&
			strings.HasSuffix(os.Args[i+1], ".png") {
			env.imgFile = os.Args[i+1]
		} else if arg == "-s" && i+1 < len(os.Args) {
			size, err := strconv.Atoi(os.Args[i+1])
			if err != nil || size < 2 || size > 31 {
				return errors.New("error invalid size value")
			}
			env.size = size
		} else if arg == "-h" && i+1 < len(os.Args) {
			if os.Args[i+1] != "md" &&
				os.Args[i+1] != "hd" &&
				os.Args[i+1] != "ed" &&
				os.Args[i+1] != "lc" {
				return errors.New("error invalid heuristic value")
			}
			env.heuristic = os.Args[i+1]
		} else if arg == "-dg" {
			env.digit = true
		} else if arg == "-g" {
			env.graph = true
		}
	}
	if env.mapFile == "" {
		if env.size >= 2 {
			env.buildMap()
		}
	}
	return nil
}

func (env *Env) readSize(reader *bufio.Reader) error {
	firstLine, err := parseLine(reader)
	if len(strings.Fields(firstLine)) == 0 ||
		len(strings.Fields(firstLine)) > 1 {
		return errors.New("error missing size value")
	}
	size, err := strconv.Atoi(firstLine)
	if err != nil || size < 3 || size > 31 {
		return errors.New("error invalid size value")
	}
	env.size = size
	return nil
}

func (env *Env) readMap(reader *bufio.Reader) error {
	env.grid.mapping = make([]*cell, env.size*env.size)
	for j := 0; j < env.size; j++ {
		line, err := parseLine(reader)
		if err != nil {
			return err
		}
		ids := strings.Fields(line)
		if len(ids) != env.size {
			return errors.New("error invalid map size")
		}
		for i, val := range ids {
			valInt, err := strconv.Atoi(val)
			if err != nil || valInt >= env.size*env.size || valInt < 0 {
				return errors.New("error invalid cell id")
			}
			if env.isPresent(valInt) {
				return errors.New("error duplicated cell id")
			}
			env.grid.mapping[valInt] = &cell{
				X: i,
				Y: j,
			}
		}
	}
	if env.isFinished() == false && env.checkSolvability(env.grid) == false {
		return errors.New("error unsolvable map")
	}
	return nil
}

func parseLine(reader *bufio.Reader) (string, error) {
	var lineStr string
	for len(lineStr) == 0 {
		line, _, err := reader.ReadLine()
		if err != nil {
			return "", errors.New("error reading map file")
		}
		lineStr = string(line)
		lineTab := strings.Split(lineStr, "#")
		lineStr = lineTab[0]
	}
	return lineStr, nil
}

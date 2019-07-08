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
		with -m map	   	   = 'map_file.map'
		     -i image      = 'image_file.png'
			 -d difficulty = 'E[asy]', 'M[edium]', 'H[ard]'
			 -a heuristic  = 'heuristic' (default 'manhattan distance')`)
}

func (env *Env) parseArgs() error {
	if len(os.Args[1:]) < 1 {
		printUsage()
	}
	for i, arg := range os.Args {
		if arg == "-m" && i+1 < len(os.Args) &&
			strings.HasSuffix(os.Args[i+1], ".map") {
			env.mapFile = os.Args[i+1]
		} else if arg == "-i" && i+1 < len(os.Args) &&
			strings.HasSuffix(os.Args[i+1], ".png") {
			env.imgFile = os.Args[i+1]
		} else if arg == "-d" && i+1 < len(os.Args) {
			env.difficulty = os.Args[i+1]
		} else if arg == "-a" && i+1 < len(os.Args) {
			env.autoMode = true
			env.heuristic = os.Args[i+1]
		}
	}
	if env.mapFile == "" {
		if env.difficulty != "" {
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
	if err != nil || size < 3 {
		return errors.New("error invalid size value")
	}
	env.size = size
	return nil
}

func (env *Env) readMap(reader *bufio.Reader) error {
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
			if err != nil || valInt > env.size*env.size || valInt < 0 {
				return errors.New("error invalid cell id")
			}
			if env.isPresent(valInt) {
				return errors.New("error duplicated cell id")
			}
			env.grid = append(env.grid, &cell{
				id: valInt,
				X:  i,
				Y:  j,
			})
		}
	}
	return nil
}

func parseLine(reader *bufio.Reader) (string, error) {
	line, _, err := reader.ReadLine()
	if err != nil {
		return "", errors.New("error reading map file")
	}
	lineStr := string(line)
	lineTab := strings.Split(lineStr, "#")
	return lineTab[0], nil
}

func (env *Env) isPresent(idToTest int) bool {
	for _, cellID := range env.grid {
		if cellID.id == idToTest {
			return true
		}
	}
	return false
}

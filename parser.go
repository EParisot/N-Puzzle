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
	// Read file
	file, err := os.Open(os.Args[1])
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
	return nil
}

func (env *Env) parseArgs() error {
	if len(os.Args[1:]) < 1 {
		fmt.Println(`Usage : N-Puzzle map_file [-m heuristic]
		with heuristic = 'md' for Manhattan Distance, WIP...`)
		return errors.New("missing argument")
	}
	for i, arg := range os.Args {
		if arg == "-m" && i+1 < len(os.Args) {
			env.mode = os.Args[i+1]
		}
	}
	return nil
}

func (env *Env) readSize(reader *bufio.Reader) error {
	firstLine, err := parseLine(reader)
	if len(strings.Fields(firstLine)) == 0 ||
		len(strings.Fields(firstLine)) > 1 {
		return errors.New("missing size value")
	}
	size, err := strconv.Atoi(firstLine)
	if err != nil || size < 3 {
		return errors.New("invalid size value")
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
			return errors.New("invalid map size")
		}
		for i, val := range ids {
			valInt, err := strconv.Atoi(val)
			if err != nil || valInt > env.size*env.size || valInt < 0 {
				return errors.New("invalid cell id")
			}
			if env.isPresent(valInt) {
				return errors.New("duplicated cell id")
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
		return "", errors.New("Error Reading map file")
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

package main

func (env *Env) isPresent(idToTest int) bool {
	for id, cell := range env.grid {
		if id == idToTest && cell != nil {
			return true
		}
	}
	return false
}

func (env *Env) isFinished() bool {
	x := 0
	y := 0
	countSide := 0
	countCell := 0
	offset := 0
	way := 0
	for id := 1; id < (env.size * env.size); id++ {
		if env.grid[id].X != x || env.grid[id].Y != y {
			return false
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
			x++
		case way == 1:
			y++
		case way == 2:
			x--
		case way == 3:
			y--
		}
	}

	return true
}

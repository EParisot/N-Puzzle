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

	return true
}

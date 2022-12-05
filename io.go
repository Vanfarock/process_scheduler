package main

type IO struct {
	id        int
	cycles    Cycles
	maxCycles Cycles
}

func NewIO(id int, maxCycles Cycles) *IO {
	return &IO{
		id:        id,
		cycles:    0,
		maxCycles: maxCycles,
	}
}

func (io IO) finished() bool {
	return io.cycles >= io.maxCycles
}

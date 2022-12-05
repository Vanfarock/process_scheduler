package main

type ProcessState int

const (
	Ready ProcessState = iota + 1
	Running
	Waiting
)

type Process struct {
	id           int
	timeSlice    Cycles
	cpuCycles    Cycles
	maxCpuCycles Cycles
	io           *IO
	ioChance     float32
	ioCycles     Cycles
	skips        Cycles
	history      []ProcessState
}

func NewProcess(id int, timeSlice Cycles, maxCpuCycles Cycles, ioChance float32) *Process {
	return &Process{
		id:           id,
		timeSlice:    timeSlice,
		cpuCycles:    0,
		maxCpuCycles: maxCpuCycles,
		io:           nil,
		ioChance:     ioChance,
		ioCycles:     0,
		skips:        0,
		history:      make([]ProcessState, 0),
	}
}

func (p Process) getTotalCycles() Cycles {
	return p.cpuCycles + p.ioCycles
}

func (p Process) getIoBound() float32 {
	totalCycles := p.getTotalCycles()
	if totalCycles > 0 {
		return float32(p.ioCycles) / float32(totalCycles)
	}
	return 0
}

// Calcula skips do processo baseado na porcentagem de ciclos IO
// [ 0%,  20%) → 0 skips
// [20%,  40%) → 2 skips
// [40%,  60%) → 4 skips
// [60%,  80%) → 6 skips
// [80%, 100%) → 8 skips
func (p Process) getMaxSkips() Cycles {
	ioBound := p.getIoBound()
	if ioBound < 0.2 {
		return 0
	} else if ioBound < 0.4 {
		return 2
	} else if ioBound < 0.6 {
		return 4
	} else if ioBound < 0.8 {
		return 6
	} else {
		return 8
	}
}

func (p Process) finished() bool {
	return p.cpuCycles >= p.maxCpuCycles
}

func (p Process) finishedTimeSlice() bool {
	return p.cpuCycles%p.timeSlice == 0
}

func (p *Process) saveHistory(state ProcessState) {
	p.history = append(p.history, state)
}

func (p Process) showHistory() {
	for _, state := range p.history {
		if state == Ready {
			PrintColor("■", White)
		} else if state == Running {
			PrintColor("■", Green)
		} else if state == Waiting {
			PrintColor("■", Red)
		}
	}
}

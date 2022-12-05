package main

type Cpu struct {
	cores int
}

func NewCpu(cores int) Cpu {
	return Cpu{
		cores: cores,
	}
}

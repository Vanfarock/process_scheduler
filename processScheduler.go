package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Cycles int

type ProcessScheduler struct {
	cpu              Cpu
	totalProcesses   int
	processes        []*Process
	chReady          chan *Process
	chRequeueReady   chan *Process
	chRunning        chan *Process
	chRequeueRunning chan *Process
	chWaiting        chan *Process
	chRequeueWaiting chan *Process
	countRunning     int
	cycleDuration    time.Duration
}

func NewProcessScheduler(cpu Cpu, processes []*Process, cycleDuration time.Duration) ProcessScheduler {
	return ProcessScheduler{
		cpu:              cpu,
		totalProcesses:   len(processes),
		processes:        processes,
		chReady:          make(chan *Process, len(processes)),
		chRequeueReady:   make(chan *Process, len(processes)),
		chRunning:        make(chan *Process, cpu.cores),
		chRequeueRunning: make(chan *Process, cpu.cores),
		chWaiting:        make(chan *Process, len(processes)),
		chRequeueWaiting: make(chan *Process, len(processes)),
		countRunning:     0,
		cycleDuration:    cycleDuration,
	}
}

// Inicia o escalonamento dos processos
// Fluxo da informação:
// Início → Ready
// Ready → Run
// Run → [Ready, Waiting, Fim]
// Waiting → Ready
func (scheduler *ProcessScheduler) Run() {
	// Inicia x threads, onde x é o número de núcleos da CPU
	for i := 0; i < scheduler.cpu.cores; i++ {
		go scheduler.RunReady()
		go scheduler.RunRunning()
		go scheduler.RunWaiting()
	}

	// Inicia threads para fazer o reenfileiramento dos processos ready
	go func() {
		for process := range scheduler.chRequeueReady {
			scheduler.chReady <- process
		}
	}()

	// Inicia threads para fazer o reenfileiramento dos processos running
	go func() {
		for process := range scheduler.chRequeueRunning {
			scheduler.chRunning <- process
		}
	}()

	// Inicia threads para fazer o reenfileiramento dos processos waiting
	go func() {
		for process := range scheduler.chRequeueWaiting {
			scheduler.chWaiting <- process
		}
	}()

	// Popula channel de processos ready com todos os processos informados
	fmt.Println("Starting up...")
	for _, process := range scheduler.processes {
		scheduler.chRequeueReady <- process
	}
}

func (scheduler *ProcessScheduler) RunReady() {
	for process := range scheduler.chReady {
		time.Sleep(scheduler.cycleDuration)
		process.saveHistory(Ready)

		// Verifica se existem núcleos de CPU ociosos
		if scheduler.canRun() {
			// Se for identificado que um processo possui muitos IOs, será atribuído
			// dinamicamente um número de skips que ele deve fazer antes de ser
			// executado. Isto é, se tiverem mais processos na situação Ready,
			// ele irá dar sua vez para eles, até um certo limite (getMaxSkips).
			// Um possível problema que pode surgir, com a entrada de processos
			// dinâmicamente, é o starvation do processo (não acontecerá no nosso caso).
			if process.skips < process.getMaxSkips() {
				process.skips += 1
				scheduler.chRequeueReady <- process
				continue
			}

			// Adiciona o processo na fila de processos running
			process.skips = 0
			scheduler.countRunning += 1
			scheduler.chRequeueRunning <- process
			fmt.Println("Moving process", process.id, "from ready to running")
			continue
		}

		// Adiciona o processo na fila de processos ready
		scheduler.chRequeueReady <- process
	}
}

func (scheduler *ProcessScheduler) RunRunning() {
	for process := range scheduler.chRunning {
		time.Sleep(scheduler.cycleDuration)
		process.saveHistory(Running)

		// Adiciona um ciclo na execução do processo
		process.cpuCycles += 1
		fmt.Println("Running process", process.id, "- Cycle", process.cpuCycles, "/", process.maxCpuCycles)

		// Verifica se o processo terminou
		if process.finished() {
			fmt.Println("Process", process.id, "finished")
			scheduler.countRunning -= 1
			scheduler.totalProcesses -= 1
			continue
		}

		// Verifica se o processo já finalizou seu time slice
		if process.finishedTimeSlice() {
			// Adiciona o processo na fila de processos ready
			scheduler.chRequeueReady <- process
			scheduler.countRunning -= 1
			fmt.Println("Moving process", process.id, "from running to ready")
			continue
		}

		// Verifica se deve rodar a IO
		dice := rand.Float32()
		if dice < process.ioChance {
			// Adiciona o processo na fila de processos waiting
			scheduler.chRequeueWaiting <- process
			scheduler.countRunning -= 1
			fmt.Println("Moving process", process.id, "from running to waiting")
			continue
		}

		// Adiciona o processo na fila de processos running
		scheduler.chRequeueRunning <- process
	}
}

func (scheduler *ProcessScheduler) RunWaiting() {
	for process := range scheduler.chWaiting {
		time.Sleep(scheduler.cycleDuration)
		process.saveHistory(Waiting)

		// Se o processo não possui IO, cria um IO com número de ciclos aleatório (min: 1, máx: 10)
		if process.io == nil {
			ioMaxCycles := rand.Float32()*10 + 1
			process.io = NewIO(process.id, Cycles(ioMaxCycles))
		}

		// Adiciona um ciclo na execução do IO
		process.io.cycles += 1
		fmt.Println("Running IO", process.io.id, "- Cycle", process.io.cycles, "/", process.io.maxCycles)

		// Verifica se o IO terminou
		if process.io.finished() {
			// Adiciona ciclos de IO no processo, para que o cálculo de IO bound do processo seja feito
			process.ioCycles += process.io.cycles
			// Adiciona o processo na fila de processos ready
			scheduler.chRequeueReady <- process
			fmt.Println("IO", process.io.id, "finished")
			fmt.Println("Moving process", process.id, "from waiting to ready")
			process.io = nil
			continue
		}

		// Adiciona o processo na fila de processos waiting
		scheduler.chRequeueWaiting <- process
	}
}

func (scheduler *ProcessScheduler) canRun() bool {
	return scheduler.countRunning < scheduler.cpu.cores
}

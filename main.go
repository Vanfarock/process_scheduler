package main

import (
	"fmt"
	"time"
)

func main() {
	processes := make([]*Process, 0)

	// Cria os processos
	// Pârametros: (id, time slice, máximo de ciclos, chance de rodar IO)
	process1 := NewProcess(1, 5, 15, 0.07)
	process2 := NewProcess(2, 3, 10, 0.2)
	process3 := NewProcess(3, 7, 20, 0.3)
	process4 := NewProcess(4, 5, 40, 0.01)
	process5 := NewProcess(5, 5, 5, 0.8)
	processes = append(processes, process1)
	processes = append(processes, process2)
	processes = append(processes, process3)
	processes = append(processes, process4)
	processes = append(processes, process5)

	// Cria o CPU
	// Pârametros: (núcleos)
	cpu := NewCpu(1)

	// Criar o escalonador de processos
	// Pârametros: (cpu, lista de processos, duração do ciclo para visualização)
	scheduler := NewProcessScheduler(cpu, processes, 1*time.Millisecond)
	scheduler.Run()

	// Loop que verifica se tem algum processo para executar
	for {
		// Se não tem mais processos rodando, exibe o histórico de execução de cada uma dos processos
		if scheduler.totalProcesses == 0 {
			fmt.Println()
			PrintColor("■ is Ready", White)
			fmt.Println()
			PrintColor("■ is Running", Green)
			fmt.Println()
			PrintColor("■ is Waiting", Red)
			fmt.Println()
			fmt.Println()
			for _, process := range processes {
				fmt.Println("Process", process.id)
				process.showHistory()
				fmt.Println()
			}
			break
		}
	}
}

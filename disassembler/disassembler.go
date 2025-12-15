package disassembler

import (
	"fmt"
	"log"
	"math"
	"nes-go/emulator"
	"nes-go/mos6502"
	"net/http"
)

type Disassembler struct {
	Instructions map[uint16]*mos6502.Instruction
	Cpu          *mos6502.CPU
	startPc      uint16
}

func NewDisassembler(cpu *mos6502.CPU) *Disassembler {
	instructions := make(map[uint16]*mos6502.Instruction)
	startPc := cpu.Pc
	logger := emulator.GetDisassemblyLogger()

	for pc := cpu.Pc; pc < math.MaxUint16; pc = cpu.Pc {
		instruction := cpu.GetNextInstruction()
		instructions[pc] = instruction
		logger.Printf("[%04X] %v", pc, instruction.InstructionText)
	}

	return &Disassembler{
		Instructions: instructions,
		Cpu:          cpu,
		startPc:      startPc,
	}
}

func (disassembler *Disassembler) Run() {
	disassembler.Cpu.Pc = disassembler.startPc

	for {
		instruction := disassembler.Instructions[disassembler.Cpu.Pc]

		if instruction == nil {
			instruction = disassembler.Cpu.GetNextInstruction()
			//disassembler.instructions[disassembler.cpu.Pc] = instruction
		}

		disassembler.Cpu.Pc = instruction.Pc + 1
		instruction.Run(disassembler.Cpu)
	}
}

func (disassembler *Disassembler) Step() {
	currentInstruction, got := disassembler.Instructions[disassembler.Cpu.Pc]
	if !got {
		currentInstruction = disassembler.Cpu.GetNextInstruction()
	}
	disassembler.Cpu.Pc = currentInstruction.Pc + 1
	currentInstruction.Run(disassembler.Cpu)
}

func (disassembler *Disassembler) Disassemble() {
	disassembler.Cpu.Pc = disassembler.startPc

	var input string
	for {
		currentInstruction := disassembler.Instructions[disassembler.Cpu.Pc]
		if currentInstruction == nil {
			currentInstruction = disassembler.Cpu.GetNextInstruction()
		}
		disassembler.Cpu.Pc = currentInstruction.Pc + 1

		fmt.Printf("%v\n", disassembler.Cpu)
		fmt.Printf("\x1b[1;33m%v\x1b[0m\n", currentInstruction)

		next := currentInstruction
		for range 10 {
			next = disassembler.Instructions[next.NextPc]
			if next != nil {
				fmt.Printf("%v\n", next)
			}
		}

		fmt.Scanln(&input)
		currentInstruction.Run(disassembler.Cpu)
	}
}

func (disassembler *Disassembler) DisassembleWeb() {
	disassembler.Cpu.Pc = disassembler.startPc

	fmt.Println("Starting web server on port http://localhost:8080...")

	http.HandleFunc("/", serveStaticSite)
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("disassembler/assets/style"))))
	http.Handle("/scripts/", http.StripPrefix("/scripts", http.FileServer(http.Dir("disassembler/assets/scripts"))))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/instructions", disassembler.GetInstructions)
	http.HandleFunc("/step", disassembler.StepHandler)
	http.HandleFunc("/continue", disassembler.ContinueHandler)
	http.HandleFunc("/cpu-state", disassembler.GetCpuState)
	http.HandleFunc("/memory-dump", disassembler.GetMemoryDump)

	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}

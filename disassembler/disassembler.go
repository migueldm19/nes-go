package disassembler

import (
	"fmt"
	"log"
	"math"
	"nes-go/emulator"
	"nes-go/mos6502"
	"net/http"
	"text/template"
)

type Disassembler struct {
	instructions map[uint16]*mos6502.Instruction
	cpu          *mos6502.CPU
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
		instructions: instructions,
		cpu:          cpu,
		startPc:      startPc,
	}
}

func (disassembler *Disassembler) Run() {
	disassembler.cpu.Pc = disassembler.startPc

	for {
		instruction := disassembler.instructions[disassembler.cpu.Pc]

		if instruction == nil {
			instruction = disassembler.cpu.GetNextInstruction()
			//disassembler.instructions[disassembler.cpu.Pc] = instruction
		}

		disassembler.cpu.Pc = instruction.Pc + 1
		instruction.Run(disassembler.cpu)
	}
}

func (disassembler *Disassembler) Disassemble() {
	disassembler.cpu.Pc = disassembler.startPc

	var input string
	for {
		currentInstruction := disassembler.instructions[disassembler.cpu.Pc]
		if currentInstruction == nil {
			currentInstruction = disassembler.cpu.GetNextInstruction()
		}
		disassembler.cpu.Pc = currentInstruction.Pc + 1

		fmt.Printf("%v\n", disassembler.cpu)
		fmt.Printf("\x1b[1;33m%v\x1b[0m\n", currentInstruction)

		next := currentInstruction
		for range 10 {
			next = disassembler.instructions[next.NextPc]
			if next != nil {
				fmt.Printf("%v\n", next)
			}
		}

		fmt.Scanln(&input)
		currentInstruction.Run(disassembler.cpu)
	}
}

type DisassemblerPage struct {
	NextInstructions   []*mos6502.Instruction
	CurrentInstruction *mos6502.Instruction
	MemoryDump         *emulator.MemoryDump
	CpuState           mos6502.StateData
}

func (disassembler *Disassembler) DisassembleWeb() {
	disassembler.cpu.Pc = disassembler.startPc

	handler := func(w http.ResponseWriter, r *http.Request) {
		currentInstruction, got := disassembler.instructions[disassembler.cpu.Pc]
		if !got {
			currentInstruction = disassembler.cpu.GetNextInstruction()
		}
		disassembler.cpu.Pc = currentInstruction.Pc + 1

		fmt.Println(currentInstruction)

		t, err := template.ParseFiles("disassembler/assets/disassembler.html")
		if err != nil {
			fmt.Println("Error parsing template:", err)
			return
		}

		instructions := make([]*mos6502.Instruction, 0)

		next := currentInstruction
		for range 10 {
			next = disassembler.instructions[next.NextPc]
			if next != nil {
				instructions = append(instructions, next)
			}
		}

		page := DisassemblerPage{
			NextInstructions:   instructions,
			CurrentInstruction: currentInstruction,
			MemoryDump:         disassembler.cpu.Dump(),
			CpuState:           disassembler.cpu.GetStateData(),
		}

		t.Execute(w, page)
		currentInstruction.Run(disassembler.cpu)
	}

	fmt.Println("Starting web server on port http://localhost:8080...")

	http.HandleFunc("/", handler)
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("disassembler/assets/style"))))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})

	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}

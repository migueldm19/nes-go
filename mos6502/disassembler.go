package mos6502

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"text/template"
)

type Disassembler struct {
	instructions map[uint16]*Instruction
	cpu          *CPU
	startPc      uint16
}

func NewDisassembler(cpu *CPU) *Disassembler {
	instructions := make(map[uint16]*Instruction)
	startPc := cpu.pc
	logger := GetLogger()

	for pc := cpu.pc; pc < math.MaxUint16; pc = cpu.pc {
		instruction := cpu.GetNextInstruction()
		instructions[pc] = instruction
		logger.Disassembly.Printf("[%04X] %v", pc, instruction.InstructionText)
	}

	return &Disassembler{
		instructions: instructions,
		cpu:          cpu,
		startPc:      startPc,
	}
}

func (disassembler *Disassembler) Run() {
	disassembler.cpu.pc = disassembler.startPc

	for {
		instruction := disassembler.instructions[disassembler.cpu.pc]

		if instruction == nil {
			instruction = disassembler.cpu.GetNextInstruction()
			//disassembler.instructions[disassembler.cpu.pc] = instruction
		}

		disassembler.cpu.pc = instruction.Pc + 1
		instruction.Run(disassembler.cpu)
	}
}

func (disassembler *Disassembler) Disassemble() {
	disassembler.cpu.pc = disassembler.startPc

	var input string
	for {
		currentInstruction := disassembler.instructions[disassembler.cpu.pc]
		if currentInstruction == nil {
			currentInstruction = disassembler.cpu.GetNextInstruction()
		}
		disassembler.cpu.pc = currentInstruction.Pc + 1

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
	NextInstructions   []*Instruction
	CurrentInstruction *Instruction
	MemoryDump         *MemoryDump
	CpuState           string
}

func (disassembler *Disassembler) DisassembleWeb() {
	disassembler.cpu.pc = disassembler.startPc

	handler := func(w http.ResponseWriter, r *http.Request) {
		currentInstruction, got := disassembler.instructions[disassembler.cpu.pc]
		if !got {
			currentInstruction = disassembler.cpu.GetNextInstruction()
		}
		disassembler.cpu.pc = currentInstruction.Pc + 1

		fmt.Println(currentInstruction)

		t, err := template.ParseFiles("disassembler/disassembler.html")
		if err != nil {
			fmt.Println("Error parsing template:", err)
			return
		}

		instructions := make([]*Instruction, 0)

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
			CpuState:           disassembler.cpu.String(),
		}

		t.Execute(w, page)
		currentInstruction.Run(disassembler.cpu)
	}

	fmt.Println("Starting web server on port 8080...")

	http.HandleFunc("/", handler)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})

	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}

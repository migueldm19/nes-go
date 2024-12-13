package mos6502

import (
	"fmt"
	"math"
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
		logger.Disassembly.Printf("[%04X] %v", pc, instruction.instructionText)
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

		disassembler.cpu.pc = instruction.pc + 1
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
		disassembler.cpu.pc = currentInstruction.pc + 1

		fmt.Printf("%v\n", disassembler.cpu)
		fmt.Printf("\x1b[1;33m[%04X] %v\x1b[0m\n", currentInstruction.pc, currentInstruction.instructionText)

		next := currentInstruction
		for range 10 {
			next = disassembler.instructions[next.nextPc]
			if next != nil {
				fmt.Printf("[%04X] %v\n", next.pc, next.instructionText)
			}
		}

		fmt.Scanln(&input)
		currentInstruction.Run(disassembler.cpu)
	}
}

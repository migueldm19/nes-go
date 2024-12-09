package mos6502

import (
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

	for pc := cpu.pc; pc < math.MaxUint16; pc = cpu.pc {
		instruction := cpu.GetNextInstruction()
		instructions[pc] = instruction
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

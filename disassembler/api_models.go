package disassembler

import (
	"nes-go/emulator"
	"nes-go/mos6502"
)

type DisassemblerPage struct {
	NextInstructions   []*mos6502.Instruction
	CurrentInstruction *mos6502.Instruction
	MemoryDump         *emulator.MemoryDump
	CpuState           mos6502.StateData
}

type InstructionsData struct {
	Instructions map[uint16]*mos6502.Instruction
	Pc           uint16
}

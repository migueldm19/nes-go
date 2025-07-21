package ppu

import (
	"nes-go/emulator"
)

/*
* Registers: from $2000 to $2007 and repeated from $2008 to $3FFF.
* A write in $3456 is the same as a write in $2006.
 */
const (
	PPUCTRL   = 0x2000
	PPUMASK   = 0x2001
	PPUSTATUS = 0x2002
	OAMADDR   = 0x2003
	OAMDATA   = 0x2004
	PPUSCROLL = 0x2005
	PPUADDR   = 0x2006
	PPUDATA   = 0x2007
	OAMDMA    = 0x4014
)


type PPU struct {
	mem *emulator.Memory
}

func NewPPU(memory *emulator.Memory) *PPU {
	return &PPU{
		mem: memory,
	}
}

func (ppu *PPU) GetPatternTable0() PatternTable {
	data := ppu.mem.RomData.ChrData
	return createPatternTable(data, PATTERN_TABLE_0_ADDRESS)
}

func (ppu *PPU) GetPatternTable1() PatternTable {
	data := ppu.mem.RomData.ChrData
	return createPatternTable(data, PATTERN_TABLE_1_ADDRESS)
}

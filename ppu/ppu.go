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
	/*
	 *  Address range 	Size 	Descriptio
	 *	$0000-$0FFF 	$1000 	Pattern table 0
	 *	$1000-$1FFF 	$1000 	Pattern table 1
	 *	$2000-$23BF 	$03c0 	Nametable 0
	 *	$23C0-$23FF 	$0040 	Attribute table 0
	 *	$2400-$27BF 	$03c0 	Nametable 1
	 *	$27C0-$27FF 	$0040 	Attribute table 1
	 *	$2800-$2BBF 	$03c0 	Nametable 2
	 *	$2BC0-$2BFF 	$0040 	Attribute table 2
	 *	$2C00-$2FBF 	$03c0 	Nametable 3
	 *	$2FC0-$2FFF 	$0040 	Attribute table 3
	 *	$3000-$3EFF 	$0F00 	Unused
	 *	$3F00-$3F1F 	$0020 	Palette RAM indexes
	 *	$3F20-$3FFF 	$00E0 	Mirrors of $3F00-$3F1F
	 */
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

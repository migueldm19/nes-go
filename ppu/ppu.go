package ppu

import (
	"nes-go/emulator"
)

/*
* Registros: de $2000 a $2007, y repetidos de $2008 a $3FFF.
* Una escritura a $3456 es lo mismo que una escritura a $2006.
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

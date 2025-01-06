package ppu

import (
	"fmt"
	"nes-go/mos6502"
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
	mem *mos6502.Memory
}

func NewPPU(memory *mos6502.Memory) *PPU {
	return &PPU{
		mem: memory,
	}
}

func GetTile(data [16]byte) [8][8]byte {
	var tile [8][8]byte

	plane0 := data[:8]
	plane1 := data[8:]

	for i := range 8 {
		for j := range 8 {
			val0 := (plane0[i] >> j) & 1
			val1 := (plane1[i] >> j) & 1
			tile[i][7-j] = val1*2 + val0
		}
	}

	return tile
}
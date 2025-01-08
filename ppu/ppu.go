package ppu

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"nes-go/mos6502"
	"os"
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

type Tile [8][8]byte
type PatternTable [256]Tile

func NewPPU(memory *mos6502.Memory) *PPU {
	return &PPU{
		mem: memory,
	}
}

func GetTile(data []byte) Tile {
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

func (ppu *PPU) getPatternTable(start int) PatternTable {
	var pattern_table PatternTable

	data := ppu.mem.RomData.ChrData

	pt_idx := 0
	for i := start; i < start+0x1000; i += 16 {
		tile := GetTile(data[i : i+16])
		pattern_table[pt_idx] = tile
		pt_idx++
	}

	return pattern_table
}

func (ppu *PPU) GetPatternTable0() PatternTable {
	return ppu.getPatternTable(0x0000)
}

func (ppu *PPU) GetPatternTable1() PatternTable {
	return ppu.getPatternTable(0x1000)
}

func printTile(img *image.RGBA, tile Tile, row, col int) {
	cyan := color.RGBA{128, 128, 128, 0xff}

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			x := col*8 + j
			y := row*8 + i
			switch tile[i][j] {
			case 0:
				img.Set(x, y, color.Transparent)
			case 1:
				img.Set(x, y, color.White)
			case 2:
				img.Set(x, y, cyan)
			case 3:
				img.Set(x, y, color.Black)
			}
		}
	}
}

func GenerateImage(path string, table PatternTable) {
	width := 8 * 16
	height := 8 * 16

	upLeft := image.Point{0, 0}
	bottomRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, bottomRight})

	for row := 0; row < 16; row++ {
		for col := 0; col < 16; col++ {
			tile := table[col*8+row]
			printTile(img, tile, row, col)
		}
	}

	f, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error creating path %v: %v", path, err)
		return
	}

	png.Encode(f, img)
}

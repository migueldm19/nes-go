package ppu

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

const (
	TILE_SIZE = 8
	TILE_RAW_BITS_SIZE = 8
	PATTERN_TABLE_N_TILES = 256
	PATTERN_TABLE_SIZE = 16
	PATTERN_TABLE_0_ADDRESS = 0x0000
	PATTERN_TABLE_1_ADDRESS = 0x1000
)

type Tile [TILE_SIZE][TILE_SIZE]byte
type PatternTable [PATTERN_TABLE_N_TILES]Tile

func GetTile(data []byte) (tile Tile) {
	plane0 := data[:TILE_RAW_BITS_SIZE]
	plane1 := data[TILE_RAW_BITS_SIZE:]

	for i := range TILE_RAW_BITS_SIZE {
		for j := range TILE_RAW_BITS_SIZE {
			val0 := (plane0[i] >> j) & 1
			val1 := (plane1[i] >> j) & 1
			tile[i][7-j] = val1*2 + val0 // combine plane0 and plane1 values
		}
	}

	return tile
}

func (ppu *PPU) getPatternTable(start int) PatternTable {
	var pattern_table PatternTable

	data := ppu.mem.RomData.ChrData
	tileRawSize := TILE_RAW_BITS_SIZE * 2 // two planes

	pt_idx := 0
	for i := start; i < start+0x1000; i += tileRawSize {
		tile := GetTile(data[i : i+tileRawSize])
		pattern_table[pt_idx] = tile
		pt_idx++
	}

	return pattern_table
}

func (ppu *PPU) GetPatternTable0() PatternTable {
	return ppu.getPatternTable(PATTERN_TABLE_0_ADDRESS)
}

func (ppu *PPU) GetPatternTable1() PatternTable {
	return ppu.getPatternTable(PATTERN_TABLE_1_ADDRESS)
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
	width := TILE_SIZE * PATTERN_TABLE_SIZE
	height := TILE_SIZE * PATTERN_TABLE_SIZE

	upLeft := image.Point{0, 0}
	bottomRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, bottomRight})

	for row := 0; row < PATTERN_TABLE_SIZE; row++ {
		for col := 0; col < PATTERN_TABLE_SIZE; col++ {
			tile := table[col*TILE_SIZE+row]
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

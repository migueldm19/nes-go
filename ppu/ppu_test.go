package ppu

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTile(t *testing.T) {
	data := [16]byte{0x41, 0xc2, 0x44, 0x48, 0x10, 0x20, 0x40, 0x80, 0x01, 0x02, 0x04, 0x08, 0x16, 0x21, 0x42, 0x87}
	expected := Tile{
		{0, 1, 0, 0, 0, 0, 0, 3},
		{1, 1, 0, 0, 0, 0, 3, 0},
		{0, 1, 0, 0, 0, 3, 0, 0},
		{0, 1, 0, 0, 3, 0, 0, 0},
		{0, 0, 0, 3, 0, 2, 2, 0},
		{0, 0, 3, 0, 0, 0, 0, 2},
		{0, 3, 0, 0, 0, 0, 2, 0},
		{3, 0, 0, 0, 0, 2, 2, 2},
	}

	tile := GetTile(data[:])
	for _, row := range tile {
		fmt.Println(row)
	}
	assert.Equal(t, expected, tile)
}

package mos6502

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddOverflow(t *testing.T) {
	var result byte
	var overflow bool

	result, overflow = addOverflow(0x10, 0x30)
	assert.Equal(t, result, byte(0x40))
	assert.Equal(t, overflow, false)

	result, overflow = addOverflow(0x00, 0x00)
	assert.Equal(t, result, byte(0x00))
	assert.Equal(t, overflow, false)

	result, overflow = addOverflow(0xff, 0x30)
	assert.Equal(t, result, byte(0x2f))
	assert.Equal(t, overflow, true)
}

func TestSubOverflow(t *testing.T) {
	var result byte
	var overflow bool

	result, overflow = subOverflow(0x30, 0x10)
	assert.Equal(t, byte(0x20), result)
	assert.Equal(t, false, overflow)

	result, overflow = subOverflow(0x00, 0x00)
	assert.Equal(t, byte(0x00), result)
	assert.Equal(t, false, overflow)

	result, overflow = subOverflow(0x10, 0x30)
	assert.Equal(t, byte(0xe0), result)
	assert.Equal(t, true, overflow)
}

func TestMemoryRead(t *testing.T) {
	rom_data := slices.Repeat([]byte{0}, 16400)
	rom_data[4] = 1
	rom := NewRom(rom_data)
	mem := NewMemory(rom)

	var addr uint16
	var val byte = 0xaa

	addr = 0x100
	mem.Data[addr] = val

	var out byte
	var err error

	out, err = mem.Read(addr)
	assert.Nil(t, err)
	assert.Equal(t, val, out)

	mem.PrgRom.Data[0x50] = val
	addr = MEMORY_SIZE + 0x50

	out, err = mem.Read(addr)
	assert.Nil(t, err)
	assert.Equal(t, val, out)
}

func TestMemoryWrite(t *testing.T) {
	rom_data := slices.Repeat([]byte{0}, 16400)
	rom_data[4] = 1
	rom := NewRom(rom_data)
	mem := NewMemory(rom)

	var addr uint16
	var val byte = 0xaa

	addr = 0x100

	var err error

	err = mem.Write(val, addr)
	assert.Nil(t, err)
	assert.Equal(t, val, mem.Data[addr])

	addr = MEMORY_SIZE + 0x50

	err = mem.Write(val, addr)
	assert.Nil(t, err)
	assert.Equal(t, val, mem.PrgRom.Data[0x50])
}

package mos6502

import (
	"nes-go/emulator"
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

func getEmptyRom() *emulator.Rom {
	rom_data := slices.Repeat([]byte{0}, 16400)
	rom_data[4] = 1
	return emulator.NewRom(rom_data)
}

func TestMemoryRead(t *testing.T) {
	rom := getEmptyRom()
	mem := emulator.NewMemory(rom)

	var addr uint16
	var val byte = 0xaa

	addr = 0x50
	mem.CPUData[addr] = val

	var out byte
	var err error

	out, err = mem.ReadCpu(addr)
	assert.Nil(t, err)
	assert.Equal(t, val, out)

	mem.RomData.PrgData[0x50] = val
	addr = emulator.CPU_MEMORY_SIZE + 0x50

	out, err = mem.ReadCpu(addr)
	assert.Nil(t, err)
	assert.Equal(t, val, out)
}

func TestMemoryWrite(t *testing.T) {
	rom := getEmptyRom()
	mem := emulator.NewMemory(rom)

	var addr uint16
	var val byte = 0xaa

	addr = 0x50

	var err error

	err = mem.WriteCpu(val, addr)
	assert.Nil(t, err)
	assert.Equal(t, val, mem.CPUData[addr])

	addr = emulator.CPU_MEMORY_SIZE + 0x50

	err = mem.WriteCpu(val, addr)
	assert.Nil(t, err)
	assert.Equal(t, val, mem.RomData.PrgData[0x50])
}

func TestRor(t *testing.T) {
	rom := getEmptyRom()
	mem := emulator.NewMemory(rom)
	cpu := NewCPU(mem)

	var addr uint16 = 0x50
	var val byte = 0xff
	var out byte

	cpu.write(val, addr)
	cpu.setFlag(FlagCarry, true)
	cpu.ror(addr)
	out = cpu.read(addr)

	assert.Equal(t, byte(0xff), out)
	assert.True(t, cpu.getFlag(FlagCarry))
	assert.False(t, cpu.getFlag(FlagZero))
	assert.True(t, cpu.getFlag(FlagNegative))

	val = 0x00
	cpu.write(val, addr)
	cpu.ror(addr)
	out = cpu.read(addr)

	assert.Equal(t, byte(0x80), out)
	assert.False(t, cpu.getFlag(FlagCarry))
	assert.False(t, cpu.getFlag(FlagZero))
	assert.True(t, cpu.getFlag(FlagNegative))

	cpu.ror(addr)
	out = cpu.read(addr)

	assert.Equal(t, byte(0x40), out)
	assert.False(t, cpu.getFlag(FlagCarry))
	assert.False(t, cpu.getFlag(FlagZero))
	assert.False(t, cpu.getFlag(FlagNegative))

	val = 0xaa
	cpu.write(val, addr)
	for range 9 {
		cpu.ror(addr)
	}
	assert.Equal(t, val, cpu.read(addr))
}

func TestRol(t *testing.T) {
	rom := getEmptyRom()
	mem := emulator.NewMemory(rom)
	cpu := NewCPU(mem)

	var addr uint16 = 0x50
	var val byte = 0xff
	var out byte

	cpu.write(val, addr)
	cpu.setFlag(FlagCarry, true)
	cpu.rol(addr)
	out = cpu.read(addr)

	assert.Equal(t, byte(0xff), out)
	assert.True(t, cpu.getFlag(FlagCarry))
	assert.False(t, cpu.getFlag(FlagZero))
	assert.True(t, cpu.getFlag(FlagNegative))

	val = 0x00
	cpu.write(val, addr)
	cpu.rol(addr)
	out = cpu.read(addr)

	assert.Equal(t, byte(0x01), out)
	assert.False(t, cpu.getFlag(FlagCarry))
	assert.False(t, cpu.getFlag(FlagZero))
	assert.False(t, cpu.getFlag(FlagNegative))

	cpu.rol(addr)
	out = cpu.read(addr)

	assert.Equal(t, byte(0x02), out)
	assert.False(t, cpu.getFlag(FlagCarry))
	assert.False(t, cpu.getFlag(FlagZero))
	assert.False(t, cpu.getFlag(FlagNegative))

	val = 0xaa
	cpu.write(val, addr)
	for range 9 {
		cpu.rol(addr)
	}
	assert.Equal(t, val, cpu.read(addr))
}

func TestLsr(t *testing.T) {
	rom := getEmptyRom()
	mem := emulator.NewMemory(rom)
	cpu := NewCPU(mem)

	var addr uint16 = 0x50
	var val byte = 0xff
	var out byte

	cpu.write(val, addr)
	cpu.lsr(addr)
	out = cpu.read(addr)

	assert.Equal(t, byte(0x7f), out)
	assert.True(t, cpu.getFlag(FlagCarry))
	assert.False(t, cpu.getFlag(FlagZero))
	assert.False(t, cpu.getFlag(FlagNegative))

	val = 0x00
	cpu.write(val, addr)
	cpu.lsr(addr)
	out = cpu.read(addr)

	assert.Equal(t, byte(0x00), out)
	assert.False(t, cpu.getFlag(FlagCarry))
	assert.True(t, cpu.getFlag(FlagZero))
	assert.False(t, cpu.getFlag(FlagNegative))

	val = 0x10
	cpu.write(val, addr)
	cpu.lsr(addr)
	out = cpu.read(addr)

	assert.Equal(t, byte(0x08), out)
	assert.False(t, cpu.getFlag(FlagCarry))
	assert.False(t, cpu.getFlag(FlagZero))
	assert.False(t, cpu.getFlag(FlagNegative))
}

func TestAsl(t *testing.T) {
	rom := getEmptyRom()
	mem := emulator.NewMemory(rom)
	cpu := NewCPU(mem)

	var addr uint16 = 0x50
	var val byte = 0xff
	var out byte

	cpu.write(val, addr)
	cpu.asl(addr)
	out = cpu.read(addr)

	assert.Equal(t, byte(0xfe), out)
	assert.True(t, cpu.getFlag(FlagCarry))
	assert.False(t, cpu.getFlag(FlagZero))
	assert.True(t, cpu.getFlag(FlagNegative))

	val = 0x00
	cpu.write(val, addr)
	cpu.asl(addr)
	out = cpu.read(addr)

	assert.Equal(t, byte(0x00), out)
	assert.False(t, cpu.getFlag(FlagCarry))
	assert.True(t, cpu.getFlag(FlagZero))
	assert.False(t, cpu.getFlag(FlagNegative))

	val = 0x10
	cpu.write(val, addr)
	cpu.asl(addr)
	out = cpu.read(addr)

	assert.Equal(t, byte(0x20), out)
	assert.False(t, cpu.getFlag(FlagCarry))
	assert.False(t, cpu.getFlag(FlagZero))
	assert.False(t, cpu.getFlag(FlagNegative))
}

func TestIndirectXAddressing(t *testing.T) {
	rom := getEmptyRom()
	mem := emulator.NewMemory(rom)
	cpu := NewCPU(mem)

	cpu.write(0x00, 0)
	cpu.write(0x02, 1)

	out, _ := cpu.nextAddress(IndirectX)
	assert.Equal(t, uint16(0x0200), out)
}

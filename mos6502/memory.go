package mos6502

import (
	"fmt"
)

const MEMORY_SIZE = 0x8000

type Memory struct {
	Data   [MEMORY_SIZE]byte
	PrgRom *Rom
}

func NewMemory(cartridge *Rom) *Memory {
	return &Memory{PrgRom: cartridge}
}

func (mem *Memory) Read(address uint16) (byte, error) {
	if address < MEMORY_SIZE {
		return mem.Data[address], nil
	}

	address -= MEMORY_SIZE
	if int(address) > len(mem.PrgRom.Data) {
		return 0, fmt.Errorf("Read index out of range: %04X", address)
	}

	return mem.PrgRom.Data[address], nil
}

func (mem *Memory) Write(value byte, address uint16) error {
	if address < MEMORY_SIZE {
		mem.Data[address] = value
		return nil
	}

	address -= MEMORY_SIZE
	if int(address) > len(mem.PrgRom.Data) {
		return fmt.Errorf("Write index out of range: %04X", address)
	}

	mem.PrgRom.Data[address] = value
	return nil
}

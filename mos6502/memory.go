package mos6502

import (
	"fmt"
)

const MEMORY_SIZE = 0x8000

type Memory struct {
	Data   [MEMORY_SIZE]byte
	PrgRom Rom
}

func NewMemory(cartridge []byte) *Memory {
	return &Memory{PrgRom: *NewRom(cartridge)}
}

func (mem *Memory) Read(address uint16) (byte, error) {
	if address < MEMORY_SIZE {
		return mem.Data[address], nil
	}

	address = address - MEMORY_SIZE
	if int(address) > len(mem.PrgRom.Data) {
		return 0, fmt.Errorf("Read index out of range: %v", address)
	}

	return mem.PrgRom.Data[address], nil
}

func (mem *Memory) Write(value byte, address uint16) error {
	if address < MEMORY_SIZE {
		mem.Data[address] = value
		return nil
	}

	address = address - MEMORY_SIZE
	if int(address) > len(mem.PrgRom.Data) {
		return fmt.Errorf("Write index out of range: %v", address)
	}

	mem.PrgRom.Data[address] = value
	return nil
}

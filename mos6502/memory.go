package mos6502

import (
	"fmt"
)

const MEMORY_SIZE = 0x8000

const ZERO_PAGE_START = 0x0000
const ZERO_PAGE_FINISH = 0x0100

const STACK_START = 0x0100
const STACK_FINISH = 0x0200

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

func (mem Memory) getDump(start, finish, step int) (dump string) {
	for i, prev := start+step, start; i <= finish; i += step {
		dump += fmt.Sprintf("%04X\t", prev)
		for _, b := range mem.Data[prev:i] {
			dump += fmt.Sprintf("%02X ", b)
		}
		dump += "\n"
		prev = i
	}
	return
}

func (mem Memory) ZeroPageDump() (dump string) {
	dump += "------- ZERO PAGE -------\n"
	dump += mem.getDump(ZERO_PAGE_START, ZERO_PAGE_FINISH, 32)
	return
}

func (mem Memory) StackDump() (dump string) {
	dump += "------- STACK -------\n"
	dump += mem.getDump(STACK_START, STACK_FINISH, 32)
	return
}

package emulator

import (
	"fmt"
)

const CPU_MEMORY_SIZE = 0x8000
const PPU_MEMORY_SIZE = 0x4000

const ZERO_PAGE_START = 0x0000
const ZERO_PAGE_FINISH = 0x0100

const STACK_START = 0x0100
const STACK_FINISH = 0x0200

type Memory struct {
	CPUData [CPU_MEMORY_SIZE]byte
	PPUData [PPU_MEMORY_SIZE]byte
	RomData *Rom
}

func NewMemory(cartridge *Rom) *Memory {
	return &Memory{RomData: cartridge}
}

func (mem *Memory) ReadCpu(address uint16) (byte, error) {
	if address < CPU_MEMORY_SIZE {
		return mem.CPUData[address], nil
	}

	address -= CPU_MEMORY_SIZE
	if int(address) > len(mem.RomData.PrgData) {
		return 0, fmt.Errorf("Read index out of range: %04X", address)
	}

	return mem.RomData.PrgData[address], nil
}

func (mem *Memory) WriteCpu(value byte, address uint16) error {
	if address < CPU_MEMORY_SIZE {
		mem.CPUData[address] = value
		return nil
	}

	address -= CPU_MEMORY_SIZE
	if int(address) > len(mem.RomData.PrgData) {
		return fmt.Errorf("Write index out of range: %04X", address)
	}

	mem.RomData.PrgData[address] = value
	return nil
}

func (mem Memory) getDump(start, finish, step int) map[int]string {
	dump := make(map[int]string, 0)

	for i, prev := start+step, start; i <= finish; i += step {
		dump_str := ""
		for _, b := range mem.CPUData[prev:i] {
			dump_str += fmt.Sprintf("%02X ", b)
		}
		dump[prev] = dump_str
		prev = i
	}

	return dump
}

func (mem Memory) ZeroPageDump() map[int]string {
	return mem.getDump(ZERO_PAGE_START, ZERO_PAGE_FINISH, 32)
}

func (mem Memory) StackDump() map[int]string {
	return mem.getDump(STACK_START, STACK_FINISH, 32)
}

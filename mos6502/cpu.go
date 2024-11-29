package mos6502

import (
	"fmt"
	"log"
)

type AdressingMode int

const (
	Immediate = iota
	ZeroPage
	ZeroPageX
	ZeroPageY
	Absolute
	AbsoluteX
	AbsoluteY
	Indirect
	IndirectX
	IndirectY
)

type Flag byte

const (
	FlagCarry = 1 << iota
	FlagZero
	FlagInterruptDisable
	FlagDecimalMode
	FlagB
	_
	FlagOverflow
	FlagNegative
)

type CPU struct {
	a   byte
	x   byte
	y   byte
	pc  uint16
	sp  byte
	p   byte
	mem Memory
}

func NewCPU(cartridge *Rom) *CPU {
	return &CPU{
		p:   0x24,
		pc:  0xc000,
		sp:  0xfd,
		mem: *NewMemory(cartridge),
	}
}

func (cpu *CPU) nextInstruction() byte {
	val, err := cpu.mem.Read(cpu.pc)

	if err != nil {
		log.Fatalf("Error geting next CPU instruction: %v", err)
	}

	cpu.pc += 1
	return val
}

func (cpu *CPU) stackPush(val byte) {
	addr := 0x0100 + uint16(cpu.sp)
	cpu.mem.Write(val, addr)

	if cpu.sp == 0 {
		log.Fatal("Stack overflow!")
	}

	cpu.sp -= 1
}

func (cpu *CPU) stackPull() (val byte) {
	if cpu.sp < 0xff {
		cpu.sp += 1
		val, _ = cpu.mem.Read(0x0100 + uint16(cpu.sp))
		return
	}

	log.Fatal("Empty stack!")
	return
}

func (cpu *CPU) stackPushCurrentPc(displacement int16) {
	addr := cpu.pc
	if displacement < 0 {
		addr -= uint16(-displacement)
	} else {
		addr += uint16(displacement)
	}

	cpu.stackPush(byte((addr & 0xff00) >> 8))
	cpu.stackPush(byte(addr & 0x00ff))
}

func (cpu *CPU) stackPullAddr() uint16 {
	addr := uint16(cpu.stackPull())
	addr += uint16(cpu.stackPull()) << 8

	return addr
}

func (cpu *CPU) setFlag(flag Flag, val bool) {
	if val {
		cpu.p |= byte(flag)
	} else {
		cpu.p &= ^byte(flag)
	}
}

func (cpu *CPU) getFlag(flag Flag) bool {
	return cpu.p&byte(flag) != 0
}

func (cpu *CPU) assignBasicFlags(val byte) {
	cpu.setFlag(FlagZero, val == 0)
	cpu.setFlag(FlagNegative, isNegative(val))
}

func (cpu *CPU) bitTest(val byte) {
	cpu.setFlag(FlagZero, (val&cpu.a) == 0)
	cpu.setFlag(FlagNegative, isNegative(val))
	cpu.setFlag(FlagOverflow, val&0x40 == 0x40)
}

func (cpu *CPU) compare(val1, val2 byte) {
	cpu.setFlag(FlagCarry, val1 >= val2)
	cpu.setFlag(FlagZero, val1 == val2)
	cpu.setFlag(FlagNegative, isNegative(val1-val2))
}

func (cpu *CPU) branchJump(displacement int8) {
	if displacement < 0 {
		cpu.pc -= uint16(uint8(-displacement))
	} else {
		cpu.pc += uint16(displacement)
	}
}

func (cpu *CPU) nextAddrHelper() uint16 {
	addr := uint16(cpu.nextInstruction())
	addr += uint16(cpu.nextInstruction()) << 8
	return addr
}

func (cpu *CPU) nextAddress(am AdressingMode) (addr, originalAddr uint16) {
	switch am {
	case ZeroPage:
		originalAddr = uint16(cpu.nextInstruction())
		addr = originalAddr
	case ZeroPageX:
		originalAddr = uint16(cpu.nextInstruction())
		addr = (originalAddr + uint16(cpu.x)) % 256
	case ZeroPageY:
		originalAddr = uint16(cpu.nextInstruction())
		addr = (originalAddr + uint16(cpu.y)) % 256
	case Absolute:
		originalAddr = cpu.nextAddrHelper()
		addr = originalAddr
	case AbsoluteX:
		originalAddr = cpu.nextAddrHelper() + uint16(cpu.x)
		addr = originalAddr
	case AbsoluteY:
		originalAddr = cpu.nextAddrHelper() + uint16(cpu.y)
		addr = originalAddr
	case IndirectX:
		addr, originalAddr = cpu.nextAddress(ZeroPageX)
		addr_1_b, _ := cpu.mem.Read(addr)
		addr_2_b, _ := cpu.mem.Read((addr + 1) % 256)
		addr = uint16(addr_2_b)<<8 + uint16(addr_1_b)
	case IndirectY:
		addr, originalAddr = cpu.nextAddress(Indirect)
		addr = addr + uint16(cpu.y)
	case Indirect:
		originalAddr := cpu.nextAddrHelper()
		addr_1_b, _ := cpu.mem.Read(originalAddr)
		addr_2_b, _ := cpu.mem.Read((originalAddr + 1) % 256)
		addr = uint16(addr_1_b) + (uint16(addr_2_b) << 8)
	}

	return
}

func (cpu *CPU) nextValue(am AdressingMode) (val byte, originalAddr uint16) {
	if am == Immediate {
		val = cpu.nextInstruction()
		return
	}

	var err error
	var addr uint16

	addr, originalAddr = cpu.nextAddress(am)
	val, err = cpu.mem.Read(addr)

	if err != nil {
		log.Fatalf("Error geting next value! %v", err)
	}

	return
}

func (cpu *CPU) write(val byte, addr uint16) {
	err := cpu.mem.Write(val, addr)

	if err != nil {
		log.Fatalf("Error in cpu write! %v", err)
	}
}

func (cpu *CPU) read(addr uint16) byte {
	val, err := cpu.mem.Read(addr)

	if err != nil {
		log.Fatalf("Error in cpu read! %v", err)
	}

	return val
}

func (cpu CPU) String() string {
	return fmt.Sprintf("A:%02X X:%02X Y:%02X P:%02X SP:%02X", cpu.a, cpu.x, cpu.y, cpu.p, cpu.sp)
}

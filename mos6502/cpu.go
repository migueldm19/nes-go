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

func NewCPU(cartridge []byte) *CPU {
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

func (cpu *CPU) nextAddress(am AdressingMode) (addr uint16) {
	switch am {
	case ZeroPage:
		addr = uint16(cpu.nextInstruction())
	case ZeroPageX:
		addr = uint16(cpu.nextInstruction() + cpu.x)
	case ZeroPageY:
		addr = uint16(cpu.nextInstruction() + cpu.y)
	case Absolute:
		addr = cpu.nextAddrHelper()
	case AbsoluteX:
		addr = cpu.nextAddrHelper() + uint16(cpu.x)
	case AbsoluteY:
		addr = cpu.nextAddrHelper() + uint16(cpu.y)
	case IndirectX:
		idx := uint16(cpu.nextInstruction() + cpu.x)
		addr_b, _ := cpu.mem.Read(idx)
		addr = uint16(addr_b) << 8
	case IndirectY:
		idx := uint16(cpu.nextInstruction())
		addr_b, _ := cpu.mem.Read(idx)
		addr = (uint16(addr_b) << 8) + uint16(cpu.y)
	case Indirect:
		idx := cpu.nextAddrHelper()
		val1, _ := cpu.mem.Read(idx)
		val2, _ := cpu.mem.Read(idx + 1)
		addr = uint16(val1) + (uint16(val2) << 8)
	}

	return
}

func (cpu *CPU) write(val byte, addr uint16) {
	err := cpu.mem.Write(val, addr)

	if err != nil {
		log.Fatalf("Error in cpu write! %v", err)
	}

	return
}

func (cpu *CPU) read(addr uint16) byte {
	val, err := cpu.mem.Read(addr)

	if err != nil {
		log.Fatalf("Error in cpu read! %v", err)
	}

	return val
}

func (cpu *CPU) nextValue(am AdressingMode) (val byte) {
	var err error

	switch am {
	case Immediate:
		val = cpu.nextInstruction()
	case ZeroPage:
		addr := uint16(cpu.nextInstruction())
		val, err = cpu.mem.Read(addr)
	case ZeroPageX:
		addr := uint16(cpu.nextInstruction() + cpu.x)
		val, err = cpu.mem.Read(addr)
	case ZeroPageY:
		addr := uint16(cpu.nextInstruction() + cpu.y)
		val, err = cpu.mem.Read(addr)
	case Absolute:
		addr := cpu.nextAddrHelper()
		val, err = cpu.mem.Read(addr)
	case AbsoluteX:
		addr := cpu.nextAddrHelper() + uint16(cpu.x)
		val, err = cpu.mem.Read(addr)
	case AbsoluteY:
		addr := cpu.nextAddrHelper() + uint16(cpu.y)
		val, err = cpu.mem.Read(addr)
	case IndirectX:
		idx := cpu.nextInstruction() + cpu.x
		val, err = cpu.mem.Read(uint16(idx))
		if err != nil {
			addr := uint16(val) << 8
			val, err = cpu.mem.Read(addr)
		}
	case IndirectY:
		addr := uint16(cpu.nextInstruction())
		val, err = cpu.mem.Read(addr)
		if err != nil {
			addr = (uint16(val) << 8) + uint16(cpu.y)
			val, err = cpu.mem.Read(addr)
		}
	}

	if err != nil {
		log.Fatalf("Error geting next value! %v", err)
	}

	return
}

func (cpu CPU) String() string {
	return fmt.Sprintf("A:%2x X:%2x Y:%2x P:%2x SP:%2x", cpu.a, cpu.x, cpu.y, cpu.p, cpu.sp)
}

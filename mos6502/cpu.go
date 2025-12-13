package mos6502

import (
	"fmt"
	"log"
	"nes-go/emulator"
)

type AdressingMode int

const (
	ZERO_PAGE_SIZE = 0x100
	BYTE_SIZE      = 8
)

const (
	Immediate AdressingMode = iota
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
	_ // Unused flag
	FlagOverflow
	FlagNegative
)

type CPU struct {
	a   byte
	x   byte
	y   byte
	Pc  uint16
	sp  byte
	p   byte
	Mem *emulator.Memory
}

func NewCPU(memory *emulator.Memory) *CPU {
	return &CPU{
		p:   0x24,
		Pc:  0xc000,
		sp:  0xfd,
		Mem: memory,
	}
}

func (cpu *CPU) Step() {
	instruction := cpu.GetNextInstruction()
	cpu.Pc = instruction.Pc + 1
	instruction.Run(cpu)
}

func (cpu *CPU) Run() {
	for {
		cpu.Step()
	}
}

func (cpu *CPU) nextInstruction() byte {
	val, err := cpu.Mem.ReadCpu(cpu.Pc)

	if err != nil {
		log.Fatalf("Error geting next CPU instruction: %v", err)
	}

	cpu.Pc += 1
	return val
}

func (cpu *CPU) stackPush(val byte) {
	addr := ZERO_PAGE_SIZE + uint16(cpu.sp)
	cpu.Mem.WriteCpu(val, addr)

	if cpu.sp == 0 {
		log.Fatal("Stack overflow!")
	}

	cpu.sp -= 1
}

func (cpu *CPU) stackPull() (val byte) {
	if cpu.sp < 0xff {
		cpu.sp += 1
		val, _ = cpu.Mem.ReadCpu(ZERO_PAGE_SIZE + uint16(cpu.sp))
		return
	}

	log.Fatal("Empty stack!")
	return
}

func (cpu *CPU) stackPushCurrentPc(displacement int16) {
	addr := cpu.Pc
	if displacement < 0 {
		addr -= uint16(-displacement)
	} else {
		addr += uint16(displacement)
	}

	cpu.stackPush(byte((addr & 0xff00) >> BYTE_SIZE))
	cpu.stackPush(byte(addr & 0x00ff))
}

func (cpu *CPU) stackPullAddr() uint16 {
	addr := uint16(cpu.stackPull())
	addr += uint16(cpu.stackPull()) << BYTE_SIZE

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
		cpu.Pc -= uint16(uint8(-displacement))
	} else {
		cpu.Pc += uint16(displacement)
	}
}

func (cpu *CPU) nextAddrHelper() uint16 {
	addr := uint16(cpu.nextInstruction())
	addr += uint16(cpu.nextInstruction()) << BYTE_SIZE
	return addr
}

func (cpu *CPU) nextAddress(am AdressingMode) (addr, originalAddr uint16) {
	switch am {
	case ZeroPage:
		originalAddr = uint16(cpu.nextInstruction())
		addr = originalAddr
	case ZeroPageX:
		originalAddr = uint16(cpu.nextInstruction())
		addr = (originalAddr + uint16(cpu.x)) % ZERO_PAGE_SIZE
	case ZeroPageY:
		originalAddr = uint16(cpu.nextInstruction())
		addr = (originalAddr + uint16(cpu.y)) % ZERO_PAGE_SIZE
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
		addr_1_b, _ := cpu.Mem.ReadCpu(addr)
		addr_2_b, _ := cpu.Mem.ReadCpu((addr + 1) % ZERO_PAGE_SIZE)
		addr = uint16(addr_1_b) + uint16(addr_2_b)<<BYTE_SIZE
	case IndirectY:
		addr, originalAddr = cpu.nextAddress(ZeroPage)
		addr_1_b, _ := cpu.Mem.ReadCpu(addr)
		addr_2_b, _ := cpu.Mem.ReadCpu((addr + 1) % ZERO_PAGE_SIZE)
		addr = uint16(addr_1_b) + uint16(addr_2_b)<<BYTE_SIZE + uint16(cpu.y)
	case Indirect:
		originalAddr = cpu.nextAddrHelper()
		addr_1_b, _ := cpu.Mem.ReadCpu(originalAddr)
		// Due to a bug in the cpu, indirect addressing can't
		// cross pages, so it goes to the beginning of the page
		if originalAddr&0x00ff == 0x00ff {
			originalAddr -= ZERO_PAGE_SIZE
		}
		addr_2_b, _ := cpu.Mem.ReadCpu(originalAddr + 1)
		addr = uint16(addr_1_b) + uint16(addr_2_b)<<BYTE_SIZE
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
	val, err = cpu.Mem.ReadCpu(addr)

	if err != nil {
		log.Fatalf("Error geting next value! %v", err)
	}

	return
}

func (cpu *CPU) write(val byte, addr uint16) {
	err := cpu.Mem.WriteCpu(val, addr)

	if err != nil {
		log.Fatalf("Error in cpu write! %v", err)
	}
}

func (cpu *CPU) read(addr uint16) byte {
	val, err := cpu.Mem.ReadCpu(addr)

	if err != nil {
		log.Fatalf("Error in cpu read! %v", err)
	}

	return val
}

func (cpu CPU) Dump() *emulator.MemoryDump {
	return emulator.NewMemoryDump(cpu.Mem)
}

func (cpu CPU) String() string {
	return fmt.Sprintf("A:%02X X:%02X Y:%02X P:%02X SP:%02X", cpu.a, cpu.x, cpu.y, cpu.p, cpu.sp)
}

type FlagData struct {
	Carry            bool
	Zero             bool
	InterruptDisable bool
	DecimalMode      bool
	B                bool
	Overflow         bool
	Negative         bool
}

type StateData struct {
	PC    uint16
	A     byte
	X     byte
	Y     byte
	SP    byte
	Flags FlagData
}

func (cpu CPU) GetStateData() StateData {
	return StateData{
		PC: cpu.Pc,
		A:  cpu.a,
		X:  cpu.x,
		Y:  cpu.y,
		SP: cpu.sp,
		Flags: FlagData{
			Carry:            cpu.getFlag(FlagCarry),
			Zero:             cpu.getFlag(FlagZero),
			InterruptDisable: cpu.getFlag(FlagInterruptDisable),
			DecimalMode:      cpu.getFlag(FlagDecimalMode),
			B:                cpu.getFlag(FlagB),
			Overflow:         cpu.getFlag(FlagOverflow),
			Negative:         cpu.getFlag(FlagNegative),
		},
	}
}

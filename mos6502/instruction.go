package mos6502

import (
	"fmt"
	"nes-go/emulator"
)

type Instruction struct {
	Pc, NextPc      uint16
	InstructionText string
	action          func()
}

func NewInstruction(pc, nextPc uint16, text string, action func()) *Instruction {
	return &Instruction{
		Pc:              pc,
		NextPc:          nextPc,
		InstructionText: text,
		action:          action,
	}
}

func (instruction Instruction) Run(cpu *CPU) {
	instructions_logger := emulator.GetInstructionsLogger()
	memory_dump_logger := emulator.GetMemoryDumpLogger()

	instruction_log := fmt.Sprintf("[PC: %04X] OPCODE %02X | %v | ", instruction.Pc, cpu.read(instruction.Pc), cpu)
	instruction_log += instruction.InstructionText
	instructions_logger.Print(instruction_log)
	memory_dump_logger.Printf("[PC: %04X]\n%v", instruction.Pc, cpu.Dump())

	instruction.action()
}

func (instruction Instruction) String() string {
	return fmt.Sprintf("[%04X] %v", instruction.Pc, instruction.InstructionText)
}

func (cpu *CPU) execByte(action func(byte), am AdressingMode) {
	val, _ := cpu.nextValue(am)
	action(val)
}

func (cpu *CPU) execAddr(action func(uint16), am AdressingMode) {
	addr, _ := cpu.nextAddress(am)
	action(addr)
}

func (cpu *CPU) GetNextInstruction() *Instruction {
	opcode := cpu.nextInstruction()
	var val byte
	var addr uint16
	var originalAddr uint16

	instruction_pc := cpu.Pc - 1

	var instruction *Instruction

	switch opcode {
	case 0x00:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "BRK", cpu.brk)

	case 0x40:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "RTI", cpu.rti)

	case 0xea:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "NOP", func() {})

	case 0xa9:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("LDA #$%02X", val), func() { cpu.execByte(cpu.lda, Immediate) })
	case 0xa5:
		_, addr = cpu.nextValue(ZeroPage)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("LDA $%02X", addr), func() { cpu.execByte(cpu.lda, ZeroPage) })
	case 0xb5:
		_, addr = cpu.nextValue(ZeroPageX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("LDA $%02X, X", addr), func() { cpu.execByte(cpu.lda, ZeroPageX) })
	case 0xad:
		_, addr = cpu.nextValue(Absolute)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("LDA $%04X", addr), func() { cpu.execByte(cpu.lda, Absolute) })
	case 0xbd:
		_, addr = cpu.nextValue(AbsoluteX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("LDA $%04X, X", addr), func() { cpu.execByte(cpu.lda, AbsoluteX) })
	case 0xb9:
		_, addr = cpu.nextValue(AbsoluteY)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("LDA $%04X, Y", addr), func() { cpu.execByte(cpu.lda, AbsoluteY) })
	case 0xa1:
		_, addr = cpu.nextValue(IndirectX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("LDA ($%02X, X)", addr), func() { cpu.execByte(cpu.lda, IndirectX) })
	case 0xb1:
		_, addr = cpu.nextValue(IndirectY)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("LDA ($%02X), Y", addr), func() { cpu.execByte(cpu.lda, IndirectY) })

	case 0xa2:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("LDX #$%02X", val), func() { cpu.execByte(cpu.ldx, Immediate) })
	case 0xa6:
		_, addr = cpu.nextValue(ZeroPage)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("LDX $%02X", addr), func() { cpu.execByte(cpu.ldx, ZeroPage) })
	case 0xb6:
		_, addr = cpu.nextValue(ZeroPageY)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("LDX $%02X, Y", addr), func() { cpu.execByte(cpu.ldx, ZeroPageY) })
	case 0xae:
		_, addr = cpu.nextValue(Absolute)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("LDX $%04X", addr), func() { cpu.execByte(cpu.ldx, Absolute) })
	case 0xbe:
		_, addr = cpu.nextValue(AbsoluteY)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("LDX $%04X, Y", addr), func() { cpu.execByte(cpu.ldx, AbsoluteY) })

	case 0xa0:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("LDY #$%02X", val), func() { cpu.execByte(cpu.ldy, Immediate) })
	case 0xa4:
		_, addr = cpu.nextValue(ZeroPage)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("LDY $%02X", addr), func() { cpu.execByte(cpu.ldy, ZeroPage) })
	case 0xb4:
		_, addr = cpu.nextValue(ZeroPageX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("LDY $%02X, X", addr), func() { cpu.execByte(cpu.ldy, ZeroPageX) })
	case 0xac:
		_, addr = cpu.nextValue(Absolute)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("LDY $%04X", addr), func() { cpu.execByte(cpu.ldy, Absolute) })
	case 0xbc:
		_, addr = cpu.nextValue(AbsoluteX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("LDY $%04X, X", addr), func() { cpu.execByte(cpu.ldy, AbsoluteX) })

	case 0x85:
		_, originalAddr = cpu.nextAddress(ZeroPage)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("STA $%02X", originalAddr), func() { cpu.execAddr(cpu.sta, ZeroPage) })
	case 0x95:
		_, originalAddr = cpu.nextAddress(ZeroPageX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("STA $%02X, X", originalAddr), func() { cpu.execAddr(cpu.sta, ZeroPageX) })
	case 0x8d:
		_, originalAddr = cpu.nextAddress(Absolute)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("STA $%04X", originalAddr), func() { cpu.execAddr(cpu.sta, Absolute) })
	case 0x9d:
		_, originalAddr = cpu.nextAddress(AbsoluteX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("STA $%04X, X", originalAddr), func() { cpu.execAddr(cpu.sta, AbsoluteX) })
	case 0x99:
		_, originalAddr = cpu.nextAddress(AbsoluteY)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("STA $%04X, Y", originalAddr), func() { cpu.execAddr(cpu.sta, AbsoluteY) })
	case 0x81:
		_, originalAddr = cpu.nextAddress(IndirectX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("STA ($%02X, X)", originalAddr), func() { cpu.execAddr(cpu.sta, IndirectX) })
	case 0x91:
		_, originalAddr = cpu.nextAddress(IndirectY)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("STA ($%02X), Y", originalAddr), func() { cpu.execAddr(cpu.sta, IndirectY) })

	case 0x86:
		_, originalAddr = cpu.nextAddress(ZeroPage)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("STX $%02X", originalAddr), func() { cpu.execAddr(cpu.stx, ZeroPage) })
	case 0x96:
		_, originalAddr = cpu.nextAddress(ZeroPageY)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("STX $%02X, Y", originalAddr), func() { cpu.execAddr(cpu.stx, ZeroPageY) })
	case 0x8e:
		_, originalAddr = cpu.nextAddress(Absolute)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("STX $%04X", originalAddr), func() { cpu.execAddr(cpu.stx, Absolute) })

	case 0x84:
		_, originalAddr = cpu.nextAddress(ZeroPage)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("STY $%02X", originalAddr), func() { cpu.execAddr(cpu.sty, ZeroPage) })
	case 0x94:
		_, originalAddr = cpu.nextAddress(ZeroPageX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("STY $%02X, X", originalAddr), func() { cpu.execAddr(cpu.sty, ZeroPageX) })
	case 0x8c:
		_, originalAddr = cpu.nextAddress(Absolute)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("STY $%04X", originalAddr), func() { cpu.execAddr(cpu.sty, Absolute) })

	case 0xaa:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "TAX", cpu.tax)
	case 0xa8:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "TAY", cpu.tay)
	case 0xba:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "TSX", cpu.tsx)
	case 0x8a:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "TXA", cpu.txa)
	case 0x9a:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "TXS", cpu.txs)
	case 0x98:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "TYA", cpu.tya)

	case 0x48:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "PHA", cpu.pha)
	case 0x08:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "PHP", cpu.php)
	case 0x68:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "PLA", cpu.pla)
	case 0x28:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "PLP", cpu.plp)

	case 0x29:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("AND #$%02X", val), func() { cpu.execByte(cpu.and, Immediate) })
	case 0x25:
		_, addr = cpu.nextValue(ZeroPage)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("AND $%02X", addr), func() { cpu.execByte(cpu.and, ZeroPage) })
	case 0x35:
		_, addr = cpu.nextValue(ZeroPageX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("AND $%02X, X", addr), func() { cpu.execByte(cpu.and, ZeroPageX) })
	case 0x2d:
		_, addr = cpu.nextValue(Absolute)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("AND $%04X", addr), func() { cpu.execByte(cpu.and, Absolute) })
	case 0x3d:
		_, addr = cpu.nextValue(AbsoluteX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("AND $%04X, X", addr), func() { cpu.execByte(cpu.and, AbsoluteX) })
	case 0x39:
		_, addr = cpu.nextValue(AbsoluteY)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("AND $%04X, Y", addr), func() { cpu.execByte(cpu.and, AbsoluteY) })
	case 0x21:
		_, addr = cpu.nextValue(IndirectX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("AND ($%02X, X)", addr), func() { cpu.execByte(cpu.and, IndirectX) })
	case 0x31:
		_, addr = cpu.nextValue(IndirectY)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("AND ($%02X), Y", addr), func() { cpu.execByte(cpu.and, IndirectY) })

	case 0x49:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("EOR #$%02X", val), func() { cpu.execByte(cpu.eor, Immediate) })
	case 0x45:
		_, addr = cpu.nextValue(ZeroPage)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("EOR $%02X", addr), func() { cpu.execByte(cpu.eor, ZeroPage) })
	case 0x55:
		_, addr = cpu.nextValue(ZeroPageX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("EOR $%02X, X", addr), func() { cpu.execByte(cpu.eor, ZeroPageX) })
	case 0x4d:
		_, addr = cpu.nextValue(Absolute)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("EOR $%04X", addr), func() { cpu.execByte(cpu.eor, Absolute) })
	case 0x5d:
		_, addr = cpu.nextValue(AbsoluteX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("EOR $%04X, X", addr), func() { cpu.execByte(cpu.eor, AbsoluteX) })
	case 0x59:
		_, addr = cpu.nextValue(AbsoluteY)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("EOR $%04X, Y", addr), func() { cpu.execByte(cpu.eor, AbsoluteY) })
	case 0x41:
		_, addr = cpu.nextValue(IndirectX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("EOR ($%02X, X)", addr), func() { cpu.execByte(cpu.eor, IndirectX) })
	case 0x51:
		_, addr = cpu.nextValue(IndirectY)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("EOR ($%02X), Y", addr), func() { cpu.execByte(cpu.eor, IndirectY) })

	case 0x09:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ORA #$%02X", val), func() { cpu.execByte(cpu.ora, Immediate) })
	case 0x05:
		_, addr = cpu.nextValue(ZeroPage)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ORA $%02X", addr), func() { cpu.execByte(cpu.ora, ZeroPage) })
	case 0x15:
		_, addr = cpu.nextValue(ZeroPageX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ORA $%02X, X", addr), func() { cpu.execByte(cpu.ora, ZeroPageX) })
	case 0x0d:
		_, addr = cpu.nextValue(Absolute)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ORA $%04X", addr), func() { cpu.execByte(cpu.ora, Absolute) })
	case 0x1d:
		_, addr = cpu.nextValue(AbsoluteX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ORA $%04X, X", addr), func() { cpu.execByte(cpu.ora, AbsoluteX) })
	case 0x19:
		_, addr = cpu.nextValue(AbsoluteY)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ORA $%04X, Y", addr), func() { cpu.execByte(cpu.ora, AbsoluteY) })
	case 0x01:
		_, addr = cpu.nextValue(IndirectX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ORA ($%02X, X)", addr), func() { cpu.execByte(cpu.ora, IndirectX) })
	case 0x11:
		_, addr = cpu.nextValue(IndirectY)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ORA ($%02X), Y", addr), func() { cpu.execByte(cpu.ora, IndirectY) })

	case 0x24:
		_, addr = cpu.nextValue(ZeroPage)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("BIT $%02X", addr), func() { cpu.execByte(cpu.bit, ZeroPage) })
	case 0x2c:
		_, addr = cpu.nextValue(Absolute)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("BIT $%04X", addr), func() { cpu.execByte(cpu.bit, Absolute) })

	case 0x69:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ADC #$%02X", val), func() { cpu.execByte(cpu.adc, Immediate) })
	case 0x65:
		_, addr = cpu.nextValue(ZeroPage)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ADC $%02X", addr), func() { cpu.execByte(cpu.adc, ZeroPage) })
	case 0x75:
		_, addr = cpu.nextValue(ZeroPageX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ADC $%02X, X", addr), func() { cpu.execByte(cpu.adc, ZeroPageX) })
	case 0x6d:
		_, addr = cpu.nextValue(Absolute)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ADC $%04X", addr), func() { cpu.execByte(cpu.adc, Absolute) })
	case 0x7d:
		_, addr = cpu.nextValue(AbsoluteX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ADC $%04X, X", addr), func() { cpu.execByte(cpu.adc, AbsoluteX) })
	case 0x79:
		_, addr = cpu.nextValue(AbsoluteY)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ADC $%04X, Y", addr), func() { cpu.execByte(cpu.adc, AbsoluteY) })
	case 0x61:
		_, addr = cpu.nextValue(IndirectX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ADC ($%02X, X)", addr), func() { cpu.execByte(cpu.adc, IndirectX) })
	case 0x71:
		_, addr = cpu.nextValue(IndirectY)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ADC ($%02X), Y", addr), func() { cpu.execByte(cpu.adc, IndirectY) })

	case 0xe9:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("SBC #$%02X", val), func() { cpu.execByte(cpu.sbc, Immediate) })
	case 0xe5:
		_, addr = cpu.nextValue(ZeroPage)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("SBC $%02X", addr), func() { cpu.execByte(cpu.sbc, ZeroPage) })
	case 0xf5:
		_, addr = cpu.nextValue(ZeroPageX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("SBC $%02X, X", addr), func() { cpu.execByte(cpu.sbc, ZeroPageX) })
	case 0xed:
		_, addr = cpu.nextValue(Absolute)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("SBC $%04X", addr), func() { cpu.execByte(cpu.sbc, Absolute) })
	case 0xfd:
		_, addr = cpu.nextValue(AbsoluteX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("SBC $%04X, X", addr), func() { cpu.execByte(cpu.sbc, AbsoluteX) })
	case 0xf9:
		_, addr = cpu.nextValue(AbsoluteY)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("SBC $%04X, Y", addr), func() { cpu.execByte(cpu.sbc, AbsoluteY) })
	case 0xe1:
		_, addr = cpu.nextValue(IndirectX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("SBC ($%02X, X)", addr), func() { cpu.execByte(cpu.sbc, IndirectX) })
	case 0xf1:
		_, addr = cpu.nextValue(IndirectY)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("SBC ($%02X), Y", addr), func() { cpu.execByte(cpu.sbc, IndirectY) })

	case 0xc9:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("CMP #$%02X", val), func() { cpu.execByte(cpu.cmp, Immediate) })
	case 0xc5:
		_, addr = cpu.nextValue(ZeroPage)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("CMP $%02X", addr), func() { cpu.execByte(cpu.cmp, ZeroPage) })
	case 0xd5:
		_, addr = cpu.nextValue(ZeroPageX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("CMP $%02X, X", addr), func() { cpu.execByte(cpu.cmp, ZeroPageX) })
	case 0xcd:
		_, addr = cpu.nextValue(Absolute)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("CMP $%04X", addr), func() { cpu.execByte(cpu.cmp, Absolute) })
	case 0xdd:
		_, addr = cpu.nextValue(AbsoluteX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("CMP $%04X, X", addr), func() { cpu.execByte(cpu.cmp, AbsoluteX) })
	case 0xd9:
		_, addr = cpu.nextValue(AbsoluteY)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("CMP $%04X, Y", addr), func() { cpu.execByte(cpu.cmp, AbsoluteY) })
	case 0xc1:
		_, addr = cpu.nextValue(IndirectX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("CMP ($%02X, X)", addr), func() { cpu.execByte(cpu.cmp, IndirectX) })
	case 0xd1:
		_, addr = cpu.nextValue(IndirectY)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("CMP ($%02X), Y", addr), func() { cpu.execByte(cpu.cmp, IndirectY) })

	case 0xe0:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("CPX #$%02X", val), func() { cpu.execByte(cpu.cpx, Immediate) })
	case 0xe4:
		_, addr = cpu.nextValue(ZeroPage)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("CPX $%02X", addr), func() { cpu.execByte(cpu.cpx, ZeroPage) })
	case 0xec:
		_, addr = cpu.nextValue(Absolute)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("CPX $%04X", addr), func() { cpu.execByte(cpu.cpx, Absolute) })

	case 0xc0:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("CPY #$%02X", val), func() { cpu.execByte(cpu.cpy, Immediate) })
	case 0xc4:
		_, addr = cpu.nextValue(ZeroPage)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("CPY $%02X", addr), func() { cpu.execByte(cpu.cpy, ZeroPage) })
	case 0xcc:
		_, addr = cpu.nextValue(Absolute)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("CPY $%04X", addr), func() { cpu.execByte(cpu.cpy, Absolute) })

	case 0xe6:
		_, originalAddr = cpu.nextAddress(ZeroPage)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("INC $%02X", originalAddr), func() { cpu.execAddr(cpu.inc, ZeroPage) })
	case 0xf6:
		_, originalAddr = cpu.nextAddress(ZeroPageX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("INC $%02X, X", originalAddr), func() { cpu.execAddr(cpu.inc, ZeroPageX) })
	case 0xee:
		_, originalAddr = cpu.nextAddress(Absolute)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("INC $%04X", originalAddr), func() { cpu.execAddr(cpu.inc, Absolute) })
	case 0xfe:
		_, originalAddr = cpu.nextAddress(AbsoluteX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("INC $%04X, X", originalAddr), func() { cpu.execAddr(cpu.inc, AbsoluteX) })

	case 0xc6:
		_, originalAddr = cpu.nextAddress(ZeroPage)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("DEC $%02X", originalAddr), func() { cpu.execAddr(cpu.dec, ZeroPage) })
	case 0xd6:
		_, originalAddr = cpu.nextAddress(ZeroPageX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("DEC $%02X, X", originalAddr), func() { cpu.execAddr(cpu.dec, ZeroPageX) })
	case 0xce:
		_, originalAddr = cpu.nextAddress(Absolute)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("DEC $%04X", originalAddr), func() { cpu.execAddr(cpu.dec, Absolute) })
	case 0xde:
		_, originalAddr = cpu.nextAddress(AbsoluteX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("DEC $%04X, X", originalAddr), func() { cpu.execAddr(cpu.dec, AbsoluteX) })

	case 0xe8:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "INX", cpu.inx)
	case 0xc8:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "INY", cpu.iny)

	case 0xca:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "DEX", cpu.dex)
	case 0x88:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "DEY", cpu.dey)

	case 0x0a:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "ASL A", cpu.asl_acc)
	case 0x06:
		_, originalAddr = cpu.nextAddress(ZeroPage)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ASL $%02X", originalAddr), func() { cpu.execAddr(cpu.asl, ZeroPage) })
	case 0x16:
		_, originalAddr = cpu.nextAddress(ZeroPageX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ASL $%02X, X", originalAddr), func() { cpu.execAddr(cpu.asl, ZeroPageX) })
	case 0x0e:
		_, originalAddr = cpu.nextAddress(Absolute)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ASL $%04X", originalAddr), func() { cpu.execAddr(cpu.asl, Absolute) })
	case 0x1e:
		_, originalAddr = cpu.nextAddress(AbsoluteX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ASL $%04X, X", originalAddr), func() { cpu.execAddr(cpu.asl, AbsoluteX) })

	case 0x4a:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "LSR A", cpu.lsr_acc)
	case 0x46:
		_, originalAddr = cpu.nextAddress(ZeroPage)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("LSR $%02X", originalAddr), func() { cpu.execAddr(cpu.lsr, ZeroPage) })
	case 0x56:
		_, originalAddr = cpu.nextAddress(ZeroPageX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("LSR $%02X, X", originalAddr), func() { cpu.execAddr(cpu.lsr, ZeroPageX) })
	case 0x4e:
		_, originalAddr = cpu.nextAddress(Absolute)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("LSR $%04X", originalAddr), func() { cpu.execAddr(cpu.lsr, Absolute) })
	case 0x5e:
		_, originalAddr = cpu.nextAddress(AbsoluteX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("LSR $%04X, X", originalAddr), func() { cpu.execAddr(cpu.lsr, AbsoluteX) })

	case 0x2a:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "ROL A", cpu.rol_acc)
	case 0x26:
		_, originalAddr = cpu.nextAddress(ZeroPage)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ROL $%02X", originalAddr), func() { cpu.execAddr(cpu.rol, ZeroPage) })
	case 0x36:
		_, originalAddr = cpu.nextAddress(ZeroPageX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ROL $%02X, X", originalAddr), func() { cpu.execAddr(cpu.rol, ZeroPageX) })
	case 0x2e:
		_, originalAddr = cpu.nextAddress(Absolute)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ROL $%04X", originalAddr), func() { cpu.execAddr(cpu.rol, Absolute) })
	case 0x3e:
		_, originalAddr = cpu.nextAddress(AbsoluteX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ROL $%04X, X", originalAddr), func() { cpu.execAddr(cpu.rol, AbsoluteX) })

	case 0x6a:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "ROR A", cpu.ror_acc)
	case 0x66:
		_, originalAddr = cpu.nextAddress(ZeroPage)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ROR $%02X", originalAddr), func() { cpu.execAddr(cpu.ror, ZeroPage) })
	case 0x76:
		_, originalAddr = cpu.nextAddress(ZeroPageX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ROR $%02X, X", originalAddr), func() { cpu.execAddr(cpu.ror, ZeroPageX) })
	case 0x6e:
		_, originalAddr = cpu.nextAddress(Absolute)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ROR $%04X", originalAddr), func() { cpu.execAddr(cpu.ror, Absolute) })
	case 0x7e:
		_, originalAddr = cpu.nextAddress(AbsoluteX)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("ROR $%04X, X", originalAddr), func() { cpu.execAddr(cpu.ror, AbsoluteX) })

	case 0x4c:
		_, originalAddr = cpu.nextAddress(Absolute)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("JMP $%04X", originalAddr), func() { cpu.execAddr(cpu.jmp, Absolute) })
	case 0x6c:
		_, originalAddr = cpu.nextAddress(Indirect)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("JMP ($%04X)", originalAddr), func() { cpu.execAddr(cpu.jmp, Indirect) })

	case 0x20:
		_, originalAddr = cpu.nextAddress(Absolute)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("JSR $%04X", originalAddr), func() { cpu.execAddr(cpu.jsr, Absolute) })

	case 0x60:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "RTS", cpu.rts)

	case 0x90:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("BCC #$%02X", val), func() { cpu.execByte(cpu.bcc, Immediate) })
	case 0xb0:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("BCS #$%02X", val), func() { cpu.execByte(cpu.bcs, Immediate) })
	case 0xf0:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("BEQ #$%02X", val), func() { cpu.execByte(cpu.beq, Immediate) })
	case 0x30:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("BMI #$%02X", val), func() { cpu.execByte(cpu.bmi, Immediate) })
	case 0xd0:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("BNE #$%02X", val), func() { cpu.execByte(cpu.bne, Immediate) })
	case 0x10:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("BPL #$%02X", val), func() { cpu.execByte(cpu.bpl, Immediate) })
	case 0x50:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("BVC #$%02X", val), func() { cpu.execByte(cpu.bvc, Immediate) })
	case 0x70:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, cpu.Pc, fmt.Sprintf("BVS #$%02X", val), func() { cpu.execByte(cpu.bvs, Immediate) })

	case 0x18:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "CLC", cpu.clc)
	case 0xd8:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "CLD", cpu.cld)
	case 0x58:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "CLI", cpu.cli)
	case 0xb8:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "CLV", cpu.clv)
	case 0x38:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "SEC", cpu.sec)
	case 0xf8:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "SED", cpu.sed)
	case 0x78:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "SEI", cpu.sei)

	// unofficial opcodes
	case 0x1a, 0x3a, 0x5a, 0x7a, 0xda, 0xfa:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "*NOP", func() {})
	case 0x04, 0x44, 0x64, 0x14, 0x34, 0x54, 0x74, 0xd4, 0xf4, 0x80:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "*NOP", func() { cpu.nextValue(Immediate) })
	case 0x0c, 0x1c, 0x3c, 0x5c, 0x7c, 0xdc, 0xfc:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "*NOP", func() { cpu.nextValue(Absolute) })

	default:
		instruction = NewInstruction(instruction_pc, cpu.Pc, "UNKNOWN", func() {})
	}

	return instruction
}

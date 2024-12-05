package mos6502

import "fmt"

type Instruction struct {
	PC              uint16
	InstructionText string
	action          func()
}

func NewInstruction(pc uint16, text string, action func()) *Instruction {
	return &Instruction{
		PC:              pc,
		InstructionText: text,
		action:          action,
	}
}

func (instruction Instruction) Run(cpu *CPU) {
	logger := GetLogger()
	instruction_log := fmt.Sprintf("[PC: %04X] OPCODE %02X | %v | ", instruction.PC, cpu.read(instruction.PC), cpu)
	instruction_log += instruction.InstructionText
	logger.Instructions.Print(instruction_log)
	instruction.action()
}

func (cpu *CPU) GetNextInstruction() *Instruction {
	logger := GetLogger()
	opcode := cpu.nextInstruction()
	var val byte
	var addr uint16
	var originalAddr uint16

	instruction_pc := cpu.pc - 1
	logger.MemoryDump.Printf("[PC: %04X]\n%v", instruction_pc, cpu.Dump())

	var instruction *Instruction

	switch opcode {
	case 0x00:
		instruction = NewInstruction(instruction_pc, "BRK", cpu.brk)

	case 0x40:
		instruction = NewInstruction(instruction_pc, "RTI", cpu.rti)

	case 0xea:
		instruction = NewInstruction(instruction_pc, "NOP", func() {})

	case 0xa9:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("LDA #$%02X", val), func() { cpu.lda(val) })
	case 0xa5:
		val, addr = cpu.nextValue(ZeroPage)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("LDA $%02X", addr), func() { cpu.lda(val) })
	case 0xb5:
		val, addr = cpu.nextValue(ZeroPageX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("LDA $%02X, X", addr), func() { cpu.lda(val) })
	case 0xad:
		val, addr = cpu.nextValue(Absolute)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("LDA $%04X", addr), func() { cpu.lda(val) })
	case 0xbd:
		val, addr = cpu.nextValue(AbsoluteX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("LDA $%04X, X", addr), func() { cpu.lda(val) })
	case 0xb9:
		val, addr = cpu.nextValue(AbsoluteY)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("LDA $%04X, Y", addr), func() { cpu.lda(val) })
	case 0xa1:
		val, addr = cpu.nextValue(IndirectX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("LDA ($%02X, X)", addr), func() { cpu.lda(val) })
	case 0xb1:
		val, addr = cpu.nextValue(IndirectY)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("LDA ($%02X), Y", addr), func() { cpu.lda(val) })

	case 0xa2:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("LDX #$%02X", val), func() { cpu.ldx(val) })
	case 0xa6:
		val, addr = cpu.nextValue(ZeroPage)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("LDX $%02X", addr), func() { cpu.ldx(val) })
	case 0xb6:
		val, addr = cpu.nextValue(ZeroPageY)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("LDX $%02X, Y", addr), func() { cpu.ldx(val) })
	case 0xae:
		val, addr = cpu.nextValue(Absolute)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("LDX $%04X", addr), func() { cpu.ldx(val) })
	case 0xbe:
		val, addr = cpu.nextValue(AbsoluteY)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("LDX $%04X, Y", addr), func() { cpu.ldx(val) })

	case 0xa0:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("LDY #$%02X", val), func() { cpu.ldy(val) })
	case 0xa4:
		val, addr = cpu.nextValue(ZeroPage)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("LDY $%02X", addr), func() { cpu.ldy(val) })
	case 0xb4:
		val, addr = cpu.nextValue(ZeroPageX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("LDY $%02X, X", addr), func() { cpu.ldy(val) })
	case 0xac:
		val, addr = cpu.nextValue(Absolute)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("LDY $%04X", addr), func() { cpu.ldy(val) })
	case 0xbc:
		val, addr = cpu.nextValue(AbsoluteX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("LDY $%04X, X", addr), func() { cpu.ldy(val) })

	case 0x85:
		addr, originalAddr = cpu.nextAddress(ZeroPage)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("STA $%02X", originalAddr), func() { cpu.sta(addr) })
	case 0x95:
		addr, originalAddr = cpu.nextAddress(ZeroPageX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("STA $%02X, X", originalAddr), func() { cpu.sta(addr) })
	case 0x8d:
		addr, originalAddr = cpu.nextAddress(Absolute)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("STA $%04X", originalAddr), func() { cpu.sta(addr) })
	case 0x9d:
		addr, originalAddr = cpu.nextAddress(AbsoluteX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("STA $%04X, X", originalAddr), func() { cpu.sta(addr) })
	case 0x99:
		addr, originalAddr = cpu.nextAddress(AbsoluteY)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("STA $%04X, Y", originalAddr), func() { cpu.sta(addr) })
	case 0x81:
		addr, originalAddr = cpu.nextAddress(IndirectX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("STA ($%02X, X)", originalAddr), func() { cpu.sta(addr) })
	case 0x91:
		addr, originalAddr = cpu.nextAddress(IndirectY)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("STA ($%02X), Y", originalAddr), func() { cpu.sta(addr) })

	case 0x86:
		addr, originalAddr = cpu.nextAddress(ZeroPage)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("STX $%02X", originalAddr), func() { cpu.stx(addr) })
	case 0x96:
		addr, originalAddr = cpu.nextAddress(ZeroPageY)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("STX $%02X, Y", originalAddr), func() { cpu.stx(addr) })
	case 0x8e:
		addr, originalAddr = cpu.nextAddress(Absolute)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("STX $%04X", originalAddr), func() { cpu.stx(addr) })

	case 0x84:
		addr, originalAddr = cpu.nextAddress(ZeroPage)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("STY $%02X", originalAddr), func() { cpu.sty(addr) })
	case 0x94:
		addr, originalAddr = cpu.nextAddress(ZeroPageX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("STY $%02X, X", originalAddr), func() { cpu.sty(addr) })
	case 0x8c:
		addr, originalAddr = cpu.nextAddress(Absolute)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("STY $%04X", originalAddr), func() { cpu.sty(addr) })

	case 0xaa:
		instruction = NewInstruction(instruction_pc, "TAX", cpu.tax)
	case 0xa8:
		instruction = NewInstruction(instruction_pc, "TAY", cpu.tay)
	case 0xba:
		instruction = NewInstruction(instruction_pc, "TSX", cpu.tsx)
	case 0x8a:
		instruction = NewInstruction(instruction_pc, "TXA", cpu.txa)
	case 0x9a:
		instruction = NewInstruction(instruction_pc, "TXS", cpu.txs)
	case 0x98:
		instruction = NewInstruction(instruction_pc, "TYA", cpu.tya)

	case 0x48:
		instruction = NewInstruction(instruction_pc, "PHA", cpu.pha)
	case 0x08:
		instruction = NewInstruction(instruction_pc, "PHP", cpu.php)
	case 0x68:
		instruction = NewInstruction(instruction_pc, "PLA", cpu.pla)
	case 0x28:
		instruction = NewInstruction(instruction_pc, "PLP", cpu.plp)

	case 0x29:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("AND #$%02X", val), func() { cpu.and(val) })
	case 0x25:
		val, addr = cpu.nextValue(ZeroPage)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("AND $%02X", addr), func() { cpu.and(val) })
	case 0x35:
		val, addr = cpu.nextValue(ZeroPageX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("AND $%02X, X", addr), func() { cpu.and(val) })
	case 0x2d:
		val, addr = cpu.nextValue(Absolute)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("AND $%04X", addr), func() { cpu.and(val) })
	case 0x3d:
		val, addr = cpu.nextValue(AbsoluteX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("AND $%04X, X", addr), func() { cpu.and(val) })
	case 0x39:
		val, addr = cpu.nextValue(AbsoluteY)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("AND $%04X, Y", addr), func() { cpu.and(val) })
	case 0x21:
		val, addr = cpu.nextValue(IndirectX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("AND ($%02X, X)", addr), func() { cpu.and(val) })
	case 0x31:
		val, addr = cpu.nextValue(IndirectY)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("AND ($%02X), Y", addr), func() { cpu.and(val) })

	case 0x49:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("EOR #$%02X", val), func() { cpu.eor(val) })
	case 0x45:
		val, addr = cpu.nextValue(ZeroPage)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("EOR $%02X", addr), func() { cpu.eor(val) })
	case 0x55:
		val, addr = cpu.nextValue(ZeroPageX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("EOR $%02X, X", addr), func() { cpu.eor(val) })
	case 0x4d:
		val, addr = cpu.nextValue(Absolute)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("EOR $%04X", addr), func() { cpu.eor(val) })
	case 0x5d:
		val, addr = cpu.nextValue(AbsoluteX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("EOR $%04X, X", addr), func() { cpu.eor(val) })
	case 0x59:
		val, addr = cpu.nextValue(AbsoluteY)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("EOR $%04X, Y", addr), func() { cpu.eor(val) })
	case 0x41:
		val, addr = cpu.nextValue(IndirectX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("EOR ($%02X, X)", addr), func() { cpu.eor(val) })
	case 0x51:
		val, addr = cpu.nextValue(IndirectY)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("EOR ($%02X), Y", addr), func() { cpu.eor(val) })

	case 0x09:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ORA #$%02X", val), func() { cpu.ora(val) })
	case 0x05:
		val, addr = cpu.nextValue(ZeroPage)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ORA $%02X", addr), func() { cpu.ora(val) })
	case 0x15:
		val, addr = cpu.nextValue(ZeroPageX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ORA $%02X, X", addr), func() { cpu.ora(val) })
	case 0x0d:
		val, addr = cpu.nextValue(Absolute)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ORA $%04X", addr), func() { cpu.ora(val) })
	case 0x1d:
		val, addr = cpu.nextValue(AbsoluteX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ORA $%04X, X", addr), func() { cpu.ora(val) })
	case 0x19:
		val, addr = cpu.nextValue(AbsoluteY)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ORA $%04X, Y", addr), func() { cpu.ora(val) })
	case 0x01:
		val, addr = cpu.nextValue(IndirectX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ORA ($%02X, X)", addr), func() { cpu.ora(val) })
	case 0x11:
		val, addr = cpu.nextValue(IndirectY)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ORA ($%02X), Y", addr), func() { cpu.ora(val) })

	case 0x24:
		val, addr = cpu.nextValue(ZeroPage)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("BIT $%02X", addr), func() { cpu.bit(val) })
	case 0x2c:
		val, addr = cpu.nextValue(Absolute)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("BIT $%04X", addr), func() { cpu.bit(val) })

	case 0x69:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ADC #$%02X", val), func() { cpu.adc(val) })
	case 0x65:
		val, addr = cpu.nextValue(ZeroPage)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ADC $%02X", addr), func() { cpu.adc(val) })
	case 0x75:
		val, addr = cpu.nextValue(ZeroPageX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ADC $%02X, X", addr), func() { cpu.adc(val) })
	case 0x6d:
		val, addr = cpu.nextValue(Absolute)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ADC $%04X", addr), func() { cpu.adc(val) })
	case 0x7d:
		val, addr = cpu.nextValue(AbsoluteX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ADC $%04X, X", addr), func() { cpu.adc(val) })
	case 0x79:
		val, addr = cpu.nextValue(AbsoluteY)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ADC $%04X, Y", addr), func() { cpu.adc(val) })
	case 0x61:
		val, addr = cpu.nextValue(IndirectX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ADC ($%02X, X)", addr), func() { cpu.adc(val) })
	case 0x71:
		val, addr = cpu.nextValue(IndirectY)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ADC ($%02X), Y", addr), func() { cpu.adc(val) })

	case 0xe9:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("SBC #$%02X", val), func() { cpu.sbc(val) })
	case 0xe5:
		val, addr = cpu.nextValue(ZeroPage)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("SBC $%02X", addr), func() { cpu.sbc(val) })
	case 0xf5:
		val, addr = cpu.nextValue(ZeroPageX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("SBC $%02X, X", addr), func() { cpu.sbc(val) })
	case 0xed:
		val, addr = cpu.nextValue(Absolute)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("SBC $%04X", addr), func() { cpu.sbc(val) })
	case 0xfd:
		val, addr = cpu.nextValue(AbsoluteX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("SBC $%04X, X", addr), func() { cpu.sbc(val) })
	case 0xf9:
		val, addr = cpu.nextValue(AbsoluteY)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("SBC $%04X, Y", addr), func() { cpu.sbc(val) })
	case 0xe1:
		val, addr = cpu.nextValue(IndirectX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("SBC ($%02X, X)", addr), func() { cpu.sbc(val) })
	case 0xf1:
		val, addr = cpu.nextValue(IndirectY)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("SBC ($%02X), Y", addr), func() { cpu.sbc(val) })

	case 0xc9:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("CMP #$%02X", val), func() { cpu.cmp(val) })
	case 0xc5:
		val, addr = cpu.nextValue(ZeroPage)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("CMP $%02X", addr), func() { cpu.cmp(val) })
	case 0xd5:
		val, addr = cpu.nextValue(ZeroPageX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("CMP $%02X, X", addr), func() { cpu.cmp(val) })
	case 0xcd:
		val, addr = cpu.nextValue(Absolute)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("CMP $%04X", addr), func() { cpu.cmp(val) })
	case 0xdd:
		val, addr = cpu.nextValue(AbsoluteX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("CMP $%04X, X", addr), func() { cpu.cmp(val) })
	case 0xd9:
		val, addr = cpu.nextValue(AbsoluteY)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("CMP $%04X, Y", addr), func() { cpu.cmp(val) })
	case 0xc1:
		val, addr = cpu.nextValue(IndirectX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("CMP ($%02X, X)", addr), func() { cpu.cmp(val) })
	case 0xd1:
		val, addr = cpu.nextValue(IndirectY)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("CMP ($%02X), Y", addr), func() { cpu.cmp(val) })

	case 0xe0:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("CPX #$%02X", val), func() { cpu.cpx(val) })
	case 0xe4:
		val, addr = cpu.nextValue(ZeroPage)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("CPX $%02X", addr), func() { cpu.cpx(val) })
	case 0xec:
		val, addr = cpu.nextValue(Absolute)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("CPX $%04X", addr), func() { cpu.cpx(val) })

	case 0xc0:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("CPY #$%02X", val), func() { cpu.cpy(val) })
	case 0xc4:
		val, addr = cpu.nextValue(ZeroPage)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("CPY $%02X", addr), func() { cpu.cpy(val) })
	case 0xcc:
		val, addr = cpu.nextValue(Absolute)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("CPY $%04X", addr), func() { cpu.cpy(val) })

	case 0xe6:
		addr, originalAddr = cpu.nextAddress(ZeroPage)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("INC $%02X", originalAddr), func() { cpu.inc(addr) })
	case 0xf6:
		addr, originalAddr = cpu.nextAddress(ZeroPageX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("INC $%02X, X", originalAddr), func() { cpu.inc(addr) })
	case 0xee:
		addr, originalAddr = cpu.nextAddress(Absolute)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("INC $%04X", originalAddr), func() { cpu.inc(addr) })
	case 0xfe:
		addr, originalAddr = cpu.nextAddress(AbsoluteX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("INC $%04X, X", originalAddr), func() { cpu.inc(addr) })

	case 0xc6:
		addr, originalAddr = cpu.nextAddress(ZeroPage)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("DEC $%02X", originalAddr), func() { cpu.dec(addr) })
	case 0xd6:
		addr, originalAddr = cpu.nextAddress(ZeroPageX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("DEC $%02X, X", originalAddr), func() { cpu.dec(addr) })
	case 0xce:
		addr, originalAddr = cpu.nextAddress(Absolute)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("DEC $%04X", originalAddr), func() { cpu.dec(addr) })
	case 0xde:
		addr, originalAddr = cpu.nextAddress(AbsoluteX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("DEC $%04X, X", originalAddr), func() { cpu.dec(addr) })

	case 0xe8:
		instruction = NewInstruction(instruction_pc, "INX", cpu.inx)
	case 0xc8:
		instruction = NewInstruction(instruction_pc, "INY", cpu.iny)

	case 0xca:
		instruction = NewInstruction(instruction_pc, "DEX", cpu.dex)
	case 0x88:
		instruction = NewInstruction(instruction_pc, "DEY", cpu.dey)

	case 0x0a:
		instruction = NewInstruction(instruction_pc, "ASL A", cpu.asl_acc)
	case 0x06:
		addr, originalAddr = cpu.nextAddress(ZeroPage)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ASL $%02X", originalAddr), func() { cpu.asl(addr) })
	case 0x16:
		addr, originalAddr = cpu.nextAddress(ZeroPageX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ASL $%02X, X", originalAddr), func() { cpu.asl(addr) })
	case 0x0e:
		addr, originalAddr = cpu.nextAddress(Absolute)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ASL $%04X", originalAddr), func() { cpu.asl(addr) })
	case 0x1e:
		addr, originalAddr = cpu.nextAddress(AbsoluteX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ASL $%04X, X", originalAddr), func() { cpu.asl(addr) })

	case 0x4a:
		instruction = NewInstruction(instruction_pc, "LSR A", cpu.lsr_acc)
	case 0x46:
		addr, originalAddr = cpu.nextAddress(ZeroPage)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("LSR $%02X", originalAddr), func() { cpu.lsr(addr) })
	case 0x56:
		addr, originalAddr = cpu.nextAddress(ZeroPageX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("LSR $%02X, X", originalAddr), func() { cpu.lsr(addr) })
	case 0x4e:
		addr, originalAddr = cpu.nextAddress(Absolute)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("LSR $%04X", originalAddr), func() { cpu.lsr(addr) })
	case 0x5e:
		addr, originalAddr = cpu.nextAddress(AbsoluteX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("LSR $%04X, X", originalAddr), func() { cpu.lsr(addr) })

	case 0x2a:
		instruction = NewInstruction(instruction_pc, "ROL A", cpu.rol_acc)
	case 0x26:
		addr, originalAddr = cpu.nextAddress(ZeroPage)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ROL $%02X", originalAddr), func() { cpu.rol(addr) })
	case 0x36:
		addr, originalAddr = cpu.nextAddress(ZeroPageX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ROL $%02X, X", originalAddr), func() { cpu.rol(addr) })
	case 0x2e:
		addr, originalAddr = cpu.nextAddress(Absolute)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ROL $%04X", originalAddr), func() { cpu.rol(addr) })
	case 0x3e:
		addr, originalAddr = cpu.nextAddress(AbsoluteX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ROL $%04X, X", originalAddr), func() { cpu.rol(addr) })

	case 0x6a:
		instruction = NewInstruction(instruction_pc, "ROR A", cpu.ror_acc)
	case 0x66:
		addr, originalAddr = cpu.nextAddress(ZeroPage)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ROR $%02X", originalAddr), func() { cpu.ror(addr) })
	case 0x76:
		addr, originalAddr = cpu.nextAddress(ZeroPageX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ROR $%02X, X", originalAddr), func() { cpu.ror(addr) })
	case 0x6e:
		addr, originalAddr = cpu.nextAddress(Absolute)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ROR $%04X", originalAddr), func() { cpu.ror(addr) })
	case 0x7e:
		addr, originalAddr = cpu.nextAddress(AbsoluteX)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("ROR $%04X, X", originalAddr), func() { cpu.ror(addr) })

	case 0x4c:
		addr, originalAddr = cpu.nextAddress(Absolute)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("JMP $%04X", originalAddr), func() { cpu.jmp(addr) })
	case 0x6c:
		addr, originalAddr = cpu.nextAddress(Indirect)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("JMP ($%04X)", originalAddr), func() { cpu.jmp(addr) })

	case 0x20:
		addr, originalAddr = cpu.nextAddress(Absolute)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("JSR $%04X", originalAddr), func() { cpu.jsr(addr) })

	case 0x60:
		instruction = NewInstruction(instruction_pc, "RTS", cpu.rts)

	case 0x90:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("BCC #$%02X", val), func() { cpu.bcc(val) })
	case 0xb0:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("BCS #$%02X", val), func() { cpu.bcs(val) })
	case 0xf0:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("BEQ #$%02X", val), func() { cpu.beq(val) })
	case 0x30:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("BMI #$%02X", val), func() { cpu.bmi(val) })
	case 0xd0:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("BNE #$%02X", val), func() { cpu.bne(val) })
	case 0x10:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("BPL #$%02X", val), func() { cpu.bpl(val) })
	case 0x50:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("BVC #$%02X", val), func() { cpu.bvc(val) })
	case 0x70:
		val, _ = cpu.nextValue(Immediate)
		instruction = NewInstruction(instruction_pc, fmt.Sprintf("BVS #$%02X", val), func() { cpu.bvs(val) })

	case 0x18:
		instruction = NewInstruction(instruction_pc, "CLC", cpu.clc)
	case 0xd8:
		instruction = NewInstruction(instruction_pc, "CLD", cpu.cld)
	case 0x58:
		instruction = NewInstruction(instruction_pc, "CLI", cpu.cli)
	case 0xb8:
		instruction = NewInstruction(instruction_pc, "CLV", cpu.clv)
	case 0x38:
		instruction = NewInstruction(instruction_pc, "SEC", cpu.sec)
	case 0xf8:
		instruction = NewInstruction(instruction_pc, "SED", cpu.sed)
	case 0x78:
		instruction = NewInstruction(instruction_pc, "SEI", cpu.sei)

	// unofficial opcodes
	case 0x1a, 0x3a, 0x5a, 0x7a, 0xda, 0xfa:
		instruction = NewInstruction(instruction_pc, "*NOP", func() {})
	case 0x04, 0x44, 0x64, 0x14, 0x34, 0x54, 0x74, 0xd4, 0xf4, 0x80:
		instruction = NewInstruction(instruction_pc, "*NOP", func() {})
		cpu.pc++
	case 0x0c, 0x1c, 0x3c, 0x5c, 0x7c, 0xdc, 0xfc:
		instruction = NewInstruction(instruction_pc, "*NOP", func() {})
		cpu.pc += 2

	default:
		instruction = NewInstruction(instruction_pc, "UNKNOWN", func() {})
	}

	return instruction
}

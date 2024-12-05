package mos6502

import "fmt"

func (cpu *CPU) Step() {
	opcode := cpu.nextInstruction()
	var val byte
	var addr uint16
	var originalAddr uint16

	logger := GetLogger()
	instruction_log := fmt.Sprintf("[PC: %04X] OPCODE %02X | %s | ", cpu.pc-1, opcode, cpu)
	logger.MemoryDump.Printf("[PC: %04X]\n%v", cpu.pc-1, cpu.Dump())

	switch opcode {
	case 0x00:
		cpu.brk()
		instruction_log += "BRK"

	case 0x40:
		cpu.rti()
		instruction_log += "RTI"

	case 0xea:
		instruction_log += fmt.Sprintf("NOP")

	case 0xa9:
		val, _ = cpu.nextValue(Immediate)
		cpu.lda(val)
		instruction_log += fmt.Sprintf("LDA #$%02X", val)
	case 0xa5:
		val, addr = cpu.nextValue(ZeroPage)
		cpu.lda(val)
		instruction_log += fmt.Sprintf("LDA $%02X", addr)
	case 0xb5:
		val, addr = cpu.nextValue(ZeroPageX)
		cpu.lda(val)
		instruction_log += fmt.Sprintf("LDA $%02X, X", addr)
	case 0xad:
		val, addr = cpu.nextValue(Absolute)
		cpu.lda(val)
		instruction_log += fmt.Sprintf("LDA $%04X", addr)
	case 0xbd:
		val, addr = cpu.nextValue(AbsoluteX)
		cpu.lda(val)
		instruction_log += fmt.Sprintf("LDA $%04X, X", addr)
	case 0xb9:
		val, addr = cpu.nextValue(AbsoluteY)
		cpu.lda(val)
		instruction_log += fmt.Sprintf("LDA $%04X, Y", addr)
	case 0xa1:
		val, addr = cpu.nextValue(IndirectX)
		cpu.lda(val)
		instruction_log += fmt.Sprintf("LDA ($%02X, X)", addr)
	case 0xb1:
		val, addr = cpu.nextValue(IndirectY)
		cpu.lda(val)
		instruction_log += fmt.Sprintf("LDA ($%02X), Y", addr)

	case 0xa2:
		val, _ = cpu.nextValue(Immediate)
		cpu.ldx(val)
		instruction_log += fmt.Sprintf("LDX #$%02X", val)
	case 0xa6:
		val, addr = cpu.nextValue(ZeroPage)
		cpu.ldx(val)
		instruction_log += fmt.Sprintf("LDX $%02X", addr)
	case 0xb6:
		val, addr = cpu.nextValue(ZeroPageY)
		cpu.ldx(val)
		instruction_log += fmt.Sprintf("LDX $%02X, Y", addr)
	case 0xae:
		val, addr = cpu.nextValue(Absolute)
		cpu.ldx(val)
		instruction_log += fmt.Sprintf("LDX $%04X", addr)
	case 0xbe:
		val, addr = cpu.nextValue(AbsoluteY)
		cpu.ldx(val)
		instruction_log += fmt.Sprintf("LDX $%04X, Y", addr)

	case 0xa0:
		val, _ = cpu.nextValue(Immediate)
		cpu.ldy(val)
		instruction_log += fmt.Sprintf("LDY #$%02X", val)
	case 0xa4:
		val, addr = cpu.nextValue(ZeroPage)
		cpu.ldy(val)
		instruction_log += fmt.Sprintf("LDY $%02X", addr)
	case 0xb4:
		val, addr = cpu.nextValue(ZeroPageX)
		cpu.ldy(val)
		instruction_log += fmt.Sprintf("LDY $%02X, X", addr)
	case 0xac:
		val, addr = cpu.nextValue(Absolute)
		cpu.ldy(val)
		instruction_log += fmt.Sprintf("LDY $%04X", addr)
	case 0xbc:
		val, addr = cpu.nextValue(AbsoluteX)
		cpu.ldy(val)
		instruction_log += fmt.Sprintf("LDY $%04X, X", addr)

	case 0x85:
		addr, originalAddr = cpu.nextAddress(ZeroPage)
		cpu.sta(addr)
		instruction_log += fmt.Sprintf("STA $%02X", originalAddr)
	case 0x95:
		addr, originalAddr = cpu.nextAddress(ZeroPageX)
		cpu.sta(addr)
		instruction_log += fmt.Sprintf("STA $%02X, X", originalAddr)
	case 0x8d:
		addr, originalAddr = cpu.nextAddress(Absolute)
		cpu.sta(addr)
		instruction_log += fmt.Sprintf("STA $%04X", originalAddr)
	case 0x9d:
		addr, originalAddr = cpu.nextAddress(AbsoluteX)
		cpu.sta(addr)
		instruction_log += fmt.Sprintf("STA $%04X, X", originalAddr)
	case 0x99:
		addr, originalAddr = cpu.nextAddress(AbsoluteY)
		cpu.sta(addr)
		instruction_log += fmt.Sprintf("STA $%04X, Y", originalAddr)
	case 0x81:
		addr, originalAddr = cpu.nextAddress(IndirectX)
		cpu.sta(addr)
		instruction_log += fmt.Sprintf("STA ($%02X, X)", originalAddr)
	case 0x91:
		addr, originalAddr = cpu.nextAddress(IndirectY)
		cpu.sta(addr)
		instruction_log += fmt.Sprintf("STA ($%02X), Y", originalAddr)

	case 0x86:
		addr, originalAddr = cpu.nextAddress(ZeroPage)
		cpu.stx(addr)
		instruction_log += fmt.Sprintf("STX $%02X", originalAddr)
	case 0x96:
		addr, originalAddr = cpu.nextAddress(ZeroPageY)
		cpu.stx(addr)
		instruction_log += fmt.Sprintf("STX $%02X, Y", originalAddr)
	case 0x8e:
		addr, originalAddr = cpu.nextAddress(Absolute)
		cpu.stx(addr)
		instruction_log += fmt.Sprintf("STX $%04X", originalAddr)

	case 0x84:
		addr, originalAddr = cpu.nextAddress(ZeroPage)
		cpu.sty(addr)
		instruction_log += fmt.Sprintf("STY $%02X", originalAddr)
	case 0x94:
		addr, originalAddr = cpu.nextAddress(ZeroPageX)
		cpu.sty(addr)
		instruction_log += fmt.Sprintf("STY $%02X, X", originalAddr)
	case 0x8c:
		addr, originalAddr = cpu.nextAddress(Absolute)
		cpu.sty(addr)
		instruction_log += fmt.Sprintf("STY $%04X", originalAddr)

	case 0xaa:
		cpu.tax()
		instruction_log += "TAX"
	case 0xa8:
		cpu.tay()
		instruction_log += "TAY"
	case 0xba:
		cpu.tsx()
		instruction_log += "TSX"
	case 0x8a:
		cpu.txa()
		instruction_log += "TXA"
	case 0x9a:
		cpu.txs()
		instruction_log += "TXS"
	case 0x98:
		cpu.tya()
		instruction_log += "TYA"

	case 0x48:
		cpu.pha()
		instruction_log += "PHA"
	case 0x08:
		cpu.php()
		instruction_log += "PHP"
	case 0x68:
		cpu.pla()
		instruction_log += "PLA"
	case 0x28:
		cpu.plp()
		instruction_log += "PLP"

	case 0x29:
		val, _ = cpu.nextValue(Immediate)
		cpu.and(val)
		instruction_log += fmt.Sprintf("AND #$%02X", val)
	case 0x25:
		val, addr = cpu.nextValue(ZeroPage)
		cpu.and(val)
		instruction_log += fmt.Sprintf("AND $%02X", addr)
	case 0x35:
		val, addr = cpu.nextValue(ZeroPageX)
		cpu.and(val)
		instruction_log += fmt.Sprintf("AND $%02X, X", addr)
	case 0x2d:
		val, addr = cpu.nextValue(Absolute)
		cpu.and(val)
		instruction_log += fmt.Sprintf("AND $%04X", addr)
	case 0x3d:
		val, addr = cpu.nextValue(AbsoluteX)
		cpu.and(val)
		instruction_log += fmt.Sprintf("AND $%04X, X", addr)
	case 0x39:
		val, addr = cpu.nextValue(AbsoluteY)
		cpu.and(val)
		instruction_log += fmt.Sprintf("AND $%04X, Y", addr)
	case 0x21:
		val, addr = cpu.nextValue(IndirectX)
		cpu.and(val)
		instruction_log += fmt.Sprintf("AND ($%02X, X)", addr)
	case 0x31:
		val, addr = cpu.nextValue(IndirectY)
		cpu.and(val)
		instruction_log += fmt.Sprintf("AND ($%02X), Y", addr)

	case 0x49:
		val, _ = cpu.nextValue(Immediate)
		cpu.eor(val)
		instruction_log += fmt.Sprintf("EOR #$%02X", val)
	case 0x45:
		val, addr = cpu.nextValue(ZeroPage)
		cpu.eor(val)
		instruction_log += fmt.Sprintf("EOR $%02X", addr)
	case 0x55:
		val, addr = cpu.nextValue(ZeroPageX)
		cpu.eor(val)
		instruction_log += fmt.Sprintf("EOR $%02X, X", addr)
	case 0x4d:
		val, addr = cpu.nextValue(Absolute)
		cpu.eor(val)
		instruction_log += fmt.Sprintf("EOR $%04X", addr)
	case 0x5d:
		val, addr = cpu.nextValue(AbsoluteX)
		cpu.eor(val)
		instruction_log += fmt.Sprintf("EOR $%04X, X", addr)
	case 0x59:
		val, addr = cpu.nextValue(AbsoluteY)
		cpu.eor(val)
		instruction_log += fmt.Sprintf("EOR $%04X, Y", addr)
	case 0x41:
		val, addr = cpu.nextValue(IndirectX)
		cpu.eor(val)
		instruction_log += fmt.Sprintf("EOR ($%02X, X)", addr)
	case 0x51:
		val, addr = cpu.nextValue(IndirectY)
		cpu.eor(val)
		instruction_log += fmt.Sprintf("EOR ($%02X), Y", addr)

	case 0x09:
		val, _ = cpu.nextValue(Immediate)
		cpu.ora(val)
		instruction_log += fmt.Sprintf("ORA #$%02X", val)
	case 0x05:
		val, addr = cpu.nextValue(ZeroPage)
		cpu.ora(val)
		instruction_log += fmt.Sprintf("ORA $%02X", addr)
	case 0x15:
		val, addr = cpu.nextValue(ZeroPageX)
		cpu.ora(val)
		instruction_log += fmt.Sprintf("ORA $%02X, X", addr)
	case 0x0d:
		val, addr = cpu.nextValue(Absolute)
		cpu.ora(val)
		instruction_log += fmt.Sprintf("ORA $%04X", addr)
	case 0x1d:
		val, addr = cpu.nextValue(AbsoluteX)
		cpu.ora(val)
		instruction_log += fmt.Sprintf("ORA $%04X, X", addr)
	case 0x19:
		val, addr = cpu.nextValue(AbsoluteY)
		cpu.ora(val)
		instruction_log += fmt.Sprintf("ORA $%04X, Y", addr)
	case 0x01:
		val, addr = cpu.nextValue(IndirectX)
		cpu.ora(val)
		instruction_log += fmt.Sprintf("ORA ($%02X, X)", addr)
	case 0x11:
		val, addr = cpu.nextValue(IndirectY)
		cpu.ora(val)
		instruction_log += fmt.Sprintf("ORA ($%02X), Y", addr)

	case 0x24:
		val, addr = cpu.nextValue(ZeroPage)
		cpu.bit(val)
		instruction_log += fmt.Sprintf("BIT $%02X", addr)
	case 0x2c:
		val, addr = cpu.nextValue(Absolute)
		cpu.bit(val)
		instruction_log += fmt.Sprintf("BIT $%04X", addr)

	case 0x69:
		val, _ = cpu.nextValue(Immediate)
		cpu.adc(val)
		instruction_log += fmt.Sprintf("ADC #$%02X", val)
	case 0x65:
		val, addr = cpu.nextValue(ZeroPage)
		cpu.adc(val)
		instruction_log += fmt.Sprintf("ADC $%02X", addr)
	case 0x75:
		val, addr = cpu.nextValue(ZeroPageX)
		cpu.adc(val)
		instruction_log += fmt.Sprintf("ADC $%02X, X", addr)
	case 0x6d:
		val, addr = cpu.nextValue(Absolute)
		cpu.adc(val)
		instruction_log += fmt.Sprintf("ADC $%04X", addr)
	case 0x7d:
		val, addr = cpu.nextValue(AbsoluteX)
		cpu.adc(val)
		instruction_log += fmt.Sprintf("ADC $%04X, X", addr)
	case 0x79:
		val, addr = cpu.nextValue(AbsoluteY)
		cpu.adc(val)
		instruction_log += fmt.Sprintf("ADC $%04X, Y", addr)
	case 0x61:
		val, addr = cpu.nextValue(IndirectX)
		cpu.adc(val)
		instruction_log += fmt.Sprintf("ADC ($%02X, X)", addr)
	case 0x71:
		val, addr = cpu.nextValue(IndirectY)
		cpu.adc(val)
		instruction_log += fmt.Sprintf("ADC ($%02X), Y", addr)

	case 0xe9:
		val, _ = cpu.nextValue(Immediate)
		cpu.sbc(val)
		instruction_log += fmt.Sprintf("SBC #$%02X", val)
	case 0xe5:
		val, addr = cpu.nextValue(ZeroPage)
		cpu.sbc(val)
		instruction_log += fmt.Sprintf("SBC $%02X", addr)
	case 0xf5:
		val, addr = cpu.nextValue(ZeroPageX)
		cpu.sbc(val)
		instruction_log += fmt.Sprintf("SBC $%02X, X", addr)
	case 0xed:
		val, addr = cpu.nextValue(Absolute)
		cpu.sbc(val)
		instruction_log += fmt.Sprintf("SBC $%04X", addr)
	case 0xfd:
		val, addr = cpu.nextValue(AbsoluteX)
		cpu.sbc(val)
		instruction_log += fmt.Sprintf("SBC $%04X, X", addr)
	case 0xf9:
		val, addr = cpu.nextValue(AbsoluteY)
		cpu.sbc(val)
		instruction_log += fmt.Sprintf("SBC $%04X, Y", addr)
	case 0xe1:
		val, addr = cpu.nextValue(IndirectX)
		cpu.sbc(val)
		instruction_log += fmt.Sprintf("SBC ($%02X, X)", addr)
	case 0xf1:
		val, addr = cpu.nextValue(IndirectY)
		cpu.sbc(val)
		instruction_log += fmt.Sprintf("SBC ($%02X), Y", addr)

	case 0xc9:
		val, _ = cpu.nextValue(Immediate)
		cpu.cmp(val)
		instruction_log += fmt.Sprintf("CMP #$%02X", val)
	case 0xc5:
		val, addr = cpu.nextValue(ZeroPage)
		cpu.cmp(val)
		instruction_log += fmt.Sprintf("CMP $%02X", addr)
	case 0xd5:
		val, addr = cpu.nextValue(ZeroPageX)
		cpu.cmp(val)
		instruction_log += fmt.Sprintf("CMP $%02X, X", addr)
	case 0xcd:
		val, addr = cpu.nextValue(Absolute)
		cpu.cmp(val)
		instruction_log += fmt.Sprintf("CMP $%04X", addr)
	case 0xdd:
		val, addr = cpu.nextValue(AbsoluteX)
		cpu.cmp(val)
		instruction_log += fmt.Sprintf("CMP $%04X, X", addr)
	case 0xd9:
		val, addr = cpu.nextValue(AbsoluteY)
		cpu.cmp(val)
		instruction_log += fmt.Sprintf("CMP $%04X, Y", addr)
	case 0xc1:
		val, addr = cpu.nextValue(IndirectX)
		cpu.cmp(val)
		instruction_log += fmt.Sprintf("CMP ($%02X, X)", addr)
	case 0xd1:
		val, addr = cpu.nextValue(IndirectY)
		cpu.cmp(val)
		instruction_log += fmt.Sprintf("CMP ($%02X), Y", addr)

	case 0xe0:
		val, _ = cpu.nextValue(Immediate)
		cpu.cpx(val)
		instruction_log += fmt.Sprintf("CPX #$%02X", val)
	case 0xe4:
		val, addr = cpu.nextValue(ZeroPage)
		cpu.cpx(val)
		instruction_log += fmt.Sprintf("CPX $%02X", addr)
	case 0xec:
		val, addr = cpu.nextValue(Absolute)
		cpu.cpx(val)
		instruction_log += fmt.Sprintf("CPX $%04X", addr)

	case 0xc0:
		val, _ = cpu.nextValue(Immediate)
		cpu.cpy(val)
		instruction_log += fmt.Sprintf("CPY #$%02X", val)
	case 0xc4:
		val, addr = cpu.nextValue(ZeroPage)
		cpu.cpy(val)
		instruction_log += fmt.Sprintf("CPY $%02X", addr)
	case 0xcc:
		val, addr = cpu.nextValue(Absolute)
		cpu.cpy(val)
		instruction_log += fmt.Sprintf("CPY $%04X", addr)

	case 0xe6:
		addr, originalAddr = cpu.nextAddress(ZeroPage)
		cpu.inc(addr)
		instruction_log += fmt.Sprintf("INC $%02X", originalAddr)
	case 0xf6:
		addr, originalAddr = cpu.nextAddress(ZeroPageX)
		cpu.inc(addr)
		instruction_log += fmt.Sprintf("INC $%02X, X", originalAddr)
	case 0xee:
		addr, originalAddr = cpu.nextAddress(Absolute)
		cpu.inc(addr)
		instruction_log += fmt.Sprintf("INC $%04X", originalAddr)
	case 0xfe:
		addr, originalAddr = cpu.nextAddress(AbsoluteX)
		cpu.inc(addr)
		instruction_log += fmt.Sprintf("INC $%04X, X", originalAddr)

	case 0xc6:
		addr, originalAddr = cpu.nextAddress(ZeroPage)
		cpu.dec(addr)
		instruction_log += fmt.Sprintf("DEC $%02X", originalAddr)
	case 0xd6:
		addr, originalAddr = cpu.nextAddress(ZeroPageX)
		cpu.dec(addr)
		instruction_log += fmt.Sprintf("DEC $%02X, X", originalAddr)
	case 0xce:
		addr, originalAddr = cpu.nextAddress(Absolute)
		cpu.dec(addr)
		instruction_log += fmt.Sprintf("DEC $%04X", originalAddr)
	case 0xde:
		addr, originalAddr = cpu.nextAddress(AbsoluteX)
		cpu.dec(addr)
		instruction_log += fmt.Sprintf("DEC $%04X, X", originalAddr)

	case 0xe8:
		cpu.inx()
		instruction_log += "INX"
	case 0xc8:
		cpu.iny()
		instruction_log += "INY"

	case 0xca:
		cpu.dex()
		instruction_log += "DEX"
	case 0x88:
		cpu.dey()
		instruction_log += "DEY"

	case 0x0a:
		cpu.asl_acc()
		instruction_log += "ASL A"
	case 0x06:
		addr, originalAddr = cpu.nextAddress(ZeroPage)
		cpu.asl(addr)
		instruction_log += fmt.Sprintf("ASL $%02X", originalAddr)
	case 0x16:
		addr, originalAddr = cpu.nextAddress(ZeroPageX)
		cpu.asl(addr)
		instruction_log += fmt.Sprintf("ASL $%02X, X", originalAddr)
	case 0x0e:
		addr, originalAddr = cpu.nextAddress(Absolute)
		cpu.asl(addr)
		instruction_log += fmt.Sprintf("ASL $%04X", originalAddr)
	case 0x1e:
		addr, originalAddr = cpu.nextAddress(AbsoluteX)
		cpu.asl(addr)
		instruction_log += fmt.Sprintf("ASL $%04X, X", originalAddr)

	case 0x4a:
		cpu.lsr_acc()
		instruction_log += "LSR A"
	case 0x46:
		addr, originalAddr = cpu.nextAddress(ZeroPage)
		cpu.lsr(addr)
		instruction_log += fmt.Sprintf("LSR $%02X", originalAddr)
	case 0x56:
		addr, originalAddr = cpu.nextAddress(ZeroPageX)
		cpu.lsr(addr)
		instruction_log += fmt.Sprintf("LSR $%02X, X", originalAddr)
	case 0x4e:
		addr, originalAddr = cpu.nextAddress(Absolute)
		cpu.lsr(addr)
		instruction_log += fmt.Sprintf("LSR $%04X", originalAddr)
	case 0x5e:
		addr, originalAddr = cpu.nextAddress(AbsoluteX)
		cpu.lsr(addr)
		instruction_log += fmt.Sprintf("LSR $%04X, X", originalAddr)

	case 0x2a:
		cpu.rol_acc()
		instruction_log += "ROL A"
	case 0x26:
		addr, originalAddr = cpu.nextAddress(ZeroPage)
		cpu.rol(addr)
		instruction_log += fmt.Sprintf("ROL $%02X", originalAddr)
	case 0x36:
		addr, originalAddr = cpu.nextAddress(ZeroPageX)
		cpu.rol(addr)
		instruction_log += fmt.Sprintf("ROL $%02X, X", originalAddr)
	case 0x2e:
		addr, originalAddr = cpu.nextAddress(Absolute)
		cpu.rol(addr)
		instruction_log += fmt.Sprintf("ROL $%04X", originalAddr)
	case 0x3e:
		addr, originalAddr = cpu.nextAddress(AbsoluteX)
		cpu.rol(addr)
		instruction_log += fmt.Sprintf("ROL $%04X, X", originalAddr)

	case 0x6a:
		cpu.ror_acc()
		instruction_log += "ROR A"
	case 0x66:
		addr, originalAddr = cpu.nextAddress(ZeroPage)
		cpu.ror(addr)
		instruction_log += fmt.Sprintf("ROR $%02X", originalAddr)
	case 0x76:
		addr, originalAddr = cpu.nextAddress(ZeroPageX)
		cpu.ror(addr)
		instruction_log += fmt.Sprintf("ROR $%02X, X", originalAddr)
	case 0x6e:
		addr, originalAddr = cpu.nextAddress(Absolute)
		cpu.ror(addr)
		instruction_log += fmt.Sprintf("ROR $%04X", originalAddr)
	case 0x7e:
		addr, originalAddr = cpu.nextAddress(AbsoluteX)
		cpu.ror(addr)
		instruction_log += fmt.Sprintf("ROR $%04X, X", originalAddr)

	case 0x4c:
		addr, originalAddr = cpu.nextAddress(Absolute)
		cpu.jmp(addr)
		instruction_log += fmt.Sprintf("JMP $%04X", originalAddr)
	case 0x6c:
		addr, originalAddr = cpu.nextAddress(Indirect)
		cpu.jmp(addr)
		instruction_log += fmt.Sprintf("JMP ($%04X)", originalAddr)

	case 0x20:
		addr, originalAddr = cpu.nextAddress(Absolute)
		cpu.jsr(addr)
		instruction_log += fmt.Sprintf("JSR $%04X", originalAddr)

	case 0x60:
		cpu.rts()
		instruction_log += "RTS"

	case 0x90:
		val, _ = cpu.nextValue(Immediate)
		cpu.bcc(val)
		instruction_log += fmt.Sprintf("BCC #$%02X", val)
	case 0xb0:
		val, _ = cpu.nextValue(Immediate)
		cpu.bcs(val)
		instruction_log += fmt.Sprintf("BCS #$%02X", val)
	case 0xf0:
		val, _ = cpu.nextValue(Immediate)
		cpu.beq(val)
		instruction_log += fmt.Sprintf("BEQ #$%02X", val)
	case 0x30:
		val, _ = cpu.nextValue(Immediate)
		cpu.bmi(val)
		instruction_log += fmt.Sprintf("BMI #$%02X", val)
	case 0xd0:
		val, _ = cpu.nextValue(Immediate)
		cpu.bne(val)
		instruction_log += fmt.Sprintf("BNE #$%02X", val)
	case 0x10:
		val, _ = cpu.nextValue(Immediate)
		cpu.bpl(val)
		instruction_log += fmt.Sprintf("BPL #$%02X", val)
	case 0x50:
		val, _ = cpu.nextValue(Immediate)
		cpu.bvc(val)
		instruction_log += fmt.Sprintf("BVC #$%02X", val)
	case 0x70:
		val, _ = cpu.nextValue(Immediate)
		cpu.bvs(val)
		instruction_log += fmt.Sprintf("BVS #$%02X", val)

	case 0x18:
		cpu.clc()
		instruction_log += "CLC"
	case 0xd8:
		cpu.cld()
		instruction_log += "CLD"
	case 0x58:
		cpu.cli()
		instruction_log += "CLI"
	case 0xb8:
		cpu.clv()
		instruction_log += "CLV"
	case 0x38:
		cpu.sec()
		instruction_log += "SEC"
	case 0xf8:
		cpu.sed()
		instruction_log += "SED"
	case 0x78:
		cpu.sei()
		instruction_log += "SEI"

	// unofficial opcodes
	case 0x1a, 0x3a, 0x5a, 0x7a, 0xda, 0xfa:
		instruction_log += "*NOP"
	case 0x04, 0x44, 0x64, 0x14, 0x34, 0x54, 0x74, 0xd4, 0xf4, 0x80:
		cpu.pc++
		instruction_log += "*NOP"
	case 0x0c, 0x1c, 0x3c, 0x5c, 0x7c, 0xdc, 0xfc:
		cpu.pc += 2
		instruction_log += "*NOP"

	default:
		instruction_log += fmt.Sprintf("UNKNOWN")
	}

	logger.Instructions.Print(instruction_log)
}

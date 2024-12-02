package mos6502

import "fmt"

func (cpu *CPU) Step() {
	opcode := cpu.nextInstruction()
	var val byte
	var addr uint16
	fmt.Printf("[PC: %04X] OPCODE %02X | %s | ", cpu.pc-1, opcode, cpu)

	switch opcode {
	case 0x00:
		cpu.brk()
		fmt.Print("BRK")

	case 0x40:
		cpu.rti()
		fmt.Print("RTI")

	case 0xea:
		fmt.Printf("NOP")

	case 0xa9:
		val, _ = cpu.nextValue(Immediate)
		cpu.lda(val)
		fmt.Printf("LDA #$%02X", val)
	case 0xa5:
		val, addr = cpu.nextValue(ZeroPage)
		cpu.lda(val)
		fmt.Printf("LDA $%02X", addr)
	case 0xb5:
		val, addr = cpu.nextValue(ZeroPageX)
		cpu.lda(val)
		fmt.Printf("LDA $%02X, X", addr)
	case 0xad:
		val, addr = cpu.nextValue(Absolute)
		cpu.lda(val)
		fmt.Printf("LDA $%04X", addr)
	case 0xbd:
		val, addr = cpu.nextValue(AbsoluteX)
		cpu.lda(val)
		fmt.Printf("LDA $%04X, X", addr)
	case 0xb9:
		val, addr = cpu.nextValue(AbsoluteY)
		cpu.lda(val)
		fmt.Printf("LDA $%04X, Y", addr)
	case 0xa1:
		val, addr = cpu.nextValue(IndirectX)
		cpu.lda(val)
		fmt.Printf("LDA ($%02X, X)", addr)
	case 0xb1:
		val, addr = cpu.nextValue(IndirectY)
		cpu.lda(val)
		fmt.Printf("LDA ($%04X), Y", addr)

	case 0xa2:
		val, _ = cpu.nextValue(Immediate)
		cpu.ldx(val)
		fmt.Printf("LDX #$%02X", val)
	case 0xa6:
		val, addr = cpu.nextValue(ZeroPage)
		cpu.ldx(val)
		fmt.Printf("LDX $%02X", addr)
	case 0xb6:
		val, addr = cpu.nextValue(ZeroPageY)
		cpu.ldx(val)
		fmt.Printf("LDX $%02X, Y", addr)
	case 0xae:
		val, addr = cpu.nextValue(Absolute)
		cpu.ldx(val)
		fmt.Printf("LDX $%04X", addr)
	case 0xbe:
		val, addr = cpu.nextValue(AbsoluteY)
		cpu.ldx(val)
		fmt.Printf("LDX $%04X, Y", addr)

	case 0xa0:
		val, _ = cpu.nextValue(Immediate)
		cpu.ldy(val)
		fmt.Printf("LDY #$%02X", val)
	case 0xa4:
		val, addr = cpu.nextValue(ZeroPage)
		cpu.ldy(val)
		fmt.Printf("LDY $%02X", addr)
	case 0xb4:
		val, addr = cpu.nextValue(ZeroPageX)
		cpu.ldy(val)
		fmt.Printf("LDY $%02X, X", addr)
	case 0xac:
		val, addr = cpu.nextValue(Absolute)
		cpu.ldy(val)
		fmt.Printf("LDY $%04X", addr)
	case 0xbc:
		val, addr = cpu.nextValue(AbsoluteX)
		cpu.ldy(val)
		fmt.Printf("LDY $%04X, X", addr)

	case 0x85:
		_, addr = cpu.nextAddress(ZeroPage)
		cpu.sta(addr)
		fmt.Printf("STA $%02X", addr)
	case 0x95:
		_, addr = cpu.nextAddress(ZeroPageX)
		cpu.sta(addr)
		fmt.Printf("STA $%02X, X", addr)
	case 0x8d:
		_, addr = cpu.nextAddress(Absolute)
		cpu.sta(addr)
		fmt.Printf("STA $%04X", addr)
	case 0x9d:
		_, addr = cpu.nextAddress(AbsoluteX)
		cpu.sta(addr)
		fmt.Printf("STA $%04X, X", addr)
	case 0x99:
		_, addr = cpu.nextAddress(AbsoluteY)
		cpu.sta(addr)
		fmt.Printf("STA $%04X, Y", addr)
	case 0x81:
		_, addr = cpu.nextAddress(IndirectX)
		cpu.sta(addr)
		fmt.Printf("STA ($%02X, X)", addr)
	case 0x91:
		_, addr = cpu.nextAddress(IndirectY)
		cpu.sta(addr)
		fmt.Printf("STA ($%04X), Y", addr)

	case 0x86:
		_, addr = cpu.nextAddress(ZeroPage)
		cpu.stx(addr)
		fmt.Printf("STX $%02X", addr)
	case 0x96:
		_, addr = cpu.nextAddress(ZeroPageY)
		cpu.stx(addr)
		fmt.Printf("STX $%02X, Y", addr)
	case 0x8e:
		_, addr = cpu.nextAddress(Absolute)
		cpu.stx(addr)
		fmt.Printf("STX $%04X", addr)

	case 0x84:
		_, addr = cpu.nextAddress(ZeroPage)
		cpu.sty(addr)
		fmt.Printf("STY $%02X", addr)
	case 0x94:
		_, addr = cpu.nextAddress(ZeroPageX)
		cpu.sty(addr)
		fmt.Printf("STY $%02X, X", addr)
	case 0x8c:
		_, addr = cpu.nextAddress(Absolute)
		cpu.sty(addr)
		fmt.Printf("STY $%04X", addr)

	case 0xaa:
		cpu.tax()
		fmt.Print("TAX")
	case 0xa8:
		cpu.tay()
		fmt.Print("TAY")
	case 0xba:
		cpu.tsx()
		fmt.Print("TSX")
	case 0x8a:
		cpu.txa()
		fmt.Print("TXA")
	case 0x9a:
		cpu.txs()
		fmt.Print("TXS")
	case 0x98:
		cpu.tya()
		fmt.Print("TYA")

	case 0x48:
		cpu.pha()
		fmt.Print("PHA")
	case 0x08:
		cpu.php()
		fmt.Print("PHP")
	case 0x68:
		cpu.pla()
		fmt.Print("PLA")
	case 0x28:
		cpu.plp()
		fmt.Print("PLP")

	case 0x29:
		val, _ = cpu.nextValue(Immediate)
		cpu.and(val)
		fmt.Printf("AND #$%02X", val)
	case 0x25:
		val, addr = cpu.nextValue(ZeroPage)
		cpu.and(val)
		fmt.Printf("AND $%02X", addr)
	case 0x35:
		val, addr = cpu.nextValue(ZeroPageX)
		cpu.and(val)
		fmt.Printf("AND $%02X, X", addr)
	case 0x2d:
		val, addr = cpu.nextValue(Absolute)
		cpu.and(val)
		fmt.Printf("AND $%04X", addr)
	case 0x3d:
		val, addr = cpu.nextValue(AbsoluteX)
		cpu.and(val)
		fmt.Printf("AND $%04X, X", addr)
	case 0x39:
		val, addr = cpu.nextValue(AbsoluteY)
		cpu.and(val)
		fmt.Printf("AND $%04X, Y", addr)
	case 0x21:
		val, addr = cpu.nextValue(IndirectX)
		cpu.and(val)
		fmt.Printf("AND ($%02X, X)", addr)
	case 0x31:
		val, addr = cpu.nextValue(IndirectY)
		cpu.and(val)
		fmt.Printf("AND ($%04X), Y", addr)

	case 0x49:
		val, _ = cpu.nextValue(Immediate)
		cpu.eor(val)
		fmt.Printf("EOR #$%02X", val)
	case 0x45:
		val, addr = cpu.nextValue(ZeroPage)
		cpu.eor(val)
		fmt.Printf("EOR $%02X", addr)
	case 0x55:
		val, addr = cpu.nextValue(ZeroPageX)
		cpu.eor(val)
		fmt.Printf("EOR $%02X, X", addr)
	case 0x4d:
		val, addr = cpu.nextValue(Absolute)
		cpu.eor(val)
		fmt.Printf("EOR $%04X", addr)
	case 0x5d:
		val, addr = cpu.nextValue(AbsoluteX)
		cpu.eor(val)
		fmt.Printf("EOR $%04X, X", addr)
	case 0x59:
		val, addr = cpu.nextValue(AbsoluteY)
		cpu.eor(val)
		fmt.Printf("EOR $%04X, Y", addr)
	case 0x41:
		val, addr = cpu.nextValue(IndirectX)
		cpu.eor(val)
		fmt.Printf("EOR ($%02X, X)", addr)
	case 0x51:
		val, addr = cpu.nextValue(IndirectY)
		cpu.eor(val)
		fmt.Printf("EOR ($%04X), Y", addr)

	case 0x09:
		val, _ = cpu.nextValue(Immediate)
		cpu.ora(val)
		fmt.Printf("ORA #$%02X", val)
	case 0x05:
		val, addr = cpu.nextValue(ZeroPage)
		cpu.ora(val)
		fmt.Printf("ORA $%02X", addr)
	case 0x15:
		val, addr = cpu.nextValue(ZeroPageX)
		cpu.ora(val)
		fmt.Printf("ORA $%02X, X", addr)
	case 0x0d:
		val, addr = cpu.nextValue(Absolute)
		cpu.ora(val)
		fmt.Printf("ORA $%04X", addr)
	case 0x1d:
		val, addr = cpu.nextValue(AbsoluteX)
		cpu.ora(val)
		fmt.Printf("ORA $%04X, X", addr)
	case 0x19:
		val, addr = cpu.nextValue(AbsoluteY)
		cpu.ora(val)
		fmt.Printf("ORA $%04X, Y", addr)
	case 0x01:
		val, addr = cpu.nextValue(IndirectX)
		cpu.ora(val)
		fmt.Printf("ORA ($%02X, X)", addr)
	case 0x11:
		val, addr = cpu.nextValue(IndirectY)
		cpu.ora(val)
		fmt.Printf("ORA ($%04X), Y", addr)

	case 0x24:
		val, addr = cpu.nextValue(ZeroPage)
		cpu.bit(val)
		fmt.Printf("BIT $%02X", addr)
	case 0x2c:
		val, addr = cpu.nextValue(Absolute)
		cpu.bit(val)
		fmt.Printf("BIT $%04X", addr)

	case 0x69:
		val, _ = cpu.nextValue(Immediate)
		cpu.adc(val)
		fmt.Printf("ADC #$%02X", val)
	case 0x65:
		val, addr = cpu.nextValue(ZeroPage)
		cpu.adc(val)
		fmt.Printf("ADC $%02X", addr)
	case 0x75:
		val, addr = cpu.nextValue(ZeroPageX)
		cpu.adc(val)
		fmt.Printf("ADC $%02X, X", addr)
	case 0x6d:
		val, addr = cpu.nextValue(Absolute)
		cpu.adc(val)
		fmt.Printf("ADC $%04X", addr)
	case 0x7d:
		val, addr = cpu.nextValue(AbsoluteX)
		cpu.adc(val)
		fmt.Printf("ADC $%04X, X", addr)
	case 0x79:
		val, addr = cpu.nextValue(AbsoluteY)
		cpu.adc(val)
		fmt.Printf("ADC $%04X, Y", addr)
	case 0x61:
		val, addr = cpu.nextValue(IndirectX)
		cpu.adc(val)
		fmt.Printf("ADC ($%02X, X)", addr)
	case 0x71:
		val, addr = cpu.nextValue(IndirectY)
		cpu.adc(val)
		fmt.Printf("ADC ($%04X), Y", addr)

	case 0xe9:
		val, _ = cpu.nextValue(Immediate)
		cpu.sbc(val)
		fmt.Printf("SBC #$%02X", val)
	case 0xe5:
		val, addr = cpu.nextValue(ZeroPage)
		cpu.sbc(val)
		fmt.Printf("SBC $%02X", addr)
	case 0xf5:
		val, addr = cpu.nextValue(ZeroPageX)
		cpu.sbc(val)
		fmt.Printf("SBC $%02X, X", addr)
	case 0xed:
		val, addr = cpu.nextValue(Absolute)
		cpu.sbc(val)
		fmt.Printf("SBC $%04X", addr)
	case 0xfd:
		val, addr = cpu.nextValue(AbsoluteX)
		cpu.sbc(val)
		fmt.Printf("SBC $%04X, X", addr)
	case 0xf9:
		val, addr = cpu.nextValue(AbsoluteY)
		cpu.sbc(val)
		fmt.Printf("SBC $%04X, Y", addr)
	case 0xe1:
		val, addr = cpu.nextValue(IndirectX)
		cpu.sbc(val)
		fmt.Printf("SBC ($%02X, X)", addr)
	case 0xf1:
		val, addr = cpu.nextValue(IndirectY)
		cpu.sbc(val)
		fmt.Printf("SBC ($%04X), Y", addr)

	case 0xc9:
		val, _ = cpu.nextValue(Immediate)
		cpu.cmp(val)
		fmt.Printf("CMP #$%02X", val)
	case 0xc5:
		val, addr = cpu.nextValue(ZeroPage)
		cpu.cmp(val)
		fmt.Printf("CMP $%02X", addr)
	case 0xd5:
		val, addr = cpu.nextValue(ZeroPageX)
		cpu.cmp(val)
		fmt.Printf("CMP $%02X, X", addr)
	case 0xcd:
		val, addr = cpu.nextValue(Absolute)
		cpu.cmp(val)
		fmt.Printf("CMP $%04X", addr)
	case 0xdd:
		val, addr = cpu.nextValue(AbsoluteX)
		cpu.cmp(val)
		fmt.Printf("CMP $%04X, X", addr)
	case 0xd9:
		val, addr = cpu.nextValue(AbsoluteY)
		cpu.cmp(val)
		fmt.Printf("CMP $%04X, Y", addr)
	case 0xc1:
		val, addr = cpu.nextValue(IndirectX)
		cpu.cmp(val)
		fmt.Printf("CMP ($%02X, X)", addr)
	case 0xd1:
		val, addr = cpu.nextValue(IndirectY)
		cpu.cmp(val)
		fmt.Printf("CMP ($%04X), Y", addr)

	case 0xe0:
		val, _ = cpu.nextValue(Immediate)
		cpu.cpx(val)
		fmt.Printf("CPX #$%02X", val)
	case 0xe4:
		val, addr = cpu.nextValue(ZeroPage)
		cpu.cpx(val)
		fmt.Printf("CPX $%02X", addr)
	case 0xec:
		val, addr = cpu.nextValue(Absolute)
		cpu.cpx(val)
		fmt.Printf("CPX $%04X", addr)

	case 0xc0:
		val, _ = cpu.nextValue(Immediate)
		cpu.cpy(val)
		fmt.Printf("CPY #$%02X", val)
	case 0xc4:
		val, addr = cpu.nextValue(ZeroPage)
		cpu.cpy(val)
		fmt.Printf("CPY $%02X", addr)
	case 0xcc:
		val, addr = cpu.nextValue(Absolute)
		cpu.cpy(val)
		fmt.Printf("CPY $%04X", addr)

	case 0xe6:
		_, addr = cpu.nextAddress(ZeroPage)
		cpu.inc(addr)
		fmt.Printf("INC $%02X", addr)
	case 0xf6:
		_, addr = cpu.nextAddress(ZeroPageX)
		cpu.inc(addr)
		fmt.Printf("INC $%02X, X", addr)
	case 0xee:
		_, addr = cpu.nextAddress(Absolute)
		cpu.inc(addr)
		fmt.Printf("INC $%04X", addr)
	case 0xfe:
		_, addr = cpu.nextAddress(AbsoluteX)
		cpu.inc(addr)
		fmt.Printf("INC $%04X, X", addr)

	case 0xc6:
		_, addr = cpu.nextAddress(ZeroPage)
		cpu.dec(addr)
		fmt.Printf("DEC $%02X", addr)
	case 0xd6:
		_, addr = cpu.nextAddress(ZeroPageX)
		cpu.dec(addr)
		fmt.Printf("DEC $%02X, X", addr)
	case 0xce:
		_, addr = cpu.nextAddress(Absolute)
		cpu.dec(addr)
		fmt.Printf("DEC $%04X", addr)
	case 0xde:
		_, addr = cpu.nextAddress(AbsoluteX)
		cpu.dec(addr)
		fmt.Printf("DEC $%04X, X", addr)

	case 0xe8:
		cpu.inx()
		fmt.Print("INX")
	case 0xc8:
		cpu.iny()
		fmt.Print("INY")

	case 0xca:
		cpu.dex()
		fmt.Print("DEX")
	case 0x88:
		cpu.dey()
		fmt.Print("DEY")

	case 0x0a:
		cpu.asl_acc()
	case 0x06:
		_, addr = cpu.nextAddress(ZeroPage)
		cpu.asl(addr)
		fmt.Printf("ASL $%02X", addr)
	case 0x16:
		_, addr = cpu.nextAddress(ZeroPageX)
		cpu.asl(addr)
		fmt.Printf("ASL $%02X, X", addr)
	case 0x0e:
		_, addr = cpu.nextAddress(Absolute)
		cpu.asl(addr)
		fmt.Printf("ASL $%04X", addr)
	case 0x1e:
		_, addr = cpu.nextAddress(AbsoluteX)
		cpu.asl(addr)
		fmt.Printf("ASL $%04X, X", addr)

	case 0x4a:
		cpu.lsr_acc()
	case 0x46:
		_, addr = cpu.nextAddress(ZeroPage)
		cpu.lsr(addr)
		fmt.Printf("LSR $%02X", addr)
	case 0x56:
		_, addr = cpu.nextAddress(ZeroPageX)
		cpu.lsr(addr)
		fmt.Printf("LSR $%02X, X", addr)
	case 0x4e:
		_, addr = cpu.nextAddress(Absolute)
		cpu.lsr(addr)
		fmt.Printf("LSR $%04X", addr)
	case 0x5e:
		_, addr = cpu.nextAddress(AbsoluteX)
		cpu.lsr(addr)
		fmt.Printf("LSR $%04X, X", addr)

	case 0x2a:
		cpu.rol_acc()
		fmt.Print("ROL A")
	case 0x26:
		_, addr = cpu.nextAddress(ZeroPage)
		cpu.rol(addr)
		fmt.Printf("ROL $%02X", addr)
	case 0x36:
		_, addr = cpu.nextAddress(ZeroPageX)
		cpu.rol(addr)
		fmt.Printf("ROL $%02X, X", addr)
	case 0x2e:
		_, addr = cpu.nextAddress(Absolute)
		cpu.rol(addr)
		fmt.Printf("ROL $%04X", addr)
	case 0x3e:
		_, addr = cpu.nextAddress(AbsoluteX)
		cpu.rol(addr)
		fmt.Printf("ROL $%04X, X", addr)

	case 0x6a:
		cpu.ror_acc()
		fmt.Print("ROR A")
	case 0x66:
		_, addr = cpu.nextAddress(ZeroPage)
		cpu.ror(addr)
		fmt.Printf("ROR $%02X", addr)
	case 0x76:
		_, addr = cpu.nextAddress(ZeroPageX)
		cpu.ror(addr)
		fmt.Printf("ROR $%02X, X", addr)
	case 0x6e:
		_, addr = cpu.nextAddress(Absolute)
		cpu.ror(addr)
		fmt.Printf("ROR $%04X", addr)
	case 0x7e:
		_, addr = cpu.nextAddress(AbsoluteX)
		cpu.ror(addr)
		fmt.Printf("ROR $%04X, X", addr)

	case 0x4c:
		_, addr = cpu.nextAddress(Absolute)
		cpu.jmp(addr)
		fmt.Printf("JMP $%04X", addr)
	case 0x6c:
		_, addr = cpu.nextAddress(Indirect)
		cpu.jmp(addr)
		fmt.Printf("JMP ($%04X)", addr)

	case 0x20:
		_, addr = cpu.nextAddress(Absolute)
		cpu.jsr(addr)
		fmt.Printf("JSR $%04X", addr)

	case 0x60:
		cpu.rts()
		fmt.Print("RTS")

	case 0x90:
		val, _ = cpu.nextValue(Immediate)
		cpu.bcc(val)
		fmt.Printf("BCC #$%02X", val)
	case 0xb0:
		val, _ = cpu.nextValue(Immediate)
		cpu.bcs(val)
		fmt.Printf("BCS #$%02X", val)
	case 0xf0:
		val, _ = cpu.nextValue(Immediate)
		cpu.beq(val)
		fmt.Printf("BEQ #$%02X", val)
	case 0x30:
		val, _ = cpu.nextValue(Immediate)
		cpu.bmi(val)
		fmt.Printf("BMI #$%02X", val)
	case 0xd0:
		val, _ = cpu.nextValue(Immediate)
		cpu.bne(val)
		fmt.Printf("BNE #$%02X", val)
	case 0x10:
		val, _ = cpu.nextValue(Immediate)
		cpu.bpl(val)
		fmt.Printf("BPL #$%02X", val)
	case 0x50:
		val, _ = cpu.nextValue(Immediate)
		cpu.bvc(val)
		fmt.Printf("BVC #$%02X", val)
	case 0x70:
		val, _ = cpu.nextValue(Immediate)
		cpu.bvs(val)
		fmt.Printf("BVS #$%02X", val)

	case 0x18:
		cpu.clc()
		fmt.Print("CLC")
	case 0xd8:
		cpu.cld()
		fmt.Print("CLD")
	case 0x58:
		cpu.cli()
		fmt.Print("CLI")
	case 0xb8:
		cpu.clv()
		fmt.Print("CLV")
	case 0x38:
		cpu.sec()
		fmt.Print("SEC")
	case 0xf8:
		cpu.sed()
		fmt.Print("SED")
	case 0x78:
		cpu.sei()
		fmt.Print("SEI")

	default:
		fmt.Printf("UNKNOWN")
	}

	fmt.Print("\n")
	fmt.Print(cpu.Dump())
}

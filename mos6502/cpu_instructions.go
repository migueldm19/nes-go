package mos6502

import "fmt"

func (cpu *CPU) brk() {
	cpu.stackPushCurrentPc(2)
	cpu.stackPush(cpu.p | FlagB)

	addr1, _ := cpu.mem.Read(0xfffe)
	addr2, _ := cpu.mem.Read(0xffff)

	cpu.pc = (uint16(addr1) << 8) + uint16(addr2)
}

func (cpu *CPU) rti() {
	cpu.p = (cpu.stackPull() | 0x20) & 0xef
	cpu.pc = cpu.stackPullAddr()
}

func (cpu *CPU) lda(val byte) {
	cpu.a = val
	cpu.assignBasicFlags(cpu.a)
}

func (cpu *CPU) ldx(val byte) {
	cpu.x = val
	cpu.assignBasicFlags(cpu.x)
}

func (cpu *CPU) ldy(val byte) {
	cpu.y = val
	cpu.assignBasicFlags(cpu.y)
}

func (cpu *CPU) sta(addr uint16) {
	cpu.write(cpu.a, addr)
}

func (cpu *CPU) stx(addr uint16) {
	cpu.write(cpu.x, addr)
}

func (cpu *CPU) sty(addr uint16) {
	cpu.write(cpu.y, addr)
}

func (cpu *CPU) tax() {
	cpu.x = cpu.a
	cpu.assignBasicFlags(cpu.x)
}

func (cpu *CPU) tay() {
	cpu.y = cpu.a
	cpu.assignBasicFlags(cpu.y)
}

func (cpu *CPU) tsx() {
	cpu.x = cpu.sp
	cpu.assignBasicFlags(cpu.x)
}

func (cpu *CPU) txa() {
	cpu.a = cpu.x
	cpu.assignBasicFlags(cpu.a)
}

func (cpu *CPU) txs() {
	cpu.sp = cpu.x
	fmt.Printf("txs")
}

func (cpu *CPU) tya() {
	cpu.a = cpu.y
	cpu.assignBasicFlags(cpu.a)
}

func (cpu *CPU) pha() {
	cpu.stackPush(cpu.a)
}

func (cpu *CPU) php() {
	cpu.stackPush(cpu.p | FlagB)
}

func (cpu *CPU) pla() {
	cpu.a = cpu.stackPull()
	cpu.assignBasicFlags(cpu.a)
}

func (cpu *CPU) plp() {
	cpu.p = (cpu.stackPull() | 0x20) & 0xef
}

func (cpu *CPU) and(val byte) {
	cpu.a &= val
	cpu.assignBasicFlags(cpu.a)
}

func (cpu *CPU) eor(val byte) {
	cpu.a ^= val
	cpu.assignBasicFlags(cpu.a)
}

func (cpu *CPU) ora(val byte) {
	cpu.a |= val
	cpu.assignBasicFlags(cpu.a)
}

func (cpu *CPU) bit(val byte) {
	cpu.bitTest(val)
}

func (cpu *CPU) adc(val byte) {
	prev := cpu.a
	val1, overflow1 := addOverflow(cpu.a, val)

	var carry byte
	if cpu.getFlag(FlagCarry) {
		carry = 1
	}

	val2, overflow2 := addOverflow(val1, carry)
	cpu.a = val2

	cpu.setFlag(FlagZero, cpu.a == 0)
	cpu.setFlag(FlagCarry, overflow1 || overflow2)
	cpu.setFlag(FlagNegative, isNegative(cpu.a))
	cpu.setFlag(FlagOverflow, ((cpu.a^prev)&(cpu.a^val)&0x80) == 0x80)
}

func (cpu *CPU) sbc(val byte) {
	prev := cpu.a
	val1, overflow1 := subOverflow(cpu.a, val)

	var carry byte
	if !cpu.getFlag(FlagCarry) {
		carry = 1
	}

	val2, overflow2 := subOverflow(val1, carry)
	cpu.a = val2

	cpu.setFlag(FlagCarry, !(overflow1 || overflow2))
	cpu.setFlag(FlagZero, cpu.a == 0)
	cpu.setFlag(FlagNegative, isNegative(cpu.a))
	cpu.setFlag(FlagOverflow, ((cpu.a^prev)&(cpu.a^(^val))&0x80) == 0x80)
}

func (cpu *CPU) cmp(val byte) {
	cpu.compare(cpu.a, val)
}

func (cpu *CPU) cpx(val byte) {
	cpu.compare(cpu.x, val)
}

func (cpu *CPU) cpy(val byte) {
	cpu.compare(cpu.y, val)
}

func (cpu *CPU) inc(addr uint16) {
	val := cpu.read(addr) + 1
	cpu.assignBasicFlags(val)
	cpu.write(val, addr)
}

func (cpu *CPU) inx() {
	cpu.x += 1
	cpu.assignBasicFlags(cpu.x)
}

func (cpu *CPU) iny() {
	cpu.y += 1
	cpu.assignBasicFlags(cpu.y)
}

func (cpu *CPU) dec(addr uint16) {
	val := cpu.read(addr) - 1
	cpu.assignBasicFlags(val)
	cpu.write(val, addr)
}

func (cpu *CPU) dex() {
	cpu.x -= 1
	cpu.assignBasicFlags(cpu.x)
}

func (cpu *CPU) dey() {
	cpu.y -= 1
	cpu.assignBasicFlags(cpu.y)
}

func (cpu *CPU) asl_acc() {
	cpu.setFlag(FlagCarry, cpu.a&0x80 == 0x80)
	cpu.a = cpu.a << 1
	cpu.assignBasicFlags(cpu.a)
}

func (cpu *CPU) asl(addr uint16) {
	val := cpu.read(addr)

	cpu.setFlag(FlagCarry, val&0x80 == 0x80)

	val = val << 1

	cpu.assignBasicFlags(val)
	cpu.write(val, addr)
}

func (cpu *CPU) lsr_acc() {
	cpu.setFlag(FlagCarry, cpu.a&0x01 == 0x01)
	cpu.a = cpu.a >> 1
	cpu.assignBasicFlags(cpu.a)
}

func (cpu *CPU) lsr(addr uint16) {
	val := cpu.read(addr)

	cpu.setFlag(FlagCarry, val&0x01 == 0x01)

	val = val >> 1

	cpu.assignBasicFlags(val)
	cpu.write(val, addr)
}

func (cpu *CPU) rol_acc() {
	prev_carry := cpu.getFlag(FlagCarry)
	cpu.setFlag(FlagCarry, isNegative(cpu.a))
	cpu.a = cpu.a << 1
	if prev_carry {
		cpu.a += 1
	}
	cpu.assignBasicFlags(cpu.a)
}

func (cpu *CPU) rol(addr uint16) {
	val := cpu.read(addr)

	prev_carry := cpu.getFlag(FlagCarry)
	cpu.setFlag(FlagCarry, isNegative(val))

	val = val << 1
	if prev_carry {
		cpu.a += 1
	}

	cpu.assignBasicFlags(val)
	cpu.write(val, addr)
}

func (cpu *CPU) ror_acc() {
	prev_carry := cpu.getFlag(FlagCarry)
	cpu.setFlag(FlagCarry, cpu.a&0x01 == 0x01)
	cpu.a = cpu.a >> 1
	if prev_carry {
		cpu.a += 0x80
	}
	cpu.assignBasicFlags(cpu.a)
}

func (cpu *CPU) ror(addr uint16) {
	val := cpu.read(addr)

	prev_carry := cpu.getFlag(FlagCarry)
	cpu.setFlag(FlagCarry, val&0x01 == 0x01)

	val = val << 1
	if prev_carry {
		cpu.a += 0x80
	}

	cpu.assignBasicFlags(val)
	cpu.write(val, addr)
}

func (cpu *CPU) jmp(addr uint16) {
	cpu.pc = addr
}

func (cpu *CPU) jsr(addr uint16) {
	cpu.stackPushCurrentPc(-1)
	cpu.pc = addr
}

func (cpu *CPU) rts() {
	cpu.pc = cpu.stackPullAddr() + 1
}

func (cpu *CPU) bcc(displacement byte) {
	if !cpu.getFlag(FlagCarry) {
		cpu.branchJump(int8(displacement))
	}
}

func (cpu *CPU) bcs(displacement byte) {
	if cpu.getFlag(FlagCarry) {
		cpu.branchJump(int8(displacement))
	}
}

func (cpu *CPU) beq(displacement byte) {
	if cpu.getFlag(FlagZero) {
		cpu.branchJump(int8(displacement))
	}
}

func (cpu *CPU) bmi(displacement byte) {
	if cpu.getFlag(FlagNegative) {
		cpu.branchJump(int8(displacement))
	}
}

func (cpu *CPU) bne(displacement byte) {
	if !cpu.getFlag(FlagZero) {
		cpu.branchJump(int8(displacement))
	}
}

func (cpu *CPU) bpl(displacement byte) {
	if !cpu.getFlag(FlagNegative) {
		cpu.branchJump(int8(displacement))
	}
}

func (cpu *CPU) bvc(displacement byte) {
	if !cpu.getFlag(FlagOverflow) {
		cpu.branchJump(int8(displacement))
	}
}

func (cpu *CPU) bvs(displacement byte) {
	if cpu.getFlag(FlagOverflow) {
		cpu.branchJump(int8(displacement))
	}
}

func (cpu *CPU) clc() {
	cpu.setFlag(FlagCarry, false)
}

func (cpu *CPU) cld() {
	cpu.setFlag(FlagDecimalMode, false)
}

func (cpu *CPU) cli() {
	cpu.setFlag(FlagInterruptDisable, false)
}

func (cpu *CPU) clv() {
	cpu.setFlag(FlagOverflow, false)
}

func (cpu *CPU) sec() {
	cpu.setFlag(FlagCarry, true)
}

func (cpu *CPU) sed() {
	cpu.setFlag(FlagDecimalMode, true)
}

func (cpu *CPU) sei() {
	cpu.setFlag(FlagInterruptDisable, true)
}

func (cpu *CPU) Run() {
	for {
		cpu.Step()
	}
}

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
		fmt.Printf("LDA ($%04X, X)", addr)
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
		addr = cpu.nextAddress(ZeroPage)
		cpu.sta(addr)
		fmt.Printf("STA $%02X", addr)
	case 0x95:
		addr = cpu.nextAddress(ZeroPageX)
		cpu.sta(addr)
		fmt.Printf("STA $%02X, X", addr)
	case 0x8d:
		addr = cpu.nextAddress(Absolute)
		cpu.sta(addr)
		fmt.Printf("STA $%04X", addr)
	case 0x9d:
		addr = cpu.nextAddress(AbsoluteX)
		cpu.sta(addr)
		fmt.Printf("STA $%04X, X", addr)
	case 0x99:
		addr = cpu.nextAddress(AbsoluteY)
		cpu.sta(addr)
		fmt.Printf("STA $%04X, Y", addr)
	case 0x81:
		addr = cpu.nextAddress(IndirectX)
		cpu.sta(addr)
		fmt.Printf("STA ($%04X, X)", addr)
	case 0x91:
		addr = cpu.nextAddress(IndirectY)
		cpu.sta(addr)
		fmt.Printf("STA ($%04X), Y", addr)

	case 0x86:
		addr = cpu.nextAddress(ZeroPage)
		cpu.stx(addr)
		fmt.Printf("STX $%02X", addr)
	case 0x96:
		addr = cpu.nextAddress(ZeroPageY)
		cpu.stx(addr)
		fmt.Printf("STX $%02X, Y", addr)
	case 0x8e:
		addr = cpu.nextAddress(Absolute)
		cpu.stx(addr)
		fmt.Printf("STX $%04X", addr)

	case 0x84:
		addr = cpu.nextAddress(ZeroPage)
		cpu.sty(addr)
		fmt.Printf("STY $%02X", addr)
	case 0x94:
		addr = cpu.nextAddress(ZeroPageX)
		cpu.sty(addr)
		fmt.Printf("STY $%02X, X", addr)
	case 0x8c:
		addr = cpu.nextAddress(Absolute)
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
		fmt.Printf("AND ($%04X, X)", addr)
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
		fmt.Printf("EOR ($%04X, X)", addr)
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
		fmt.Printf("ORA ($%04X, X)", addr)
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
		fmt.Printf("ADC ($%04X, X)", addr)
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
		fmt.Printf("SBC ($%04X, X)", addr)
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
		fmt.Printf("CMP ($%04X, X)", addr)
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
		val, addr = cpu.nextValue(Immediate)
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
		addr = cpu.nextAddress(ZeroPage)
		cpu.inc(addr)
		fmt.Printf("INC $%02X", addr)
	case 0xf6:
		addr = cpu.nextAddress(ZeroPageX)
		cpu.inc(addr)
		fmt.Printf("INC $%02X, X", addr)
	case 0xee:
		addr = cpu.nextAddress(Absolute)
		cpu.inc(addr)
		fmt.Printf("INC $%04X", addr)
	case 0xfe:
		addr = cpu.nextAddress(AbsoluteX)
		cpu.inc(addr)
		fmt.Printf("INC $%04X, X", addr)

	case 0xc6:
		addr = cpu.nextAddress(ZeroPage)
		cpu.dec(addr)
		fmt.Printf("DEC $%02X", addr)
	case 0xd6:
		addr = cpu.nextAddress(ZeroPageX)
		cpu.dec(addr)
		fmt.Printf("DEC $%02X, X", addr)
	case 0xce:
		addr = cpu.nextAddress(Absolute)
		cpu.dec(addr)
		fmt.Printf("DEC $%04X", addr)
	case 0xde:
		addr = cpu.nextAddress(AbsoluteX)
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
		addr = cpu.nextAddress(ZeroPage)
		cpu.asl(addr)
		fmt.Printf("ASL $%02X", addr)
	case 0x16:
		addr = cpu.nextAddress(ZeroPageX)
		cpu.asl(addr)
		fmt.Printf("ASL $%02X, X", addr)
	case 0x0e:
		addr = cpu.nextAddress(Absolute)
		cpu.asl(addr)
		fmt.Printf("ASL $%04X", addr)
	case 0x1e:
		addr = cpu.nextAddress(AbsoluteX)
		cpu.asl(addr)
		fmt.Printf("ASL $%04X, X", addr)

	case 0x4a:
		cpu.lsr_acc()
	case 0x46:
		addr = cpu.nextAddress(ZeroPage)
		cpu.lsr(addr)
		fmt.Printf("LSR $%02X", addr)
	case 0x56:
		addr = cpu.nextAddress(ZeroPageX)
		cpu.lsr(addr)
		fmt.Printf("LSR $%02X, X", addr)
	case 0x4e:
		addr = cpu.nextAddress(Absolute)
		cpu.lsr(addr)
		fmt.Printf("LSR $%04X", addr)
	case 0x5e:
		addr = cpu.nextAddress(AbsoluteX)
		cpu.lsr(addr)
		fmt.Printf("LSR $%04X, X", addr)

	case 0x2a:
		cpu.rol_acc()
		fmt.Print("ROL A")
	case 0x26:
		addr = cpu.nextAddress(ZeroPage)
		cpu.rol(addr)
		fmt.Printf("ROL $%02X", addr)
	case 0x36:
		addr = cpu.nextAddress(ZeroPageX)
		cpu.rol(addr)
		fmt.Printf("ROL $%02X, X", addr)
	case 0x2e:
		addr = cpu.nextAddress(Absolute)
		cpu.rol(addr)
		fmt.Printf("ROL $%04X", addr)
	case 0x3e:
		addr = cpu.nextAddress(AbsoluteX)
		cpu.rol(addr)
		fmt.Printf("ROL $%04X, X", addr)

	case 0x6a:
		cpu.ror_acc()
		fmt.Print("ROR A")
	case 0x66:
		addr = cpu.nextAddress(ZeroPage)
		cpu.ror(addr)
		fmt.Printf("ROR $%02X", addr)
	case 0x76:
		addr = cpu.nextAddress(ZeroPageX)
		cpu.ror(addr)
		fmt.Printf("ROR $%02X, X", addr)
	case 0x6e:
		addr = cpu.nextAddress(Absolute)
		cpu.ror(addr)
		fmt.Printf("ROR $%04X", addr)
	case 0x7e:
		addr = cpu.nextAddress(AbsoluteX)
		cpu.ror(addr)
		fmt.Printf("ROR $%04X, X", addr)

	case 0x4c:
		addr = cpu.nextAddress(Absolute)
		cpu.jmp(addr)
		fmt.Printf("JMP $%04X", addr)
	case 0x6c:
		addr = cpu.nextAddress(Indirect)
		cpu.jmp(addr)
		fmt.Printf("JMP ($%04X)", addr)

	case 0x20:
		addr = cpu.nextAddress(Absolute)
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
}

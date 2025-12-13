package mos6502

func (cpu *CPU) brk() {
	cpu.stackPushCurrentPc(2)
	cpu.stackPush(cpu.p | FlagB)

	addr1, _ := cpu.Mem.ReadCpu(0xfffe)
	addr2, _ := cpu.Mem.ReadCpu(0xffff)

	cpu.Pc = (uint16(addr1) << 8) + uint16(addr2)
}

func (cpu *CPU) rti() {
	cpu.p = (cpu.stackPull() | 0x20) & 0xef
	cpu.Pc = cpu.stackPullAddr()
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
		val += 1
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

	val = val >> 1
	if prev_carry {
		val += 0x80
	}

	cpu.assignBasicFlags(val)
	cpu.write(val, addr)
}

func (cpu *CPU) jmp(addr uint16) {
	cpu.Pc = addr
}

func (cpu *CPU) jsr(addr uint16) {
	cpu.stackPushCurrentPc(-1)
	cpu.Pc = addr
}

func (cpu *CPU) rts() {
	cpu.Pc = cpu.stackPullAddr() + 1
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

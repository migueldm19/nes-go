package mos6502

import "fmt"

func (cpu *CPU) brk() {
	cpu.stackPushCurrentPc(2)
	cpu.stackPush(cpu.p | FlagB)

	addr1, _ := cpu.mem.Read(0xfffe)
	addr2, _ := cpu.mem.Read(0xffff)

	cpu.pc = (uint16(addr1) << 8) + uint16(addr2)
	fmt.Printf("brk")
}

func (cpu *CPU) rti() {
	cpu.p = (cpu.stackPull() | 0x20) & 0xef
	cpu.pc = cpu.stackPullAddr()
	fmt.Printf("rti")
}

func (cpu *CPU) lda(val byte) {
	cpu.a = val
	cpu.assignBasicFlags(cpu.a)
	fmt.Printf("lda %x", cpu.a)
}

func (cpu *CPU) ldx(val byte) {
	cpu.x = val
	cpu.assignBasicFlags(cpu.x)
	fmt.Printf("ldx %x", cpu.x)
}

func (cpu *CPU) ldy(val byte) {
	cpu.y = val
	cpu.assignBasicFlags(cpu.y)
	fmt.Printf("ldy %x", cpu.y)
}

func (cpu *CPU) sta(addr uint16) {
	cpu.write(cpu.a, addr)
	fmt.Printf("sta %x %x", addr, cpu.a)
}

func (cpu *CPU) stx(addr uint16) {
	cpu.write(cpu.x, addr)
	fmt.Printf("stx %x %x", addr, cpu.x)
}

func (cpu *CPU) sty(addr uint16) {
	cpu.write(cpu.y, addr)
	fmt.Printf("sty %x %x", addr, cpu.y)
}

func (cpu *CPU) tax() {
	cpu.x = cpu.a
	cpu.assignBasicFlags(cpu.x)
	fmt.Printf("tax")
}

func (cpu *CPU) tay() {
	cpu.y = cpu.a
	cpu.assignBasicFlags(cpu.y)
	fmt.Printf("tay")
}

func (cpu *CPU) tsx() {
	cpu.x = cpu.sp
	cpu.assignBasicFlags(cpu.x)
	fmt.Printf("tsx")
}

func (cpu *CPU) txa() {
	cpu.a = cpu.x
	cpu.assignBasicFlags(cpu.a)
	fmt.Printf("txa")
}

func (cpu *CPU) txs() {
	cpu.sp = cpu.x
	fmt.Printf("txs")
}

func (cpu *CPU) tya() {
	cpu.a = cpu.y
	cpu.assignBasicFlags(cpu.a)
	fmt.Printf("tya")
}

func (cpu *CPU) pha() {
	cpu.stackPush(cpu.a)
	fmt.Printf("pha %x", cpu.a)
}

func (cpu *CPU) php() {
	cpu.stackPush(cpu.p | FlagB)
	fmt.Printf("php %x", cpu.p|FlagB)
}

func (cpu *CPU) pla() {
	cpu.a = cpu.stackPull()
	cpu.assignBasicFlags(cpu.a)
	fmt.Printf("pla %x", cpu.a)
}

func (cpu *CPU) plp() {
	cpu.p = (cpu.stackPull() | 0x20) & 0xef
	fmt.Printf("plp")
}

func (cpu *CPU) and(val byte) {
	cpu.a &= val
	cpu.assignBasicFlags(cpu.a)
	fmt.Printf("and %x %x", val, cpu.a)
}

func (cpu *CPU) eor(val byte) {
	cpu.a ^= val
	cpu.assignBasicFlags(cpu.a)
	fmt.Printf("eor %x %x", val, cpu.a)
}

func (cpu *CPU) ora(val byte) {
	cpu.a |= val
	cpu.assignBasicFlags(cpu.a)
	fmt.Printf("ora %x %x", val, cpu.a)
}

func (cpu *CPU) bit(val byte) {
	cpu.bitTest(val)
	fmt.Printf("bit %x", val)
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
	fmt.Printf("adc %x %x", val, cpu.a)
}

// TODO: Revisar
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
	fmt.Printf("sbc %x %x", val, cpu.a)
}

func (cpu *CPU) cmp(val byte) {
	cpu.compare(cpu.a, val)
	fmt.Printf("cmp %x, %x", cpu.a, val)
}

func (cpu *CPU) cpx(val byte) {
	cpu.compare(cpu.x, val)
	fmt.Printf("cpx %x, %x", cpu.x, val)
}

func (cpu *CPU) cpy(val byte) {
	cpu.compare(cpu.y, val)
	fmt.Printf("cpy %x, %x", cpu.y, val)
}

func (cpu *CPU) inc(addr uint16) {
	val := cpu.read(addr) + 1
	cpu.assignBasicFlags(val)
	cpu.write(val, addr)
	fmt.Printf("inc %x %x", addr, val-1)
}

func (cpu *CPU) inx() {
	cpu.x += 1
	cpu.assignBasicFlags(cpu.x)
	fmt.Printf("inx %v", cpu.x-1)
}

func (cpu *CPU) iny() {
	cpu.y += 1
	cpu.assignBasicFlags(cpu.y)
	fmt.Printf("iny %v", cpu.y-1)
}

func (cpu *CPU) dec(addr uint16) {
	val := cpu.read(addr) - 1
	cpu.assignBasicFlags(val)
	cpu.write(val, addr)
	fmt.Printf("dec %x %x", addr, val+1)
}

func (cpu *CPU) dex() {
	cpu.x -= 1
	cpu.assignBasicFlags(cpu.x)
	fmt.Printf("inx %v", cpu.x+1)
}

func (cpu *CPU) dey() {
	cpu.y -= 1
	cpu.assignBasicFlags(cpu.y)
	fmt.Printf("iny %v", cpu.y+1)
}

func (cpu *CPU) asl_acc() {
	cpu.setFlag(FlagCarry, cpu.a&0x80 == 0x80)
	cpu.a = cpu.a << 1
	cpu.assignBasicFlags(cpu.a)
	fmt.Printf("asl acc %v", cpu.a)
}

func (cpu *CPU) asl(addr uint16) {
	val := cpu.read(addr)

	cpu.setFlag(FlagCarry, val&0x80 == 0x80)

	val = val << 1

	cpu.assignBasicFlags(val)
	cpu.write(val, addr)
	fmt.Printf("asl %x %x", addr, val)
}

func (cpu *CPU) lsr_acc() {
	cpu.setFlag(FlagCarry, cpu.a&0x01 == 0x01)
	cpu.a = cpu.a >> 1
	cpu.assignBasicFlags(cpu.a)
	fmt.Printf("lsr acc %v", cpu.a)
}

func (cpu *CPU) lsr(addr uint16) {
	val := cpu.read(addr)

	cpu.setFlag(FlagCarry, val&0x01 == 0x01)

	val = val >> 1

	cpu.assignBasicFlags(val)
	cpu.write(val, addr)
	fmt.Printf("lsr %x %x", addr, val)
}

func (cpu *CPU) rol_acc() {
	prev_carry := cpu.getFlag(FlagCarry)
	cpu.setFlag(FlagCarry, isNegative(cpu.a))
	cpu.a = cpu.a << 1
	if prev_carry {
		cpu.a += 1
	}
	cpu.assignBasicFlags(cpu.a)
	fmt.Printf("rol acc %v", cpu.a)
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
	fmt.Printf("rol %x %x", addr, val)
}

func (cpu *CPU) ror_acc() {
	prev_carry := cpu.getFlag(FlagCarry)
	cpu.setFlag(FlagCarry, cpu.a&0x01 == 0x01)
	cpu.a = cpu.a >> 1
	if prev_carry {
		cpu.a += 0x80
	}
	cpu.assignBasicFlags(cpu.a)
	fmt.Printf("ror acc %v", cpu.a)
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
	fmt.Printf("ror %x %x", addr, val)
}

func (cpu *CPU) jmp(addr uint16) {
	cpu.pc = addr
	fmt.Printf("jmp %x", addr)
}

func (cpu *CPU) jsr(addr uint16) {
	cpu.stackPushCurrentPc(-1)
	cpu.pc = addr
	fmt.Printf("jsr %x", addr)
}

func (cpu *CPU) rts() {
	cpu.pc = cpu.stackPullAddr() + 1
	fmt.Printf("rts %x", cpu.pc)
}

func (cpu *CPU) bcc(displacement byte) {
	if !cpu.getFlag(FlagCarry) {
		cpu.branchJump(int8(displacement))
	}
	fmt.Printf("bcc %v", displacement)
}

func (cpu *CPU) bcs(displacement byte) {
	if cpu.getFlag(FlagCarry) {
		cpu.branchJump(int8(displacement))
	}
	fmt.Printf("bcs %v", displacement)
}

func (cpu *CPU) beq(displacement byte) {
	if cpu.getFlag(FlagZero) {
		cpu.branchJump(int8(displacement))
	}
	fmt.Printf("beq %v Zero flag[%v]", displacement, cpu.getFlag(FlagZero))
}

func (cpu *CPU) bmi(displacement byte) {
	if cpu.getFlag(FlagNegative) {
		cpu.branchJump(int8(displacement))
	}
	fmt.Printf("bmi %v", displacement)
}

func (cpu *CPU) bne(displacement byte) {
	if !cpu.getFlag(FlagZero) {
		cpu.branchJump(int8(displacement))
	}
	fmt.Printf("bne %v", displacement)
}

func (cpu *CPU) bpl(displacement byte) {
	if !cpu.getFlag(FlagNegative) {
		cpu.branchJump(int8(displacement))
	}
	fmt.Printf("bpl %v", displacement)
}

func (cpu *CPU) bvc(displacement byte) {
	if !cpu.getFlag(FlagOverflow) {
		cpu.branchJump(int8(displacement))
	}
	fmt.Printf("bvc %v", displacement)
}

func (cpu *CPU) bvs(displacement byte) {
	if cpu.getFlag(FlagOverflow) {
		cpu.branchJump(int8(displacement))
	}
	fmt.Printf("bvs %v", displacement)
}

func (cpu *CPU) clc() {
	cpu.setFlag(FlagCarry, false)
	fmt.Printf("clc")
}

func (cpu *CPU) cld() {
	cpu.setFlag(FlagDecimalMode, false)
	fmt.Printf("cld")
}

func (cpu *CPU) cli() {
	cpu.setFlag(FlagInterruptDisable, false)
	fmt.Printf("cli")
}

func (cpu *CPU) clv() {
	cpu.setFlag(FlagOverflow, false)
	fmt.Printf("clv")
}

func (cpu *CPU) sec() {
	cpu.setFlag(FlagCarry, true)
	fmt.Printf("sec")
}

func (cpu *CPU) sed() {
	cpu.setFlag(FlagDecimalMode, true)
	fmt.Printf("sed")
}

func (cpu *CPU) sei() {
	cpu.setFlag(FlagInterruptDisable, true)
	fmt.Printf("sei")
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
	fmt.Printf("[PC: %4x] OPCODE %2x | %s |", cpu.pc-1, opcode, cpu)

	switch opcode {
	case 0x00:
		cpu.brk()

	case 0x40:
		cpu.rti()

	case 0xea:
		fmt.Printf("nop")

	case 0xa9:
		val = cpu.nextValue(Immediate)
		cpu.lda(val)
	case 0xa5:
		val = cpu.nextValue(ZeroPage)
		cpu.lda(val)
	case 0xb5:
		val = cpu.nextValue(ZeroPageX)
		cpu.lda(val)
	case 0xad:
		val = cpu.nextValue(Absolute)
		cpu.lda(val)
	case 0xbd:
		val = cpu.nextValue(AbsoluteX)
		cpu.lda(val)
	case 0xb9:
		val = cpu.nextValue(AbsoluteY)
		cpu.lda(val)
	case 0xa1:
		val = cpu.nextValue(IndirectX)
		cpu.lda(val)
	case 0xb1:
		val = cpu.nextValue(IndirectY)
		cpu.lda(val)

	case 0xa2:
		val = cpu.nextValue(Immediate)
		cpu.ldx(val)
	case 0xa6:
		val = cpu.nextValue(ZeroPage)
		cpu.ldx(val)
	case 0xb6:
		val = cpu.nextValue(ZeroPageY)
		cpu.ldx(val)
	case 0xae:
		val = cpu.nextValue(Absolute)
		cpu.ldx(val)
	case 0xbe:
		val = cpu.nextValue(AbsoluteY)
		cpu.ldx(val)

	case 0xa0:
		val = cpu.nextValue(Immediate)
		cpu.ldy(val)
	case 0xa4:
		val = cpu.nextValue(ZeroPage)
		cpu.ldy(val)
	case 0xb4:
		val = cpu.nextValue(ZeroPageX)
		cpu.ldy(val)
	case 0xac:
		val = cpu.nextValue(Absolute)
		cpu.ldy(val)
	case 0xbc:
		val = cpu.nextValue(AbsoluteX)
		cpu.ldy(val)

	case 0x85:
		addr = cpu.nextAddress(ZeroPage)
		cpu.sta(addr)
	case 0x95:
		addr = cpu.nextAddress(ZeroPageX)
		cpu.sta(addr)
	case 0x8d:
		addr = cpu.nextAddress(Absolute)
		cpu.sta(addr)
	case 0x9d:
		addr = cpu.nextAddress(AbsoluteX)
		cpu.sta(addr)
	case 0x99:
		addr = cpu.nextAddress(AbsoluteY)
		cpu.sta(addr)
	case 0x81:
		addr = cpu.nextAddress(IndirectX)
		cpu.sta(addr)
	case 0x91:
		addr = cpu.nextAddress(IndirectY)
		cpu.sta(addr)

	case 0x86:
		addr = cpu.nextAddress(ZeroPage)
		cpu.stx(addr)
	case 0x96:
		addr = cpu.nextAddress(ZeroPageY)
		cpu.stx(addr)
	case 0x8e:
		addr = cpu.nextAddress(Absolute)
		cpu.stx(addr)

	case 0x84:
		addr = cpu.nextAddress(ZeroPage)
		cpu.sty(addr)
	case 0x94:
		addr = cpu.nextAddress(ZeroPageY)
		cpu.sty(addr)
	case 0x8c:
		addr = cpu.nextAddress(Absolute)
		cpu.sty(addr)

	case 0xaa:
		cpu.tax()
	case 0xa8:
		cpu.tay()
	case 0xba:
		cpu.tsx()
	case 0x8a:
		cpu.txa()
	case 0x9a:
		cpu.txs()
	case 0x98:
		cpu.tya()

	case 0x48:
		cpu.pha()
	case 0x08:
		cpu.php()
	case 0x68:
		cpu.pla()
	case 0x28:
		cpu.plp()

	case 0x29:
		val = cpu.nextValue(Immediate)
		cpu.and(val)
	case 0x25:
		val = cpu.nextValue(ZeroPage)
		cpu.and(val)
	case 0x35:
		val = cpu.nextValue(ZeroPageX)
		cpu.and(val)
	case 0x2d:
		val = cpu.nextValue(Absolute)
		cpu.and(val)
	case 0x3d:
		val = cpu.nextValue(AbsoluteX)
		cpu.and(val)
	case 0x39:
		val = cpu.nextValue(AbsoluteY)
		cpu.and(val)
	case 0x21:
		val = cpu.nextValue(IndirectX)
		cpu.and(val)
	case 0x31:
		val = cpu.nextValue(IndirectY)
		cpu.and(val)

	case 0x49:
		val = cpu.nextValue(Immediate)
		cpu.eor(val)
	case 0x45:
		val = cpu.nextValue(ZeroPage)
		cpu.eor(val)
	case 0x55:
		val = cpu.nextValue(ZeroPageX)
		cpu.eor(val)
	case 0x4d:
		val = cpu.nextValue(Absolute)
		cpu.eor(val)
	case 0x5d:
		val = cpu.nextValue(AbsoluteX)
		cpu.eor(val)
	case 0x59:
		val = cpu.nextValue(AbsoluteY)
		cpu.eor(val)
	case 0x41:
		val = cpu.nextValue(IndirectX)
		cpu.eor(val)
	case 0x51:
		val = cpu.nextValue(IndirectY)
		cpu.eor(val)

	case 0x09:
		val = cpu.nextValue(Immediate)
		cpu.ora(val)
	case 0x05:
		val = cpu.nextValue(ZeroPage)
		cpu.ora(val)
	case 0x15:
		val = cpu.nextValue(ZeroPageX)
		cpu.ora(val)
	case 0x0d:
		val = cpu.nextValue(Absolute)
		cpu.ora(val)
	case 0x1d:
		val = cpu.nextValue(AbsoluteX)
		cpu.ora(val)
	case 0x19:
		val = cpu.nextValue(AbsoluteY)
		cpu.ora(val)
	case 0x01:
		val = cpu.nextValue(IndirectX)
		cpu.ora(val)
	case 0x11:
		val = cpu.nextValue(IndirectY)
		cpu.ora(val)

	case 0x24:
		val = cpu.nextValue(ZeroPage)
		cpu.bit(val)
	case 0x2c:
		val = cpu.nextValue(Absolute)
		cpu.bit(val)

	case 0x69:
		val = cpu.nextValue(Immediate)
		cpu.adc(val)
	case 0x65:
		val = cpu.nextValue(ZeroPage)
		cpu.adc(val)
	case 0x75:
		val = cpu.nextValue(ZeroPageX)
		cpu.adc(val)
	case 0x6d:
		val = cpu.nextValue(Absolute)
		cpu.adc(val)
	case 0x7d:
		val = cpu.nextValue(AbsoluteX)
		cpu.adc(val)
	case 0x79:
		val = cpu.nextValue(AbsoluteY)
		cpu.adc(val)
	case 0x61:
		val = cpu.nextValue(IndirectX)
		cpu.adc(val)
	case 0x71:
		val = cpu.nextValue(IndirectY)
		cpu.adc(val)

	case 0xe9:
		val = cpu.nextValue(Immediate)
		cpu.sbc(val)
	case 0xe5:
		val = cpu.nextValue(ZeroPage)
		cpu.sbc(val)
	case 0xf5:
		val = cpu.nextValue(ZeroPageX)
		cpu.sbc(val)
	case 0xed:
		val = cpu.nextValue(Absolute)
		cpu.sbc(val)
	case 0xfd:
		val = cpu.nextValue(AbsoluteX)
		cpu.sbc(val)
	case 0xf9:
		val = cpu.nextValue(AbsoluteY)
		cpu.sbc(val)
	case 0xe1:
		val = cpu.nextValue(IndirectX)
		cpu.sbc(val)
	case 0xf1:
		val = cpu.nextValue(IndirectY)
		cpu.sbc(val)

	case 0xc9:
		val = cpu.nextValue(Immediate)
		cpu.cmp(val)
	case 0xc5:
		val = cpu.nextValue(ZeroPage)
		cpu.cmp(val)
	case 0xd5:
		val = cpu.nextValue(ZeroPageX)
		cpu.cmp(val)
	case 0xcd:
		val = cpu.nextValue(Absolute)
		cpu.cmp(val)
	case 0xdd:
		val = cpu.nextValue(AbsoluteX)
		cpu.cmp(val)
	case 0xd9:
		val = cpu.nextValue(AbsoluteY)
		cpu.cmp(val)
	case 0xc1:
		val = cpu.nextValue(IndirectX)
		cpu.cmp(val)
	case 0xd1:
		val = cpu.nextValue(IndirectY)
		cpu.cmp(val)

	case 0xe0:
		val = cpu.nextValue(Immediate)
		cpu.cpx(val)
	case 0xe4:
		val = cpu.nextValue(ZeroPage)
		cpu.cpx(val)
	case 0xec:
		val = cpu.nextValue(Absolute)
		cpu.cpx(val)

	case 0xc0:
		val = cpu.nextValue(Immediate)
		cpu.cpy(val)
	case 0xc4:
		val = cpu.nextValue(ZeroPage)
		cpu.cpy(val)
	case 0xcc:
		val = cpu.nextValue(Absolute)
		cpu.cpy(val)

	case 0xe6:
		addr = cpu.nextAddress(ZeroPage)
		cpu.inc(addr)
	case 0xf6:
		addr = cpu.nextAddress(ZeroPageX)
		cpu.inc(addr)
	case 0xee:
		addr = cpu.nextAddress(Absolute)
		cpu.inc(addr)
	case 0xfe:
		addr = cpu.nextAddress(AbsoluteX)
		cpu.inc(addr)

	case 0xe8:
		cpu.inx()
	case 0xc8:
		cpu.iny()

	case 0xc6:
		addr = cpu.nextAddress(ZeroPage)
		cpu.dec(addr)
	case 0xd6:
		addr = cpu.nextAddress(ZeroPageX)
		cpu.dec(addr)
	case 0xce:
		addr = cpu.nextAddress(Absolute)
		cpu.dec(addr)
	case 0xde:
		addr = cpu.nextAddress(AbsoluteX)
		cpu.dec(addr)

	case 0xca:
		cpu.dex()
	case 0x88:
		cpu.dey()

	case 0x0a:
		cpu.asl_acc()
	case 0x06:
		addr = cpu.nextAddress(ZeroPage)
		cpu.asl(addr)
	case 0x16:
		addr = cpu.nextAddress(ZeroPageX)
		cpu.asl(addr)
	case 0x0e:
		addr = cpu.nextAddress(Absolute)
		cpu.asl(addr)
	case 0x1e:
		addr = cpu.nextAddress(AbsoluteX)
		cpu.asl(addr)

	case 0x4a:
		cpu.lsr_acc()
	case 0x46:
		addr = cpu.nextAddress(ZeroPage)
		cpu.lsr(addr)
	case 0x56:
		addr = cpu.nextAddress(ZeroPageX)
		cpu.lsr(addr)
	case 0x4e:
		addr = cpu.nextAddress(Absolute)
		cpu.lsr(addr)
	case 0x5e:
		addr = cpu.nextAddress(AbsoluteX)
		cpu.lsr(addr)

	case 0x2a:
		cpu.rol_acc()
	case 0x26:
		addr = cpu.nextAddress(ZeroPage)
		cpu.rol(addr)
	case 0x36:
		addr = cpu.nextAddress(ZeroPageX)
		cpu.rol(addr)
	case 0x2e:
		addr = cpu.nextAddress(Absolute)
		cpu.rol(addr)
	case 0x3e:
		addr = cpu.nextAddress(AbsoluteX)
		cpu.rol(addr)

	case 0x6a:
		cpu.ror_acc()
	case 0x66:
		addr = cpu.nextAddress(ZeroPage)
		cpu.ror(addr)
	case 0x76:
		addr = cpu.nextAddress(ZeroPageX)
		cpu.ror(addr)
	case 0x6e:
		addr = cpu.nextAddress(Absolute)
		cpu.ror(addr)
	case 0x7e:
		addr = cpu.nextAddress(AbsoluteX)
		cpu.ror(addr)

	case 0x4c:
		addr = cpu.nextAddress(Absolute)
		cpu.jmp(addr)
	case 0x6c:
		addr = cpu.nextAddress(Indirect)
		cpu.jmp(addr)

	case 0x20:
		addr = cpu.nextAddress(Absolute)
		cpu.jsr(addr)

	case 0x60:
		cpu.rts()

	case 0x90:
		val = cpu.nextValue(Immediate)
		cpu.bcc(val)
	case 0xb0:
		val = cpu.nextValue(Immediate)
		cpu.bcs(val)
	case 0xf0:
		val = cpu.nextValue(Immediate)
		cpu.beq(val)
	case 0x30:
		val = cpu.nextValue(Immediate)
		cpu.bmi(val)
	case 0xd0:
		val = cpu.nextValue(Immediate)
		cpu.bne(val)
	case 0x10:
		val = cpu.nextValue(Immediate)
		cpu.bpl(val)
	case 0x50:
		val = cpu.nextValue(Immediate)
		cpu.bvc(val)
	case 0x70:
		val = cpu.nextValue(Immediate)
		cpu.bvs(val)

	case 0x18:
		cpu.clc()
	case 0xd8:
		cpu.cld()
	case 0x58:
		cpu.cli()
	case 0xb8:
		cpu.clv()
	case 0x38:
		cpu.sec()
	case 0xf8:
		cpu.sed()
	case 0x78:
		cpu.sei()

	default:
		fmt.Printf("Unknown")
	}

	fmt.Print("\n")
}

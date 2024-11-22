package mos6502

import "fmt"

func (cpu *CPU) brk() {
	cpu.stackPushCurrentPc()
	cpu.stackPush(cpu.p)

	addr1, _ := cpu.mem.Read(0xfffe)
	addr2, _ := cpu.mem.Read(0xffff)

	cpu.pc = (uint16(addr1) << 8) + uint16(addr2)
	fmt.Printf("brk")
}

func (cpu *CPU) rti() {
	cpu.p = cpu.stackPull()
	cpu.pc = cpu.stackPullAddr()
	fmt.Printf("rti")
}

func (cpu *CPU) lda(am AdressingMode) {
	cpu.a = cpu.nextValue(am)
	cpu.assignBasicFlags(cpu.a)
	fmt.Printf("lda AdressingMode[%v] %x", am, cpu.a)
}

func (cpu *CPU) ldx(am AdressingMode) {
	cpu.x = cpu.nextValue(am)
	cpu.assignBasicFlags(cpu.x)
	fmt.Printf("ldx AdressingMode[%v] %v", am, cpu.x)
}

func (cpu *CPU) ldy(am AdressingMode) {
	cpu.y = cpu.nextValue(am)
	cpu.assignBasicFlags(cpu.y)
	fmt.Printf("ldy AdressingMode[%v] %v", am, cpu.y)
}

func (cpu *CPU) sta(am AdressingMode) {
	cpu.write(cpu.a, am)
	fmt.Printf("sta AdressingMode[%v] %v", am, cpu.a)
}

func (cpu *CPU) stx(am AdressingMode) {
	cpu.write(cpu.x, am)
	fmt.Printf("stx AdressingMode[%v] %v", am, cpu.x)
}

func (cpu *CPU) sty(am AdressingMode) {
	cpu.write(cpu.y, am)
	fmt.Printf("sty AdressingMode[%v] %v", am, cpu.y)
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
	fmt.Printf("pha")
}

func (cpu *CPU) php() {
	cpu.stackPush(cpu.p)
	fmt.Printf("php")
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

func (cpu *CPU) and(am AdressingMode) {
	cpu.a &= cpu.nextValue(am)
	cpu.assignBasicFlags(cpu.a)
	fmt.Printf("and AdressingMode[%v] %x", am, cpu.a)
}

func (cpu *CPU) eor(am AdressingMode) {
	cpu.a ^= cpu.nextValue(am)
	cpu.assignBasicFlags(cpu.a)
	fmt.Printf("eor AdressingMode[%v] %v", am, cpu.a)
}

func (cpu *CPU) ora(am AdressingMode) {
	cpu.a |= cpu.nextValue(am)
	cpu.assignBasicFlags(cpu.a)
	fmt.Printf("ora AdressingMode[%v] %v", am, cpu.a)
}

func (cpu *CPU) bit(am AdressingMode) {
	val := cpu.nextValue(am)
	cpu.bitTest(val)
	fmt.Printf("bit AdressingMode[%v] %v", am, val)
}

func (cpu *CPU) adc(am AdressingMode) {
	prev := cpu.a
	memory := cpu.nextValue(am)
	val1, overflow1 := addOverflow(cpu.a, memory)

	var carry byte
	if cpu.getFlag(FlagCarry) {
		carry = 1
	} else {
		carry = 0
	}

	val2, overflow2 := addOverflow(val1, carry)
	cpu.a = val2

	cpu.setFlag(FlagZero, cpu.a == 0)
	cpu.setFlag(FlagCarry, overflow1 || overflow2)
	cpu.setFlag(FlagNegative, isNegative(cpu.a))
	cpu.setFlag(FlagOverflow, ((cpu.a^prev)&(cpu.a^memory)&0x80) == 0x80)
	fmt.Printf("adc AdressingMode[%v] %x", am, cpu.a)
}

// TODO: Revisar
func (cpu *CPU) sbc(am AdressingMode) {
	prev := cpu.a
	memory := cpu.nextValue(am)
	val1, overflow1 := subOverflow(cpu.a, memory)

	var carry byte
	if !cpu.getFlag(FlagCarry) {
		carry = 1
	}

	val2, overflow2 := subOverflow(val1, carry)

	cpu.a = val2

	cpu.setFlag(FlagCarry, !(overflow1 || overflow2))
	cpu.setFlag(FlagZero, cpu.a == 0)
	cpu.setFlag(FlagNegative, isNegative(cpu.a))
	cpu.setFlag(FlagOverflow, ((cpu.a^prev)&(cpu.a^(^memory))&0x80) == 0x80)
	fmt.Printf("sbc AdressingMode[%v] %v", am, cpu.a)
}

func (cpu *CPU) cmp(am AdressingMode) {
	val := cpu.nextValue(am)
	cpu.compare(cpu.a, val)
	fmt.Printf("cmp AdressingMode[%v] %x, %x", am, cpu.a, val)
}

func (cpu *CPU) cpx(am AdressingMode) {
	val := cpu.nextValue(am)
	cpu.compare(cpu.x, val)
	fmt.Printf("cpx AdressingMode[%v] %v, %v", am, cpu.x, val)
}

func (cpu *CPU) cpy(am AdressingMode) {
	val := cpu.nextValue(am)
	cpu.compare(cpu.y, val)
	fmt.Printf("cpy AdressingMode[%v] %v, %v", am, cpu.y, val)
}

func (cpu *CPU) inc(am AdressingMode) {
	pc_temp := cpu.pc
	val := cpu.nextValue(am) + 1
	cpu.assignBasicFlags(val)
	cpu.pc = pc_temp
	cpu.write(val, am)
	fmt.Printf("inc AdressingMode[%v] %v", am, val-1)
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

func (cpu *CPU) dec(am AdressingMode) {
	pc_temp := cpu.pc
	val := cpu.nextValue(am) - 1
	cpu.assignBasicFlags(val)
	cpu.pc -= pc_temp
	cpu.write(val, am)
	fmt.Printf("dec AdressingMode[%v] %v", am, val+1)
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

func (cpu *CPU) asl(am AdressingMode) {
	pc_temp := cpu.pc
	val := cpu.nextValue(am)
	cpu.pc -= pc_temp
	cpu.setFlag(FlagCarry, val&0x80 == 0x80)
	val = val << 1
	cpu.assignBasicFlags(val)
	cpu.write(val, am)
	fmt.Printf("asl AdressingMode[%v] %v", am, val)
}

func (cpu *CPU) lsr_acc() {
	cpu.setFlag(FlagCarry, cpu.a&0x01 == 0x01)
	cpu.a = cpu.a >> 1
	cpu.assignBasicFlags(cpu.a)
	fmt.Printf("lsr acc %v", cpu.a)
}

func (cpu *CPU) lsr(am AdressingMode) {
	pc_temp := cpu.pc
	val := cpu.nextValue(am)
	cpu.pc -= pc_temp
	cpu.setFlag(FlagCarry, val&0x01 == 0x01)
	val = val >> 1
	cpu.assignBasicFlags(val)
	cpu.write(val, am)
	fmt.Printf("lsr AdressingMode[%v] %v", am, val)
}

func (cpu *CPU) rol_acc() {
	cpu.setFlag(FlagCarry, cpu.a&0x80 == 0x80)
	cpu.a = cpu.a << 1
	if cpu.getFlag(FlagCarry) {
		cpu.a += 1
	}
	cpu.assignBasicFlags(cpu.a)
	fmt.Printf("rol acc %v", cpu.a)
}

func (cpu *CPU) rol(am AdressingMode) {
	pc_temp := cpu.pc
	val := cpu.nextValue(am)
	cpu.pc -= pc_temp
	cpu.setFlag(FlagCarry, val&0x80 == 0x80)
	val = val << 1
	if cpu.getFlag(FlagCarry) {
		cpu.a += 1
	}
	cpu.assignBasicFlags(val)
	cpu.write(val, am)
	fmt.Printf("rol AdressingMode[%v] %v", am, val)
}

func (cpu *CPU) ror_acc() {
	cpu.setFlag(FlagCarry, cpu.a&0x01 == 0x01)
	cpu.a = cpu.a >> 1
	if cpu.getFlag(FlagCarry) {
		cpu.a += 0x80
	}
	cpu.assignBasicFlags(cpu.a)
	fmt.Printf("ror acc %v", cpu.a)
}

func (cpu *CPU) ror(am AdressingMode) {
	pc_temp := cpu.pc
	val := cpu.nextValue(am)
	cpu.pc -= pc_temp
	cpu.setFlag(FlagCarry, val&0x01 == 0x01)
	val = val << 1
	if cpu.getFlag(FlagCarry) {
		cpu.a += 1
	}
	cpu.assignBasicFlags(val)
	cpu.write(val, am)
	fmt.Printf("ror AdressingMode[%v] %v", am, val)
}

func (cpu *CPU) jmp_absolute() {
	addr := cpu.nextAddr()
	cpu.pc = addr
	fmt.Printf("jmp absolute %x", addr)
}

func (cpu *CPU) jmp_indirect() {
	addr := cpu.nextAddrIndirect()
	cpu.pc = addr
	fmt.Printf("jmp indirect %x", addr)
}

func (cpu *CPU) jsr() {
	addr := cpu.nextAddr()
	cpu.stackPushCurrentPc()
	cpu.pc = addr
	fmt.Printf("jsr %x", addr)
}

func (cpu *CPU) rts() {
	cpu.pc = cpu.stackPullAddr()
	fmt.Printf("rts %x", cpu.pc)
}

// TODO: comprobar si funciona bien. PC + 1?
func (cpu *CPU) bcc() {
	displacement := int8(cpu.nextValue(Immediate))
	if !cpu.getFlag(FlagCarry) {
		cpu.branchJump(displacement)
	}
	fmt.Printf("bcc %v", displacement)
}

func (cpu *CPU) bcs() {
	displacement := int8(cpu.nextValue(Immediate))
	if cpu.getFlag(FlagCarry) {
		cpu.branchJump(displacement)
	}
	fmt.Printf("bcs %v", displacement)
}

func (cpu *CPU) beq() {
	displacement := int8(cpu.nextValue(Immediate))
	if cpu.getFlag(FlagZero) {
		cpu.branchJump(displacement)
	}
	fmt.Printf("beq %v Zero flag[%v]", displacement, cpu.getFlag(FlagZero))
}

func (cpu *CPU) bmi() {
	displacement := int8(cpu.nextValue(Immediate))
	if cpu.getFlag(FlagNegative) {
		cpu.branchJump(displacement)
	}
	fmt.Printf("bmi %v", displacement)
}

func (cpu *CPU) bne() {
	displacement := int8(cpu.nextValue(Immediate))
	if !cpu.getFlag(FlagZero) {
		cpu.branchJump(displacement)
	}
	fmt.Printf("bne %v", displacement)
}

func (cpu *CPU) bpl() {
	displacement := int8(cpu.nextValue(Immediate))
	if !cpu.getFlag(FlagNegative) {
		cpu.branchJump(displacement)
	}
	fmt.Printf("bpl %v", displacement)
}

func (cpu *CPU) bvc() {
	displacement := int8(cpu.nextValue(Immediate))
	if !cpu.getFlag(FlagOverflow) {
		cpu.branchJump(displacement)
	}
	fmt.Printf("bvc %v", displacement)
}

func (cpu *CPU) bvs() {
	displacement := int8(cpu.nextValue(Immediate))
	if cpu.getFlag(FlagOverflow) {
		cpu.branchJump(displacement)
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
	fmt.Printf("[PC: %4x] OPCODE %2x | ", cpu.pc-1, opcode)

	switch opcode {
	case 0x00:
		cpu.brk()

	case 0x40:
		cpu.rti()

	case 0xea:
		fmt.Printf("nop")

	case 0xa9:
		cpu.lda(Immediate)
	case 0xa5:
		cpu.lda(ZeroPage)
	case 0xb5:
		cpu.lda(ZeroPageX)
	case 0xad:
		cpu.lda(Absolute)
	case 0xbd:
		cpu.lda(AbsoluteX)
	case 0xb9:
		cpu.lda(AbsoluteY)
	case 0xa1:
		cpu.lda(IndirectX)
	case 0xb1:
		cpu.lda(IndirectY)

	case 0xa2:
		cpu.ldx(Immediate)
	case 0xa6:
		cpu.ldx(ZeroPage)
	case 0xb6:
		cpu.ldx(ZeroPageY)
	case 0xae:
		cpu.ldx(Absolute)
	case 0xbe:
		cpu.ldx(AbsoluteY)

	case 0xa0:
		cpu.ldy(Immediate)
	case 0xa4:
		cpu.ldy(ZeroPage)
	case 0xb4:
		cpu.ldy(ZeroPageX)
	case 0xac:
		cpu.ldy(Absolute)
	case 0xbc:
		cpu.ldy(AbsoluteX)

	case 0x85:
		cpu.sta(ZeroPage)
	case 0x95:
		cpu.sta(ZeroPageX)
	case 0x8d:
		cpu.sta(Absolute)
	case 0x9d:
		cpu.sta(AbsoluteX)
	case 0x99:
		cpu.sta(AbsoluteY)
	case 0x81:
		cpu.sta(IndirectX)
	case 0x91:
		cpu.sta(IndirectY)

	case 0x86:
		cpu.stx(ZeroPage)
	case 0x96:
		cpu.stx(ZeroPageY)
	case 0x8e:
		cpu.stx(Absolute)

	case 0x84:
		cpu.sty(ZeroPage)
	case 0x94:
		cpu.sty(ZeroPageY)
	case 0x8c:
		cpu.sty(Absolute)

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
		cpu.and(Immediate)
	case 0x25:
		cpu.and(ZeroPage)
	case 0x35:
		cpu.and(ZeroPageX)
	case 0x2d:
		cpu.and(Absolute)
	case 0x3d:
		cpu.and(AbsoluteX)
	case 0x39:
		cpu.and(AbsoluteY)
	case 0x21:
		cpu.and(IndirectX)
	case 0x31:
		cpu.and(IndirectY)

	case 0x49:
		cpu.eor(Immediate)
	case 0x45:
		cpu.eor(ZeroPage)
	case 0x55:
		cpu.eor(ZeroPageX)
	case 0x4d:
		cpu.eor(Absolute)
	case 0x5d:
		cpu.eor(AbsoluteX)
	case 0x59:
		cpu.eor(AbsoluteY)
	case 0x41:
		cpu.eor(IndirectX)
	case 0x51:
		cpu.eor(IndirectY)

	case 0x09:
		cpu.ora(Immediate)
	case 0x05:
		cpu.ora(ZeroPage)
	case 0x15:
		cpu.ora(ZeroPageX)
	case 0x0d:
		cpu.ora(Absolute)
	case 0x1d:
		cpu.ora(AbsoluteX)
	case 0x19:
		cpu.ora(AbsoluteY)
	case 0x01:
		cpu.ora(IndirectX)
	case 0x11:
		cpu.ora(IndirectY)

	case 0x24:
		cpu.bit(ZeroPage)
	case 0x2c:
		cpu.bit(Absolute)

	case 0x69:
		cpu.adc(Immediate)
	case 0x65:
		cpu.adc(ZeroPage)
	case 0x75:
		cpu.adc(ZeroPageX)
	case 0x6d:
		cpu.adc(Absolute)
	case 0x7d:
		cpu.adc(AbsoluteX)
	case 0x79:
		cpu.adc(AbsoluteY)
	case 0x61:
		cpu.adc(IndirectX)
	case 0x71:
		cpu.adc(IndirectY)

	case 0xe9:
		cpu.sbc(Immediate)
	case 0xe5:
		cpu.sbc(ZeroPage)
	case 0xf5:
		cpu.sbc(ZeroPageX)
	case 0xed:
		cpu.sbc(Absolute)
	case 0xfd:
		cpu.sbc(AbsoluteX)
	case 0xf9:
		cpu.sbc(AbsoluteY)
	case 0xe1:
		cpu.sbc(IndirectX)
	case 0xf1:
		cpu.sbc(IndirectY)

	case 0xc9:
		cpu.cmp(Immediate)
	case 0xc5:
		cpu.cmp(ZeroPage)
	case 0xd5:
		cpu.cmp(ZeroPageX)
	case 0xcd:
		cpu.cmp(Absolute)
	case 0xdd:
		cpu.cmp(AbsoluteX)
	case 0xd9:
		cpu.cmp(AbsoluteY)
	case 0xc1:
		cpu.cmp(IndirectX)
	case 0xd1:
		cpu.cmp(IndirectY)

	case 0xe0:
		cpu.cpx(Immediate)
	case 0xe4:
		cpu.cpx(ZeroPage)
	case 0xec:
		cpu.cpx(Absolute)

	case 0xc0:
		cpu.cpy(Immediate)
	case 0xc4:
		cpu.cpy(ZeroPage)
	case 0xcc:
		cpu.cpy(Absolute)

	case 0xe6:
		cpu.inc(ZeroPage)
	case 0xf6:
		cpu.inc(ZeroPageX)
	case 0xee:
		cpu.inc(Absolute)
	case 0xfe:
		cpu.inc(AbsoluteX)

	case 0xe8:
		cpu.inx()
	case 0xc8:
		cpu.iny()

	case 0xc6:
		cpu.dec(ZeroPage)
	case 0xd6:
		cpu.dec(ZeroPageX)
	case 0xce:
		cpu.dec(Absolute)
	case 0xde:
		cpu.dec(AbsoluteX)

	case 0xca:
		cpu.dex()
	case 0x88:
		cpu.dey()

	case 0x0a:
		cpu.asl_acc()
	case 0x06:
		cpu.asl(ZeroPage)
	case 0x16:
		cpu.asl(ZeroPageX)
	case 0x0e:
		cpu.asl(Absolute)
	case 0x1e:
		cpu.asl(AbsoluteX)

	case 0x4a:
		cpu.lsr_acc()
	case 0x46:
		cpu.lsr(ZeroPage)
	case 0x56:
		cpu.lsr(ZeroPageX)
	case 0x4e:
		cpu.lsr(Absolute)
	case 0x5e:
		cpu.lsr(AbsoluteX)

	case 0x2a:
		cpu.rol_acc()
	case 0x26:
		cpu.rol(ZeroPage)
	case 0x36:
		cpu.rol(ZeroPageX)
	case 0x2e:
		cpu.rol(Absolute)
	case 0x3e:
		cpu.rol(AbsoluteX)

	case 0x6a:
		cpu.ror_acc()
	case 0x66:
		cpu.ror(ZeroPage)
	case 0x76:
		cpu.ror(ZeroPageX)
	case 0x6e:
		cpu.ror(Absolute)
	case 0x7e:
		cpu.ror(AbsoluteX)

	case 0x4c:
		cpu.jmp_absolute()
	case 0x6c:
		cpu.jmp_indirect()

	case 0x20:
		cpu.jsr()

	case 0x60:
		cpu.rts()

	case 0x90:
		cpu.bcc()
	case 0xb0:
		cpu.bcs()
	case 0xf0:
		cpu.beq()
	case 0x30:
		cpu.bmi()
	case 0xd0:
		cpu.bne()
	case 0x10:
		cpu.bpl()
	case 0x50:
		cpu.bvc()
	case 0x70:
		cpu.bvs()

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

	cpu.PrintState()
}

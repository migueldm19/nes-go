package main

import (
	"flag"
	"log"
	"nes-go/mos6502"
	"nes-go/ppu"
	"os"
)

func main() {
	disassemble_activated := flag.Bool("disassemble", false, "Run disassembler")
	flag.Parse()

	flag_tail := flag.Args()

	var rom_path string

	if len(flag_tail) == 0 {
		rom_path = "nestest.nes"
	} else {
		rom_path = flag_tail[0]
	}

	cart, err := os.ReadFile(rom_path)

	if err != nil {
		log.Fatalf("Error reading cartridge: %v", err)
	}

	rom := mos6502.NewRom(cart)
	memory := mos6502.NewMemory(rom)
	cpu := mos6502.NewCPU(memory)

	mppu := ppu.NewPPU(memory)
	pt0 := mppu.GetPatternTable0()
	pt1 := mppu.GetPatternTable1()

	ppu.GenerateImage("pt0.png", pt0)
	ppu.GenerateImage("pt1.png", pt1)

	if *disassemble_activated {
		disassembler := mos6502.NewDisassembler(cpu)
		disassembler.DisassembleWeb()
	} else {
		cpu.Run()
	}
}

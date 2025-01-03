package main

import (
	"flag"
	"log"
	"nes-go/mos6502"
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
	cpu := mos6502.NewCPU(rom)

	if *disassemble_activated {
		disassembler := mos6502.NewDisassembler(cpu)
		disassembler.DisassembleWeb()
	} else {
		cpu.Run()
	}
}

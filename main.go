package main

import (
	"log"
	"nes-go/mos6502"
	"os"
)

func main() {
	cart, err := os.ReadFile("nestest.nes")

	if err != nil {
		log.Fatalf("Error reading cartridge: %v", err)
	}

	rom := mos6502.NewRom(cart)
	cpu := mos6502.NewCPU(rom)
	cpu.Run()
}

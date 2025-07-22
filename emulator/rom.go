package emulator

import "log"

const (
	HEADER_SIZE        = 16
	PRG_BYTES_UNITS    = 16
	CHR_BYTES_UNITS    = 8
	BYTES_IN_KILOBYTES = 1024
	TRAINER_SIZE       = 512
)

type NametableArrangement byte

const (
	VERTICAL NametableArrangement = iota
	HORIZONTAL
)

type Rom struct {
	PrgRomSize uint16
	ChrRomSize uint16
	PrgData    []byte
	ChrData    []byte
	Trainer    []byte

	NtArrangement NametableArrangement
	HasPrgRam     bool
}

func getBit(val byte, idx int) bool {
	return ((val >> idx) & 1) != 0
}

func NewRom(cartridge []byte) *Rom {
	header := cartridge[:HEADER_SIZE]
	var prgSize uint16 = uint16(header[4]) * PRG_BYTES_UNITS * BYTES_IN_KILOBYTES
	var chrSize uint16 = uint16(header[5]) * CHR_BYTES_UNITS * BYTES_IN_KILOBYTES

	var ntArrangement NametableArrangement
	flags6 := header[6]
	if getBit(flags6, 0) {
		ntArrangement = HORIZONTAL
	} else {
		ntArrangement = VERTICAL
	}

	startPrg := HEADER_SIZE
	var trainer []byte
	if getBit(flags6, 2) {
		startPrg += TRAINER_SIZE
		trainer = cartridge[HEADER_SIZE : HEADER_SIZE+TRAINER_SIZE]
	}

	startChr := HEADER_SIZE + prgSize

	prgData := cartridge[HEADER_SIZE:startChr]

	if prgSize >= 0x4000 {
		prgData = append(prgData, prgData...)
	}

	chrData := cartridge[startChr : startChr+chrSize]

	if len(chrData) != 0x2000 {
		log.Printf("[Warning] Chr data should be 0x2000 bytes long (len: %v)\n", len(chrData))
	}

	return &Rom{
		PrgRomSize:    prgSize,
		ChrRomSize:    chrSize,
		PrgData:       prgData,
		ChrData:       chrData,
		Trainer:       trainer,
		NtArrangement: ntArrangement,
		HasPrgRam:     getBit(flags6, 1),
	}
}

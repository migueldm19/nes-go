package emulator

import "log"

const(
	HEADER_SIZE = 16
	PRG_BYTES_UNITS = 16
	CHR_BYTES_UNITS = 8
	BYTES_IN_KILOBYTES = 1024
)

type Rom struct {
	PrgRomSize uint16
	ChrRomSize uint16
	PrgData    []byte
	ChrData    []byte
}

func NewRom(cartridge []byte) *Rom {
	header := cartridge[:HEADER_SIZE]
	var prgSize uint16 = uint16(header[4]) * PRG_BYTES_UNITS * BYTES_IN_KILOBYTES
	var chrSize uint16 = uint16(header[5]) * CHR_BYTES_UNITS * BYTES_IN_KILOBYTES

	// TODO: read header and add trainer

	start_chr := HEADER_SIZE + prgSize

	prgData := cartridge[HEADER_SIZE:start_chr]

	if prgSize >= 0x4000 {
		prgData = append(prgData, prgData...)
	}

	chrData := cartridge[start_chr : start_chr+chrSize]

	if len(chrData) != 0x2000 {
		log.Printf("[Warning] Chr data should be 0x2000 bytes long (len: %v)\n", len(chrData))
	}

	return &Rom{
		PrgRomSize: prgSize,
		ChrRomSize: chrSize,
		PrgData:    prgData,
		ChrData:    chrData,
	}
}

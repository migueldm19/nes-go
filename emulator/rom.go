package emulator

import "fmt"

type Rom struct {
	PrgRomSize uint16
	ChrRomSize uint16
	PrgData    []byte
	ChrData    []byte
}

func NewRom(cartridge []byte) *Rom {
	header := cartridge[:16]
	var prgSize uint16 = uint16(header[4]) * 16384
	var chrSize uint16 = uint16(header[5]) * 8192

	start_chr := 16 + prgSize

	prgData := cartridge[16:start_chr]

	if prgSize >= 0x4000 {
		prgData = append(prgData, prgData...)
	}

	chrData := cartridge[start_chr : start_chr+chrSize]

	if len(chrData) != 0x2000 {
		panic(fmt.Sprintf("Chr data should be 0x2000 bytes long (len: %v)", len(chrData)))
	}

	return &Rom{
		PrgRomSize: prgSize,
		ChrRomSize: chrSize,
		PrgData:    prgData,
		ChrData:    chrData,
	}
}

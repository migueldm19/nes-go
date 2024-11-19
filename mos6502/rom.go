package mos6502

type Rom struct {
	PrgRomSize uint16
	ChrRomSize uint16
	Data       []byte
}

func NewRom(cartridge []byte) *Rom {
	header := cartridge[:16]
	var prgSize uint16 = uint16(header[4]) * 16384
	var chrSize uint16 = uint16(header[5]) * 8192

	prgData := cartridge[16 : 16+prgSize]

	if prgSize >= 0x4000 {
		prgData = append(prgData, prgData...)
	}

	return &Rom{
		PrgRomSize: prgSize,
		ChrRomSize: chrSize,
		Data:       prgData,
	}
}

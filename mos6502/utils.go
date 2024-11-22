package mos6502

import "math"

// TODO: probar funcionamiento
func addOverflow(a, b byte) (byte, bool) {
	sum := a + b
	return sum, a > math.MaxUint8-b
}

// TODO: probar funcionamiento
func subOverflow(a, b byte) (byte, bool) {
	sub := a - b
	return sub, a < b
}

func isNegative(val byte) bool {
	return (val & 0x80) == 0x80
}

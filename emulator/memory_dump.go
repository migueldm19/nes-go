package emulator

import (
	"fmt"
)

type Dump map[int]string

type MemoryDump struct {
	ZeroPage Dump
	Stack    Dump
}

func NewMemoryDump(mem *Memory) *MemoryDump {
	return &MemoryDump{
		ZeroPage: mem.ZeroPageDump(),
		Stack:    mem.StackDump(),
	}
}

func (dump Dump) String() string {
	buff := ""

	for addr, data := range dump {
		buff += fmt.Sprintf("[%04X] %v\n", addr, data)
	}

	return buff
}

func (dump MemoryDump) String() string {
	buff := ""

	buff += fmt.Sprintf("%v", dump.ZeroPage)
	buff += fmt.Sprintf("%v", dump.Stack)

	return buff
}

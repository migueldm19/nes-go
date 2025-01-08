package mos6502

import (
	"fmt"
	"log"
	"os"
)

var _instructions_logger *log.Logger
var _memory_dump_logger *log.Logger
var _disassembly_logger *log.Logger

func createLogger(name string) *log.Logger {
	file, err := os.OpenFile(fmt.Sprintf("%v.log", name), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("Couldn't create %v log file: %v. Using stdout.\n", name, err)
		file = os.Stdout
	}
	return log.New(file, "", 0)
}

func GetInstructionsLogger() *log.Logger {
	if _instructions_logger == nil {
		_instructions_logger = createLogger("instructions")
	}

	return _instructions_logger
}

func GetMemoryDumpLogger() *log.Logger {
	if _memory_dump_logger == nil {
		_memory_dump_logger = createLogger("memory_dump")
	}

	return _memory_dump_logger
}

func GetDisassemblyLogger() *log.Logger {
	if _disassembly_logger == nil {
		_disassembly_logger = createLogger("disassembly")
	}

	return _disassembly_logger
}

package mos6502

import (
	"fmt"
	"log"
	"os"
)

type LoggerMos6502 struct {
	Instructions *log.Logger
	MemoryDump   *log.Logger
}

var logger *LoggerMos6502

func createLogger(name string) *log.Logger {
	file, err := os.OpenFile(fmt.Sprintf("%v.log", name), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("Couldn't create %v log file: %v. Using stdout.\n", name, err)
		file = os.Stdout
	}
	return log.New(file, "", 0)
}

func GetLogger() *LoggerMos6502 {
	if logger == nil {
		logger = &LoggerMos6502{}

		logger.Instructions = createLogger("instructions")
		logger.MemoryDump = createLogger("memory_dump")
	}

	return logger
}

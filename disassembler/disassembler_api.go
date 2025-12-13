package disassembler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (disassembler *Disassembler) GetInstructions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	instructionsData := InstructionsData{
		Instructions: disassembler.Instructions,
		Pc:           disassembler.Cpu.Pc,
	}
	json.NewEncoder(w).Encode(instructionsData)
}

func (disassembler *Disassembler) StepHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentInstruction, got := disassembler.Instructions[disassembler.Cpu.Pc]
	if !got {
		currentInstruction = disassembler.Cpu.GetNextInstruction()
	}
	disassembler.Cpu.Pc = currentInstruction.Pc + 1

	fmt.Println(currentInstruction)

	currentInstruction.Run(disassembler.Cpu)

	w.WriteHeader(http.StatusOK)
}

func (disassembler *Disassembler) GetCpuState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cpuStateData := disassembler.Cpu.GetStateData()
	json.NewEncoder(w).Encode(cpuStateData)
}

func serveStaticSite(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "disassembler/assets/new-disassembler.html")
}

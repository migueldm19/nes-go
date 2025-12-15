package disassembler

import (
	"encoding/json"
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

	disassembler.Step()

	w.WriteHeader(http.StatusOK)
}

func (disassembler *Disassembler) GetCpuState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cpuStateData := disassembler.Cpu.GetStateData()
	json.NewEncoder(w).Encode(cpuStateData)
}

func (disassembler *Disassembler) GetMemoryDump(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dump := disassembler.Cpu.Dump()
	json.NewEncoder(w).Encode(dump)
}

func (disassembler *Disassembler) ContinueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		Breakpoints []uint16 `json:"breakpoints"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	disassembler.Step()

outerLoop:
	for {
		for _, bp := range requestData.Breakpoints {
			if disassembler.Cpu.Pc == bp {
				break outerLoop
			}
		}

		disassembler.Step()
	}

	w.WriteHeader(http.StatusOK)
}

func serveStaticSite(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "disassembler/assets/disassembler.html")
}

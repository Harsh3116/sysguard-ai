package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type ShutdownAnalysis struct {
	Cause          string `json:"cause"`
	Severity       string `json:"severity"`
	Explanation    string `json:"explanation"`
	Recommendation string `json:"recommendation"`
}

func withCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		handler(w, r)
	}
}

func shutdownAnalysisHandler(w http.ResponseWriter, r *http.Request) {
	// PowerShell command to fetch last critical shutdown
	cmd := exec.Command(
		"powershell",
		"-Command",
		"Get-WinEvent -FilterHashtable @{LogName='System'; Id=41} -MaxEvents 1 | Select-Object -ExpandProperty Message",
	)

	output, err := cmd.Output()
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Unable to read system shutdown logs",
		})
		return
	}

	message := strings.ToLower(string(output))

	analysis := ShutdownAnalysis{
		Cause:          "Unknown",
		Severity:       "Low",
		Explanation:    "No critical shutdown detected recently.",
		Recommendation: "No action required.",
	}

	if strings.Contains(message, "power") {
		analysis = ShutdownAnalysis{
			Cause:          "Unexpected power loss",
			Severity:       "High",
			Explanation:    "The system shut down due to a sudden power interruption or forced shutdown.",
			Recommendation: "Check power adapter and battery health.",
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analysis)
}

func main() {
	http.HandleFunc("/shutdown-analysis", withCORS(shutdownAnalysisHandler))

	log.Println("SysGuard Agent running on http://127.0.0.1:7878")
	log.Fatal(http.ListenAndServe("127.0.0.1:7878", nil))
}

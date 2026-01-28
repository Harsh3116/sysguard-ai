package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

/* =======================
   DATA STRUCTURES
======================= */

type HealthScoreResponse struct {
	Score   int      `json:"score"`
	Status  string   `json:"status"`
	Reasons []string `json:"reasons"`
}

type ShutdownAnalysis struct {
	Cause          string `json:"cause"`
	Severity       string `json:"severity"`
	Explanation    string `json:"explanation"`
	Recommendation string `json:"recommendation"`
}

type StartupItem struct {
	Name   string `json:"name"`
	Impact string `json:"impact"`
}

type StartupSummary struct {
	OverallStatus  string   `json:"overall_status"`
	Message        string   `json:"message"`
	HighImpactApps []string `json:"high_impact_apps"`
	Recommendation string   `json:"recommendation"`
}

/* =======================
   COMMON MIDDLEWARE
======================= */

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

/* =======================
   HEALTH SCORE
======================= */

func healthScoreHandler(w http.ResponseWriter, r *http.Request) {
	score := 100
	reasons := []string{}

	cpuPercents, _ := cpu.Percent(time.Second, false)
	if cpuPercents[0] > 80 {
		score -= 20
		reasons = append(reasons, "High CPU usage")
	}

	memStat, _ := mem.VirtualMemory()
	if memStat.UsedPercent > 80 {
		score -= 20
		reasons = append(reasons, "High memory usage")
	}

	diskStat, _ := disk.Usage("C:\\")
	if diskStat.UsedPercent > 85 {
		score -= 20
		reasons = append(reasons, "Low disk space")
	}

	status := "Healthy"
	if score < 70 {
		status = "Critical"
	} else if score < 85 {
		status = "Warning"
	}

	json.NewEncoder(w).Encode(HealthScoreResponse{
		Score:   score,
		Status:  status,
		Reasons: reasons,
	})
}

/* =======================
   SHUTDOWN ANALYSIS
======================= */

func shutdownAnalysisHandler(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command(
		"powershell",
		"-Command",
		"Get-WinEvent -FilterHashtable @{LogName='System'; Id=41} -MaxEvents 1 | Select-Object -ExpandProperty Message",
	)

	output, _ := cmd.Output()
	message := strings.ToLower(string(output))

	result := ShutdownAnalysis{
		Cause:          "Unknown",
		Severity:       "Low",
		Explanation:    "No critical shutdown detected recently.",
		Recommendation: "No action required.",
	}

	if strings.Contains(message, "power") {
		result = ShutdownAnalysis{
			Cause:          "Unexpected power loss",
			Severity:       "High",
			Explanation:    "System shut down due to power failure or forced shutdown.",
			Recommendation: "Check power supply and battery health.",
		}
	}

	json.NewEncoder(w).Encode(result)
}

/* =======================
   STARTUP ANALYSIS
======================= */

func classifyImpact(name string) string {
	n := strings.ToLower(name)

	if strings.Contains(n, "steam") ||
		strings.Contains(n, "epic") ||
		strings.Contains(n, "discord") ||
		strings.Contains(n, "spotify") ||
		strings.Contains(n, "adobe") {
		return "High"
	}

	if strings.Contains(n, "onedrive") ||
		strings.Contains(n, "dropbox") ||
		strings.Contains(n, "helper") ||
		strings.Contains(n, "assistant") ||
		strings.Contains(n, "updater") {
		return "Medium"
	}

	return "Low"
}

func startupSummaryHandler(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command(
		"powershell",
		"-Command",
		"Get-CimInstance Win32_StartupCommand | Select-Object Name",
	)

	output, err := cmd.Output()
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Unable to read startup programs",
		})
		return
	}

	lines := strings.Split(string(output), "\n")
	highImpact := []string{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "Name") {
			continue
		}

		if classifyImpact(line) == "High" {
			highImpact = append(highImpact, line)
		}
	}

	summary := StartupSummary{
		OverallStatus:  "Good",
		Message:        "No high-impact startup programs detected.",
		HighImpactApps: highImpact,
		Recommendation: "No action required.",
	}

	if len(highImpact) > 0 {
		summary = StartupSummary{
			OverallStatus:  "Poor",
			Message:        "High-impact startup programs detected.",
			HighImpactApps: highImpact,
			Recommendation: "Disable high-impact apps to improve boot time.",
		}
	}

	json.NewEncoder(w).Encode(summary)
}

/* =======================
   MAIN
======================= */

func main() {
	http.HandleFunc("/health-score", withCORS(healthScoreHandler))
	http.HandleFunc("/shutdown-analysis", withCORS(shutdownAnalysisHandler))
	http.HandleFunc("/startup-summary", withCORS(startupSummaryHandler))

	log.Println("SysGuard Agent running on http://127.0.0.1:7878")
	log.Fatal(http.ListenAndServe("127.0.0.1:7878", nil))
}

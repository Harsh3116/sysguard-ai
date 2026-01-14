package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type MetricsResponse struct {
	CPUPercent  float64 `json:"cpu_percent"`
	MemoryUsed  uint64  `json:"memory_used"`
	MemoryTotal uint64  `json:"memory_total"`
	DiskUsed    uint64  `json:"disk_used"`
	DiskTotal   uint64  `json:"disk_total"`
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

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	cpuPercents, _ := cpu.Percent(time.Second, false)
	memStat, _ := mem.VirtualMemory()
	diskStat, _ := disk.Usage("C:\\")

	resp := MetricsResponse{
		CPUPercent:  cpuPercents[0],
		MemoryUsed:  memStat.Used,
		MemoryTotal: memStat.Total,
		DiskUsed:    diskStat.Used,
		DiskTotal:   diskStat.Total,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/metrics", withCORS(metricsHandler))

	log.Println("SysGuard Agent running on http://127.0.0.1:7878")
	log.Fatal(http.ListenAndServe("127.0.0.1:7878", nil))
}

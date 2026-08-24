// Package server exposes buoy-qc analysis via HTTP/JSON endpoints.
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"buoy-qc/internal/export"
	"buoy-qc/internal/obs"
	"buoy-qc/internal/qc"
	"buoy-qc/internal/sea"
)

// Config holds server configuration.
type Config struct {
	Addr string
}

// New creates a configured http.ServeMux with all routes registered.
func New(cfg Config) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/analyze", handleAnalyze)
	mux.HandleFunc("/api/qc", handleQC)
	mux.HandleFunc("/api/seastate", handleSeaState)
	mux.HandleFunc("/api/report", handleReport)
	return mux
}

// ListenAndServe starts the HTTP server.
func ListenAndServe(cfg Config) error {
	mux := New(cfg)
	return http.ListenAndServe(cfg.Addr, mux)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// readingInput is the JSON representation of a buoy reading for API input.
type readingInput struct {
	Time       string  `json:"time"`
	Buoy       string  `json:"buoy"`
	WindSpd    float64 `json:"windspd"`
	WindDir    float64 `json:"winddir"`
	AirTemp    float64 `json:"airtemp"`
	Pressure   float64 `json:"pressure"`
	WaveHt     float64 `json:"waveht"`
	WavePer    float64 `json:"waveper"`
	SST        float64 `json:"sst"`
	Salinity   float64 `json:"salinity"`
	CurrentSpd float64 `json:"currentspd"`
	CurrentDir float64 `json:"currentdir"`
}

func toReading(ri readingInput) obs.Reading {
	return obs.Reading{
		Time:       ri.Time,
		Buoy:       ri.Buoy,
		WindSpd:    ri.WindSpd,
		WindDir:    ri.WindDir,
		AirTemp:    ri.AirTemp,
		Pressure:   ri.Pressure,
		WaveHt:     ri.WaveHt,
		WavePer:    ri.WavePer,
		SST:        ri.SST,
		Salinity:   ri.Salinity,
		CurrentSpd: ri.CurrentSpd,
		CurrentDir: ri.CurrentDir,
	}
}

func toReadings(inputs []readingInput) []obs.Reading {
	out := make([]obs.Reading, len(inputs))
	for i, ri := range inputs {
		out[i] = toReading(ri)
	}
	return out
}

type analyzeRequest struct {
	Readings []readingInput `json:"readings"`
	BuoyCode string         `json:"buoy_code"`
}

type buoyAnalysis struct {
	Buoy       string  `json:"buoy"`
	Count      int     `json:"count"`
	PassRate   float64 `json:"pass_rate"`
	Hs         float64 `json:"hs"`
	MeanPeriod float64 `json:"mean_period"`
	Extreme50  float64 `json:"extreme_50y"`
}

func handleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req analyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if len(req.Readings) == 0 {
		httpError(w, http.StatusBadRequest, "readings array is empty")
		return
	}
	readings := toReadings(req.Readings)
	if req.BuoyCode != "" {
		filtered := make([]obs.Reading, 0)
		for _, rd := range readings {
			if rd.Buoy == req.BuoyCode {
				filtered = append(filtered, rd)
			}
		}
		readings = filtered
	}
	if len(readings) == 0 {
		httpError(w, http.StatusBadRequest, "no readings after filter")
		return
	}

	byBuoy := obs.ByBuoy(readings)
	results := make([]buoyAnalysis, 0, len(byBuoy))
	for _, code := range sortedKeys(byBuoy) {
		list := byBuoy[code]
		waveht := obs.ExtractField(list, func(r obs.Reading) float64 { return r.WaveHt })
		waveper := obs.ExtractField(list, func(r obs.Reading) float64 { return r.WavePer })
		scores := qc.ScoreReadings(list, 2.0, 0.05, 2.0, 3.0)
		results = append(results, buoyAnalysis{
			Buoy:       code,
			Count:      len(list),
			PassRate:   qc.PassRate(scores),
			Hs:         sea.SignificantWaveHt(waveht),
			MeanPeriod: sea.MeanWavePeriod(waveper),
			Extreme50:  sea.ExtremeReturn(waveht, 50),
		})
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"buoys": results})
}

type qcRequest struct {
	Readings []readingInput `json:"readings"`
	SpikeK   float64        `json:"spike_k"`
	FlatEps  float64        `json:"flat_eps"`
	GradMax  float64        `json:"grad_max"`
	TempRate float64        `json:"temp_rate"`
}

type qcResult struct {
	Index   int     `json:"index"`
	Flags   int     `json:"flags"`
	Quality float64 `json:"quality"`
	Range   bool    `json:"range_flag"`
	Spike   bool    `json:"spike_flag"`
	Flat    bool    `json:"flat_flag"`
	Grad    bool    `json:"grad_flag"`
	Temp    bool    `json:"temp_flag"`
}

func handleQC(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req qcRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if len(req.Readings) == 0 {
		httpError(w, http.StatusBadRequest, "readings array is empty")
		return
	}
	if req.SpikeK <= 0 {
		req.SpikeK = 2.0
	}
	if req.FlatEps <= 0 {
		req.FlatEps = 0.05
	}
	if req.GradMax <= 0 {
		req.GradMax = 2.0
	}
	if req.TempRate <= 0 {
		req.TempRate = 3.0
	}
	readings := toReadings(req.Readings)
	scores := qc.ScoreReadings(readings, req.SpikeK, req.FlatEps, req.GradMax, req.TempRate)
	results := make([]qcResult, len(scores))
	for i, s := range scores {
		results[i] = qcResult{
			Index:   s.Index,
			Flags:   s.Total,
			Quality: s.Quality,
			Range:   s.RangeFlag,
			Spike:   s.SpikeFlag,
			Flat:    s.FlatFlag,
			Grad:    s.GradFlag,
			Temp:    s.TempFlag,
		}
	}
	passRate := qc.PassRate(scores)
	avgQ := qc.AverageQuality(scores)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"pass_rate":       passRate,
		"average_quality": avgQ,
		"scores":          results,
	})
}

type seaStateRequest struct {
	Readings   []readingInput `json:"readings"`
	ReturnYrs  float64        `json:"return_years"`
}

type seaStateResult struct {
	Hs         float64    `json:"hs"`
	MeanPeriod float64    `json:"mean_period"`
	Extreme    float64    `json:"extreme_wave"`
	WindRose   [8]float64 `json:"wind_rose"`
}

func handleSeaState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req seaStateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if len(req.Readings) == 0 {
		httpError(w, http.StatusBadRequest, "readings array is empty")
		return
	}
	if req.ReturnYrs <= 1 {
		req.ReturnYrs = 50
	}
	readings := toReadings(req.Readings)
	waveht := obs.ExtractField(readings, func(r obs.Reading) float64 { return r.WaveHt })
	waveper := obs.ExtractField(readings, func(r obs.Reading) float64 { return r.WavePer })

	res := seaStateResult{
		Hs:         sea.SignificantWaveHt(waveht),
		MeanPeriod: sea.MeanWavePeriod(waveper),
		Extreme:    sea.ExtremeReturn(waveht, req.ReturnYrs),
		WindRose:   sea.WindRose(readings),
	}
	writeJSON(w, http.StatusOK, res)
}

func handleReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req struct {
		Readings []readingInput `json:"readings"`
		Format   string         `json:"format"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if len(req.Readings) == 0 {
		httpError(w, http.StatusBadRequest, "readings array is empty")
		return
	}
	readings := toReadings(req.Readings)
	byBuoy := obs.ByBuoy(readings)
	rpt := export.NewFullReport("Buoy QC Report")
	for _, code := range sortedKeys(byBuoy) {
		list := byBuoy[code]
		waveht := obs.ExtractField(list, func(r obs.Reading) float64 { return r.WaveHt })
		waveper := obs.ExtractField(list, func(r obs.Reading) float64 { return r.WavePer })
		scores := qc.ScoreReadings(list, 2.0, 0.05, 2.0, 3.0)
		hs := sea.SignificantWaveHt(waveht)
		mwp := sea.MeanWavePeriod(waveper)
		ext := sea.ExtremeReturn(waveht, 50)
		br := export.MakeBuoyReport(code, scores, hs, mwp, ext)
		rpt.AddBuoy(br)
	}
	if strings.EqualFold(req.Format, "text") {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		rpt.WriteText(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	rpt.WriteJSON(w)
}

func sortedKeys(m map[string][]obs.Reading) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func httpError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

// ParsePort extracts port number from addr string.
func ParsePort(addr string) int {
	parts := strings.Split(addr, ":")
	if len(parts) < 2 {
		return 0
	}
	p, _ := strconv.Atoi(parts[len(parts)-1])
	return p
}

// FormatAddr produces a display-friendly address.
func FormatAddr(addr string) string {
	port := ParsePort(addr)
	if port == 0 {
		return addr
	}
	return fmt.Sprintf("http://0.0.0.0:%d", port)
}

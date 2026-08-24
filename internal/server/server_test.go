package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func sampleReadings() []readingInput {
	return []readingInput{
		{Time: "2024-01-01T00:00:00Z", Buoy: "B1", WindSpd: 5, WindDir: 90, AirTemp: 15, Pressure: 1013, WaveHt: 1.2, WavePer: 6, SST: 18, Salinity: 35, CurrentSpd: 0.5, CurrentDir: 180},
		{Time: "2024-01-01T01:00:00Z", Buoy: "B1", WindSpd: 6, WindDir: 100, AirTemp: 14, Pressure: 1012, WaveHt: 1.4, WavePer: 7, SST: 17.5, Salinity: 35, CurrentSpd: 0.6, CurrentDir: 185},
		{Time: "2024-01-01T02:00:00Z", Buoy: "B1", WindSpd: 7, WindDir: 110, AirTemp: 14, Pressure: 1011, WaveHt: 1.6, WavePer: 7, SST: 17, Salinity: 34.8, CurrentSpd: 0.7, CurrentDir: 190},
		{Time: "2024-01-01T03:00:00Z", Buoy: "B1", WindSpd: 8, WindDir: 120, AirTemp: 13, Pressure: 1010, WaveHt: 1.8, WavePer: 8, SST: 16.5, Salinity: 34.5, CurrentSpd: 0.8, CurrentDir: 195},
		{Time: "2024-01-01T04:00:00Z", Buoy: "B1", WindSpd: 9, WindDir: 130, AirTemp: 13, Pressure: 1009, WaveHt: 2.0, WavePer: 8, SST: 16, Salinity: 34, CurrentSpd: 0.9, CurrentDir: 200},
		{Time: "2024-01-01T05:00:00Z", Buoy: "B2", WindSpd: 4, WindDir: 270, AirTemp: 16, Pressure: 1015, WaveHt: 0.8, WavePer: 5, SST: 19, Salinity: 36, CurrentSpd: 0.3, CurrentDir: 90},
	}
}

func TestHealthEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var body map[string]string
	json.Unmarshal(rec.Body.Bytes(), &body)
	if body["status"] != "ok" {
		t.Fatalf("expected ok, got %q", body["status"])
	}
}

func TestAnalyzeEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := analyzeRequest{Readings: sampleReadings()}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/analyze", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Buoys []buoyAnalysis `json:"buoys"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if len(resp.Buoys) != 2 {
		t.Fatalf("expected 2 buoys, got %d", len(resp.Buoys))
	}
}

func TestAnalyzeEndpoint_Empty(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	body := []byte(`{"readings":[]}`)
	req := httptest.NewRequest(http.MethodPost, "/api/analyze", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestQCEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := qcRequest{Readings: sampleReadings()}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/qc", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		PassRate float64    `json:"pass_rate"`
		Scores   []qcResult `json:"scores"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if len(resp.Scores) != 6 {
		t.Fatalf("expected 6 scores, got %d", len(resp.Scores))
	}
}

func TestSeaStateEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := seaStateRequest{Readings: sampleReadings(), ReturnYrs: 50}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/seastate", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp seaStateResult
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Hs <= 0 {
		t.Errorf("expected Hs > 0, got %.4f", resp.Hs)
	}
}

func TestReportEndpoint_JSON(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := map[string]interface{}{
		"readings": sampleReadings(),
		"format":   "json",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/report", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestReportEndpoint_Text(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := map[string]interface{}{
		"readings": sampleReadings(),
		"format":   "text",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/report", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	ct := rec.Header().Get("Content-Type")
	if ct != "text/plain; charset=utf-8" {
		t.Errorf("expected text/plain, got %q", ct)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	endpoints := []string{"/api/analyze", "/api/qc", "/api/seastate", "/api/report"}
	for _, ep := range endpoints {
		req := httptest.NewRequest(http.MethodGet, ep, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s: expected 405, got %d", ep, rec.Code)
		}
	}
}

func TestParsePort(t *testing.T) {
	if p := ParsePort(":8080"); p != 8080 {
		t.Errorf("expected 8080, got %d", p)
	}
}

func TestFormatAddr(t *testing.T) {
	got := FormatAddr(":8080")
	if got != "http://0.0.0.0:8080" {
		t.Errorf("expected http://0.0.0.0:8080, got %q", got)
	}
}

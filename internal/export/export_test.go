package export

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"buoy-qc/internal/qc"
)

func TestMakeBuoyReport(t *testing.T) {
	scores := []qc.Score{{Total: 0, Quality: 1.0}, {Total: 1, Quality: 0.8}}
	br := MakeBuoyReport("B1", scores, 2.5, 8.0, 12.0)
	if br.Buoy != "B1" {
		t.Errorf("buoy: %s", br.Buoy)
	}
	if br.RecordCount != 2 {
		t.Errorf("count: %d", br.RecordCount)
	}
	if br.Hs != 2.5 {
		t.Errorf("Hs: %f", br.Hs)
	}
}

func TestFullReport_WriteText(t *testing.T) {
	r := NewFullReport("Test Report")
	r.AddBuoy(BuoyReport{Buoy: "B1", RecordCount: 10, PassRate: 0.9, AvgQuality: 0.95})
	var buf bytes.Buffer
	r.WriteText(&buf)
	if !strings.Contains(buf.String(), "B1") {
		t.Error("text should contain buoy code")
	}
	if !strings.Contains(buf.String(), "90.0%") {
		t.Error("text should contain pass rate")
	}
}

func TestFullReport_WriteJSON(t *testing.T) {
	r := NewFullReport("JSON Test")
	r.AddBuoy(BuoyReport{Buoy: "B2", RecordCount: 5})
	var buf bytes.Buffer
	if err := r.WriteJSON(&buf); err != nil {
		t.Fatalf("write json: %v", err)
	}
	var decoded FullReport
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if decoded.Title != "JSON Test" {
		t.Errorf("title: %s", decoded.Title)
	}
}

func TestFullReport_String(t *testing.T) {
	r := NewFullReport("Str Test")
	r.AddBuoy(BuoyReport{Buoy: "B3"})
	s := r.String()
	if !strings.Contains(s, "B3") {
		t.Error("String() should contain buoy")
	}
}

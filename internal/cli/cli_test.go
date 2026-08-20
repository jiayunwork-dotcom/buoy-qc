package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_Help(t *testing.T) {
	var out bytes.Buffer
	code := Run([]string{"help"}, &out, &bytes.Buffer{})
	if code != 0 {
		t.Errorf("help should exit 0, got %d", code)
	}
	if !strings.Contains(out.String(), "usage") {
		t.Error("help output should contain usage")
	}
}

func TestRun_NoArgs(t *testing.T) {
	code := Run(nil, &bytes.Buffer{}, &bytes.Buffer{})
	if code != 2 {
		t.Errorf("no args should exit 2, got %d", code)
	}
}

func TestCmdAnalyze_NoReadings(t *testing.T) {
	var errBuf bytes.Buffer
	code := Run([]string{"analyze"}, &bytes.Buffer{}, &errBuf)
	if code != 2 {
		t.Errorf("missing -readings should exit 2, got %d", code)
	}
}

func TestCmdAnalyze_WithFile(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := Run([]string{"analyze", "-readings", "../../example/readings.csv"}, &out, &errBuf)
	if code != 0 {
		t.Errorf("analyze should exit 0, got %d; stderr: %s", code, errBuf.String())
	}
	if !strings.Contains(out.String(), "Buoy") {
		t.Error("output should contain Buoy")
	}
}

func TestCmdReport_Text(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := Run([]string{"report", "-readings", "../../example/readings.csv"}, &out, &errBuf)
	if code != 0 {
		t.Errorf("report should exit 0, got %d; stderr: %s", code, errBuf.String())
	}
	if !strings.Contains(out.String(), "Pass Rate") {
		t.Error("text report should contain Pass Rate")
	}
}

func TestCmdReport_JSON(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := Run([]string{"report", "-readings", "../../example/readings.csv", "-json"}, &out, &errBuf)
	if code != 0 {
		t.Errorf("report json should exit 0, got %d; stderr: %s", code, errBuf.String())
	}
	if !strings.Contains(out.String(), "\"title\"") {
		t.Error("JSON should contain title")
	}
}

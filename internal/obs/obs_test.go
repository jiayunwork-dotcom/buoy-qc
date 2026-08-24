package obs

import (
	"os"
	"path/filepath"
	"testing"
)

// TestParseReadings verifies error handling on missing/malformed input and
// correct slicing on valid input.
func TestParseReadings(t *testing.T) {
	// error: missing file
	if _, err := ParseReadings(filepath.Join(t.TempDir(), "nope.csv")); err == nil {
		t.Fatal("expected error for missing file")
	}

	dir := t.TempDir()

	// error: non-numeric field
	bad := filepath.Join(dir, "bad.csv")
	badContent := "time,buoy,windspd,winddir,airtemp,pressure,waveht,waveper,sst,salinity,currentspd,currentdir\n" +
		"2020,B1,abc,0,0,1000,1,5,10,35,0,0\n"
	if err := os.WriteFile(bad, []byte(badContent), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := ParseReadings(bad); err == nil {
		t.Fatal("expected error for non-numeric field")
	}

	// error: missing column
	missing := filepath.Join(dir, "missing.csv")
	missingContent := "time,buoy,windspd,winddir,airtemp,pressure,waveht,waveper,sst,salinity,currentspd\n" +
		"2020,B1,5,90,20,1000,1,5,10,35,0\n"
	if err := os.WriteFile(missing, []byte(missingContent), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := ParseReadings(missing); err == nil {
		t.Fatal("expected error for missing column")
	}

	// success: valid file
	good := filepath.Join(dir, "good.csv")
	goodContent := "time,buoy,windspd,winddir,airtemp,pressure,waveht,waveper,sst,salinity,currentspd,currentdir\n" +
		"2020,B1,5,90,20,1010,1.2,5,18,35,0.2,180\n"
	if err := os.WriteFile(good, []byte(goodContent), 0644); err != nil {
		t.Fatal(err)
	}
	rs, err := ParseReadings(good)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rs) != 1 {
		t.Fatalf("expected 1 reading, got %d", len(rs))
	}
	if rs[0].Buoy != "B1" || rs[0].Pressure != 1010 || rs[0].WaveHt != 1.2 {
		t.Fatalf("unexpected parsed reading: %+v", rs[0])
	}
}

// TestByBuoy verifies nil-safety and correct grouping by buoy code.
func TestByBuoy(t *testing.T) {
	// nil-safe
	m := ByBuoy(nil)
	if m == nil {
		t.Fatal("expected non-nil map for nil input")
	}
	if len(m) != 0 {
		t.Fatal("expected empty map")
	}
	rs := []Reading{{Buoy: "A"}, {Buoy: "B"}, {Buoy: "A"}}
	m = ByBuoy(rs)
	if len(m["A"]) != 2 || len(m["B"]) != 1 {
		t.Fatalf("unexpected grouping: %v", m)
	}
}

// Package obs provides buoy observation reading types and CSV import.
package obs

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Reading is a single ocean-buoy observation record.
type Reading struct {
	Time, Buoy                                                string
	WindSpd, WindDir, AirTemp, Pressure, WaveHt, WavePer, SST float64
	Salinity, CurrentSpd, CurrentDir                          float64
}

// required columns, in any order in the header row.
var columns = []string{
	"time", "buoy", "windspd", "winddir", "airtemp", "pressure",
	"waveht", "waveper", "sst", "salinity", "currentspd", "currentdir",
}

// ParseReadings reads buoy readings from a CSV file at path.
// It returns an error if the file is missing, the CSV is malformed, a
// required column is absent, a row has too few fields, or any numeric field
// is non-numeric.
func ParseReadings(path string) ([]Reading, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, stringifyParseErr(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	r.TrimLeadingSpace = true
	recs, err := r.ReadAll()
	if err != nil {
		return nil, stringifyParseErr(fmt.Errorf("malformed CSV: %w", err))
	}
	if len(recs) == 0 {
		return nil, stringifyParseErr(fmt.Errorf("empty file"))
	}

	header := recs[0]
	if len(header) > 0 {
		header[0] = strings.TrimPrefix(header[0], "\ufeff")
	}
	idx := make(map[string]int, len(header))
	for i, h := range header {
		idx[h] = i
	}
	for _, c := range columns {
		if _, ok := idx[c]; !ok {
			return nil, stringifyParseErr(fmt.Errorf("missing column: %s", c))
		}
	}

	out := make([]Reading, 0, len(recs)-1)
	for ri, row := range recs[1:] {
		line := ri + 2 // 1-based, after header
		if len(row) < len(columns) {
			return nil, stringifyParseErr(fmt.Errorf("line %d: expected %d fields, got %d", line, len(columns), len(row)))
		}
		rd := Reading{Time: row[idx["time"]], Buoy: row[idx["buoy"]]}
		var perr error
		if rd.WindSpd, perr = parseField(row, idx, "windspd", line); perr != nil {
			return nil, perr
		}
		if rd.WindDir, perr = parseField(row, idx, "winddir", line); perr != nil {
			return nil, perr
		}
		if rd.AirTemp, perr = parseField(row, idx, "airtemp", line); perr != nil {
			return nil, perr
		}
		if rd.Pressure, perr = parseField(row, idx, "pressure", line); perr != nil {
			return nil, perr
		}
		if rd.WaveHt, perr = parseField(row, idx, "waveht", line); perr != nil {
			return nil, perr
		}
		if rd.WavePer, perr = parseField(row, idx, "waveper", line); perr != nil {
			return nil, perr
		}
		if rd.SST, perr = parseField(row, idx, "sst", line); perr != nil {
			return nil, perr
		}
		if rd.Salinity, perr = parseField(row, idx, "salinity", line); perr != nil {
			return nil, perr
		}
		if rd.CurrentSpd, perr = parseField(row, idx, "currentspd", line); perr != nil {
			return nil, perr
		}
		if rd.CurrentDir, perr = parseField(row, idx, "currentdir", line); perr != nil {
			return nil, perr
		}
		out = append(out, rd)
	}
	return out, nil
}

// parseField parses the named column from row as a float64.
func parseField(row []string, idx map[string]int, col string, line int) (float64, error) {
	s := row[idx[col]]
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, stringifyParseErr(fmt.Errorf("line %d: column %s non-numeric %q", line, col, s))
	}
	return v, nil
}

// ByBuoy groups readings by buoy code. It is nil-safe and returns a
// non-nil (possibly empty) map when the input is nil or empty.
func ByBuoy(readings []Reading) map[string][]Reading {
	m := make(map[string][]Reading)
	if readings == nil {
		return m
	}
	for _, r := range readings {
		m[r.Buoy] = append(m[r.Buoy], r)
	}
	return m
}

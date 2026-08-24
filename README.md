# buoy-qc

`buoy-qc` is a pure-Go command-line tool and library for **ocean buoy
observational data quality control (QC) and sea-state analysis**.

It imports buoy readings from CSV (time, buoy, windspd, winddir, airtemp,
pressure, waveht, waveper, sst, salinity, currentspd, currentdir), applies
QC rules (range / spike / flat / gradient), and computes sea-state metrics:

- Significant wave height `Hs = 4 * sqrt(variance(elevation))`
- Mean wave period
- Wind rose (8 sectors of 45° each)
- Extreme return level for wave height via the Gumbel distribution

The module uses only the Go standard library. No network, no third-party
dependencies.

## Packages

- `internal/obs` — `Reading`, `ParseReadings`, `ByBuoy`
- `internal/qc`  — `RangeCheck`, `SpikeDetect`, `FlatCheck`, `GradientCheck`
- `internal/sea` — `SignificantWaveHt`, `MeanWavePeriod`, `WindRose`, `ExtremeReturn`

## CLI usage

```
go run . -readings <path> [-buoy <code>]
```

- `-readings` is required. If omitted, usage is printed to stderr and the
  process exits with code 2.
- `-buoy` is optional and filters to a single buoy.
- Malformed input exits with code 1 and a clear message.
- A successful run prints per-buoy QC flag counts, Hs, mean wave period, the
  wind rose, and the 50-year extreme return level for wave height, then exits 0.

## Examples

A sample dataset is provided in `example/readings.csv` (14 rows, two buoys,
including one out-of-range pressure value and one visible wave-height spike):

```
go run . -readings example/readings.csv
go run . -readings example/readings.csv -buoy B1
```

## Build / test

```
go vet ./...
go build ./...
go test ./...
```

package stats

import (
	"math"
	"testing"
)

func TestPercentile(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	p50 := Percentile(data, 50)
	if math.Abs(p50-5.5) > 0.01 {
		t.Errorf("P50 expected 5.5, got %f", p50)
	}
	p0 := Percentile(data, 0)
	if p0 != 1 {
		t.Errorf("P0 expected 1, got %f", p0)
	}
}

func TestQuantiles(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	q25, q50, q75 := Quantiles(data)
	if q25 >= q50 || q50 >= q75 {
		t.Errorf("quantile ordering broken: %f %f %f", q25, q50, q75)
	}
}

func TestIQR(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	iqr := IQR(data)
	if iqr <= 0 {
		t.Errorf("IQR should be > 0, got %f", iqr)
	}
}

func TestOutliers(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 100}
	outliers := Outliers(data, 1.5)
	if len(outliers) == 0 {
		t.Error("expected at least 1 outlier (100)")
	}
}

func TestHistogram(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	edges, counts := Histogram(data, 5)
	if len(edges) != 6 {
		t.Fatalf("expected 6 edges, got %d", len(edges))
	}
	total := 0
	for _, c := range counts {
		total += c
	}
	if total != 10 {
		t.Errorf("total counts should be 10, got %d", total)
	}
}

func TestMovingAverage(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	out := MovingAverage(data, 3)
	// at index 2: (2+3+4)/3 = 3
	if math.Abs(out[2]-3.0) > 0.001 {
		t.Errorf("MA[2] expected 3.0, got %f", out[2])
	}
}

func TestCorrelation(t *testing.T) {
	a := []float64{1, 2, 3, 4, 5}
	b := []float64{2, 4, 6, 8, 10}
	r := Correlation(a, b)
	if math.Abs(r-1.0) > 0.001 {
		t.Errorf("perfect correlation expected 1.0, got %f", r)
	}
}

func TestLinearRegression(t *testing.T) {
	x := []float64{0, 1, 2, 3, 4}
	y := []float64{1, 3, 5, 7, 9} // y = 2x + 1
	slope, intercept := LinearRegression(x, y)
	if math.Abs(slope-2.0) > 0.001 {
		t.Errorf("slope expected 2.0, got %f", slope)
	}
	if math.Abs(intercept-1.0) > 0.001 {
		t.Errorf("intercept expected 1.0, got %f", intercept)
	}
}

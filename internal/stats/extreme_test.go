package stats

import (
	"math"
	"testing"
)

func TestFitGumbel(t *testing.T) {
	data := []float64{3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	g := FitGumbel(data)
	if g.Scale <= 0 {
		t.Error("scale should be > 0")
	}
	if g.Location <= 0 {
		t.Error("location should be > 0 for positive data")
	}
}

func TestGumbelReturnLevel(t *testing.T) {
	data := []float64{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	g := FitGumbel(data)
	rl50 := g.ReturnLevel(50)
	rl100 := g.ReturnLevel(100)
	if rl100 <= rl50 {
		t.Errorf("RL100 should exceed RL50: rl50=%f rl100=%f", rl50, rl100)
	}
}

func TestGumbelCDF(t *testing.T) {
	g := GumbelFit{Location: 10, Scale: 2}
	// CDF at location should be exp(-1) ≈ 0.368
	cdf := g.CDF(10)
	expected := math.Exp(-1)
	if math.Abs(cdf-expected) > 0.01 {
		t.Errorf("CDF(location) expected %f, got %f", expected, cdf)
	}
}

func TestGumbelPDF(t *testing.T) {
	g := GumbelFit{Location: 10, Scale: 2}
	pdf := g.PDF(10)
	if pdf <= 0 {
		t.Error("PDF should be > 0")
	}
}

func TestFitWeibull(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	w := FitWeibull(data)
	if w.Shape <= 0 || w.Scale <= 0 {
		t.Errorf("Weibull params should be positive: k=%f lambda=%f", w.Shape, w.Scale)
	}
}

func TestWeibullCDF(t *testing.T) {
	w := WeibullFit{Shape: 2, Scale: 5}
	cdf := w.CDF(5)
	// CDF at scale: 1 - exp(-1) ≈ 0.632
	expected := 1 - math.Exp(-1)
	if math.Abs(cdf-expected) > 0.01 {
		t.Errorf("Weibull CDF(scale) expected %f, got %f", expected, cdf)
	}
}

func TestAnnualMaxima(t *testing.T) {
	data := []float64{1, 5, 3, 8, 2, 9}
	years := []int{2020, 2020, 2020, 2021, 2021, 2021}
	maxima := AnnualMaxima(data, years)
	if len(maxima) != 2 {
		t.Fatalf("expected 2 annual maxima, got %d", len(maxima))
	}
	if maxima[0] != 5 {
		t.Errorf("2020 max expected 5, got %f", maxima[0])
	}
	if maxima[1] != 9 {
		t.Errorf("2021 max expected 9, got %f", maxima[1])
	}
}

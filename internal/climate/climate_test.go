package climate

import (
	"math"
	"testing"
)

func TestMonthlyStats(t *testing.T) {
	data := []float64{10, 12, 11, 20, 22, 21}
	months := []int{1, 1, 1, 7, 7, 7}
	stats := MonthlyStats(data, months)
	if len(stats) != 2 {
		t.Fatalf("expected 2 months, got %d", len(stats))
	}
	// January
	if stats[0].Month != 1 {
		t.Errorf("first month: %d", stats[0].Month)
	}
	if math.Abs(stats[0].Mean-11) > 0.01 {
		t.Errorf("Jan mean expected 11, got %f", stats[0].Mean)
	}
}

func TestMonthToSeason(t *testing.T) {
	tests := []struct {
		month  int
		season Season
	}{
		{1, Winter}, {3, Spring}, {6, Summer}, {9, Autumn}, {12, Winter},
	}
	for _, tt := range tests {
		if MonthToSeason(tt.month) != tt.season {
			t.Errorf("month %d: got %s, want %s", tt.month, MonthToSeason(tt.month), tt.season)
		}
	}
}

func TestSeasonalMean(t *testing.T) {
	data := []float64{5, 6, 20, 21}
	months := []int{1, 2, 7, 8}
	sm := SeasonalMean(data, months)
	if math.Abs(sm[Winter]-5.5) > 0.01 {
		t.Errorf("Winter mean expected 5.5, got %f", sm[Winter])
	}
	if math.Abs(sm[Summer]-20.5) > 0.01 {
		t.Errorf("Summer mean expected 20.5, got %f", sm[Summer])
	}
}

func TestExceedanceRate(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	rate := ExceedanceRate(data, 5)
	if math.Abs(rate-0.5) > 0.01 {
		t.Errorf("exceedance rate expected 0.5, got %f", rate)
	}
}

func TestExceedanceCurve(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	th, rates := ExceedanceCurve(data, 4)
	if len(th) != 4 || len(rates) != 4 {
		t.Fatalf("expected 4 points, got %d", len(th))
	}
	// rates should decrease
	if rates[0] < rates[3] {
		t.Error("exceedance should decrease with threshold")
	}
}

func TestReturnPeriod(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	rp := ReturnPeriod(data, 9)
	// 1 out of 10 exceeds 9 → RP = 10
	if math.Abs(rp-10) > 0.01 {
		t.Errorf("return period expected 10, got %f", rp)
	}
}

func TestSeason_String(t *testing.T) {
	if Winter.String() != "Winter" {
		t.Error("Winter string failed")
	}
	if Summer.String() != "Summer" {
		t.Error("Summer string failed")
	}
}

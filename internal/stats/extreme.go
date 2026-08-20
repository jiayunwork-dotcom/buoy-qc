package stats

import (
	"math"
	"sort"
)

// GumbelFit fits a Gumbel Type I distribution (location a, scale b) to data
// using the method of moments.
type GumbelFit struct {
	Location float64 // a (mu)
	Scale    float64 // b (beta)
}

// FitGumbel estimates Gumbel parameters from data using method of moments.
func FitGumbel(data []float64) GumbelFit {
	if len(data) < 2 {
		return GumbelFit{}
	}
	const euler = 0.5772156649
	m := meanf(data)
	s := stdPopf(data)
	b := s * math.Sqrt(6) / math.Pi
	a := m - euler*b
	return GumbelFit{Location: a, Scale: b}
}

// ReturnLevel computes the return level for return period rp.
// x_rp = a - b * ln(-ln(1 - 1/rp))
func (g GumbelFit) ReturnLevel(rp float64) float64 {
	if rp <= 1 || g.Scale <= 0 {
		return 0
	}
	p := 1 - 1.0/rp
	return g.Location - g.Scale*math.Log(-math.Log(p))
}

// CDF returns the Gumbel CDF: P(X <= x) = exp(-exp(-(x-a)/b)).
func (g GumbelFit) CDF(x float64) float64 {
	if g.Scale <= 0 {
		return 0
	}
	z := (x - g.Location) / g.Scale
	return math.Exp(-math.Exp(-z))
}

// PDF returns the Gumbel PDF: f(x) = (1/b) * exp(-(z + exp(-z))).
func (g GumbelFit) PDF(x float64) float64 {
	if g.Scale <= 0 {
		return 0
	}
	z := (x - g.Location) / g.Scale
	return (1 / g.Scale) * math.Exp(-(z + math.Exp(-z)))
}

// WeibullFit fits a 2-parameter Weibull distribution using the method of moments.
type WeibullFit struct {
	Shape float64 // k
	Scale float64 // lambda
}

// FitWeibull estimates Weibull parameters. Uses a simple iterative approximation.
func FitWeibull(data []float64) WeibullFit {
	if len(data) < 2 {
		return WeibullFit{}
	}
	m := meanf(data)
	s := stdPopf(data)
	if m <= 0 {
		return WeibullFit{}
	}
	cv := s / m // coefficient of variation
	// Approximate k from CV: k ≈ 1.086 / cv (empirical for moderate CV)
	k := 1.086 / cv
	if k < 0.5 {
		k = 0.5
	}
	if k > 10 {
		k = 10
	}
	// lambda from mean: lambda = m / Gamma(1 + 1/k)
	lambda := m / math.Gamma(1+1/k)
	return WeibullFit{Shape: k, Scale: lambda}
}

// WeibullCDF returns P(X <= x) = 1 - exp(-(x/lambda)^k).
func (w WeibullFit) CDF(x float64) float64 {
	if x <= 0 || w.Scale <= 0 || w.Shape <= 0 {
		return 0
	}
	return 1 - math.Exp(-math.Pow(x/w.Scale, w.Shape))
}

// AnnualMaxima extracts annual maxima from data given years[i] for each data[i].
func AnnualMaxima(data []float64, years []int) []float64 {
	if len(data) == 0 || len(data) != len(years) {
		return nil
	}
	byYear := map[int]float64{}
	first := map[int]bool{}
	for i, y := range years {
		if !first[y] || data[i] > byYear[y] {
			byYear[y] = data[i]
			first[y] = true
		}
	}
	yearList := make([]int, 0, len(byYear))
	for y := range byYear {
		yearList = append(yearList, y)
	}
	sort.Ints(yearList)
	maxima := make([]float64, len(yearList))
	for i, y := range yearList {
		maxima[i] = byYear[y]
	}
	return maxima
}

func meanf(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

func stdPopf(data []float64) float64 {
	n := len(data)
	if n == 0 {
		return 0
	}
	m := meanf(data)
	var acc float64
	for _, v := range data {
		d := v - m
		acc += d * d
	}
	return math.Sqrt(acc / float64(n))
}

package expvar

import (
	"strconv"
	"testing"

	"github.com/go-kit/kit/metrics/teststat"
)

func TestCounter(t *testing.T) {
	counter := NewCounter("expvar_counter").With("label values", "not supported").(*Counter)
	value := func() float64 { f, _ := strconv.ParseFloat(counter.f.String(), 64); return f }
	if err := teststat.TestCounter(counter, value); err != nil {
		t.Fatal(err)
	}
}

func TestGauge(t *testing.T) {
	gauge := NewGauge("expvar_gauge").With("label values", "not supported").(*Gauge)
	value := func() float64 { f, _ := strconv.ParseFloat(gauge.f.String(), 64); return f }
	if err := teststat.TestGauge(gauge, value); err != nil {
		t.Fatal(err)
	}
}

func TestHistogram(t *testing.T) {
	histogram := NewHistogram("expvar_histogram", 50).With("label values", "not supported").(*Histogram)
	quantiles := func() (float64, float64, float64, float64) {
		p50, _ := strconv.ParseFloat(histogram.p50.String(), 64)
		p90, _ := strconv.ParseFloat(histogram.p90.String(), 64)
		p95, _ := strconv.ParseFloat(histogram.p95.String(), 64)
		p99, _ := strconv.ParseFloat(histogram.p99.String(), 64)
		return p50, p90, p95, p99
	}
	if err := teststat.TestHistogram(histogram, quantiles, 0.01); err != nil {
		t.Fatal(err)
	}
}

func TestGauge(t *testing.T) {
	var (
		name  = "xyz"
		value = 54321
		delta = 12345
		g     = expvar.NewGauge(name).With(metrics.Field{Key: "ignored", Value: "field"})
	)
	g.Set(float64(value))
	g.Add(float64(delta))
	if want, have := fmt.Sprint(value+delta), stdexpvar.Get(name).String(); want != have {
		t.Errorf("want %q, have %q", want, have)
	}
}

func TestInvalidQuantile(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("expected panic, got none")
		} else {
			t.Logf("got expected panic: %v", err)
		}
	}()
	expvar.NewHistogram("foo", 0.0, 100.0, 3, 50, 90, 95, 99, 101)
}

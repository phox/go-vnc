package metrics

import (
	"math"
	"testing"
)

func TestCounter(t *testing.T) {
	reset()

	c := NewCounter("test")
	if got, want := c.Value(), uint64(0); got != want {
		t.Errorf("initial value incorrect; got = %v, want = %v", got, want)
	}

	c.Increment()
	if got, want := c.Value(), uint64(1); got != want {
		t.Errorf("incremented value incorrect; got = %v, want = %v", got, want)
	}

	c.Reset()
	if got, want := c.Value(), uint64(0); got != want {
		t.Errorf("reset value incorrect; got = %v, want = %v", got, want)
	}

	if got, want := c.Name(), "test"; got != want {
		t.Errorf("name incorrect; got = %v, want = %v", got, want)
	}
}

func TestGauge(t *testing.T) {
	reset()

	c := NewGauge("test")
	if got, want := c.Value(), uint64(0); got != want {
		t.Errorf("initial value incorrect; got = %v, want = %v", got, want)
	}

	c.Adjust(123)
	if got, want := c.Value(), uint64(123); got != want {
		t.Errorf("incremented value incorrect; got = %v, want = %v", got, want)
	}

	c.Adjust(-23)
	if got, want := c.Value(), uint64(100); got != want {
		t.Errorf("decremented value incorrect; got = %v, want = %v", got, want)
	}

	c.Adjust(-456)
	if got, want := c.Value(), uint64(0); got != want {
		t.Errorf("minimum value incorrect; got = %v, want = %v", got, want)
	}

	c.Adjust(math.MaxInt64)
	c.Adjust(math.MaxInt64)
	c.Adjust(math.MaxInt64)
	if got, want := c.Value(), uint64(math.MaxUint64); got != want {
		t.Errorf("maximum value incorrect; got = %v, want = %v", got, want)
	}

	c.Reset()
	if got, want := c.Value(), uint64(0); got != want {
		t.Errorf("reset value incorrect; got = %v, want = %v", got, want)
	}

	if got, want := c.Name(), "test"; got != want {
		t.Errorf("name incorrect; got = %v, want = %v", got, want)
	}
}

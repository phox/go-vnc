/*
The metrics package provides support for tracking various metrics.
*/
package metrics

import (
	"log"
	"math"
)

type Metric interface {
	// Adjust increments or decrements the metric value.
	Adjust(int64)

	// Increment increases the metric value by one.
	Increment()

	// Value returns the current metric value.
	Value() uint64
}

// Counter provides a simple monotonically incrementing counter.
type Counter struct {
	val uint64
}

func (c *Counter) Adjust(val int64) {
	log.Fatal("A Counter metric cannot be adjusted")
}

func (c *Counter) Increment() {
	c.val += 1
}

func (c *Counter) Value() uint64 {
	return c.val
}

// The Gauge type represents a non-negative integer, which may increase or
// decrease, but shall never exceed the maximum value.
type Gauge struct {
	val uint64
}

// Adjust allows one to increase or decrease a metric.
func (g *Gauge) Adjust(val int64) {
	// The value is positive.
	if val > 0 {
		if g.val == math.MaxUint64 {
			return
		}
		v := g.val + uint64(val)
		if v > g.val {
			g.val = v
			return
		}
		// The value wrapped, so set to maximum allowed value.
		g.val = math.MaxUint64
		return
	}

	// The value is negative.
	v := g.val - uint64(-val)
	if v < g.val {
		g.val = v
		return
	}
	// The value wrapped, so set to zero.
	g.val = 0
}

func (g *Gauge) Increment() {
	log.Fatal("A Gauge metric cannot be adjusted")
}

func (g *Gauge) Value() uint64 {
	return g.val
}

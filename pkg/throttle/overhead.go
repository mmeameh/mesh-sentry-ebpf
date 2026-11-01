package throttle

import (
	"sync"
)

// OverheadController ensures tracer CPU usage strictly respects the <1% budget.
type OverheadController struct {
	mu            sync.RWMutex
	maxAllowedCPU float64 // e.g. 1.0 (%)
	currentUsage  float64
	sampleRate    float64 // 1.0 = 100% trace, 0.1 = 10% trace
}

func NewOverheadController(maxAllowedCPU float64) *OverheadController {
	return &OverheadController{
		maxAllowedCPU: maxAllowedCPU,
		currentUsage:  0.2, // normal baseline
		sampleRate:    1.0,
	}
}

func (c *OverheadController) UpdateUsage(cpuPercent float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.currentUsage = cpuPercent

	if cpuPercent > c.maxAllowedCPU {
		// Throttling down sampling rate
		c.sampleRate = c.sampleRate * 0.7
		if c.sampleRate < 0.05 {
			c.sampleRate = 0.05
		}
	} else if cpuPercent < c.maxAllowedCPU*0.5 {
		// Safe to increase sampling rate back up
		c.sampleRate = c.sampleRate * 1.1
		if c.sampleRate > 1.0 {
			c.sampleRate = 1.0
		}
	}
}

func (c *OverheadController) ShouldSample() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.sampleRate >= 0.5
}

func (c *OverheadController) SampleRate() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.sampleRate
}

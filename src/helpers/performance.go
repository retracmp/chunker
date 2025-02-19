package helpers

import (
	"fmt"
	"time"
)

type PerformanceTimer struct {
	start time.Time
}

func NewPerformanceTimer() *PerformanceTimer {
	return &PerformanceTimer{
		start: time.Now(),
	}
}

func (t *PerformanceTimer) EndTimer() {
	elapsed := time.Since(t.start)
	fmt.Printf("Time::%s\n", elapsed)
}
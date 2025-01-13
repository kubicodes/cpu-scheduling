package types

import "time"

type SchedulingQueue interface {
	// Core queue operations
	Enqueue(p Process) error
	Dequeue() (Process, error)
	Peek() (Process, error)

	// Queue state
	IsEmpty() bool
	Size() int

	// Metrics
	GetMetrics() SchedulingMetrics
}

type SchedulingMetrics struct {
	AverageWaitTime   time.Duration
	AverageTurnaround time.Duration
	ThroughputPerMin  float64
}

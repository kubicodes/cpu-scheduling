package scheduler

import (
	"cpu-scheduling/core/internal/types"
	"fmt"
	"sync"
	"time"
)

type FCFSQueue struct {
	processes []types.Process
	mu        sync.RWMutex

	// Metrics tracking
	totalWaitTime   time.Duration
	totalTurnaround time.Duration
	processedCount  int
	startTime       time.Time
}

func NewFCFSQueue() *FCFSQueue {
	return &FCFSQueue{
		processes: make([]types.Process, 0),
		startTime: time.Now(),
	}
}

func (q *FCFSQueue) Enqueue(p types.Process) error {
	if p == nil {
		return fmt.Errorf("cannot enqueue nil process")
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	q.processes = append(q.processes, p)
	return nil
}

func (q *FCFSQueue) Dequeue() (types.Process, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.processes) == 0 {
		return nil, fmt.Errorf("queue is empty")
	}

	// Get first process (FIFO)
	p := q.processes[0]
	q.processes = q.processes[1:]

	// Update metrics
	q.processedCount++
	q.totalWaitTime += p.GetTimeInState()
	q.totalTurnaround += p.GetTotalTime()

	return p, nil
}

func (q *FCFSQueue) Peek() (types.Process, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if len(q.processes) == 0 {
		return nil, fmt.Errorf("queue is empty")
	}

	return q.processes[0], nil
}

func (q *FCFSQueue) IsEmpty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return len(q.processes) == 0
}

func (q *FCFSQueue) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return len(q.processes)
}

func (q *FCFSQueue) GetMetrics() types.SchedulingMetrics {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.processedCount == 0 {
		return types.SchedulingMetrics{}
	}

	// Calculate averages
	avgWait := q.totalWaitTime / time.Duration(q.processedCount)
	avgTurnaround := q.totalTurnaround / time.Duration(q.processedCount)

	// Calculate throughput (processes per minute)
	elapsedMinutes := time.Since(q.startTime).Minutes()
	var throughput float64
	if elapsedMinutes > 0 {
		throughput = float64(q.processedCount) / elapsedMinutes
	}

	return types.SchedulingMetrics{
		AverageWaitTime:   avgWait,
		AverageTurnaround: avgTurnaround,
		ThroughputPerMin:  throughput,
	}
}

package queue

import (
	"cpu-scheduling/core/internal/types"
	"fmt"
	"sync"
	"time"
)

type RoundRobinQueue struct {
	processes   []types.Process
	mu          sync.RWMutex
	timeQuantum time.Duration

	// Current process tracking
	currentIndex int

	// Metrics tracking
	totalWaitTime   time.Duration
	totalTurnaround time.Duration
	processedCount  int
	startTime       time.Time
}

func NewRoundRobinQueue(timeQuantum time.Duration) *RoundRobinQueue {
	if timeQuantum <= 0 {
		timeQuantum = 100 * time.Millisecond // Default time quantum
	}

	return &RoundRobinQueue{
		processes:    make([]types.Process, 0),
		timeQuantum:  timeQuantum,
		currentIndex: 0,
		startTime:    time.Now(),
	}
}

func (q *RoundRobinQueue) Enqueue(p types.Process) error {
	if p == nil {
		return fmt.Errorf("cannot enqueue nil process")
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	q.processes = append(q.processes, p)
	return nil
}

func (q *RoundRobinQueue) Dequeue() (types.Process, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.processes) == 0 {
		return nil, fmt.Errorf("queue is empty")
	}

	// Get current process
	p := q.processes[q.currentIndex]

	// Remove it from the queue
	q.processes = append(q.processes[:q.currentIndex], q.processes[q.currentIndex+1:]...)

	// Update metrics
	q.processedCount++
	q.totalWaitTime += p.GetTimeInState()
	q.totalTurnaround += p.GetTotalTime()

	// Adjust current index if necessary
	if len(q.processes) > 0 {
		q.currentIndex = q.currentIndex % len(q.processes)
	} else {
		q.currentIndex = 0
	}

	return p, nil
}

func (q *RoundRobinQueue) Peek() (types.Process, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if len(q.processes) == 0 {
		return nil, fmt.Errorf("queue is empty")
	}

	return q.processes[q.currentIndex], nil
}

func (q *RoundRobinQueue) IsEmpty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return len(q.processes) == 0
}

func (q *RoundRobinQueue) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return len(q.processes)
}

func (q *RoundRobinQueue) GetMetrics() types.SchedulingMetrics {
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

// Round Robin specific methods

func (q *RoundRobinQueue) GetTimeQuantum() time.Duration {
	return q.timeQuantum
}

func (q *RoundRobinQueue) MoveToNext() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.processes) > 0 {
		q.currentIndex = (q.currentIndex + 1) % len(q.processes)
	}
}

func (q *RoundRobinQueue) RequeueProcess(p types.Process) error {
	return q.Enqueue(p) // Simply add to end of queue
}

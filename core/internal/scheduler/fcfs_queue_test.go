package scheduler

import (
	"cpu-scheduling/core/internal/process"
	"testing"
	"time"
)

func TestNewFCFSQueue(t *testing.T) {
	t.Run("should create empty queue", func(t *testing.T) {
		queue := NewFCFSQueue()

		if !queue.IsEmpty() {
			t.Error("new queue should be empty")
		}

		if queue.Size() != 0 {
			t.Errorf("expected size 0, got %d", queue.Size())
		}
	})
}

func TestFCFSQueue_Enqueue(t *testing.T) {
	t.Run("should enqueue process successfully", func(t *testing.T) {
		queue := NewFCFSQueue()
		p := process.NewPCB(1, process.NewTask(func() (any, error) { return nil, nil }))

		err := queue.Enqueue(p)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if queue.Size() != 1 {
			t.Errorf("expected size 1, got %d", queue.Size())
		}
	})

	t.Run("should return error for nil process", func(t *testing.T) {
		queue := NewFCFSQueue()
		err := queue.Enqueue(nil)

		if err == nil {
			t.Error("expected error for nil process")
		}
	})

	t.Run("should maintain FIFO order", func(t *testing.T) {
		queue := NewFCFSQueue()
		p1 := process.NewPCB(1, process.NewTask(func() (any, error) { return nil, nil }))
		p2 := process.NewPCB(2, process.NewTask(func() (any, error) { return nil, nil }))

		queue.Enqueue(p1)
		queue.Enqueue(p2)

		first, _ := queue.Peek()
		if first.GetPID() != 1 {
			t.Errorf("expected first process to have PID 1, got %d", first.GetPID())
		}
	})
}

func TestFCFSQueue_Dequeue(t *testing.T) {
	t.Run("should dequeue in FIFO order", func(t *testing.T) {
		queue := NewFCFSQueue()
		p1 := process.NewPCB(1, process.NewTask(func() (any, error) { return nil, nil }))
		p2 := process.NewPCB(2, process.NewTask(func() (any, error) { return nil, nil }))

		queue.Enqueue(p1)
		queue.Enqueue(p2)

		first, _ := queue.Dequeue()
		if first.GetPID() != 1 {
			t.Errorf("expected first dequeued process to have PID 1, got %d", first.GetPID())
		}

		second, _ := queue.Dequeue()
		if second.GetPID() != 2 {
			t.Errorf("expected second dequeued process to have PID 2, got %d", second.GetPID())
		}
	})

	t.Run("should return error when empty", func(t *testing.T) {
		queue := NewFCFSQueue()
		_, err := queue.Dequeue()

		if err == nil {
			t.Error("expected error when dequeuing from empty queue")
		}
	})
}

func TestFCFSQueue_Peek(t *testing.T) {
	t.Run("should return first process without removing", func(t *testing.T) {
		queue := NewFCFSQueue()
		p := process.NewPCB(1, process.NewTask(func() (any, error) { return nil, nil }))

		queue.Enqueue(p)

		first, err := queue.Peek()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if first.GetPID() != 1 {
			t.Errorf("expected PID 1, got %d", first.GetPID())
		}

		if queue.Size() != 1 {
			t.Error("peek should not remove process from queue")
		}
	})

	t.Run("should return error when empty", func(t *testing.T) {
		queue := NewFCFSQueue()
		_, err := queue.Peek()

		if err == nil {
			t.Error("expected error when peeking empty queue")
		}
	})
}

func TestFCFSQueue_GetMetrics(t *testing.T) {
	t.Run("should return zero metrics for empty queue", func(t *testing.T) {
		queue := NewFCFSQueue()
		metrics := queue.GetMetrics()

		if metrics.AverageWaitTime != 0 || metrics.AverageTurnaround != 0 || metrics.ThroughputPerMin != 0 {
			t.Error("expected zero metrics for empty queue")
		}
	})

	t.Run("should calculate metrics after processing", func(t *testing.T) {
		queue := NewFCFSQueue()
		p := process.NewPCB(1, process.NewTask(func() (any, error) { return nil, nil }))

		queue.Enqueue(p)
		time.Sleep(10 * time.Millisecond)
		_, err := queue.Dequeue()

		if err != nil {
			t.Fatalf("failed to dequeue: %v", err)
		}

		metrics := queue.GetMetrics()
		if metrics.AverageWaitTime < 10*time.Millisecond {
			t.Errorf("expected average wait time >= 10ms, got %v", metrics.AverageWaitTime)
		}

		if metrics.AverageTurnaround < 10*time.Millisecond {
			t.Errorf("expected average turnaround >= 10ms, got %v", metrics.AverageTurnaround)
		}

		if metrics.ThroughputPerMin <= 0 {
			t.Error("expected non-zero throughput")
		}
	})
}

func TestFCFSQueue_Concurrency(t *testing.T) {
	t.Run("should handle concurrent operations safely", func(t *testing.T) {
		queue := NewFCFSQueue()
		done := make(chan bool)

		// Concurrent enqueues
		go func() {
			for i := 0; i < 100; i++ {
				p := process.NewPCB(i, process.NewTask(func() (any, error) { return nil, nil }))
				queue.Enqueue(p)
			}
			done <- true
		}()

		// Concurrent dequeues
		go func() {
			for i := 0; i < 50; i++ {
				queue.Dequeue()
			}
			done <- true
		}()

		// Concurrent peeks and size checks
		go func() {
			for i := 0; i < 100; i++ {
				queue.Peek()
				queue.Size()
				queue.IsEmpty()
			}
			done <- true
		}()

		// Wait for all goroutines
		for i := 0; i < 3; i++ {
			<-done
		}

		// Final size should be 50 (100 enqueued - 50 dequeued)
		if queue.Size() != 50 {
			t.Errorf("expected final size 50, got %d", queue.Size())
		}
	})
}

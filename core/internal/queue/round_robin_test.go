package queue

import (
	"cpu-scheduling/core/internal/process"
	"testing"
	"time"
)

func TestNewRoundRobinQueue(t *testing.T) {
	t.Run("should create empty queue with specified time quantum", func(t *testing.T) {
		timeQuantum := 200 * time.Millisecond
		queue := NewRoundRobinQueue(timeQuantum)

		if !queue.IsEmpty() {
			t.Error("new queue should be empty")
		}

		if queue.GetTimeQuantum() != timeQuantum {
			t.Errorf("expected time quantum %v, got %v", timeQuantum, queue.GetTimeQuantum())
		}
	})

	t.Run("should use default time quantum when zero provided", func(t *testing.T) {
		queue := NewRoundRobinQueue(0)
		expected := 100 * time.Millisecond

		if queue.GetTimeQuantum() != expected {
			t.Errorf("expected default time quantum %v, got %v", expected, queue.GetTimeQuantum())
		}
	})
}

func TestRoundRobinQueue_Enqueue(t *testing.T) {
	t.Run("should enqueue process successfully", func(t *testing.T) {
		queue := NewRoundRobinQueue(100 * time.Millisecond)
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
		queue := NewRoundRobinQueue(100 * time.Millisecond)
		err := queue.Enqueue(nil)

		if err == nil {
			t.Error("expected error for nil process")
		}
	})
}

func TestRoundRobinQueue_Dequeue(t *testing.T) {
	t.Run("should dequeue processes in round robin order", func(t *testing.T) {
		queue := NewRoundRobinQueue(100 * time.Millisecond)
		p1 := process.NewPCB(1, process.NewTask(func() (any, error) { return nil, nil }))
		p2 := process.NewPCB(2, process.NewTask(func() (any, error) { return nil, nil }))
		p3 := process.NewPCB(3, process.NewTask(func() (any, error) { return nil, nil }))

		queue.Enqueue(p1)
		queue.Enqueue(p2)
		queue.Enqueue(p3)

		// First round
		first, _ := queue.Dequeue()
		if first.GetPID() != 1 {
			t.Errorf("expected first PID 1, got %d", first.GetPID())
		}

		// Requeue first process
		queue.RequeueProcess(first)

		// Second process
		second, _ := queue.Dequeue()
		if second.GetPID() != 2 {
			t.Errorf("expected second PID 2, got %d", second.GetPID())
		}
	})

	t.Run("should return error when empty", func(t *testing.T) {
		queue := NewRoundRobinQueue(100 * time.Millisecond)
		_, err := queue.Dequeue()

		if err == nil {
			t.Error("expected error when dequeuing from empty queue")
		}
	})
}

func TestRoundRobinQueue_MoveToNext(t *testing.T) {
	t.Run("should cycle through processes", func(t *testing.T) {
		queue := NewRoundRobinQueue(100 * time.Millisecond)
		p1 := process.NewPCB(1, process.NewTask(func() (any, error) { return nil, nil }))
		p2 := process.NewPCB(2, process.NewTask(func() (any, error) { return nil, nil }))

		queue.Enqueue(p1)
		queue.Enqueue(p2)

		// Check first process
		first, _ := queue.Peek()
		if first.GetPID() != 1 {
			t.Errorf("expected first PID 1, got %d", first.GetPID())
		}

		// Move to next
		queue.MoveToNext()

		// Check second process
		second, _ := queue.Peek()
		if second.GetPID() != 2 {
			t.Errorf("expected second PID 2, got %d", second.GetPID())
		}

		// Move to next (should wrap around)
		queue.MoveToNext()

		// Should be back to first process
		wrapped, _ := queue.Peek()
		if wrapped.GetPID() != 1 {
			t.Errorf("expected wrapped PID 1, got %d", wrapped.GetPID())
		}
	})
}

func TestRoundRobinQueue_RequeueProcess(t *testing.T) {
	t.Run("should add process to end of queue", func(t *testing.T) {
		queue := NewRoundRobinQueue(100 * time.Millisecond)
		p1 := process.NewPCB(1, process.NewTask(func() (any, error) { return nil, nil }))
		p2 := process.NewPCB(2, process.NewTask(func() (any, error) { return nil, nil }))

		queue.Enqueue(p1)
		queue.Enqueue(p2)

		// Dequeue first process
		first, _ := queue.Dequeue()
		if first.GetPID() != 1 {
			t.Errorf("expected first PID 1, got %d", first.GetPID())
		}

		// Requeue it
		queue.RequeueProcess(first)

		// Should now be at the end
		// Dequeue p2
		second, _ := queue.Dequeue()
		if second.GetPID() != 2 {
			t.Errorf("expected second PID 2, got %d", second.GetPID())
		}

		// Dequeue p1 again
		requeued, _ := queue.Dequeue()
		if requeued.GetPID() != 1 {
			t.Errorf("expected requeued PID 1, got %d", requeued.GetPID())
		}
	})
}

func TestRoundRobinQueue_GetMetrics(t *testing.T) {
	t.Run("should calculate metrics correctly with time quantum", func(t *testing.T) {
		queue := NewRoundRobinQueue(5 * time.Millisecond)
		p1 := process.NewPCB(1, process.NewTask(func() (any, error) { return nil, nil }))
		p2 := process.NewPCB(2, process.NewTask(func() (any, error) { return nil, nil }))

		queue.Enqueue(p1)
		queue.Enqueue(p2)

		time.Sleep(2 * time.Millisecond)

		// Process p1
		first, _ := queue.Dequeue()
		time.Sleep(5 * time.Millisecond) // Just one time quantum
		queue.RequeueProcess(first)

		// Process p2
		_, _ = queue.Dequeue()           // Ignore second process, just simulate processing
		time.Sleep(5 * time.Millisecond) // Just one time quantum

		metrics := queue.GetMetrics()
		if metrics.AverageWaitTime < 2*time.Millisecond {
			t.Errorf("expected average wait time >= 2ms, got %v", metrics.AverageWaitTime)
		}

		if metrics.AverageTurnaround < 4*time.Millisecond {
			t.Errorf("expected average turnaround >= 4ms, got %v", metrics.AverageTurnaround)
		}

		if metrics.ThroughputPerMin <= 0 {
			t.Error("expected non-zero throughput")
		}
	})
}

func TestRoundRobinQueue_Concurrency(t *testing.T) {
	t.Run("should handle concurrent operations safely", func(t *testing.T) {
		queue := NewRoundRobinQueue(1 * time.Millisecond)
		done := make(chan bool)
		const numProcesses = 20
		processed := make(chan int, numProcesses*2)

		// Concurrent enqueues
		go func() {
			for i := 0; i < numProcesses; i++ {
				p := process.NewPCB(i, process.NewTask(func() (any, error) { return nil, nil }))
				if err := queue.Enqueue(p); err != nil {
					t.Errorf("enqueue failed: %v", err)
				}
			}
			done <- true
		}()

		// Concurrent dequeue/requeue with timeout
		go func() {
			timeout := time.After(2 * time.Second)
			seen := make(map[int]bool)
			for len(seen) < numProcesses {
				select {
				case <-timeout:
					t.Error("test timed out")
					done <- true
					return
				default:
					if p, err := queue.Dequeue(); err == nil {
						pid := p.GetPID()
						if !seen[pid] {
							processed <- pid
							seen[pid] = true
						}
						if err := queue.RequeueProcess(p); err != nil {
							t.Errorf("requeue failed: %v", err)
						}
					}
				}
			}
			close(processed)
			done <- true
		}()

		// Wait for completion
		for i := 0; i < 2; i++ {
			<-done
		}

		// Verify each process was processed exactly once
		seen := make(map[int]bool)
		for pid := range processed {
			seen[pid] = true
		}

		for i := 0; i < numProcesses; i++ {
			if !seen[i] {
				t.Errorf("process %d was not processed", i)
			}
		}

		// Verify final queue state
		if queue.Size() != numProcesses {
			t.Errorf("expected final queue size %d, got %d", numProcesses, queue.Size())
		}
	})
}

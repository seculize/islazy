package async

import (
	"sync"
	"testing"
	"time"
)

func TestNewQueue(t *testing.T) {
	executed := false
	logic := func(job Job) {
		executed = true
	}

	q := NewQueue(1, logic)
	if q == nil {
		t.Fatal("NewQueue() returned nil")
	}

	q.Add("test job")
	q.WaitDone()
	q.Stop()

	if !executed {
		t.Error("Job was not executed")
	}
}

func TestNewQueueAutoScaling(t *testing.T) {
	logic := func(job Job) {}

	// Test with 0 workers (should auto-scale to NumCPU)
	q := NewQueue(0, logic)
	if q == nil {
		t.Fatal("NewQueue() returned nil")
	}
	q.Stop()

	// Test with negative workers (should auto-scale to NumCPU)
	q = NewQueue(-1, logic)
	if q == nil {
		t.Fatal("NewQueue() with negative workers returned nil")
	}
	q.Stop()
}

func TestWorkQueueMultipleJobs(t *testing.T) {
	var mu sync.Mutex
	var results []int

	logic := func(job Job) {
		mu.Lock()
		defer mu.Unlock()
		results = append(results, job.(int))
	}

	q := NewQueue(2, logic)

	for i := 0; i < 10; i++ {
		q.Add(i)
	}

	q.WaitDone()
	q.Stop()

	if len(results) != 10 {
		t.Errorf("Expected 10 jobs to be processed, got %d", len(results))
	}
}

func TestNewBufferedQueue(t *testing.T) {
	var count int
	var mu sync.Mutex

	logic := func(job Job) {
		mu.Lock()
		count++
		mu.Unlock()
		time.Sleep(10 * time.Millisecond)
	}

	q := NewBufferedQueue(2, logic, 5)

	// Add jobs to buffered queue
	for i := 0; i < 5; i++ {
		q.Add(i)
	}

	q.WaitDone()
	q.Stop()

	if count != 5 {
		t.Errorf("Expected 5 jobs to be processed, got %d", count)
	}
}

func TestWithTimeout(t *testing.T) {
	t.Run("completes before timeout", func(t *testing.T) {
		result, err := WithTimeout(100*time.Millisecond, func() interface{} {
			return "success"
		})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result != "success" {
			t.Errorf("Expected 'success', got %v", result)
		}
	})

	t.Run("timeout occurs", func(t *testing.T) {
		result, err := WithTimeout(10*time.Millisecond, func() interface{} {
			time.Sleep(100 * time.Millisecond)
			return "should not reach here"
		})

		if err != ErrTimeout {
			t.Errorf("Expected ErrTimeout, got %v", err)
		}

		if result != nil {
			t.Errorf("Expected nil result on timeout, got %v", result)
		}
	})

	t.Run("callback returns error", func(t *testing.T) {
		expectedErr := ErrTimeout
		result, err := WithTimeout(100*time.Millisecond, func() interface{} {
			return expectedErr
		})

		if err != expectedErr {
			t.Errorf("Expected error to be returned, got %v", err)
		}

		if result != nil {
			t.Errorf("Expected nil result when error returned, got %v", result)
		}
	})

	t.Run("callback returns nil", func(t *testing.T) {
		result, err := WithTimeout(100*time.Millisecond, func() interface{} {
			return nil
		})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})
}

func TestWorkQueueStop(t *testing.T) {
	var count int
	var mu sync.Mutex

	logic := func(job Job) {
		mu.Lock()
		count++
		mu.Unlock()
	}

	q := NewQueue(2, logic)
	q.Add(1)
	q.Add(2)

	q.Stop()

	mu.Lock()
	finalCount := count
	mu.Unlock()

	if finalCount != 2 {
		t.Errorf("Expected 2 jobs processed before stop, got %d", finalCount)
	}
}

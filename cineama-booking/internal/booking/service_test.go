package booking

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func TestConcurrentBooking_ExactlyOneWins(t *testing.T) {
	store := NewRedisStore(redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	}))
	// store := NewConcurentStore()
	svc := NewService(store)

	const numGoroutines = 10_000

	var (
		successes atomic.Int64
		failures  atomic.Int64
		wg        sync.WaitGroup
	)

	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(userNum int) {
			defer wg.Done()

			_, err := svc.Book(Booking{
				MovieID: "movie-123",
				SeatID:  "seat-1",
				UserID:  uuid.New().String(),
			})

			if err != nil {
				failures.Add(1)
			} else {
				successes.Add(1)
			}
		}(i)
	}

	wg.Wait()

	if successes.Load() != 1 {
		t.Fatalf("expected exactly 1 success, got %d", successes.Load())
	}

	t.Logf("Success: %d, Failures: %d", successes.Load(), failures.Load())
}

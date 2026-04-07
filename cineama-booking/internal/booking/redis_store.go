package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const defaultHoldTTL = 2 * time.Minute

type RedisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) *RedisStore {
	return &RedisStore{rdb: rdb}
}

func (s *RedisStore) Book(b Booking) error {
	session, err := s.hold(b)
	if err != nil {
		return err
	}

	log.Printf("session booked %+v", session)
	return nil
}
func parseSession(val string) (Booking, error) {
	var data Booking
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return Booking{}, err
	}
	return Booking{
		ID:        data.ID,
		MovieID:   data.MovieID,
		SeatID:    data.SeatID,
		UserID:    data.UserID,
		Status:    data.Status,
		// ExpiresAt: data.ExpiresAt,
	}, nil

}

func (s *RedisStore) hold(b Booking) (Booking, error) {
	ctx := context.Background()
	now := time.Now()

	key := fmt.Sprintf("seat:%s:%s", b.MovieID, b.SeatID)

	// enrich booking
	b.ID = uuid.New().String()
	b.Status = "held"
	b.ExpiresAt = now.Add(defaultHoldTTL)

	val, err := json.Marshal(b)
	if err != nil {
		return Booking{}, err
	}

	res := s.rdb.SetArgs(ctx, key, val, redis.SetArgs{
		Mode: "NX",
		TTL:  defaultHoldTTL,
	})

	if res.Err() != nil {
		return Booking{}, res.Err()
	}

	if res.Val() != "OK" {
		return Booking{}, ErrSeatAlreadyBooked
	}

	return b, nil
}

func (s *RedisStore) ListBookings(MovieID string) []Booking {
	pattern := fmt.Sprintf("seat:%s:*", MovieID)
	var sessions []Booking
	ctx := context.Background()
	iter := s.rdb.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		val, err := s.rdb.Get(ctx, iter.Val()).Result()
		if err != nil {
			log.Printf("error fetching booking for key %s: %v", iter.Val(), err)
			continue
		}

		session, err := parseSession(val)
		if err != nil {
			log.Printf("error parsing booking for key %s: %v", iter.Val(), err)
			continue
		}
		sessions = append(sessions, session)

	}
	return sessions
}

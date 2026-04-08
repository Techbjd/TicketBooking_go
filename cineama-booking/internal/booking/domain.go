package booking

import (
	"context"
	"errors"
	"time"
)

var (
	ErrSeatAlreadyBooked = errors.New("Seat is already booked")
)

type Booking struct {
	ID        string
	MovieID   string
	SeatID    string
	UserID    string
	ShowTime  string
	Status    string
	ExpiresAt time.Time
}

type BookingStore interface {
	Book(b Booking) (Booking, error)
	ListBookings(MovieID, showTime string) []Booking
	Confirm(ctx context.Context, sessionID, userID string) (Booking, error)
	Release(ctx context.Context, sessionID, userID string) error
}

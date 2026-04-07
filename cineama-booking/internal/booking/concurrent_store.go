package booking

import "sync"

type ConcurentStore struct {
	bookings map[string]Booking
	sync.RWMutex
}

func NewMemoryStore() *MemoryStory {
	return &MemoryStory{
		bookings: map[string]Booking{},
	}
}

func NewConcurentStore() *ConcurentStore {
	return &ConcurentStore{
		bookings: map[string]Booking{},
	}
}

func (s *ConcurentStore) Book(b Booking) error {
	s.Lock()
	defer s.Unlock()
	if _, exists := s.bookings[b.SeatID]; exists {
		return ErrSeatAlreadyBooked
	}
	s.bookings[b.ID] = b
	return nil
}

func (s *ConcurentStore) ListBookings(MovieID string) []Booking {
	s.RLock()
	defer s.RUnlock()
	var bookings []Booking
	for _, b := range s.bookings {
		if b.MovieID == MovieID {
			bookings = append(bookings, b)
		}
	}
	return bookings
}

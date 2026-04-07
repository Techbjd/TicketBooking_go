package booking

type MemoryStory struct {
	bookings map[string]Booking
}
/* func NewMemoryStore() *MemoryStory {
	return &MemoryStory{
		bookings: map[string]Booking{},
	}
} */


func (s *MemoryStory) Book(b Booking) error {

	if _, exists := s.bookings[b.SeatID]; exists {
		return ErrSeatAlreadyBooked
	}
	s.bookings[b.ID] = b
	return nil
}

func (s *MemoryStory) ListBookings(MovieID string) []Booking {
	var bookings []Booking
	for _, b := range s.bookings {
		if b.MovieID == MovieID {
			bookings = append(bookings, b)
		}
	}
	return bookings
}

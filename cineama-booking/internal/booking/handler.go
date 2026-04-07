package booking

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Techbjd/cinemabooking/internal/booking/utils"
)

type Handler struct {
	svc Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: *svc}
}

func (h *Handler) ListSeats(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	bookings := h.svc.ListBookings(movieID)

	seats := make([]seatInfo, 0, len(bookings))
	for _, b := range bookings {
		seats = append(seats, seatInfo{
			SeatID: b.SeatID,
			UserID: b.UserID,
			Booked: true,
		})
	}

	utils.WriteJSON(w, http.StatusOK, seats)

}

type holdRequest struct {
	UserID string `json:"user_id"`
}
type holdResponse struct {
	SeatsID   string `json:"seats_id"`
	MovieID   string `json:"movie_id"`
	SessionId string `json:"session_id"`
	ExpiresAt string `json:"expires_at"`
}

func (h *Handler) Hold(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	seatsID := r.PathValue("seatsID")

	var req holdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err)
		return
	}
	data := Booking{
		UserID:  req.UserID,
		SeatID:  seatsID,
		MovieID: movieID,
	}
	session, err := h.svc.Book(data)
	if err != nil {
		log.Println(err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, holdResponse{
		SeatsID:   seatsID,
		MovieID:   movieID,
		SessionId: session.ID,
		ExpiresAt: session.ExpiresAt.Format(time.RFC3339),
	})

}

type seatInfo struct {
	SeatID string `json:"seat_id"`
	UserID string `json:"user_id"`
	Booked bool   `json:"booked"`
}

func (h *Handler) ConfirmSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")
	var req holdSeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err)
		return
	}
	if req.UserID == "" {
		return
	}
	session, err := h.svc.ConfirmSeat(r.Context(), sessionID, req.UserID)
	if err != nil {
		log.Println(err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, sessionResponse{
		SessionID: sessionID,
		MovieID:   session.MovieID,
		SeatID:    session.SeatID,
		UserID:    session.UserID,
		Status:    session.Status,
	})
}

func (h *Handler) ReleaseSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")
	var req holdSeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err)
		return
	}
	if req.UserID == "" {
		return
	}
	err := h.svc.ReleaseSeat(r.Context(), sessionID, req.UserID)
	if err != nil {
		log.Println(err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "session released successfully",
	})
}

type holdSeatRequest struct {
	UserID string `json:"user_id"`
}
type sessionResponse struct {
	SessionID string `json:"session_id"`
	MovieID   string `json:"movie_id"`
	SeatID    string `json:"seat_id"`
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
}

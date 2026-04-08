package booking

import (
	"encoding/json"
	"fmt"
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
	if movieID == "" {
		utils.WriteError(w, http.StatusBadRequest, "invalid_request", "movie ID is required")
		return
	}

	showTime := r.URL.Query().Get("show_time")
	rows := r.URL.Query().Get("rows")
	seatsPerRow := r.URL.Query().Get("seats_per_row")

	var totalRows, totalSeatsPerRow int
	if rows != "" && seatsPerRow != "" {
		fmt.Sscanf(rows, "%d", &totalRows)
		fmt.Sscanf(seatsPerRow, "%d", &totalSeatsPerRow)
	}

	if totalRows == 0 {
		totalRows = 5
	}
	if totalSeatsPerRow == 0 {
		totalSeatsPerRow = 8
	}

	bookings := h.svc.ListBookings(movieID, showTime)
	bookedSeats := make(map[string]bool)
	for _, b := range bookings {
		bookedSeats[b.SeatID] = true
	}

	seats := make([]seatInfo, 0, totalRows*totalSeatsPerRow)
	for row := 1; row <= totalRows; row++ {
		for seat := 1; seat <= totalSeatsPerRow; seat++ {
			seatID := fmt.Sprintf("R%dS%d", row, seat)
			seats = append(seats, seatInfo{
				SeatID: seatID,
				Booked: bookedSeats[seatID],
			})
		}
	}

	utils.WriteJSON(w, http.StatusOK, seats)

}

type holdRequest struct {
	UserID string `json:"user_id"`
}
type holdResponse struct {
	SeatsID   string `json:"seats_id"`
	MovieID   string `json:"movie_id"`
	ShowTime  string `json:"show_time"`
	SessionId string `json:"session_id"`
	ExpiresAt string `json:"expires_at"`
}

func (h *Handler) Hold(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	seatsID := r.PathValue("seatsID")
	showTime := r.URL.Query().Get("show_time")

	if movieID == "" || seatsID == "" {
		utils.WriteError(w, http.StatusBadRequest, "invalid_request", "movie ID and seat ID are required")
		return
	}

	var req holdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}
	if req.UserID == "" {
		utils.WriteError(w, http.StatusBadRequest, "invalid_request", "user_id is required")
		return
	}
	data := Booking{
		UserID:   req.UserID,
		SeatID:   seatsID,
		MovieID:  movieID,
		ShowTime: showTime,
	}
	session, err := h.svc.Book(data)
	if err != nil {
		if err == ErrSeatAlreadyBooked {
			utils.WriteError(w, http.StatusConflict, "seat_already_booked", "this seat is already booked")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to book seat")
		return
	}
	utils.WriteJSON(w, http.StatusOK, holdResponse{
		SeatsID:   seatsID,
		MovieID:   movieID,
		ShowTime:  showTime,
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
	if sessionID == "" {
		utils.WriteError(w, http.StatusBadRequest, "invalid_request", "session ID is required")
		return
	}
	var req holdSeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}
	if req.UserID == "" {
		utils.WriteError(w, http.StatusBadRequest, "invalid_request", "user_id is required")
		return
	}
	session, err := h.svc.ConfirmSeat(r.Context(), sessionID, req.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "session_not_found", "session not found or already confirmed")
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
	if sessionID == "" {
		utils.WriteError(w, http.StatusBadRequest, "invalid_request", "session ID is required")
		return
	}
	var req holdSeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}
	if req.UserID == "" {
		utils.WriteError(w, http.StatusBadRequest, "invalid_request", "user_id is required")
		return
	}
	err := h.svc.ReleaseSeat(r.Context(), sessionID, req.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "session_not_found", "session not found or already released")
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

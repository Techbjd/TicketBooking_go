package booking

import "net/http"



type Handler struct {
	svc Service
}



func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ListSeats(w http.ResponseWriter,r *http.Request){
	h.svc.ListBookings()
}
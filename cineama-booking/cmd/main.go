package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Techbjd/cinemabooking/internal/adapters/redis"
	"github.com/Techbjd/cinemabooking/internal/booking"
)

type moviesResponse struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Rows       int    `json:"rows"`
	SeatPerRow int    `json:"seat_per_row"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /movies", listMovies)
	mux.Handle("GET /", http.FileServer(http.Dir("static")))
	store := booking.NewRedisStore(redis.NewClient("localhost:6379"))
	svc := booking.NewService(store)
	bookingHandler := booking.NewHandler(svc)
	mux.HandleFunc("GET /movies/{movieID}/seats", bookingHandler.ListSeats)
	mux.HandleFunc("POST /movies/{movieID}/seats/{seatsID}", bookingHandler.Hold)

	mux.HandleFunc("PUT /session/{sessionID}/confirm", bookingHandler.ConfirmSession)
	mux.HandleFunc("DELETE /session/{sessionID}", bookingHandler.ReleaseSession)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
	fmt.Println("The server is running on the port 8080")
}

var movies = []moviesResponse{
	{ID: "1", Title: "Movie 1", Rows: 5, SeatPerRow: 8},
	{ID: "2", Title: "Movie 2", Rows: 4, SeatPerRow: 6},
}

func listMovies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, movies)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to write json response: %v", err)
	}
}

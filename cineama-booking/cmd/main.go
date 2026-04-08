package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Techbjd/cinemabooking/internal/adapters/redis"
	"github.com/Techbjd/cinemabooking/internal/booking"
)

type Movie struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Duration    int      `json:"duration"` // in minutes
	Genre       string   `json:"genre"`
	Rating      float64  `json:"rating"`
	Rows        int      `json:"rows"`
	SeatPerRow  int      `json:"seat_per_row"`
	ShowTimes   []string `json:"show_times"`
	ImageURL    string   `json:"image_url"`
}

var movies = []Movie{
	{
		ID:          "1",
		Title:       "Avatar: The Way of Water",
		Description: "Jake Sully lives with his newfound family formed on the planet of Pandora. When a familiar threat returns to finish what was previously started, Jake must work with Neytiri and the army of the Na'vi race to protect their planet.",
		Duration:    192,
		Genre:       "Sci-Fi",
		Rating:      7.6,
		Rows:        5,
		SeatPerRow:  8,
		ShowTimes:   []string{"10:00", "13:30", "17:00", "20:30"},
		ImageURL:    "/static/images/avatar.jpg",
	},
	{
		ID:          "2",
		Title:       "The Dark Knight",
		Description: "When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological and physical tests of his ability to fight injustice.",
		Duration:    152,
		Genre:       "Action",
		Rating:      9.0,
		Rows:        4,
		SeatPerRow:  6,
		ShowTimes:   []string{"11:00", "14:00", "17:00", "20:00"},
		ImageURL:    "/static/images/dark-knight.jpg",
	},
	{
		ID:          "3",
		Title:       "Inception",
		Description: "A thief who steals corporate secrets through the use of dream-sharing technology is given the inverse task of planting an idea into the mind of a C.E.O.",
		Duration:    148,
		Genre:       "Thriller",
		Rating:      8.8,
		Rows:        6,
		SeatPerRow:  10,
		ShowTimes:   []string{"10:30", "13:00", "16:00", "19:00", "22:00"},
		ImageURL:    "/static/images/inception.jpg",
	},
	{
		ID:          "4",
		Title:       "Interstellar",
		Description: "A team of explorers travel through a wormhole in space in an attempt to ensure humanity's survival.",
		Duration:    169,
		Genre:       "Sci-Fi",
		Rating:      8.7,
		Rows:        5,
		SeatPerRow:  10,
		ShowTimes:   []string{"12:00", "15:30", "19:00", "22:30"},
		ImageURL:    "/static/images/interstellar.jpg",
	},
	{
		ID:          "5",
		Title:       "The Matrix",
		Description: "When a beautiful stranger leads computer hacker Neo to a forbidding underworld, he discovers the shocking truth--the life he knows is the elaborate deception of an evil cyber-intelligence.",
		Duration:    136,
		Genre:       "Sci-Fi",
		Rating:      8.7,
		Rows:        4,
		SeatPerRow:  8,
		ShowTimes:   []string{"10:00", "13:00", "16:00", "19:00", "22:00"},
		ImageURL:    "/static/images/matrix.jpg",
	},
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /movies", listMovies)
	mux.HandleFunc("GET /movies/{movieID}", getMovieDetails)
	mux.Handle("GET /", http.FileServer(http.Dir("static")))
	store := booking.NewRedisStore(redis.NewClient("localhost:6379"))
	svc := booking.NewService(store)
	bookingHandler := booking.NewHandler(svc)
	mux.HandleFunc("GET /movies/{movieID}/seats", bookingHandler.ListSeats)
	mux.HandleFunc("POST /movies/{movieID}/seats/{seatsID}", bookingHandler.Hold)

	mux.HandleFunc("PUT /session/{sessionID}/confirm", bookingHandler.ConfirmSession)
	mux.HandleFunc("DELETE /session/{sessionID}", bookingHandler.ReleaseSession)

	handler := corsMiddleware(mux)

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
	fmt.Println("The server is running on the port 8080")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func listMovies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, movies)
}

func getMovieDetails(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	if movieID == "" {
		http.Error(w, `{"error": "invalid_request", "message": "movie ID is required"}`, http.StatusBadRequest)
		return
	}

	for _, m := range movies {
		if m.ID == movieID {
			writeJSON(w, http.StatusOK, m)
			return
		}
	}
	http.Error(w, `{"error": "not_found", "message": "movie not found"}`, http.StatusNotFound)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to write json response: %v", err)
	}
}

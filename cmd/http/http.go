package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"what-to-watch/handlers"
)

// Server holds the HTTP server instance
type Server struct {
	port int
}

// NewServer creates a new HTTP server
func NewServer(port int) *Server {
	return &Server{port: port}
}

// Start begins listening for HTTP requests
func (s *Server) Start() error {
	http.HandleFunc("/api/shows", handleGetShows)
	http.HandleFunc("/api/shows/mark", handleMarkShowWatched)
	http.HandleFunc("/api/films", handleGetFilms)
	http.HandleFunc("/health", handleHealth)

	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("HTTP server listening on port %d\n", s.port)
	return http.ListenAndServe(addr, nil)
}

func handleGetShows(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodError(w, http.MethodGet)
		return
	}

	shows, err := handlers.GetCurrentlyWatchingShows()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, shows)
}

func handleMarkShowWatched(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodError(w, http.MethodPost)
		return
	}

	idx := r.URL.Query().Get("index")
	if idx == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("index query parameter is required"))
		return
	}

	showIdx, err := strconv.Atoi(idx)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("index must be a valid integer"))
		return
	}

	isCompleted, err := handlers.MarkShowWatched(showIdx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, isCompleted)
}

func handleGetFilms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodError(w, http.MethodGet)
		return
	}

	films, err := handlers.GetAllFilms()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, films)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodError(w, http.MethodGet)
		return
	}

	writeJSON(w, http.StatusOK, nil)
}

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	fmt.Printf("%s\n", err)
}

func methodError(w http.ResponseWriter, allowedMethod string) {
	writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("only %s is allowed", allowedMethod))
}

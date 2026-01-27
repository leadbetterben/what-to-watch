package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"what-to-watch/data"
	"what-to-watch/handlers"
)

// Handler defines the interface for business logic functions
type Handler interface {
	GetCurrentlyWatchingShows() ([]data.Show, error)
	MarkShowWatched(idx int) (bool, error)
	GetAllFilms() ([]data.Film, error)
}

// defaultHandler uses the handlers package functions
type defaultHandler struct{}

func (h *defaultHandler) GetCurrentlyWatchingShows() ([]data.Show, error) {
	return handlers.GetCurrentlyWatchingShows()
}

func (h *defaultHandler) MarkShowWatched(idx int) (bool, error) {
	return handlers.MarkShowWatched(idx)
}

func (h *defaultHandler) GetAllFilms() ([]data.Film, error) {
	return handlers.GetAllFilms()
}

// Server holds the HTTP server instance
type Server struct {
	port    int
	handler Handler
}

// NewServer creates a new HTTP server
func NewServer(port int) *Server {
	return &Server{
		port:    port,
		handler: &defaultHandler{},
	}
}

// NewServerWithHandler creates a new HTTP server with a custom handler (for testing)
func NewServerWithHandler(port int, handler Handler) *Server {
	return &Server{
		port:    port,
		handler: handler,
	}
}

// Start begins listening for HTTP requests
func (s *Server) Start() error {
	http.HandleFunc("/shows", func(w http.ResponseWriter, r *http.Request) {
		s.handleGetShows(w, r)
	})
	http.HandleFunc("/shows/watch", func(w http.ResponseWriter, r *http.Request) {
		s.handleMarkShowWatched(w, r)
	})
	http.HandleFunc("/films", func(w http.ResponseWriter, r *http.Request) {
		s.handleGetFilms(w, r)
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		s.handleHealth(w, r)
	})

	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("HTTP server listening on port %d\n", s.port)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleGetShows(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodError(w, http.MethodGet)
		return
	}

	shows, err := s.handler.GetCurrentlyWatchingShows()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, shows)
}

func (s *Server) handleMarkShowWatched(w http.ResponseWriter, r *http.Request) {
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

	isCompleted, err := s.handler.MarkShowWatched(showIdx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, isCompleted)
}

func (s *Server) handleGetFilms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodError(w, http.MethodGet)
		return
	}

	films, err := s.handler.GetAllFilms()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, films)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
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

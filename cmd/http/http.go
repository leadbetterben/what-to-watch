package http

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"what-to-watch/data"
	"what-to-watch/handlers"
)

// Handler defines the interface for business logic functions
type Handler interface {
	GetCurrentlyWatchingShows() ([]data.Show, error)
	MarkShowWatched(idx int) (bool, error)
	GetAllFilms() ([]data.Film, error)
	GetAvailableGenres() ([]string, error)
	GetUnwatchedShowsByGenre(genre string) ([]data.Show, error)
	GetHalloweenEpisodes() ([]data.ShowSpecialEpisode, error)
	GetChristmasEpisodes() ([]data.ShowSpecialEpisode, error)
	GetRandomHalloweenEpisode() (data.ShowSpecialEpisode, error)
	GetRandomChristmasEpisode() (data.ShowSpecialEpisode, error)
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

func (h *defaultHandler) GetAvailableGenres() ([]string, error) {
	return handlers.GetAvailableGenres()
}

func (h *defaultHandler) GetUnwatchedShowsByGenre(genre string) ([]data.Show, error) {
	return handlers.GetUnwatchedShowsByGenre(genre)
}

func (h *defaultHandler) GetHalloweenEpisodes() ([]data.ShowSpecialEpisode, error) {
	return handlers.GetSpecialEpisodes(data.SpecialCategoryHalloween)
}

func (h *defaultHandler) GetChristmasEpisodes() ([]data.ShowSpecialEpisode, error) {
	return handlers.GetSpecialEpisodes(data.SpecialCategoryChristmas)
}

func (h *defaultHandler) GetRandomHalloweenEpisode() (data.ShowSpecialEpisode, error) {
	return handlers.GetRandomSpecialEpisode(data.SpecialCategoryHalloween)
}

func (h *defaultHandler) GetRandomChristmasEpisode() (data.ShowSpecialEpisode, error) {
	return handlers.GetRandomSpecialEpisode(data.SpecialCategoryChristmas)
}

// Server holds the HTTP server instance
type Server struct {
	handler        Handler
	httpServer     *http.Server
	listenAndServe func(*http.Server) error
}

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 30 * time.Second
	idleTimeout       = 60 * time.Second
)

// NewServer creates a new HTTP server
func NewServer(port int) *Server {
	return newServer(port, &defaultHandler{})
}

// NewServerWithHandler creates a new HTTP server with a custom handler (for testing)
func NewServerWithHandler(port int, handler Handler) *Server {
	return newServer(port, handler)
}

func newServer(port int, handler Handler) *Server {
	s := &Server{handler: handler}
	s.httpServer = &http.Server{
		Addr:              net.JoinHostPort("127.0.0.1", strconv.Itoa(port)),
		Handler:           s.routes(),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}
	s.listenAndServe = func(server *http.Server) error {
		return server.ListenAndServe()
	}
	return s
}

// Start begins listening for HTTP requests
func (s *Server) Start() error {
	fmt.Printf("HTTP server listening on %s\n", s.httpServer.Addr)
	return s.listenAndServe(s.httpServer)
}

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/shows", func(w http.ResponseWriter, r *http.Request) {
		s.handleGetShows(w, r)
	})
	mux.HandleFunc("/shows/watch", func(w http.ResponseWriter, r *http.Request) {
		s.handleMarkShowWatched(w, r)
	})
	mux.HandleFunc("/films", func(w http.ResponseWriter, r *http.Request) {
		s.handleGetFilms(w, r)
	})
	mux.HandleFunc("/genres", func(w http.ResponseWriter, r *http.Request) {
		s.handleGetGenres(w, r)
	})
	mux.HandleFunc("/specials/halloween", func(w http.ResponseWriter, r *http.Request) {
		s.handleGetHalloweenEpisodes(w, r)
	})
	mux.HandleFunc("/specials/christmas", func(w http.ResponseWriter, r *http.Request) {
		s.handleGetChristmasEpisodes(w, r)
	})
	mux.HandleFunc("/specials/halloween/random", func(w http.ResponseWriter, r *http.Request) {
		s.handleGetRandomHalloweenEpisode(w, r)
	})
	mux.HandleFunc("/specials/christmas/random", func(w http.ResponseWriter, r *http.Request) {
		s.handleGetRandomChristmasEpisode(w, r)
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		s.handleHealth(w, r)
	})

	return mux
}

func (s *Server) handleGetShows(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodError(w, http.MethodGet)
		return
	}

	genre := r.URL.Query().Get("genre")
	if genre != "" {
		shows, err := s.handler.GetUnwatchedShowsByGenre(genre)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		writeJSON(w, http.StatusOK, shows)
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
		writeMethodError(w, http.MethodPost)
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
		writeMethodError(w, http.MethodGet)
		return
	}

	films, err := s.handler.GetAllFilms()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, films)
}

func (s *Server) handleGetGenres(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodError(w, http.MethodGet)
		return
	}

	genres, err := s.handler.GetAvailableGenres()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, genres)
}

func (s *Server) handleGetHalloweenEpisodes(w http.ResponseWriter, r *http.Request) {
	s.handleGetSpecialEpisodes(w, r, s.handler.GetHalloweenEpisodes)
}

func (s *Server) handleGetChristmasEpisodes(w http.ResponseWriter, r *http.Request) {
	s.handleGetSpecialEpisodes(w, r, s.handler.GetChristmasEpisodes)
}

func (s *Server) handleGetRandomHalloweenEpisode(w http.ResponseWriter, r *http.Request) {
	s.handleGetRandomSpecialEpisode(w, r, s.handler.GetRandomHalloweenEpisode)
}

func (s *Server) handleGetRandomChristmasEpisode(w http.ResponseWriter, r *http.Request) {
	s.handleGetRandomSpecialEpisode(w, r, s.handler.GetRandomChristmasEpisode)
}

func (s *Server) handleGetSpecialEpisodes(w http.ResponseWriter, r *http.Request, get func() ([]data.ShowSpecialEpisode, error)) {
	if r.Method != http.MethodGet {
		writeMethodError(w, http.MethodGet)
		return
	}

	episodes, err := get()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, episodes)
}

func (s *Server) handleGetRandomSpecialEpisode(w http.ResponseWriter, r *http.Request, get func() (data.ShowSpecialEpisode, error)) {
	if r.Method != http.MethodGet {
		writeMethodError(w, http.MethodGet)
		return
	}

	episode, err := get()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, episode)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodError(w, http.MethodGet)
		return
	}

	writeJSON(w, http.StatusOK, nil)
}

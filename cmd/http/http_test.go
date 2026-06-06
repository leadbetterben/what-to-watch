package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"what-to-watch/data"
)

// mockHandler implements the Handler interface for testing
type mockHandler struct {
	getShowsFunc            func() ([]data.Show, error)
	markShowWatchedFunc     func(idx int) (bool, error)
	getFilmsFunc            func() ([]data.Film, error)
	getGenresFunc           func() ([]string, error)
	getUnwatchedByGenreFunc func(genre string) ([]data.Show, error)
}

func (m *mockHandler) GetCurrentlyWatchingShows() ([]data.Show, error) {
	return m.getShowsFunc()
}

func (m *mockHandler) MarkShowWatched(idx int) (bool, error) {
	return m.markShowWatchedFunc(idx)
}

func (m *mockHandler) GetAllFilms() ([]data.Film, error) {
	return m.getFilmsFunc()
}

func (m *mockHandler) GetAvailableGenres() ([]string, error) {
	return m.getGenresFunc()
}

func (m *mockHandler) GetUnwatchedShowsByGenre(genre string) ([]data.Show, error) {
	return m.getUnwatchedByGenreFunc(genre)
}

func TestHandleGetShows(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		genreParam     string
		mockShows      []data.Show
		mockErr        error
		expectedStatus int
		expectShowLen  int
	}{
		{
			name:   "successful get shows",
			method: http.MethodGet,
			mockShows: []data.Show{
				{
					Name:     "Breaking Bad",
					Genre:    "Drama",
					Provider: "Netflix",
					Episodes: []int{1, 2, 3},
				},
				{
					Name:     "The Crown",
					Genre:    "Drama",
					Provider: "Netflix",
					Episodes: []int{1, 2},
				},
			},
			mockErr:        nil,
			expectedStatus: http.StatusOK,
			expectShowLen:  2,
		},
		{
			name:           "empty shows list",
			method:         http.MethodGet,
			mockShows:      []data.Show{},
			mockErr:        nil,
			expectedStatus: http.StatusOK,
			expectShowLen:  0,
		},
		{
			name:           "handler error",
			method:         http.MethodGet,
			mockShows:      nil,
			mockErr:        fmt.Errorf("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectShowLen:  0,
		},
		{
			name:           "invalid method POST",
			method:         http.MethodPost,
			mockShows:      []data.Show{},
			mockErr:        nil,
			expectedStatus: http.StatusMethodNotAllowed,
			expectShowLen:  0,
		},
		{
			name:           "invalid method DELETE",
			method:         http.MethodDelete,
			mockShows:      []data.Show{},
			mockErr:        nil,
			expectedStatus: http.StatusMethodNotAllowed,
			expectShowLen:  0,
		},
		{
			name:       "successful get shows by genre",
			method:     http.MethodGet,
			genreParam: "Drama",
			mockShows: []data.Show{
				{
					Name:     "Breaking Bad",
					Genre:    "Drama",
					Provider: "Netflix",
					Episodes: []int{1, 2, 3},
				},
				{
					Name:     "The Crown",
					Genre:    "Drama",
					Provider: "Netflix",
					Episodes: []int{1, 2},
				},
			},
			mockErr:        nil,
			expectedStatus: http.StatusOK,
			expectShowLen:  2,
		},
		{
			name:           "empty shows for genre",
			method:         http.MethodGet,
			genreParam:     "Horror",
			mockShows:      []data.Show{},
			mockErr:        nil,
			expectedStatus: http.StatusOK,
			expectShowLen:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockHandler{
				getShowsFunc: func() ([]data.Show, error) {
					return tt.mockShows, tt.mockErr
				},
				getUnwatchedByGenreFunc: func(genre string) ([]data.Show, error) {
					return tt.mockShows, tt.mockErr
				},
			}

			server := NewServerWithHandler(8080, mock)
			url := "/shows"
			if tt.genreParam != "" {
				url = url + fmt.Sprintf("?genre=%s", tt.genreParam)
			}
			req := httptest.NewRequest(tt.method, url, nil)
			w := httptest.NewRecorder()

			server.handleGetShows(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			body, _ := io.ReadAll(w.Body)
			var shows []data.Show
			if err := json.Unmarshal(body, &shows); err != nil {
				if tt.expectShowLen > 0 {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
			}
			if len(shows) != tt.expectShowLen {
				t.Errorf("expected %d shows, got %d", tt.expectShowLen, len(shows))
			}
		})
	}
}

func TestHandleMarkShowWatched(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		indexParam     string
		mockCompleted  bool
		mockErr        error
		expectedStatus int
		expectBody     bool
	}{
		{
			name:           "successful mark show watched",
			method:         http.MethodPost,
			indexParam:     "0",
			mockCompleted:  false,
			mockErr:        nil,
			expectedStatus: http.StatusOK,
			expectBody:     false,
		},
		{
			name:           "mark show as series complete",
			method:         http.MethodPost,
			indexParam:     "1",
			mockCompleted:  true,
			mockErr:        nil,
			expectedStatus: http.StatusOK,
			expectBody:     true,
		},
		{
			name:           "missing index parameter",
			method:         http.MethodPost,
			indexParam:     "",
			mockCompleted:  false,
			mockErr:        nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid index parameter (non-integer)",
			method:         http.MethodPost,
			indexParam:     "invalid",
			mockCompleted:  false,
			mockErr:        nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid index parameter (negative)",
			method:         http.MethodPost,
			indexParam:     "-1",
			mockCompleted:  false,
			mockErr:        nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid method GET",
			method:         http.MethodGet,
			indexParam:     "0",
			mockCompleted:  false,
			mockErr:        nil,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "handler error",
			method:         http.MethodPost,
			indexParam:     "0",
			mockCompleted:  false,
			mockErr:        fmt.Errorf("failed to update show"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockHandler{
				markShowWatchedFunc: func(idx int) (bool, error) {
					return tt.mockCompleted, tt.mockErr
				},
			}

			server := NewServerWithHandler(8080, mock)
			url := "/shows/watch"
			if tt.indexParam != "" {
				url += "?index=" + tt.indexParam
			}

			req := httptest.NewRequest(tt.method, url, nil)
			w := httptest.NewRecorder()

			server.handleMarkShowWatched(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			body, _ := io.ReadAll(w.Body)
			if w.Code == http.StatusOK {
				var result bool
				if err := json.Unmarshal(bytes.TrimSpace(body), &result); err != nil {
					t.Errorf("failed to unmarshal response: %v", err)
				}
				if result != tt.expectBody {
					t.Errorf("expected body %v, got %v", tt.expectBody, result)
				}
			}
		})
	}
}

func TestHandleGetFilms(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		mockFilms      []data.Film
		mockErr        error
		expectedStatus int
		expectFilmLen  int
	}{
		{
			name:   "successful get films",
			method: http.MethodGet,
			mockFilms: []data.Film{
				{
					Name:     "Inception",
					Genre:    "Sci-Fi",
					Provider: "Netflix",
				},
				{
					Name:     "The Matrix",
					Genre:    "Sci-Fi",
					Provider: "Prime Video",
				},
			},
			mockErr:        nil,
			expectedStatus: http.StatusOK,
			expectFilmLen:  2,
		},
		{
			name:           "empty films list",
			method:         http.MethodGet,
			mockFilms:      []data.Film{},
			mockErr:        nil,
			expectedStatus: http.StatusOK,
			expectFilmLen:  0,
		},
		{
			name:           "handler error",
			method:         http.MethodGet,
			mockFilms:      nil,
			mockErr:        fmt.Errorf("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectFilmLen:  0,
		},
		{
			name:           "invalid method POST",
			method:         http.MethodPost,
			mockFilms:      []data.Film{},
			mockErr:        nil,
			expectedStatus: http.StatusMethodNotAllowed,
			expectFilmLen:  0,
		},
		{
			name:           "invalid method PUT",
			method:         http.MethodPut,
			mockFilms:      []data.Film{},
			mockErr:        nil,
			expectedStatus: http.StatusMethodNotAllowed,
			expectFilmLen:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockHandler{
				getFilmsFunc: func() ([]data.Film, error) {
					return tt.mockFilms, tt.mockErr
				},
			}

			server := NewServerWithHandler(8080, mock)
			req := httptest.NewRequest(tt.method, "/films", nil)
			w := httptest.NewRecorder()

			server.handleGetFilms(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			body, _ := io.ReadAll(w.Body)
			var films []data.Film
			if err := json.Unmarshal(body, &films); err != nil {
				if tt.expectFilmLen > 0 {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
			}
			if len(films) != tt.expectFilmLen {
				t.Errorf("expected %d films, got %d", tt.expectFilmLen, len(films))
			}
		})
	}
}

func TestHandleGetGenres(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		mockGenres     []string
		mockErr        error
		expectedStatus int
		expectGenreLen int
	}{
		{
			name:           "successful get genres",
			method:         http.MethodGet,
			mockGenres:     []string{"Drama", "Comedy", "Sci-Fi"},
			mockErr:        nil,
			expectedStatus: http.StatusOK,
			expectGenreLen: 3,
		},
		{
			name:           "empty genres list",
			method:         http.MethodGet,
			mockGenres:     []string{},
			mockErr:        nil,
			expectedStatus: http.StatusOK,
			expectGenreLen: 0,
		},
		{
			name:           "handler error",
			method:         http.MethodGet,
			mockGenres:     nil,
			mockErr:        fmt.Errorf("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectGenreLen: 0,
		},
		{
			name:           "invalid method POST",
			method:         http.MethodPost,
			mockGenres:     []string{},
			mockErr:        nil,
			expectedStatus: http.StatusMethodNotAllowed,
			expectGenreLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockHandler{
				getGenresFunc: func() ([]string, error) {
					return tt.mockGenres, tt.mockErr
				},
			}

			server := NewServerWithHandler(8080, mock)
			req := httptest.NewRequest(tt.method, "/genres", nil)
			w := httptest.NewRecorder()

			server.handleGetGenres(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			body, _ := io.ReadAll(w.Body)
			var genres []string
			if err := json.Unmarshal(body, &genres); err != nil {
				if tt.expectGenreLen > 0 {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
			}
			if len(genres) != tt.expectGenreLen {
				t.Errorf("expected %d genres, got %d", tt.expectGenreLen, len(genres))
			}
		})
	}
}

func TestHandleHealth(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "successful health check",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid method POST",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "invalid method PUT",
			method:         http.MethodPut,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "invalid method DELETE",
			method:         http.MethodDelete,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockHandler{}
			server := NewServerWithHandler(8080, mock)

			req := httptest.NewRequest(tt.method, "/health", nil)
			w := httptest.NewRecorder()

			server.handleHealth(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestServerRoutes(t *testing.T) {
	var watchedIndex int
	mock := &mockHandler{
		getShowsFunc: func() ([]data.Show, error) {
			return []data.Show{{Name: "Slow Horses", Genre: "Drama", Provider: "Apple TV+"}}, nil
		},
		getUnwatchedByGenreFunc: func(genre string) ([]data.Show, error) {
			return []data.Show{{Name: "Unstarted " + genre, Genre: genre, Provider: "Netflix"}}, nil
		},
		markShowWatchedFunc: func(idx int) (bool, error) {
			watchedIndex = idx
			return true, nil
		},
		getFilmsFunc: func() ([]data.Film, error) {
			return []data.Film{{Name: "Arrival", Genre: "Sci-Fi", Provider: "Netflix"}}, nil
		},
		getGenresFunc: func() ([]string, error) {
			return []string{"Drama", "Comedy"}, nil
		},
	}

	server := NewServerWithHandler(8080, mock)
	testServer := httptest.NewServer(server.routes())
	defer testServer.Close()

	t.Run("GET /shows", func(t *testing.T) {
		resp, body := getRoute(t, testServer.URL+"/shows")
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var shows []data.Show
		unmarshalBody(t, body, &shows)
		if len(shows) != 1 || shows[0].Name != "Slow Horses" {
			t.Fatalf("unexpected shows response: %+v", shows)
		}
	})

	t.Run("GET /shows with genre", func(t *testing.T) {
		resp, body := getRoute(t, testServer.URL+"/shows?genre=Drama")
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var shows []data.Show
		unmarshalBody(t, body, &shows)
		if len(shows) != 1 || shows[0].Name != "Unstarted Drama" {
			t.Fatalf("unexpected genre response: %+v", shows)
		}
	})

	t.Run("POST /shows/watch", func(t *testing.T) {
		resp, body := postRoute(t, testServer.URL+"/shows/watch?index=3")
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}
		if watchedIndex != 3 {
			t.Fatalf("expected watched index 3, got %d", watchedIndex)
		}

		var completed bool
		unmarshalBody(t, body, &completed)
		if !completed {
			t.Fatalf("expected completed response to be true")
		}
	})

	t.Run("GET /films", func(t *testing.T) {
		resp, body := getRoute(t, testServer.URL+"/films")
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var films []data.Film
		unmarshalBody(t, body, &films)
		if len(films) != 1 || films[0].Name != "Arrival" {
			t.Fatalf("unexpected films response: %+v", films)
		}
	})

	t.Run("GET /genres", func(t *testing.T) {
		resp, body := getRoute(t, testServer.URL+"/genres")
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var genres []string
		unmarshalBody(t, body, &genres)
		if len(genres) != 2 || genres[0] != "Drama" || genres[1] != "Comedy" {
			t.Fatalf("unexpected genres response: %+v", genres)
		}
	})

	t.Run("GET /health", func(t *testing.T) {
		resp, body := getRoute(t, testServer.URL+"/health")
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}
		if got := bytes.TrimSpace(body); !bytes.Equal(got, []byte("null")) {
			t.Fatalf("expected health body null, got %q", string(got))
		}
	})
}

func TestServerRoutesMethodAndNotFoundResponses(t *testing.T) {
	mock := &mockHandler{}
	server := NewServerWithHandler(8080, mock)
	testServer := httptest.NewServer(server.routes())
	defer testServer.Close()

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{name: "unsupported route method", method: http.MethodPost, path: "/shows", wantStatus: http.StatusMethodNotAllowed},
		{name: "unknown route", method: http.MethodGet, path: "/missing", wantStatus: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, testServer.URL+tt.path, nil)
			if err != nil {
				t.Fatalf("creating request: %v", err)
			}

			resp, err := testServer.Client().Do(req)
			if err != nil {
				t.Fatalf("sending request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}
		})
	}
}

func getRoute(t *testing.T, url string) (*http.Response, []byte) {
	t.Helper()

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		t.Fatalf("reading body: %v", err)
	}

	return resp, body
}

func postRoute(t *testing.T, url string) (*http.Response, []byte) {
	t.Helper()

	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		t.Fatalf("POST %s: %v", url, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		t.Fatalf("reading body: %v", err)
	}

	return resp, body
}

func unmarshalBody(t *testing.T, body []byte, v interface{}) {
	t.Helper()

	if err := json.Unmarshal(body, v); err != nil {
		t.Fatalf("unmarshaling %q: %v", string(body), err)
	}
}

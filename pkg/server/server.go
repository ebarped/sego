package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ebarped/sego/pkg/engine"
	"github.com/go-chi/chi/v5"
)

type Result struct {
	Query     string   `json:"query"`
	Documents []string `json:"documents"`
}

type Server struct {
	router *chi.Mux
	port   int
	engine *engine.Engine
}

func New(p string, e *engine.Engine) Server {
	port, err := strconv.Atoi(p)
	if err != nil {
		log.Fatalf("error parsing port of server: %s\n", err)
	}

	s := Server{
		port:   port,
		router: chi.NewRouter(),
		engine: e,
	}

	s.router.Get("/search", s.handleSearch())

	return s
}

// Start starts a server
func (s *Server) Start() {
	log.Println("starting server on port " + fmt.Sprint(s.port))
	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(s.port), s.router))
}

// will handle the /search?query="example of a query" route
func (s Server) handleSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			response    string
			resultCount int
			err         error
		)

		query := r.URL.Query().Get("query")
		if query == "" || query == "\"\"" {
			response = "You have to provide a query. Example: curl 'localhost:4000/search?query=memory%20management'"
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(response))
			return
		}

		resultCountStr := r.URL.Query().Get("count")
		if resultCountStr == "" {
			resultCount = 5
		} else {
			resultCount, err = strconv.Atoi(resultCountStr)
			if err != nil {
				fmt.Printf("error parsing integer from %s: %s\n", resultCountStr, err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("error parsing the result count:%s", err)))
				return
			}
		}

		log.Printf("searching for %q, get %d results.\n", query, resultCount)

		res := s.engine.Search(query, resultCount)
		for _, doc := range res {
			response += doc + "\n"
		}
		result := Result{
			Query:     query,
			Documents: res,
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			log.Printf("error writing results to client: %s\n", err)
		}

		return
	}
}

package server

import (
	"fmt"
	"jukebox/internal/db"
	"jukebox/internal/routes"
	"net/http"

	"github.com/tmc/langchaingo/llms"
)

type server struct {
	port     int
	database *db.DB
	llm      llms.Model
}

func New(port int, database *db.DB, llm llms.Model) *server {
	return &server{
		port:     port,
		database: database,
		llm:      llm,
	}
}

func (s *server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/conversations", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			routes.GetConversations(s.database, w, r)
			return
		} else if r.Method == http.MethodPost {
			routes.CreateConversation(s.database, w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/api/v1/conversations/{id}/messages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			routes.GetChatMessages(s.database, w, r)
			return
		} else if r.Method == http.MethodPost {
			routes.CreateChatMessage(s.database, s.llm, w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/api/v1/conversations/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			routes.GetConversation(s.database, w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), mux)
	if err != nil {
		return err
	}

	return nil
}

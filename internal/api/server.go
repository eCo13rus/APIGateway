package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	addr   string
}

func NewServer(handler *Handler, addr string) *Server {
	router := mux.NewRouter()

	router.Use(RequestIDMiddleware)
	router.Use(LoggingMiddleware)

	router.HandleFunc("/api/news", handler.GetNews).Methods(http.MethodGet)
	router.HandleFunc("/api/news/{id:[0-9]+}", handler.GetNewsDetails).Methods(http.MethodGet)
	router.HandleFunc("/api/news/{id:[0-9]+}/comments", handler.GetNewsComments).Methods(http.MethodGet)
	router.HandleFunc("/api/news/{id:[0-9]+}/comments", handler.AddComment).Methods(http.MethodPost)
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("webapp"))))

	return &Server{
		router: router,
		addr:   addr,
	}
}

func (s *Server) Start() error {
	log.Printf("API Gateway запущен на %s", s.addr)
	return http.ListenAndServe(s.addr, s.router)
}

package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type AbstractServer interface {
	ListenAndServe(*AbstractServer)
}

type Server struct {
	router *mux.Router
}

func CreateServer() *Server {
	server := Server{}

	server.router = mux.NewRouter()
	server.router.HandleFunc("/upload", upload).Methods("POST")

	return &server
}

func (s *Server) ListenAndServe(address string) error {
	return http.ListenAndServe(address, s.router)
}

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

	// Добавление скриптов
	jsHandler := http.StripPrefix("/js/", http.FileServer(http.Dir("./ui/js/")))
	server.router.PathPrefix("/js").Methods("GET").Handler(jsHandler)

	// Добавление html файлов
	htmlHandler := http.StripPrefix("/", http.FileServer(http.Dir("./ui/html/")))
	server.router.PathPrefix("/").Methods("GET").Handler(htmlHandler)

	return &server
}

func (s *Server) ListenAndServe(address string) error {
	return http.ListenAndServe(address, s.router)
}

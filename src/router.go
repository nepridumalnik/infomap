package server

import (
	"net/http"

	"github.com/glebarez/sqlite"
	"github.com/gorilla/mux"

	"gorm.io/gorm"
)

type AbstractServer interface {
	ListenAndServe(*AbstractServer)
}

type Server struct {
	router *mux.Router
	db     *gorm.DB
}

func CreateServer() (*Server, error) {
	db, err := gorm.Open(sqlite.Open("data.db"))

	if err != nil {
		return nil, err
	}

	server := Server{
		db:     db,
		router: mux.NewRouter(),
	}

	server.db.AutoMigrate(&TableRow{})

	server.router = mux.NewRouter()
	server.router.HandleFunc("/upload", upload).Methods("POST")

	// Добавление скриптов
	jsHandler := http.StripPrefix("/js/", http.FileServer(http.Dir("./ui/js/")))
	server.router.PathPrefix("/js").Methods("GET").Handler(jsHandler)

	// Добавление html файлов
	htmlHandler := http.StripPrefix("/", http.FileServer(http.Dir("./ui/html/")))
	server.router.PathPrefix("/").Methods("GET").Handler(htmlHandler)

	return &server, nil
}

func (s *Server) ListenAndServe(address string) error {
	return http.ListenAndServe(address, s.router)
}

package app

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type App struct {
	router  *mux.Router
	storage *storage
	address string
}

func CreateApp(address string) (*App, error) {
	storage, err := NewStorage()

	if err != nil {
		return nil, err
	}

	app := App{
		storage: storage,
		router:  mux.NewRouter(),
		address: address,
	}

	app.storage.RegisterHandlers(app.router.Path("/upload").Methods("POST"))

	// Добавление скриптов
	jsHandler := http.StripPrefix("/js/", http.FileServer(http.Dir("./ui/js/")))
	app.router.PathPrefix("/js").Methods("GET").Handler(jsHandler)

	// Добавление html файлов
	htmlHandler := http.StripPrefix("/", http.FileServer(http.Dir("./ui/html/")))
	app.router.PathPrefix("/").Methods("GET").Handler(htmlHandler)

	return &app, nil
}

func (s *App) Run() error {
	file, err := os.OpenFile("error_log.txt", os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return err
	}

	server := &http.Server{
		Handler:      s.router,
		Addr:         s.address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		ErrorLog:     log.New(file, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	return server.ListenAndServe()
}

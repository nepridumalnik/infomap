package app

import (
	"net/http"

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

	app.router = mux.NewRouter()

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
	return http.ListenAndServe(s.address, s.router)
}

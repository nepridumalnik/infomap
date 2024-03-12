package app

import (
	"net/http"

	"github.com/glebarez/sqlite"
	"github.com/gorilla/mux"

	"gorm.io/gorm"
)

type App struct {
	router  *mux.Router
	db      *gorm.DB
	address string
}

func CreateApp(address string) (*App, error) {
	db, err := gorm.Open(sqlite.Open("data.db"))

	if err != nil {
		return nil, err
	}

	app := App{
		db:      db,
		router:  mux.NewRouter(),
		address: address,
	}

	app.db.AutoMigrate(&TableRow{})

	app.router = mux.NewRouter()
	app.router.HandleFunc("/upload", upload).Methods("POST")

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

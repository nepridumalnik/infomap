package app

import (
	"html/template"
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

	app.registerHandlers()

	return &app, nil
}

// Регистрация всех обработчиков
func (app *App) registerHandlers() {
	// Проверка авторизации
	mw := &middleware{storage: app.storage}
	app.router.Use(mw.authMiddleware)

	// Обработчик авторизации
	app.router.HandleFunc("/auth", mw.authHandler).Methods("GET", "POST")
	app.router.HandleFunc("/unauth", mw.unauthHandler).Methods("POST")

	// Загрузка файлов
	app.storage.RegisterHandlers(app.router.Path("/upload").Methods("POST"))

	// Загрузка скриптов
	jsHandler := http.StripPrefix("/js/", http.FileServer(http.Dir("./ui/js/")))
	app.router.PathPrefix("/js").Methods("GET").Handler(jsHandler)

	// Загрузка стилей
	cssHandler := http.StripPrefix("/styles/", http.FileServer(http.Dir("./ui/styles/")))
	app.router.PathPrefix("/styles").Methods("GET").Handler(cssHandler)

	// Загрузка html
	app.router.HandleFunc("/", app.commonHandler).Methods("GET", "POST")
}

func (app *App) commonHandler(w http.ResponseWriter, r *http.Request) {
	const defaultPath = "./ui/html/index.html"

	var tmpl *template.Template
	var err error

	if r.URL.Path == "/" {
		tmpl, err = template.ParseFiles(defaultPath)
	} else {
		tmpl, err = template.ParseFiles(r.URL.Path)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := "Test Template Data"
	err = tmpl.Execute(w, data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *App) Run() error {
	file, err := os.OpenFile("error_log.txt", os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return err
	}

	server := &http.Server{
		Handler:      app.router,
		Addr:         app.address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		ErrorLog:     log.New(file, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	return server.ListenAndServe()
}

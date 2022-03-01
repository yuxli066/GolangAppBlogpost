package app

import (
	"log"
	"net/http"

	"leo-blog-post/src/github.com/yuxli066/blogposts/app/handler"

	"github.com/gorilla/mux"
)

// App struct holds key app components, such as the api router
type App struct {
	Router *mux.Router
}

// Initialize server application
func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.setRouters()
}

// The setRouters function specifies different backend routes for the api
func (a *App) setRouters() {
	a.Get("/api/ping", a.handleRequest(handler.GetHealthCheck))
	a.Get("/api/posts", a.handleRequest(handler.GetPosts))
}

// HTTP CRUD wrapper function for HTTP GET
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// HTTP CRUD wrapper function for HTTP POST
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// The Run function runs the api on specified port number
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}

// The RequestHandlerFunction handles incoming http requests
type RequestHandlerFunction func(w http.ResponseWriter, r *http.Request)

// The handleRequest function handles http requests using the defined RequestHandlerFunction above
func (a *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}
}

// The GetHTTPHandler function returns the http handler for App
// USED FOR TESTING PURPOSES ONLY
func (a *App) GetHTTPHandler() http.Handler {
	return a.Router
}

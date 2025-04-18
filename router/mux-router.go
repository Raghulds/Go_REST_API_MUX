package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Router interface
type Router interface {
	PING(uri string, f func(responseWriter http.ResponseWriter, request *http.Request))
	GET(uri string, f func(responseWriter http.ResponseWriter, request *http.Request))
	POST(uri string, f func(responseWriter http.ResponseWriter, request *http.Request))
	SERVE(port string)
}

var gorillaMuxRouter = mux.NewRouter()

type muxRouter struct{}

func NewMuxRouter() Router {
	return &muxRouter{}
}

func (mr *muxRouter) PING(uri string, f func(responseWriter http.ResponseWriter, request *http.Request)) {
	gorillaMuxRouter.HandleFunc(uri, f).Methods("GET")
}

func (mr *muxRouter) GET(uri string, f func(responseWriter http.ResponseWriter, request *http.Request)) {
	gorillaMuxRouter.HandleFunc(uri, f).Methods("GET")
}

func (mr *muxRouter) POST(uri string, f func(responseWriter http.ResponseWriter, request *http.Request)) {
	gorillaMuxRouter.HandleFunc(uri, f).Methods("POST")
}

func (mr *muxRouter) SERVE(port string) {
	log.Println("Starting server at Port: ", port)
	err := http.ListenAndServe(port, gorillaMuxRouter)
	if err != nil {
		panic(err)
	}
}

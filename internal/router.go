package internal

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Router mux.Router

func NewRouter() *Router {
	return (*Router)(mux.NewRouter())
}

func (router *Router) inner() *mux.Router {
	return (*mux.Router)(router)
}

func routerHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		handler(writer, request)
	}
}

func (router *Router) Get(path string, handler http.HandlerFunc) *Router {
	router.inner().Methods("GET").Path(path).HandlerFunc(routerHandler(handler))
	return router
}

func (router *Router) Post(path string, handler http.HandlerFunc) *Router {
	router.inner().Methods("POST").Path(path).HandlerFunc(routerHandler(handler))
	return router
}

func (router *Router) Put(path string, handler http.HandlerFunc) *Router {
	router.inner().Methods("PUT").Path(path).HandlerFunc(routerHandler(handler))
	return router
}

func (router *Router) Delete(path string, handler http.HandlerFunc) *Router {
	router.inner().Methods("DELETE").Path(path).HandlerFunc(routerHandler(handler))
	return router
}

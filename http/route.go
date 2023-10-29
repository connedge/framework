package http

import (
	"context"
	"net/http"
)

type GroupFunc func(routes Route)

type Engine interface {
	Route
	Fallback(handler HandlerFunc)
	GlobalMiddleware(middlewares ...Middleware)
	Run(port string) error
	RunTLS(host ...string) error
	RunTLSWithCert(host, certFile, keyFile string) error
	ServeHTTP(writer http.ResponseWriter, request *http.Request)
	Shutdown(ctx context.Context) error
}

type Route interface {
	Group(handler GroupFunc)
	Prefix(addr string) Route
	Middleware(middlewares ...Middleware) Route

	Any(relativePath string, handler HandlerFunc)
	Get(relativePath string, handler HandlerFunc)
	Post(relativePath string, handler HandlerFunc)
	Delete(relativePath string, handler HandlerFunc)
	Patch(relativePath string, handler HandlerFunc)
	Put(relativePath string, handler HandlerFunc)
	Options(relativePath string, handler HandlerFunc)

	Static(relativePath, root string)
	StaticFile(relativePath, filepath string)
	StaticFS(relativePath string, fs http.FileSystem)
}

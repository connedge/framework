package http

import (
	"bytes"
	"context"
	"net/http"
)

type HttpEngineType string

const (
	GinServer HttpEngineType = "GinServer"
)

type Middleware func(Context)

type HandlerFunc func(Context) error

type Context interface {
	context.Context
	Context() context.Context
	WithValue(key string, value any)
	Request() Request
	Response() Response
}

type Request interface {
	Header(key string, defaultValue ...string) string
	Headers() http.Header
	Method() string
	Path() string
	Url() string
	FullUrl() string
	Ip() string
	Host() string

	// All Retrieve json, form and query
	All() map[string]any
	// Bind Retrieve json and bind to obj
	Bind(obj any) error
	// Route Retrieve an input item from the request: /users/{id}
	Route(key string) string
	RouteInt(key string) int
	RouteInt64(key string) int64
	// Query Retrieve a query string item form the request: /users?id=1
	Query(key string, defaultValue ...string) string
	QueryInt(key string, defaultValue ...int) int
	QueryInt64(key string, defaultValue ...int64) int64
	QueryBool(key string, defaultValue ...bool) bool
	QueryArray(key string) []string
	QueryMap(key string) map[string]string
	Queries() map[string]string

	// Input Retrieve data by order: json, form, query, route
	Input(key string, defaultValue ...string) string
	InputInt(key string, defaultValue ...int) int
	InputInt64(key string, defaultValue ...int64) int64
	InputBool(key string, defaultValue ...bool) bool

	//File(name string) (filesystem.File, error)

	AbortWithStatus(code int)
	AbortWithStatusJson(code int, jsonObj any)

	Next()
	Origin() *http.Request
}

type FormRequest interface {
	Authorize(ctx Context) error
	Rules(ctx Context) map[string]string
	Messages(ctx Context) map[string]string
	Attributes(ctx Context) map[string]string
}

type Response interface {
	Data(code int, contentType string, data []byte)
	Download(filepath, filename string)
	File(filepath string)
	Header(key, value string) Response
	Json(code int, obj any)
	Origin() ResponseOrigin
	Redirect(code int, location string)
	String(code int, format string, values ...any)
	Success() ResponseSuccess
	Writer() http.ResponseWriter
}

type ResponseSuccess interface {
	Data(contentType string, data []byte)
	Json(obj any)
	String(format string, values ...any)
}

type ResponseOrigin interface {
	Body() *bytes.Buffer
	Header() http.Header
	Size() int
	Status() int
}

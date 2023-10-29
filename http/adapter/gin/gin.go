package gin

import (
	"context"
	"errors"
	"fmt"
	"github.com/connedge/framework/config"
	"github.com/connedge/framework/http"
	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
	nethttp "net/http"
)

type Gin struct {
	http.Route
	instance *gin.Engine
	server   *nethttp.Server
	config   config.Config
}

func NewGin(config config.Config) *Gin {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	return &Gin{
		Route: NewGinGroup(engine.Group("/"),
			"",
			[]http.Middleware{},
			[]http.Middleware{GinResponseMiddleware()},
		),
		config:   config,
		instance: engine,
		server: &nethttp.Server{
			Addr:    ":8080",
			Handler: engine,
		},
	}
}

func (r *Gin) Fallback(handler http.HandlerFunc) {
	r.instance.NoRoute(handlerToGinHandler(handler))
}

func (r *Gin) GlobalMiddleware(middlewares ...http.Middleware) {
	if len(middlewares) > 0 {
		r.instance.Use(middlewaresToGinHandlers(middlewares)...)
	}
	r.Route = NewGinGroup(r.instance.Group("/"),
		"",
		[]http.Middleware{},
		[]http.Middleware{GinResponseMiddleware()},
	)
}

func (r *Gin) Run(port string) error {
	color.Yellowln("[HTTP] Listening and serving HTTP on " + port)

	if port != "" {
		r.server.Addr = port
	}

	return r.server.ListenAndServe()

}

func (r *Gin) RunTLS(host ...string) error {
	if len(host) == 0 {
		defaultHost := r.config.GetString("http.tls.host")
		if defaultHost == "" {
			return errors.New("host can't be empty")
		}

		defaultPort := r.config.GetString("http.tls.port")
		if defaultPort == "" {
			return errors.New("port can't be empty")
		}
		completeHost := defaultHost + ":" + defaultPort
		host = append(host, completeHost)
	}

	certFile := r.config.GetString("http.tls.ssl.cert")
	keyFile := r.config.GetString("http.tls.ssl.key")

	return r.RunTLSWithCert(host[0], certFile, keyFile)
}

func (r *Gin) RunTLSWithCert(host, certFile, keyFile string) error {
	if host == "" {
		return errors.New("host can't be empty")
	}
	if certFile == "" || keyFile == "" {
		return errors.New("certificate can't be empty")
	}

	r.outputRoutes()
	color.Greenln("[HTTPS] Listening and serving HTTPS on " + host)

	return r.instance.RunTLS(host, certFile, keyFile)
}

func (r *Gin) outputRoutes() {
	if r.config.GetBool("app.debug") {
		for _, item := range r.instance.Routes() {
			fmt.Printf("%-10s %s\n", item.Method, colonToBracket(item.Path))
		}
	}
}

func (r *Gin) ServeHTTP(writer nethttp.ResponseWriter, request *nethttp.Request) {
	r.instance.ServeHTTP(writer, request)
}

func (r *Gin) Shutdown(ctx context.Context) error {
	return r.server.Shutdown(ctx)
}

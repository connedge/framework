package gin

import (
	"github.com/connedge/framework/http"
	"github.com/gin-gonic/gin"
	"strings"
)

func handlerToGinHandler(handler http.HandlerFunc) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		handler(NewGinContext(ginCtx))
	}
}

func middlewaresToGinHandlers(middlewares []http.Middleware) []gin.HandlerFunc {
	var ginHandlers []gin.HandlerFunc
	for _, item := range middlewares {
		ginHandlers = append(ginHandlers, middlewareToGinHandler(item))
	}

	return ginHandlers
}

func middlewareToGinHandler(handler http.Middleware) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		handler(NewGinContext(ginCtx))
	}
}

func colonToBracket(relativePath string) string {
	arr := strings.Split(relativePath, "/")
	var newArr []string
	for _, item := range arr {
		if strings.HasPrefix(item, ":") {
			item = "{" + strings.ReplaceAll(item, ":", "") + "}"
		}
		newArr = append(newArr, item)
	}

	return strings.Join(newArr, "/")
}

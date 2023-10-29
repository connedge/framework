package gin

import (
	"github.com/connedge/framework/http"
	"github.com/gin-gonic/gin"
	nethttp "net/http"
)

type GinGroup struct {
	instance          gin.IRouter
	originPrefix      string
	prefix            string
	originMiddlewares []http.Middleware
	middlewares       []http.Middleware
	lastMiddlewares   []http.Middleware
}

func NewGinGroup(
	instance gin.IRouter,
	prefix string,
	originMiddlewares []http.Middleware,
	lastMiddlewares []http.Middleware) http.Route {
	return &GinGroup{
		instance:          instance,
		originPrefix:      prefix,
		originMiddlewares: originMiddlewares,
		lastMiddlewares:   lastMiddlewares,
	}
}

func (r *GinGroup) Group(handler http.GroupFunc) {
	var middlewares []http.Middleware
	middlewares = append(middlewares, r.originMiddlewares...)
	middlewares = append(middlewares, r.middlewares...)
	r.middlewares = []http.Middleware{}
	prefix := r.originPrefix + "/" + r.prefix
	r.prefix = ""

	handler(NewGinGroup(r.instance, prefix, middlewares, r.lastMiddlewares))
}

func (r *GinGroup) Prefix(addr string) http.Route {
	r.prefix += "/" + addr
	return r
}

func (r *GinGroup) Middleware(middlewares ...http.Middleware) http.Route {
	r.middlewares = append(r.middlewares, middlewares...)

	return r
}

func (r *GinGroup) Any(relativePath string, handler http.HandlerFunc) {
	r.getGinRoutesWithMiddlewares().Any(relativePath, []gin.HandlerFunc{handlerToGinHandler(handler)}...)
	r.clearMiddlewares()
}

func (r *GinGroup) Get(relativePath string, handler http.HandlerFunc) {
	r.getGinRoutesWithMiddlewares().GET(relativePath, []gin.HandlerFunc{handlerToGinHandler(handler)}...)
	r.clearMiddlewares()
}

func (r *GinGroup) Post(relativePath string, handler http.HandlerFunc) {
	r.getGinRoutesWithMiddlewares().POST(relativePath, []gin.HandlerFunc{handlerToGinHandler(handler)}...)
	r.clearMiddlewares()
}

func (r *GinGroup) Delete(relativePath string, handler http.HandlerFunc) {
	r.getGinRoutesWithMiddlewares().DELETE(relativePath, []gin.HandlerFunc{handlerToGinHandler(handler)}...)
	r.clearMiddlewares()
}

func (r *GinGroup) Patch(relativePath string, handler http.HandlerFunc) {
	r.getGinRoutesWithMiddlewares().PATCH(relativePath, []gin.HandlerFunc{handlerToGinHandler(handler)}...)
	r.clearMiddlewares()
}

func (r *GinGroup) Put(relativePath string, handler http.HandlerFunc) {
	r.getGinRoutesWithMiddlewares().PUT(relativePath, []gin.HandlerFunc{handlerToGinHandler(handler)}...)
	r.clearMiddlewares()
}

func (r *GinGroup) Options(relativePath string, handler http.HandlerFunc) {
	r.getGinRoutesWithMiddlewares().OPTIONS(relativePath, []gin.HandlerFunc{handlerToGinHandler(handler)}...)
	r.clearMiddlewares()
}

func (r *GinGroup) Static(relativePath, root string) {
	r.getGinRoutesWithMiddlewares().Static(relativePath, root)
	r.clearMiddlewares()
}

func (r *GinGroup) StaticFile(relativePath, filepath string) {
	r.getGinRoutesWithMiddlewares().StaticFile(relativePath, filepath)
	r.clearMiddlewares()
}

func (r *GinGroup) StaticFS(relativePath string, fs nethttp.FileSystem) {
	r.getGinRoutesWithMiddlewares().StaticFS(relativePath, fs)
	r.clearMiddlewares()
}

func (r *GinGroup) clearMiddlewares() {
	r.middlewares = []http.Middleware{}
}

func (r *GinGroup) getGinRoutesWithMiddlewares() gin.IRoutes {
	var middlewares []gin.HandlerFunc

	prefix := r.originPrefix + "/" + r.prefix
	r.prefix = ""
	ginGroup := r.instance.Group(prefix)

	ginOriginMiddlewares := middlewaresToGinHandlers(r.originMiddlewares)
	ginMiddlewares := middlewaresToGinHandlers(r.middlewares)
	ginLastMiddlewares := middlewaresToGinHandlers(r.lastMiddlewares)
	middlewares = append(middlewares, ginOriginMiddlewares...)
	middlewares = append(middlewares, ginMiddlewares...)
	middlewares = append(middlewares, ginLastMiddlewares...)
	if len(middlewares) > 0 {
		return ginGroup.Use(middlewares...)
	} else {
		return ginGroup
	}
}

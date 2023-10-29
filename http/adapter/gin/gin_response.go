package gin

import (
	"bytes"
	"github.com/connedge/framework/http"
	"github.com/gin-gonic/gin"
	nethttp "net/http"
)

type GinResponse struct {
	instance *gin.Context
	origin   http.ResponseOrigin
}

func NewGinResponse(instance *gin.Context, origin http.ResponseOrigin) *GinResponse {
	return &GinResponse{instance, origin}
}

func (r *GinResponse) Data(code int, contentType string, data []byte) {
	r.instance.Data(code, contentType, data)
}

func (r *GinResponse) Download(filepath, filename string) {
	r.instance.FileAttachment(filepath, filename)
}

func (r *GinResponse) File(filepath string) {
	r.instance.File(filepath)
}

func (r *GinResponse) Header(key, value string) http.Response {
	r.instance.Header(key, value)

	return r
}

func (r *GinResponse) Json(code int, obj any) {
	r.instance.JSON(code, obj)
}

func (r *GinResponse) Origin() http.ResponseOrigin {
	return r.origin
}

func (r *GinResponse) Redirect(code int, location string) {
	r.instance.Redirect(code, location)
}

func (r *GinResponse) String(code int, format string, values ...any) {
	r.instance.String(code, format, values...)
}

func (r *GinResponse) Success() http.ResponseSuccess {
	return NewGinSuccess(r.instance)
}

func (r *GinResponse) Writer() nethttp.ResponseWriter {
	return r.instance.Writer
}

type GinSuccess struct {
	instance *gin.Context
}

func NewGinSuccess(instance *gin.Context) http.ResponseSuccess {
	return &GinSuccess{instance}
}

func (r *GinSuccess) Data(contentType string, data []byte) {
	r.instance.Data(nethttp.StatusOK, contentType, data)
}

func (r *GinSuccess) Json(obj any) {
	r.instance.JSON(nethttp.StatusOK, obj)
}

func (r *GinSuccess) String(format string, values ...any) {
	r.instance.String(nethttp.StatusOK, format, values...)
}

func GinResponseMiddleware() http.Middleware {
	return func(ctx http.Context) {
		blw := &BodyWriter{body: bytes.NewBufferString("")}
		switch ctx := ctx.(type) {
		case *GinContext:
			blw.ResponseWriter = ctx.Instance().Writer
			ctx.Instance().Writer = blw
		}

		ctx.WithValue("responseOrigin", blw)
		ctx.Request().Next()
	}
}

type BodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *BodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)

	return w.ResponseWriter.Write(b)
}

func (w *BodyWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)

	return w.ResponseWriter.WriteString(s)
}

func (w *BodyWriter) Body() *bytes.Buffer {
	return w.body
}

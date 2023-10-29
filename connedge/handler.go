package connedge

import "github.com/connedge/framework/http"

type Handler interface {
	RegisterRoute(engine http.Engine)
}

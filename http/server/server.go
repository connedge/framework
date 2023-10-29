package server

import (
	"github.com/connedge/framework/config"
	"github.com/connedge/framework/http"
	"github.com/connedge/framework/http/adapter/gin"
)

type Server struct {
	engine http.Engine
}

type Config struct {
	HttpEngine http.HttpEngineType
	Config     config.Config
}

func New(s Config) http.Engine {
	switch s.HttpEngine {
	case http.GinServer:
		return gin.NewGin(s.Config)
	}
	return nil
}

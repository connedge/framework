package connedge

import (
	"context"
	"fmt"
	"github.com/connedge/framework/config"
	"github.com/connedge/framework/database"
	"github.com/connedge/framework/http"
	"github.com/connedge/framework/http/server"
	"github.com/connedge/framework/ioc"
	"github.com/gookit/color"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	Config         config.Config
	Container      *ioc.Container
	Engine         http.Engine
	AsyncProviders []Provider
	AppProviders   []Provider
	Database       database.Database
	OnShutdown     func(container *ioc.Container) error
}

const (
	defaultConfigName     = "local"
	defaultConfigType     = "json"
	defaultConfigFilePath = "./configs/"
)

var internalProviders = []Provider{
	database.Provider{},
}

type Config struct {

	// ConfigName is a string variable that holds the name of the configuration file.
	// This will likely be used to load configuration settings from a file.
	// The default configuration name is "local".
	ConfigName string

	// ConfigType is a string variable that holds the type of the configuration file,
	// such as "json" or "yaml".
	// The default configuration type is "json".
	ConfigType string

	// ConfigFilePath is the default file path where configuration files are stored.
	// The default value is "/configs/".
	// Configuration files will be loaded from this path unless explicitly set to a different path.
	ConfigFilePath string

	// HttpEngineType specifies the HTTP engine that connedge uses for
	// server and routing.
	// connedge uses the gin framework by default for its HTTP engine.
	HttpEngineType http.HttpEngineType

	// AsyncRegistries allows building a modular application structure.
	// AsyncRegistries are used with the container (inversion of control) library
	// This enables building the app out of modular, reusable components.
	AsyncProviders []Provider

	// AppRegistries allows building a modular application structure.
	// AppRegistries are used with the container (inversion of control) library
	// This enables building the app out of modular, reusable components.
	AppProviders []Provider

	// OnShutdown is a function that is called when the application shuts down.
	OnShutdown func(container *ioc.Container) error
}

func New(c Config) *App {
	container := ioc.New()
	return &App{
		Config: config.New(
			config.WithConfigName(c.ConfigName, defaultConfigName),
			config.WithConfigType(c.ConfigType, defaultConfigType),

			config.WithFilePath(c.ConfigFilePath, defaultConfigFilePath),
		),
		Engine: server.New(server.Config{
			HttpEngine: c.HttpEngineType,
		}),
		AsyncProviders: c.AsyncProviders,
		AppProviders:   c.AppProviders,
		Container:      container,
		OnShutdown:     c.OnShutdown,
	}
}

func (a *App) Start() {
	color.Cyanf("=============== PROVIDERS ===============\n")
	a.handleInternalAsyncProviders()
	a.handleInitInternalAsyncProviders()
	a.handleAsyncRegister()
	a.handleAppRegister()
	fmt.Println("")
	color.Cyanf("=============== SERVER ===================\n")
	a.gracefullyShutDown()
	_ = a.Engine.Run(fmt.Sprintf(":%s", a.getServerPort()))

}

// getServerPort returns default port or config port
func (a *App) getServerPort() string {
	defaultPort := "3000"
	if p := a.Config.GetString("server.port"); p == "" {
		p = defaultPort
	}
	return defaultPort
}

func (a *App) handleInternalAsyncProviders() {
	for _, provider := range internalProviders {
		a.printProviderName(provider)
		provider.Register(a)
	}
}

func (a *App) handleInitInternalAsyncProviders() {
	a.Database = ioc.MustInvoke[database.Database](a.Container)
}

func (a *App) handleAsyncRegister() {
	for _, provider := range a.AsyncProviders {
		a.printProviderName(provider)
		provider.Register(a)
	}
}

func (a *App) handleAppRegister() {
	for _, provider := range a.AppProviders {
		a.printProviderName(provider)
		provider.Register(a)
	}
}

func (a *App) gracefullyShutDown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		color.Cyanf("\n")
		color.Cyanf("gracefully shutting down...")

		if err := a.OnShutdown(a.Container); err != nil {
			log.Fatalf("app.OnShutdown error: %v", err)
		}

		if err := a.Engine.Shutdown(context.Background()); err != nil {
			log.Fatalf("app.server.shutdown error: %v", err)
		}
	}()
}

func (a *App) printProviderName(registry Provider) {
	if registryWithOptional, ok := registry.(RegistryOptional); ok {
		color.Yellowln("Loaded provider: ", registryWithOptional.Name())
	}
}

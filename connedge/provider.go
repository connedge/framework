package connedge

type Provider interface {

	// The Register method should only be used to bind things into the registry.
	Register(app *App)

	// Handler allows the registry to register any http api handler,
	// When a handler is registered using the 'Handler' method, the corresponding API endpoints will be activated.
	Handler(app *App) []Handler
}

type RegistryOptional interface {
	Name() string
}

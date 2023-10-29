package ioc

type Service[T any] interface {
	getName() string
	getInstance(container *Container) (T, error)
}

type Shutdown interface {
	Shutdown() error
}

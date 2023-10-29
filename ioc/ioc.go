package ioc

import (
	"log"
)

func Bind[T any](i *Container, provider Provider[T]) {
	name := generateServiceName[T]()
	BindNamed[T](i, name, provider)
}

func BindValue[T any](c *Container, value T) {
	name := generateServiceName[T]()
	BindNamedValue[T](c, name, value)
}

func BindNamedValue[T any](c *Container, name string, value T) {
	container := getContainerOrDefault(c)
	if container.exists(name) {
		log.Fatalf("IOC: BindNamedValue `%s` has already been declared", name)
	}

	service := newInstance(name, value)
	container.set(name, service)
	log.Printf("IOC: BindNamedValue service %s injected", name)
}

func BindNamed[T any](c *Container, name string, provider Provider[T]) {
	container := getContainerOrDefault(c)
	if container.exists(name) {
		log.Fatalf("IOC: BindNamed service `%s` has already been declared", name)
	}
	service := newSingleton(name, provider)
	container.set(name, service)

	log.Printf("IOC: BindNamed service %s injected", name)
}

func MustInvoke[T any](c *Container) T {
	s, err := Invoke[T](c)

	name := generateServiceName[T]()

	if err != nil {
		log.Printf("IOC: MustInvoke `%s` failed: %s", name, err)
		return empty[T]()
	}
	return s
}

func Invoke[T any](i *Container) (T, error) {
	name := generateServiceName[T]()
	return InvokeNamed[T](i, name)
}

func InvokeNamed[T any](i *Container, name string) (T, error) {
	return invoke[T](i, name)
}

func invoke[T any](c *Container, name string) (T, error) {
	container := getContainerOrDefault(c)

	serviceAny, ok := container.get(name)
	if !ok {
		return empty[T](), container.serviceNotFound(name)
	}

	service, ok := serviceAny.(Service[T])
	if !ok {
		return empty[T](), container.serviceNotFound(name)
	}

	instance, err := service.getInstance(container)
	if err != nil {
		return empty[T](), err
	}

	return instance, nil
}

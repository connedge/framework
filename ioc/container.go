package ioc

import (
	"fmt"
	"log"
)

var DefaultContainer = New()

type Container struct {
	services map[string]any
}

func getContainerOrDefault(c *Container) *Container {
	if c != nil {
		return c
	}
	return DefaultContainer
}

func (c *Container) set(name string, service any) {
	c.services[name] = service
}

func (c *Container) get(name string) (any, bool) {
	s, ok := c.services[name]
	return s, ok
}

func (c *Container) exists(name string) bool {
	_, ok := c.services[name]
	return ok
}

func (c *Container) Shutdown(name string) error {
	err := c.shutdown(name)
	if err != nil {
		return err
	}
	return nil
}

func (c *Container) serviceNotFound(name string) error {
	return fmt.Errorf("IOC: could not find service `%s`", name)
}

func (c *Container) shutdown(name string) error {
	serviceAny, ok := c.services[name]
	if !ok {
		return fmt.Errorf("container.shutdown could not find service %s", name)
	}

	service, ok := serviceAny.(Shutdown)

	if ok {
		log.Printf("Container.shutdown: requested shutdown for service %s", name)
		err := service.Shutdown()
		if err != nil {
			return err
		}
	}
	delete(c.services, name)
	return nil
}

func New() *Container {
	return &Container{
		services: make(map[string]any),
	}
}

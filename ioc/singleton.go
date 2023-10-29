package ioc

import (
	"sync"
)

type Provider[T any] func(container *Container) (T, error)

type Singleton[T any] struct {
	mu       sync.RWMutex
	name     string
	instance T
	built    bool
	provider Provider[T]
}

func newSingleton[T any](name string, provider Provider[T]) Service[T] {
	return &Singleton[T]{
		name:     name,
		built:    false,
		provider: provider,
	}
}

func (s *Singleton[T]) getName() string {
	return s.name
}

func (s *Singleton[T]) Shutdown() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := any(s.instance).(Shutdown); ok {
		return any(s.instance).(Shutdown).Shutdown()
	} else {
		return nil
	}
}

func (s *Singleton[T]) getInstance(i *Container) (T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.built {
		err := s.build(i)
		if err != nil {
			return empty[T](), err
		}
	}

	return s.instance, nil
}

func (s *Singleton[T]) build(i *Container) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()

	instance, err := s.provider(i)
	if err != nil {
		return err
	}

	s.instance = instance
	s.built = true

	return nil
}

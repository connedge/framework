package ioc

import "sync"

type Instance[T any] struct {
	mu       sync.RWMutex
	name     string
	instance T
	built    bool
}

func newInstance[T any](name string, instance T) *Singleton[T] {
	return &Singleton[T]{
		name:     name,
		built:    false,
		instance: instance,
	}
}

func (s *Instance[T]) getName() string {
	return s.name
}

func (s *Instance[T]) getInstance() (T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.instance, nil
}

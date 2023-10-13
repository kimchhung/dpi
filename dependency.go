package dpi

import (
	"reflect"
	"sync"
)

type Dependency struct {
	name  string
	value any
}

func (s *Dependency) Name() string {
	return s.name
}

func (s *Dependency) Value() any {
	return s.value
}

var dependencyPool = sync.Pool{
	New: func() interface{} {
		return &Dependency{}
	},
}

func toDependency(dependency any) *Dependency {
	if val, ok := dependency.(*Dependency); ok {
		return val
	}

	return WithName("", dependency)
}

// For Manual Mapping `inject:"true", name:"myDep1"`
func WithName(name string, dependency any) *Dependency {
	dep := dependencyPool.Get().(*Dependency)
	dep.name = name
	dep.value = dependency

	if name == "" {
		dep.name = reflect.TypeOf(dependency).String()
	}

	return dep
}

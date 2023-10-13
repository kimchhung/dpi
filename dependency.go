package dpi

import "reflect"

type Dependency map[string]any

func (s Dependency) Name() string {
	for k := range s {
		return k
	}

	return ""
}

func (s Dependency) Value() any {
	for _, v := range s {
		return v
	}

	return nil
}

func toDependency(dependency any) *Dependency {
	if val, ok := dependency.(*Dependency); ok {
		return val
	}

	return WithName("", dependency)
}

// For Manual Mapping `inject:"true", name:"myDep1"`
func WithName(name string, dependency any) *Dependency {
	if name == "" {
		return &Dependency{
			reflect.TypeOf(dependency).String(): dependency,
		}
	}

	return &Dependency{
		name: dependency,
	}
}

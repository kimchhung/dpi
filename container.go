package dpi

import (
	"context"
	"log"
	"os"

	"sync"
)

const prefixName = "\033[35m[dpi]\033[0m "

type Container struct {
	ctx   context.Context
	store map[string]any
	mu    sync.RWMutex
	wg    sync.WaitGroup
	log   *log.Logger
}

type ContainerKey string

var containerKey ContainerKey = "containerKey"

func NewContainer(ctx context.Context) *Container {
	c := &Container{}
	c.log = log.New(os.Stdout, prefixName, 0)
	c.store = make(map[string]any)
	c.ctx = context.WithValue(ctx, containerKey, c)

	go func() {
		<-ctx.Done()
		c.Flush()
	}()

	return c
}

func FromContext(ctx context.Context) *Container {
	if c, ok := ctx.Value(containerKey).(*Container); ok {
		return c
	}
	return NewContainer(ctx)
}

func (c *Container) set(value any) {
	dep := toDependency(value)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[dep.Name()] = dep.Value()
}

func (c *Container) Get(key string) any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.store[key]
}

// wait for lazy injection, and validate
func (c *Container) Wait() {
	c.wg.Wait() // Wait for all goroutines to finish
	count := 0

	for _, dep := range c.store {
		if err := Validate(dep, "lazy"); err != nil {
			panic(err)
		}
		count++
	}

	c.log.Printf("Dependencies: %d", count)
}

func (c *Container) Context() context.Context {
	return c.ctx
}

func (c *Container) Provide(dependencies ...any) *Container {
	for _, _dep := range dependencies {
		c.set(_dep)
	}

	return c
}

func (c *Container) Flush() {
	for key, dep := range c.store {
		if dep != nil {
			dep = nil
		}
		delete(c.store, key)
	}
}

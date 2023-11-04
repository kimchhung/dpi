package dpi

import (
	"context"
	"log"
	"os"

	"sync"
)

const prefixName = "\033[35m[dpi]\033[0m "

type DPI struct {
	ctx   context.Context
	store map[string]any
	mu    sync.RWMutex
	wg    sync.WaitGroup
	log   *log.Logger
}

type DPIKey string
type CleanupFunc func()

var containerKey DPIKey = "containerKey"

func New(ctx context.Context) (*DPI, CleanupFunc) {
	ctx, cancel := context.WithCancel(ctx)
	return NewDPI(ctx), CleanupFunc(cancel)
}

func NewDPI(ctx context.Context) *DPI {
	c := &DPI{}
	c.log = log.New(os.Stdout, prefixName, 0)
	c.store = make(map[string]any)
	c.ctx = context.WithValue(ctx, containerKey, c)

	go func() {
		<-ctx.Done()
		c.Flush()
	}()

	return c
}

func FromContext(ctx context.Context) *DPI {
	if c, ok := ctx.Value(containerKey).(*DPI); ok {
		return c
	}
	return NewDPI(ctx)
}

func (c *DPI) set(value any) {
	dep := toDependency(value)
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[dep.Name()] = dep.Value()
}

func (c *DPI) get(key string) any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.store[key]
}

func (c *DPI) Get(dependency any) any {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if depName, ok := dependency.(string); ok {
		return c.store[depName]
	}

	return c.store[toDependency(dependency).Name()]
}

// wait for lazy injection, and validate
func (c *DPI) Wait() {
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

func (c *DPI) Context() context.Context {
	return c.ctx
}

func (c *DPI) Provide(dependencies ...any) *DPI {
	for _, _dep := range dependencies {
		c.set(_dep)
	}

	return c
}

func (c *DPI) Flush() {
	for key, dep := range c.store {
		if dep != nil {
			dep = nil
		}
		delete(c.store, key)
	}
}

package dpi_test

import (
	"context"
	"testing"

	"github.com/kimchhung/dpi"
)

type DBConn struct {
	Name string
}

func NewDBConn(name string) *DBConn {
	return &DBConn{
		Name: name,
	}
}

type DBComsumer struct {
	DB  *DBConn `inject:"true"`
	DB1 *DBConn `inject:"true" name:"anotherDB"`
}

func NewDBComsumer(ctx context.Context) *DBComsumer {
	s, err := dpi.InjectFromContext(ctx, &DBComsumer{})
	if err != nil {
		panic(err)
	}

	return s
}

type ServiceA struct {
	ServiceB *ServiceB `inject:"true,lazy"`
}

type ServiceB struct {
	ServiceA *ServiceA `inject:"true,lazy"`
}

func NewServiceB(ctx context.Context) *ServiceB {
	s, err := dpi.InjectFromContext(ctx, &ServiceB{})
	if err != nil {
		panic(err)
	}

	return s
}
func NewServiceA(ctx context.Context) *ServiceA {
	s := &ServiceA{}
	if _, err := dpi.InjectFromContext(ctx, s); err != nil {
		panic(err)
	}

	return s
}

func TestLazyInjection(t *testing.T) {
	c, _ := dpi.New(context.Background())
	ctx := c.Context()

	serviceA := NewServiceA(ctx)
	serviceB := NewServiceB(ctx)
	ctx = dpi.ProvideWithContext(ctx, serviceA, serviceB)

	dpi.FromContext(ctx).Wait()

	if serviceA.ServiceB == nil || serviceB.ServiceA == nil {
		t.Errorf("TestLazyInjection did not return services")
	}
}

func BenchmarkInjection(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		c, cleanup := dpi.New(context.Background())
		defer cleanup()

		c.Provide(NewDBConn("defuault"), dpi.WithName("anotherDB", NewDBConn("anotherDB")))
		c.Provide(
			NewDBComsumer(c.Context()),
		)
	}
}

func BenchmarkLazyInjection(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		c, cleanup := dpi.New(context.Background())
		defer cleanup()

		c.Provide(
			NewServiceA(c.Context()),
			NewServiceB(c.Context()),
		)
		c.Wait()
	}
}

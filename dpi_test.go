package dpi_test

import (
	"context"
	"fmt"
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
	injected := dpi.MustInjectFromContext(ctx, new(DBComsumer))
	fmt.Printf("DB: %v \n", injected.DB.Name)
	fmt.Printf("DB1: %v \n", injected.DB1.Name)
	return injected
}

type ServiceA struct {
	ServiceB *ServiceB `inject:"true,lazy"`
}

type ServiceB struct {
	ServiceA *ServiceA `inject:"true,lazy"`
}

func NewServiceB(ctx context.Context) *ServiceB {
	return dpi.MustInjectFromContext(ctx, new(ServiceB))
}

func NewServiceA(ctx context.Context) *ServiceA {
	return dpi.MustInjectFromContext(ctx, new(ServiceA))
}

func TestLazyInjection(t *testing.T) {
	c, _ := dpi.New(context.Background())
	ctx := c.Context()
	ctx = dpi.ProvideWithContext(ctx, NewDBConn("db1"))
	ctx = dpi.ProvideWithContext(ctx, dpi.WithName("db2", NewDBConn("db2")))

	serviceA := NewServiceA(ctx)
	serviceB := NewServiceB(ctx)

	ctx = dpi.ProvideWithContext(ctx,
		NewServiceA(ctx), serviceB)

	dpi.FromContext(ctx).Wait()

	if serviceA.ServiceB == nil || serviceB.ServiceA == nil {
		t.Errorf("TestLazyInjection did not return services")
	}

	t.Run("test extract from ctx", func(t *testing.T) {
		db := dpi.ExtractFromContext[*DBConn](ctx, &DBConn{})
		db2 := dpi.ExtractFromContext[*DBConn](ctx, "db2")

		if db.Name != "db1" {
			t.Errorf("this should be db1")
		}
		if db2.Name != "db2" {
			t.Errorf("this should be db2")
		}
	})
}

func BenchmarkInjection(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		c, cleanup := dpi.New(context.Background())
		defer cleanup()

		c.Provide(
			NewDBConn("hahaha"),
			dpi.WithName("anotherDB", NewDBConn("anotherDB")),
		)
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

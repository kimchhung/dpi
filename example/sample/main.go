package main

import (
	"context"
	"log"

	"github.com/kimchhung/dpi"
)

type DBConn struct {
}

func (d *DBConn) Name() string {
	return "this is DB"
}

type ServiceA struct {
	DB       *DBConn   `inject:"true"`
	ServiceB *ServiceB `inject:"true,lazy"`
}

func NewServiceA(ctx context.Context) *ServiceA {
	s := &ServiceA{}
	if _, err := dpi.InjectFromContext(ctx, s); err != nil {
		panic(err)
	}

	return s
}

func (s *ServiceA) Name() string {
	return "this is Service A"
}

func (s *ServiceA) Print() string {
	log.Printf("from [A]: %s,%s ", s.DB.Name(), s.ServiceB.Name())
	return s.DB.Name()
}

type ServiceB struct {
	DB       *DBConn   `inject:"true"`
	ServiceA *ServiceA `inject:"true,lazy"`
}

func NewServiceB(ctx context.Context) *ServiceB {
	s, err := dpi.InjectFromContext(ctx, &ServiceB{})
	if err != nil {
		panic(err)
	}

	return s
}

func (s *ServiceB) Name() string {
	return "this is Service B"
}

func (s *ServiceB) Print() string {
	log.Printf("from [B]: %s,%s ", s.DB.Name(), s.ServiceA.Name())
	return s.DB.Name()
}

type API struct {
	ServiceA *ServiceA `inject:"true"`
	ServiceB *ServiceB `inject:"true"`
}

func NewAPI(ctx context.Context) *API {
	s := &API{}
	if _, err := dpi.InjectFromContext(ctx, s); err != nil {
		panic(err)
	}

	return s
}

func (api *API) Print() {
	api.ServiceA.Print()
	api.ServiceB.Print()
}

func main() {
	ctx := dpi.ProvideWithContext(context.Background(), &DBConn{})
	ctx = dpi.ProvideWithContext(ctx,
		NewServiceA(ctx),
		NewServiceB(ctx),
	)

	// wait for lazy injection
	dpi.FromContext(ctx).Wait()

	api := NewAPI(ctx)
	api.Print()
}

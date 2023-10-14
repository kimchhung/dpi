package main

import (
	"context"
	"log"
	"sample/database"
	"sample/services"

	"github.com/kimchhung/dpi"
)

type API struct {
	ServiceA *services.ServiceA `inject:"true"`
	ServiceB *services.ServiceB `inject:"true"`

	DBNo *database.DBConn `inject:"true"`
	DB   *database.DBConn `inject:"true" name:"myAnotherDB"`
}

func NewAPI(ctx context.Context) *API {
	return dpi.MustInjectFromContext(ctx, new(API))
}

func (api *API) Print() {
	log.Printf("db no name: %v", api.DBNo.Name())
	log.Printf("db name: %v", api.DB.Name())
	api.ServiceA.Print()
	api.ServiceB.Print()
}

func main() {
	c, cleanup := dpi.New(context.Background())
	defer cleanup()
	// cleanup

	ctx := dpi.ProvideWithContext(
		c.Context(),
		// eg consumer: DB   *database.DBConn `inject:"true"`
		database.New("no name"),

		// eg consumer: DB   *database.DBConn `inject:"true" name:"myAnotherDB"`
		dpi.WithName("myAnotherDB", database.New("with name")),
	)

	ctx = dpi.ProvideWithContext(ctx,
		// eg inject B to A, ServiceB *ServiceB        `inject:"true,lazy"`
		services.NewServiceA(ctx),

		// eg inject A to B, ServiceA *ServiceA        `inject:"true,lazy"`
		services.NewServiceB(ctx),
	)

	api := NewAPI(ctx)

	// wait for lazy injection
	dpi.FromContext(ctx).Wait()

	api.Print()
}

func AlternativeExample() {
	c, cleanup := dpi.New(context.Background())
	defer cleanup()
	// cleanup

	c.Provide(
		database.New("no name"),

		// eg consumer: DB   *database.DBConn `inject:"true" name:"myAnotherDB"`
		dpi.WithName("myAnotherDB", database.New("with name")),
	)

	c.Provide(
		// eg inject B to A, ServiceB *ServiceB        `inject:"true,lazy"`
		services.NewServiceA(c.Context()),

		// eg inject A to B, ServiceA *ServiceA        `inject:"true,lazy"`
		services.NewServiceB(c.Context()),
	)

	api := NewAPI(c.Context())

	// wait for lazy injection
	dpi.FromContext(c.Context()).Wait()
	api.Print()
}

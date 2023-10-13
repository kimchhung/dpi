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
	s := &API{}
	if _, err := dpi.InjectFromContext(ctx, s); err != nil {
		panic(err)
	}

	return s
}

func (api *API) Print() {
	log.Printf("db no name: %v", api.DBNo.Name())
	log.Printf("db name: %v", api.DB.Name())
	api.ServiceA.Print()
	api.ServiceB.Print()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// cleanup

	ctx = dpi.ProvideWithContext(
		ctx,
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

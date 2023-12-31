package services

import (
	"context"
	"log"
	"sample/database"

	"github.com/kimchhung/dpi"
)

type ServiceA struct {
	DB       *database.DBConn `inject:"true"`
	ServiceB *ServiceB        `inject:"true,lazy"`
}

func NewServiceA(ctx context.Context) *ServiceA {
	return dpi.MustInjectFromContext(ctx, new(ServiceA))
}

func (s *ServiceA) Name() string {
	return "this is Service A"
}

func (s *ServiceA) Print() string {
	log.Printf("from [A]: %s,%s ", s.DB.Name(), s.ServiceB.Name())
	return s.DB.Name()
}

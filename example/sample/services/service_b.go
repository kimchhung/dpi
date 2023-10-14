package services

import (
	"context"
	"log"
	"sample/database"

	"github.com/kimchhung/dpi"
)

type ServiceB struct {
	DB       *database.DBConn `inject:"true"`
	ServiceA *ServiceA        `inject:"true,lazy"`
}

func NewServiceB(ctx context.Context) *ServiceB {
	return dpi.MustInjectFromContext(ctx, new(ServiceB))
}

func (s *ServiceB) Name() string {
	return "this is Service B"
}

func (s *ServiceB) Print() string {
	log.Printf("from [B]: %s,%s ", s.DB.Name(), s.ServiceA.Name())
	return s.DB.Name()
}

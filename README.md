# dpi
Simple dependency injection base on context.Context from golang


main.go

```
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

  // consumer
	api := NewAPI(ctx)

	// wait for lazy injection
	dpi.FromContext(ctx).Wait()

	api.Print()
}

```

Services Provider | Consumer
```
type ServiceA struct {
	DB       *database.DBConn `inject:"true"`
	ServiceB *ServiceB        `inject:"true,lazy"` // circular dependency injection
}

func NewServiceA(ctx context.Context) *ServiceA {
	s := &ServiceA{}
	if _, err := dpi.InjectFromContext(ctx, s); err != nil {
		panic(err)
	}

	return s
}

```

```
type ServiceB struct {
	DB       *database.DBConn `inject:"true"`
	ServiceA *ServiceA        `inject:"true,lazy"` // circular dependency injection
}

// use InjectFromContext to extract dependencies from context
func NewServiceA(ctx context.Context) *ServiceA {
	s := &ServiceA{}
	if _, err := dpi.InjectFromContext(ctx, s); err != nil {
		panic(err)
	}

	return s
}
```

Consumer

```
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
```

Benchmark

two simple injections
```
   68794             18141 ns/op            1979 B/op         34 allocs/op
PASS
ok      github.com/kimchhung/dpi        1.506s

```


two lazy injections
```
   41380             27898 ns/op            2054 B/op         41 allocs/op
PASS
ok      github.com/kimchhung/dpi        1.699s

```


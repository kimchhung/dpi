# dpi
Simple dependency injection


// use ProvideWithContext to provice dependencies to context
main.go

```
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// cleanup

  // Provider 
	ctx = dpi.ProvideWithContext(
		ctx,
		database.New("no name"),

		// with custom name
		dpi.WithName("myAnotherDB", database.New("with name")),
	)

  // Provider | Consumer
	ctx = dpi.ProvideWithContext(ctx,
		services.NewServiceA(ctx),
		services.NewServiceB(ctx),
	)
```

Services Provider | Consumer
```
type ServiceB struct {
	DB       *database.DBConn `inject:"true"`
	ServiceA *ServiceA        `inject:"true,lazy"`
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

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

        // for circular dependency injection, need to use inject:"true,lazy"
	ctx = dpi.ProvideWithContext(ctx,
		services.NewServiceA(ctx),
		services.NewServiceB(ctx),
	)
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

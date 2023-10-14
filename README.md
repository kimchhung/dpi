# dpi
Simple dependency injection base on context.Context from golang


main.go

```
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// cleanup

  // as provider
	ctx = dpi.ProvideWithContext(
		ctx,
		// eg consumer: DB   *database.DBConn `inject:"true"`
		database.New("no name"),

		// eg consumer: DB   *database.DBConn `inject:"true" name:"myAnotherDB"`
		dpi.WithName("myAnotherDB", database.New("with name")),
	)

  // as provider and consumer
	ctx = dpi.ProvideWithContext(ctx,
		// eg inject B to A, ServiceB *ServiceB        `inject:"true,lazy"`
		services.NewServiceA(ctx),

		// eg inject A to B, ServiceA *ServiceA        `inject:"true,lazy"`
		services.NewServiceB(ctx),
	)

  // as consumer
	api := NewAPI(ctx)

	// wait for lazy injection
	dpi.FromContext(ctx).Wait()

	api.Print()
}

```
Services Provider | Consumer, A <-> B
```
// consume *ServiceB
type ServiceA struct {
	DB       *database.DBConn `inject:"true"`
	ServiceB *ServiceB        `inject:"true,lazy"` // A <-> B circular dependency injection
}

func NewServiceA(ctx context.Context) *ServiceA {
	return dpi.MustInjectFromContext(ctx, new(ServiceA))
}

// consume *ServiceA
type ServiceB struct {
	DB       *database.DBConn `inject:"true"`
	ServiceA *ServiceA        `inject:"true,lazy"` // A <-> B circular dependency injection
}

// use InjectFromContext to extract dependencies from context
func NewServiceB(ctx context.Context) *ServiceB {
	return dpi.MustInjectFromContext(ctx, new(ServiceB))
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
	return dpi.MustInjectFromContext(ctx, new(API))
}
```

Benchmark

two simple injections
```
[dpi] dpi_test.DBComsumer <- 1/2 `42ns` *dpi_test.DBConn
[dpi] dpi_test.DBComsumer <- 2/2 `41ns` anotherDB
42 ns/op per extraction from ctx

   78619             16272 ns/op            1991 B/op         34 allocs/op
```


two lazy injections
```
[dpi] dpi_test.ServiceB <- 1/1 `42ns` (Lazy) *dpi_test.ServiceA
[dpi] dpi_test.ServiceA <- 1/1 `83ns` (Lazy) *dpi_test.ServiceB
[dpi] dpi_test.ServiceB <- 1/1 `41ns` (Lazy) *dpi_test.ServiceA
[dpi] dpi_test.ServiceA <- 1/1 `84ns` (Lazy) *dpi_test.ServiceB
[dpi] dpi_test.ServiceB <- 1/1 `125ns` (Lazy) *dpi_test.ServiceA
[dpi] dpi_test.ServiceA <- 1/1 `42ns` (Lazy) *dpi_test.ServiceB
<200 ns/op per extraction from ctx for lazy

  48290             22996 ns/op            1818 B/op         41 allocs/op
```


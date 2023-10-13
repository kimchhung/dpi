package dpi

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

// Store dependencies in context automatically, use injection.WithName to store and map base on name
func ProvideWithContext(ctx context.Context, dependencies ...any) context.Context {
	c := FromContext(ctx)
	return c.Provide(dependencies...).Context()
}

// Validate only with tag  "inject:"true"", type ="true|lazy"
func Validate[T any](r T, injectType ...string) error {
	rv := reflect.ValueOf(r).Elem()
	rt := reflect.TypeOf(r).Elem()

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		tag := field.Tag.Get("inject")
		typeInject := "true"
		for _, it := range injectType {
			typeInject = it
		}

		if strings.Contains(tag, "true") {
			isLazyField := strings.Contains(tag, "lazy")
			if typeInject == "true" && !isLazyField {
				if rv.Field(i).IsNil() || rv.Field(i).IsZero() {
					return fmt.Errorf("inject [%v] %v is missing", rt.Name(), field.Name)
				}
			} else if typeInject == "lazy" && isLazyField {
				if rv.Field(i).IsNil() || rv.Field(i).IsZero() {
					return fmt.Errorf("inject [%v] %v is missing", rt.Name(), field.Name)
				}
			}
		}
	}

	return nil
}

// Assign dynamically into fields from context:
// Auto `inject:"true"`,
// Manual `inject:"true", name:"myDep1"`
// Lazy for circle injection `inject:"true,lazy"`,
func InjectFromContext[T any](ctx context.Context, to T) (T, error) {
	c := FromContext(ctx)
	toType := reflect.TypeOf(to).Elem()
	toValue := reflect.ValueOf(to).Elem()

	maxInjection := []uint{0, 0}
	maxLazyInjection := []uint{0, 0}

	cacheField := make(map[int]reflect.StructField)

	for i := 0; i < toType.NumField(); i++ {
		cacheField[i] = toType.Field(i)
		tagName := cacheField[i].Tag.Get("inject")
		if strings.Contains(tagName, "true") {
			if strings.Contains(tagName, "lazy") {
				maxLazyInjection[1]++
			} else {
				maxInjection[1]++
			}
		}
	}

	for i := 0; i < toType.NumField(); i++ {
		tag := cacheField[i].Tag.Get("inject")
		isInject := strings.Contains(tag, "true")

		if isInject {
			tagName := cacheField[i].Tag.Get("name")
			isLazy := strings.Contains(tag, "lazy")

			if tagName == "" {
				tagName = cacheField[i].Type.String()
			}

			assignValue := func(toFieldNumber int, service any, startTime time.Time) {
				sv := reflect.ValueOf(service)
				toValue.Field(toFieldNumber).Set(sv)

				if isLazy {
					maxLazyInjection[0]++
					c.log.Printf("%s <- %d/%d `%v` (Lazy) %s", toType, maxLazyInjection[0], maxLazyInjection[1], time.Since(startTime), tagName)
				} else {
					maxInjection[0]++
					c.log.Printf("%s <- %d/%d `%v` %s", toType, maxInjection[0], maxInjection[1], time.Since(startTime), tagName)
				}
			}

			getValueFromContext := func(fieldNumber int, startTime time.Time) any {
				if service := c.Get(tagName); service != nil {
					assignValue(fieldNumber, service, startTime)
					return service
				}

				return nil
			}

			if isLazy {
				c.wg.Add(1)
				go func(_i int, _getValueFromContext func(i int, startTime time.Time) any) {
					defer c.wg.Done()
					startTime := time.Now()
					for service := _getValueFromContext(_i, startTime); service == nil; service = _getValueFromContext(_i, startTime) {
						time.Sleep(time.Duration(rand.Intn(1000)) * time.Nanosecond)
					}
				}(i, getValueFromContext)
			} else {
				getValueFromContext(i, time.Now())
			}
		}
	}

	return to, Validate(to)
}

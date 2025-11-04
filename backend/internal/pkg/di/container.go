// Package di is a lightweight yet powerful IoC container for Go projects.
// It provides an easy-to-use interface and performance-in-mind container to be your ultimate requirement.
package di

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

// binding holds a resolver and a concrete (if already resolved).
// It is the break for the Container wall!
type binding struct {
	resolver    interface{} // resolver is the function that is responsible for making the concrete.
	concrete    interface{} // concrete is the stored instance for singleton bindings.
	isSingleton bool        // isSingleton is true if the binding is a singleton.
}

// make resolves the binding if needed and returns the resolved concrete.
func (b *binding) make(c *Container) (interface{}, error) {
	if b.concrete != nil {
		return b.concrete, nil
	}

	retVal, err := c.invoke(b.resolver)
	if b.isSingleton {
		b.concrete = retVal
	}

	return retVal, err
}

// Container holds the bindings and provides methods to interact with them.
// It is the entry point in the package.
type Container struct {
	father   *Container
	bindings map[reflect.Type]map[string]*binding
}

type NamedParam string

// New creates a new concrete of the Container.
func New() *Container {
	return &Container{
		father:   nil,
		bindings: make(map[reflect.Type]map[string]*binding),
	}
}

// Child creates a child container that contains all the fathers bindings.
func (c *Container) Child() *Container {
	return &Container{
		father:   c,
		bindings: make(map[reflect.Type]map[string]*binding),
	}
}

// bind maps an abstraction to concrete and instantiates if it is a singleton binding.
func (c *Container) bind(resolver interface{}, name string, isSingleton bool, isLazy bool) error {
	reflectedResolver := reflect.TypeOf(resolver)
	if reflectedResolver.Kind() != reflect.Func {
		return errors.New("container: the resolver must be a function")
	}

	if reflectedResolver.NumOut() > 0 {
		if _, exist := c.bindings[reflectedResolver.Out(0)]; !exist {
			c.bindings[reflectedResolver.Out(0)] = make(map[string]*binding)
		}
	}

	if err := c.validateResolverFunction(reflectedResolver); err != nil {
		return err
	}

	var concrete interface{}
	if !isLazy {
		var err error
		concrete, err = c.invoke(resolver)
		if err != nil {
			return err
		}
	}

	if isSingleton {
		c.bindings[reflectedResolver.Out(0)][name] = &binding{resolver: resolver, concrete: concrete, isSingleton: isSingleton}
	} else {
		c.bindings[reflectedResolver.Out(0)][name] = &binding{resolver: resolver, isSingleton: isSingleton}
	}

	return nil
}

func (c *Container) validateResolverFunction(funcType reflect.Type) error {
	retCount := funcType.NumOut()

	if retCount == 0 || retCount > 2 {
		return errors.New("container: resolver function signature is invalid - it must return abstract, or abstract and error")
	}

	resolveType := funcType.Out(0)
	for i := 0; i < funcType.NumIn(); i++ {
		if funcType.In(i) == resolveType {
			return fmt.Errorf("container: resolver function signature is invalid - depends on abstract it returns")
		}
	}

	return nil
}

// invoke calls a function and its returned values.
// It only accepts one value and an optional error.
func (c *Container) invoke(function interface{}) (interface{}, error) {
	arguments, err := c.arguments(function)
	if err != nil {
		return nil, err
	}

	values := reflect.ValueOf(function).Call(arguments)
	if len(values) == 2 && values[1].CanInterface() {
		if err, ok := values[1].Interface().(error); ok {
			return values[0].Interface(), err
		}
	}
	return values[0].Interface(), nil
}

func (c *Container) findConcrete(abstraction reflect.Type, name string) (interface{}, error) {
	ctnr := c
	for ctnr != nil {
		if concrete, exist := ctnr.bindings[abstraction][name]; exist {
			instance, err := concrete.make(ctnr)
			if err != nil {
				return nil, fmt.Errorf("container: encountered error while making concrete for: %s. Error encountered: %w", abstraction.String(), err)
			}

			return instance, nil
		}
		ctnr = ctnr.father
	}

	return nil, errors.New("container: no concrete found for: " + abstraction.String())
}

// arguments returns the list of resolved arguments for a function.
func (c *Container) arguments(function interface{}) ([]reflect.Value, error) {
	reflectedFunction := reflect.TypeOf(function)
	argumentsCount := reflectedFunction.NumIn()
	arguments := make([]reflect.Value, argumentsCount)

	for i := 0; i < argumentsCount; i++ {
		abstraction := reflectedFunction.In(i)
		if i+1 >= argumentsCount || reflectedFunction.In(i+1).Kind() != reflect.String {
			instance, err := c.findConcrete(abstraction, "")
			if err != nil {
				return nil, err
			}
			arguments[i] = reflect.ValueOf(instance)
		} else {
			abstractionName := reflectedFunction.In(i + 1)
			instance, err := c.findConcrete(abstraction, abstractionName.String())
			if err != nil {
				return nil, err
			}
			arguments[i] = reflect.ValueOf(instance)
			i++
		}
	}

	return arguments, nil
}

// Reset deletes all the existing bindings and empties the container.
func (c *Container) reset() {
	for k := range c.bindings {
		delete(c.bindings, k)
	}
}

// Call takes a receiver function with one or more arguments of the abstractions (interfaces).
// It invokes the receiver function and passes the related concretes.
func (c *Container) call(function interface{}) error {
	receiverType := reflect.TypeOf(function)
	if receiverType == nil || receiverType.Kind() != reflect.Func {
		return errors.New("container: invalid function")
	}

	arguments, err := c.arguments(function)
	if err != nil {
		return err
	}

	result := reflect.ValueOf(function).Call(arguments)

	if len(result) == 0 {
		return nil
	} else if len(result) == 1 && result[0].CanInterface() {
		if result[0].IsNil() {
			return nil
		}
		if err, ok := result[0].Interface().(error); ok {
			return err
		}
	}

	return errors.New("container: receiver function signature is invalid")
}

// Resolve takes an abstraction (reference of an interface type) and fills it with the related concrete.
func (c *Container) resolve(abstraction interface{}) error {
	return c.namedResolve(abstraction, "")
}

// NamedResolve takes abstraction and its name and fills it with the related concrete.
func (c *Container) namedResolve(abstraction interface{}, name string) error {
	receiverType := reflect.TypeOf(abstraction)
	if receiverType == nil {
		return errors.New("container: invalid abstraction")
	}

	if receiverType.Kind() == reflect.Pointer {
		elem := receiverType.Elem()

		instance, err := c.findConcrete(elem, name)
		if err != nil {
			return err
		}

		reflect.ValueOf(abstraction).Elem().Set(reflect.ValueOf(instance))
		return nil
	}

	return errors.New("container: invalid abstraction")
}

// Fill takes a struct and resolves the fields with the tag `container:"inject"`
func (c *Container) fill(structure interface{}) error {
	receiverType := reflect.TypeOf(structure)
	if receiverType == nil {
		return errors.New("container: invalid structure")
	}

	if receiverType.Kind() == reflect.Pointer {
		elem := receiverType.Elem()
		if elem.Kind() == reflect.Struct {
			s := reflect.ValueOf(structure).Elem()

			for i := 0; i < s.NumField(); i++ {
				f := s.Field(i)

				if t, exist := s.Type().Field(i).Tag.Lookup("container"); exist {
					var name string

					switch t {
					case "type":
						name = ""
					case "name":
						name = s.Type().Field(i).Name
					default:
						return fmt.Errorf("container: %v has an invalid struct tag", s.Type().Field(i).Name)
					}

					instance, err := c.findConcrete(f.Type(), name)
					if err != nil {
						return fmt.Errorf("container: cannot make %v field: %w", s.Type().Field(i).Name, err)
					}

					ptr := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
					ptr.Set(reflect.ValueOf(instance))
				}
			}

			return nil
		}
	}

	return errors.New("container: invalid structure")
}

// SingletonLazy binds an abstraction to concrete lazily in singleton mode.
// The concrete is resolved only when the abstraction is resolved for the first time.
// It takes a resolver function that returns the concrete, and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have been declared in the Container already.
func (c *Container) singletonLazy(resolver interface{}) error {
	return c.bind(resolver, "", true, true)
}

// NamedSingleton binds a named abstraction to concrete lazily in singleton mode.
// The concrete is resolved only when the abstraction is resolved for the first time.
func (c *Container) namedSingletonLazy(name string, resolver interface{}) error {
	return c.bind(resolver, name, true, true)
}

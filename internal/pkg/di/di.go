package di

// Reset deletes all the existing bindings and empties the container.
func Reset(ctnr *Container) {
	ctnr.reset()
}

// Provide binds an abstraction to concrete lazily in singleton mode.
// The concrete is resolved only when the abstraction is resolved for the first time.
// It takes a resolver function that returns the concrete, and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have been declared in the Container already.
func Provide(ctnr *Container, resolver any) {
	err := ctnr.singletonLazy(resolver)
	if err != nil {
		panic(err)
	}
}

// ProvideNamed binds a named abstraction to concrete lazily in singleton mode.
// The concrete is resolved only when the abstraction is resolved for the first time.
func ProvideNamed(ctnr *Container, name string, resolver any) {
	err := ctnr.namedSingletonLazy(name, resolver)
	if err != nil {
		panic(err)
	}
}

// NamedResolve takes abstraction and its name and fills it with the related concrete.
func NamedResolve[T any](ctnr *Container, name string) T {
	var obj T
	err := ctnr.namedResolve(&obj, name)
	if err != nil {
		panic(err)
	}

	return obj
}

// Resolve takes an abstraction (reference of an interface type) and fills it with the related concrete.
func Resolve[T any](ctnr *Container) T {
	var obj T
	err := ctnr.resolve(&obj)
	if err != nil {
		panic(err)
	}

	return obj
}

// Fill takes a struct and resolves the fields with the tag `container:"inject"`
func Fill(ctnr *Container, structure any) {
	err := ctnr.fill(structure)
	if err != nil {
		panic(err)
	}
}

// Call takes a receiver function with one or more arguments of the abstractions (interfaces).
// It invokes the receiver function and passes the related concretes.
func Call(ctnr *Container, function any) error {
	return ctnr.call(function)
}

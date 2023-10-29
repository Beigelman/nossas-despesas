package di_test

import (
	"errors"
	"testing"

	"github.com/Beigelman/ludaapi/internal/pkg/di"
	"github.com/stretchr/testify/assert"
)

type Shape interface {
	SetArea(int)
	GetArea() int
}

type Circle struct {
	a int
}

func (c *Circle) SetArea(a int) {
	c.a = a
}

func (c *Circle) GetArea() int {
	return c.a
}

type Database interface {
	Connect() bool
}

type MySQL struct{}

func (m MySQL) Connect() bool {
	return true
}

var container = di.New()

func TestContainer_Provide(t *testing.T) {
	defer func() {
		r := recover()
		assert.Nil(t, r)
	}()

	di.Provide(container, func() Shape {
		return &Circle{a: 13}
	})

	assert.NoError(t, di.Call(container, func(s1 Shape) {
		s1.SetArea(666)
	}))

	assert.NoError(t, di.Call(container, func(s1 Shape) {
		a := s1.GetArea()
		assert.Equal(t, a, 666)
	}))
}

func TestContainer_SingletonLazy(t *testing.T) {
	defer func() {
		r := recover()
		assert.Nil(t, r)
	}()

	di.Provide(container, func() Shape {
		return &Circle{a: 13}
	})

	assert.NoError(t, di.Call(container, func(s1 Shape) {
		s1.SetArea(666)
	}))

	assert.NoError(t, di.Call(container, func(s2 Shape) {
		a := s2.GetArea()
		assert.Equal(t, a, 666)
	}))
}

func TestContainer_SingletonLazy_With_Resolve_That_Returns_Nothing(t *testing.T) {
	defer func() {
		r := recover()
		assert.Error(t, r.(error), "container: resolver function signature is invalid")
	}()

	di.Provide(container, func() {})
}

func TestContainer_SingletonLazy_With_Resolve_That_Returns_Error(t *testing.T) {
	defer func() {
		r := recover()
		assert.Error(t, r.(error), "app: error")
	}()

	di.Provide(container, func() (Shape, error) {
		return nil, errors.New("app: error")
	})

	s := di.Resolve[Shape](container)
	s.GetArea()
}

func TestContainer_SingletonLazy_With_NonFunction_Resolver_It_Should_Fail(t *testing.T) {
	defer func() {
		r := recover()
		assert.Error(t, r.(error), "container: the resolver must be a function")
	}()

	di.Provide(container, "STRING!")
}

func TestContainer_SingletonLazy_With_Resolvable_Arguments(t *testing.T) {
	defer func() {
		r := recover()
		assert.Nil(t, r)
	}()

	di.Provide(container, func() Shape {
		return &Circle{a: 666}
	})

	di.Provide(container, func(s Shape) Database {
		assert.Equal(t, s.GetArea(), 666)
		return &MySQL{}
	})

	s := di.Resolve[Shape](container)
	s.GetArea()
}

func TestContainer_SingletonLazy_With_Non_Resolvable_Arguments(t *testing.T) {
	defer func() {
		r := recover()
		assert.EqualError(t, r.(error), "container: resolver function signature is invalid - depends on abstract it returns")
	}()

	di.Reset(container)

	di.Provide(container, func(s Shape) Shape {
		return &Circle{a: s.GetArea()}
	})
}

func TestContainer_NamedSingletonLazy(t *testing.T) {
	defer func() {
		r := recover()
		assert.Nil(t, r)
	}()

	di.ProvideNamed(container, "theCircle", func() Shape {
		return &Circle{a: 13}
	})

	sh := di.NamedResolve[Shape](container, "theCircle")
	assert.Equal(t, sh.GetArea(), 13)
}

func TestContainer_Call_With_Multiple_Resolving(t *testing.T) {
	defer func() {
		r := recover()
		assert.Nil(t, r)
	}()

	di.Provide(container, func() Shape {
		return &Circle{a: 5}
	})

	di.Provide(container, func() Database {
		return &MySQL{}
	})

	assert.NoError(t, di.Call(container, func(s Shape, m Database) {
		if _, ok := s.(*Circle); !ok {
			t.Error("Expected Circle")
		}

		if _, ok := m.(*MySQL); !ok {
			t.Error("Expected MySQL")
		}
	}))
}

func TestContainer_Call_With_Dependency_Missing_In_Chain(t *testing.T) {
	defer func() {
		r := recover()
		assert.Error(t, r.(error), "container: no concrete found for: di_test.Shape")
	}()

	var container = di.New()
	di.Provide(container, func() (Database, error) {
		di.Resolve[Shape](container)

		return &MySQL{}, nil
	})

	assert.NoError(t, di.Call(container, func(m Database) {
		if _, ok := m.(*MySQL); !ok {
			t.Error("Expected MySQL")
		}
	}))
}

func TestContainer_Call_With_Unsupported_Receiver_It_Should_Fail(t *testing.T) {
	err := di.Call(container, "STRING!")
	assert.Error(t, err, "container: invalid function")
}

func TestContainer_Call_With_Second_UnBounded_Argument(t *testing.T) {
	di.Reset(container)

	di.Provide(container, func() Shape {
		return &Circle{}
	})

	err := di.Call(container, func(s Shape, d Database) {})
	assert.EqualError(t, err, "container: no concrete found for: di_test.Database")
}

func TestContainer_Call_With_A_Returning_Error(t *testing.T) {
	di.Reset(container)

	di.Provide(container, func() Shape {
		return &Circle{}
	})

	err := di.Call(container, func(s Shape) error {
		return errors.New("app: some context error")
	})
	assert.EqualError(t, err, "app: some context error")
}

func TestContainer_Call_With_A_Returning_Nil_Error(t *testing.T) {
	di.Reset(container)

	di.Provide(container, func() Shape {
		return &Circle{}
	})

	err := di.Call(container, func(s Shape) error {
		return nil
	})
	assert.Nil(t, err)
}

func TestContainer_Call_With_Invalid_Signature(t *testing.T) {
	di.Reset(container)

	di.Provide(container, func() Shape {
		return &Circle{}
	})

	err := di.Call(container, func(s Shape) (int, error) {
		return 13, errors.New("app: some context error")
	})
	assert.EqualError(t, err, "container: receiver function signature is invalid")
}

func TestContainer_Resolve_With_Reference_As_Resolver(t *testing.T) {
	di.Provide(container, func() Shape {
		return &Circle{a: 5}
	})

	di.Provide(container, func() Database {
		return &MySQL{}
	})

	s := di.Resolve[Shape](container)
	if _, ok := s.(*Circle); !ok {
		t.Error("Expected Circle")
	}

	d := di.Resolve[Database](container)
	if _, ok := d.(*MySQL); !ok {
		t.Error("Expected MySQL")
	}
}

func TestContainer_Resolve_With_Unsupported_Receiver_It_Should_Fail(t *testing.T) {
	defer func() {
		r := recover()
		assert.Error(t, r.(error), "container: invalid abstraction")
	}()

	di.Resolve[string](container)
}

func TestContainer_Resolve_With_NonReference_Receiver_It_Should_Fail(t *testing.T) {
	defer func() {
		r := recover()
		assert.Error(t, r.(error), "container: invalid abstraction")
	}()

	s := di.Resolve[Shape](di.New())
	s.GetArea()
}

func TestContainer_Resolve_With_UnBounded_Reference_It_Should_Fail(t *testing.T) {
	defer func() {
		r := recover()
		assert.Error(t, r.(error), "container: no concrete found for: di_test.Shape")
	}()

	di.Reset(container)

	s := di.Resolve[Shape](container)
	s.GetArea()
}

func TestContainer_Fill_With_Struct_Pointer(t *testing.T) {
	defer func() {
		r := recover()
		assert.Nil(t, r)
	}()

	di.Provide(container, func() Shape {
		return &Circle{a: 5}
	})

	di.ProvideNamed(container, "C", func() Shape {
		return &Circle{a: 5}
	})

	di.Provide(container, func() Database {
		return &MySQL{}
	})

	myApp := struct {
		S Shape    `container:"type"`
		D Database `container:"type"`
		C Shape    `container:"name"`
		X string
	}{}

	di.Fill(container, &myApp)

	assert.IsType(t, &Circle{}, myApp.S)
	assert.IsType(t, &MySQL{}, myApp.D)
}

func TestContainer_Fill_Unexported_With_Struct_Pointer(t *testing.T) {
	defer func() {
		r := recover()
		assert.Nil(t, r)
	}()

	di.Provide(container, func() Shape {
		return &Circle{a: 5}
	})

	di.Provide(container, func() Database {
		return &MySQL{}
	})

	myApp := struct {
		s Shape    `container:"type"`
		d Database `container:"type"`
		y int
	}{}

	di.Fill(container, &myApp)

	assert.IsType(t, &Circle{}, myApp.s)
	assert.IsType(t, &MySQL{}, myApp.d)
}

func TestContainer_Fill_With_Invalid_Field_It_Should_Fail(t *testing.T) {
	defer func() {
		r := recover()
		assert.Error(t, r.(error), "container: cannot make S field")
	}()

	di.ProvideNamed(container, "C", func() Shape {
		return &Circle{a: 5}
	})

	type App struct {
		S string `container:"name"`
	}

	myApp := App{}

	di.Fill(container, &myApp)
}

func TestContainer_Fill_With_Invalid_Tag_It_Should_Fail(t *testing.T) {
	defer func() {
		r := recover()
		assert.Error(t, r.(error), "container: S has an invalid struct tag")
	}()

	type App struct {
		S string `container:"invalid"`
	}

	myApp := App{}

	di.Fill(container, &myApp)
}

func TestContainer_Fill_With_Invalid_Field_Name_It_Should_Fail(t *testing.T) {
	defer func() {
		r := recover()
		assert.Error(t, r.(error), "container: cannot make S field")
	}()

	type App struct {
		S string `container:"name"`
	}

	myApp := App{}

	di.Fill(container, &myApp)
}

func TestContainer_Fill_With_Invalid_Struct_It_Should_Fail(t *testing.T) {
	defer func() {
		r := recover()
		assert.Error(t, r.(error), "container: invalid structure")
	}()

	invalidStruct := 0
	di.Fill(container, &invalidStruct)
}

func TestContainer_Fill_With_Invalid_Pointer_It_Should_Fail(t *testing.T) {
	defer func() {
		r := recover()
		assert.Error(t, r.(error), "container: invalid structure")
	}()

	var s Shape
	di.Fill(container, s)
}

func TestContainer_Fill_With_Dependency_Missing_In_Chain(t *testing.T) {
	defer func() {
		r := recover()
		assert.Error(t, r.(error), "container: no concrete found for: di_test.Shape")
	}()

	var container = di.New()
	di.Provide(container, func() Shape {
		return &Circle{a: 5}
	})

	di.ProvideNamed(container, "C", func() (Shape, error) {
		s := di.NamedResolve[Shape](container, "foo")
		s.GetArea()

		return &Circle{a: 5}, nil
	})

	di.Provide(container, func() Database {
		return &MySQL{}
	})

	myApp := struct {
		S Shape    `container:"type"`
		D Database `container:"type"`
		C Shape    `container:"name"`
		X string
	}{}

	di.Fill(container, &myApp)
}

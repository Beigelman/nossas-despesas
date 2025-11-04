package eon

import (
	"context"
	"testing"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/di"
	"github.com/stretchr/testify/suite"
)

type PrintTest func() string

var module1 = NewModule("module1", func(ctx context.Context, c *di.Container, lfcm LifeCycleManager, app Info) {
	di.Provide(c, func() PrintTest {
		return func() string {
			return "Test"
		}
	})
})

type PrintBanana func() string

var module2 = NewModule("module2", func(ctx context.Context, c *di.Container, lfcm LifeCycleManager, app Info) {
	di.Provide(c, func() PrintBanana {
		return func() string {
			return "Banana"
		}
	})
})

var errorModule = NewModule("error", func(ctx context.Context, c *di.Container, lfcm LifeCycleManager, app Info) {
	di.Provide(c, func(a string) PrintBanana {
		return func() string {
			return "Banana" + a
		}
	})

	lfcm.OnBooted(HookOrders.APPEND, func() error {
		a := di.Resolve[PrintBanana](c)
		a()
		return nil
	})
})

type EonAppTestSuite struct {
	suite.Suite
	app *Application
}

func TestEonApp(t *testing.T) {
	suite.Run(t, new(EonAppTestSuite))
}

func (suite *EonAppTestSuite) SetupTest() {
	suite.app = New("Test", WithShutdownTime(100*time.Millisecond))
}

func (suite *EonAppTestSuite) TestEonApp_StartingTheApplication() {
	suite.app.BootStrap(module1, module2)

	err := suite.app.StartTest()
	suite.NoError(err)

	printTest := di.Resolve[PrintTest](suite.app.ctnr)
	suite.Equal(printTest(), "Test")

	printBanana := di.Resolve[PrintBanana](suite.app.ctnr)
	suite.Equal(printBanana(), "Banana")
}

func (suite *EonAppTestSuite) TestEonApp_LoadingABadModule() {
	defer func() {
		r := recover()
		suite.NotNil(r)
		suite.ErrorContains(r.(error), "container: no concrete found for: string")
	}()

	err := suite.app.BootStrap(module1, module2, errorModule).StartTest()
	suite.EqualError(err, "Error loading module")
}

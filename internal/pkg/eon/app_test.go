package eon

import (
	"context"
	"testing"
	"time"

	"github.com/Beigelman/ludaapi/internal/pkg/di"
	"github.com/Beigelman/ludaapi/internal/pkg/env"
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
})

type EonAppTestSuite struct {
	suite.Suite
	app Application
}

func TestEonApp(t *testing.T) {
	suite.Run(t, new(EonAppTestSuite))
}

func (suite *EonAppTestSuite) SetupTest() {
	suite.app = NewApp(Config{
		ShutdownTime: 2 * time.Second,
		ServiceName:  "Test",
		Env:          env.Development,
		IoC:          di.New(),
	})
}

func (suite *EonAppTestSuite) TestStartingTheApplication() {
	suite.app.BootStrap(module1, module2)

	err := suite.app.StartTest()
	suite.NoError(err)

	printTest := di.Resolve[PrintTest](suite.app.container())
	suite.Equal(printTest(), "Test")

	printBanana := di.Resolve[PrintBanana](suite.app.container())
	suite.Equal(printBanana(), "Banana")
}

func (suite *EonAppTestSuite) TestLoadingABadModule() {
	defer func() {
		r := recover()
		suite.NotNil(r)
		suite.ErrorContains(r.(error), "container: no concrete found for: string")
	}()

	err := suite.app.BootStrap(module1, module2, errorModule).Start()
	suite.EqualError(err, "Error loading module")
}

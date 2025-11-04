package eon

import (
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type LifeCycleManagerTestSuite struct {
	suite.Suite
	lfcm *lifeCycleManager
}

func TestLifeCycleManager(t *testing.T) {
	suite.Run(t, new(LifeCycleManagerTestSuite))
}

func (suite *LifeCycleManagerTestSuite) SetupTest() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	suite.lfcm = newLifeCycleManager(2*time.Second, logger)
}

func (suite *LifeCycleManagerTestSuite) TestLifeCycleManager_SuccessfulStart() {
	suite.lfcm.OnBooting(HookOrders.APPEND, func() error {
		return nil
	})

	err := suite.lfcm.start()
	suite.NoError(err)
}

func (suite *LifeCycleManagerTestSuite) TestLifeCycleManager_WithDoubleStart() {
	suite.NoError(suite.lfcm.start())
	err := suite.lfcm.start()
	suite.EqualError(err, "the application has already started")
}

func (suite *LifeCycleManagerTestSuite) TestLifeCycleManager_WithStartError() {
	suite.lfcm.OnBooting(HookOrders.APPEND, func() error {
		return errors.New("error running cycle")
	})
	err := suite.lfcm.start()
	suite.EqualError(err, "transitioning the life cycles: processing onBooting: error running cycle")
}

func (suite *LifeCycleManagerTestSuite) TestLifeCycleManager_StartingAndStoppingTheApp() {
	suite.lfcm.OnBooting(HookOrders.APPEND, func() error {
		return nil
	})

	suite.lfcm.OnDisposing(HookOrders.APPEND, func() error {
		return nil
	})

	err := suite.lfcm.start()
	suite.NoError(err)
	err = suite.lfcm.stop()
	suite.NoError(err)
}

func (suite *LifeCycleManagerTestSuite) TestLifeCycleManager_StoppingTheApplicationThatHasNotStarted() {
	err := suite.lfcm.stop()
	suite.EqualError(err, "the application is not running")
}

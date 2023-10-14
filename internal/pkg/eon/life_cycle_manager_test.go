package eon

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/truepay/go-commons/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type LifeCycleManagerTestSuite struct {
	suite.Suite
	observedZapCore zapcore.Core
	observedLogs    *observer.ObservedLogs
	observedLogger  *zap.Logger
	lfcm            lifeCycleManager
}

func TestLifeCycleManager(t *testing.T) {
	suite.Run(t, new(LifeCycleManagerTestSuite))
}

func (suite *LifeCycleManagerTestSuite) SetupTest() {
	suite.observedZapCore, suite.observedLogs = observer.New(zap.InfoLevel)
	suite.observedLogger = zap.New(suite.observedZapCore)
	logger, err := logger.NewCustomLogger("TestLifeCycle", suite.observedLogger)
	suite.NoError(err)
	suite.lfcm = newLifeCycleManager(2*time.Second, logger, context.TODO())
}

func (suite *LifeCycleManagerTestSuite) TestStartingTheApplication() {
	loadCycles(suite.lfcm, suite.observedLogger, nil)

	err := suite.lfcm.start(context.TODO())
	suite.NoError(err)

	logs := suite.observedLogs.All()

	suite.Len(logs, 11)
	suite.Equal("Application STARTING", logs[0].Message)
	suite.Equal("Application STARTED", logs[10].Message)
}

func (suite *LifeCycleManagerTestSuite) TestStartingTheApplicationWithDoubleStart() {
	suite.lfcm.start(context.TODO())
	err := suite.lfcm.start(context.TODO())
	suite.EqualError(err, "the application has already started")
}

func (suite *LifeCycleManagerTestSuite) TestStartingTheApplicationWithStartError() {
	loadCycles(suite.lfcm, suite.observedLogger, errors.New("error running cycle"))

	err := suite.lfcm.start(context.TODO())
	suite.EqualError(err, "error running cycle")

	logs := suite.observedLogs.All()
	suite.Len(logs, 12)
	suite.Equal("Application STARTING", logs[0].Message)
	suite.Equal("Error stopping the application on start failure", logs[11].Message)
}

func (suite *LifeCycleManagerTestSuite) TestLifeCycleManagerLoadingEveryModuleInOrder() {
	loadCycles(suite.lfcm, suite.observedLogger, nil)

	err := suite.lfcm.start(context.TODO())
	suite.NoError(err)
	err = suite.lfcm.stop(context.TODO())
	suite.NoError(err)

	logs := suite.observedLogs.All()

	suite.Len(logs, 18)
	suite.Equal("Application STARTING", logs[0].Message)
	suite.Equal("Application exited with success", logs[len(logs)-1].Message)
}

func (suite *LifeCycleManagerTestSuite) TestLifeCycleManagerWithStartError() {
	loadCycles(suite.lfcm, suite.observedLogger, errors.New("error running cycle"))

	err := suite.lfcm.start(context.TODO())
	suite.EqualError(err, "error running cycle")

	logs := suite.observedLogs.All()

	suite.Len(logs, 12)
	suite.Equal(logs[6].Message, "Error transitioning the life cycles")
}

func (suite *LifeCycleManagerTestSuite) TestLifeCycleManagerStoppingTheApplicationThatHasNotStarted() {
	err := suite.lfcm.stop(context.TODO())
	suite.EqualError(err, "the application is not running")
}

var logMessages = []string{
	"Running onBooting",
	"Running onBooted 1",
	"Running onBooted 2",
	"Running onReady",
	"Running onRunning",
	"Application started",
	"Running onDisposing",
	"Running onDisposed",
}

func loadCycles(lfcm lifeCycleManager, logger *zap.Logger, err error) {
	lfcm.OnBooting(HookOrders.APPEND, func() error {
		logger.Info(logMessages[0])
		return nil
	})

	lfcm.OnBooted(HookOrders.APPEND, func() error {
		logger.Info(logMessages[2])
		return err
	})

	lfcm.OnBooted(HookOrders.PREPEND, func() error {
		logger.Info(logMessages[1])
		return nil
	})

	lfcm.OnReady(HookOrders.APPEND, func() error {
		logger.Info(logMessages[3])
		return nil
	})

	lfcm.OnRunning(HookOrders.APPEND, func() error {
		logger.Info(logMessages[4])
		return nil
	})

	lfcm.OnDisposing(HookOrders.APPEND, func() error {
		logger.Info(logMessages[6])
		return err
	})

	lfcm.OnDisposed(HookOrders.APPEND, func() error {
		logger.Info(logMessages[7])
		return nil
	})
}

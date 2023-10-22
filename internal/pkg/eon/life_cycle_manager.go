package eon

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"
)

type appState string

var appStates = struct {
	IDLE     appState
	STARTING appState
	STARTED  appState
	STOPPING appState
	STOPPED  appState
}{
	IDLE:     "IDLE",
	STARTING: "STARTING",
	STARTED:  "STARTED",
	STOPPING: "STOPPING",
	STOPPED:  "STOPPED ",
}

type lifeCycleManager interface {
	getState() appState
	start() error
	stop() error
	OnBooting(order HookOrder, fn ...HookFn)
	OnBooted(order HookOrder, fn ...HookFn)
	OnReady(order HookOrder, fn ...HookFn)
	OnRunning(order HookOrder, fn ...HookFn)
	OnDisposing(order HookOrder, fn ...HookFn)
	OnDisposed(order HookOrder, fn ...HookFn)
	shutdown()
}

func newLifeCycleManager(shutdownTime time.Duration, logger Logger) lifeCycleManager {
	return &lifeCycleManagerImpl{
		state:         appStates.IDLE,
		hooks:         newHookStore(),
		logger:        logger,
		forceShutdown: false,
		shutdownTime:  shutdownTime,
	}
}

type lifeCycleManagerImpl struct {
	state         appState
	hooks         hookStore
	logger        Logger
	forceShutdown bool
	shutdownTime  time.Duration
}

func (lfcm *lifeCycleManagerImpl) getState() appState {
	return lfcm.state
}

func (lfcm *lifeCycleManagerImpl) status(newStatus appState) HookFn {
	return func() error {
		lfcm.logger.Info(fmt.Sprintf("Application %s", newStatus))
		lfcm.state = newStatus
		return nil
	}
}

func (lfcm *lifeCycleManagerImpl) transition(lifeCycle hook) HookFn {
	return func() error {
		lfcm.logger.Info(fmt.Sprintf("Processing on%s ", lifeCycle))
		if err := hooksChain(lfcm.hooks.Get(lifeCycle)...); err != nil {
			return fmt.Errorf("processing on%s: %w", lifeCycle, err)
		}
		return nil
	}
}

func (lfcm *lifeCycleManagerImpl) start() error {
	if lfcm.state != appStates.IDLE {
		lfcm.logger.Warn("The application has already started.")
		return errors.New("the application has already started")
	}

	err := hooksChain(
		lfcm.status(appStates.STARTING),
		lfcm.transition(hooks.BOOTING),
		lfcm.transition(hooks.BOOTED),
		lfcm.transition(hooks.READY),
		lfcm.transition(hooks.RUNNING),
		lfcm.status(appStates.STARTED),
	)
	if err != nil {
		if stopErr := lfcm.stop(); stopErr != nil {
			return fmt.Errorf("shutting down the application: %w", stopErr)
		}

		return fmt.Errorf("transitioning the life cycles: %w", err)
	}

	return nil
}

func (lfcm *lifeCycleManagerImpl) stop() error {
	if lfcm.state == appStates.IDLE {
		lfcm.logger.Warn("The application is not running.")
		return errors.New("the application is not running")
	}

	err := hooksChain(
		lfcm.status(appStates.STOPPING),
		lfcm.transition(hooks.DISPOSING),
		lfcm.transition(hooks.DISPOSED),
		lfcm.status(appStates.STOPPED),
	)
	if err != nil {
		return fmt.Errorf("shutting down the application: %w", err)
	}

	lfcm.logger.Info("Application exited with success")
	return nil
}

func (lfcm *lifeCycleManagerImpl) OnBooting(order HookOrder, fn ...HookFn) {
	lfcm.on(hooks.BOOTING, order, fn...)
}

func (lfcm *lifeCycleManagerImpl) OnBooted(order HookOrder, fn ...HookFn) {
	lfcm.on(hooks.BOOTED, order, fn...)
}

func (lfcm *lifeCycleManagerImpl) OnReady(order HookOrder, fn ...HookFn) {
	lfcm.on(hooks.READY, order, fn...)
}

func (lfcm *lifeCycleManagerImpl) OnRunning(order HookOrder, fn ...HookFn) {
	lfcm.on(hooks.RUNNING, order, fn...)
}

func (lfcm *lifeCycleManagerImpl) OnDisposing(order HookOrder, fn ...HookFn) {
	lfcm.on(hooks.DISPOSING, order, fn...)
}

func (lfcm *lifeCycleManagerImpl) OnDisposed(order HookOrder, fn ...HookFn) {
	lfcm.on(hooks.DISPOSED, order, fn...)
}

func (lfcm *lifeCycleManagerImpl) on(lfc hook, order HookOrder, fn ...HookFn) {
	if order == HookOrders.APPEND {
		lfcm.hooks.Append(lfc, fn...)
	} else {
		lfcm.hooks.Prepend(lfc, fn...)
	}
}

func (lfcm *lifeCycleManagerImpl) shutdown() {
	time.AfterFunc(lfcm.shutdownTime, func() {
		lfcm.logger.Warn("OK, my patience is over #ragequit")
		os.Exit(1)
	})

	if lfcm.state == appStates.STOPPING || lfcm.state == appStates.STOPPED {
		if lfcm.forceShutdown {
			lfcm.terminate(syscall.SIGKILL)
			return
		}

		lfcm.logger.Warn("The application is yet to finishing the shutdown process. Repeat the command to force exit")
		lfcm.forceShutdown = true
		return
	}

	if err := lfcm.stop(); err != nil {
		lfcm.logger.Error("Failed to stop the application", "err", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func (lfcm *lifeCycleManagerImpl) terminate(signal syscall.Signal) {
	// first arg is the process id
	arg0 := os.Args[0]
	val0, _ := strconv.ParseInt(arg0, 10, 32)
	pid := int(val0)

	err := syscall.Kill(pid, signal)

	if err != nil {
		lfcm.logger.Error("Failed to kill the process", "err", err)
		os.Exit(1)
	}
}

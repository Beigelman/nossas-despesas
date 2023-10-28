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

func newLifeCycleManager(shutdownTime time.Duration, logger Logger) *lifeCycleManager {
	return &lifeCycleManager{
		state:         appStates.IDLE,
		hooks:         newHookStore(),
		logger:        logger,
		forceShutdown: false,
		shutdownTime:  shutdownTime,
	}
}

type lifeCycleManager struct {
	state         appState
	hooks         *hookStore
	logger        Logger
	forceShutdown bool
	shutdownTime  time.Duration
}

func (lfcm *lifeCycleManager) getState() appState {
	return lfcm.state
}

func (lfcm *lifeCycleManager) status(newStatus appState) HookFn {
	return func() error {
		lfcm.logger.Info(fmt.Sprintf("[EON] Application %s", newStatus))
		lfcm.state = newStatus
		return nil
	}
}

func (lfcm *lifeCycleManager) transition(lifeCycle hook) HookFn {
	return func() error {
		lfcm.logger.Info(fmt.Sprintf("[EON] Processing on%s ", lifeCycle))
		if err := hooksChain(lfcm.hooks.get(lifeCycle)...); err != nil {
			return fmt.Errorf("processing on%s: %w", lifeCycle, err)
		}
		return nil
	}
}

func (lfcm *lifeCycleManager) start() error {
	if lfcm.state != appStates.IDLE {
		lfcm.logger.Warn("[EON] The application has already started.")
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

func (lfcm *lifeCycleManager) stop() error {
	if lfcm.state == appStates.IDLE {
		lfcm.logger.Warn("[EON] The application is not running.")
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

	lfcm.logger.Info("[EON] Application exited with success")
	return nil
}

func (lfcm *lifeCycleManager) OnBooting(order HookOrder, fn ...HookFn) {
	lfcm.on(hooks.BOOTING, order, fn...)
}

func (lfcm *lifeCycleManager) OnBooted(order HookOrder, fn ...HookFn) {
	lfcm.on(hooks.BOOTED, order, fn...)
}

func (lfcm *lifeCycleManager) OnReady(order HookOrder, fn ...HookFn) {
	lfcm.on(hooks.READY, order, fn...)
}

func (lfcm *lifeCycleManager) OnRunning(order HookOrder, fn ...HookFn) {
	lfcm.on(hooks.RUNNING, order, fn...)
}

func (lfcm *lifeCycleManager) OnDisposing(order HookOrder, fn ...HookFn) {
	lfcm.on(hooks.DISPOSING, order, fn...)
}

func (lfcm *lifeCycleManager) OnDisposed(order HookOrder, fn ...HookFn) {
	lfcm.on(hooks.DISPOSED, order, fn...)
}

func (lfcm *lifeCycleManager) on(lfc hook, order HookOrder, fn ...HookFn) {
	if order == HookOrders.APPEND {
		lfcm.hooks.append(lfc, fn...)
	} else {
		lfcm.hooks.prepend(lfc, fn...)
	}
}

func (lfcm *lifeCycleManager) shutdown() {
	time.AfterFunc(lfcm.shutdownTime, func() {
		lfcm.logger.Warn("[EON] OK, my patience is over #ragequit")
		os.Exit(1)
	})

	if lfcm.state == appStates.STOPPING || lfcm.state == appStates.STOPPED {
		if lfcm.forceShutdown {
			lfcm.terminate(syscall.SIGKILL)
			return
		}

		lfcm.logger.Warn("[EON] The application is yet to finishing the shutdown process. Repeat the command to force exit")
		lfcm.forceShutdown = true
		return
	}

	if err := lfcm.stop(); err != nil {
		lfcm.logger.Error("[EON] Failed to stop the application", "err", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func (lfcm *lifeCycleManager) terminate(signal syscall.Signal) {
	// first arg is the process id
	arg0 := os.Args[0]
	val0, _ := strconv.ParseInt(arg0, 10, 32)
	pid := int(val0)

	err := syscall.Kill(pid, signal)

	if err != nil {
		lfcm.logger.Error("[EON] Failed to kill the process", "err", err)
		os.Exit(1)
	}
}

package eon

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/truepay/go-commons/logger"
	"github.com/truepay/go-commons/logger/field"
	"golang.org/x/sys/unix"
)

// HookStore
type HookOrder string

var HookOrders = struct {
	APPEND  HookOrder
	PREPEND HookOrder
}{
	APPEND:  "APPEND",
	PREPEND: "PREPEND",
}

type HookFn func() error

type hookStore interface {
	Get(lfc hook) []HookFn
	Append(lfc hook, fn ...HookFn)
	Prepend(lfc hook, fn ...HookFn)
}

func newHookStore() hookStore {
	return &hookStoreImpl{
		hooks: map[hook][]HookFn{},
	}
}

type hookStoreImpl struct {
	hooks map[hook][]HookFn
}

func (store *hookStoreImpl) Get(lfc hook) []HookFn {
	return store.hooks[lfc]
}

func (store *hookStoreImpl) Append(lfc hook, fn ...HookFn) {
	store.hooks[lfc] = append(store.hooks[lfc], fn...)
}

func (store *hookStoreImpl) Prepend(lfc hook, fn ...HookFn) {
	store.hooks[lfc] = append(fn, store.hooks[lfc]...)
}

// Application
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

// LifeCycle
type hook string

var hooks = struct {
	BOOTING   hook
	BOOTED    hook
	READY     hook
	RUNNING   hook
	DISPOSING hook
	DISPOSED  hook
}{
	BOOTING:   "Booting",
	BOOTED:    "Booted",
	READY:     "Ready",
	RUNNING:   "Running",
	DISPOSING: "Disposing",
	DISPOSED:  "Disposed",
}

type lifeCycleManager interface {
	getState() appState
	start(ctx context.Context) error
	stop(ctx context.Context) error
	OnBooting(order HookOrder, fn ...HookFn)
	OnBooted(order HookOrder, fn ...HookFn)
	OnReady(order HookOrder, fn ...HookFn)
	OnRunning(order HookOrder, fn ...HookFn)
	OnDisposing(order HookOrder, fn ...HookFn)
	OnDisposed(order HookOrder, fn ...HookFn)
	shutdown(ctx context.Context)
}

func newLifeCycleManager(shutdownTime time.Duration, logger *logger.Logger, ctx context.Context) lifeCycleManager {
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
	logger        *logger.Logger
	forceShutdown bool
	shutdownTime  time.Duration
}

func (lfcm *lifeCycleManagerImpl) getState() appState {
	return lfcm.state
}

func (lfcm *lifeCycleManagerImpl) status(ctx context.Context, newStatus appState) HookFn {
	return func() error {
		lfcm.logger.Info(ctx, fmt.Sprintf("Application %s", newStatus))
		lfcm.state = newStatus
		return nil
	}
}

func (lfcm *lifeCycleManagerImpl) transition(ctx context.Context, lifeCycle hook) HookFn {
	return func() error {
		lfcm.logger.Info(ctx, fmt.Sprintf("Processing on%s ", lifeCycle))
		return hooksChain(lfcm.hooks.Get(lifeCycle)...)
	}
}

func (lfcm *lifeCycleManagerImpl) start(ctx context.Context) error {
	if lfcm.state != appStates.IDLE {
		lfcm.logger.Warn(ctx, "The application has already started.")
		return errors.New("the application has already started")
	}

	err := hooksChain(
		lfcm.status(ctx, appStates.STARTING),
		lfcm.transition(ctx, hooks.BOOTING),
		lfcm.transition(ctx, hooks.BOOTED),
		lfcm.transition(ctx, hooks.READY),
		lfcm.transition(ctx, hooks.RUNNING),
		lfcm.status(ctx, appStates.STARTED),
	)
	if err != nil {
		lfcm.logger.Error(ctx, "Error transitioning the life cycles", field.Error(err))

		if stopErr := lfcm.stop(ctx); stopErr != nil {
			lfcm.logger.Error(ctx, "Error stopping the application on start failure", field.Error(err))
			return stopErr
		}

		return err
	}

	return nil
}

func (lfcm *lifeCycleManagerImpl) stop(ctx context.Context) error {
	if lfcm.state == appStates.IDLE {
		lfcm.logger.Warn(ctx, "The application is not running.")
		return errors.New("the application is not running")
	}

	err := hooksChain(
		lfcm.status(ctx, appStates.STOPPING),
		lfcm.transition(ctx, hooks.DISPOSING),
		lfcm.transition(ctx, hooks.DISPOSED),
		lfcm.status(ctx, appStates.STOPPED),
	)
	if err != nil {
		lfcm.logger.Error(ctx, "Error shutting down the application", field.Error(err))
		return err
	}

	lfcm.logger.Info(ctx, "Application exited with success")

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

func (lfcm *lifeCycleManagerImpl) shutdown(ctx context.Context) {
	time.AfterFunc(lfcm.shutdownTime, func() {
		lfcm.logger.Warn(ctx, "OK, my patience is over #ragequit")
		if err := lfcm.logger.Flush(); err != nil {
			lfcm.logger.Error(ctx, "Error flushing the logger", field.Error(err))
		}

		os.Exit(1)
	})

	if lfcm.state == appStates.STOPPING || lfcm.state == appStates.STOPPED {
		if lfcm.forceShutdown {
			lfcm.terminate(ctx, syscall.SIGKILL)
		}

		lfcm.logger.Warn(ctx, "The application is yet to finishing the shutdown process. Repeat the command to force exit")
		lfcm.forceShutdown = true
		return
	}

	if err := lfcm.stop(ctx); err != nil {
		lfcm.logger.Fatal(ctx, "Failed to stop the application", field.Error(err))
	}

	if err := lfcm.logger.Flush(); !shouldIgnoreLoggerSyncError(err) {
		lfcm.logger.Error(ctx, "Error flushing the logger", field.Error(err))
	}

	os.Exit(0)
}

func (lfcm *lifeCycleManagerImpl) terminate(ctx context.Context, signal syscall.Signal) {
	// first arg is the process id
	arg0 := os.Args[0]
	val0, _ := strconv.ParseInt(arg0, 10, 32)
	pid := int(val0)

	err := syscall.Kill(pid, signal)

	if err != nil {
		lfcm.logger.Fatal(ctx, "Failed to kill the process", field.Error(err))
	}
}

// Workaround to ignore specific errors for unix os.
// The error only happens when stdout/stderr point to the console, but not if they are redirected to a file (since files support sync).
// Issue: https://github.com/uber-go/zap/issues/880
func shouldIgnoreLoggerSyncError(err error) bool {
	if err == nil {
		return true
	}

	errorsToIgnore := []syscall.Errno{
		unix.ENOTTY,
		unix.EINVAL,
	}

	for _, errorToIgnore := range errorsToIgnore {
		if errors.Is(err, errorToIgnore) {
			return true
		}
	}

	return false
}

// Utils
func hooksChain(hooks ...HookFn) error {
	for _, fn := range hooks {
		err := fn()
		if err != nil {
			return err
		}
	}
	return nil
}

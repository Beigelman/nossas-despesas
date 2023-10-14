package eon

import (
	"context"

	"github.com/Beigelman/ludaapi/internal/pkg/di"
	"github.com/Beigelman/ludaapi/internal/pkg/env"
)

// Both the boot and stopping processes are defined as a sequence of lifecycle events.
// These events exist in order to make these processes explicit and allow the modules to hook into them and properly integrate them into the application execution.
type LifeCycleManager interface {
	// When this event is dispatched it's a message to let the modules know that all the modules were already invoked and had already hooked into the other lifecycle events.
	// This is a good place to register error handlers because every module has already registered its routes when they were invoked.
	// Use this event to do anything you might need after all the module constructors are done running.
	OnBooted(order HookOrder, fn ...HookFn)
	// This lifecycle event happens after all the listeners for the booted event were run.
	// This is the proper place to actually start things, like starting the server or make queue consumers start listening to messages.
	OnReady(order HookOrder, fn ...HookFn)
	// After everything from the ready event is done, the app is now actually running.
	// A good usage for this lifecycle event is to know if the app is already prepared to be accessed during the setup of an automated test or offer info about the process in the console.
	OnRunning(order HookOrder, fn ...HookFn)
	// It's during this lifecycle event that the cleanup functions returned by the modules will be run.
	// To make the cleanup process consistent, the cleanup functions are run in the inverse order their modules were passed to the bootstrap function. So if your app uses `bootstrap(database, server)`, during the disposing process the cleanup function of the server module will be called first and then the database one.
	// As an example, this is where the server is stopped and the database connections are closed.
	// It's intended to be used to revert everything initialized during Booting lifecycle event.
	OnDisposing(order HookOrder, fn ...HookFn)
	// By the time Disposed event is dispatched, we expect that everything that keeps the process open is already finished, leaving it in a safe state to be terminated.
	// You could use this event to clean temporary files, for instance.
	OnDisposed(order HookOrder, fn ...HookFn)
}

type Info struct {
	ServiceName string
	Env         env.Environment
}

type BootFn func(ctx context.Context, c *di.Container, lc LifeCycleManager, info Info)

// An Eon module is a encapsulated piece of code that is responsible for booting a specific part of the application.
// It's expected to be used to separate the application in different parts, like the server, database, different bounded contexts, etc.
// The function that is passed to the module is expected to receive the dependency injection container, the lifecycle manager and the application info.
type Module struct {
	name   string
	bootFn BootFn
}

// NewModule creates an Eon module.
func NewModule(name string, bootFn BootFn) Module {
	return Module{
		name:   name,
		bootFn: bootFn,
	}
}

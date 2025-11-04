package eon

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHookStore(t *testing.T) {
	store := newHookStore()
	assert.NotNil(t, store)
	assert.NotNil(t, store.hooks)
	assert.Equal(t, 0, len(store.hooks))
}

func TestHookStoreAppend(t *testing.T) {
	store := newHookStore()

	dummyFn := func() error {
		return errors.New("an error")
	}

	store.append(hooks.BOOTING, dummyFn)

	hookFns := store.get(hooks.BOOTING)
	assert.Equal(t, 1, len(hookFns))
	assert.Equal(t, dummyFn().Error(), hookFns[0]().Error())
}

func TestHookStorePrepend(t *testing.T) {
	store := newHookStore()

	dummyFn1 := func() error {
		return errors.New("an error 1")
	}

	dummyFn2 := func() error {
		return errors.New("an error 2")
	}

	store.append(hooks.BOOTED, dummyFn1)
	store.prepend(hooks.BOOTED, dummyFn2)

	hookFns := store.get(hooks.BOOTED)
	assert.Equal(t, 2, len(hookFns))
	assert.Equal(t, dummyFn2().Error(), hookFns[0]().Error())
	assert.Equal(t, dummyFn1().Error(), hookFns[1]().Error())
}

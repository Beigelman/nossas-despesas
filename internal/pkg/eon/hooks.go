package eon

// HookOrder HookStore
type HookOrder string

var HookOrders = struct {
	APPEND  HookOrder
	PREPEND HookOrder
}{
	APPEND:  "APPEND",
	PREPEND: "PREPEND",
}

type HookFn func() error

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

type hookStore struct {
	hooks map[hook][]HookFn
}

func newHookStore() *hookStore {
	return &hookStore{
		hooks: map[hook][]HookFn{},
	}
}

func (store *hookStore) get(lfc hook) []HookFn {
	return store.hooks[lfc]
}

func (store *hookStore) append(lfc hook, fn ...HookFn) {
	store.hooks[lfc] = append(store.hooks[lfc], fn...)
}

func (store *hookStore) prepend(lfc hook, fn ...HookFn) {
	store.hooks[lfc] = append(fn, store.hooks[lfc]...)
}

func hooksChain(hooks ...HookFn) error {
	for _, fn := range hooks {
		err := fn()
		if err != nil {
			return err
		}
	}
	return nil
}

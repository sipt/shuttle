package global

import (
	"context"
	"sync"
)

func init() {
	namespace = make(map[string]*Namespace)
}

const defaultName = "default"

var namespace map[string]*Namespace
var mutex = &sync.RWMutex{}

func AddNamespace(name string, ctx context.Context, profile *Profile) {
	mutex.Lock()
	defer mutex.Unlock()
	n := &Namespace{
		profile: profile,
	}
	n.ctx, n.cancel = context.WithCancel(ctx)
	namespace[name] = n
}

func RemoveNamespace(name string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(namespace, name)
}

type Namespace struct {
	ctx     context.Context
	cancel  context.CancelFunc
	profile *Profile
}

func (n *Namespace) Profile() *Profile {
	return n.profile
}

func (n *Namespace) Cancel() {
	n.cancel()
}

func (n *Namespace) Context() context.Context {
	return n.ctx
}

func NamespaceWithContext(ctx context.Context) *Namespace {
	mutex.RLock()
	defer mutex.RUnlock()
	name, ok := ctx.Value("namespace").(string)
	if !ok || len(name) == 0 {
		return namespace[defaultName]
	}
	return namespace[name]
}

func NamespaceWithName(name ...string) *Namespace {
	mutex.RLock()
	defer mutex.RUnlock()
	if len(name) == 0 || len(name[0]) == 0 {
		return namespace[defaultName]
	}
	return namespace[name[0]]
}

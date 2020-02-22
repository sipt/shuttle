package namespace

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/sipt/shuttle/constant/typ"

	"github.com/sipt/shuttle/constant"
	"github.com/sipt/shuttle/global"
)

func init() {
	namespace = make(map[string]*Namespace)
}

const defaultName = "default"

var namespace map[string]*Namespace
var mutex = &sync.RWMutex{}

func AddNamespace(name string, ctx context.Context, profile *global.Profile, runtime typ.Runtime) {
	mutex.Lock()
	defer mutex.Unlock()
	n := &Namespace{
		profile: profile,
		mode:    constant.ModeRule,
		runtime: runtime,
	}
	needSave := true
	if mode, ok := runtime.Get("mode").(string); ok && len(mode) > 0 {
		switch mode {
		case constant.ModeRule, constant.ModeDirect, constant.ModeGlobal:
			n.mode = mode
			needSave = false
		default:
		}
	}
	if needSave {
		err := runtime.Set("mode", n.mode)
		if err != nil {
			logrus.WithError(err).Error("set namespace mode failed")
		}
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
	profile *global.Profile
	runtime typ.Runtime
	mode    string
}

func (n *Namespace) Profile() *global.Profile {
	return n.profile
}

func (n *Namespace) Cancel() {
	n.cancel()
}

func (n *Namespace) Context() context.Context {
	return n.ctx
}

func (n *Namespace) Mode() string {
	return n.mode
}

func (n *Namespace) SetMode(mode string) {
	switch mode {
	case constant.ModeRule, constant.ModeDirect, constant.ModeGlobal:
		n.mode = mode
		err := n.runtime.Set("mode", n.mode)
		if err != nil {
			logrus.WithError(err).Error("set mode failed")
		}
	default:
	}
}

func (n *Namespace) Runtime() typ.Runtime {
	return n.runtime
}

func NamespaceWithContext(ctx context.Context) *Namespace {
	mutex.RLock()
	defer mutex.RUnlock()
	name, ok := ctx.Value(constant.KeyNamespace).(string)
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

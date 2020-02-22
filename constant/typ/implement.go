package typ

func NewRuntime(name string, rt Runtime) Runtime {
	r := &runtime{
		name:    name,
		Runtime: rt,
	}
	var ok bool
	r.current, ok = rt.Get(name).(map[string]interface{})
	if !ok {
		r.current = make(map[string]interface{})
	}
	return r
}

type runtime struct {
	name    string
	current map[string]interface{}
	Runtime
}

func (r *runtime) Get(key string) interface{} {
	return r.current[key]
}
func (r *runtime) Set(key string, value interface{}) error {
	r.current[key] = value
	return r.Runtime.Set(r.name, r.current)
}

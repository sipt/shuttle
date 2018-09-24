package shuttle

func Recover(fs ...func()) {
	if err := recover(); err != nil {
		log.Logger.Errorf("[PANIC] %v", err)
		for _, f := range fs {
			f()
		}
	}
}

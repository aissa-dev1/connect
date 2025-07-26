package hasher

var global *Hasher

func SetGlobalHash(h *Hasher) {
	global = h
}

func Global() *Hasher {
	if global == nil {
		panic("Global hasher not initialized")
	}

	return global
}

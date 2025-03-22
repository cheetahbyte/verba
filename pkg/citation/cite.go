package citation

var usedKeys []string

func Register(key string) {
	for _, k := range usedKeys {
		if k == key {
			return
		}
	}
	usedKeys = append(usedKeys, key)
}

func UsedKeys() []string {
	return usedKeys
}

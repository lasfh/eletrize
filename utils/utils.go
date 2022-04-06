package utils

func Contains(list []string, name string) bool {
	for _, item := range list {
		if item == name {
			return true
		}
	}

	return false
}

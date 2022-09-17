package utils

func Contains[T comparable](list []T, name T) bool {
	for _, item := range list {
		if item == name {
			return true
		}
		
	}

	return false
}


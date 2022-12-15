package util

// returns if all elements of a given slice match a given condition.
func All[T any](slice []T, condition func(e T) bool) bool {
	for _, element := range slice {
		if !condition(element) {
			return false
		}
	}
	return true
}

// returns if any elements of a given slice match a given condition.
func Any[T any](slice []T, condition func(e T) bool) bool {
	for _, element := range slice {
		if condition(element) {
			return true
		}
	}
	return false
}

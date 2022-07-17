package utils

// IsInSlice check whether the slice contains an element based on the predicate
func IsInSlice[T any](heystack []T, pred func(T) bool) bool {
	return FindInSlice(heystack, pred) != nil
}

// FindInSlice get an element from the slice based on the predicate
func FindInSlice[T any](heystack []T, pred func(T) bool) *T {
	for _, item := range heystack {
		if pred(item) {
			return &item
		}
	}

	return nil
}

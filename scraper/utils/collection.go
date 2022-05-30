package utils

func IsInSlice[T any](heystack []T, pred func(T) bool) bool {
	return FindInSlice(heystack, pred) != nil
}

func FindInSlice[T any](heystack []T, pred func(T) bool) *T {
	for _, item := range heystack {
		if pred(item) {
			return &item
		}
	}

	return nil
}

package collut

// PtrGetDefault if the value is nil, return default value
func PtrGetDefault[T any](opt *T, defaultVal T) T {
	if opt == nil {
		return defaultVal
	}
	return *opt
}

// Ptr convert the value to the pointer
func Ptr[T any](value T) *T {
	return &value
}

// PtrNilType converts the type provided to the nil pointer with that associated type
func PtrNilType[T any]() *T {
	return (*T)(nil)
}

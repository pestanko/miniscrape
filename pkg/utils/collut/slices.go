package collut

import "golang.org/x/exp/slices"

// SliceMap takes a slice of items and a function and applies
// the function to each item and takes a result
// It takes two type parameters. S is the type of the input slice,
// and R is the type of the result slice. Both of them can be of any type.
func SliceMap[S any, R any](slice []S, mapper func(item S) R) []R {
	result := make([]R, len(slice))
	for idx, item := range slice {
		result[idx] = mapper(item)
	}
	return result
}

// SliceFilter takes a slice of items and a predicate to filter the elements
// S is the type parameter name that represents the type of the input slice
func SliceFilter[S any](slice []S, predicate func(item S) bool) (result []S) {
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return
}

// SliceContains takes a slice and a predicate function that will be called
// on the elements of the slice, if the predicate function returns true,
// the slice contains at least one element that matches the predicate
// T is an element type
// NOTE: This function is more readable than directly call the IndexFunc
func SliceContains[T any](slice []T, pred func(elem T) bool) bool {
	return slices.IndexFunc(slice, pred) != -1
}

// SliceFoldl left fold function for the slice
// takes an initial value, slice and fold function
// the fold function takes accumulator as a first parameter and the item from slice
// as the second
// T is the type of the elements of the slice, R is a result type
func SliceFoldl[T any, R any](init R, slice []T, fold func(acc R, next T) R) R {
	for _, value := range slice {
		init = fold(init, value)
	}

	return init
}

// SliceFoldr right fold function for the slice
// takes an initial value, slice and fold function
// the fold function takes accumulator as a first parameter and the item from slice
// as the second
// T is the type of the elements of the slice, R is a result type
func SliceFoldr[T any, R any](init R, list []T, fold func(acc R, next T) R) R {
	for idx := len(list) - 1; idx >= 0; idx-- {
		init = fold(init, list[idx])
	}

	return init
}

// SliceLengthLimit limit a length of slice, if the `limit` is greater than the length
// the slice, the original slice will be returned, otherwise the length will be cut
func SliceLengthLimit[T any](slice []T, limit int) []T {
	if len(slice) < limit {
		return slice
	}

	return slice[:limit]
}

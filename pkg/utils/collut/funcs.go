package collut

// Identity return the identity
func Identity[I any](item I) I {
	return item
}

// FuncCompose function composition
// apply the first function and then apply the second
// F is an input parameter type for the first function
// S is an input parameter type for the second func. and result type of the first
// R is a result type of the second function and the result function in general
// The composed function will be mapping from the F -> R (type-wise)
func FuncCompose[F, S, R any](fst func(F) S, snd func(S) R) func(item F) R {
	return func(item F) R {
		return snd(fst(item))
	}
}

// Zero return "zero"/"empty" value for a generic type
func Zero[T any]() T {
	var zero T
	return zero
}

// OpsApplyAll takes an item and returns a new item copy with applied operations
// BE AWARE: it does not update the original item (it takes it as value, so it is not possible)
// If you want to modify the original object, use OpsApplyAllRef
func OpsApplyAll[T any](original T, ops ...func(*T)) T {
	OpsApplyAllRef(&original, ops...)
	return original
}

// OpsApplyAllRef apply all operations to the provided item
func OpsApplyAllRef[T any](item *T, ops ...func(*T)) {
	for _, op := range ops {
		op(item)
	}
}

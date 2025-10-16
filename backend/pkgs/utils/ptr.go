package utils

func Ptr[T any](v T) *T {
	return &v
}

// UnPtr returns the value of the pointer or the zero value of the type
// if the pointer is nil.
func UnPtr[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}

// PtrOrNil returns T if T is not it's zero value. If T is it's zero value nil is returned.
// This is helpful when converting an API that uses non-pointer types into a database representation
// that may rely on null values to support certain operations.
func PtrOrNil[T comparable](v T) *T {
	var zero T
	if v == zero {
		return nil
	}

	return &v
}

func PtrInt32[T ~int](v *T) *int32 {
	if v == nil {
		return nil
	}
	i := int32(*v)
	return &i
}

func PtrToSlice[T any](v *[]T) []T {
	if v == nil {
		return nil
	}

	return *v
}

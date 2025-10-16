package utils

import (
	"fmt"
	"slices"
)

func ToNotNil[T any](v []T) []T {
	if v == nil {
		return []T{}
	}
	return v
}

func ToStrs[T fmt.Stringer](ids []T) []string {
	return Map(ids, func(v T) string {
		return v.String()
	})
}

// Map applies a function to each element of a slice and returns a new slice with the results
// the function f should not modify the slice elements.
func Map[T, R any](s []T, f func(T) R) []R {
	result := make([]R, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

// Set returns a new slice with only the unique elements of the input slice.
// The input slice is not modified.
func Set[T comparable](s []T) []T {
	seen := make([]T, 0, len(s))
	result := make([]T, 0, 10) // chose 10 because we're mostly using this for smaller slices

	for _, v := range s {
		found := slices.Contains(seen, v)
		if !found {
			seen = append(seen, v)
			result = append(result, v)
		}
	}

	return result
}

// Strings converts a slice of strings-like types to a slice of string.
func Strings[T ~string](slice []T) []string {
	result := make([]string, len(slice))
	for i, s := range slice {
		result[i] = string(s)
	}
	return result
}

// StringsToEnums converts a slice of strings to a slice of enums, it assumes the enum provides some
// parse function that can be used. This works in conjunction with out enum generator so anything
// generated is compatible with this function.
func StringsToEnums[T ~string](slice []string, factory func(string) (T, error)) ([]T, error) {
	result := make([]T, len(slice))

	for i, s := range slice {
		enum, err := factory(s)
		if err != nil {
			return result, err // Early return on first error
		}
		result[i] = enum
	}

	return result, nil
}

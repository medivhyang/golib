package slices

// Contains check if a slice contains a target element, alias of ContainAll
func Contains[T comparable](slice []T, target ...T) bool {
	return ContainAll(slice, target)
}

// ContainAll check if a slice contains all of the target elements
func ContainAll[T comparable](slice []T, target []T) bool {
	for _, targetItem := range target {
		find := false
		for _, item := range slice {
			if item == targetItem {
				find = true
				break
			}
		}
		if !find {
			return false
		}
	}
	return true
}

// ContainAny check if a slice contains any of the target elements
func ContainAny[T comparable](slice []T, target ...T) bool {
	for _, item := range slice {
		for _, targetItem := range target {
			if item == targetItem {
				return true
			}
		}
	}
	return false
}

// Filter filter a slice to a new slice
func Filter[T any](slice []T, fn func(target T) bool) []T {
	if fn == nil {
		return nil
	}
	var result []T
	for _, item := range slice {
		if fn(item) {
			result = append(result, item)
		}
	}
	return result
}

// Count count the number of elements in a slice
func Count[T any](slice []T, fn func(target T) bool) int {
	if fn == nil {
		return -1
	}
	var count int
	for _, item := range slice {
		if fn(item) {
			count++
		}
	}
	return count
}

// Map map a slice to a new slice
func Map[S any, D any](slice []S, fn func(s S) D) []D {
	if len(slice) == 0 {
		return nil
	}
	if fn == nil {
		return nil
	}
	var result []D
	for _, item := range slice {
		result = append(result, fn(item))
	}
	return result
}

// Reduce reduce a slice to a single value
func Reduce[T any](slice []T, fn func(a, b T) T) T {
	var result T
	if len(slice) == 0 {
		return result
	}
	result = slice[0]
	if fn == nil {
		return result
	}
	for _, v := range slice[1:] {
		result = fn(result, v)
	}
	return result
}

package lygo_array

import (
	"math/rand"
	"sort"
)

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func Sort(array interface{}) {
	if a, b := array.([]string); b {
		sort.Strings(a)
	} else if a, b := array.([]float64); b {
		sort.Float64s(a)
	} else if a, b := array.([]int); b {
		sort.Ints(a)
	}
}

// Copy a slice and return new slice with same items
func Copy(array interface{}) interface{} {
	if a, b := array.([]interface{}); b {
		response := make([]interface{}, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]string); b {
		response := make([]string, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]byte); b {
		response := make([]byte, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]int); b {
		response := make([]int, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]int8); b {
		response := make([]int8, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]int16); b {
		response := make([]int16, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]int32); b {
		response := make([]int32, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]int64); b {
		response := make([]int64, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]uint); b {
		response := make([]uint, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]uint8); b {
		response := make([]uint8, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]uint16); b {
		response := make([]uint16, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]uint32); b {
		response := make([]uint32, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]uint64); b {
		response := make([]uint64, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]uintptr); b {
		response := make([]uintptr, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]float32); b {
		response := make([]float32, len(a))
		copy(response, a)
		return response
	} else if a, b := array.([]float64); b {
		response := make([]float64, len(a))
		copy(response, a)
		return response
	}
	return nil
}

// Group a slice in batch.
// Returns a slice of slice.
// usage:
// response := Group([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
// fmt.Println(response) // [[0 1 2] [3 4 5] [6 7 8] [9]]
func Group(groupSize int, array interface{}) interface{} {
	if a, b := array.([]interface{}); b {
		groups := make([][]interface{}, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]string); b {
		groups := make([][]string, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]byte); b {
		groups := make([][]byte, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]int); b {
		groups := make([][]int, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]int8); b {
		groups := make([][]int8, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]int16); b {
		groups := make([][]int16, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]int32); b {
		groups := make([][]int32, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]int64); b {
		groups := make([][]int64, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]uint); b {
		groups := make([][]uint, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]uint8); b {
		groups := make([][]uint8, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]uint16); b {
		groups := make([][]uint16, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]uint32); b {
		groups := make([][]uint32, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]uint64); b {
		groups := make([][]uint64, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]uintptr); b {
		groups := make([][]uintptr, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]float32); b {
		groups := make([][]float32, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	} else if a, b := array.([]float64); b {
		groups := make([][]float64, 0, (len(a)+groupSize-1)/groupSize)
		for groupSize < len(a) {
			a, groups = a[groupSize:], append(groups, a[0:groupSize:groupSize])
		}
		groups = append(groups, a)
		return groups
	}
	return nil
}

func Reverse(array interface{}) {
	if a, b := array.([]interface{}); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]string); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]byte); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]int); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]int8); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]int16); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]int32); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]int64); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]uint); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]uint8); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]uint16); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]uint32); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]uint64); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]uintptr); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]float32); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	} else if a, b := array.([]float64); b {
		for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
			a[left], a[right] = a[right], a[left]
		}
	}
}

func Sub(array interface{}, start, end int) interface{} {
	if start >= end {
		start = 0
	}
	if a, b := array.([]interface{}); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]interface{}, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]string); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]string, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]byte); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]byte, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]int); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]int, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]int8); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]int8, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]int16); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]int16, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]int32); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]int32, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]int64); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]int64, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]uint); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]uint, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]uint8); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]uint8, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]uint16); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]uint16, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]uint32); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]uint32, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]uint64); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]uint64, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]uintptr); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]uintptr, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]float32); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]float32, 0)
		result = append(result, a[start:end+1]...)
		return result
	} else if a, b := array.([]float64); b {
		if end > len(a) {
			end = len(a)
		}
		result := make([]float64, 0)
		result = append(result, a[start:end+1]...)
		return result
	}
	return nil
}

func Shuffle(array interface{}) {
	if a, b := array.([]interface{}); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]string); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]byte); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]int); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]int8); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]int16); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]int32); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]int64); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]uint); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]uint8); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]uint16); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]uint32); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]uint64); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]uintptr); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]float32); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	} else if a, b := array.([]float64); b {
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

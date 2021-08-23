package utils

import "time"

// - a < b
// + a > b
// 0 a == b
type Comparator func(a, b interface{}) int

func StringComparator(a, b interface{}) int {
	s1 := a.(string)
	s2 := b.(string)
	// 获取遍历长度
	min := len(s2)
	if len(s1) < len(s2) {
		min = len(s1)
	}

	diff := 0
	for i := 0; i < min && diff == 0; i++ {
		diff = int(s1[i]) - int(s2[i])
	}
	if diff == 0 {
		diff = len(s1) - len(s2)
	}
	if diff < 0 {
		return -1
	}
	if diff > 0 {
		return 1
	}
	return 0

}

func IntComparator(a, b interface{}) int {
	aAsserted := a.(int)
	bAsserted := b.(int)
	switch {
	case aAsserted > bAsserted:
		return 1
	case aAsserted < bAsserted:
		return -1
	default:
		return 0
	}
}

func Int8comparator(a, b interface{}) int {
	aasserted := a.(int8)
	basserted := b.(int8)
	switch {
	case aasserted > basserted:
		return 1
	case aasserted < basserted:
		return -1
	default:
		return 0
	}
}
func Int16comparator(a, b interface{}) int {
	aasserted := a.(int16)
	basserted := b.(int16)
	switch {
	case aasserted > basserted:
		return 1
	case aasserted < basserted:
		return -1
	default:
		return 0
	}
}

func Int32Comparator(a, b interface{}) int {
	aAsserted := a.(int32)
	bAsserted := b.(int32)
	switch {
	case aAsserted > bAsserted:
		return 1
	case aAsserted < bAsserted:
		return -1
	default:
		return 0
	}
}
func Int64Comparator(a, b interface{}) int {
	aAsserted := a.(int64)
	bAsserted := b.(int64)
	switch {
	case aAsserted > bAsserted:
		return 1
	case aAsserted < bAsserted:
		return -1
	default:
		return 0
	}
}

func UIntComparator(a, b interface{}) int {
	aAsserted := a.(uint)
	bAsserted := b.(uint)
	switch {
	case aAsserted > bAsserted:
		return 1
	case aAsserted < bAsserted:
		return -1
	default:
		return 0
	}
}

func UInt8comparator(a, b interface{}) int {
	aasserted := a.(uint8)
	basserted := b.(uint8)
	switch {
	case aasserted > basserted:
		return 1
	case aasserted < basserted:
		return -1
	default:
		return 0
	}
}
func UInt16comparator(a, b interface{}) int {
	aasserted := a.(uint16)
	basserted := b.(uint16)
	switch {
	case aasserted > basserted:
		return 1
	case aasserted < basserted:
		return -1
	default:
		return 0
	}
}

func UInt32Comparator(a, b interface{}) int {
	aAsserted := a.(uint32)
	bAsserted := b.(uint32)
	switch {
	case aAsserted > bAsserted:
		return 1
	case aAsserted < bAsserted:
		return -1
	default:
		return 0
	}
}

func UInt64Comparator(a, b interface{}) int {
	aAsserted := a.(uint64)
	bAsserted := b.(uint64)
	switch {
	case aAsserted > bAsserted:
		return 1
	case aAsserted < bAsserted:
		return -1
	default:
		return 0
	}
}

func Float32Comparator(a, b interface{}) int {
	aAsserted := a.(float32)
	bAsserted := b.(float32)
	switch {
	case aAsserted > bAsserted:
		return 1
	case aAsserted < bAsserted:
		return -1
	default:
		return 0
	}
}

func Float64Comparator(a, b interface{}) int {
	aAsserted := a.(float64)
	bAsserted := b.(float64)
	switch {
	case aAsserted > bAsserted:
		return 1
	case aAsserted < bAsserted:
		return -1
	default:
		return 0
	}
}

func ByteComparator(a, b interface{}) int {
	aAsserted := a.(byte)
	bAsserted := b.(byte)
	switch {
	case aAsserted > bAsserted:
		return 1
	case aAsserted < bAsserted:
		return -1
	default:
		return 0
	}
}

func RuneComparator(a, b interface{}) int {
	aAsserted := a.(rune)
	bAsserted := b.(rune)
	switch {
	case aAsserted > bAsserted:
		return 1
	case aAsserted < bAsserted:
		return -1
	default:
		return 0
	}
}

func TimeComparator(a, b interface{}) int {
	aAsserted := a.(time.Time)
	bAsserted := b.(time.Time)
	switch {
	case aAsserted.After(bAsserted):
		return 1
	case aAsserted.Before(bAsserted):
		return -1
	default:
		return 0
	}
}

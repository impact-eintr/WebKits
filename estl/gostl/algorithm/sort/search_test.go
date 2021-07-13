package sort

import (
	"github.com/liyue201/gostl/ds/deque"
	"github.com/liyue201/gostl/utils/comparator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSearch(t *testing.T) {
	a := deque.New()
	a.PushBack(1)
	a.PushBack(3)
	a.PushFront(5)
	a.PushFront(4)
	a.PushBack(7)
	a.PushFront(3)
	a.PushBack(15)

	Sort(a.Begin(), a.End())

	t.Logf("a: %v", a)
	assert.True(t, BinarySearch(a.Begin(), a.End(), 5, comparator.BuiltinTypeComparator))
	assert.False(t, BinarySearch(a.Begin(), a.End(), 10))

	iter := LowerBound(a.Begin(), a.End(), 3, comparator.BuiltinTypeComparator)
	assert.Equal(t, 3, iter.Value())
	assert.Equal(t, 3, iter.Clone().Next().Value())

	iter = LowerBound(a.Begin(), a.End(), 4)
	assert.Equal(t, 4, iter.Value())

	iter = UpperBound(a.Begin(), a.End(), 4, comparator.BuiltinTypeComparator)
	assert.Equal(t, 5, iter.Value())

	iter = UpperBound(a.Begin(), a.End(), 15)
	assert.True(t, iter.Equal(a.End()))
}

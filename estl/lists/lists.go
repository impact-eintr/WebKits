package lists

import (
	"github.com/impact-eintr/WebKits/estl/container"
	"github.com/impact-eintr/WebKits/estl/utils"
)

type List interface {
	Get(index int) (interface{}, bool)
	Remove(index int)
	Add(values ...interface{})
	Contains(values ...interface{})
	Sort(comparator utils.Comparator)
	Swap(idx1, idx2 int)
	Insert(index int, values ...interface{})
	Set(index int, value interface{})

	container.Container
	// Empty() bool
	// Size() int
	// Clear()
	// Values() []interface{}
}

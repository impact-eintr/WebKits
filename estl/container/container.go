package container

import "github.com/impact-eintr/WebKits/estl/utils"

type Container interface {
	Empty() bool
	Size() int
	Clear()
	Values() []interface{}
}

func GetSortedValues(container Container, comparator utils.Comparator) []interface{} {
	values := container.Values()
	if len(values) < 2 {
		return values
	}
	utils.Sort(values, comparator)
	return values
}

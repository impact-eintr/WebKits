package arraylist

import "github.com/impact-eintr/WebKits/estl/utils"

type List struct {
	elements []interface{}
	size     int
}

const (
	growthFactor = float32(2.0)
	shrinkFactor = float32(0.25) // shrink 收缩
)

func New(values ...interface{}) *List {
	list := &List{}
	if len(values) > 0 {
		list.Add(values...)
	}
	return list
}

func (list *List) Get(index int) (interface{}, bool) {

}

func (list *List) Remove(index int) {

}

func (list *List) Add(values ...interface{}) {

}

func (list *List) Contains(values ...interface{}) {

}

func (list *List) Sort(comparator utils.Comparator) {

}

func (list *List) Swap(idx1, idx2 int) {

}

func (list *List) Insert(index int, values ...interface{}) {

}

func (list *List) Set(index int, value interface{}) {

}

func (list *List) Empty() bool {

}

func (list *List) Size() int {

}

func (list *List) Clear() {

}

func (list *List) Values() []interface{} {

}

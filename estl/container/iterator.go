package container

// 正反 根据索引/键值为联系的迭代器

type IteratorWithIndex interface {
	Next() bool
	Value() interface{}
	Index() int
	Begin()
	First() bool
}

type IteratorWithKey interface {
	Next() bool
	Value() interface{}
	Key() interface{}
	Begin()
	First() bool
}

type ReverseInteratorWithIndex interface {
	Prev() bool
	End()
	Last() bool
	IteratorWithIndex
}

type ReverseInterratorWithKey interface {
	Prev() bool
	End()
	Last() bool
	IteratorWithKey
}

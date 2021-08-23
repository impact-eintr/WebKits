package container

type EnumberableWithIndex interface {
	Each(func(index int, value interface{}))
	Any(func(index int, value interface{}) bool) bool
	All(func(index int, value interface{}) bool) bool
	Find(func(index int, value interface{}) bool) (int, interface{})
}

type EnumberableWithKey interface {
	Each(func(key interface{}, value interface{}))
	Any(func(key interface{}, value interface{}) bool) bool
	All(func(key interface{}, value interface{}) bool) bool
	Find(func(key interface{}, value interface{}) bool) (int, interface{})
}

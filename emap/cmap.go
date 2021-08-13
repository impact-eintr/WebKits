package emap

// ConcurrentMap 代表并发安全的字典的接口
type ConcurrentMap interface {
	// Concurrency() 返回并发量
	Concurrency()
	// Put 会推送一个 k-v
	Put(key string, value interface{}) (bool, error)
	// Get 可以获取与指定键向关联的那个元素
	Get(key string) interface{}
	// Delete 会删除 k-v
	Delete(key string) bool
	// Len 会返回k-v的数量
	Len()
}

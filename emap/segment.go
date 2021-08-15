package emap

import "sync"

// 并发安全字典的实现类型 [散列段接口]
type Segment interface {
	Put(p Pair) (bool, error)
	Get(key string) Pair
	GetWithHash(key string, keyHash uint64)
}

type segment struct {
	// buckets 代表散列桶切片
	buckets []Bucket
	// buckets 代表散列桶切片的长度
	bucketsLen int
	// pairTotal 代表k-v 总数
	paitTotal uint64
	// pairRedistributor k-v 再分布器
	pairRedistributor PairRedistributor
	lock              sync.Mutex
}

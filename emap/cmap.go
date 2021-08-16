package emap

import (
	"math"
	"sync/atomic"
)

// ConcurrentMap 代表并发安全的字典的接口
type ConcurrentMap interface {
	// Concurrency() 返回并发量
	Concurrency() int
	// Put 会推送一个 k-v
	Put(key string, value interface{}) (bool, error)
	// Get 可以获取与指定键向关联的那个元素
	Get(key string) interface{}
	// Delete 会删除 k-v
	Delete(key string) bool
	// Len 会返回k-v的数量
	Len() uint64
}

type myConcurrentMap struct {
	concurrency int
	segments    []Segment
	total       uint64
}

func NewConcurrentMap(concurrency int, prd PairRedistributor) (ConcurrentMap, error) {
	if concurrency <= 0 {
		return nil, newIllegalParameterError("concurrency is too small")
	}

	if concurrency > MAX_CONCURRENCY {
		return nil, newIllegalParameterError("concurrency is too large")
	}

	cmap := &myConcurrentMap{}
	cmap.concurrency = concurrency
	cmap.segments = make([]Segment, concurrency)
	for i := 0; i < concurrency; i++ {
		cmap.segments[i] = newSegment(DEFAULT_BUCKET_NUMBER, prd)
	}
	return cmap, nil

}

func (c *myConcurrentMap) Concurrency() int {
	return c.concurrency
}

func (c *myConcurrentMap) Put(key string, value interface{}) (bool, error) {
	p, err := newPair(key, value)
	if err != nil {
		return false, err
	}
	s := c.findSegment(p.Hash()) // 先找到应该放在哪个段内
	ok, err := s.Put(p)
	if ok {
		atomic.AddUint64(&c.total, 1)
	}
	return ok, err
}

func (c *myConcurrentMap) Get(key string) interface{} {
	keyHash := hash(key)
	s := c.findSegment(keyHash)
	pair := s.GetWithHash(key, keyHash)
	if pair == nil {
		return nil
	}
	return pair.Element

}

func (c *myConcurrentMap) Delete(key string) bool {
	s := c.findSegment(hash(key))
	if s.Delete(key) {
		atomic.AddUint64(&c.total, ^uint64(0))
		return true
	}
	return false
}

func (c *myConcurrentMap) Len() uint64 {
	return atomic.LoadUint64(&c.total)
}

// 会根据给定参数寻找并返回对应散列段
func (c *myConcurrentMap) findSegment(keyHash uint64) Segment {
	if c.concurrency == 1 {
		return c.segments[0] // 不并发就只有一个散列段
	}
	var keyHashHigh int // 键值hash高位
	if keyHash > math.MaxUint32 {
		keyHashHigh = int(keyHash >> 48)
	} else {
		keyHashHigh = int(keyHash >> 16)
	}
	return c.segments[keyHashHigh%c.concurrency]

}

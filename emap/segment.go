package emap

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Segment 代表并发安全的散列段的接口。
type Segment interface {
	// Put 会根据参数放入一个键-元素对。
	// 第一个返回值表示是否新增了键-元素对。
	Put(p Pair) (bool, error)
	// Get 会根据给定参数返回对应的键-元素对。
	// 该方法会根据给定的键计算哈希值。
	Get(key string) Pair
	// GetWithHash 会根据给定参数返回对应的键-元素对。
	// 注意！参数keyHash应该是基于参数key计算得出哈希值。
	GetWithHash(key string, keyHash uint64) Pair
	// Delete 会删除指定键的键-元素对。
	// 若返回值为true则说明已删除，否则说明未找到该键。
	Delete(key string) bool
	// Size 用于获取当前段的尺寸（其中包含的散列桶的数量）。
	Size() uint64
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

func newSegment(bucketNumber int, pairRedistributor PairRedistributor) Segment {
	if bucketNumber <= 0 {
		bucketNumber = DEFAULT_BUCKET_NUMBER
	}
	if pairRedistributor == nil {
		pairRedistributor = newDefaultPairRedistributor(DEFAULT_BUCKET_LOAD_FACTOR, bucketNumber)
	}
	// 初始化散列桶数组
	buckets := make([]Bucket, bucketNumber)
	for i := 0; i < bucketNumber; i++ {
		buckets[i] = newBucket()
	}
	return &segment{
		buckets:           buckets,
		bucketsLen:        bucketNumber,
		pairRedistributor: pairRedistributor,
	}

}

func (s *segment) Put(p Pair) (bool, error) {
	s.lock.Lock()
	b := s.buckets[int(p.Hash()%uint64(s.bucketsLen))]
	ok, err := b.Put(p, nil)
	if ok {
		newTotal := atomic.AddUint64(&s.paitTotal, 1)
		s.redistribute(newTotal, b.Size())
	}
	s.lock.Unlock()
	return ok, err
}

func (s *segment) Get(key string) Pair {
	return s.GetWithHash(key, hash(key))
}

func (s *segment) GetWithHash(key string, keyHash uint64) Pair {
	s.lock.Lock()
	b := s.buckets[int(keyHash%uint64(s.bucketsLen))]
	s.lock.Unlock()
	return b.Get(key)
}

func (s *segment) Delete(key string) bool {

}

func (s *segment) Size() uint64 {

}

// 会检查给定参数并设定相应的阈值和计数
// 并在必要时候重新分配所有散列桶中的所有键值对
func (s *segment) redistribute(pairTotal uint64, bucketSize uint64) (err error) {
	defer func() {
		if p := recover(); p != nil {
			if pErr, ok := p.(error); ok {
				err = newPairRedistributorError(pErr.Error())
			} else {
				err = newPairRedistributorError(fmt.Sprintf("%s", p))
			}
		}
	}()

	s.pairRedistributor.UpdateThreshold(pairTotal, s.bucketsLen)
	bucketStatus := s.pairRedistributor.CheckBucketStatus(pairTotal, bucketSize)
	newBuckets, changed := s.pairRedistributor.Redistribute(bucketStatus, s.buckets)
	if changed {
		s.buckets = newBuckets
		s.bucketsLen = len(s.buckets)
	}
	return nil

}

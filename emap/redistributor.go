package emap

import (
	"log"
	"sync/atomic"
)

type BucketStatus int

const (
	BUCKET_STATUS_NORMAL BucketStatus = iota
	BUCKET_STATUS_UNDERWEIGHT
	BUCKET_STATUS_OVERWEIGHT
)

// 代表针对键值对的再分布器
// 用于当散列段内的键值对分布不均时进行重新分布
type PairRedistributor interface {
	// 根据键值对总数和散列桶总数计算并更新阈值
	UpdateThreshold(pairTotal uint64, bucketNumber int)
	// 用于检查散列桶的状态
	CheckBucketStatus(pairTotal uint64, bucketSize uint64) (bucketStatus BucketStatus)
	// 用于实施键值对的再分布
	Redistribute(bucketStatus BucketStatus, buckets []Bucket) (newBuckets []Bucket, changed bool)
}

type myPairRedistributor struct {
	loadFactor            float64 // 装载因子
	upperThreshold        uint64  // 当触发散列桶重量阈值后会进行再散列
	overweightBucketCount uint64  // 统计过重的散列桶
	emptyBucketCount      uint64  // 统计空桶
}

// 创建一个PairRedistributor类型的实例
// 参数loadFactor散列桶的负载因子
// 参数bucketNumber散列桶的数量
func newDefaultPairRedistributor(loadFactor float64, bucketNumber int) PairRedistributor {
	if loadFactor <= 0 {
		loadFactor = DEFAULT_BUCKET_LOAD_FACTOR
	}
	pr := &myPairRedistributor{loadFactor: loadFactor}
	pr.UpdateThreshold(0, bucketNumber)
	return pr
}

// 调试试用散列桶状态信息模板
var bucketCountTemplate = `Bucket count:
  pairTotal: %d
  bucketNumber: %d
  average: %f
  upperThreshold: %d
  emptyBucketCount: %d
`

func (p *myPairRedistributor) UpdateThreshold(pairTotal uint64, bucketNumber int) {
	var average float64
	average = float64(pairTotal / uint64(bucketNumber))
	if average < 100 {
		average = 100
	}
	//defer func() {
	//	fmt.Printf(bucketCountTemplate, pairTotal, bucketNumber, average,
	//		atomic.LoadUint64(&p.upperThreshold),
	//		atomic.LoadUint64(&p.emptyBucketCount))
	//}()
	atomic.StoreUint64(&p.upperThreshold, uint64(average*p.loadFactor))

}

// 散列桶状态信息模板
var bucketStatusTemplate = `Bucket count:
  pairTotal: %d
  bucketSize: %d
  upperThreshold: %d
  overweightBucketCount: %d
  emptyBucketCount: %d
  bucketStatus: %d
`

func (p *myPairRedistributor) CheckBucketStatus(pairTotal uint64, bucketSize uint64) (bucketStatus BucketStatus) {
	//defer func() {
	//	fmt.Printf(bucketStatusTemplate,
	//		pairTotal,
	//		bucketSize,
	//		atomic.LoadUint64(&p.upperThreshold),
	//		atomic.LoadUint64(&p.overweightBucketCount),
	//		atomic.LoadUint64(&p.emptyBucketCount),
	//		bucketStatus)
	//}()
	if bucketSize > DEFAULT_BUCKET_MAX_SIZE ||
		bucketSize >= atomic.LoadUint64(&p.upperThreshold) {
		atomic.AddUint64(&p.overweightBucketCount, 1)
		bucketStatus = BUCKET_STATUS_OVERWEIGHT
		return
	}
	if bucketSize == 0 {
		atomic.AddUint64(&p.emptyBucketCount, 1)
	}
	return

}

var redistributionTemplate = `Bucket count:
  bucketStatus: %d
  currentNumber: %f
  newNumber: %d
`

func (p *myPairRedistributor) Redistribute(bucketStatus BucketStatus, buckets []Bucket) (newBuckets []Bucket, changed bool) {
	currentNumber := uint64(len(buckets))
	newNumber := currentNumber
	//defer func() {
	//	fmt.Printf(redistributionTemplate,
	//		bucketStatus, currentNumber, newNumber)
	//}()
	switch bucketStatus {
	case BUCKET_STATUS_OVERWEIGHT:
		if atomic.LoadUint64(&p.overweightBucketCount)*4 < currentNumber {
			return nil, false
		}
		newNumber = currentNumber << 1 // 桶数翻番
	case BUCKET_STATUS_UNDERWEIGHT:
		if currentNumber < 100 ||
			atomic.LoadUint64(&p.emptyBucketCount)*4 < currentNumber {
			return nil, false
		}
		newNumber = currentNumber >> 1 // 桶数砍半
		if newNumber < 2 {
			newNumber = 2
		}
	default:
		return nil, false
	}

	var pairs []Pair
	for _, b := range buckets {
		for e := b.GetFirstPair(); e != nil; e = e.Next() {
			pairs = append(pairs, e)
		}
	}
	log.Println(pairs)

	if newNumber == currentNumber { // 0 == 0
		atomic.StoreUint64(&p.overweightBucketCount, 0)
		atomic.StoreUint64(&p.emptyBucketCount, 0)
		return nil, false
	}
	if newNumber > currentNumber {
		for i := uint64(0); i < currentNumber; i++ {
			buckets[i].Clear(nil)
		}

		for j := newNumber - currentNumber; j > 0; j-- {
			buckets = append(buckets, newBucket())
		}
	} else {
		buckets = make([]Bucket, newNumber)
		for i := uint64(0); i < newNumber; i++ {
			buckets[i] = newBucket()
		}
	}

	var count int
	for _, p := range pairs {
		index := int(p.Hash() % newNumber)
		b := buckets[index]
		b.Put(p, nil)
		count++
	}
	// 扩容完毕
	atomic.StoreUint64(&p.overweightBucketCount, 0)
	atomic.StoreUint64(&p.emptyBucketCount, 0)
	return buckets, true

}

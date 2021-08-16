package emap

const (
	DEFAULT_BUCKET_LOAD_FACTOR float64 = 0.75 // 代表默认的装载因子
	DEFAULT_BUCKET_NUMBER      int     = 16   // 代表一个散列段包含的散列桶的默认数量
	DEFAULT_BUCKET_MAX_SIZE    uint64  = 1000
)

const (
	MAX_CONCURRENCY int = 1 << 16
)

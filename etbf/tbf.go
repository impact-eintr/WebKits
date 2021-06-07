package etbf

import (
	"sync"
)

type TBF interface {
	Fetchtoken()
	Returntoken()
	Destory()
}

type tbf_st struct {
	cps   int64
	burst int64
	token int64
	pos   int64
	mut   sync.RWMutex
	cond  sync.Cond
}

const (
	MYTBF_MAX int64 = 1024
)

var (
	job = make([]*tbf_st, MYTBF_MAX)
)

func get_free_pos_unlocked() int64 {
	for i := int64(0); i < MYTBF_MAX; i++ {
		if job[i] == nil {
			return i
		}
	}
	return -1
}

func Newtbf(cps, burst int64) *TBF {
	tbf := &tbf_st{
		cps:   cps,
		burst: burst,
		token: 0,
	}

	tbf.mut.Lock()
	// 将新的tbf装载到任务组中
	pos := get_free_pos_unlocked()
	if pos == -1 {
		tbf.mut.Unlock()
		return nil
	}

	tbf.pos = pos
	job[pos] = tbf
	tbf.mut.Unlock()

	return tbf
}

// 初始化
func init() {
	mod_load()
}

func mod_load() {
	go handler()
}

func mod_unload() {

}

// 每秒派发一次令牌
func handler() {
	for {

	}
}

func (t *tbf_st) Fectchtoken() {

}

func (t *tbf_st) Returntoken() {

}

func (t *tbf_st) Destory() {

}

package etbf

import (
	"errors"
	"sync"
	"time"
)

type TBF interface {
	Fetchtoken(int64) (int64, error)
	Returntoken(int64) (int64, error)
	Destory() error
}

type tbf_st struct {
	cps   int64
	burst int64
	token int64
	pos   int64
	mut   sync.Mutex
	cond  *sync.Cond
}

const (
	MYTBF_MAX int64 = 1024
)

var (
	job   = make([]*tbf_st, MYTBF_MAX)
	Mutex = sync.Mutex{}
)

// 初始化
func init() {
	mod_load()
}

func mod_load() {
	go handler()
}

func mod_unload() {
	for _, tbf := range job {
		if tbf != nil {
			tbf.Destory()
		}
	}
}

// 每秒派发一次令牌
func handler() {
	for {
		Mutex.Lock()
		for _, tbf := range job {
			if tbf != nil {
				tbf.mut.Lock()
				tbf.token += tbf.cps
				if tbf.token > tbf.burst {
					tbf.token = tbf.burst
				}
				tbf.cond.Broadcast()
				tbf.mut.Unlock()
			}
		}

		Mutex.Unlock()
		time.Sleep(time.Second)
	}
	return

}

func get_free_pos_unlocked() int64 {
	for i := int64(0); i < MYTBF_MAX; i++ {
		if job[i] == nil {
			return i
		}
	}
	return -1
}

func Newtbf(cps, burst int64) TBF {
	tbf := &tbf_st{
		cps:   cps,
		burst: burst,
		token: 0,
	}

	tbf.cond = sync.NewCond(&tbf.mut)

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

func (t *tbf_st) Fetchtoken(size int64) (int64, error) {
	if size <= 0 {
		return 0, errors.New("非法的参数")
	}

	t.mut.Lock()

	for t.token <= 0 {
		t.cond.Wait()
	}

	var n int64
	if t.token < size {
		n = t.token
	} else {
		n = size
	}
	t.token -= n

	t.mut.Unlock()

	return n, nil

}

func (t *tbf_st) Returntoken(size int64) (int64, error) {
	if size <= 0 {
		return 0, errors.New("非法的参数")
	}

	t.mut.Lock()

	t.token += size
	if t.token > t.burst {
		t.token = t.burst
	}

	t.mut.Unlock()

	return size, nil

}

func (t *tbf_st) Destory() error {
	t.mut.Lock()
	job[t.pos] = nil
	t.mut.Unlock()

	return nil
}

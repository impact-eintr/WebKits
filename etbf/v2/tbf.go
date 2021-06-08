package v2

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
	cps    int64
	burst  int64
	token  int64
	pos    int64
	mut    sync.Mutex
	cond   *sync.Cond
	ticker *time.Ticker
	Exitch chan struct{}
}

const (
	MYTBF_MAX int64 = 1024
)

var (
	job   = make([]*tbf_st, MYTBF_MAX)
	Mutex = sync.Mutex{}
)

func get_free_pos_unlocked() int64 {
	for i := int64(0); i < MYTBF_MAX; i++ {
		if job[i] == nil {
			return i
		}
	}
	return -1
}

func Newtbf(fillInterval time.Duration, cps, burst int64) TBF {
	tbf := &tbf_st{
		cps:    cps,
		burst:  burst,
		token:  0,
		ticker: time.NewTicker(fillInterval),
		Exitch: make(chan struct{}),
	}

	tbf.cond = sync.NewCond(&tbf.mut)

	Mutex.Lock()
	// 将新的tbf装载到任务组中
	pos := get_free_pos_unlocked()
	if pos == -1 {
		Mutex.Unlock()
		return nil
	}

	tbf.pos = pos
	job[pos] = tbf
	Mutex.Unlock()

	go func(s int64) {
		for {
			select {
			case <-job[s].ticker.C:
				job[s].mut.Lock()
				job[s].token += job[s].cps
				if job[s].token > job[s].burst {
					job[s].token = job[s].burst
				}
				job[s].cond.Broadcast()
				job[s].mut.Unlock()
			case <-job[s].Exitch:
				return
			}
		}
	}(pos)

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
	job[t.pos].Exitch <- struct{}{}
	t.mut.Lock()
	job[t.pos] = nil
	t.mut.Unlock()

	return nil
}

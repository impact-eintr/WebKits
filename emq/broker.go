package emq

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Broker interface {
	publish(topic string, msg interface{}) error
	subscribe(topic string) (<-chan interface{}, error)
	unsubscribe(topic string, sub <-chan interface{}) error
	close()
	broadcast(msg interface{}, subscribers []chan interface{})
	setConditions(cap int)
}

type StdBroker struct {
	exit chan struct{}
	cap  int

	topics map[string][]chan interface{}
	sync.RWMutex
}

func NewBroker() *StdBroker {
	return &StdBroker{
		exit:   make(chan struct{}),
		topics: make(map[string][]chan interface{}),
	}
}

// publish：进行消息的推送，有两个参数即topic、msg，分别是订阅的主题、要传递的消息
func (b *StdBroker) publish(topic string, msg interface{}) error {
	select {
	case <-b.exit:
		return errors.New("broker closed")
	default:
	}

	b.RLock()
	subscribers, ok := b.topics[topic]
	b.RUnlock()
	if !ok {
		return nil
	}

	b.broadcast(msg, subscribers)
	return nil
}

// subscribe：消息的订阅，传入订阅的主题，即可完成订阅，并返回对应的channel通道用来接收数据
func (b *StdBroker) subscribe(topic string) (<-chan interface{}, error) {
	select {
	case <-b.exit:
		return nil, errors.New("broker closed")
	default:
	}

	ch := make(chan interface{}, b.cap)
	b.Lock()
	b.topics[topic] = append(b.topics[topic], ch)
	b.Unlock()
	fmt.Println("订阅: ", topic)
	return ch, nil
}

// subscribe：取消订阅，传入订阅的主题和对应的通道
func (b *StdBroker) unsubscribe(topic string, sub <-chan interface{}) error {
	select {
	case <-b.exit:
		return errors.New("broker close")
	default:
	}

	b.RLock()
	subscribers, ok := b.topics[topic]
	b.RUnlock()

	if !ok {
		return nil
	}

	b.Lock()
	var newSubs []chan interface{}
	for _, subscriber := range subscribers {
		if subscriber == sub {
			continue
		}
		newSubs = append(newSubs, subscriber)
	}
	b.topics[topic] = newSubs
	b.Unlock()
	return nil

}

// close：关闭消息队列
func (b *StdBroker) close() {
	select {
	case <-b.exit:
		return
	default:
		close(b.exit)
		b.Lock()
		b.topics = make(map[string][]chan interface{})
		b.Unlock()
	}
	return
}

// broadCast：这个属于内部方法，作用是进行广播，对推送的消息进行广播，保证每一个订阅者都可以收到
func (b *StdBroker) broadcast(msg interface{}, subscribers []chan interface{}) {
	count := len(subscribers)
	concurrency := 1

	switch {
	case count > 1000:
		concurrency = 3
	case count > 100:
		concurrency = 2
	default:
		concurrency = 1
	}

	pub := func(start int) {
		idleDuration := 5 * time.Millisecond
		idleTimeout := time.NewTimer(idleDuration)
		defer idleTimeout.Stop()

		for j := start; j < count; j += concurrency {
			if !idleTimeout.Stop() {
				select {
				case <-idleTimeout.C:
				default:
				}
			}
			idleTimeout.Reset(idleDuration)
			select {
			case subscribers[j] <- msg:
			case <-idleTimeout.C:
			case <-b.exit:
				return
			}
		}
	}

	for i := 0; i < concurrency; i++ {
		go pub(i)
	}
}

// setConditions：这里是用来设置条件，条件就是消息队列的容量，这样我们就可以控制消息队列的大小了
func (b *StdBroker) setConditions(cap int) {
	b.cap = cap
}

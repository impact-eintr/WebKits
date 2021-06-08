package gotickets

import (
	"errors"
	"fmt"

	lg "github.com/impact-eintr/WebKits/loadgenerator"
)

type gotickets struct {
	total    uint32        // 票的总数
	ticketCh chan struct{} // 票的容器
	active   bool          // 票池是否已经被激活
}

func NewGoTickets(total uint32) (lg.GoTickets, error) {
	g := &gotickets{}
	if !g.init(total) {
		errMsg := fmt.Sprintf(
			"The goroutine ticket pool can NOT be initialized! (total=%d)\n",
			total)
		return nil, errors.New(errMsg)
	}

	return g, nil

}

func (g *gotickets) init(total uint32) bool {
	if g.active {
		return false
	}
	if total == 0 {
		return false
	}
	ch := make(chan struct{}, total)
	n := int(total)
	for i := 0; i < n; i++ {
		ch <- struct{}{}
	}
	g.ticketCh = ch
	g.total = total
	g.active = true
	return true

}

func (g *gotickets) Fetch() {
	<-g.ticketCh
}

func (g *gotickets) Return() {
	g.ticketCh <- struct{}{}
}

// 票池是否已经被激活
func (g *gotickets) Active() bool {
	return g.active
}

func (g *gotickets) Total() uint32 {
	return g.total
}

func (g *gotickets) Remainder() uint32 {
	return uint32(len(g.ticketCh))
}

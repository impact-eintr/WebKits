package emap

import (
	"bytes"
	"fmt"
	"sync/atomic"
	"unsafe"
)

type linkedPair interface {
	Next() Pair
	SetNext(nextPair Pair) error
}

type Pair interface {
	linkedPair
	Key() string
	Hash() uint64
	Element() interface{}
	SetElement(element interface{}) error
	Copy() Pair
	String() string
}

type pair struct {
	key     string
	hash    uint64 // 代表键的hash值
	element unsafe.Pointer
	next    unsafe.Pointer
}

func newPair(key string, element interface{}) (Pair, error) {
	p := &pair{
		key:  key,
		hash: hash(key),
	}
	if element == nil {
		return nil, newIllegalParameterError("element is nil")
	}

	p.element = unsafe.Pointer(&element)
	return p, nil

}

func (p *pair) Key() string {
	return p.key

}

func (p *pair) Hash() uint64 {
	return p.hash

}

func (p *pair) Element() interface{} {
	pointer := atomic.LoadPointer(&p.element)
	if pointer == nil {
		return nil
	}
	return *(*interface{})(pointer)

}

func (p *pair) SetElement(element interface{}) error {
	if element == nil {
		return newIllegalParameterError("element is nil")
	}
	atomic.StorePointer(&p.element, unsafe.Pointer(&element))
	return nil

}

func (p *pair) Copy() Pair {
	pCopy, _ := newPair(p.Key(), p.Element())
	return pCopy
}

func (p *pair) String() string {
	return p.genString(false)
}

func (p *pair) genString(nextDetail bool) string {
	var buf bytes.Buffer
	msg := fmt.Sprintf("pair{key: %s,hash: %d, element: %+v,", p.Key(), p.Hash(), p.Element())
	buf.WriteString(msg)
	if nextDetail {
		msg = "next: "
		if next := p.Next(); next != nil {
			if npp, ok := next.(*pair); ok {
				msg += npp.genString(nextDetail)
			} else {
				msg += "<ignore>"
			}
			buf.WriteString(msg)
		}
	} else {
		msg = "nextKey: "
		if next := p.Next(); next != nil {
			msg += next.Key()
		}
		buf.WriteString(msg)
	}
	buf.WriteString("}")
	return buf.String()

}

func (p *pair) Next() Pair {
	pointer := atomic.LoadPointer(&p.next)
	if pointer == nil {
		return nil
	}
	return (*pair)(pointer)

}

func (p *pair) SetNext(nextPair Pair) error {
	if nextPair == nil {
		atomic.StorePointer(&p.next, nil)
		return nil
	}
	pp, ok := nextPair.(*pair)
	if !ok {
		return newIllegalPairTypeError(nextPair)
	}
	atomic.StorePointer(&p.next, unsafe.Pointer(pp))
	return nil

}

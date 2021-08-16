package emap

import (
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

func newPair(key string, element interface{}) (Pairm, error) {
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

}

func (p *pair) Copy() Pair {

}

func (p *pair) String() string {
}

func (p *pair) Next() Pair {

}

func (p *pair) SetNext(nextPair Pair) error {

}

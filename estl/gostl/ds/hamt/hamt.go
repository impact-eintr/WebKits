package hamt

import (
	"github.com/liyue201/gostl/utils/sync"
	"github.com/liyue201/gostl/utils/visitor"
	"hash/fnv"
	"math/bits"
	gosync "sync"
)

// Some constants
const (
	BITMAP_NODE = 0
	KV_NODE     = 1
	Fanout      = 6 //each bitmap node has 6 bits, so the max depth of tree is 64/6 = 10.666 = 11
	Mask        = (1 << Fanout) - 1
)

// Key is a redefinition of []byte
type Key []byte

var (
	defaultLocker sync.FakeLocker
)

// Options holds Hamt's options
type Options struct {
	locker sync.Locker
}

// Option is a function type used to set Options
type Option func(option *Options)

// WithGoroutineSafe is used to config a Hamt with goroutine-safe
func WithGoroutineSafe() Option {
	return func(option *Options) {
		option.locker = &gosync.RWMutex{}
	}
}

// Entry is a tree node interface
type Entry interface {
	// Type returns the node type
	Type() int

	// BitPosNum returns number from a bit position
	BitPosNum(depth int) uint64
}

// BitmapNode defines Hamt's bitmap node
type BitmapNode struct {
	bitmap   uint64
	children []Entry
	pos      uint8 //position in parent array, in range [0, 64)
}

// KvPair is a list node with actually value
type KvPair struct {
	key   Key
	value interface{}
	next  *KvPair
}

//KvNode is Hamt's key-value node
type KvNode struct {
	hash   uint64
	kvList *KvPair
}

// Hamt is an implementation of hash-array-mapped-trie
type Hamt struct {
	root   BitmapNode
	locker sync.Locker
}

// Type returns the node type
func (h *BitmapNode) Type() int {
	return BITMAP_NODE
}

// BitPosNum returns the number from a bit position
func (h *BitmapNode) BitPosNum(int) uint64 {
	return uint64(1) << h.pos
}

// Index returns the index of a bitPos int bitmap
func (h *BitmapNode) Index(bitPos uint64) int {
	return bits.OnesCount64((bitPos - 1) & h.bitmap)
}

func (h *BitmapNode) insert(depth int, hash uint64, kv *KvPair) {
	pos := pos(hash, depth) //hash in current node's position
	bitPos := bitPos(pos)   //hash in current bitmap's position in bit
	if bitPos&h.bitmap == 0 {
		h.bitmap |= bitPos
		newChildren := make([]Entry, len(h.children)+1)
		kvNode := &KvNode{
			hash:   hash,
			kvList: kv,
		}
		index := h.Index(bitPos)
		newChildren[index] = kvNode
		for _, entry := range h.children {
			index = h.Index(entry.BitPosNum(depth))
			newChildren[index] = entry
		}
		h.children = newChildren
	} else {
		index := h.Index(bitPos)
		entry := h.children[index]
		if entry.Type() == KV_NODE {
			kvNode := entry.(*KvNode)
			if kvNode.hash == hash {
				for iter := kvNode.kvList; iter != nil; iter = iter.next {
					if string(iter.key) == string(kv.key) {
						iter.value = kv.value
						return
					}
				}
				kv.next = kvNode.kvList
				kvNode.kvList = kv
			} else {
				bitmapNode := &BitmapNode{
					pos: pos,
				}
				bitmapNode.insert(depth+1, kvNode.hash, kvNode.kvList)
				bitmapNode.insert(depth+1, hash, kv)
				h.children[index] = bitmapNode
			}
		} else {
			entry.(*BitmapNode).insert(depth+1, hash, kv)
		}
	}
}

func (h *BitmapNode) find(depth int, hash uint64, key Key) interface{} {
	pos := pos(hash, depth) //hash in current node's position
	bitPos := bitPos(pos)   //hash in current bitmap's position in bit
	if bitPos&h.bitmap == 0 {
		return nil
	}
	index := h.Index(bitPos)
	entry := h.children[index]
	if entry.Type() == KV_NODE {
		kvNode := entry.(*KvNode)
		if kvNode.hash != hash {
			return nil
		}

		for iter := kvNode.kvList; iter != nil; iter = iter.next {
			if string(iter.key) == string(key) {
				return iter.value
			}
		}
	} else {
		return entry.(*BitmapNode).find(depth+1, hash, key)
	}
	return nil
}

func (h *BitmapNode) traversal(visitor visitor.KvVisitor) {
	for _, entry := range h.children {
		if entry.Type() == BITMAP_NODE {
			entry.(*BitmapNode).traversal(visitor)
		} else {
			node := entry.(*KvNode)
			for kv := node.kvList; kv != nil; kv = kv.next {
				if !visitor(kv.key, kv.value) {
					return
				}
			}
		}
	}
}

func (h *BitmapNode) erase(depth int, hash uint64, key Key) bool {
	pos := pos(hash, depth) //hash in current node's position
	bitPos := bitPos(pos)   //hash in current bitmap's position in bit
	if bitPos&h.bitmap == 0 {
		return false
	}
	index := h.Index(bitPos)
	entry := h.children[index]
	if entry.Type() == KV_NODE {
		kvNode := entry.(*KvNode)
		if kvNode.hash != hash {
			return false
		}
		iter := kvNode.kvList
		var preIter *KvPair
		found := false
		for ; iter != nil; iter = iter.next {
			if string(iter.key) == string(key) {
				found = true
				break
			}
			preIter = iter
		}
		if found {
			// remove iter
			if preIter != nil {
				preIter.next = iter.next
			} else {
				kvNode.kvList = iter.next
			}
			if kvNode.kvList == nil {
				h.children[index] = nil
				h.bitmap &= ^bitPos
				newChildren := make([]Entry, len(h.children)-1)
				for _, entry := range h.children {
					if entry != nil {
						newIndex := h.Index(entry.BitPosNum(depth))
						newChildren[newIndex] = entry
					}
				}
				h.children = newChildren
			}
			return true
		}
		return false
	}

	bitmapNode := entry.(*BitmapNode)
	ok := bitmapNode.erase(depth+1, hash, key)
	// change bitmapNode to kvNode, if the a bitmapNode has only one kvNode
	if ok && len(bitmapNode.children) == 1 && bitmapNode.children[0].Type() == KV_NODE {
		child := bitmapNode.children[0].(*KvNode)
		h.children[index] = child
	}
	return ok
}

// Type returns the node type
func (h *KvNode) Type() int {
	return KV_NODE
}

// BitPosNum returns the bit position
func (h *KvNode) BitPosNum(depth int) uint64 {
	return uint64(1) << pos(h.hash, depth)
}

// New creates a Hamt(hash array mapped trie) instance
func New(opts ...Option) *Hamt {
	option := Options{
		locker: defaultLocker,
	}
	for _, opt := range opts {
		opt(&option)
	}
	return &Hamt{locker: option.locker}
}

// Insert inserts a key-value pair into the hamt
func (h *Hamt) Insert(key Key, value interface{}) {
	keyHash := hash(key)

	h.locker.Lock()
	defer h.locker.Unlock()

	h.root.insert(0, keyHash, &KvPair{key: key, value: value})
}

// Get returns the value by the passed key if the key is in the hamt, otherwise returns nil
func (h *Hamt) Get(key Key) interface{} {
	keyHash := hash(key)

	h.locker.RLock()
	defer h.locker.RUnlock()

	return h.root.find(0, keyHash, key)
}

// Erase erases the key-value pair in hamt, and returns true if succeed.
func (h *Hamt) Erase(key Key) bool {
	keyHash := hash(key)

	h.locker.Lock()
	defer h.locker.Unlock()

	return h.root.erase(0, keyHash, key)
}

// Keys returns keys in Hamt
func (h *Hamt) Keys() []Key {
	h.locker.RLock()
	defer h.locker.RUnlock()

	keys := make([]Key, 0)
	h.root.traversal(func(key, value interface{}) bool {
		keys = append(keys, key.(Key))
		return true
	})
	return keys
}

// StringKeys returns keys in Hamt
func (h *Hamt) StringKeys() []string {
	h.locker.RLock()
	defer h.locker.RUnlock()

	keys := make([]string, 0)
	h.root.traversal(func(key, value interface{}) bool {
		keys = append(keys, string(key.(Key)))
		return true
	})
	return keys
}

// Traversal traversals elements in Hamt, it will not stop until to the end or the visitor returns false
func (h *Hamt) Traversal(visitor visitor.KvVisitor) {
	h.locker.RLock()
	defer h.locker.RUnlock()

	h.root.traversal(visitor)
}

func hash(a []byte) uint64 {
	h := fnv.New64()
	h.Write(a)
	return h.Sum64()
}

func pos(hash uint64, depth int) uint8 {
	return uint8((hash >> (uint64(depth) * Fanout)) & Mask)
}

func bitPos(pos uint8) uint64 {
	return uint64(1) << pos
}

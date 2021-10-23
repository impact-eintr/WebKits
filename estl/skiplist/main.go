package main

const Level = 7

type Node struct {
	Value int   // 存储值
	Prev  *Node // 同层前节点
	Next  *Node // 同层后节点
	Down  *Node // 下层同节点
}

// 跳表是可以实现二分查找的有序链表。
type SkipList struct {
	Level       int
	HandNodeArr []*Node
}

// 是否包含节点

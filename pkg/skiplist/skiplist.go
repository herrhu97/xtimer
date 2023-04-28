package skiplist

// 参考 https://mp.weixin.qq.com/s/fvfz6bdvsZJtGsdL0MPYoA 实现

import "math/rand"

type node struct {
	nexts    []*node // 其长度对应为当前节点的高度
	key, val int
}

type Skiplist struct {
	head *node
}

func (s *Skiplist) Get(key int) (int, bool) {
	_node := s.search(key)
	if _node == nil {
		return -1, false
	}
	return _node.val, true
}

func (s *Skiplist) search(key int) *node {
	move := s.head

	for level := len(s.head.nexts) - 1; level >= 0; level-- {
		for move.nexts[level] != nil && move.nexts[level].key < key {
			move = move.nexts[level]
		}

		if move.nexts[level] != nil && move.nexts[level].key == key {
			return move.nexts[level]
		}
	}

	return nil
}

// roll 0层概率1/2,1层概率1/4,2层概率1/8...
func (s *Skiplist) roll() int {
	var level int
	for rand.Int()%2 == 0 {
		level++
	}
	return level
}

func (s *Skiplist) Put(key, val int) {
	if n := s.search(key); n != nil {
		n.val = val
		return
	}

	level := s.roll()

	for level+1 > len(s.head.nexts) {
		s.head.nexts = append(s.head.nexts, nil)
	}

	newnode := node{
		key:   key,
		val:   val,
		nexts: make([]*node, level+1),
	}

	move := s.head
	for level := level; level >= 0; level-- {
		for move.nexts[level] != nil && move.nexts[level].key < key {
			move = move.nexts[level]
		}

		newnode.nexts[level] = move.nexts[level] // 在move.nexts[level]节点左侧插入
		move.nexts[level] = &newnode
	}
}

func (s *Skiplist) Del(key int) {
	if _node := s.search(key); _node == nil {
		return
	}

	move := s.head
	for level := len(s.head.nexts) - 1; level >= 0; level-- {
		for move.nexts[level] != nil && move.nexts[level].key < key {
			move = move.nexts[level]
		}

		if move.nexts[level] == nil || move.nexts[level].key > key {
			continue
		}

		move.nexts[level] = move.nexts[level].nexts[level]
	}

	var diff int
	for level := len(s.head.nexts) - 1; level > 0 && s.head.nexts[level] == nil; level-- {
		diff++
	}

	s.head.nexts = s.head.nexts[:len(s.head.nexts)]
}

// 找到 skiplist 当中 ≥ start，且 ≤ end 的 kv 对
func (s *Skiplist) Range(start, end int) [][2]int {
	// 首先通过 ceiling 方法，找到 skiplist 中 key 值大于等于 start 且最接近于 start 的节点 ceilNode
	ceilNode := s.ceiling(start)
	// 如果不存在，直接返回
	if ceilNode == nil {
		return [][2]int{}
	}

	// 从 ceilNode 首层出发向右遍历，把所有位于 [start,end] 区间内的节点统统返回
	var res [][2]int
	for move := ceilNode; move != nil && move.key <= end; move = move.nexts[0] {
		res = append(res, [2]int{move.key, move.val})
	}
	return res
}

// 找到 key 值大于等于 target 且 key 值最接近于 target 的节点
func (s *Skiplist) ceiling(target int) *node {
	move := s.head

	// 自上而下，找到 key 值小于 target 且最接近 target 的 kv 对
	for level := len(s.head.nexts) - 1; level >= 0; level-- {
		for move.nexts[level] != nil && move.nexts[level].key < target {
			move = move.nexts[level]
		}
		// 如果 key 值等于 targe 的 kv 对存在，则直接返回
		if move.nexts[level] != nil && move.nexts[level].key == target {
			return move.nexts[level]
		}
	}

	// 此时 move 已经对应于在首层 key 值小于 key 且最接近于 key 的节点，其右侧第一个节点即为所寻找的目标节点
	return move.nexts[0]
}

// 找到 skiplist 中，key 值大于等于 target 且最接近于 target 的 key-value 对
func (s *Skiplist) Ceiling(target int) ([2]int, bool) {
	if ceilNode := s.ceiling(target); ceilNode != nil {
		return [2]int{ceilNode.key, ceilNode.val}, true
	}

	return [2]int{}, false
}

// 找到 skiplist 中，key 值小于等于 target 且最接近于 target 的 key-value 对
func (s *Skiplist) Floor(target int) ([2]int, bool) {
	// 引用 floor 方法，取 floorNode 值进行返回
	if floorNode := s.floor(target); floorNode != nil {
		return [2]int{floorNode.key, floorNode.val}, true
	}

	return [2]int{}, false
}

// 找到 key 值小于等于 target 且 key 值最接近于 target 的节点
func (s *Skiplist) floor(target int) *node {
	move := s.head

	// 自上而下，找到 key 值小于 target 且最接近 target 的 kv 对
	for level := len(s.head.nexts) - 1; level >= 0; level-- {
		for move.nexts[level] != nil && move.nexts[level].key < target {
			move = move.nexts[level]
		}
		// 如果 key 值等于 targe 的 kv对存在，则直接返回
		if move.nexts[level] != nil && move.nexts[level].key == target {
			return move.nexts[level]
		}
	}

	// move 是首层中 key 值小于 target 且最接近 target 的节点，直接返回 move 即可
	return move
}

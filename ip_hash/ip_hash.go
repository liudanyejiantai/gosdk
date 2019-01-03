// Copyright 2018 yejiantai Authors
//
// package ip_hash 一致性HASH(适合缓存和变化不大的)
package ip_hash

import (
	"crypto/sha1"
	"fmt"
	"sync"

	"math"
	"sort"
	"strconv"
)

const (
	//默认虚拟点个数
	DefaultVirualSpots = 512
)

//节点信息
type node struct {
	NodeKey   string //节点key
	spotValue uint32 //节点虚拟值
}

type NodesArray []node

func (p NodesArray) Len() int           { return len(p) }
func (p NodesArray) Less(i, j int) bool { return p[i].spotValue < p[j].spotValue }
func (p NodesArray) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p NodesArray) Sort()              { sort.Sort(p) }

// 哈希环存储节点和权重
type HashRing struct {
	virualSpots int
	Nodes       NodesArray
	weights     map[string]int
	mu          sync.RWMutex
}

func Test() {
	const (
		node1 = "192.168.122.121"
		node2 = "192.168.122.122"
		node3 = "192.168.122.123"
		node4 = "192.168.122.124"
		node5 = "192.168.122.125"
		node6 = "192.168.122.126"
		node7 = "192.168.122.127"
	)
	nodeWeight := make(map[string]int)
	nodeWeight[node1] = 2
	nodeWeight[node2] = 2
	nodeWeight[node3] = 3
	nodeWeight[node4] = 3
	nodeWeight[node5] = 3
	nodeWeight[node6] = 3
	nodeWeight[node7] = 3
	vitualSpots := 100
	hash := NewHashRing(vitualSpots)
	hash.AddNodes(nodeWeight)

	for i := 0; i < 100; i++ {
		str := fmt.Sprintf("%d", i) //str长度任意
		str2 := "KEY[" + str + "]分布存储在服务器[" + hash.GetNode(str) + "]上"
		fmt.Println(str2)
	}
}

// 创建具有虚拟点的哈希环
func NewHashRing(spots int) *HashRing {
	if spots == 0 {
		spots = DefaultVirualSpots
	}

	h := &HashRing{
		virualSpots: spots,
		weights:     make(map[string]int),
	}
	return h
}

// 添加节点到哈希环
func (h *HashRing) AddNodes(nodeWeight map[string]int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for nodeKey, w := range nodeWeight {
		h.weights[nodeKey] = w
	}
	h.generate()
}

// 添加节点到哈希环
func (h *HashRing) AddNode(nodeKey string, weight int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.weights[nodeKey] = weight
	h.generate()
}

// 删除哈希环中的节点
func (h *HashRing) RemoveNode(nodeKey string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.weights, nodeKey)
	h.generate()
}

// 修改哈希环中节点的权重
func (h *HashRing) UpdateNode(nodeKey string, weight int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.weights[nodeKey] = weight
	h.generate()
}

func (h *HashRing) generate() {
	var totalW int
	for _, w := range h.weights {
		totalW += w
	}

	totalVirtualSpots := h.virualSpots * len(h.weights)
	h.Nodes = NodesArray{}

	for nodeKey, w := range h.weights {
		spots := int(math.Floor(float64(w) / float64(totalW) * float64(totalVirtualSpots)))
		for i := 1; i <= spots; i++ {
			hash := sha1.New()
			hash.Write([]byte(nodeKey + ":" + strconv.Itoa(i)))
			hashBytes := hash.Sum(nil)
			n := node{
				NodeKey:   nodeKey,
				spotValue: genValue(hashBytes[6:10]),
			}
			h.Nodes = append(h.Nodes, n)
			hash.Reset()
		}
	}
	h.Nodes.Sort()
}

func genValue(bs []byte) uint32 {
	if len(bs) < 4 {
		return 0
	}
	v := (uint32(bs[3]) << 24) | (uint32(bs[2]) << 16) | (uint32(bs[1]) << 8) | (uint32(bs[0]))
	return v
}

// 获得节点值
func (h *HashRing) GetNode(s string) string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if len(h.Nodes) == 0 {
		return ""
	}

	hash := sha1.New()
	hash.Write([]byte(s))
	hashBytes := hash.Sum(nil)
	v := genValue(hashBytes[6:10])
	i := sort.Search(len(h.Nodes), func(i int) bool { return h.Nodes[i].spotValue >= v })

	if i == len(h.Nodes) {
		i = 0
	}
	return h.Nodes[i].NodeKey
}

package hnsw

import (
	"math"
	"math/rand"
	"sync"
)

type HNSWIndex struct {
	nodes    map[string]*HNSWNode
	dim      int
	m        int
	ef       int
	maxLayer int
	mutex    sync.RWMutex
}

type HNSWNode struct {
	ID        string
	Vector    []float64
	Neighbors map[int][]Neighbor
	Layer     int
}

type Neighbor struct {
	ID         string
	Similarity float64
}

func NewHNSWIndex(dim, m, ef int) *HNSWIndex {
	return &HNSWIndex{
		nodes:    make(map[string]*HNSWNode),
		dim:      dim,
		m:        m,
		ef:       ef,
		maxLayer: 0,
	}
}

func (idx *HNSWIndex) Add(id string, vector []float64) {
	idx.mutex.Lock()
	defer idx.mutex.Unlock()

	layer := int(math.Floor(-math.Log(rand.Float64()) * float64(idx.m)))
	if layer > idx.maxLayer {
		idx.maxLayer = layer
	}

	node := &HNSWNode{
		ID:        id,
		Vector:    vector,
		Neighbors: make(map[int][]Neighbor),
		Layer:     layer,
	}
	idx.nodes[id] = node

	if len(idx.nodes) == 1 {
		return
	}

	for l := idx.maxLayer; l >= 0; l-- {
		if l > layer {
			continue
		}
		neighbors := idx.searchLayer(vector, l, idx.ef)
		node.Neighbors[l] = neighbors
		for _, n := range neighbors {
			neighborNode := idx.nodes[n.ID]
			neighborNode.Neighbors[l] = append(neighborNode.Neighbors[l], Neighbor{ID: id, Similarity: n.Similarity})
			if len(neighborNode.Neighbors[l]) > idx.m {
				neighborNode.Neighbors[l] = neighborNode.Neighbors[l][:idx.m]
			}
		}
	}
}

func (idx *HNSWIndex) Remove(id string) {
	idx.mutex.Lock()
	defer idx.mutex.Unlock()

	node, exists := idx.nodes[id]
	if !exists {
		return
	}

	for l := 0; l <= node.Layer; l++ {
		for _, n := range node.Neighbors[l] {
			neighbor := idx.nodes[n.ID]
			for i, nb := range neighbor.Neighbors[l] {
				if nb.ID == id {
					neighbor.Neighbors[l] = append(neighbor.Neighbors[l][:i], neighbor.Neighbors[l][i+1:]...)
					break
				}
			}
		}
	}
	delete(idx.nodes, id)
}

func (idx *HNSWIndex) Search(query []float64, k int) []Neighbor {
	idx.mutex.RLock()
	defer idx.mutex.RUnlock()

	if len(idx.nodes) == 0 {
		return nil
	}

	var entryPoint *HNSWNode
	for _, node := range idx.nodes {
		entryPoint = node
		break
	}

	for l := idx.maxLayer; l > 0; l-- {
		entryPoint = idx.searchLayerEntry(query, entryPoint, l)
	}

	return idx.searchLayer(query, 0, k)
}

func (idx *HNSWIndex) searchLayerEntry(query []float64, entry *HNSWNode, layer int) *HNSWNode {
	current := entry
	for {
		best := current
		bestSim := cosineSimilarity(query, current.Vector)
		for _, n := range current.Neighbors[layer] {
			sim := cosineSimilarity(query, idx.nodes[n.ID].Vector)
			if sim > bestSim {
				best = idx.nodes[n.ID]
				bestSim = sim
			}
		}
		if best == current {
			break
		}
		current = best
	}
	return current
}

func (idx *HNSWIndex) searchLayer(query []float64, layer, k int) []Neighbor {
	visited := make(map[string]bool)
	candidates := make([]Neighbor, 0)
	result := make([]Neighbor, 0, k)

	var entry *HNSWNode
	for _, node := range idx.nodes {
		entry = node
		break
	}
	candidates = append(candidates, Neighbor{ID: entry.ID, Similarity: cosineSimilarity(query, entry.Vector)})
	visited[entry.ID] = true

	for len(candidates) > 0 && len(result) < k {
		closest := candidates[0]
		candidates = candidates[1:]

		if len(result) > 0 && closest.Similarity < result[len(result)-1].Similarity {
			break
		}

		result = insertSorted(result, closest, k)
		for _, n := range idx.nodes[closest.ID].Neighbors[layer] {
			if !visited[n.ID] {
				visited[n.ID] = true
				sim := cosineSimilarity(query, idx.nodes[n.ID].Vector)
				candidates = append(candidates, Neighbor{ID: n.ID, Similarity: sim})
			}
		}
	}
	return result
}

func insertSorted(list []Neighbor, n Neighbor, k int) []Neighbor {
	list = append(list, n)
	for i := len(list) - 1; i > 0 && list[i].Similarity > list[i-1].Similarity; i-- {
		list[i], list[i-1] = list[i-1], list[i]
	}
	if len(list) > k {
		return list[:k]
	}
	return list
}

func cosineSimilarity(v1, v2 []float64) float64 {
	if len(v1) != len(v2) {
		return 0
	}
	var dotProduct, norm1, norm2 float64
	for i := range v1 {
		dotProduct += v1[i] * v2[i]
		norm1 += v1[i] * v1[i]
		norm2 += v2[i] * v2[i]
	}
	if norm1 == 0 || norm2 == 0 {
		return 0
	}
	return dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2))
}

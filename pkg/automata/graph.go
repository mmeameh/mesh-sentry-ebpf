package automata

import (
	"sync"
	"time"
)

type Edge struct {
	FromService string
	ToService   string
	Latency     time.Duration
	Timestamp   time.Time
}

type DependencyGraph struct {
	mu        sync.RWMutex
	adjList   map[string][]string
	edgeTimes map[string]time.Time
}

func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		adjList:   make(map[string][]string),
		edgeTimes: make(map[string]time.Time),
	}
}

func (g *DependencyGraph) AddDependency(from, to string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, existing := range g.adjList[from] {
		if existing == to {
			return
		}
	}
	g.adjList[from] = append(g.adjList[from], to)
}

func (g *DependencyGraph) RemoveDependency(from, to string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	edges := g.adjList[from]
	for i, existing := range edges {
		if existing == to {
			g.adjList[from] = append(edges[:i], edges[i+1:]...)
			return
		}
	}
}

func (g *DependencyGraph) CloneAdjList() map[string][]string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	clone := make(map[string][]string)
	for k, v := range g.adjList {
		copied := make([]string, len(v))
		copy(copied, v)
		clone[k] = copied
	}
	return clone
}

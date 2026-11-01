package automata

// DeadlockDetector analyzes asynchronous call graphs for circular waits and cascading deadlocks.
type DeadlockDetector struct{}

func NewDeadlockDetector() *DeadlockDetector {
	return &DeadlockDetector{}
}

// FindCycles executes Tarjan's SCC algorithm and returns all circular dependency chains.
func (d *DeadlockDetector) FindCycles(graph *DependencyGraph) [][]string {
	adj := graph.CloneAdjList()

	index := 0
	indices := make(map[string]int)
	lowlink := make(map[string]int)
	onStack := make(map[string]bool)
	var stack []string
	var cycles [][]string

	var strongConnect func(node string)
	strongConnect = func(node string) {
		indices[node] = index
		lowlink[node] = index
		index++
		stack = append(stack, node)
		onStack[node] = true

		for _, neighbor := range adj[node] {
			if _, visited := indices[neighbor]; !visited {
				strongConnect(neighbor)
				if lowlink[neighbor] < lowlink[node] {
					lowlink[node] = lowlink[neighbor]
				}
			} else if onStack[neighbor] {
				if indices[neighbor] < lowlink[node] {
					lowlink[node] = indices[neighbor]
				}
			}
		}

		if lowlink[node] == indices[node] {
			var scc []string
			for {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				onStack[top] = false
				scc = append(scc, top)
				if top == node {
					break
				}
			}

			// A strongly connected component with >1 node (or a self-loop) indicates a circular wait / deadlock
			if len(scc) > 1 {
				cycles = append(cycles, scc)
			} else if len(scc) == 1 {
				// Check for self-loop
				for _, neighbor := range adj[node] {
					if neighbor == node {
						cycles = append(cycles, scc)
						break
					}
				}
			}
		}
	}

	for node := range adj {
		if _, visited := indices[node]; !visited {
			strongConnect(node)
		}
	}

	return cycles
}

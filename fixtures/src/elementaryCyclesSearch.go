package main

import (
	"math"
)

/**
 * Calculates the adjacency-list for a given adjacency-matrix.
 *
 * This is based on the Java implementation of :
 *
 * @author Frank Meyer, web@normalisiert.de
 * @version 1.0, 26.08.2006
 *
 */

// GetAdjacencyList : static method
/**
 * Calculates an adjacency-list for a given array of an adjacency-matrix.
 *
 * @param adjacencyMatrix array with the adjacency-matrix that represents
 * the graph
 * @return int[][]-array of the adjacency-list of given nodes. The first
 * dimension in the array represents the same node as in the given
 * adjacency, the second dimension represents the indices of those nodes,
 * that are direct successor nodes of the node.
 */
func GetAdjacencyList(adjacencyMatrix [][]bool) [][]int {
	var list [][]int
	list = make([][]int, len(adjacencyMatrix))

	for i := 0; i < len(adjacencyMatrix); i++ {
		var v []int
		v = make([]int, 0)
		for j := 0; j < len(adjacencyMatrix[i]); j++ {
			if adjacencyMatrix[i][j] {
				v = append(v, j) //v.add(j)
			}
		}

		list[i] = make([]int, len(v))
		for j := 0; j < len(v); j++ {
			in := v[j]
			list[i][j] = in
		}
	}

	return list
}

//===========================================================================================================================================

/**
 * This is a helpclass for the search of all elementary cycles in a graph
 * with the algorithm of Johnson. For this it searches for strong connected
 * components, using the algorithm of Tarjan. The constructor gets an
 * adjacency-list of a graph. Based on this graph, it gets a nodenumber s,
 * for which it calculates the subgraph, containing all nodes
 * {s, s + 1, ..., n}, where n is the highest nodenumber in the original
 * graph (e.g. it builds a subgraph with all nodes with higher or same
 * nodenumbers like the given node s). It returns the strong connected
 * component of this subgraph which contains the lowest nodenumber of all
 * nodes in the subgraph.
 *
 * For a description of the algorithm for calculating the strong connected
 * components see:
 * Robert Tarjan: Depth-first search and linear graph algorithms. In: SIAM
 * Journal on Computing. Volume 1, Nr. 2 (1972), pp. 146-160.
 * For a description of the algorithm for searching all elementary cycles in
 * a directed graph see:
 * Donald B. Johnson: Finding All the Elementary Circuits of a Directed Graph.
 * SIAM Journal on Computing. Volumne 4, Nr. 1 (1975), pp. 77-84.
 *
 * This is based on the Java implementation of :
 * @author Frank Meyer, web_at_normalisiert_dot_de
 * @version 1.1, 22.03.2009
 *
 */

// StrongConnectedComponents : represent a set of scc's
type StrongConnectedComponents struct {
	/** Adjacency-list of original graph */
	adjListOriginal [][]int

	/** Adjacency-list of currently viewed subgraph */
	adjList [][]int

	/** Helpattribute for finding scc's */
	visited []bool

	/** Helpattribute for finding scc's */
	stack []int

	/** Helpattribute for finding scc's */
	lowlink []int

	/** Helpattribute for finding scc's */
	number []int

	/** Helpattribute for finding scc's */
	sccCounter int

	/** Helpattribute for finding scc's */
	currentSCCs [][]int
}

// NewStrongConnectedComponents :
/**
 * Constructor.
 *
 * @param adjList adjacency-list of the graph
 */
func NewStrongConnectedComponents(adjList [][]int) *StrongConnectedComponents {
	this := new(StrongConnectedComponents)
	this.adjListOriginal = adjList
	return this
}

// getAdjacencyList :
/**
 * This method returns the adjacency-structure of the strong connected
 * component with the least vertex in a subgraph of the original graph
 * induced by the nodes {s, s + 1, ..., n}, where s is a given node. Note
 * that trivial strong connected components with just one node will not
 * be returned.
 *
 * @param node node s
 * @return SCCResult with adjacency-structure of the strong
 * connected component; null, if no such component exists
 */
func (this *StrongConnectedComponents) getAdjacencyList(node int) *SCCResult {

	this.visited = make([]bool, len(this.adjListOriginal))
	this.lowlink = make([]int, len(this.adjListOriginal))
	this.number = make([]int, len(this.adjListOriginal))
	this.visited = make([]bool, len(this.adjListOriginal))
	this.stack = make([]int, 0)
	this.currentSCCs = make([][]int, 0)

	this.makeAdjListSubgraph(node)

	for i := node; i < len(this.adjListOriginal); i++ {
		if !this.visited[i] {
			this.getStrongConnectedComponents(i)
			var nodes []int
			nodes = this.getLowestIdComponent()
			if nodes != nil && !contains(nodes, node) && !contains(nodes, node+1) {
				return this.getAdjacencyList(node + 1)
			} else {
				var adjacencyList [][]int
				adjacencyList = this.getAdjList(nodes)
				if adjacencyList != nil {
					for j := 0; j < len(this.adjListOriginal); j++ {
						if len(adjacencyList[j]) > 0 {
							result := NewSCCResult(adjacencyList, j)
							return result
						}
					}
				}
			}
		}
	}

	return nil
}

// makeAdjListSubgraph :
/**
 * Builds the adjacency-list for a subgraph containing just nodes
 * >= a given index.
 *
 * @param node Node with lowest index in the subgraph
 */
func (this *StrongConnectedComponents) makeAdjListSubgraph(node int) {
	this.adjList = make([][]int, len(this.adjListOriginal)) // = new int[this.adjListOriginal.length][0];
	for i := range this.adjList {
		this.adjList[i] = make([]int, 0)
	}

	for i := node; i < len(this.adjList); i++ {
		var successors []int
		successors = make([]int, 0)
		for j := 0; j < len(this.adjListOriginal[i]); j++ {
			if this.adjListOriginal[i][j] >= node {
				successors = append(successors, this.adjListOriginal[i][j])
			}
		}
		if len(successors) > 0 {
			this.adjList[i] = make([]int, len(successors))
			for j := 0; j < len(successors); j++ {
				var succ int
				succ = successors[j]
				this.adjList[i][j] = succ
			}
		}
	}
}

// getLowestIdComponent :
/**
 * Calculates the strong connected component out of a set of scc's, that
 * contains the node with the lowest index.
 *
 * @return Vector::Integer of the scc containing the lowest nodenumber
 */
func (this *StrongConnectedComponents) getLowestIdComponent() []int {
	min := len(this.adjList)
	var currScc []int
	currScc = nil

	for i := 0; i < len(this.currentSCCs); i++ {
		var scc []int
		scc = this.currentSCCs[i]
		for j := 0; j < len(scc); j++ {
			var node int
			node = scc[j]
			if node < min {
				currScc = scc
				min = node
			}
		}
	}

	return currScc
}

// getAdjList
/**
 * @return Vector[]::Integer representing the adjacency-structure of the
 * strong connected component with least vertex in the currently viewed
 * subgraph
 */
func (this *StrongConnectedComponents) getAdjList(nodes []int) [][]int {
	//Vector[] lowestIdAdjacencyList = null;
	var lowestIdAdjacencyList [][]int
	lowestIdAdjacencyList = nil

	if nodes != nil {
		lowestIdAdjacencyList = make([][]int, len(this.adjList))
		for i := 0; i < len(lowestIdAdjacencyList); i++ {
			lowestIdAdjacencyList[i] = make([]int, 0)
		}
		for i := 0; i < len(nodes); i++ {
			node := nodes[i]
			for j := 0; j < len(this.adjList[node]); j++ {
				succ := this.adjList[node][j]
				if contains(nodes, succ) {
					lowestIdAdjacencyList[node] = append(lowestIdAdjacencyList[node], succ)
				}
			}
		}
	}

	return lowestIdAdjacencyList
}

// getStrongConnectedComponents :
/**
 * Searches for strong connected components reachable from a given node.
 *
 * @param root node to start from.
 */
func (this *StrongConnectedComponents) getStrongConnectedComponents(root int) {
	this.sccCounter++
	this.lowlink[root] = this.sccCounter
	this.number[root] = this.sccCounter
	this.visited[root] = true
	this.stack = append(this.stack, root)

	for i := 0; i < len(this.adjList[root]); i++ {
		w := this.adjList[root][i]
		if !this.visited[w] {
			this.getStrongConnectedComponents(w)
			this.lowlink[root] = int(math.Min(float64(this.lowlink[root]), float64(this.lowlink[w])))
		} else if this.number[w] < this.number[root] {
			if contains(this.stack, w) {
				this.lowlink[root] = int(math.Min(float64(this.lowlink[root]), float64(this.number[w])))
			}
		}
	}

	// found scc
	if (this.lowlink[root] == this.number[root]) && (len(this.stack) > 0) {
		next := -1
		var scc []int
		scc = make([]int, 0)

		// do while equivalent
		for ok := true; ok; ok = (this.number[next] > this.number[root]) {
			next = this.stack[len(this.stack)-1]
			this.stack = this.stack[:len(this.stack)-1]
			scc = append(scc, next)
		}

		// simple scc's with just one node will not be added
		if len(scc) > 1 {
			this.currentSCCs = append(this.currentSCCs, scc)
		}
	}
}

//==================================================================================================================================

// SCCResult : represents the adjacency structure of a set of Strongly Connected Components
type SCCResult struct {
	nodeIDsOfSCC []int //private Set nodeIDsOfSCC = null; // Set is a list without duplicates || apparently works with ordinary list, but we should probably replicate the behaviour of Set
	adjList      [][]int
	lowestNodeId int //private int lowestNodeId = -1;
}

// NewSCCResult : Constructor
func NewSCCResult(adjList [][]int, lowestNodeId int) *SCCResult {
	this := new(SCCResult)
	this.adjList = adjList
	this.lowestNodeId = lowestNodeId
	this.nodeIDsOfSCC = make([]int, 0) //this.nodeIDsOfSCC = new HashSet(); // TODO : replicate behaviour of Set
	if this.adjList != nil {
		for i := this.lowestNodeId; i < len(this.adjList); i++ {
			if len(this.adjList[i]) > 0 {
				this.nodeIDsOfSCC = append(this.nodeIDsOfSCC, i)
			}
		}
	}
	return this
}

func (this *SCCResult) getAdjList() [][]int {
	return this.adjList
}

func (this *SCCResult) getLowestNodeId() int {
	return this.lowestNodeId
}

//==================================================================================================================================

/**
 * Searchs all elementary cycles in a given directed graph. The implementation
 * is independent from the concrete objects that represent the graphnodes, it
 * just needs an array of the objects representing the nodes the graph
 * and an adjacency-matrix of type boolean, representing the edges of the
 * graph. It then calculates based on the adjacency-matrix the elementary
 * cycles and returns a list, which contains lists itself with the objects of the
 * concrete graphnodes-implementation. Each of these lists represents an
 * elementary cycle.
 *
 * The implementation uses the algorithm of Donald B. Johnson for the search of
 * the elementary cycles. For a description of the algorithm see:
 * Donald B. Johnson: Finding All the Elementary Circuits of a Directed Graph.
 * SIAM Journal on Computing. Volumne 4, Nr. 1 (1975), pp. 77-84.
 *
 * The algorithm of Johnson is based on the search for strong connected
 * components in a graph. For a description of this part see:
 * Robert Tarjan: Depth-first search and linear graph algorithms. In: SIAM
 * Journal on Computing. Volume 1, Nr. 2 (1972), pp. 146-160.
 *
 * This is based on the Java implementation of :
 * @author Frank Meyer, web_at_normalisiert_dot_de
 * @version 1.2, 22.03.2009
 *
 */
type ElementaryCyclesSearch struct {
	/** List of cycles */
	cycles [][]int

	/** Adjacency-list of graph */
	adjList [][]int

	/** Graphnodes */
	graphNodes []int

	/** Blocked nodes, used by the algorithm of Johnson */
	blocked []bool

	/** B-Lists, used by the algorithm of Johnson */
	B [][]int

	/** Stack for nodes, used by the algorithm of Johnson */
	stack []int
}

// NewElementaryCyclesSearch : Constructor
/**
 * @param matrix adjacency-matrix of the graph
 * @param graphNodes array of the graphnodes of the graph; this is used to
 * build sets of the elementary cycles containing the objects of the original
 * graph-representation
 */
func NewElementaryCyclesSearch(matrix [][]bool, graphNodes []int) *ElementaryCyclesSearch {
	ecs := new(ElementaryCyclesSearch)
	ecs.graphNodes = graphNodes
	ecs.adjList = GetAdjacencyList(matrix)
	return ecs
}

// GetElementaryCycles :
/**
 * Returns List::List::Object with the Lists of nodes of all elementary
 * cycles in the graph.
 *
 * @return List::List::Object with the Lists of the elementary cycles.
 */
func (this *ElementaryCyclesSearch) GetElementaryCycles() [][]int {
	this.cycles = make([][]int, 0)
	this.blocked = make([]bool, len(this.adjList))
	this.B = make([][]int, len(this.adjList))
	this.stack = make([]int, 0)
	var sccs *StrongConnectedComponents
	sccs = NewStrongConnectedComponents(this.adjList)
	s := 0

	for true {
		var sccResult *SCCResult
		sccResult = sccs.getAdjacencyList(s)
		if sccResult != nil && sccResult.getAdjList() != nil {
			var scc [][]int
			scc = sccResult.getAdjList()
			s = sccResult.getLowestNodeId()
			for j := 0; j < len(scc); j++ {
				if (scc[j] != nil) && (len(scc[j]) > 0) {
					this.blocked[j] = false
					this.B[j] = make([]int, 0)
				}
			}

			this.findCycles(s, s, scc)
			s++
		} else {
			break
		}
	}

	return this.cycles
}

// findCycles :
/**
 * Calculates the cycles containing a given node in a strongly connected
 * component. The method calls itself recursivly.
 *
 * @param v
 * @param s
 * @param adjList adjacency-list with the subgraph of the strongly
 * connected component s is part of.
 * @return true, if cycle found; false otherwise
 */
func (this *ElementaryCyclesSearch) findCycles(v int, s int, adjList [][]int) bool {
	f := false
	//this.stack[len(this.stack)-1] = v
	this.stack = append(this.stack, v)
	this.blocked[v] = true

	for i := 0; i < len(adjList[v]); i++ {
		w := adjList[v][i]
		// found cycle
		if w == s {
			var cycle []int
			cycle = make([]int, 0)
			for j := 0; j < len(this.stack); j++ {
				index := this.stack[j]
				cycle = append(cycle, this.graphNodes[index])
			}
			this.cycles = append(this.cycles, cycle)
			f = true
		} else if !this.blocked[w] {
			if this.findCycles(w, s, adjList) {
				f = true
			}
		}
	}

	if f {
		this.unblock(v)
	} else {
		for i := 0; i < len(adjList[v]); i++ {
			w := adjList[v][i]
			if !contains(this.B[w], v) {
				this.B[w] = append(this.B[w], v)
			}
		}
	}

	this.stack = remove(this.stack, v) //this.stack.remove(v) : v is the object to remove

	return f
}

// remove : remove an element from a slice
func remove(s []int, r int) []int {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

// contains : returns true if slice contains element, else false
func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// unblock :
/**
 * Unblocks recursivly all blocked nodes, starting with a given node.
 *
 * @param node node to unblock
 */
func (this *ElementaryCyclesSearch) unblock(node int) {
	this.blocked[node] = false
	Bnode := this.B[node]
	for len(Bnode) > 0 {
		w := Bnode[0]
		Bnode = append(Bnode[:0], Bnode[0+1:]...) // remove first element - TODO : check if element are shifted correctly (i.e index 1 becomes 0 etc.)
		if this.blocked[w] {
			this.unblock(w)
		}
	}
}

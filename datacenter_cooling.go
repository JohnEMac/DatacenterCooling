// Author: Chernoy Viacheslav
// email: vchenoy@gmail.com
//
// The program solves "Datacenter Cooling" Problem presented by Quora
// http://www.quora.com/challenges

package main

import (
    "fmt"
)

type (
    // the 0/1/2/3-matrix type
    grid_t [][] int

    vertex_t int

    // the graph represented by adjacency list
    // the list type (maximum 4 (four) adjacent verticies exist
    graph_t [][] vertex_t
)

var (
    count int
    visited [] bool
    graph graph_t
    src, dst vertex_t
    R [] vertex_t
)


// list sub-system
func index(list [] vertex_t, x vertex_t) int {
    i := len(list)-1
    for (i >= 0 && list[i] != x) {
        i -= 1
    }

    return i
}

// reads 0/1/2/3 matrix from the standard-input
func read_input(datacenter *grid_t) {
    var (
        w, h int
    )

    _, _ = fmt.Scanf("%d %d", &w, &h)
    *datacenter = make(grid_t, h)
    for i := range *datacenter {
        (*datacenter)[i] = make([]int, w)
        for j := range (*datacenter)[i] {
            _, _ = fmt.Scanf("%d", &(*datacenter)[i][j])
        }
    }
}

// converts the pair (i, j) to a vertex
func vert(w int, i, j int) vertex_t {
    return vertex_t(i * w + j)
}

// datacenter is the 0/1/2/3-grid
// graph is adjacency list contructed from the matrix datacenter
// builds the graph representation of datacenter
// also returns source and destination vertices
func build_graph(datacenter grid_t, graph *graph_t, src, dst *vertex_t) {
    // four directions up, down, left, right
    const (
        // the room types
        OWN       = 0
        START     = 2
        END       = 3
        DONOT_OWN = 1
    )

    Dh := [4]int{-1, 1, 0, 0}
    Dw := [4]int{0, 0, -1, 1}

    // create empty graph
    n := len(datacenter) * len(datacenter[0])
    *graph = make(graph_t, 0, n)

    for i := range datacenter {
        for j := range datacenter[i] {
            // the current vertex v
            v := vert(len(datacenter[i]), i, j)
            // adjacency list of v
            adj := make([]vertex_t, 0, 4)
            if datacenter[i][j] == OWN || datacenter[i][j] == START {
                // find all verticies adjacent to v
                for k := range Dh {
                    i1 := i + Dh[k]
                    j1 := j + Dw[k]
                    if (i1 >= 0) && (i1 < len(datacenter)) && (j1 >= 0) && (j1 < len(datacenter[i])) && (datacenter[i1][j1] != DONOT_OWN) {
                        u := vert(len(datacenter[i]), i1, j1)
                        adj = append(adj, u)
                    }
                }
            }
            *graph = append(*graph, adj)
            if datacenter[i][j] == START {
                *src = v
            } else if datacenter[i][j] == END {
                *dst = v
            }
        }
    }
}

// graph represented by adjacency list
// L is the list of vertices to check
// u is a destination vertex
// checks whether the vertex u is reachable from all the vertices of L
func connected(graph graph_t, L [] vertex_t, u vertex_t) bool {
    // global: R

    // C[v] == 0 means it is not connected
    // for every k >= 2, {v|C[v] == k} are connected to v = L[1-k] (and between each other)
    C := make([]int, len(graph))
    C[u] = 1
    mark := 2
    for _, v := range L {
        // check whether v can reach u
        if C[v] == 0 {
            R = R[:1]
            R[0] = v
            lb:
            for {
                if len(R) == 0 {
                    // v cannot reach u
                    return false
                }
                // x is reachable from v
                x := R[len(R)-1]
                R = R[:len(R)-1]
                // check its adjacent ones
                for _, y := range graph[x] {
                    if C[y] == 0 {
                        // y is reachable from v
                        C[y] = mark
                        R = append(R, y)
                    } else if C[y] < mark {
                        // u is reachable from y,
                        // so we are done with v
                        break lb 
                    }
                }
            }
            mark += 1
        }
    }
    return true
}

// v is the current vertex
// l is the number of steps remained to do
// backtracking, computes the number of possible paths of length l from v to dst
func search(v vertex_t, l int) {
    // global: dst, count, visited, graph

    X := make([]vertex_t, 0, 4)
    Y := make([]vertex_t, 0, 4)
    for _, u := range graph[v] {
        if !visited[u] {
            Y = append(Y, u)
        }
    }
    // remove all the edges to v
    for _, u := range Y {
        if index(graph[u], v) >= 0 {
            if len(graph[u]) <= 1 {
                // u is isolated, rollback
                for _, u := range X {
                    graph[u] = append(graph[u], v)
                }
                return
            }
            // updating the graph
            i := index(graph[u], v)
            graph[u][i] = graph[u][len(graph[u])-1]
            graph[u] = graph[u][:len(graph[u])-1]

            // save, for later restore it
            X = append(X, u)
        }
    }
    visited[v] = true
    // check that dst is reachable from all the vertices of X 
    if len(X) == 0 || connected(graph, X, dst) {
        // one step less is remained to dst
        l -= 1
        // go over all possible steps that can be done
        // i.e., check every unvisited adjacent vertex
        for _, u := range Y {
            if l == 0 && u == dst {
                // we reached dst and passed all verticies
                count += 1
                //if (count % 1000 == 0) {
                //    fmt.Printf("%d\n", count)
                //}
            } else if l > 0 && u != dst {
                // not yet reached dst, let's check the step
                search(u, l)
            }
        }
    }
    visited[v] = false

    // restore the graph
    for _, u := range X {
        graph[u] = append(graph[u], v)
    }
}

func main() {
    var (
        datacenter grid_t
    )

    read_input(&datacenter)
    build_graph(datacenter, &graph, &src, &dst)
    // compute the number of verticies (the path length) we have to pass to reach dst
    n := len(graph)
    R = make([]vertex_t, 0, n)
    length := 0
    for _, adj := range graph {
        if len(adj) > 0 {
            length += 1
        }
    }
    visited = make([]bool, n)
    // the number of pathes
    count = 0
    search(src, length)
    fmt.Printf("%d\n", count)
}


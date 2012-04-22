// Author: Chernoy Viacheslav
// email: vchenoy@gmail.com
//
// The program solves "Datacenter Cooling" Problem presented by Quora
// http://www.quora.com/challenges

package main

import (
    "fmt"
    "os"
)

const (
    MAX_DIM = 10
    MAX_N = MAX_DIM * MAX_DIM
)

type (
    // the 0/1/2/3-matrix type
    grid_t struct {
        h, w int
        arr [MAX_DIM][MAX_DIM] int
    }

    vertex_t int

    // the list type (maximum 4 (four) adjacent verticies exist
    list_t struct {
        len int
        elems [4] vertex_t
    }

    // the graph represented by adjacency list
    graph_t struct {
        n int
        adjs [MAX_N] list_t
    }
)

func assert(cond bool) {
    if !cond  {
        fmt.Printf("assert! \n")
        os.Exit(1)
    }
}

// list sub-system
func list__init(list *list_t) {
    list.len = 0
}

func list__append(list *list_t, x vertex_t) {
    list.elems[list.len] = x
    list.len += 1
}

func list__pop(list *list_t) vertex_t {
    x := list.elems[0]
    list.elems[0] = list.elems[list.len-1]
    list.len -= 1

    return x
}

func list__remove(list *list_t, x vertex_t) {
    i := list.len-1
    for i >= 0 && list.elems[i] != x  {
        i -= 1
    }

    list.elems[i] = list.elems[list.len-1]
    list.len -= 1
}

func list__contains(list *list_t, x vertex_t) bool {
    i := list.len-1
    for (i >= 0 && list.elems[i] != x) {
        i -= 1
    }

    return i >= 0
}

// reads 0/1/2/3 matrix from the standard-input
func read_input(datacenter *grid_t) {
    _, _ = fmt.Scanf("%d %d", &datacenter.w, &datacenter.h)
    for i := 0; i < datacenter.h; i++ {
        for j := 0; j < datacenter.w; j++ {
            _, _ = fmt.Scanf("%d", &datacenter.arr[i][j])
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
func build_graph(datacenter *grid_t, graph *graph_t, src, dst *vertex_t) {
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
    graph.n = datacenter.w * datacenter.h
    for i := 0; i < MAX_N; i++ {
        list__init(&graph.adjs[i])
    }

    for i := 0; i < datacenter.h; i++ {
        for j := 0; j < datacenter.w; j++ {
            // the current vertex v
            v := vert(datacenter.w, i, j)
            // adjacency list of v
            adj := &graph.adjs[v]
            if datacenter.arr[i][j] == OWN || datacenter.arr[i][j] == START {
                // find all verticies adjacent to v
                for k := 0; k < 4; k++ {
                    i1 := i + Dh[k]
                    j1 := j + Dw[k]
                    if (i1 >= 0) && (i1 < datacenter.h) && (j1 >= 0) && (j1 < datacenter.w) && (datacenter.arr[i1][j1] != DONOT_OWN) {
                        u := vert(datacenter.w, i1, j1)
                        list__append(adj, u)
                    }
                }
            }
            if datacenter.arr[i][j] == START {
                *src = v
            }else if datacenter.arr[i][j] == END {
                *dst = v
            }
        }
    }
}

// graph represented by adjacency list
// L is the list of vertices to check
// u is a destination vertex
// checks whether the vertex u is reachable from all the vertices of L
func connected(graph *graph_t, L list_t, u vertex_t) bool {
    // C[v] == 0 means it is not connected
    // for every k >= 2, {v|C[v] == k} are connected to v = L[1-k] (and between each other)
    var (
        i int
        C [MAX_N] int
        // R is a list of verticies reachable from v
        R struct  {
            len int
            elems [MAX_N] vertex_t
        }
    )

    for i = MAX_N-1; i >= 0; i-- {
        C[i] = 0
    }

    C[u] = 1
    mark := 2
    for L.len > 0 {
        v := list__pop(&L)
        // check whether v can reach u
        if C[v] == 0 {
            R.elems[0] = v
            R.len = 1
            for true {
                if R.len == 0 {
                    // v cannot reach u
                    return false
                }
                // x is reachable from v
                x := R.elems[R.len-1]
                R.len -= 1
                // check its adjacent ones
                for i = graph.adjs[x].len-1; i >= 0; i-- {
                    y := graph.adjs[x].elems[i]
                    if C[y] == 0 {
                        // y is reachable from v
                        C[y] = mark
                        R.elems[R.len] = y
                        R.len += 1
                    }else if C[y] < mark {
                        // u is reachable from y,
                        // so we are done with v
                        break
                    }
                }
                if i >= 0 {
                    break
                }
            }
            mark += 1
        }
    }
    return true
}

var (
    count int
    visited [MAX_N] bool
    graph graph_t
    src, dst vertex_t
)

// v is the current vertex
// l is the number of steps remained to do
// backtracking, computes the number of possible paths of length l from v to dst
func search(v vertex_t, l int) {
    // global dst, count, visited, graph
    var (
        X, Y list_t
    )

    list__init(&X)
    list__init(&Y)
    for i := graph.adjs[v].len-1; i >= 0; i-- {
        u := graph.adjs[v].elems[i]
        if !visited[u] {
            list__append(&Y, u)
        }
    }
    // remove all the edges to v
    for i := Y.len-1; i >= 0; i-- {
        u := Y.elems[i]
        if list__contains(&graph.adjs[u], v) {
            if graph.adjs[u].len <= 1 {
                // u is isolated, rollback
                for i = X.len-1; i >= 0; i-- {
                    list__append(&graph.adjs[X.elems[i]], v)
                }
                return
            }
            // updating the graph
            list__remove(&graph.adjs[u], v)
            // save, for later restore it
            list__append(&X, u)
        }
    }
    visited[v] = true
    // check that dst is reachable from all the vertices of X 
    if X.len == 0 || connected(&graph, X, dst) {
        // one step less is remained to dst
        l -= 1
        // go over all possible steps that can be done
        // i.e., check every unvisited adjacent vertex
        for i := Y.len-1; i >= 0; i-- {
            u := Y.elems[i]
            if l == 0 && u == dst {
                // we reached dst and passed all verticies
                count += 1
                //if (count % 1000 == 0) {
                //    printf("%d\n", count)
                //}
            }else if l > 0 && u != dst {
                // not yet reached dst, let's check the step
                search(u, l)
            }
        }
    }
    visited[v] = false

    // restore the graph
    for i := X.len-1; i >= 0; i-- {
        list__append(&graph.adjs[X.elems[i]], v)
    }
}

func main() {
    var (
        datacenter grid_t
    )

    read_input(&datacenter)

    build_graph(&datacenter, &graph, &src, &dst)
    // compute the number of verticies (the path length) we have to pass to reach dst
    length := 0
    for v := 0; v < graph.n; v++ {
        if graph.adjs[v].len > 0 {
            length += 1
        }
    }
    for v := 0; v < MAX_N; v++ {
        visited[v] = false
    }
    // the number of pathes
    count = 0
    search(src, length)
    fmt.Printf("%d\n", count)
}


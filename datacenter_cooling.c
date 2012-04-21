// Author: Chernoy Viacheslav
// Email: VCHERNOY@gmail.com
// Skype: vchernoy
//
// The program solves "Datacenter Cooling" Problem presented by Quora
// http://www.quora.com/challenges

#include <stdio.h>
#include <stdlib.h>

#define MAX_DIM (10)
#define MAX_N (MAX_DIM * MAX_DIM)

//#define ASSERT(c) if (c) {} else {printf("assert: "#c" in %d \n", __LINE__); *((char*)0) = 0;}
#define ASSERT(c)

#define ABS(x) ((x) >= 0 ? (x) : -(x))

typedef int bool_t;
#define TRUE  (1)
#define FALSE (0)

// the 0/1/2/3-matrix type
typedef struct{
    int h, w;
    int arr[MAX_DIM][MAX_DIM];
} grid_t;

typedef int vertex_t;

// the list type (maximum 4 (four) adjacent vertices exist
typedef struct{
    int len;
    vertex_t elems[4];
} list_t;

// the graph represented by adjacency list
typedef struct{
    int n;
    list_t adjs[MAX_N];
} graph_t;


// list sub-system
inline static void list__init(list_t* list){
    list->len = 0;
}

inline static void list__append(list_t* list, vertex_t x){
    ASSERT(list->len < 4);

    list->elems[list->len] = x;
    list->len += 1;
}

inline static vertex_t list__pop(list_t* list){
    vertex_t x;

    ASSERT(list->len > 0);

    x = list->elems[0];
    list->elems[0] = list->elems[list->len-1];
    list->len -= 1;

    return x;
}

static void list__remove(list_t* list, vertex_t x){
    register int i;

    i = list->len-1;
    while (i >= 0 && list->elems[i] != x){
        i -= 1;
    }
    ASSERT(i >= 0);
 
    list->elems[i] = list->elems[list->len-1];
    list->len -= 1;
}

static bool_t list__contains(const list_t* list, vertex_t x){
    register int i;

    i = list->len-1;
    while (i >= 0 && list->elems[i] != x){
        i -= 1;
    }

    return i >= 0;
}

// the room types
static const int OWN       = 0;
static const int START     = 2;
static const int END       = 3;
static const int DONOT_OWN = 1;
                    
// reads from the stdin 0/1/2/3-grid representing datacenter 
static void read_input(grid_t* datacenter){
    int i, j;
    int rc;

    rc = scanf("%d %d\n", &datacenter->w, &datacenter->h);
    ASSERT(rc == 2);
    for (i = 0; i < datacenter->h; i++){
        for (j = 0; j < datacenter->w; j++){
            rc = scanf("%d", &datacenter->arr[i][j]);
            ASSERT(rc == 1);
            if (rc != 1){
                datacenter->arr[i][j] = OWN;
            }
        }
    }
}

// converts the pair (i, j) to a vertex
inline static vertex_t vert(int w, int i, int j){
    return i * w + j;
}

// datacenter is the 0/1/2/3-grid
// graph is adjacency list contructed from the matrix datacenter
// builds the graph representation of datacenter
// also returns source and destination vertices
static void build_graph(const grid_t* datacenter, graph_t* graph, vertex_t* src, vertex_t* dst){
    // four directions up, down, left, right
    static const int Dh[] = {-1, 1, 0, 0};
    static const int Dw[] = {0, 0, -1, 1};

    int i, j, i1, j1, k;
    vertex_t v, u;
    list_t* adj;

    // create empty graph
    graph->n = datacenter->w * datacenter->h;
    for (i = 0; i < MAX_N; i++){
        list__init(&graph->adjs[i]);
    }

    for (i = 0; i < datacenter->h; i++){
        for (j = 0; j < datacenter->w; j++){
            // the current vertex v
            v = vert(datacenter->w, i, j);
            // adjacency list of v
            adj = &graph->adjs[v];
            if (datacenter->arr[i][j] == OWN || datacenter->arr[i][j] == START){
                // find all vertices adjacent to v
                for (k = 0; k < 4; k++){
                    i1 = i + Dh[k];
                    j1 = j + Dw[k];
                    if ((i1 >= 0) && (i1 < datacenter->h) && (j1 >= 0) && (j1 < datacenter->w) && (datacenter->arr[i1][j1] != DONOT_OWN)){
                        u = vert(datacenter->w, i1, j1);
                        list__append(adj, u);
                    }
                }
            }
            if (datacenter->arr[i][j] == START){
                *src = v;
            }else if (datacenter->arr[i][j] == END){
                *dst = v;
            }
        }
    }
}

// graph represented by adjacency list
// L is the list of vertices to check
// u is a destination vertex
// checks whether the vertex u is reachable from all the vertices of L
static bool_t connected(const graph_t* graph, list_t L, vertex_t u){
    // C[v] == 0 means it is not connected
    // for every k >= 2, {v|C[v] == k} are connected to v = L[1-k] (and between each other)
    int C[MAX_N] = {0};
    vertex_t v, x, y;
    register int i;
    register int mark;
    // R is a list of vertices reachable from v
    struct {
        int len;
        vertex_t elems[MAX_N];
    } R;

    C[u] = 1;
    mark = 2;
    while (L.len > 0){
        v = list__pop(&L);
        // check whether v can reach u
        if (C[v] == 0){
            R.elems[0] = v;
            R.len = 1;
            while (TRUE){
                if (R.len == 0){
                    // v cannot reach u
                    return FALSE;
                }

                R.len -= 1;
                x = R.elems[R.len];
                // x is reachable from v
                // check its adjacent ones
                for (i = graph->adjs[x].len-1; i >= 0; i--){
                    y = graph->adjs[x].elems[i];
                    if (C[y] == 0){
                        // y is reachable from v
                        C[y] = mark;
                        ASSERT(R.len < MAX_N);
                        R.elems[R.len] = y;
                        R.len += 1;
                    }else if (C[y] < mark){
                        // u is reachable from y,
                        // so we are done with v
                        break;
                    }
                }
                if (i >= 0){
                    break;
                }
            }
            mark += 1;
        }
    }
    return TRUE;
}

static int count;
static bool_t visited[MAX_N] = {FALSE};
static graph_t graph;
static vertex_t src, dst;

// v is the current vertex
// l is the number of steps remained to do
// backtracking, computes the number of possible paths of length l from v to dst
static void search(vertex_t v, int l){
    // global dst, count, visited, graph
    list_t X, Y;
    register int i;
    register vertex_t u;

    list__init(&X);
    list__init(&Y);
    for (i = graph.adjs[v].len-1; i >= 0; i--){
        u = graph.adjs[v].elems[i];
        if (!visited[u]){
            list__append(&Y, u);
        }
    }
    // remove all the edges to v
    for (i = Y.len-1; i >= 0; i--){
        u = Y.elems[i];
        if (list__contains(&graph.adjs[u], v)){
            if (graph.adjs[u].len <= 1){
                // u is isolated, rollback
                for (i = X.len-1; i >= 0; i--){
                    list__append(&graph.adjs[X.elems[i]], v);
                }
                return;
            }
            // updating the graph
            list__remove(&graph.adjs[u], v);
            // save, for restoring it later 
            list__append(&X, u);
        }
    }
    visited[v] = TRUE;
    // check that dst is reachable from all the vertices of X 
    if (X.len == 0 || connected(&graph, X, dst)){
        // one step less is remained to dst
        l -= 1;
        // go over all possible steps that can be done
        // i.e., check every unvisited adjacent vertex
        for (i = Y.len-1; i >= 0; i--){
            u = Y.elems[i];
            if (l > 0 && u != dst){
                // not yet reached dst, let's check the step
                search(u, l);
            } else if (l == 0 && u == dst){
                // we reached dst and passed all vertices
                count += 1;
                //if (count % 1000 == 0){
                //    printf("%d\n", count);
                //}
            }
        }
    }
    visited[v] = FALSE;

    // restore the graph
    for (i = X.len-1; i >= 0; i--){
        list__append(&graph.adjs[X.elems[i]], v);
    }
}

int main(){
    grid_t datacenter;
    int length;
    vertex_t v;

    read_input(&datacenter);
    build_graph(&datacenter, &graph, &src, &dst);

    // compute the number of vertices (the path length) we have to pass to reach dst
    length = 0;
    for (v = 0; v < graph.n; v++){
        if (graph.adjs[v].len > 0){
            length += 1;
        }
    }
    // the number of pathes
    count = 0;
    search(src, length);
    printf("%d\n", count);

    return 0;
}


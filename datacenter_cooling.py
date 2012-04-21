# Author: Chernoy Viacheslav
# email: vchenoy@gmail.com
#
# The program solves "Datacenter Cooling" Problem presented by Quora
# http://www.quora.com/challenges

# reads and returns 0/1/2/3 matrix from the standard-input
def read_input():
    datacenter = []
    w, h = [int(x) for x in raw_input().split()]
    for i in xrange(h):
        l = [int(x) for x in raw_input().split()]
        assert(len(l) == w)

        datacenter.append(l)

    return datacenter

# datacenter is 0/1/2/3-matrix
# creates and returns the graph representation of datacenter
# also returns source and destination verticies
def build_graph(datacenter):
    # converts the pair (i, j) to a vertex
    def vert(i, j):
        return i * w + j
    
    OWN       = 0
    START     = 2
    END       = 3
    DONOT_OWN = 1

    h = len(datacenter)
    w = len(datacenter[0])
    # the number of verticies
    n = w * h
    # create empty graph
    graph = [[] for k in xrange(n)]
    # four directions up, down, left, right
    D = [(-1,0), (1,0), (0,-1), (0,1)]
    for i in xrange(h):
        for j in xrange(w):
            # the current vertex v
            v = vert(i, j)
            # adjacency list of v
            adj = []
            if datacenter[i][j] in [OWN, START]:
                # find all verticies adjacent to v
                for d in D:
                    i1 = i + d[0]
                    j1 = j + d[1]
                    if (i1 >= 0) and (i1 < h) and (j1 >= 0) and (j1 < w) and (datacenter[i1][j1] != DONOT_OWN):
                        u = vert(i1, j1)
                        adj.append(u)

            if datacenter[i][j] == START:
                src = v
            elif datacenter[i][j] == END:
                dst = v 
         
            graph[v] = adj

    return src, dst, graph

# graph is a graph
# L is the list of vertices to check
# u is a destination vertex
# checks whether the vertex u is reachable from all the vertices of L
def connected(graph, L, u):
    n = len(graph)
    L = L[:]
    # C[v] == 0 means it is not connected
    # for every k >= 2, {v|C[v] == k} are connected to v = L[1-k] (and between each other)
    C = [0] * n
    C[u] = 1
    mark = 2

    while len(L) > 0:
        v = L.pop()
        # check whether v can reach u
        if C[v] == 0:
            # R is a list of verticies reachable form v
            R = [v]
            while True:
                if len(R) == 0:
                    # v cannot reach u
                    return False

                # x is reachable from v
                x = R.pop()
                # check its adjacent ones
                for y in graph[x]:
                    if C[y] == 0:
                        # y is reachable from v
                        C[y] = mark 
                        R.append(y)
                    elif C[y] < mark:
                        # u is reachable from y,
                        # so we are done with v
                        break

                if C[y] < mark:
                    break

            mark += 1

    return True    

# v is the current vertex
# l is the number of steps remained to do
# backtracking, computes the number of possible paths of length l from v to dst
def search(v, l):
    global dst, count, visited, graph

    Y = [u for u in graph[v] if not visited[u]]
    X = []
    # remove all the edges to v
    for u in Y:
        if v in graph[u]:
            if len(graph[u]) <= 1:
                # u is isolated, rollback
                for w in X:
                    graph[w].append(v)

                return

            # updating the graph
            graph[u].remove(v)
            # save, for later restore it
            X.append(u)

    visited[v] = True
    # check that dst is reachable from all the vertices of X 
    if (len(X) == 0) or connected(graph, X, dst):
        # one step less is remained to dst
        l -= 1
        # go over all possible steps that can be done
        # i.e., check every unvisited adjacent vertex
        for u in Y:
            if (l == 0) and (u == dst):
                # we reached dst and passed all verticies
                count += 1
                #if count % 1000 == 0:
                #    print count
            elif (l > 0) and (u != dst):
                # not yet reached dst, let's check the step
                search(u, l)

    visited[v] = False

    # restore the graph
    for w in X:
        graph[w].append(v)



datacenter = read_input()
src, dst, graph = build_graph(datacenter)
n = len(datacenter) * len(datacenter[0])
# compute the number of verticies (the path length) we have to pass to reach dst
length = len([v for v in xrange(n) if len(graph[v]) > 0])
visited = [False] * n
# the number of pathes
count = 0
search(src, length)
print count


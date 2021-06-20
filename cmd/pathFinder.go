package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func newDirections() ([]int, []int) {
	var dx = []int{1, -1, 0, 0}
	var dy = []int{0, 0, 1, -1}

	return dx, dy
}

type path struct {
	Grid       [][]string
	n, m       int
	start, end pair
}

func newPath(arr [][]string, ptr *path) {
	tmp := path{}
	tmp.Grid = arr
	tmp.n = len(arr)
	tmp.m = len(arr[0])

	tmp1 := pair{-1, -1}
	tmp2 := pair{-1, -1}

	for i := 0; i < tmp.n; i++ {
		for j := 0; j < tmp.m; j++ {
			if arr[i][j] == "." {
				if tmp1.x == -1 {
					tmp1 = pair{i, j}
				} else {
					tmp2 = pair{i, j}
				}
			}
		}
	}

	tmp.start = tmp1
	tmp.end = tmp2

	ptr = &tmp
}

type pair struct {
	x, y int
}

func (oth *path) bfs() int {
	var queue []pair

	cost := make([][]int, oth.n)
	for i := range cost {
		cost[i] = make([]int, oth.m)
	}

	for i := 0; i < oth.n; i++ {
		for j := 0; j < oth.m; j++ {
			cost[i][j] = 0
		}
	}

	dx, dy := newDirections()
	queue = append(queue, oth.start)

	for len(queue) > 0 {
		curx, cury := queue[0].x, queue[0].y
		queue = queue[1:]

		for i := 0; i < 4; i++ {
			if oth.valid(curx+dx[i], cury+dy[i]) && cost[curx][cury] != 0 {
				queue = append(queue, pair{curx + dx[i], cury + dy[i]})
				cost[curx+dx[i]][cury+dy[i]] = cost[curx][cury] + 1
			}
		}

		if curx == oth.end.x && cury == oth.end.y {
			break
		}
	}

	fmt.Println(cost[oth.end.x][oth.end.y])
	return cost[oth.end.x][oth.end.y]
}

func (oth *path) valid(x, y int) bool {
	if x >= 0 && y >= 0 && x < oth.n && y < oth.m && oth.Grid[x][y] != "#" {
		return true
	} else {
		return false
	}
}

type pathFinder struct {
	Path       *path
	BestValues []int
}

func (oth *pathFinder) HandleReq(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		oth.get(w, r)
		return
	case "POST":
		oth.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
}

func newPathFinder() pathFinder {
	return pathFinder{}
}

type output struct {
	Answer []int `json:"ans"`
}

func (oth *pathFinder) get(w http.ResponseWriter, r *http.Request) {
	res := output{oth.BestValues}
	if res.Answer == nil {
		jsonBytes, _ := json.Marshal(fmt.Sprintf("%d", 0))
		w.Write(jsonBytes)
	} else {
		jsonBytes, _ := json.Marshal(res)
		w.Write(jsonBytes)
	}
	w.Header().Add("content-type", "aplication/json")
	w.WriteHeader(http.StatusOK)
}

type myGrid struct {
	Array [][]string `json:"grid"`
}

func (oth *pathFinder) post(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var grid myGrid

	if err := decoder.Decode(&grid); err != nil {
		panic(err)
	}

	fmt.Println(grid.Array)
	newPath(grid.Array, oth.Path)
	oth.BestValues = append(oth.BestValues, oth.Path.bfs())
}

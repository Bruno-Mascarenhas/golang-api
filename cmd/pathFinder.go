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

func newPath() path {
	tmp := path{}

	return tmp
}

type pair struct {
	x, y int
}

func (this *path) bfs() int {
	var queue []pair

	cost := make([][]int, this.n)
	for i := range cost {
		cost[i] = make([]int, this.m)
	}

	for i := 0; i < this.n; i++ {
		for j := 0; j < this.m; j++ {
			cost[i][j] = 0
		}
	}

	dx, dy := newDirections()
	queue = append(queue, this.start)

	for len(queue) > 0 {
		curx, cury := queue[0].x, queue[0].y
		queue = queue[1:]

		for i := 0; i < 4; i++ {
			if this.valid(curx+dx[i], cury+dy[i]) && cost[curx][cury] != 0 {
				queue = append(queue, pair{curx + dx[i], cury + dy[i]})
				cost[curx+dx[i]][cury+dy[i]] = cost[curx][cury] + 1
			}
		}

		if curx == this.end.x && cury == this.end.y {
			break
		}
	}

	return cost[this.end.x][this.end.y]
}

func (this *path) valid(x, y int) bool {
	if x >= 0 && y >= 0 && x < this.n && y < this.m && this.Grid[x][y] != "#" {
		return true
	} else {
		return false
	}
}

type pathFinder struct {
	Path       *path
	BestValues []float64
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

type output struct {
	Answer []float64 `json:"ans"`
}

func (oth *pathFinder) get(w http.ResponseWriter, r *http.Request) {
	res := output{oth.BestValues}
	jsonBytes, _ := json.Marshal(res)
	w.Write(jsonBytes)
	w.Header().Add("content-type", "aplication/json")
	w.WriteHeader(http.StatusOK)
}

type myGrid struct {
	Array [][]string `json:"grid"`
}

func (oth *pathFinder) post(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Body)
	/*
		num, _ := strconv.Atoi("0")

		jsonExample := `{"grid":[["#","*","*","*","*","."],
								  ["#","*","*","*","*","."],
								  ["#","*","*","*","*","."],
								  ["#","*","*","*","*","."],
								  ["#","*","*","*","*","."]]
						}`

		var test myGrid

		json.Unmarshal([]byte(jsonExample), &test)
	*/
}

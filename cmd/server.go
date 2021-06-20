package main

import (
	"net/http"
)

func main() {
	rot := newRot()
	http.HandleFunc("/crypt/rot", rot.handle)

	chatServer := newChatServer()
	chatServer.handleRequests()

	pathFinder := newPathFinder()
	http.HandleFunc("/maze", pathFinder.HandleReq)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

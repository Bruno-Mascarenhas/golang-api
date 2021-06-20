package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	/*
	coasterHandlers := newCoasterHandlers()  //cada um faz o seu handler
	http.HandleFunc("/coasters", coasterHandlers.coasters)    // ai coloca as rotas do seu handler
	http.HandleFunc("/coasters/", coasterHandlers.getCoaster) // se for usar get / post ai cria tudo 
															  // no arquivo separado e so add aqui
	*/
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}


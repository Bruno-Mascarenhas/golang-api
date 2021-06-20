package main

import (
	"fmt"
	"strings"
	"math"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type Rot struct {
	Route	float64	`json:"route"`
	Message	string	`json:"message"`
}

func  newRot() Rot {
	return Rot{}
}

func (u Rot) handle(w http.ResponseWriter, re *http.Request){
	reqBody, _ := ioutil.ReadAll(re.Body)

	var rot Rot
	json.Unmarshal(reqBody, &rot)

	if  rot.Route < -26 || rot.Route > 26 {

		fmt.Fprintf(w, "Invalid arguments!")

	} else {
		var route rune

		if rot.Route < 0 { route = rune(26-math.Abs(rot.Route))
		} else { route = rune(rot.Route) }

		shift := func(c rune) rune {
			if c >= 'A' && c <= 'Z'{
				return 'A' + rune( (c-'A'+route)%26 )
			} else if c >= 'a' && c <= 'z' {
				return 'a' + rune( (c-'a'+route)%26 )
			}

			return c

		}

		fmt.Fprintf(w, strings.Map(shift, rot.Message))
	}
}

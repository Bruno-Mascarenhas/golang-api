package main

import (
	"fmt"
	"strings"
	"math"
	"net/http"
	"strconv"
)

type rot struct {
	route int
	message string
}

func  newRot() rot {
	return rot{}
}

func (u rot) handle(w http.ResponseWriter, re *http.Request){
	query := strings.Split(re.URL.String(), "?")

	if len(query)!=2 {

		fmt.Fprintf(w, "Invalid query!")

	} else {

		args := strings.Split(query[1], "&")

		var rs string = strings.Split(args[0],"=")[1]
		var r, _ = strconv.ParseFloat(rs, 64)
		var message string = strings.Split(args[1],"=")[1]

		if  r < -26 || r > 26 || len(args)<2 {
			
			fmt.Fprintf(w, "Invalid arguments!")

		} else {
			var route rune
			
			if r < 0 { route = rune(26-math.Abs(r)) } else { route = rune(r) }

			shift := func(c rune) rune {
				if c >= 'A' && c <= 'Z'{
					return 'A' + rune( (c-'A'+route)%26 )
				} else if c >= 'a' && c <= 'z' {
					return 'a' + rune( (c-'a'+route)%26 )
				}

				return c
						
			}

			fmt.Fprintf(w, strings.Map(shift, message))
		}
	}
}

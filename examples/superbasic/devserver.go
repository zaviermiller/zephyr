package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Serving at http://localhost:9000")
	http.ListenAndServe(`:9000`, http.FileServer(http.Dir(`.`)))
}

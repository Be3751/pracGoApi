package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(rw, "Hello, world!")
	})

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

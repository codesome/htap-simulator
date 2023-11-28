package main

import (
	"fmt"
	"os"
)

func main() {
	w, err := NewWriter()
	handleErr(err)
	q, err := NewQuerier()
	handleErr(err)

	server := NewServer(w, q)
	err = server.ListenAndServe()
	handleErr(err)
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

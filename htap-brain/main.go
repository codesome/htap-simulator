package main

import (
	"fmt"
	"os"
)

func main() {
	os.RemoveAll("wal")

	h, err := NewHTAPBrain()
	handleErr(err)
	server := NewServer(h)
	handleErr(server.ListenAndServe())
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

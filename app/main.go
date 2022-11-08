package main

import "fmt"

func main() {
	port := ":6379"
	fmt.Printf("Redis sever is running on %s port", port)
	ListenAndServe(port)
}

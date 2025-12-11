// Package main provides a compatibility entry point for the Discovery Tree REST API
// 
// DEPRECATED: This file is kept for backward compatibility.
// The main application entry point has been moved to cmd/api/main.go
// 
// To run the API server, use:
//   go run cmd/api/main.go
//
// Or build and run:
//   go build -o discovery-tree-api cmd/api/main.go
//   ./discovery-tree-api
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("DEPRECATED: This entry point has been moved to cmd/api/main.go")
	fmt.Println("")
	fmt.Println("To run the Discovery Tree API server, use:")
	fmt.Println("  go run cmd/api/main.go")
	fmt.Println("")
	fmt.Println("Or build and run:")
	fmt.Println("  go build -o discovery-tree-api cmd/api/main.go")
	fmt.Println("  ./discovery-tree-api")
	fmt.Println("")
	fmt.Println("For configuration options, see the README.md file.")
	
	os.Exit(1)
}
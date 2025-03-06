package main

import (
	"fmt"
	"os"

	"repeater/internal/request"
	"repeater/internal/response"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <request_file>")
		return
	}

	fileName := os.Args[1]

	// Validate the file
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		fmt.Println("Error: File does not exist")
		return
	}

	// Read the raw request from the file
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	defer f.Close()

	err := request.ReadHTTP2FromFile(f)
	if err != nil {
		fmt.Println("Error reading HTTP from file:", err)
		return
	}

	response.HandleResponse(conn[0].Response)
}

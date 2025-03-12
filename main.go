package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	// For example, test the "create_network" action.
	parameters := map[string]interface{}{
		"action": "create_network",
		"name":   "test-network",
	}

	result, err := Handler(context.Background(), parameters)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Handler error: %v\n", err)
		os.Exit(1)
	}

	resJSON, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println("Result:", string(resJSON))
}

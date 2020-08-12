package main

import (
	"fmt"

	"github.com/cloudboltsoftware/terraform-provider-onefuse/onefuse"
)

func main() {
	config := onefuse.NewConfig("http", "localhost", "8000", "admin", "admin", true)
	client := config.NewOneFuseApiClient()

	// Add your stuff here

	fmt.Printf("Complete!")
}

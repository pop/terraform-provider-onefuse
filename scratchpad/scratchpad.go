package main

import (
	"fmt"
	"path"

	"github.com/cloudboltsoftware/terraform-provider-onefuse/onefuse"
)

func main() {
	config := onefuse.NewConfig("http", "localhost", "8000", "admin", "admin", true)
	_ = config.NewOneFuseApiClient()

	fmt.Printf("%s/%s/", "http://foo.com:123", path.Join("bar", "baz"))

	// Add your stuff here
	fmt.Printf("Complete!")
}

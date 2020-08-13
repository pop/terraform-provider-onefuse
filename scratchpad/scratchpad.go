package main

import (
	"fmt"

	"github.com/cloudboltsoftware/terraform-provider-onefuse/onefuse"
)

func main() {
	config := onefuse.NewConfig("http", "localhost", "8000", "admin", "admin", true)
	client := config.NewOneFuseApiClient()

	// Add your stuff here
    _, err := client.UpdateMicrosoftAdPolicy(6, &onefuse.MicrosoftAdPolicy{ComputerNameLetterCase: "UPPER", OU: "OU=NotOneFuse,DC=sovlabs,DC=net", Name: "NewName", Workspace: "/api/v3/onefuse/workspaces/1/"})
    if err != nil {
        fmt.Println(err)
    }

	fmt.Printf("Complete!")
}

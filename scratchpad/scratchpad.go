package main

import (
	"fmt"

	"github.com/cloudboltsoftware/terraform-provider-onefuse/onefuse"
)

func main() {
	config := onefuse.NewConfig("http", "localhost", "8000", "admin", "admin", true)
	client := config.NewOneFuseApiClient()

    desiredPolicy := onefuse.MicrosoftADPolicy{
        ComputerNameLetterCase: "UPPER",
        OU: "OU=NotOneFuse,DC=sovlabs,DC=net",
        Name: "NewName4",
        WorkspaceURL: "/api/v3/onefuse/workspaces/1/",
        MicrosoftEndpointID: 13,
    }

	// Add your stuff here
	policy, err := client.CreateMicrosoftADPolicy(&desiredPolicy)
    if err != nil {
        fmt.Println("Create failed:", err)
        return
    } else {
        fmt.Println("Create policy:", policy)
    }
    policy, err = client.GetMicrosoftADPolicy(policy.ID)
    if err != nil {
        fmt.Println("Get failed:", err)
        return
    } else {
        fmt.Println("Get policy:", policy)
    }
    policy, err = client.UpdateMicrosoftADPolicy(policy.ID, &desiredPolicy) // TODO: Update only works if name changes
    if err != nil {
        fmt.Println("Update failed:", err)
        return
    } else {
        fmt.Println("Update policy:", policy)
    }
    err = client.DeleteMicrosoftADPolicy(policy.ID)
	if err != nil {
        fmt.Println("Delete failed:", err)
        return
	}
	fmt.Printf("Complete!")
}

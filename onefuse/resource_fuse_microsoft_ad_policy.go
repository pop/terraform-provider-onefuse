// Copyright 2020 CloudBolt Software
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package onefuse

import (
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func resourceMicrosoftADPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMicrosoftADPolicyCreate,
		Read:   resourceMicrosoftADPolicyRead,
		Update: resourceMicrosoftADPolicyUpdate,
		Delete: resourceMicrosoftADPolicyDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"microsoft_endpoint_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"computer_name_letter_case": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Either Lowercase or Uppercase",
			},
			"ou": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workspace_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func bindMicrosoftADPolicyResource(d *schema.ResourceData, policy *MicrosoftADPolicy) error {
	d.SetId(strconv.Itoa(policy.ID))

	if err := d.Set("name", policy.Name); err != nil {
		return errors.WithMessage(err, "cannot set name")
	}
	if err := d.Set("description", policy.Description); err != nil {
		return errors.WithMessage(err, "cannot set description")
	}
	if err := d.Set("workspace_url", policy.Links.Workspace.Href); err != nil {
		return errors.WithMessage(err, "cannot set workspace")
	}
	if err := d.Set("computer_name_letter_case", policy.ComputerNameLetterCase); err != nil {
		return errors.WithMessage(err, "cannot set computer_name_letter_case")
	}
	if err := d.Set("ou", policy.OU); err != nil {
		return errors.WithMessage(err, "cannot set OU")
	}

	log.Println("[!!] Policy Links: ", policy.Links)
	log.Println("[!!] Policy Links MicrosoftEndpoint: ", policy.Links.MicrosoftEndpoint)
	microsoftEndpointURLSplit := strings.Split(policy.Links.MicrosoftEndpoint.Href, "/")
	log.Println("[!!] splitURL: ", microsoftEndpointURLSplit)
	microsoftEndpointID := microsoftEndpointURLSplit[len(microsoftEndpointURLSplit)-2]
	microsoftEndpointIDInt, err := strconv.Atoi(microsoftEndpointID)
	if err != nil {
		return errors.WithMessage(err, "Failed to convert string to int...")
	}
	log.Println("[!!] Endpoint ID: ", microsoftEndpointIDInt)
	if err := d.Set("microsoft_endpoint_id", microsoftEndpointIDInt); err != nil {
		return errors.WithMessage(err, "cannot set microsoft_endpoint_id")
	}
	// TODO: set OU and letter_case, too.

	return nil
}

func resourceMicrosoftADPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("[!!] calling resourceMicrosoftADPolicyCreate")

	newPolicy := MicrosoftADPolicy{
		Name:                   d.Get("name").(string),
		Description:            d.Get("description").(string),
		OU:                     d.Get("ou").(string),
		MicrosoftEndpointID:    d.Get("microsoft_endpoint_id").(int),
		ComputerNameLetterCase: d.Get("computer_name_letter_case").(string),
		WorkspaceURL:           d.Get("workspace_url").(string),
	}
	log.Println("[!!] newPolicy: ", newPolicy)

	config := m.(Config)
	policy, err := config.NewOneFuseApiClient().CreateMicrosoftADPolicy(&newPolicy)
	if err != nil {
		return err
	}
	log.Println("[!!] receivedPolicy: ", policy)

	err = bindMicrosoftADPolicyResource(d, &policy)
	return err
}

func resourceMicrosoftADPolicyRead(d *schema.ResourceData, m interface{}) error {
	config := m.(Config)
	id := d.Id()
	int_id, _ := strconv.Atoi(id)
	policy, err := config.NewOneFuseApiClient().GetMicrosoftADPolicy(int_id)
	bindMicrosoftADPolicyResource(d, &policy)
	return err
}

func resourceMicrosoftADPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	// Determine if a change is needed
	changed := (d.HasChange("name") ||
		d.HasChange("description") ||
		d.HasChange("microsoft_endpoint_id") ||
		d.HasChange("computer_name_letter_case") ||
		d.HasChange("workspace_url") ||
		d.HasChange("ou"))
	if !changed {
		return nil
	}

	log.Printf("[!!] stuff has changed")

	// Create the desired AD Policy object
	id := d.Id()
	desiredPolicy := MicrosoftADPolicy{
		Name:                   d.Get("name").(string),
		Description:            d.Get("description").(string),
		MicrosoftEndpointID:    d.Get("microsoft_endpoint_id").(int),
		ComputerNameLetterCase: d.Get("computer_name_letter_case").(string),
		WorkspaceURL:           d.Get("workspace_url").(string),
		OU:                     d.Get("ou").(string),
	}

	log.Println("[!!]", id, desiredPolicy)

	// Make the API call to update the policy
	config := m.(Config)
	int_id, _ := strconv.Atoi(id)
	updatedPolicy, err := config.NewOneFuseApiClient().UpdateMicrosoftADPolicy(int_id, &desiredPolicy)
	if err != nil {
		return err
	}

	log.Println("[!!]", updatedPolicy, err)

	// Update Terraform's state
	d.Set("name", updatedPolicy.Name)
	d.Set("description", updatedPolicy.Description)
	d.Set("microsoft_endpoint_id", strconv.Itoa(updatedPolicy.MicrosoftEndpointID))
	d.Set("computer_name_letter_case", updatedPolicy.ComputerNameLetterCase)
	d.Set("workspace_url", updatedPolicy.WorkspaceURL)
	d.Set("ou", updatedPolicy.OU)

	return nil
}

func resourceMicrosoftADPolicyDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(Config)
	id := d.Id()
	int_id, _ := strconv.Atoi(id)
	return config.NewOneFuseApiClient().DeleteMicrosoftADPolicy(int_id)
}

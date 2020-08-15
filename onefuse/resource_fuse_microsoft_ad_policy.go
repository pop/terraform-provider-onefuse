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
				Computed: true,
				Optional: true,
			},
		},
	}
}

func bindMicrosoftADPolicyResource(d *schema.ResourceData, policy *MicrosoftADPolicy) error {
	log.Println("onefuse.bindMicrosoftADPolicyResource")

	if err := d.Set("name", policy.Name); err != nil {
		return errors.WithMessage(err, "Cannot set name: "+policy.Name)
	}

	if err := d.Set("description", policy.Description); err != nil {
		return errors.WithMessage(err, "Cannot set description: "+policy.Description)
	}

	log.Printf("[!!] Policy workspace URL: %#v\n", policy.Links.Workspace.Href)
	if err := d.Set("workspace_url", policy.Links.Workspace.Href); err != nil {
		return errors.WithMessage(err, "Cannot set workspace: "+policy.Links.Workspace.Href)
	}
	log.Printf("[!!] Policy as set on URL: %#v\n", d.Get("workspace_url"))

	if err := d.Set("computer_name_letter_case", policy.ComputerNameLetterCase); err != nil {
		return errors.WithMessage(err, "Cannot set computer_name_letter_case: "+policy.ComputerNameLetterCase)
	}

	if err := d.Set("ou", policy.OU); err != nil {
		return errors.WithMessage(err, "Cannot set OU: "+policy.OU)
	}

	microsoftEndpointURLSplit := strings.Split(policy.Links.MicrosoftEndpoint.Href, "/")
	microsoftEndpointID := microsoftEndpointURLSplit[len(microsoftEndpointURLSplit)-2]
	microsoftEndpointIDInt, err := strconv.Atoi(microsoftEndpointID)
	if err != nil {
		return errors.WithMessage(err, "Expected to convert "+microsoftEndpointID+" to int value.")
	}
	if err := d.Set("microsoft_endpoint_id", microsoftEndpointIDInt); err != nil {
		return errors.WithMessage(err, "Cannot set microsoft_endpoint_id")
	}

	return nil
}

func resourceMicrosoftADPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceMicrosoftADPolicyCreate")

	config := m.(Config)

	newPolicy := MicrosoftADPolicy{
		Name:                   d.Get("name").(string),
		Description:            d.Get("description").(string),
		OU:                     d.Get("ou").(string),
		MicrosoftEndpointID:    d.Get("microsoft_endpoint_id").(int),
		ComputerNameLetterCase: d.Get("computer_name_letter_case").(string),
		WorkspaceURL:           d.Get("workspace_url").(string),
	}

	policy, err := config.NewOneFuseApiClient().CreateMicrosoftADPolicy(&newPolicy)
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(policy.ID))

	return resourceMicrosoftADPolicyRead(d, m)
}

func resourceMicrosoftADPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceMicrosoftADPolicyRead")

	config := m.(Config)

	id := d.Id()
	int_id, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	policy, err := config.NewOneFuseApiClient().GetMicrosoftADPolicy(int_id)
	if err != nil {
		return err
	}

	log.Printf("[!!] Read Policy: %#v\n", policy)
	log.Printf("[!!] Read Policy Worksapce: %#v\n", policy.Links.Workspace)

	return bindMicrosoftADPolicyResource(d, policy)
}

func resourceMicrosoftADPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceMicrosoftADPolicyUpdate")

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

	// Make the API call to update the policy
	config := m.(Config)

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

	int_id, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	_, err = config.NewOneFuseApiClient().UpdateMicrosoftADPolicy(int_id, &desiredPolicy)
	if err != nil {
		return err
	}

	return resourceMicrosoftADPolicyRead(d, m)
}

func resourceMicrosoftADPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceMicrosoftADPolicyDelete")

	config := m.(Config)

	id := d.Id()
	int_id, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	return config.NewOneFuseApiClient().DeleteMicrosoftADPolicy(int_id)
}

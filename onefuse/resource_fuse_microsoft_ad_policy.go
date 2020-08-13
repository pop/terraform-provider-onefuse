// Copyright 2020 CloudBolt Software
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package onefuse

import (
	// "log" // TODO: Un-comment when data source implemented

	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func resourceMicrosoftAdPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMicrosoftAdPolicyCreate,
		Read:   resourceMicrosoftAdPolicyRead,
		Update: resourceMicrosoftAdPolicyUpdate,
		Delete: resourceMicrosoftAdPolicyDelete,
		Schema: map[string]*schema.Schema{
			"microsoft_ad_policy_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"microsoft_endpoint": {
				Type:     schema.TypeString,
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
			"workspace_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func bindMicrosoftADPolicyResource(d *schema.ResourceData, policy *MicrosoftAdPolicy) error {
	d.SetId(fmt.Sprintf("%v", policy.ID))

	if err := d.Set("name", policy.Name); err != nil {
		return errors.WithMessage(err, "cannot set name")
	}
	if err := d.Set("description", policy.Description); err != nil {
		return errors.WithMessage(err, "cannot set description")
	}
	if err := d.Set("workspace", policy.Links.Workspace); err != nil {
		return errors.WithMessage(err, "cannot set workspace")
	}
	if err := d.Set("microsoft_endpoint", policy.Links.MicrosoftEndpoint); err != nil {
		return errors.WithMessage(err, "cannot set microsoft_endpoint")
	}

	// TODO: set OU and letter_case, too.

	return nil
}

func resourceMicrosoftAdPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("calling resourceMicrosoftAdPolicyCreate")

	// TODO: parse out the id and/or name from workspaceUrl
	//  and will have to update when OneFuse supports multiple workspaces.
	workspaceURL := d.Get("workspace").(string)
	workspace := Workspace{
		Name: workspaceURL,
		ID:   workspaceURL,
	}

	newPolicy := MicrosoftAdPolicy{
		Name:                   d.Get("name").(string),
		Description:            d.Get("description").(string),
		OU:                     d.Get("ou").(string),
		MicrosoftEndpoint:      d.Get("microsoftEndpointId").(string),
		ComputerNameLetterCase: d.Get("computerNameLetterCase").(string),
	}
	newPolicy.Links.Workspace = workspace
	config := m.(Config)
	policy, err := config.NewOneFuseApiClient().CreateMicrosoftAdPolicy(newPolicy)
	if err != nil {
		return err
	}
	err = bindMicrosoftADPolicyResource(d, &policy)
	return err
}

func resourceMicrosoftAdPolicyRead(d *schema.ResourceData, m interface{}) error {
	config := m.(Config)
	id := d.Get("microsoft_ad_policy_id").(int)
	policy, err := config.NewOneFuseApiClient().GetMicrosoftAdPolicy(id)
	bindMicrosoftADPolicyResource(d, &policy)
	return err
}

func resourceMicrosoftAdPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	return errors.New("Not implemented yet")
}

func resourceMicrosoftAdPolicyDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(Config)
	id := d.Get("microsoft_ad_policy_id").(int)
	return config.NewOneFuseApiClient().DeleteMicrosoftAdPolicy(id)
}

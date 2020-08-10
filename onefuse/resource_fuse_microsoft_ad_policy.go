// Copyright 2020 CloudBolt Software
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package onefuse

import (
	"errors" // DELETE ME once func's are implemented
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

func resourceMicrosoftAdPolicyCreate(d *schema.ResourceData, m interface{}) error {
	return errors.New("Not implemented yet")
}

func resourceMicrosoftAdPolicyRead(d *schema.ResourceData, m interface{}) error {
	return errors.New("Not implemented yet")
}

func resourceMicrosoftAdPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	return errors.New("Not implemented yet")
}

func resourceMicrosoftAdPolicyDelete(d *schema.ResourceData, m interface{}) error {
	return errors.New("Not implemented yet")
}

// Copyright 2020 CloudBolt Software
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package onefuse

import (
	// "log" // TODO: Un-comment when data source implemented

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func dataSourceMicrosoftEndpoint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMicrosoftEndpointRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMicrosoftEndpointRead(d *schema.ResourceData, meta interface{}) error {
	return errors.New("Not implemented yet")
}

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNsxtVpcServiceProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtVpcServiceProfileRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"context":      getContextSchemaExtended(true, false, false, true),
		},
	}
}

func dataSourceNsxtVpcServiceProfileRead(d *schema.ResourceData, m interface{}) error {
	_, err := policyDataSourceResourceRead(d, getPolicyConnector(m), getSessionContext(d, m), "VpcServiceProfile", nil)
	if err != nil {
		return err
	}
	return nil
}

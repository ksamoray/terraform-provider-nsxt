// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNsxtPolicySpoofGuardProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicySpoofGuardProfileRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"context":      getContextSchema(false, false, false),
		},
	}
}

func dataSourceNsxtPolicySpoofGuardProfileRead(d *schema.ResourceData, m interface{}) error {
	_, err := policyDataSourceResourceRead(d, getPolicyConnector(m), getSessionContext(d, m), "SpoofGuardProfile", nil)
	if err != nil {
		return err
	}
	return nil
}

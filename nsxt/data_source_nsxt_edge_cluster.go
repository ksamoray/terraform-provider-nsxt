// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/go-vmware-nsxt/manager"
)

func dataSourceNsxtEdgeCluster() *schema.Resource {
	return &schema.Resource{
		Read:               dataSourceNsxtEdgeClusterRead,
		DeprecationMessage: mpObjectDataSourceDeprecationMessage,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Unique ID of this resource",
				Optional:    true,
				Computed:    true,
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "The display name of this resource",
				Optional:    true,
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
				Computed:    true,
			},
			"deployment_type": {
				Type:        schema.TypeString,
				Description: "The deployment type of edge cluster members (UNKNOWN/VIRTUAL_MACHINE|PHYSICAL_MACHINE)",
				Optional:    true,
				Computed:    true,
			},
			"member_node_type": {
				Type:        schema.TypeString,
				Description: "Type of transport nodes",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtEdgeClusterRead(d *schema.ResourceData, m interface{}) error {
	// Read an edge cluster by name or id
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return dataSourceNotSupportedError()
	}

	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)
	var obj manager.EdgeCluster
	if objID != "" {
		// Get by id
		objGet, resp, err := nsxClient.NetworkTransportApi.ReadEdgeCluster(nsxClient.Context, objID)

		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Edge cluster %s was not found", objID)
		}
		if err != nil {
			return fmt.Errorf("Error while reading edge cluster %s: %v", objID, err)
		}
		obj = objGet

	} else if objName == "" {
		return fmt.Errorf("Error obtaining edge cluster ID or name during read")
	} else {
		// Get by full name/prefix
		// TODO use 2nd parameter localVarOptionals for paging
		objList, _, err := nsxClient.NetworkTransportApi.ListEdgeClusters(nsxClient.Context, nil)
		if err != nil {
			return fmt.Errorf("Error while reading edge clusters: %v", err)
		}
		// go over the list to find the correct one (prefer a perfect match. If not - prefix match)
		var perfectMatch []manager.EdgeCluster
		var prefixMatch []manager.EdgeCluster
		for _, objInList := range objList.Results {
			if strings.HasPrefix(objInList.DisplayName, objName) {
				prefixMatch = append(prefixMatch, objInList)
			}
			if objInList.DisplayName == objName {
				perfectMatch = append(perfectMatch, objInList)
			}
		}
		if len(perfectMatch) > 0 {
			if len(perfectMatch) > 1 {
				return fmt.Errorf("Found multiple edge clusters with name '%s'", objName)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("Found multiple edge clusters with name starting with '%s'", objName)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("Edge cluster with name '%s' was not found", objName)
		}
	}

	d.SetId(obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("deployment_type", obj.DeploymentType)
	d.Set("member_node_type", obj.MemberNodeType)

	return nil
}

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyVirtualNetworkApplianceClusterRealization() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyVirtualNetworkApplianceClusterRealizationRead,

		Schema: map[string]*schema.Schema{
			"id": getDataSourceIDSchema(),
			"path": {
				Type:         schema.TypeString,
				Description:  "The policy path of the VirtualNetworkApplianceCluster resource",
				Required:     true,
				ValidateFunc: validatePolicyPath(),
			},
			"state": {
				Type:        schema.TypeString,
				Description: "Current realization state of the cluster",
				Computed:    true,
			},
			"timeout": {
				Type:         schema.TypeInt,
				Description:  "Realization timeout in seconds",
				Optional:     true,
				Default:      1800,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"delay": {
				Type:         schema.TypeInt,
				Description:  "Initial delay before starting realization checks, in seconds",
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntAtLeast(0),
			},
		},
	}
}

func dataSourceNsxtPolicyVirtualNetworkApplianceClusterRealizationRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	clusterPath := d.Get("path").(string)
	delay := d.Get("delay").(int)
	timeout := d.Get("timeout").(int)

	siteID := getResourceIDFromResourcePath(clusterPath, "sites")
	if siteID == "" {
		return fmt.Errorf("error obtaining site ID from path %s", clusterPath)
	}
	epID := getResourceIDFromResourcePath(clusterPath, "enforcement-points")
	if epID == "" {
		return fmt.Errorf("error obtaining enforcement-point ID from path %s", clusterPath)
	}
	clusterID := getResourceIDFromResourcePath(clusterPath, "virtual-network-appliance-clusters")
	if clusterID == "" {
		return fmt.Errorf("error obtaining cluster ID from path %s", clusterPath)
	}

	id := d.Get("id").(string)
	if id == "" {
		d.SetId(newUUID())
	}

	stateClient := cliVNAClusterStateClient(sessionContext, connector)
	if stateClient == nil {
		return policyResourceNotSupportedError()
	}

	pendingStates := []string{
		model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_UNINITIALIZED,
		model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_IN_PROGRESS,
		model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_SANDBOXED_REALIZATION_PENDING,
		model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_UNKNOWN,
	}
	targetStates := []string{
		model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_SUCCESS,
		model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_ERROR,
		model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_DOWN,
		model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_DEGRAGED, //nolint:misspell
		model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_DISABLED,
	}

	stateConf := &resource.StateChangeConf{
		Pending: pendingStates,
		Target:  targetStates,
		Refresh: func() (interface{}, string, error) {
			clusterState, err := stateClient.Get(siteID, epID, clusterID)
			if err != nil {
				return clusterState, model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_ERROR,
					logAPIError("Error while waiting for realization of VirtualNetworkApplianceCluster", err)
			}
			if clusterState.ConsolidatedStatus == nil {
				log.Printf("[DEBUG] VirtualNetworkApplianceCluster %s realization state is unknown", clusterID)
				return clusterState, model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_UNKNOWN, nil
			}
			log.Printf("[DEBUG] VirtualNetworkApplianceCluster %s realization state: %s", clusterID, *clusterState.ConsolidatedStatus)
			d.Set("state", *clusterState.ConsolidatedStatus)
			return clusterState, *clusterState.ConsolidatedStatus, nil
		},
		Timeout:    time.Duration(timeout) * time.Second,
		MinTimeout: 5 * time.Second,
		Delay:      time.Duration(delay) * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("failed to get realization information for VirtualNetworkApplianceCluster %s: %v", clusterPath, err)
	}
	return nil
}

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyVirtualNetworkAppliance_basic(t *testing.T) {
	testResourceName := "nsxt_policy_virtual_network_appliance.test"
	testDataSourceName := "data.nsxt_policy_virtual_network_appliance.test"
	displayName := getAccTestResourceName()
	edgeTransportNodeName := getEdgeTransportNodeName()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccVNAPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNsxtPolicyVirtualNetworkApplianceCheckDestroy(testResourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyVirtualNetworkApplianceDataSourceTemplate(displayName, edgeTransportNodeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(testDataSourceName, "id", testResourceName, "id"),
					resource.TestCheckResourceAttrPair(testDataSourceName, "path", testResourceName, "path"),
					resource.TestCheckResourceAttr(testDataSourceName, "display_name", displayName),
				),
			},
		},
	})
}

func testAccNsxtPolicyVirtualNetworkApplianceDataSourceTemplate(displayName, edgeTransportNodeName string) string {
	return testAccNsxtPolicyVirtualNetworkApplianceClusterTemplate(displayName, edgeTransportNodeName) + fmt.Sprintf(`
resource "nsxt_policy_virtual_network_appliance" "test" {
  display_name = "%s"
  description  = "Acceptance test VNA for data source"
  cluster_path = nsxt_policy_virtual_network_appliance_cluster.test.path
  hostname     = "%s"

  management_interface {
    network_id = "%s"

    ip_assignment {
      dhcp_v4 = true
    }
  }

  vm_deployment_config {
    compute_manager_id          = data.nsxt_compute_manager.test.id
    cluster_or_resource_pool_id = data.nsxt_compute_collection.test.cm_local_id
    datastore_id                = "%s"
  }
}

data "nsxt_policy_virtual_network_appliance" "test" {
  display_name = "%s"
  cluster_path = nsxt_policy_virtual_network_appliance_cluster.test.path

  depends_on = [nsxt_policy_virtual_network_appliance.test]
}
`, displayName, getVNAHostname(), getVNAPortgroupID(), getDatastoreID(), displayName)
}

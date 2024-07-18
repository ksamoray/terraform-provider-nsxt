/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var accTestPolicyTransitGatewayAttachmentCreateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform created",
}

var accTestPolicyTransitGatewayAttachmentUpdateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform updated",
}

func TestAccResourceNsxtPolicyTransitGatewayAttachment_basic(t *testing.T) {
	testResourceName := "nsxt_policy_transit_gateway_attachment.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTransitGatewayAttachmentCheckDestroy(state, accTestPolicyTransitGatewayAttachmentUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTransitGatewayAttachmentTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayAttachmentExists(accTestPolicyTransitGatewayAttachmentCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyTransitGatewayAttachmentCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyTransitGatewayAttachmentCreateAttributes["description"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyTransitGatewayAttachmentTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayAttachmentExists(accTestPolicyTransitGatewayAttachmentUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyTransitGatewayAttachmentUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyTransitGatewayAttachmentUpdateAttributes["description"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyTransitGatewayAttachmentMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayAttachmentExists(accTestPolicyTransitGatewayAttachmentCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTransitGatewayAttachment_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_transit_gateway_attachment.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTransitGatewayAttachmentCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTransitGatewayAttachmentMinimalistic(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicyTransitGatewayAttachmentExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy TransitGatewayAttachment resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy TransitGatewayAttachment resource ID not set in resources")
		}

		parentPath := rs.Primary.Attributes["parent_path"]
		exists, err := resourceNsxtPolicyTransitGatewayAttachmentExists(testAccGetSessionContext(), parentPath, resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy TransitGatewayAttachment %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyTransitGatewayAttachmentCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_transit_gateway_attachment" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		parentPath := rs.Primary.Attributes["parent_path"]
		exists, err := resourceNsxtPolicyTransitGatewayAttachmentExists(testAccGetSessionContext(), parentPath, resourceID, connector)
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy TransitGatewayAttachment %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyTransitGatewayAttachmentTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyTransitGatewayAttachmentCreateAttributes
	} else {
		attrMap = accTestPolicyTransitGatewayAttachmentUpdateAttributes
	}
	return testAccNsxtPolicyGatewayConnectionTemplate(true) + testAccNsxtVpcConnectivityProfilePrerequisite() + fmt.Sprintf(`
data "nsxt_policy_project" "test" {
  id = "%s"
}

resource "nsxt_policy_share" "share" {
  display_name = "share"
  shared_with  = [data.nsxt_policy_project.test.path]
}

resource "nsxt_policy_shared_resource" "test" {
  share_path = nsxt_policy_share.share.path
  display_name = "shared_resource"
  resource_object {
    resource_path = nsxt_policy_gateway_connection.test.path
  }
}

resource "nsxt_policy_transit_gateway_attachment" "test" {
  parent_path  = nsxt_transit_gateway.test.path
  connection_path = nsxt_policy_gateway_connection.test.path
  display_name = "%s"
  description  = "%s"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  depends_on = [nsxt_policy_shared_resource.test]
}`, os.Getenv("NSXT_PROJECT_ID"), attrMap["display_name"], attrMap["description"])
}

func testAccNsxtPolicyTransitGatewayAttachmentMinimalistic() string {
	return testAccNsxtVpcConnectivityProfilePrerequisite() + fmt.Sprintf(`
resource "nsxt_policy_transit_gateway_attachment" "test" {
  parent_path  = nsxt_transit_gateway.test.path
  display_name = "%s"

}`, accTestPolicyTransitGatewayAttachmentUpdateAttributes["display_name"])
}

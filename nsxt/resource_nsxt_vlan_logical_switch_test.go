/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"regexp"
	"testing"
)

func TestAccResourceNsxtVlanLogicalSwitch_basic(t *testing.T) {
	switchName := "test-nsx-logical-switch-vlan"
	updateSwitchName := fmt.Sprintf("%s-update", switchName)
	transportZoneName := getVlanTransportZoneName()
	resourceName := "testvlan"
	testResourceName := fmt.Sprintf("nsxt_vlan_logical_switch.%s", resourceName)

	origvlan := "1"
	updatedvlan := "2"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXLogicalSwitchCheckDestroy(state, switchName, "nsxt_vlan_logical_switch")
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXVlanLogicalSwitchCreateTemplate(resourceName, switchName, transportZoneName, origvlan),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalSwitchExists(switchName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", switchName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "admin_state", "UP"),
					resource.TestCheckResourceAttr(testResourceName, "vlan", origvlan),
				),
			},
			{
				Config: testAccNSXVlanLogicalSwitchUpdateTemplate(resourceName, updateSwitchName, transportZoneName, updatedvlan),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalSwitchExists(updateSwitchName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateSwitchName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Update"),
					resource.TestCheckResourceAttr(testResourceName, "admin_state", "DOWN"),
					resource.TestCheckResourceAttr(testResourceName, "vlan", updatedvlan),
				),
			},
		},
	})

}

func TestAccResourceNsxtVlanLogicalSwitch_withProfiles(t *testing.T) {
	switchName := fmt.Sprintf("test-nsx-logical-switch-with-profiles")
	updateSwitchName := fmt.Sprintf("%s-update", switchName)
	resourceName := "test_profiles"
	testResourceName := fmt.Sprintf("nsxt_vlan_logical_switch.%s", resourceName)
	transportZoneName := getVlanTransportZoneName()
	customProfileName := "terraform_test_LS_profile"
	oobProfileName := "nsx-default-switch-security-vif-profile"
	profileType := "SwitchSecuritySwitchingProfile"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			// Verify that the LS was deleted
			err := testAccNSXLogicalSwitchCheckDestroy(state, switchName, "nsxt_vlan_logical_switch")
			if err != nil {
				return err
			}
			// Delete the created switching profile
			return testAccDataSourceNsxtSwitchingProfileDeleteByName(customProfileName)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					// Create a custom switching profile
					if err := testAccDataSourceNsxtSwitchingProfileCreate(customProfileName, profileType); err != nil {
						panic(err)
					}
				},
				// Create a logical switch to use the custom switching profile
				Config: testAccNSXVlanLogicalSwitchCreateWithProfilesTemplate(resourceName, switchName, transportZoneName, customProfileName),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalSwitchExists(switchName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", switchName),
					resource.TestCheckResourceAttr(testResourceName, "switching_profile_id.#", "1"),
				),
			},
			{
				// Replace the custom switching profile with OOB one
				Config:             testAccNSXVlanLogicalSwitchUpdateWithProfilesTemplate(resourceName, updateSwitchName, transportZoneName, customProfileName, oobProfileName),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalSwitchExists(updateSwitchName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateSwitchName),
					// Counting only custom profiles so count should be 0
					resource.TestCheckResourceAttr(testResourceName, "switching_profile_id.#", "0"),
				),
			},
			{
				// remove the data source for the custom switching profile
				Config: testAccNSXSwitchingNoProfileTemplate(),
			},
		},
	})
}

func TestAccResourceNsxtVlanLogicalSwitch_withMacPool(t *testing.T) {
	switchName := fmt.Sprintf("test-nsx-logical-switch-with-mac")
	resourceName := "test_mac_pool"
	testResourceName := fmt.Sprintf("nsxt_vlan_logical_switch.%s", resourceName)
	transportZoneName := getVlanTransportZoneName()
	macPoolName := getMacPoolName()
	novlan := "0"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXLogicalSwitchCheckDestroy(state, switchName, "nsxt_vlan_logical_switch")
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccNSXLogicalSwitchNoTZIDTemplate(switchName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: testAccNSXVlanLogicalSwitchCreateWithMacTemplate(resourceName, switchName, transportZoneName, macPoolName, novlan),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalSwitchExists(switchName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", switchName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "admin_state", "UP"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "vlan", novlan),
				),
			},
		},
	})
}

func TestAccResourceNsxtVlanLogicalSwitch_importBasic(t *testing.T) {
	switchName := fmt.Sprintf("test-nsx-logical-switch")
	resourceName := "test"
	testResourceName := fmt.Sprintf("nsxt_vlan_logical_switch.%s", resourceName)
	vlan := "5"
	transportZoneName := getVlanTransportZoneName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXLogicalSwitchCheckDestroy(state, switchName, "nsxt_vlan_logical_switch")
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccNSXLogicalSwitchNoTZIDTemplate(switchName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: testAccNSXVlanLogicalSwitchCreateTemplate(resourceName, switchName, transportZoneName, vlan),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNSXVlanLogicalSwitchNoTZIDTemplate(switchName string) string {
	return fmt.Sprintf(`
resource "nsxt_vlan_logical_switch" "error" {
  display_name     = "%s"
  admin_state      = "UP"
  description      = "Acceptance Test"
}`, switchName)
}

func testAccNSXVlanLogicalSwitchNoVlanTemplate(switchName string, transportZoneName string) string {
	return fmt.Sprintf(`
data "nsxt_transport_zone" "TZ1" {
  display_name = "%s"
}

resource "nsxt_vlan_logical_switch" "error" {
  display_name      = "%s"
  admin_state       = "UP"
  description       = "Acceptance Test"
  transport_zone_id = "${data.nsxt_transport_zone.TZ1.id}"
}`, transportZoneName, switchName)
}

func testAccNSXVlanLogicalSwitchCreateTemplate(resourceName string, switchName string, transportZoneName string, vlan string) string {
	return fmt.Sprintf(`
data "nsxt_transport_zone" "TZ1" {
  display_name = "%s"
}

resource "nsxt_vlan_logical_switch" "%s" {
  display_name      = "%s"
  admin_state       = "UP"
  description       = "Acceptance Test"
  transport_zone_id = "${data.nsxt_transport_zone.TZ1.id}"
  vlan              = "%s"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, transportZoneName, resourceName, switchName, vlan)
}

func testAccNSXVlanLogicalSwitchUpdateTemplate(resourceName string, switchUpdateName string, transportZoneName string, vlan string) string {
	return fmt.Sprintf(`
data "nsxt_transport_zone" "TZ1" {
  display_name = "%s"
}

resource "nsxt_vlan_logical_switch" "%s" {
  display_name      = "%s"
  admin_state       = "DOWN"
  description       = "Acceptance Test Update"
  transport_zone_id = "${data.nsxt_transport_zone.TZ1.id}"
  vlan              = "%s"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  tag {
    scope = "scope2"
    tag   = "tag2"
  }
}`, transportZoneName, resourceName, switchUpdateName, vlan)
}

func testAccNSXVlanLogicalSwitchCreateWithProfilesTemplate(resourceName string, switchName string, transportZoneName string, profileName string) string {
	return fmt.Sprintf(`
data "nsxt_transport_zone" "TZ1" {
  display_name = "%s"
}

data "nsxt_switching_profile" "test1" {
  display_name = "%s"
}

resource "nsxt_vlan_logical_switch" "%s" {
  display_name      = "%s"
  admin_state       = "UP"
  transport_zone_id = "${data.nsxt_transport_zone.TZ1.id}"
  vlan              = 1

  switching_profile_id {
    key   = "${data.nsxt_switching_profile.test1.resource_type}"
    value = "${data.nsxt_switching_profile.test1.id}"
  }
}`, transportZoneName, profileName, resourceName, switchName)
}

func testAccNSXVlanLogicalSwitchUpdateWithProfilesTemplate(resourceName string, switchUpdateName string, transportZoneName string, profileName1 string, profileName2 string) string {
	return fmt.Sprintf(`
data "nsxt_transport_zone" "TZ1" {
  display_name = "%s"
}

data "nsxt_switching_profile" "test1" {
  display_name = "%s"
}

data "nsxt_switching_profile" "test2" {
  display_name = "%s"
}

resource "nsxt_vlan_logical_switch" "%s" {
  display_name      = "%s"
  admin_state       = "UP"
  transport_zone_id = "${data.nsxt_transport_zone.TZ1.id}"
  vlan              = 1

  switching_profile_id {
    key   = "${data.nsxt_switching_profile.test2.resource_type}"
    value = "${data.nsxt_switching_profile.test2.id}"
  }

}`, transportZoneName, profileName1, profileName2, resourceName, switchUpdateName)
}

func testAccNSXVlanLogicalSwitchCreateWithMacTemplate(resourceName string, switchName string, transportZoneName string, macPoolName string, vlan string) string {
	return fmt.Sprintf(`
data "nsxt_transport_zone" "TZ1" {
  display_name = "%s"
}

data "nsxt_mac_pool" "MAC1" {
  display_name = "%s"
}

resource "nsxt_vlan_logical_switch" "%s" {
  display_name      = "%s"
  admin_state       = "UP"
  description       = "Acceptance Test"
  transport_zone_id = "${data.nsxt_transport_zone.TZ1.id}"
  vlan              = "%s"
  mac_pool_id       = "${data.nsxt_mac_pool.MAC1.id}"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, transportZoneName, macPoolName, resourceName, switchName, vlan)
}

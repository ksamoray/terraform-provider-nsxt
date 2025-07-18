// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicySpoofGuardProfile_basic(t *testing.T) {
	// Use existing system defined profile
	name := "default-spoofguard-profile"
	testResourceName := "data.nsxt_policy_spoofguard_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicySpoofGuardProfileReadTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicySpoofGuardProfile_multitenancy(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_spoofguard_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyMultitenancy(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicySpoofGuardProfileMultitenancyTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}

func testAccNsxtPolicySpoofGuardProfileReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_spoofguard_profile" "test" {
  display_name = "%s"
}`, name)
}

func testAccNsxtPolicySpoofGuardProfileMultitenancyTemplate(name string) string {
	context := testAccNsxtPolicyMultitenancyContext()
	return fmt.Sprintf(`
resource "nsxt_policy_spoof_guard_profile" "test" {
%s
  display_name = "%s"

}
data "nsxt_policy_spoofguard_profile" "test" {
%s
  display_name = nsxt_policy_spoof_guard_profile.test.display_name
}`, context, name, context)
}

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	clientLayer "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/transit_gateways/nat"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
)

var transitGatewayNatPathExample = "/orgs/[org]/projects/[project]/transit-gateways/[gateway]/nat/[type]"

func resourceNsxtPolicyTransitGatewayNatRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyTransitGatewayNatRuleCreate,
		Read:   resourceNsxtPolicyTransitGatewayNatRuleRead,
		Update: resourceNsxtPolicyTransitGatewayNatRuleUpdate,
		Delete: resourceNsxtPolicyTransitGatewayNatRuleDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtParentPathResourceImporter,
		},
		// Today, transit gateway nat rule schema is equal to VPC nat rule schema
		// If in future they diverge, we'll introduce new schema here
		Schema: metadata.GetSchemaFromExtendedSchema(getPolicyVpcNatRuleSchema(true)),
	}
}

func resourceNsxtPolicyTransitGatewayNatRuleExists(sessionContext utl.SessionContext, parentPath string, id string, connector client.Connector) (bool, error) {
	var err error
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 4, transitGatewayNatPathExample)
	if pathErr != nil {
		return false, pathErr
	}
	client := clientLayer.NewNatRulesClient(connector)
	_, err = client.Get(parents[0], parents[1], parents[2], parents[3], id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtPolicyTransitGatewayNatRuleCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id, err := getOrGenerateIDWithParent(d, m, resourceNsxtPolicyTransitGatewayNatRuleExists)
	if err != nil {
		return err
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 4, transitGatewayNatPathExample)
	if pathErr != nil {
		return pathErr
	}
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.TransitGatewayNatRule{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, getPolicyVpcNatRuleSchema(true), "", nil); err != nil {
		return err
	}

	log.Printf("[INFO] Creating PolicyTransitGatewayNatRule with ID %s", id)

	client := clientLayer.NewNatRulesClient(connector)
	err = client.Patch(parents[0], parents[1], parents[2], parents[3], id, obj)
	if err != nil {
		return handleCreateError("PolicyTransitGatewayNatRule", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyTransitGatewayNatRuleRead(d, m)
}

func resourceNsxtPolicyTransitGatewayNatRuleRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining PolicyTransitGatewayNatRule ID")
	}

	client := clientLayer.NewNatRulesClient(connector)
	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 4, transitGatewayNatPathExample)
	if pathErr != nil {
		return pathErr
	}
	obj, err := client.Get(parents[0], parents[1], parents[2], parents[3], id)
	if err != nil {
		return handleReadError(d, "PolicyTransitGatewayNatRule", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	elem := reflect.ValueOf(&obj).Elem()
	return metadata.StructToSchema(elem, d, getPolicyVpcNatRuleSchema(true), "", nil)
}

func resourceNsxtPolicyTransitGatewayNatRuleUpdate(d *schema.ResourceData, m interface{}) error {

	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining PolicyTransitGatewayNatRule ID")
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 4, transitGatewayNatPathExample)
	if pathErr != nil {
		return pathErr
	}
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)

	revision := int64(d.Get("revision").(int))

	obj := model.TransitGatewayNatRule{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Revision:    &revision,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, getPolicyVpcNatRuleSchema(true), "", nil); err != nil {
		return err
	}
	client := clientLayer.NewNatRulesClient(connector)
	_, err := client.Update(parents[0], parents[1], parents[2], parents[3], id, obj)
	if err != nil {
		return handleUpdateError("PolicyTransitGatewayNatRule", id, err)
	}

	return resourceNsxtPolicyTransitGatewayNatRuleRead(d, m)
}

func resourceNsxtPolicyTransitGatewayNatRuleDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining PolicyTransitGatewayNatRule ID")
	}

	connector := getPolicyConnector(m)
	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 4, transitGatewayNatPathExample)
	if pathErr != nil {
		return pathErr
	}

	client := clientLayer.NewNatRulesClient(connector)
	err := client.Delete(parents[0], parents[1], parents[2], parents[3], id)

	if err != nil {
		return handleDeleteError("PolicyTransitGatewayNatRule", id, err)
	}

	return nil
}

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/go-vmware-nsxt/manager"
)

func resourceNsxtIgmpTypeNsService() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtIgmpTypeNsServiceCreate,
		Read:   resourceNsxtIgmpTypeNsServiceRead,
		Update: resourceNsxtIgmpTypeNsServiceUpdate,
		Delete: resourceNsxtIgmpTypeNsServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: mpObjectResourceDeprecationMessage,
		Schema: map[string]*schema.Schema{
			"revision": getRevisionSchema(),
			"description": {
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "The display name of this resource. Defaults to ID if not set",
				Optional:    true,
				Computed:    true,
			},
			"tag": getTagsSchema(),
			"default_service": {
				Type:        schema.TypeBool,
				Description: "A boolean flag which reflects whether this is a default NSServices which can't be modified/deleted",
				Computed:    true,
			},
		},
	}
}

func resourceNsxtIgmpTypeNsServiceCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return resourceNotSupportedError()
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)

	nsService := manager.IgmpTypeNsService{
		NsService: manager.NsService{
			Description: description,
			DisplayName: displayName,
			Tags:        tags,
		},
		NsserviceElement: manager.IgmpTypeNsServiceEntry{
			ResourceType: "IGMPTypeNSService",
		},
	}

	nsService, resp, err := nsxClient.GroupingObjectsApi.CreateIgmpTypeNSService(nsxClient.Context, nsService)

	if err != nil {
		return fmt.Errorf("Error during NsService create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during NsService create: %v", resp.StatusCode)
	}
	d.SetId(nsService.Id)
	return resourceNsxtIgmpTypeNsServiceRead(d, m)
}

func resourceNsxtIgmpTypeNsServiceRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return resourceNotSupportedError()
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining ns service id")
	}

	nsService, resp, err := nsxClient.GroupingObjectsApi.ReadIgmpTypeNSService(nsxClient.Context, id)
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] NsService %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during NsService read: %v", err)
	}

	d.Set("revision", nsService.Revision)
	d.Set("description", nsService.Description)
	d.Set("display_name", nsService.DisplayName)
	setTagsInSchema(d, nsService.Tags)
	d.Set("default_service", nsService.DefaultService)

	return nil
}

func resourceNsxtIgmpTypeNsServiceUpdate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return resourceNotSupportedError()
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining ns service id")
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	revision := int64(d.Get("revision").(int))

	nsService := manager.IgmpTypeNsService{
		NsService: manager.NsService{
			Description: description,
			DisplayName: displayName,
			Tags:        tags,
			Revision:    revision,
		},
		NsserviceElement: manager.IgmpTypeNsServiceEntry{
			ResourceType: "IGMPTypeNSService",
		},
	}

	_, resp, err := nsxClient.GroupingObjectsApi.UpdateIgmpTypeNSService(nsxClient.Context, id, nsService)
	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during NsService update: %v %v", err, resp)
	}

	return resourceNsxtIgmpTypeNsServiceRead(d, m)
}

func resourceNsxtIgmpTypeNsServiceDelete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return resourceNotSupportedError()
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining ns service id")
	}

	localVarOptionals := make(map[string]interface{})
	localVarOptionals["force"] = true
	resp, err := nsxClient.GroupingObjectsApi.DeleteNSService(nsxClient.Context, id, localVarOptionals)
	if err != nil {
		return fmt.Errorf("Error during NsService delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] NsService %s not found", id)
		d.SetId("")
	}
	return nil
}

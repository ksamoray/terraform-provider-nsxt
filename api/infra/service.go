//nolint:revive
package infra

// The following file has been autogenerated. Please avoid any changes!
import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	client1 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra"
	model1 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	client2 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/infra"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type ServiceClientContext utl.ClientContext

func NewServicesClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *ServiceClientContext {
	var client interface{}

	switch sessionContext.ClientType {

	case utl.Local:
		client = client0.NewServicesClient(connector)

	case utl.Global:
		client = client1.NewServicesClient(connector)

	case utl.Multitenancy:
		client = client2.NewServicesClient(connector)

	default:
		return nil
	}
	return &ServiceClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID, VPCID: sessionContext.VPCID}
}

func (c ServiceClientContext) Get(serviceIdParam string) (model0.Service, error) {
	var obj model0.Service
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.ServicesClient)
		obj, err = client.Get(serviceIdParam)
		if err != nil {
			return obj, err
		}

	case utl.Global:
		client := c.Client.(client1.ServicesClient)
		gmObj, err1 := client.Get(serviceIdParam)
		if err1 != nil {
			return obj, err1
		}
		var rawObj interface{}
		rawObj, err = utl.ConvertModelBindingType(gmObj, model1.ServiceBindingType(), model0.ServiceBindingType())
		obj = rawObj.(model0.Service)

	case utl.Multitenancy:
		client := c.Client.(client2.ServicesClient)
		obj, err = client.Get(utl.DefaultOrgID, c.ProjectID, serviceIdParam)
		if err != nil {
			return obj, err
		}

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c ServiceClientContext) Patch(serviceIdParam string, serviceParam model0.Service) error {
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.ServicesClient)
		err = client.Patch(serviceIdParam, serviceParam)

	case utl.Global:
		client := c.Client.(client1.ServicesClient)
		gmObj, err1 := utl.ConvertModelBindingType(serviceParam, model0.ServiceBindingType(), model1.ServiceBindingType())
		if err1 != nil {
			return err1
		}
		err = client.Patch(serviceIdParam, gmObj.(model1.Service))

	case utl.Multitenancy:
		client := c.Client.(client2.ServicesClient)
		err = client.Patch(utl.DefaultOrgID, c.ProjectID, serviceIdParam, serviceParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c ServiceClientContext) Update(serviceIdParam string, serviceParam model0.Service) (model0.Service, error) {
	var err error
	var obj model0.Service

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.ServicesClient)
		obj, err = client.Update(serviceIdParam, serviceParam)

	case utl.Global:
		client := c.Client.(client1.ServicesClient)
		gmObj, err := utl.ConvertModelBindingType(serviceParam, model0.ServiceBindingType(), model1.ServiceBindingType())
		if err != nil {
			return obj, err
		}
		gmObj, err = client.Update(serviceIdParam, gmObj.(model1.Service))
		if err != nil {
			return obj, err
		}
		obj1, err1 := utl.ConvertModelBindingType(gmObj, model1.ServiceBindingType(), model0.ServiceBindingType())
		if err1 != nil {
			return obj, err1
		}
		obj = obj1.(model0.Service)

	case utl.Multitenancy:
		client := c.Client.(client2.ServicesClient)
		obj, err = client.Update(utl.DefaultOrgID, c.ProjectID, serviceIdParam, serviceParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c ServiceClientContext) Delete(serviceIdParam string) error {
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.ServicesClient)
		err = client.Delete(serviceIdParam)

	case utl.Global:
		client := c.Client.(client1.ServicesClient)
		err = client.Delete(serviceIdParam)

	case utl.Multitenancy:
		client := c.Client.(client2.ServicesClient)
		err = client.Delete(utl.DefaultOrgID, c.ProjectID, serviceIdParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c ServiceClientContext) List(cursorParam *string, defaultServiceParam *bool, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model0.ServiceListResult, error) {
	var err error
	var obj model0.ServiceListResult

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.ServicesClient)
		obj, err = client.List(cursorParam, defaultServiceParam, includeMarkForDeleteObjectsParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)

	case utl.Global:
		client := c.Client.(client1.ServicesClient)
		gmObj, err := client.List(cursorParam, defaultServiceParam, includeMarkForDeleteObjectsParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)
		if err != nil {
			return obj, err
		}
		obj1, err1 := utl.ConvertModelBindingType(gmObj, model1.ServiceListResultBindingType(), model0.ServiceListResultBindingType())
		if err1 != nil {
			return obj, err1
		}
		obj = obj1.(model0.ServiceListResult)

	case utl.Multitenancy:
		client := c.Client.(client2.ServicesClient)
		obj, err = client.List(utl.DefaultOrgID, c.ProjectID, cursorParam, defaultServiceParam, includeMarkForDeleteObjectsParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return obj, err
}

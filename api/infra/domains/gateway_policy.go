//nolint:revive
package domains

// The following file has been autogenerated. Please avoid any changes!
import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	client1 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra/domains"
	model1 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/domains"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	client2 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/infra/domains"
	client3 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/vpcs"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type GatewayPolicyClientContext utl.ClientContext

func NewGatewayPoliciesClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *GatewayPolicyClientContext {
	var client interface{}

	switch sessionContext.ClientType {

	case utl.Local:
		client = client0.NewGatewayPoliciesClient(connector)

	case utl.Global:
		client = client1.NewGatewayPoliciesClient(connector)

	case utl.Multitenancy:
		client = client2.NewGatewayPoliciesClient(connector)

	case utl.VPC:
		client = client3.NewGatewayPoliciesClient(connector)

	default:
		return nil
	}
	return &GatewayPolicyClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID, VPCID: sessionContext.VPCID}
}

func (c GatewayPolicyClientContext) Get(domainIdParam string, gatewayPolicyIdParam string) (model0.GatewayPolicy, error) {
	var obj model0.GatewayPolicy
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.GatewayPoliciesClient)
		obj, err = client.Get(domainIdParam, gatewayPolicyIdParam)
		if err != nil {
			return obj, err
		}

	case utl.Global:
		client := c.Client.(client1.GatewayPoliciesClient)
		gmObj, err1 := client.Get(domainIdParam, gatewayPolicyIdParam)
		if err1 != nil {
			return obj, err1
		}
		var rawObj interface{}
		rawObj, err = utl.ConvertModelBindingType(gmObj, model1.GatewayPolicyBindingType(), model0.GatewayPolicyBindingType())
		obj = rawObj.(model0.GatewayPolicy)

	case utl.Multitenancy:
		client := c.Client.(client2.GatewayPoliciesClient)
		obj, err = client.Get(utl.DefaultOrgID, c.ProjectID, domainIdParam, gatewayPolicyIdParam)
		if err != nil {
			return obj, err
		}

	case utl.VPC:
		client := c.Client.(client3.GatewayPoliciesClient)
		obj, err = client.Get(utl.DefaultOrgID, c.ProjectID, c.VPCID, gatewayPolicyIdParam)
		if err != nil {
			return obj, err
		}

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c GatewayPolicyClientContext) Patch(domainIdParam string, gatewayPolicyIdParam string, gatewayPolicyParam model0.GatewayPolicy) error {
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.GatewayPoliciesClient)
		err = client.Patch(domainIdParam, gatewayPolicyIdParam, gatewayPolicyParam)

	case utl.Global:
		client := c.Client.(client1.GatewayPoliciesClient)
		gmObj, err1 := utl.ConvertModelBindingType(gatewayPolicyParam, model0.GatewayPolicyBindingType(), model1.GatewayPolicyBindingType())
		if err1 != nil {
			return err1
		}
		err = client.Patch(domainIdParam, gatewayPolicyIdParam, gmObj.(model1.GatewayPolicy))

	case utl.Multitenancy:
		client := c.Client.(client2.GatewayPoliciesClient)
		err = client.Patch(utl.DefaultOrgID, c.ProjectID, domainIdParam, gatewayPolicyIdParam, gatewayPolicyParam)

	case utl.VPC:
		client := c.Client.(client3.GatewayPoliciesClient)
		err = client.Patch(utl.DefaultOrgID, c.ProjectID, c.VPCID, gatewayPolicyIdParam, gatewayPolicyParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c GatewayPolicyClientContext) Update(domainIdParam string, gatewayPolicyIdParam string, gatewayPolicyParam model0.GatewayPolicy) (model0.GatewayPolicy, error) {
	var err error
	var obj model0.GatewayPolicy

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.GatewayPoliciesClient)
		obj, err = client.Update(domainIdParam, gatewayPolicyIdParam, gatewayPolicyParam)

	case utl.Global:
		client := c.Client.(client1.GatewayPoliciesClient)
		gmObj, err := utl.ConvertModelBindingType(gatewayPolicyParam, model0.GatewayPolicyBindingType(), model1.GatewayPolicyBindingType())
		if err != nil {
			return obj, err
		}
		gmObj, err = client.Update(domainIdParam, gatewayPolicyIdParam, gmObj.(model1.GatewayPolicy))
		if err != nil {
			return obj, err
		}
		obj1, err1 := utl.ConvertModelBindingType(gmObj, model1.GatewayPolicyBindingType(), model0.GatewayPolicyBindingType())
		if err1 != nil {
			return obj, err1
		}
		obj = obj1.(model0.GatewayPolicy)

	case utl.Multitenancy:
		client := c.Client.(client2.GatewayPoliciesClient)
		obj, err = client.Update(utl.DefaultOrgID, c.ProjectID, domainIdParam, gatewayPolicyIdParam, gatewayPolicyParam)

	case utl.VPC:
		client := c.Client.(client3.GatewayPoliciesClient)
		obj, err = client.Update(utl.DefaultOrgID, c.ProjectID, c.VPCID, gatewayPolicyIdParam, gatewayPolicyParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c GatewayPolicyClientContext) Delete(domainIdParam string, gatewayPolicyIdParam string) error {
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.GatewayPoliciesClient)
		err = client.Delete(domainIdParam, gatewayPolicyIdParam)

	case utl.Global:
		client := c.Client.(client1.GatewayPoliciesClient)
		err = client.Delete(domainIdParam, gatewayPolicyIdParam)

	case utl.Multitenancy:
		client := c.Client.(client2.GatewayPoliciesClient)
		err = client.Delete(utl.DefaultOrgID, c.ProjectID, domainIdParam, gatewayPolicyIdParam)

	case utl.VPC:
		client := c.Client.(client3.GatewayPoliciesClient)
		err = client.Delete(utl.DefaultOrgID, c.ProjectID, c.VPCID, gatewayPolicyIdParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c GatewayPolicyClientContext) List(domainIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includeRuleCountParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model0.GatewayPolicyListResult, error) {
	var err error
	var obj model0.GatewayPolicyListResult

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.GatewayPoliciesClient)
		obj, err = client.List(domainIdParam, cursorParam, includeMarkForDeleteObjectsParam, includeRuleCountParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)

	case utl.Global:
		client := c.Client.(client1.GatewayPoliciesClient)
		gmObj, err := client.List(domainIdParam, cursorParam, includeMarkForDeleteObjectsParam, includeRuleCountParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)
		if err != nil {
			return obj, err
		}
		obj1, err1 := utl.ConvertModelBindingType(gmObj, model1.GatewayPolicyListResultBindingType(), model0.GatewayPolicyListResultBindingType())
		if err1 != nil {
			return obj, err1
		}
		obj = obj1.(model0.GatewayPolicyListResult)

	case utl.Multitenancy:
		client := c.Client.(client2.GatewayPoliciesClient)
		obj, err = client.List(utl.DefaultOrgID, c.ProjectID, domainIdParam, cursorParam, includeMarkForDeleteObjectsParam, includeRuleCountParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)

	case utl.VPC:
		client := c.Client.(client3.GatewayPoliciesClient)
		obj, err = client.List(utl.DefaultOrgID, c.ProjectID, c.VPCID, cursorParam, includeMarkForDeleteObjectsParam, includeRuleCountParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return obj, err
}

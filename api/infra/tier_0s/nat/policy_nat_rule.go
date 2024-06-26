//nolint:revive
package nat

// The following file has been autogenerated. Please avoid any changes!
import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	client1 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra/tier_0s/nat"
	model1 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s/nat"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type PolicyNatRuleClientContext utl.ClientContext

func NewNatRulesClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *PolicyNatRuleClientContext {
	var client interface{}

	switch sessionContext.ClientType {

	case utl.Local:
		client = client0.NewNatRulesClient(connector)

	case utl.Global:
		client = client1.NewNatRulesClient(connector)

	default:
		return nil
	}
	return &PolicyNatRuleClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID, VPCID: sessionContext.VPCID}
}

func (c PolicyNatRuleClientContext) Get(tier0IdParam string, natIdParam string, natRuleIdParam string) (model0.PolicyNatRule, error) {
	var obj model0.PolicyNatRule
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.NatRulesClient)
		obj, err = client.Get(tier0IdParam, natIdParam, natRuleIdParam)
		if err != nil {
			return obj, err
		}

	case utl.Global:
		client := c.Client.(client1.NatRulesClient)
		gmObj, err1 := client.Get(tier0IdParam, natIdParam, natRuleIdParam)
		if err1 != nil {
			return obj, err1
		}
		var rawObj interface{}
		rawObj, err = utl.ConvertModelBindingType(gmObj, model1.PolicyNatRuleBindingType(), model0.PolicyNatRuleBindingType())
		obj = rawObj.(model0.PolicyNatRule)

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c PolicyNatRuleClientContext) Patch(tier0IdParam string, natIdParam string, natRuleIdParam string, policyNatRuleParam model0.PolicyNatRule) error {
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.NatRulesClient)
		err = client.Patch(tier0IdParam, natIdParam, natRuleIdParam, policyNatRuleParam)

	case utl.Global:
		client := c.Client.(client1.NatRulesClient)
		gmObj, err1 := utl.ConvertModelBindingType(policyNatRuleParam, model0.PolicyNatRuleBindingType(), model1.PolicyNatRuleBindingType())
		if err1 != nil {
			return err1
		}
		err = client.Patch(tier0IdParam, natIdParam, natRuleIdParam, gmObj.(model1.PolicyNatRule))

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c PolicyNatRuleClientContext) Update(tier0IdParam string, natIdParam string, natRuleIdParam string, policyNatRuleParam model0.PolicyNatRule) (model0.PolicyNatRule, error) {
	var err error
	var obj model0.PolicyNatRule

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.NatRulesClient)
		obj, err = client.Update(tier0IdParam, natIdParam, natRuleIdParam, policyNatRuleParam)

	case utl.Global:
		client := c.Client.(client1.NatRulesClient)
		gmObj, err := utl.ConvertModelBindingType(policyNatRuleParam, model0.PolicyNatRuleBindingType(), model1.PolicyNatRuleBindingType())
		if err != nil {
			return obj, err
		}
		gmObj, err = client.Update(tier0IdParam, natIdParam, natRuleIdParam, gmObj.(model1.PolicyNatRule))
		if err != nil {
			return obj, err
		}
		obj1, err1 := utl.ConvertModelBindingType(gmObj, model1.PolicyNatRuleBindingType(), model0.PolicyNatRuleBindingType())
		if err1 != nil {
			return obj, err1
		}
		obj = obj1.(model0.PolicyNatRule)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c PolicyNatRuleClientContext) Delete(tier0IdParam string, natIdParam string, natRuleIdParam string) error {
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.NatRulesClient)
		err = client.Delete(tier0IdParam, natIdParam, natRuleIdParam)

	case utl.Global:
		client := c.Client.(client1.NatRulesClient)
		err = client.Delete(tier0IdParam, natIdParam, natRuleIdParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c PolicyNatRuleClientContext) List(tier0IdParam string, natIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model0.PolicyNatRuleListResult, error) {
	var err error
	var obj model0.PolicyNatRuleListResult

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.NatRulesClient)
		obj, err = client.List(tier0IdParam, natIdParam, cursorParam, includeMarkForDeleteObjectsParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)

	case utl.Global:
		client := c.Client.(client1.NatRulesClient)
		gmObj, err := client.List(tier0IdParam, natIdParam, cursorParam, includeMarkForDeleteObjectsParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)
		if err != nil {
			return obj, err
		}
		obj1, err1 := utl.ConvertModelBindingType(gmObj, model1.PolicyNatRuleListResultBindingType(), model0.PolicyNatRuleListResultBindingType())
		if err1 != nil {
			return obj, err1
		}
		obj = obj1.(model0.PolicyNatRuleListResult)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return obj, err
}

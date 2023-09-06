//nolint:revive
package segments

// The following file has been autogenerated. Please avoid any changes!
import (
	"errors"

	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	client0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/segments"
	model0 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	client1 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/infra/segments"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type SegmentPortClientContext utl.ClientContext

func NewPortsClient(sessionContext utl.SessionContext, connector vapiProtocolClient_.Connector) *SegmentPortClientContext {
	var client interface{}

	switch sessionContext.ClientType {

	case utl.Local:
		client = client0.NewPortsClient(connector)

	case utl.Multitenancy:
		client = client1.NewPortsClient(connector)

	default:
		return nil
	}
	return &SegmentPortClientContext{Client: client, ClientType: sessionContext.ClientType, ProjectID: sessionContext.ProjectID}
}

func (c SegmentPortClientContext) Get(segmentIdParam string, portIdParam string) (model0.SegmentPort, error) {
	var obj model0.SegmentPort
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.PortsClient)
		obj, err = client.Get(segmentIdParam, portIdParam)
		if err != nil {
			return obj, err
		}

	case utl.Multitenancy:
		client := c.Client.(client1.PortsClient)
		obj, err = client.Get(utl.DefaultOrgID, c.ProjectID, segmentIdParam, portIdParam)
		if err != nil {
			return obj, err
		}

	default:
		return obj, errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c SegmentPortClientContext) Patch(segmentIdParam string, portIdParam string, segmentPortParam model0.SegmentPort) error {
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.PortsClient)
		err = client.Patch(segmentIdParam, portIdParam, segmentPortParam)

	case utl.Multitenancy:
		client := c.Client.(client1.PortsClient)
		err = client.Patch(utl.DefaultOrgID, c.ProjectID, segmentIdParam, portIdParam, segmentPortParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c SegmentPortClientContext) Update(segmentIdParam string, portIdParam string, segmentPortParam model0.SegmentPort) (model0.SegmentPort, error) {
	var err error
	var obj model0.SegmentPort

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.PortsClient)
		obj, err = client.Update(segmentIdParam, portIdParam, segmentPortParam)

	case utl.Multitenancy:
		client := c.Client.(client1.PortsClient)
		obj, err = client.Update(utl.DefaultOrgID, c.ProjectID, segmentIdParam, portIdParam, segmentPortParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return obj, err
}

func (c SegmentPortClientContext) Delete(segmentIdParam string, portIdParam string) error {
	var err error

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.PortsClient)
		err = client.Delete(segmentIdParam, portIdParam)

	case utl.Multitenancy:
		client := c.Client.(client1.PortsClient)
		err = client.Delete(utl.DefaultOrgID, c.ProjectID, segmentIdParam, portIdParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return err
}

func (c SegmentPortClientContext) List(segmentIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model0.SegmentPortListResult, error) {
	var err error
	var obj model0.SegmentPortListResult

	switch c.ClientType {

	case utl.Local:
		client := c.Client.(client0.PortsClient)
		obj, err = client.List(segmentIdParam, cursorParam, includeMarkForDeleteObjectsParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)

	case utl.Multitenancy:
		client := c.Client.(client1.PortsClient)
		obj, err = client.List(utl.DefaultOrgID, c.ProjectID, segmentIdParam, cursorParam, includeMarkForDeleteObjectsParam, includedFieldsParam, pageSizeParam, sortAscendingParam, sortByParam)

	default:
		err = errors.New("invalid infrastructure for model")
	}
	return obj, err
}
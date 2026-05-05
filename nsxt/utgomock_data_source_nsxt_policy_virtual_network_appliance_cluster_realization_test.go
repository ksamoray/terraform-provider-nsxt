//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	enforcementpoints "github.com/vmware/terraform-provider-nsxt/api/infra/sites/enforcement_points"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	epmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/sites/enforcement_points"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func setupVNAClusterStateDataSourceMock(t *testing.T, ctrl *gomock.Controller) (*epmocks.MockVirtualNetworkApplianceClusterStateClient, func()) {
	t.Helper()
	mockState := epmocks.NewMockVirtualNetworkApplianceClusterStateClient(ctrl)
	wrapper := &enforcementpoints.VirtualNetworkApplianceClusterStateClientContext{
		Client:     mockState,
		ClientType: utl.Local,
	}
	orig := cliVNAClusterStateClient
	cliVNAClusterStateClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *enforcementpoints.VirtualNetworkApplianceClusterStateClientContext {
		return wrapper
	}
	return mockState, func() { cliVNAClusterStateClient = orig }
}

var vnaClusterRealizationPath = "/infra/sites/default/enforcement-points/default/virtual-network-appliance-clusters/vna-cluster-1"

func TestMockDataSourceNsxtPolicyVirtualNetworkApplianceClusterRealizationRead(t *testing.T) {
	util.NsxVersion = "9.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockState, restore := setupVNAClusterStateDataSourceMock(t, ctrl)
	defer restore()

	ds := dataSourceNsxtPolicyVirtualNetworkApplianceClusterRealization()

	statusSuccess := model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_SUCCESS
	statusError := model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_ERROR

	t.Run("Read_success", func(t *testing.T) {
		mockState.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID).Return(
			model.VirtualNetworkApplianceClusterState{ConsolidatedStatus: &statusSuccess}, nil,
		)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"path": vnaClusterRealizationPath,
		})
		m := newGoMockProviderClient()
		err := dataSourceNsxtPolicyVirtualNetworkApplianceClusterRealizationRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, statusSuccess, d.Get("state"))
	})

	t.Run("Read_stops_on_error_state", func(t *testing.T) {
		mockState.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID).Return(
			model.VirtualNetworkApplianceClusterState{ConsolidatedStatus: &statusError}, nil,
		)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"path": vnaClusterRealizationPath,
		})
		m := newGoMockProviderClient()
		// ERROR is a target state so WaitForState should return without error
		err := dataSourceNsxtPolicyVirtualNetworkApplianceClusterRealizationRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, statusError, d.Get("state"))
	})

	t.Run("Read_fails_on_api_error", func(t *testing.T) {
		mockState.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID).Return(
			model.VirtualNetworkApplianceClusterState{}, vapiErrors.InternalServerError{},
		)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"path":    vnaClusterRealizationPath,
			"timeout": 1,
		})
		m := newGoMockProviderClient()
		err := dataSourceNsxtPolicyVirtualNetworkApplianceClusterRealizationRead(d, m)
		require.Error(t, err)
	})

	t.Run("Read_fails_when_path_has_no_site", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"path": "/infra/no-sites/default/enforcement-points/default/virtual-network-appliance-clusters/id",
		})
		m := newGoMockProviderClient()
		err := dataSourceNsxtPolicyVirtualNetworkApplianceClusterRealizationRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "site ID")
	})
}

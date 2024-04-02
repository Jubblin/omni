// Copyright (c) 2024 Sidero Labs, Inc.
//
// Use of this software is governed by the Business Source License
// included in the LICENSE file.

package omni_test

import (
	"fmt"
	"net"
	"strconv"
	"testing"

	"github.com/cosi-project/runtime/pkg/resource/rtestutils"
	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/siderolabs/omni/client/api/omni/specs"
	"github.com/siderolabs/omni/client/pkg/omni/resources"
	"github.com/siderolabs/omni/client/pkg/omni/resources/omni"
	omnictrl "github.com/siderolabs/omni/internal/backend/runtime/omni/controllers/omni"
	"github.com/siderolabs/omni/internal/pkg/siderolink"
)

type DiscoveryServiceConfigPatchSuite struct {
	OmniSuite
}

func (suite *DiscoveryServiceConfigPatchSuite) TestReconcile() {
	suite.startRuntime()

	port := 1234
	controller := omnictrl.NewDiscoveryServiceConfigPatchController(true, port)

	suite.Require().NoError(suite.runtime.RegisterQController(controller))

	cluster1 := omni.NewCluster(resources.DefaultNamespace, "test-cluster-1")
	newClusterPatchID := omnictrl.DiscoveryServiceConfigPatchPrefix + cluster1.Metadata().ID()

	cluster1.TypedSpec().Value.TalosVersion = "1.6.0"
	cluster1.TypedSpec().Value.Features = &specs.ClusterSpec_Features{
		UseEmbeddedDiscoveryService: true,
	}

	suite.Require().NoError(suite.state.Create(suite.ctx, cluster1))

	// assert that the new cluster is marked to use the embedded discovery service
	rtestutils.AssertResource[*omni.ConfigPatch](suite.ctx, suite.T(), suite.state, newClusterPatchID, func(r *omni.ConfigPatch, assertion *assert.Assertions) {
		assertion.Contains(r.TypedSpec().Value.Data, fmt.Sprintf("http://"+net.JoinHostPort(siderolink.ListenHost, strconv.Itoa(port))))
	})

	_, err := safe.StateUpdateWithConflicts[*omni.Cluster](suite.ctx, suite.state, cluster1.Metadata(), func(res *omni.Cluster) error {
		res.TypedSpec().Value.Features.UseEmbeddedDiscoveryService = false

		return nil
	})
	suite.Require().NoError(err)

	// assert that the endpoint is removed from the config patch
	rtestutils.AssertResource[*omni.ConfigPatch](suite.ctx, suite.T(), suite.state, newClusterPatchID, func(r *omni.ConfigPatch, assertion *assert.Assertions) {
		assertion.Empty(r.TypedSpec().Value.Data)
	})

	// cluster2 is a cluster created after this feature is rolled out, but it has a talos version that does not support the embedded discovery service yet.
	cluster2 := omni.NewCluster(resources.DefaultNamespace, "test-cluster-2")
	unsupportedNewClusterPatchID := omnictrl.DiscoveryServiceConfigPatchPrefix + cluster2.Metadata().ID()

	cluster2.TypedSpec().Value.TalosVersion = "1.4.8"
	cluster2.TypedSpec().Value.Features = &specs.ClusterSpec_Features{
		UseEmbeddedDiscoveryService: true,
	}

	suite.Require().NoError(suite.state.Create(suite.ctx, cluster2))

	// assert that the cluster with the unsupported version has the config patch created with empty data.
	rtestutils.AssertResource[*omni.ConfigPatch](suite.ctx, suite.T(), suite.state, unsupportedNewClusterPatchID, func(r *omni.ConfigPatch, assertion *assert.Assertions) {
		assertion.Empty(r.TypedSpec().Value.Data)
	})

	// update cluster with the unsupported Talos version to a supported version
	_, err = safe.StateUpdateWithConflicts[*omni.Cluster](suite.ctx, suite.state, cluster2.Metadata(), func(res *omni.Cluster) error {
		res.TypedSpec().Value.TalosVersion = "1.6.1"

		return nil
	})
	suite.Require().NoError(err)

	// assert that the config patch was updated to use the discovery service endpoint
	rtestutils.AssertResource[*omni.ConfigPatch](suite.ctx, suite.T(), suite.state, unsupportedNewClusterPatchID, func(r *omni.ConfigPatch, assertion *assert.Assertions) {
		endpoint := fmt.Sprintf("http://" + net.JoinHostPort(siderolink.ListenHost, strconv.Itoa(port)))

		assertion.Contains(r.TypedSpec().Value.Data, endpoint)
	})
}

func TestDiscoveryServiceConfigPatchSuite(t *testing.T) {
	suite.Run(t, new(DiscoveryServiceConfigPatchSuite))
}

// Copyright (c) 2024 Sidero Labs, Inc.
//
// Use of this software is governed by the Business Source License
// included in the LICENSE file.

package omni

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/blang/semver"
	"github.com/cosi-project/runtime/pkg/controller"
	"github.com/cosi-project/runtime/pkg/controller/generic/qtransform"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"github.com/siderolabs/omni/client/pkg/omni/resources"
	"github.com/siderolabs/omni/client/pkg/omni/resources/omni"
	"github.com/siderolabs/omni/internal/pkg/siderolink"
)

// DiscoveryServiceConfigPatchPrefix is the prefix for the system config patch that contains the discovery service endpoint config.
const DiscoveryServiceConfigPatchPrefix = "950-discovery-service-"

// DiscoveryServiceConfigPatchController creates the system config patch that contains the discovery service endpoint config.
//
// It sets the cluster discovery service endpoint to the embedded discovery service endpoint if the cluster is new and the embedded discovery service feature is enabled.
type DiscoveryServiceConfigPatchController = qtransform.QController[*omni.Cluster, *omni.ConfigPatch]

// NewDiscoveryServiceConfigPatchController initializes DiscoveryServiceConfigPatchController.
func NewDiscoveryServiceConfigPatchController(enabled bool, port int) *DiscoveryServiceConfigPatchController {
	supportedTalosVersion := semver.MustParse("1.5.0")

	endpoint := "http://" + net.JoinHostPort(siderolink.ListenHost, strconv.Itoa(port))
	patch := map[string]any{
		"cluster": map[string]any{
			"discovery": map[string]any{
				"registries": map[string]any{
					"service": map[string]any{
						"endpoint": endpoint,
					},
				},
			},
		},
	}

	return qtransform.NewQController(
		qtransform.Settings[*omni.Cluster, *omni.ConfigPatch]{
			Name: "DiscoveryServiceConfigPatchController",
			MapMetadataFunc: func(cluster *omni.Cluster) *omni.ConfigPatch {
				return omni.NewConfigPatch(resources.DefaultNamespace, DiscoveryServiceConfigPatchPrefix+cluster.Metadata().ID())
			},
			UnmapMetadataFunc: func(configPatch *omni.ConfigPatch) *omni.Cluster {
				id := strings.TrimPrefix(configPatch.Metadata().ID(), DiscoveryServiceConfigPatchPrefix)

				return omni.NewCluster(resources.DefaultNamespace, id)
			},
			TransformFunc: func(_ context.Context, _ controller.Reader, logger *zap.Logger, cluster *omni.Cluster, configPatch *omni.ConfigPatch) error {
				configPatch.Metadata().Labels().Set(omni.LabelSystemPatch, "")
				configPatch.Metadata().Labels().Set(omni.LabelCluster, cluster.Metadata().ID())

				if !enabled || !cluster.TypedSpec().Value.GetFeatures().GetUseEmbeddedDiscoveryService() {
					configPatch.TypedSpec().Value.Data = ""

					return nil
				}

				clusterVersion, err := semver.ParseTolerant(cluster.TypedSpec().Value.TalosVersion)
				if err != nil {
					return fmt.Errorf("failed to parse cluster version: %w", err)
				}

				if clusterVersion.LT(supportedTalosVersion) {
					logger.Info("cluster version does not support discovery service endpoint with http:// prefix", zap.String("version", cluster.TypedSpec().Value.TalosVersion))

					configPatch.TypedSpec().Value.Data = ""

					return nil
				}

				var sb strings.Builder

				encoder := yaml.NewEncoder(&sb)

				encoder.SetIndent(2)

				if err = encoder.Encode(patch); err != nil {
					return fmt.Errorf("failed to encode patch: %w", err)
				}

				configPatch.TypedSpec().Value.Data = sb.String()

				return nil
			},
		},
		qtransform.WithConcurrency(2),
		qtransform.WithOutputKind(controller.OutputShared),
	)
}

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package provision

import (
	"context"
	"slices"
	"strings"

	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/siderolabs/gen/xslices"
	"github.com/siderolabs/image-factory/pkg/schematic"
	"go.uber.org/zap"

	"github.com/siderolabs/omni/client/api/omni/specs"
	"github.com/siderolabs/omni/client/pkg/omni/resources/infra"
	"github.com/siderolabs/omni/client/pkg/omni/resources/omni"
)

// FactoryClient ensures that the given schematic exists in the image factory.
type FactoryClient interface {
	EnsureSchematic(context.Context, schematic.Schematic) (string, error)
}

// SchematicOptions is used during schematic ID generation.
type SchematicOptions struct {
	overlay              *schematic.Overlay
	kernelArgs           []string
	extensions           []string
	metaValues           []schematic.MetaValue
	skipConnectionParams bool
}

// SchematicOption is the optional argument to the GetSchematicID method.
type SchematicOption func(*SchematicOptions)

// WithoutConnectionParams generates the schematic without embedding connection params into the kernel args.
// This flag might be useful for providers which use PXE to boot the machines, so the schematics won't need
// the parameters for Omni connection. This can allow to minimize the amount of schematics needed to be generated for the provider.
func WithoutConnectionParams() SchematicOption {
	return func(so *SchematicOptions) {
		so.skipConnectionParams = true
	}
}

// WithExtraExtensions adds more extensions to the schematic.
// The provider can detect the hardware and install some extensions automatically using this method.
func WithExtraExtensions(extensions ...string) SchematicOption {
	return func(so *SchematicOptions) {
		so.extensions = extensions
	}
}

// WithMetaValues adds meta values to the generated schematic.
// If the meta values with the same names are already set they are overwritten.
func WithMetaValues(values ...schematic.MetaValue) SchematicOption {
	return func(so *SchematicOptions) {
		so.metaValues = values
	}
}

// WithExtraKernelArgs adds kernel args to the schematic.
// This method doesn't remove duplicate kernel arguments.
func WithExtraKernelArgs(args ...string) SchematicOption {
	return func(so *SchematicOptions) {
		so.kernelArgs = args
	}
}

// WithOverlay sets the overlay on the schematic.
func WithOverlay(overlay schematic.Overlay) SchematicOption {
	return func(so *SchematicOptions) {
		so.overlay = &overlay
	}
}

// NewContext creates a new provision context.
func NewContext[T resource.Resource](
	machineRequest *infra.MachineRequest,
	machineRequestStatus *infra.MachineRequestStatus,
	state T,
	connectionParams string,
	imageFactory FactoryClient,
) Context[T] {
	return Context[T]{
		machineRequest:       machineRequest,
		MachineRequestStatus: machineRequestStatus,
		State:                state,
		ConnectionParams:     connectionParams,
		imageFactory:         imageFactory,
	}
}

// Context keeps all context which might be required for the provision calls.
type Context[T resource.Resource] struct {
	machineRequest       *infra.MachineRequest
	imageFactory         FactoryClient
	MachineRequestStatus *infra.MachineRequestStatus
	State                T
	ConnectionParams     string
}

// GetRequestID returns machine request id.
func (context *Context[T]) GetRequestID() string {
	return context.machineRequest.Metadata().ID()
}

// GetTalosVersion returns Talos version from the machine request.
func (context *Context[T]) GetTalosVersion() string {
	return context.machineRequest.TypedSpec().Value.TalosVersion
}

// SetMachineUUID in the machine request status.
func (context *Context[T]) SetMachineUUID(value string) {
	context.MachineRequestStatus.TypedSpec().Value.Id = value
}

// SetMachineInfraID in the machine request status.
func (context *Context[T]) SetMachineInfraID(value string) {
	context.MachineRequestStatus.Metadata().Labels().Set(omni.LabelMachineInfraID, value)
}

// GenerateSchematicID generate the final schematic out of the machine request.
// This method also calls the image factory and uploads the schematic there.
func (context *Context[T]) GenerateSchematicID(ctx context.Context, logger *zap.Logger, opts ...SchematicOption) (string, error) {
	var schematicOptions SchematicOptions

	for _, o := range opts {
		o(&schematicOptions)
	}

	res := schematic.Schematic{
		Customization: schematic.Customization{
			ExtraKernelArgs: context.machineRequest.TypedSpec().Value.KernelArgs,
			Meta: xslices.Map(context.machineRequest.TypedSpec().Value.MetaValues, func(v *specs.MetaValue) schematic.MetaValue {
				return schematic.MetaValue{
					Key:   uint8(v.Key),
					Value: v.Value,
				}
			}),
			SystemExtensions: schematic.SystemExtensions{
				OfficialExtensions: context.machineRequest.TypedSpec().Value.Extensions,
			},
		},
	}

	for _, extension := range schematicOptions.extensions {
		if slices.Index(res.Customization.SystemExtensions.OfficialExtensions, extension) != -1 {
			continue
		}

		res.Customization.SystemExtensions.OfficialExtensions = append(res.Customization.SystemExtensions.OfficialExtensions, extension)
	}

	slices.Sort(res.Customization.SystemExtensions.OfficialExtensions)

	for _, metaValue := range schematicOptions.metaValues {
		index := slices.IndexFunc(res.Customization.Meta, func(v schematic.MetaValue) bool {
			return v.Key == metaValue.Key
		})

		if index == -1 {
			res.Customization.Meta = append(res.Customization.Meta, metaValue)

			continue
		}

		res.Customization.Meta[index] = metaValue
	}

	switch {
	case schematicOptions.overlay != nil:
		res.Overlay = *schematicOptions.overlay
	case context.machineRequest.TypedSpec().Value.Overlay != nil:
		res.Overlay = schematic.Overlay{
			Image: context.machineRequest.TypedSpec().Value.Overlay.Image,
			Name:  context.machineRequest.TypedSpec().Value.Overlay.Name,
		}
	}

	if !schematicOptions.skipConnectionParams {
		res.Customization.ExtraKernelArgs = append(
			res.Customization.ExtraKernelArgs,
			strings.Split(context.ConnectionParams, " ")...,
		)
	}

	slices.Sort(res.Customization.ExtraKernelArgs)

	res.Customization.ExtraKernelArgs = append(res.Customization.ExtraKernelArgs, schematicOptions.kernelArgs...)

	logger.Info("creating schematic", zap.Reflect("schematic", res))

	return context.imageFactory.EnsureSchematic(ctx, res)
}
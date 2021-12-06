package processors

import (
	"context"
	"fmt"

	prisma_api "github.com/aporeto-se/prisma-sdk-go-v2/api"
	prisma_types "github.com/aporeto-se/prisma-sdk-go-v2/types"
	"go.uber.org/zap"

	"github.com/aporeto-se/cloud-operator/common/types"
)

// NamespaceProcessor processor
type NamespaceProcessor struct {
	kube                []string
	compute             []string
	cloudOperatorConfig *types.CloudOperatorConfig
	prismaClient        *prisma_api.Client
}

// NewNamespaceProcessor returns new entity instance
func NewNamespaceProcessor(cloudOperatorConfig *types.CloudOperatorConfig, prismaClient *prisma_api.Client) (*NamespaceProcessor, error) {

	if cloudOperatorConfig == nil {
		return nil, fmt.Errorf("entity CloudOperatorConfig is required")
	}

	if prismaClient == nil {
		return nil, fmt.Errorf("entity PrismaClient is required")
	}

	return &NamespaceProcessor{
		cloudOperatorConfig: cloudOperatorConfig,
		prismaClient:        prismaClient,
	}, nil
}

func (t *NamespaceProcessor) has(v string) bool {

	for _, namespace := range t.kube {
		if namespace == v {
			return true
		}
	}

	for _, namespace := range t.compute {
		if namespace == v {
			return true
		}
	}

	return false
}

// AddKube ...
func (t *NamespaceProcessor) AddKube(v ...string) {
	t.kube = append(t.kube, v...)
}

// AddCompute ...
func (t *NamespaceProcessor) AddCompute(v ...string) {
	t.compute = append(t.compute, v...)
}

// CloudEntityTypeInvalid CloudEntityType = "INVALID"
// CloudEntityTypeCompute CloudEntityType = "COMPUTE"
// CloudEntityTypeKubernetes CloudEntityType = "KUBERNETES"
// CloudEntityTypeDefault CloudEntityType = "DEFAULT"

// Process returns Report
func (t *NamespaceProcessor) Process(ctx context.Context) *types.NamespaceReports {

	zap.L().Debug("entering Process")

	report := types.NewNamespaceReports()

	namespaceDeleteEnabled := false
	rogueNamespaceDeleteEnabled := false
	computeNamespaceDeleteEnabled := false
	kubeNamespaceDeleteEnabled := false

	if t.cloudOperatorConfig.HasOp(types.OpNamespaceRogueDelete) {
		zap.L().Debug(fmt.Sprintf("Op %s is enabled", types.OpNamespaceComputeDelete))
		namespaceDeleteEnabled = true
		rogueNamespaceDeleteEnabled = true
	} else {
		zap.L().Debug(fmt.Sprintf("Op %s is disabled", types.OpNamespaceComputeDelete))
	}

	if t.cloudOperatorConfig.HasOp(types.OpNamespaceComputeDelete) {
		zap.L().Debug(fmt.Sprintf("Op %s is enabled", types.OpNamespaceComputeDelete))
		namespaceDeleteEnabled = true
		computeNamespaceDeleteEnabled = true
	} else {
		zap.L().Debug(fmt.Sprintf("Op %s is disabled", types.OpNamespaceComputeDelete))
	}

	if t.cloudOperatorConfig.HasOp(types.OpNamespaceKubeDelete) {
		zap.L().Debug(fmt.Sprintf("Op %s is enabled", types.OpNamespaceKubeDelete))
		namespaceDeleteEnabled = true
		kubeNamespaceDeleteEnabled = true
	} else {
		zap.L().Debug(fmt.Sprintf("Op %s is disabled", types.OpNamespaceKubeDelete))
	}

	if namespaceDeleteEnabled {

		for _, namespace := range t.prismaClient.GetNamespaces() {

			namespaceType, err := t.namespaceType(namespace)

			if err != nil {
				report.AddNamespaces(types.NewNamespaceReport(namespace.Name).
					SetOperation(types.NamespaceOperationDelete).
					SetStatus(types.OperationStatusFailed))
				continue
			}

			switch namespaceType {

			case types.CloudEntityTypeDefault:

				if t.has(namespace.Name) {
					zap.L().Debug(fmt.Sprintf("namespace %s type %s is in use", namespace.Name, namespaceType))
				} else {

					zap.L().Debug(fmt.Sprintf("namespace %s type %s is in not in use", namespace.Name, namespaceType))

					if rogueNamespaceDeleteEnabled {
						report.AddNamespaces(t.namespaceDeleteReport(ctx, namespace, namespaceType))
					}
				}

			case types.CloudEntityTypeCompute:

				if t.has(namespace.Name) {
					zap.L().Debug(fmt.Sprintf("namespace %s type %s is in use", namespace.Name, namespaceType))
				} else {

					zap.L().Debug(fmt.Sprintf("namespace %s type %s is in not in use", namespace.Name, namespaceType))

					if computeNamespaceDeleteEnabled {
						report.AddNamespaces(t.namespaceDeleteReport(ctx, namespace, namespaceType))
					}
				}

			case types.CloudEntityTypeKubernetes:

				if t.has(namespace.Name) {
					zap.L().Debug(fmt.Sprintf("compute namespace %s is active not be deleted", namespace.Name))
				} else {
					if kubeNamespaceDeleteEnabled {
						report.AddNamespaces(t.namespaceDeleteReport(ctx, namespace, namespaceType))
					}
				}

			}

		}

	}

	if t.cloudOperatorConfig.HasOp(types.OpNamespaceComputeCreate) {

		zap.L().Debug(fmt.Sprintf("Op %s is enabled", types.OpNamespaceComputeCreate))

		for _, namespace := range t.compute {
			report.AddNamespaces(t.namespaceCreateReport(ctx, namespace, types.CloudEntityTypeCompute))
		}
	} else {
		zap.L().Debug(fmt.Sprintf("Op %s is disabled", types.OpNamespaceComputeCreate))
	}

	if t.cloudOperatorConfig.HasOp(types.OpNamespaceKubeCreate) {

		zap.L().Debug(fmt.Sprintf("Op %s is enabled", types.OpNamespaceKubeCreate))

		for _, namespace := range t.kube {
			report.AddNamespaces(t.namespaceCreateReport(ctx, namespace, types.CloudEntityTypeKubernetes))
		}

	} else {
		zap.L().Debug(fmt.Sprintf("Op %s is disabled", types.OpNamespaceKubeCreate))
	}

	zap.L().Debug("returning Process")
	return report
}

func (t *NamespaceProcessor) namespaceCreateReport(ctx context.Context, name string, ptype types.CloudEntityType) *types.NamespaceReport {

	zap.L().Debug("entering namespaceCreateReport")

	// Set status to failed. This is so we an bail on error. If and when the status
	// changes we will update it. We also set the type to Compute. We will update if
	// necessary.
	report := types.NewNamespaceReport(name).
		SetOperation(types.NamespaceOperationCreate).
		SetStatus(types.OperationStatusFailed).
		SetType(ptype)

	if t.prismaClient.HasNamespace(name) {
		zap.L().Debug(fmt.Sprintf("namespace %s already exist", name))
		return report.SetStatus(types.OperationStatusAlreadyExist)
	}

	zap.L().Debug(fmt.Sprintf("namespace %s does not exist; creating", name))
	_, err := t.prismaClient.CreateNamespace(ctx,
		prisma_types.NewNamespace(name).
			SetNamespaceType(prisma_types.NamespaceTypeGroup).
			AddAnnotation(namespaceAnnotationKey, []string{string(ptype)}).
			SetDefaultPUIncomingTrafficAction(prisma_types.TrafficActionInherit).
			SetDefaultPUOutgoingTrafficAction(prisma_types.TrafficActionInherit))

	if err != nil {
		zap.L().Debug("returning namespaceCreateReport with error(s)")
		return report.SetError(err)
	}

	zap.L().Debug("returning namespaceCreateReport")
	return report.SetStatus(types.OperationStatusCompleted)
}

func (t *NamespaceProcessor) namespaceDeleteReport(ctx context.Context, namespace *prisma_types.Namespace, ptype types.CloudEntityType) *types.NamespaceReport {

	zap.L().Debug("entering namespaceDeleteReport")

	// Set status to failed. This is so we an bail on error. If and when the status
	// changes we will update it. We also set the type to Compute. We will update if
	// necessary.
	report := types.NewNamespaceReport(namespace.Name).
		SetOperation(types.NamespaceOperationDelete).
		SetStatus(types.OperationStatusFailed).
		SetType(ptype)

	zap.L().Debug(fmt.Sprintf("deleting namespace %s", namespace.Name))
	err := t.prismaClient.DeleteNamespace(ctx, namespace.Name)

	if err != nil {
		zap.L().Debug("returning namespaceDeleteReport with error(s)")
		return report.SetError(err)
	}

	zap.L().Debug("returning namespaceDeleteReport")
	return report.SetStatus(types.OperationStatusCompleted)
}

func (t *NamespaceProcessor) namespaceType(namespace *prisma_types.Namespace) (types.CloudEntityType, error) {

	if namespace.Annotations == nil {
		return types.CloudEntityTypeDefault, nil
	}

	entry := namespace.Annotations[namespaceAnnotationKey]

	if entry == nil {
		return types.CloudEntityTypeDefault, nil
	}

	for _, v := range entry {

		ptype, err := types.CloudEntityTypeFromString(v)

		if err != nil {
			return types.CloudEntityTypeInvalid, err
		}

		return ptype, nil
	}

	return types.CloudEntityTypeInvalid, fmt.Errorf("namespace %s has an unexpected and invalid annotation", namespace.Name)
}

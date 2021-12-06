package types

import (
	"fmt"
	"strings"
)

// ================================================================================================

// LogLevel is the log level
type LogLevel string

const (

	// LogLevelInvalid is invalid
	LogLevelInvalid LogLevel = "Invalid"

	// LogLevelError is the Error LogLevel
	LogLevelError LogLevel = "ERROR"

	// LogLevelWarn is the Warn LogLevel
	LogLevelWarn LogLevel = "WARN"

	// LogLevelInfo is the Info LogLevel
	LogLevelInfo LogLevel = "INFO"

	// LogLevelDebug is the Debug LogLevel
	LogLevelDebug LogLevel = "DEBUG"
)

// LogLevelFromString returns the LogLevel Enum from the provided string. If the
// string is not a valid LogLevel then LogLevelInfo and an error will be returned.
func LogLevelFromString(s string) (LogLevel, error) {

	switch strings.ToUpper(s) {
	case "ERROR":
		return LogLevelError, nil

	case "WARN":
		return LogLevelWarn, nil

	case "INFO":
		return LogLevelInfo, nil

	case "DEBUG":
		return LogLevelDebug, nil

	}

	return LogLevelInvalid, fmt.Errorf("string %s is not a valid type", s)
}

// ================================================================================================

// CloudEntityType is the type of cloud entity such as compute (vm/host) or Kubernetes.
type CloudEntityType string

const (
	// CloudEntityTypeInvalid invalid
	CloudEntityTypeInvalid CloudEntityType = "INVALID"

	// CloudEntityTypeCompute compute
	CloudEntityTypeCompute CloudEntityType = "COMPUTE"

	// CloudEntityTypeKubernetes kubernetes
	CloudEntityTypeKubernetes CloudEntityType = "KUBERNETES"

	// CloudEntityTypeDefault default
	CloudEntityTypeDefault CloudEntityType = "DEFAULT"
)

// CloudEntityTypeFromString returns type from string or error
func CloudEntityTypeFromString(s string) (CloudEntityType, error) {

	switch strings.ToUpper(s) {

	case string(CloudEntityTypeCompute):
		return CloudEntityTypeCompute, nil

	case string(CloudEntityTypeKubernetes):
		return CloudEntityTypeKubernetes, nil

	case string(CloudEntityTypeDefault):
		return CloudEntityTypeDefault, nil

	}

	return CloudEntityTypeInvalid, fmt.Errorf("string %s is not a valid type", s)
}

// ================================================================================================

// Op operations
type Op string

const (
	// OpInvalid invalid
	OpInvalid Op = "INVALID"

	// OpDHCP DHCP
	OpDHCP Op = "DHCP"

	// OpNamespaceRogueDelete Namespace Rogue Delete
	OpNamespaceRogueDelete Op = "NS_ROGUE_DELETE"

	// OpNamespaceComputeCreate Namespace Compute Create
	OpNamespaceComputeCreate Op = "NS_COMPUTE_CREATE"

	// OpNamespaceComputeDelete Namespace Compute Delete
	OpNamespaceComputeDelete Op = "NS_COMPUTE_DELETE"

	// OpNamespaceKubeCreate Namespace Kube Create
	OpNamespaceKubeCreate Op = "NS_KUBE_CREATE"

	// OpNamespaceKubeDelete Namespace Kube Delete
	OpNamespaceKubeDelete Op = "NS_KUBE_DELETE"

	// OpComputeAuth Compute Auth
	OpComputeAuth Op = "COMPUTE_AUTH"

	// OpKubeAuth Kube Auth
	OpKubeAuth Op = "KUBE_AUTH"

	// OpKubeAPINet API Network List
	OpKubeAPINet Op = "KUBE_API_NET"

	// OpKubeDNSNet DNS Network List
	OpKubeDNSNet Op = "KUBE_DNS_NET"

	// OpKubeNodesNet Nodes (Workers) Network List
	OpKubeNodesNet Op = "KUBE_NODES_NET"

	// OpKubeEnforcer Enforcer (Daemonset) Install
	OpKubeEnforcer Op = "KUBE_ENFORCER"
)

// OpFromString returns type from string or error
func OpFromString(s string) (Op, error) {

	switch strings.ToUpper(s) {

	case string(OpDHCP):
		return OpDHCP, nil

	case string(OpNamespaceRogueDelete):
		return OpNamespaceRogueDelete, nil

	case string(OpNamespaceComputeCreate):
		return OpNamespaceComputeCreate, nil

	case string(OpNamespaceComputeDelete):
		return OpNamespaceComputeDelete, nil

	case string(OpNamespaceKubeCreate):
		return OpNamespaceKubeCreate, nil

	case string(OpNamespaceKubeDelete):
		return OpNamespaceKubeDelete, nil

	case string(OpComputeAuth):
		return OpComputeAuth, nil

	case string(OpKubeAuth):
		return OpKubeAuth, nil

	case string(OpKubeAPINet):
		return OpKubeAPINet, nil

	case string(OpKubeDNSNet):
		return OpKubeDNSNet, nil

	case string(OpKubeNodesNet):
		return OpKubeNodesNet, nil

	case string(OpKubeEnforcer):
		return OpKubeEnforcer, nil

	}

	return OpInvalid, fmt.Errorf("string %s is not a valid type", s)
}

// ================================================================================================

// OpStatus status
type OpStatus string

const (
	// OpStatusInvalid invalid
	OpStatusInvalid OpStatus = "INVALID"

	// OpStatusCompleted Completed
	OpStatusCompleted OpStatus = "COMPLETED/CREATED"

	// OpStatusFailed failed
	OpStatusFailed OpStatus = "FAILED"

	// OpStatusNothingToDo Nothing To Do
	OpStatusNothingToDo OpStatus = "ALREADY_EXIST/NOTHING_TO_DO"

	// OpStatusNotReady failed
	OpStatusNotReady OpStatus = "NOT_READY"
)

// OpStatusFromString returns type from string or error
func OpStatusFromString(s string) (OpStatus, error) {

	switch strings.ToUpper(s) {

	case string(OpStatusCompleted):
		return OpStatusCompleted, nil

	case string(OpStatusFailed):
		return OpStatusFailed, nil

	case string(OpStatusNothingToDo):
		return OpStatusNothingToDo, nil

	case string(OpStatusNotReady):
		return OpStatusNotReady, nil

	}

	return OpStatusInvalid, fmt.Errorf("string %s is not a valid type", s)
}

// ================================================================================================

// NamespaceOperation NamespaceOperation
type NamespaceOperation string

const (
	// NamespaceOperationInvalid invalid
	NamespaceOperationInvalid NamespaceOperation = "INVALID"

	// NamespaceOperationCreate create
	NamespaceOperationCreate NamespaceOperation = "CREATE"

	// NamespaceOperationDelete delete
	NamespaceOperationDelete NamespaceOperation = "DELETE"

	// NamespaceOperationIgnore ignore
	NamespaceOperationIgnore NamespaceOperation = "IGNORE"
)

// NamespaceOperationFromString returns type from string or error
func NamespaceOperationFromString(s string) (NamespaceOperation, error) {

	switch strings.ToUpper(s) {

	case string(NamespaceOperationCreate):
		return NamespaceOperationCreate, nil

	case string(NamespaceOperationDelete):
		return NamespaceOperationDelete, nil

	case string(NamespaceOperationIgnore):
		return NamespaceOperationIgnore, nil

	}

	return NamespaceOperationInvalid, fmt.Errorf("string %s is not a valid type", s)
}

// ================================================================================================

// OperationStatus OperationStatus
type OperationStatus string

const (
	// OperationStatusInvalid invalid
	OperationStatusInvalid OperationStatus = "INVALID"

	// OperationStatusCompleted completed
	OperationStatusCompleted OperationStatus = "COMPLETED"

	// OperationStatusAlreadyExist already exist
	OperationStatusAlreadyExist OperationStatus = "ALREADY_EXIST/NOTHING_TO_DO"

	// OperationStatusFailed failed
	OperationStatusFailed OperationStatus = "FAILED"
)

// OperationStatusromString returns type from string or error
func OperationStatusromString(s string) (OperationStatus, error) {

	switch strings.ToUpper(s) {

	case string(OperationStatusCompleted):
		return OperationStatusCompleted, nil

	case string(OperationStatusAlreadyExist):
		return OperationStatusAlreadyExist, nil

	case string(OperationStatusFailed):
		return OperationStatusFailed, nil

	}

	return OperationStatusInvalid, fmt.Errorf("string %s is not a valid type", s)
}

// ================================================================================================

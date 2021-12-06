package operator

import (
	logging "github.com/aporeto-se/cloud-operator/common/logging"
	"github.com/aporeto-se/cloud-operator/common/types"
)

// InitLogging delegation for operator InitLogging
func InitLogging(logLevel types.LogLevel) error {
	return logging.InitLogging(logLevel)
}

package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/aporeto-se/cloud-operator/common/types"
)

// InitLogging initializes logging with specified log level. If not successful an error is returned.
func InitLogging(logLevel types.LogLevel) error {

	zapConfig := &zap.Config{
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		EncoderConfig: zap.NewProductionEncoderConfig(),
	}

	zapConfig.Encoding = "json"
	zapConfig.OutputPaths = append(zapConfig.OutputPaths, "stderr")
	zapConfig.ErrorOutputPaths = append(zapConfig.ErrorOutputPaths, "stderr")

	switch logLevel {

	case types.LogLevelError:
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)

	case types.LogLevelWarn:
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)

	case types.LogLevelInfo:
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	case types.LogLevelDebug:
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)

	default:
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	}

	logger, err := zapConfig.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(logger)

	return nil
}

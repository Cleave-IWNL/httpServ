package logger

import "go.uber.org/zap"

func New(level string) (*zap.Logger, error) {
	atomicLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = atomicLevel

	return cfg.Build()
}

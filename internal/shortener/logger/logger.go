package logger

import "go.uber.org/zap"

// Initialize missing godoc.
func Initialize(level string) (*zap.Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	zCnf := zap.NewProductionConfig()
	zCnf.Level = lvl
	zLogger, err := zCnf.Build()
	if err != nil {
		return nil, err
	}

	return zLogger, nil
}

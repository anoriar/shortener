package logger

import "go.uber.org/zap"

// MENTOR: Целесообразно ли делать свою обертку - логгер над запом?
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

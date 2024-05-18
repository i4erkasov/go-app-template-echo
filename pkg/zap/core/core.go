package core

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Core interface {
	Create(cfg *viper.Viper) (zapcore.Core, error)
}

func Create(cfg *viper.Viper, path string) (zapcore.Core, error) {
	cfgCore := cfg.Sub(path)
	if cfgCore == nil {
		return nil, fmt.Errorf("core config at path '%s' not found", path)
	}

	cfgCore.SetDefault("type", TypeCoreStream)
	cfgCore.SetDefault("level", "info")

	var core Core
	coreType := cfgCore.GetString("type")
	switch coreType {
	case TypeCoreStream:
		core = &StreamCore{}
	case TypeCoreFile:
		core = &FileCore{}
	default:
		return nil, fmt.Errorf("unsupported core type: %s", coreType)
	}

	return core.Create(cfgCore)
}

func getEncoder(cfg *viper.Viper) (zapcore.Encoder, error) {
	cfgEncoder := zap.NewProductionEncoderConfig()
	if cfg.IsSet("zap.development") && cfg.GetBool("zap.development") {
		cfgEncoder = zap.NewDevelopmentEncoderConfig()
	}

	if cfg.IsSet("time_encoder") {
		if err := cfgEncoder.EncodeTime.UnmarshalText(
			[]byte(cfg.GetString("time_encoder")),
		); err != nil {
			return nil, err
		}
	}

	if cfg.IsSet("message_key") {
		cfgEncoder.MessageKey = cfg.GetString("message_key")
	}

	var encoding = "json"
	if cfg.IsSet("encoding") {
		encoding = cfg.GetString("encoding")
	}

	var encoder zapcore.Encoder
	switch encoding {
	case "json":
		encoder = zapcore.NewJSONEncoder(cfgEncoder)
	case "console":
		encoder = zapcore.NewConsoleEncoder(cfgEncoder)
	default:
		return nil, fmt.Errorf(`encoding "%s" is not supported`, encoding)
	}

	return encoder, nil
}

func getLevel(cfg *viper.Viper) (zapcore.LevelEnabler, error) {
	var level = zap.NewAtomicLevel()
	if cfg.IsSet("level") {
		if err := level.UnmarshalText([]byte(cfg.GetString("level"))); err != nil {
			return nil, err
		}
	}

	return level, nil
}

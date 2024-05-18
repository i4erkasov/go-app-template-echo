package zap

import (
	"github.com/i4erkasov/go-ddd-cqrs/pkg/zap/core"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(cfg *viper.Viper) (*zap.Logger, error) {
	// Retrieve logger core configuration from Viper
	cfgCore := cfg.Sub("zap.cores")
	settings := cfgCore.AllSettings()
	cores := make([]zapcore.Core, 0, len(settings))

	for name := range settings {
		c, err := core.Create(cfg, "zap.cores."+name)
		if err != nil {
			return nil, err
		}

		cores = append(cores, c)
	}

	// Add logger options
	var opts []zap.Option
	if cfg.IsSet("zap.development") && cfg.GetBool("zap.development") {
		opts = append(opts, zap.Development())
	}

	if cfg.GetBool("zap.caller") {
		opts = append(opts, zap.AddCaller())
	}
	if cfg.IsSet("zap.stacktrace") {
		var level = zap.NewAtomicLevel()
		if err := level.UnmarshalText([]byte(cfg.GetString("zap.stacktrace"))); err != nil {
			return nil, err
		}

		opts = append(opts, zap.AddStacktrace(level))
	}

	return zap.New(zapcore.NewTee(cores...), opts...), nil
}

package core

import (
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type FileCore struct{}

const TypeCoreFile = "file"

const (
	defaultPath       = "/var/log/myapp" // Base path to the file without date
	defaultMaxBackups = 3                // Maximum number of backup files
	defaultMaxAge     = 30               // Maximum retention period in days
	defaultMaxSize    = 32               // Maximum file size in MB
)

func (c *FileCore) Create(cfg *viper.Viper) (_ zapcore.Core, err error) {
	var encoder zapcore.Encoder
	if encoder, err = getEncoder(cfg); err != nil {
		return nil, err
	}

	var level zapcore.LevelEnabler
	if level, err = getLevel(cfg); err != nil {
		return nil, err
	}

	cfg.SetDefault("file.path", defaultPath)
	cfg.SetDefault("file.max_backups", defaultMaxBackups)
	cfg.SetDefault("file.max_age", defaultMaxAge)   // Base path to the file without date
	cfg.SetDefault("file.max_size", defaultMaxSize) // Base path to the file without date

	// Формирование пути файла с текущей датой
	today := time.Now().Format("2006-01-02")
	fileName := cfg.GetString("path") + "-" + today + ".log"

	return zapcore.NewCore(encoder, zapcore.AddSync(&lumberjack.Logger{
		Filename:   fileName,
		MaxBackups: cfg.GetInt("max_backups"), // How many old versions of logs to keep
		MaxAge:     cfg.GetInt("max_age"),     // Delete old versions older than this many days
		MaxSize:    cfg.GetInt("max_size"),    // Maximum size in megabytes
	}), level), nil
}

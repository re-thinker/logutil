// Package logutil 日志基本库
package logutil

import (
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

const (
	defaultLogTimeFormat = "2006/01/02 15:04:05.000"
	defaultLogMaxSize    = 10 // MB
	defaultLogFormat     = "text"
	defaultLogLevel      = log.InfoLevel
)

// LogConfig 日志文件配置
type LogConfig struct {
	Filename   string `json:"filename"`    // 日志文件名
	LogRotate  bool   `json:"log_rotate"`  // 是否进行滚动
	MaxSize    int    `json:"max_size"`    // 单个文件的最大大小，MB
	MaxDays    int    `json:"max_days"`    // 保留的最大天数
	MaxBackups int    `json:"max_backups"` // 最大文件数量
	Level      string `json:"level"`       // 日志级别
	Format     string `json:"format"`      // 日志格式 json, text, console
}

func stringToLogFormatter(format string) log.Formatter {
	switch strings.ToLower(format) {
	/*case "text":
	return &textFormatter{
		DisableTimestamp: true,
	}*/
	case "json":
		return &log.JSONFormatter{
			TimestampFormat:  defaultLogTimeFormat,
			DisableTimestamp: false,
		}
	case "console":
		return &log.TextFormatter{
			FullTimestamp:    true,
			TimestampFormat:  defaultLogTimeFormat,
			DisableTimestamp: false,
		}
	/*case "highlight":
	return &textFormatter{
		DisableTimestamp: disableTimestamp,
		EnableColors:     true,
	}*/
	default:
		return &log.TextFormatter{
			FullTimestamp:    true,
			TimestampFormat:  defaultLogTimeFormat,
			DisableTimestamp: false,
		}
	}
}

func stringToLogLevel(level string) log.Level {
	switch strings.ToLower(level) {
	case "fatel":
		return log.FatalLevel
	case "error":
		return log.ErrorLevel
	case "warn":
		return log.WarnLevel
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	}
	return defaultLogLevel
}

// InitLogger 根据配置信息初使化日志
func InitLogger(cfg *LogConfig) {
	log.SetLevel(stringToLogLevel(cfg.Level))
	if cfg.Format == "" {
		cfg.Format = defaultLogFormat
	}
	formatter := stringToLogFormatter(cfg.Format)
	log.SetFormatter(formatter)
}

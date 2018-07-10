// Package logutil 日志基本库
// 支持日志回滚，json、toml配置格式，信号改变日志级别、日志压缩
package logutil

import (
	"errors"
	"os"
	"strings"

	"github.com/natefinch/lumberjack"
	log "github.com/sirupsen/logrus"
	"os/signal"
	"syscall"
)

const (
	defaultLogTimeFormat = "2006/01/02 15:04:05.000"
	defaultLogMaxSize    = 10 // MB
	defaultLogMaxDays    = 90 // 默认保存90天
	defaultLogMaxBackups = 5  // 默认5个文件循环
	defaultLogFormat     = "json"
	defaultLogLevel      = log.InfoLevel
)

// LogConfig 日志文件配置
type LogConfig struct {
	Filename   string `json:"filename" toml:"filename"`       // 日志文件名
	MaxSize    int    `json:"max_size" toml:"max_size"`       // 单个文件的最大大小，MB
	MaxDays    int    `json:"max_days" toml:"max_days"`       // 保留的最大天数
	MaxBackups int    `json:"max_backups" toml:"max_backups"` // 最大文件数量
	Level      string `json:"level" toml:"level"`             // 日志级别
	Format     string `json:"format" toml:"format"`           // 日志格式 json, text, console
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
	case "fatal":
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

func singalChangeLogLevel() {
	signalUser1 := make(chan os.Signal)
	signalUser2 := make(chan os.Signal)
	signal.Notify(signalUser1, syscall.SIGUSR1)
	signal.Notify(signalUser2, syscall.SIGUSR2)
	for {
		select {
		case <-signalUser1:
			log.SetLevel(log.DebugLevel)
		case <-signalUser2:
			log.SetLevel(log.ErrorLevel)
		}
	}
}

// InitLogger 根据配置信息初使化日志
func InitLogger(cfg *LogConfig) error {
	if st, err := os.Stat(cfg.Filename); err == nil {
		if st.IsDir() {
			return errors.New("can't use directory as log file name")
		}
	}
	log.SetLevel(stringToLogLevel(cfg.Level))
	log.SetFormatter(stringToLogFormatter(cfg.Format))

	if cfg.MaxSize == 0 {
		cfg.MaxSize = defaultLogMaxSize
	}
	if cfg.MaxDays == 0 {
		cfg.MaxDays = defaultLogMaxDays
	}
	if cfg.MaxBackups == 0 {
		cfg.MaxBackups = defaultLogMaxBackups
	}

	// 用lumberjack 进行回滚
	output := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxDays,
		MaxBackups: cfg.MaxBackups,
		LocalTime:  true,
		Compress:   true,
	}
	go singalChangeLogLevel()
	log.SetOutput(output)

	return nil
}

package logutil

import (
	"bufio"
	"io"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestStringToLogLevel(t *testing.T) {
	assert.Equal(t, stringToLogLevel("fatal"), log.FatalLevel)
	assert.Equal(t, stringToLogLevel("error"), log.ErrorLevel)
	assert.Equal(t, stringToLogLevel("Debug"), log.DebugLevel)
	assert.Equal(t, stringToLogLevel("Warn"), log.WarnLevel)
	assert.Equal(t, stringToLogLevel("Info"), log.InfoLevel)
	assert.Equal(t, stringToLogLevel("whatlevel"), log.InfoLevel)
}

func TestStringToLogFormatter(t *testing.T) {
	assert.IsType(t,
		&log.JSONFormatter{
			TimestampFormat:  defaultLogTimeFormat,
			DisableTimestamp: false},
		stringToLogFormatter("json"))
	assert.IsType(t,
		&log.TextFormatter{
			FullTimestamp:    true,
			TimestampFormat:  defaultLogTimeFormat,
			DisableTimestamp: false},
		stringToLogFormatter("console"))
	assert.IsType(t,
		&log.TextFormatter{
			FullTimestamp:    true,
			TimestampFormat:  defaultLogTimeFormat,
			DisableTimestamp: false},
		stringToLogFormatter("waht"))
}

func TestInitLogger(t *testing.T) {
	fileName := "test_log.log"
	logCfg := &LogConfig{
		Filename: fileName,
	}
	defer os.Remove(fileName)
	InitLogger(logCfg)

	log.Info("test")
	log.WithFields(log.Fields{"name": "jack", "sex": "male"}).Info("ok")
	f, err := os.Open(fileName)
	assert.Nil(t, err)
	r := bufio.NewReader(f)
	for {
		var str string
		str, err = r.ReadString('\n')
		if err != nil {
			break
		}
		assert.Contains(t, str, "info")
	}
	assert.Equal(t, io.EOF, err)
}

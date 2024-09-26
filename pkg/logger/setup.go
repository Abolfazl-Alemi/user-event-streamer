package logger

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"strings"
	"user-event-streamer/configs"
)

var Zap zap.Logger

func init() {
	Zap = NewLogger()
}

func NewLogger() zap.Logger {
	// Open the go.mod file
	file, err := os.Open("go.mod")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing go.mod:", err)
		}
	}(file)

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	var moduleName string = "myProject"
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module") {
			// The module name is the second part of the line
			moduleName = strings.Fields(line)[1]
			fmt.Println("Module name:", moduleName)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading go.mod:", err)
	}

	cfgLogger := configs.GetConfig().Logger

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"pid":     os.Getpid(),
			"service": moduleName,
		},
	}

	if cfgLogger.Level == "development" {
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	logger, err := config.Build()

	if err != nil {
		log.Fatal(err)
	}
	return *logger
}

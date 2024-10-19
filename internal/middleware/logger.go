package middleware

import (
	"SkinRest/config"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(cfg *config.Config) *zap.Logger {
	stdout := zapcore.AddSync(os.Stdout)

	var filePath string

	if cfg.Server.ApiEnv == "local" { // set local or release path to config
		filePath = "logs/app.log"
	} else {
		filePath = "/root/logs/app.log"
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// setup configuration settings
	level := zapcore.ErrorLevel

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(file), level),
		zapcore.NewCore(consoleEncoder, stdout, zapcore.DebugLevel),
	)

	return zap.New(core, zap.AddCaller())
}

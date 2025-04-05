package logger

import (
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	Log *zap.SugaredLogger
}

func NewZapLogger() *ZapLogger {
	cfg := zap.Config{
		Encoding:    "console", 
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths: []string{"stderr"},

		EncoderConfig: zapcore.EncoderConfig{
				MessageKey: "message",
				TimeKey:    "time",
				LevelKey:   "level",
				CallerKey:  "caller",
				EncodeCaller: zapcore.FullCallerEncoder,
				EncodeLevel:  CustomLevelEncoder,  
				EncodeTime:   CustomTimeEncoder,  
		},
	}
	
	logger, err := cfg.Build()
	if err != nil {
			log.Fatalf(fmt.Sprintf("failed to build logger: %v", err))
	}

	return &ZapLogger{
		Log: logger.Sugar(),
	}
}
// Formatting Time and Level Logging

func CustomTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(t.Format("2006-01-02 15:04:05"))
}

func CustomLevelEncoder(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString("[" + level.CapitalString() + "]")
}


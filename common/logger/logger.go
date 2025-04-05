package logger

import (
	"fmt"
	"log"
	"strings"
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
				EncodeCaller: ColorCallerEncoder,
				EncodeLevel:  zapcore.CapitalColorLevelEncoder,  
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
	encoder.AppendString("[" + t.Format("15:04:05") + "]")
}

func ColorCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	const cyan = "\033[36m"
	const reset = "\033[0m"
	const root = "Goph-Chat/" // only keep path after this

	shortPath := caller.FullPath()
	// Shorten the caller path to root directory
	if idx := strings.Index(shortPath, root); idx != -1 {
		shortPath = shortPath[idx+len(root):]
	}

	enc.AppendString(cyan + shortPath + reset)
}

// func CustomLevelEncoder(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
// 	encoder.AppendString("[" + level.CapitalString() + "]")
// }


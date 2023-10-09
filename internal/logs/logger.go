package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	*zap.SugaredLogger
}

var NewLogger = newLogger()

func newLogger() *Logger {
	level := zapcore.DebugLevel

	core := zapcore.NewCore(getEncoder(), getWriter(), level)

	log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	defer log.Sync()

	return &Logger{log.Sugar()}
}

func getEncoder() zapcore.Encoder {
	var encoderConfig = zapcore.EncoderConfig{
		MessageKey: "message",

		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,

		TimeKey:    "time",
		EncodeTime: zapcore.ISO8601TimeEncoder,

		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    200,
		MaxAge:     30,
		MaxBackups: 10,
		Compress:   false,
	}

	return zapcore.AddSync(lumberJackLogger)
}

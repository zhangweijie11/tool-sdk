package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"tool-sdk/global"
)

const (
	logTypeFile   = "file"
	logTypeStdout = "stdout"
)

var logger *zap.Logger

func Setup(level string) (err error) {
	var loggerLevel zapcore.Level
	if loggerLevel.UnmarshalText([]byte(level)) != nil {
		loggerLevel = zapcore.DebugLevel
	}

	var writeSyncer zapcore.WriteSyncer

	// 日志写入文件
	if global.Config.Server.LogType == logTypeFile {
		logPath := filepath.Join(global.Config.Server.RootDir, global.Config.Server.LogFilePath)

		lumberJackLogger := &lumberjack.Logger{
			Filename:   logPath, // 文件位置
			MaxSize:    10,      // 进行切割之前,日志文件的最大大小(MB为单位)
			MaxAge:     90,      // 保留旧文件的最大天数
			MaxBackups: 100,     // 保留旧文件的最大个数
			Compress:   false,   // 是否压缩/归档旧文件
		}

		writeSyncer = zapcore.AddSync(lumberJackLogger)
	} else {
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	core := zapcore.NewCore(zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(global.TimeFormatSecond),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}), writeSyncer, loggerLevel)
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return err
}

func Info(message string, fields ...zapcore.Field) {
	logger.Info(message, fields...)
}

func Warn(message string, fields ...zapcore.Field) {
	logger.Warn(message, fields...)
}

func Error(message string, err error, fields ...zapcore.Field) {
	logger.Error(message, append(fields, zap.Error(err))...)
}

func Panic(message string, err error, fields ...zapcore.Field) {
	logger.Panic(message, append(fields, zap.Error(err))...)
}

func Get() *zap.Logger {
	return logger
}

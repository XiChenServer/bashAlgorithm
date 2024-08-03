package logger

import "os"

func NewExample(options ...Option) *Logger {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",                         // 日志内容key:val， 前面的key设为msg
		LevelKey:       "level",                       // 日志级别的key设为level
		NameKey:        "logger",                      // 日志名
		EncodeLevel:    zapcore.LowercaseLevelEncoder, //日志级别，默认小写
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // 日志时间
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), os.Stdout, DebugLevel)
	return New(core).WithOptions(options...)
}

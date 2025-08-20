package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger 日志接口
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
}

// LogrusLogger logrus日志实现
type LogrusLogger struct {
	logger *logrus.Logger
}

// NewLogger 创建新的日志器
func NewLogger(level string) *LogrusLogger {
	logger := logrus.New()

	// 设置输出格式
	logger.SetFormatter(&logrus.JSONFormatter{})

	// 设置输出到标准输出
	logger.SetOutput(os.Stdout)

	// 设置日志级别
	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	return &LogrusLogger{
		logger: logger,
	}
}

// Debug 调试日志
func (l *LogrusLogger) Debug(msg string, fields ...interface{}) {
	l.logger.WithFields(parseFields(fields...)).Debug(msg)
}

// Info 信息日志
func (l *LogrusLogger) Info(msg string, fields ...interface{}) {
	l.logger.WithFields(parseFields(fields...)).Info(msg)
}

// Warn 警告日志
func (l *LogrusLogger) Warn(msg string, fields ...interface{}) {
	l.logger.WithFields(parseFields(fields...)).Warn(msg)
}

// Error 错误日志
func (l *LogrusLogger) Error(msg string, fields ...interface{}) {
	l.logger.WithFields(parseFields(fields...)).Error(msg)
}

// Fatal 致命错误日志
func (l *LogrusLogger) Fatal(msg string, fields ...interface{}) {
	l.logger.WithFields(parseFields(fields...)).Fatal(msg)
}

// parseFields 解析字段
func parseFields(fields ...interface{}) logrus.Fields {
	logrusFields := logrus.Fields{}

	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			if key, ok := fields[i].(string); ok {
				logrusFields[key] = fields[i+1]
			}
		}
	}

	return logrusFields
}

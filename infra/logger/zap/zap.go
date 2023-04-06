package zap

import (
	"os"

	"github.com/sourcenetwork/orbis-go/config"
	"github.com/sourcenetwork/orbis-go/infra/logger"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	level       zapcore.Level
	encoding    string
	sugarLogger *zap.SugaredLogger
}

// New returns a new zap logger.
func New(cfg config.Logger) *zapLogger {

	l := &zapLogger{
		level:    toZapLevel(cfg.Zap.Level),
		encoding: cfg.Zap.Encoding,
	}

	l.config()

	return l
}

func (l *zapLogger) SetLevel(lv string) {

	l.level = toZapLevel(string(lv))
	l.config()
	l.Sync()
}

func (l *zapLogger) Debugf(fmt string, args ...interface{}) {
	l.sugarLogger.Debugf(fmt, args...)
}

func (l *zapLogger) Infof(fmt string, args ...interface{}) {
	l.sugarLogger.Infof(fmt, args...)
}

func (l *zapLogger) Warnf(fmt string, args ...interface{}) {
	l.sugarLogger.Warnf(fmt, args...)
}

func (l *zapLogger) Errorf(fmt string, args ...interface{}) {
	l.sugarLogger.Errorf(fmt, args...)
}

func (l *zapLogger) Panicf(fmt string, args ...interface{}) {
	l.sugarLogger.Panicf(fmt, args...)
}

func (l *zapLogger) Fatalf(fmt string, args ...interface{}) {
	l.sugarLogger.Fatalf(fmt, args...)
}

func (l *zapLogger) Named(name string) logger.Logger {
	return &zapLogger{
		level:       l.level,
		encoding:    l.encoding,
		sugarLogger: l.sugarLogger.Named(name),
	}
}

func (l *zapLogger) Sync() error {
	return l.sugarLogger.Sync()
}

func (l *zapLogger) config() {

	encoder := setupEncoder(l.encoding)
	writer := zapcore.AddSync(os.Stdout)
	core := zapcore.NewCore(encoder, writer, zap.NewAtomicLevelAt(l.level))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	l.sugarLogger = logger.Sugar()
}

func setupEncoder(encoding string) zapcore.Encoder {

	var encoderCfg zapcore.EncoderConfig

	switch encoding {
	case "dev":
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	default:
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	return zapcore.NewJSONEncoder(encoderCfg)
}

func toZapLevel(lv string) zapcore.Level {

	m := map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
		"panic": zapcore.PanicLevel,
		"fatal": zapcore.FatalLevel,
	}

	level, exist := m[lv]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

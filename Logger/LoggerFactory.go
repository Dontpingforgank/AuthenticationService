package Logger

import (
	"fmt"
	"github.com/Dontpingforgank/AuthenticationService/Models"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"time"
)

type LoggerFactory interface {
	NewLogger() (*zap.Logger, func(), error)
}

type loggerFactory struct {
	configuration *Models.Config
}

func (logFactory loggerFactory) NewLogger() (*zap.Logger, func(), error) {
	now := time.Now()
	logFile := path.Join(logFactory.configuration.Zap.Path, fmt.Sprintf("%s.log", now.Format("2006-04-02-15")))

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	pe := zap.NewProductionEncoderConfig()

	fileEncoder := zapcore.NewJSONEncoder(pe)
	pe.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(pe)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(file), logFactory.configuration.Zap.Level),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), logFactory.configuration.Zap.Level),
	)

	log := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))
	c := func() {
		err := log.Sync()
		if err != nil {
			return
		}
		err = file.Close()
		if err != nil {
			return
		}
	}

	return log, c, nil
}

func NewLoggerFactory(configuration *Models.Config) LoggerFactory {
	return &loggerFactory{
		configuration: configuration,
	}
}

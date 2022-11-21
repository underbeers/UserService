package service

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"os"
	"time"
)

const logPermissions = 0o600

func NewLogger(isLocal bool) *zap.SugaredLogger {
	loggerConfig := zap.NewProductionEncoderConfig()
	loggerConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(loggerConfig)
	timeStamp := time.Now().Format("02-01-2006")
	logFile, err := os.OpenFile(fmt.Sprintf("./logs/user-service-%s.log", timeStamp),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, logPermissions)
	if err != nil {
		log.Printf("Error while creating NewLogger %s", err.Error())
	}
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	if isLocal {
		consoleEncoder := zapcore.NewConsoleEncoder(loggerConfig)
		core := zapcore.NewTee(
			zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
		)
		logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
		sugar := logger.Sugar()

		return sugar
	}
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	sugar := logger.Sugar()

	return sugar
}

func (srv *service) error(w http.ResponseWriter, code int, err error, ctx context.Context) {
	srv.respond(w, code, map[string]string{"message": err.Error()})
	u, ok := err.(interface { //nolint: errorlint
		ErrorEx() string
	})
	if !ok {
		srv.Logger.Error(ctx.Value(requestID), err.Error())
	} else {
		srv.Logger.Error(ctx.Value(requestID), u.ErrorEx())
	}
}

func (srv *service) warning(w http.ResponseWriter, code int, err error) {
	srv.respond(w, code, map[string]string{"message": err.Error()})
	u, ok := err.(interface { //nolint: errorlint
		ErrorEx() string
	})
	if !ok {
		srv.Logger.Warn(err.Error())
	} else {
		srv.Logger.Warn(u.ErrorEx())
	}
}

func (srv *service) respond(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			srv.Logger.Errorf("failed to encode json %v", err)
		}
	}
}

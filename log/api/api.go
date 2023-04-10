package api

import (
	"map_chan/log/internal"
	"sync"
)

type LogApi struct {
	sync.Mutex
	logger *internal.Logger
}

var logApi LogApi = LogApi{
	logger: internal.NewLogger(internal.LOG_LEVEL_DBG, 4),
}

func LoggerInit(l uint8, isout uint8, out_ffp string) *internal.Logger {
	return internal.NewLoggerWithOutter(internal.LOG_LEVEL(l), 1, out_ffp)
}

func LoggerDeinit(l *internal.Logger) {
	internal.DelLogger(l)
}

func (log *LogApi) Debugf(fmt string, v ...any) {
	log.Lock()
	defer log.Unlock()

	log.logger.Debugf(fmt, v...)
}

func Debugf(fmt string, v ...any) {
	logApi.Lock()
	defer logApi.Unlock()

	logApi.logger.Debugf(fmt, v...)
}

func (log *LogApi) Infof(fmt string, v ...any) {
	log.Lock()
	defer log.Unlock()

	log.logger.Infof(fmt, v...)
}

func Infof(fmt string, v ...any) {
	logApi.Lock()
	defer logApi.Unlock()

	logApi.logger.Infof(fmt, v...)
}

func (log *LogApi) Warnf(fmt string, v ...any) {
	log.Lock()
	defer log.Unlock()
	log.logger.Warnf(fmt, v...)
}

func Warnf(fmt string, v ...any) {
	logApi.Lock()
	defer logApi.Unlock()
	logApi.logger.Warnf(fmt, v...)
}

func (log *LogApi) Errf(fmt string, v ...any) {
	log.Lock()
	defer log.Unlock()
	log.logger.Errf(fmt, v...)
}

func Errf(fmt string, v ...any) {
	logApi.Lock()
	defer logApi.Unlock()
	logApi.logger.Errf(fmt, v...)
}

func (log *LogApi) Fatalf(fmt string, v ...any) {
	log.Lock()
	defer log.Unlock()
	log.logger.Fatalf(fmt, v...)
}

func Fatalf(fmt string, v ...any) {
	logApi.Lock()
	defer logApi.Unlock()
	logApi.logger.Fatalf(fmt, v...)
}

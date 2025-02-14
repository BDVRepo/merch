package smart_context

import (
	"context"
	"time"

	gorm_logger "gorm.io/gorm/logger"
)

var _ gorm_logger.Interface = (*GormLoggerWrapper)(nil)

type GormLoggerWrapper struct {
	logger ISmartContext
}

func NewGormLoggerWrapper(logger ISmartContext) *GormLoggerWrapper {
	return &GormLoggerWrapper{logger: logger}
}

func (l *GormLoggerWrapper) LogMode(level gorm_logger.LogLevel) gorm_logger.Interface {
	return l
}

func (l *GormLoggerWrapper) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Infof(msg, data...)
}

func (l *GormLoggerWrapper) Debug(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Debugf(msg, data...)
}

func (l *GormLoggerWrapper) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Warnf(msg, data...)
}

func (l *GormLoggerWrapper) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Errorf(msg, data...)
}

func (l *GormLoggerWrapper) Trace(ctx context.Context, beginTime time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	elapsed := time.Since(beginTime)

	if err != nil {
		l.logger.Errorf("gorm query error (rows=%d, time=%s): sql=[%s]. Error: %v",
			rows, elapsed, sql, err)
	} else {
		l.logger.Debugf("gorm query success (rows=%d, time=%s): sql=[%s]",
			rows, elapsed, sql)
	}
}

// Package redislog provides logdriver emits log messages into redis
package redislog

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/daemon/logger"
)

// Name is the name of the file that the redis logs to.
const Name = "redis"

// redislog is Logger implementation for default Docker logging.
type redislog struct {
}

func init() {
	if err := logger.RegisterLogDriver(Name, New); err != nil {
		logrus.Fatal(err)
	}
	if err := logger.RegisterLogOptValidator(Name, ValidateLogOpt); err != nil {
		logrus.Fatal(err)
	}
}

// New creates new redislog which writes to filename passed in
// on given context.
func New(ctx logger.Context) (logger.Logger, error) {
	// TODO
	return &redislog{}, nil
}

// ValidateLogOpt looks for log options `server` (e.g. redis.smec.sap.corp:6379)
func ValidateLogOpt(cfg map[string]string) error {
	for key := range cfg {
		switch key {
		case "server":
		case "port":
		case "database":
		default:
			return fmt.Errorf("unknown log opt '%s' for redis log driver", key)
		}
	}
	return initAndValidate(cfg["server"], cfg["port"], cfg["database"])
}

// Log converts logger.Message to jsonlog.JSONLog and serializes it to file.
func (l *redislog) Log(msg *logger.Message) error {
	return nil
}

// Close closes underlying file and signals all readers to stop.
func (l *redislog) Close() error {
	return nil
}

// Name returns name of this logger.
func (l *redislog) Name() string {
	return Name
}

func initAndValidate(server, port, database string) error {
	logrus.Infof("Initialize redis connection to %s:%s.%s", server, port, database)
	// TODO:
	return nil
}

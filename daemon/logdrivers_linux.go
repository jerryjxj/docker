package daemon

import (
	_ "github.com/docker/docker/daemon/logger/jsonfilelog"
	_ "github.com/docker/docker/daemon/logger/loghub"
	_ "github.com/docker/docker/daemon/logger/redislog"
)

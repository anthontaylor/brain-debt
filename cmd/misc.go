package main

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type logAdapter struct{ log.Logger }

func (a logAdapter) Error(msg string) {
	level.Error(a.Logger).Log("component", "Jaeger", "msg", msg)
}

func (a logAdapter) Infof(msg string, args ...interface{}) {
	level.Info(a.Logger).Log("component", "Jaeger", "msg", fmt.Sprintf(msg, args...))
}

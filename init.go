package goxgo

import (
	"fmt"
	zmq "github.com/alecthomas/gozmq"
)

var (
	Context *zmq.Context
)

func init() {
	var err error
	Context, err = zmq.NewContext()
	if err != nil {
		panic(fmt.Sprintf("Could not acquire ZMQ context: %+v", err.Error()))
	}
}

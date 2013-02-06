// Client to make calls on external - in this first case - python services
// providing NLP/scientific math functions

package goxgo

import (
	"fmt"
	zmq "github.com/alecthomas/gozmq"
)

/*
Structure containing the connection endpoint

For now the only supported protocol option is "tcp" which will be a ZMQ
connection.

TODO: http protocol support. will need some refactoring
*/
type DSN struct {
	Protocol string
	Host	 string
	Port	 int
}

/*
Conn - ZMQ connection structure
*/
type Conn struct {
	Dsn		  *DSN
	Context   zmq.Context
	Socket    zmq.Socket
	connected bool
}

/*
Set up the connection to a goxgo service specified by the DSN
*/
func (c *Conn) Connect(dsn *DSN) (err error) {
	c.Context, err = zmq.NewContext()
	if err != nil {
		panic(fmt.Sprintf("Could not acquire ZMQ context: %+v", err.Error()))
	}
	c.Socket, err = c.Context.NewSocket(zmq.REQ)
	if err != nil {
		panic(fmt.Sprintf("Could not acquire ZMQ socket: %+v", err.Error()))
	}
	// using a global context and then making a lot of calls from separate
	// goroutines produces funky behaviour - running out of fds and
	// null pointer dereferences
	// c.Context = ZmqContext

	// TODO: add a conn/conf parameter to set a timeout
	// c.Socket.SetSockOptInt(zmq.LINGER, 0)
	if err == nil {
		c.connected = true
	}
	// fmt.Println("Connecting to: " + fmt.Sprintf("%v://%v:%v", dsn.Protocol, dsn.Host, dsn.Port) )
	c.Socket.Connect(fmt.Sprintf("%v://%v:%v", dsn.Protocol, dsn.Host, dsn.Port))
	return
}


/*
Close the connections zmq socket and zmq context
*/
func (c *Conn) Close() {
	if c.connected {
		c.Socket.Close()
		// fmt.Println("socket closed...")
		c.Context.Close()
		// fmt.Println("context closed...")
	}
	c.connected = false
	return
}

/*
Serialize the given payload, send it over the wire and return the
response data
*/
func (c *Conn) Send(payload interface{}) (response []byte) {
	msg, err := Serialize(&payload)
	if err != nil {
		panic(fmt.Sprintf("Could not serialize payload: %+v\n%+v", payload, err.Error()))
	}
	c.Socket.Send([]byte(msg), 0)
	response, _ = c.Socket.Recv(0)
	// fmt.Println(string(response))
	return
}

/*
Convinience function that will create a connection,
send a payload and Unserialize the reponse into a response
structure.
*/
func Call( dsn *DSN, request interface{}, response interface{}) {
	c := Conn{Dsn: dsn}
	var err error
	err = c.Connect(dsn)
	if err != nil {
		fmt.Println("Shit hit the fan: %v.", err)
	}
	defer c.Close()

	Unserialize(c.Send(&request), &response)
	return
}

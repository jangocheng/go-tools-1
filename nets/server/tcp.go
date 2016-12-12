// The simple TCP/UDP server.
package server

import (
	"errors"
	"net"

	"github.com/xgfone/go-tools/nets"
)

type THandle interface {
	Handle(conn *net.TCPConn)
}

// Wrap the function handler to the interface THandle.
type THandleFunc (func(*net.TCPConn))

func (h THandleFunc) Handle(conn *net.TCPConn) {
	h(conn)
}

// Wrap a panic, only print it, but ignore it.
func TCPWrapError(conn *net.TCPConn, handler THandle) {
	defer func() {
		if err := recover(); err != nil {
			_logger.Error("Get a error: %v", err)
		}

		if conn != nil {
			conn.Close()
		}
	}()

	handler.Handle(conn)
}

// Start a TCP server and never return. Return an error if returns.
//
// addr is like "host:port", such as "127.0.0.1:8000", and host or port may be omitted.
// size is the number of the pool. If it's 0, it's infinite.
// handle is the handler to handle the connection came from the client.
// handle is either a function whose type is func(*net.TCPConn), or a struct
// which implements the interface, THandle. Of course, you may wrap it by THandleFunc.
func TCPServerForever(addr string, handle interface{}) error {
	var handler THandle
	if _handler, ok := handle.(THandle); ok {
		handler = _handler
	} else if _handler, ok := handle.(func(*net.TCPConn)); ok {
		handler = THandleFunc(_handler)
	} else {
		panic("Don't support the handler")
	}

	var ln *net.TCPListener
	if _addr, err := net.ResolveTCPAddr("tcp", addr); err != nil {
		return err
	} else {
		if ln, err = net.ListenTCP("tcp", _addr); err != nil {
			return err
		}
	}

	defer ln.Close()

	_logger.Info("Listening on %v", addr)

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			_logger.Error("Failed to AcceptTCP: %v", err)
		} else {
			_logger.Debug("Get a connection from %v", conn.RemoteAddr())
			go TCPWrapError(conn, handler)
		}
	}

	// Never execute forever.
	return nil
}

// DialTCP is the same as DialTCPWithAddr, but it joins host and port firstly.
func DialTCP(host, port interface{}) (*net.TCPConn, error) {
	addr := nets.JoinHostPort(host, port)
	return DialTCPWithAddr(addr)
}

// DialTCPWithAddr dials a tcp connection to addr.
func DialTCPWithAddr(addr string) (*net.TCPConn, error) {
	var err error
	_conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	conn, ok := _conn.(*net.TCPConn)
	if !ok {
		_conn.Close()
		return nil, errors.New("not the tcp connection")
	}

	return conn, nil
}
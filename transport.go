package main

import "net"

type Peer interface {
	net.Conn
	Send([]byte) error
	CloseStream(int)
}

type Transport interface {
	Addr() string
	Dial(string, string, DialOptions) error
	ListenAndAccept(string, DialOptions) error
	Consume() <-chan msg
	Close() error
}

// type CommHandler interface {
// 	ClientComm
// 	ServerComm
// }

package tcptest
import (
	"net"
	"strings"
)
func C2S(addr string, server func(c net.Conn), client func(c net.Conn)) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	go func() {
		laddr := l.Addr().String()
		netConn, err := net.Dial("tcp", laddr)
		if err != nil {
			panic(err)
		}
		client(netConn)
		netConn.Close()
		l.Close()
	}()

	for {
		netConn, err := l.Accept()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				panic(err)
			}
			return
		}
		server(netConn)
	}
}

func C2L(addr string, server func(l net.Listener), client func(c net.Conn)) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	go func() {
		laddr := l.Addr().String()
		netConn, err := net.Dial("tcp", laddr)
		if err != nil {
			panic(err)
		}
		client(netConn)
		netConn.Close()
		l.Close()
	}()

	server(l)
}

func D2S(addr string, server func(c net.Conn), client func(laddr net.Addr)) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	go func() {
		client(l.Addr())
		l.Close()
	}()

	for {
		netConn, err := l.Accept()
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				panic(err)
			}
			return
		}
		server(netConn)
	}
}

func D2L(addr string, server func(l net.Listener), client func(c net.Addr)) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	go func() {
		client(l.Addr())
		l.Close()
	}()

	server(l)
}

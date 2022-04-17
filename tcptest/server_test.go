package tcptest

import (
	"io"
	"net"
	"testing"
	"time"

	"github.com/456vv/vconn"
)

// 由于服务端连接缓冲区【有】数据等待读取，关闭连接，服务端向客户端发送RST。客户端得到：wsarecv: An existing connection was forcibly closed by the remote host.
func Test_C2S_1(t *testing.T) {
	var err error
	C2S("127.0.0.1:0", func(c net.Conn) {
		c.Close()
		time.Sleep(1e9)
	}, func(c net.Conn) {
		defer c.Close()
		if _, err = c.Write([]byte("123")); err != nil {
			return
		}
		p := make([]byte, 2)
		_, err = c.Read(p)
		if op, ok := err.(*net.OpError); ok {
			switch op.Err.Error() {
			case "wsarecv: An established connection was aborted by the software in your host machine.":
				err = nil
			case "wsarecv: An existing connection was forcibly closed by the remote host.":
				err = nil
			}
			return
		}
	})
	if err != nil {
		t.Fatal(err)
	}
}

// 如果服务端连接缓冲区【没有】数据等待读取，关闭连接，服务端向客户端发送FIN。关客户端得到：EOF
func Test_C2S_2(t *testing.T) {
	var err error
	C2S("127.0.0.1:0", func(c net.Conn) {
		c.SetReadDeadline(time.Now().Add(1e9)) // 设置读取超时，防止io.Copy阻塞
		io.Copy(io.Discard, c)                 // 读完缓冲区数据
		c.Close()
		time.Sleep(1e9)
	}, func(c net.Conn) {
		defer c.Close()

		if _, err = c.Write([]byte("123")); err != nil {
			return
		}
		p := make([]byte, 2)
		if _, err = c.Read(p); err == io.EOF {
			// EOF
			err = nil
			return
		}
	})

	if err != nil {
		t.Fatal(err)
	}
}

// 服务器连接已经关闭情况下，给服务器发送数据，返回RST
func Test_C2S_3(t *testing.T) {
	var err error
	C2S("127.0.0.1:0", func(c net.Conn) {
		c.Close()
	}, func(c net.Conn) {
		defer c.Close()

		// 服务器连接已经关闭，客户端依然发送
		c.Write([]byte("123"))

		// 客户端收到wsarecv错误
		p := make([]byte, 2)
		_, err = c.Read(p)
		if op, ok := err.(*net.OpError); ok {
			switch op.Err.Error() {
			case "wsarecv: An established connection was aborted by the software in your host machine.":
				err = nil
			case "wsarecv: An existing connection was forcibly closed by the remote host.":
				err = nil
			}
			return
		}
	})

	if err != nil {
		t.Fatal(err)
	}
}

// 使用vonn库来检查
// 服务器连接已经关闭情况下，给服务器发送数据，返回RST
func Test_C2S_4(t *testing.T) {
	var err error
	C2S("127.0.0.1:0", func(c net.Conn) {
		c.Close()
	}, func(c net.Conn) {
		defer c.Close()

		// 服务器连接已经关闭，客户端依然发送
		// 稍候服务回复RST指令
		c.Write([]byte("123"))

		// 客户端收到wsarecv错误
		conn := vconn.New(c)
		select {
		case err = <-conn.CloseNotify():
			if op, ok := err.(*net.OpError); ok {
				switch op.Err.Error() {
				case "wsarecv: An established connection was aborted by the software in your host machine.":
					err = nil
				case "wsarecv: An existing connection was forcibly closed by the remote host.":
					err = nil
				}
				return
			}
		case <-time.After(1e6):
		}
	})

	if err != nil {
		t.Fatal(err)
	}
}

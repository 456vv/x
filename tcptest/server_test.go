package tcptest

import (
	"io"
	"net"
	"testing"
	"time"
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

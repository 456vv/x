package smtp
import (
	"github.com/456vv/vconn"
	"encoding/json"
	"os"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"errors"
	"crypto/tls"
	"sync"
//	"time"
)

//其中使得SMTP工作的基本的命令有7个，分别为：HELO﹑MAIL﹑RCPT﹑DATA﹑REST﹑NOOP和QUIT．下面分别介绍如下。
//
//HELO--发件方问候收件方，后面是发件人的服务器地址或标识。收件方回答OK时标识自己的身份。问候和确认过程表明两台机器可以进行通信，同时状态参量被复位，缓冲区被清空。
//
//MAIL--这个命令用来开始传送邮件，它的后面跟随发件方邮件地址（返回邮件地址）。它也用来当邮件无法送达时，发送失败通知。为保证邮件的成功发送，发件方的地址应是被对方或中间转发方同意接受的。这个命令会清空有关的缓冲区，为新的邮件做准备。
//
//RCPT --这个命令告诉收件方收件人的邮箱。当有多个收件人时，需要多次使用该命令，每次只能指明一个人。如果接收方服务器不同意转发这个地址的邮件，它必须报 550错误代码通知发件方。如果服务器同意转发，它要更改邮件发送路径，把最开始的目的地（该服务器）换成下一个服务器。
//
//DATA--收件方把该命令之后的数据作为发送的数据。数据被加入数据缓冲区中，以单独一行是"."的行结束数据。结束行对于接收方同时意味立即开始缓冲区内的数据传送，传送结束后清空缓冲区。如果传送接受，接收方回复OK。
//
//REST--这个命令用来通知收件方复位，所有已存入缓冲区的收件人数据，发件人数据和待传送的数据都必须清除，接收放必须回答OK.
//
//NOOP--这个命令不影响任何参数，只是要求接收放回答OK, 不会影响缓冲区的数据。
//
//QUIT--SMTP要求接收放必须回答OK，然后中断传输；在收到这个命令并回答OK前，收件方不得中断连接，即使传输出现错误。发件方在发出这个命令并收到OK答复前，也不得中断连接。

var errSmtpInfo 		= errors.New("配置信息未加载")
var errSmtpInfoServer	= errors.New("配置文件 Server 字段不正确！格式：host:port")
var errSmtpClosed		= errors.New("smtp已经被关闭")
type plainAuth struct {
	identity, username, password string
	host                         string
}
func PlainAuth(identity, username, password, host string) smtp.Auth {
	return &plainAuth{identity, username, password, host}
}
func (a *plainAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	if server.Name != a.host {
		return "", nil, errors.New("wrong host name")
	}
	resp := []byte(a.identity + "\x00" + a.username + "\x00" + a.password)
	return "PLAIN", resp, nil
}
func (a *plainAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		return nil, errors.New("unexpected server challenge")
	}
	return nil, nil
}
type Info struct{
	Server		string
	Host		string
	UserName	string
	Password	string
	FromEmail	string
	TLS			*tls.Config
	Extra		map[string]interface{}
}
type Smtp struct{
	Info 		*Info
	conn		net.Conn
	c			*smtp.Client
	connClosed	atomicBool
	closed		atomicBool
	inited		atomicBool
	mu			sync.Mutex
}

//打开配置文件
func (T *Smtp) OpenConfig(file string) error {
	osFile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer osFile.Close()
	return json.NewDecoder(osFile).Decode(&T.Info)
}
func (T *Smtp) connection() error {
	if T.closed.isTrue() {
		return errSmtpClosed
	}
	
	if T.connClosed.isTrue() || !T.inited.isTrue() {
		return T.reconnection()
	}
	return nil
}
func (T *Smtp) reconnection() error {
	if T.Info == nil {
		return errSmtpInfo
	}
	var(
		conn		net.Conn
		err			error
		host		= T.Info.Host
  		tlsconfig 	= T.Info.TLS
  	)
	if host == "" {
		host, _, err =net.SplitHostPort(T.Info.Server)
		if err != nil {
			return errSmtpInfoServer
		}
	}
	if tlsconfig == nil {
		conn, err = net.Dial("tcp", T.Info.Server)
	}else{
		conn, err = tls.DialWithDialer(new(net.Dialer), "tcp", T.Info.Server, tlsconfig)
	}
	if err != nil {
		return err
	}
	T.conn = vconn.NewConn(conn)
	
	c, err := smtp.NewClient(T.conn , host)
	if err != nil {
		return err
	}
	
	if tlsconfig == nil {
		if err := c.Hello("localhost"); err != nil {
			return err
		}
		
		//STARTTLS AUTH LOGIN PLAIN AUTH
		if ok, _ := c.Extension("STARTTLS"); ok {
	  		tlsconfig = &tls.Config{ServerName: host}
	  		if err = c.StartTLS(tlsconfig); err != nil {
	  			return err
	  		}
	  	}
	}
	
  	//map[PIPELINING: SIZE:73400320 AUTH:LOGIN PLAIN AUTH=LOGIN: MAILCOMPRESS: 8BITMIME:]
  	if ok, _ := c.Extension("AUTH"); ok {
	  	auth := PlainAuth("", T.Info.UserName, T.Info.Password, host)
	  	if err := c.Auth(auth); err != nil {
			return err
	  	}
  	}
  	
	T.c = c
	T.inited.setTrue()
	T.connClosed.setFalse()
	go T.notifyConn()
	return err
}
func (T *Smtp) Close() error {
	T.mu.Lock()
	defer T.mu.Unlock()
	
	T.closed.setTrue()
	if T.c != nil {
		err := T.c.Close()
		T.c 		= nil
		T.conn		= nil
		return err
	}
	return nil
}
func (T *Smtp) notifyConn(){
	if cn, ok := T.conn.(vconn.CloseNotifier); ok {
		select{
		case <-cn.CloseNotify():
			T.connClosed.setTrue()
		}
	}
}
func (T *Smtp) Send(to []string, title, body string) error {
	T.mu.Lock()
	defer T.mu.Unlock()
	
	if err := T.connection(); err != nil {
		return err
	}
	if err := T.c.Reset(); err != nil {
		return err
	}
	if err := T.c.Mail(T.Info.FromEmail); err != nil {
		return err
	}
	for _, addr := range to {
  		if err := T.c.Rcpt(addr); err != nil {
  			return err
  		}
  	}
  	w, err := T.c.Data()
  	if err != nil {
  		return err
  	}
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", T.Info.FromEmail, strings.Join(to, ";"), title, body)
  	_, err = w.Write([]byte(msg))
  	if err != nil {
  		return err
  	}
	err = w.Close()
  	if err != nil {
  		return err
  	}
  	return T.c.Quit()
}

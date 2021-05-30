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
	client		*smtp.Client
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
	T.mu.Lock()
	defer T.mu.Unlock()
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
  	
	T.client = c
	T.inited.setTrue()
	T.connClosed.setFalse()
	go T.notifyConn()
	return err
}
func (T *Smtp) Close() error {
	T.mu.Lock()
	defer T.mu.Unlock()
	
	T.closed.setTrue()
	if T.client != nil {
		err := T.client.Close()
		T.client 	= nil
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
	if err := T.connection(); err != nil {
		return err
	}
	if err := T.client.Mail(T.Info.FromEmail); err != nil {
		return err
	}
	for _, addr := range to {
  		if err := T.client.Rcpt(addr); err != nil {
  			return err
  		}
  	}
  	w, err := T.client.Data()
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
  	return nil
}
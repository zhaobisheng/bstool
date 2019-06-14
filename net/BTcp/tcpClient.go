package BTcp

import (
	"errors"
	"fmt"
	"net"
	"nodeServer/Utils"
	"sync"
	"time"
)

var t *TcpClient

type TcpClient struct {
	Timeout time.Duration
	Conn    net.Conn
	Send    chan []byte
	Read    chan []byte
	Status  bool
	Lock    sync.RWMutex
}

// 初始化一个Client
func New() *TcpClient {
	t = &TcpClient{
		Timeout: 10 * time.Second,
		Send:    make(chan []byte),
		Read:    make(chan []byte),
		Status:  false,
	}
	return t
}

// 设置超时时间，注意这个时间是每一次数据接收响应的超时时间，而不是总共的超时时间
func (t *TcpClient) SetTimeout(timeout time.Duration) *TcpClient {
	t.Timeout = timeout
	return t
}

func (t *TcpClient) Connect(protocol, addr string) error {
	var err error
	t.Conn, err = net.DialTimeout(protocol, addr, t.Timeout)
	if err != nil {
		fmt.Println("Connect-error:", err)
		return err
	} else {
		fmt.Println(t.Conn.LocalAddr(), "-", t.Conn.RemoteAddr())
		t.Status = true
		//go t.ReadLoop()
		//go t.WriteLoop()
		return nil
	}

}

func (t *TcpClient) Disconnect() error {
	t.Status = false
	return t.Conn.Close()
}

func (t *TcpClient) ReadLoop() {
	for {
		var data = make([]byte, 512)
		bufLen, err := t.Conn.Read(data[:])
		if bufLen > 0 {
			if err == nil {
				t.Lock.RLock()
				t.Read <- data[:bufLen]
				t.Lock.RUnlock()
			} else {
				fmt.Println(t.Conn.RemoteAddr(), " read error:", err)
			}
		}
	}
}

func (t *TcpClient) DoSend(msg []byte) {
	//fmt.Println("DoSend:", string(msg))
	_, err := t.Conn.Write(msg)
	if err != nil {
		fmt.Println(t.Conn.RemoteAddr(), " write error:", err)
	}
}

func (t *TcpClient) Start(cryptFlag bool, ReadCallback func([]byte)) {
	go t.ReadLoop()
	for {
		select {
		case message, _ := <-t.Send:
			if cryptFlag {
				message = Utils.Encrypt(message)
			}
			t.DoSend(message)
		case message, _ := <-t.Read:
			t.Lock.Lock()
			if cryptFlag {
				message = Utils.Decrypt(message)
			}
			t.Lock.Unlock()
			ReadCallback(message)
		}
	}
}

func (t *TcpClient) SendMsg(content string) error {
	if t.Status {
		t.Lock.Lock()
		t.Send <- []byte(content)
		t.Lock.Unlock()
		return nil
	} else {
		return errors.New("连接处于关闭状态!")
	}
}

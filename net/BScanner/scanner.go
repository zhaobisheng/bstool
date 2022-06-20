package BScanner

import (
	"bstool/net/Bipv4"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type Scanner struct {
	Timeout time.Duration
}

func Ip2long(ip string) uint32 {
	return Bipv4.Ip2long(ip)
}

func Long2ip(proper uint32) string {
	return Bipv4.Long2ip(proper)
}

// 初始化一个扫描器
func New() *Scanner {
	return &Scanner{
		6 * time.Second,
	}
}

// 设置超时时间，注意这个时间是每一次扫描的超时时间，而不是总共的超时时间
func (s *Scanner) SetTimeout(t time.Duration) *Scanner {
	s.Timeout = t
	return s
}

// 异步TCP扫描网段及端口，如果扫描的端口是打开的，那么将链接给定给回调函数进行调用
// 注意startIp和endIp需要是同一个网段，否则会报错,并且回调函数不会执行
func (s *Scanner) ScanIp(startIp string, endIp string, port int, callback func(net.Conn)) error {
	if callback == nil {
		return errors.New("callback function should not be nil")
	}
	var waitGroup sync.WaitGroup
	startIplong := Bipv4.Ip2long(startIp)
	endIplong := Bipv4.Ip2long(endIp)
	result := endIplong - startIplong
	if startIplong == 0 || endIplong == 0 {
		return errors.New("invalid startip or endip: ipv4 string should be given")
	}
	if result < 0 || result > 255 {
		return errors.New("invalid startip and endip: startip and endip should be in the same ip segment")
	}

	for i := startIplong; i <= endIplong; i++ {
		waitGroup.Add(1)
		go func(ip string) {
			//fmt.Println("scanning:", ip)
			// 这里必需设置超时时间
			conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), s.Timeout)
			if err == nil {
				callback(conn)
				conn.Close()
			}
			//fmt.Println("scanning:", ip, "done")
			waitGroup.Done()
		}(Bipv4.Long2ip(i))
	}
	waitGroup.Wait()
	return nil
}

// 注意startIp和endIp需要是同一个网段，否则会报错,并且回调函数不会执行
func (s *Scanner) ScanIpPortList(portList, startIp string, mask uint, callback func(net.Conn)) error {
	if callback == nil {
		return errors.New("callback function should not be nil")
	}
	startIplong := Bipv4.Ip2long(startIp)
	if startIplong == 0 {
		return errors.New("invalid startip or endip: ipv4 string should be given")
	}
	offset := 1 << mask
	offset = offset - 1
	allFinish := true
	finish := make(chan bool, 15)
	finishMode := 0
	for mode := 15; mode > 1; mode-- {
		if offset%mode == 0 {
			finishMode = mode
			allFinish = false
			//fmt.Println("ScanPortSplit-mode:", mode)
			split := offset / mode
			for start := 0; start < mode; start++ {
				startSplit := startIplong + uint32(start*split)
				endSplit := startIplong + uint32((start+1)*split)
				go s.ScanPortSplit(portList, startSplit, endSplit, callback, finish)
			}
			break
		} else {
			if mode == 2 {
				endIplong := startIplong + uint32(offset) - 1
				for i := startIplong; i <= endIplong; i++ {
					target := Bipv4.Long2ip(i)
					//go s.ScanPort(target, 80, callback)
					/*go s.ScanPort(target, 8333, callback)
					go s.ScanPort(target, 30303, callback)
					s.ScanPort(target, 9876, callback)*/
					s.DoScanPort(target, portList, callback)
				}
			}
		}
	}
	if !allFinish {
		finishNum := 0
		for {
			select {
			case flag, ok := <-finish:
				if flag && ok {
					finishNum++
					if finishNum >= finishMode {
						return nil
					}
				}
			}
		}
	}
	return nil
}

func (s *Scanner) ScanPortSplit(portList string, startIplong, endIplong uint32, callback func(net.Conn), finish chan bool) {
	for i := startIplong; i <= endIplong; i++ {
		target := Bipv4.Long2ip(i)
		//go s.ScanPort(target, 80, callback)
		s.DoScanPort(target, portList, callback)
	}
	finish <- true
	//fmt.Println(Bipv4.Long2ip(startIplong), " is ok!")
}

func (s *Scanner) DoScanPort(target, portList string, callback func(net.Conn)) bool {
	portNum := strings.Count(portList, ",")
	ports := strings.Split(portList, ",")
	for index := 0; index < portNum; index++ {
		go s.ScanPort(target, ports[index], callback)
		//go s.ScanPort(target, 8333, callback)
		//go s.ScanPort(target, 30303, callback)

	}
	s.ScanPort(target, ports[portNum], callback)

	return true
}

func (s *Scanner) ScanPort(ip, port string, callback func(net.Conn)) error {
	addr := fmt.Sprintf("%s:%s", ip, port)
	//fmt.Println("addr:", addr)
	if callback == nil {
		return errors.New("callback function should not be nil")
	}
	conn, err := net.DialTimeout("tcp", addr, s.Timeout)
	if err == nil {
		fmt.Println("found ", conn.RemoteAddr().String(), " is open!")
		callback(conn)
		conn.Close()
		return nil
	} else {
		return err
	}
}

// 扫描目标主机打开的端口列表
func (s *Scanner) ScanAllPort(ip string, callback func(net.Conn)) error {
	if callback == nil {
		return errors.New("callback function should not be nil")
	}
	finishChan := make(chan bool, 1024)
	var finishNum = 0
	//var waitGroup sync.WaitGroup
	for i := 0; i < 65536; i += 64 {
		go s.DoScanPortRange(ip, i, i+64, callback, finishChan)
		//waitGroup.Add(1)
		//fmt.Println("scanning:", i)
		/*go func(port int) {
		port := i
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), s.Timeout)
		waitGroup.Done()
		if err == nil {
			callback(conn)
			conn.Close()
			} else {
			return
			}
			}(i)
		}*/
	}
	//waitGroup.Wait()
	for {
		select {
		case <-finishChan:
			finishNum++
			if finishNum >= 1024 {
				break
			}
		}
	}
	return nil
}

func (s *Scanner) DoScanPortRange(target string, portStart, portEnd int, callback func(net.Conn), finish chan bool) bool {
	for i := portStart; i <= portEnd; i++ {
		s.ScanPort(target, fmt.Sprintf("%d", i), callback)
	}
	finish <- true
	return true
}

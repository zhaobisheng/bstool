package BScanner

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func Test_ScanPort(t *testing.T) {
	scan := &Scanner{
		Timeout: 5 * time.Second,
	}
	err := scan.ScanAllPort("222.186.175.120", func(conn net.Conn) { fmt.Println(conn.RemoteAddr()) })
	if err == nil {
		t.Log("第1个测试通过了")
	} else {
		t.Error("ScanPort-error:", err)
	}

}

/*func Test_ScanIp(t *testing.T) {
	scan := &Scanner{
		timeout: 5 * time.Second,
	}
	err := scan.ScanIp("192.168.31.1", "192.168.31.255", 3389, func(conn net.Conn) { fmt.Println(conn.RemoteAddr()) })
	if err == nil {
		t.Log("第2个测试通过了")
	} else {
		t.Error("Test_ScanIp-error:", err)
	}
}*/

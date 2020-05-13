package BScanner

import (
	"fmt"
	"net"
	"os"

	//"path/filepath"
	"testing"
	"time"
)

/*func Test_ScanPort(t *testing.T) {
	scan := &Scanner{
		Timeout: 5 * time.Second,
	}
	err := scan.ScanAllPort("222.186.175.120", func(conn net.Conn) { fmt.Println(conn.RemoteAddr()) })
	if err == nil {
		t.Log("第1个测试通过了")
	} else {
		t.Error("ScanPort-error:", err)
	}

}*/

/*func Test_ScanIp(t *testing.T) {
	scan := &Scanner{
		Timeout: 5 * time.Second,
	}
	err := scan.ScanIp("192.168.0.1", "192.168.0.255", 80, func(conn net.Conn) { fmt.Println(conn.RemoteAddr()) })
	if err == nil {
		t.Log("第2个测试通过了")
	} else {
		t.Error("Test_ScanIp-error:", err)
	}
}*/
func Test_ScanIpList(t *testing.T) {
	scan := &Scanner{
		Timeout: 3 * time.Second,
	}
	err := scan.ScanIpPortList("80,8080,7002", "172.31.100.1", 8, func(conn net.Conn) { GenerateReport(conn.RemoteAddr().String()) })
	if err == nil {
		t.Log("第2个测试通过了")
	} else {
		t.Error("Test_ScanIp-error:", err)
	}
}

func GenerateReport(content string) bool {
	content += "\r\n"
	now := time.Now()
	lastWeek := time.Unix(now.Unix()-(7*86400), 0)
	Today := fmt.Sprintf("%d月%d号-%d月%d号", int(lastWeek.Month()), lastWeek.Day(), int(now.Month()), now.Day())
	LOGDir := Today + "扫描结果"
	if !PathExists(LOGDir) {
		os.Mkdir(LOGDir, os.ModePerm)
	}
	//dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir := "/Users/fc/Downloads/poc"
	file, err := os.OpenFile(dir+"/"+LOGDir+"/"+Today+"扫描结果.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	if err != nil {

		fmt.Println("生成扫描结果文件失败:", err)
		return false
	} else {
		_, err = file.WriteString(content)
		if err != nil {
			fmt.Println("生成扫描结果文件失败:", err)
			return false
		}
	}
	fmt.Println("生成扫描结果文件成功!")
	return true
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

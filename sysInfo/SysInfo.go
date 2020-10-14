package BSysInfo

import (
	"math"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

func GetCPUPercent() int {
	cpuNumber, _ := cpu.Counts(true)
	cu, err := cpu.Percent(time.Second, true)
	if err == nil {
		var cpuUsage float64 = 0
		for _, usage := range cu {
			cpuUsage += usage
		}
		return int(math.Floor(cpuUsage/float64(cpuNumber) + 0.5))
	} else {
		return 0
	}
}

func GetMemoryPercent() int {
	vm, err := mem.VirtualMemory()
	if err == nil {
		return int(math.Floor(vm.UsedPercent + 0.5))
	} else {
		return 0
	}
}

func GetDiskPercent() int {
	di, err := disk.Usage("/")
	if err == nil {
		return int(math.Floor(di.UsedPercent + 0.5))
	} else {
		return 0
	}
}

func GetHostname() string {
	host, err := os.Hostname()
	if err != nil {
		return ""
	} else {
		return host
	}
}

func GetDNS() ([]string, error) {
	dns := make([]string, 2)
	cmd := exec.Command("cat", "/etc/resolv.conf")
	buf, err := cmd.Output()
	if err != nil {
		return dns, err
	}
	dnsContent := string(buf)
	var reg *regexp.Regexp
	reg = regexp.MustCompile(`(2(5[0-5]{1}|[0-4]\d{1})|[0-1]?\d{1,2})(\.(2(5[0-5]{1}|[0-4]\d{1})|[0-1]?\d{1,2})){3}`)
	rs := reg.FindAllString(dnsContent, -1)
	if len(rs) == 1 {
		dns = make([]string, 1)
		dns[0] = rs[0]
	} else if len(rs) == 0 {
		dns = make([]string, 0)
	} else {
		dns[0] = rs[0]
		dns[1] = rs[1]
	}
	return dns, nil
}

func SetDNS(newDns []string) error {
	dns, err := GetDNS()
	if err != nil {
		return err
	}
	filename := "/etc/resolv.conf"
	readFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer readFile.Close()
	var buf []byte = make([]byte, 512)
	bufLen, err := readFile.Read(buf)
	if err != nil {
		return err
	}
	fileContent := string(buf[:bufLen])
	tagFlag := "nameserver"
	tagFlagLen := len(tagFlag)
	index := strings.Index(fileContent, tagFlag)
	if len(dns) == 1 {
		fileContent = strings.Replace(fileContent, dns[0], newDns[0], -1)
		if len(newDns) > 1 {
			if index > 0 {
				tempStr := fileContent[index+tagFlagLen:]
				index1 := strings.Index(tempStr, "\n")
				if index1 > 0 {
					tempStr1 := fileContent[index+tagFlagLen+index1+1:]
					fileContent = fileContent[:index+tagFlagLen+index1+1] + tagFlag + " " + newDns[1] + "\n" + tempStr1
				} else {
					fileContent = fileContent + "\n" + tagFlag + " " + newDns[1]
				}
			}
		}
	} else if len(dns) == 0 {
		fileContent = fileContent + "\n" + tagFlag + " " + newDns[0]
		fileContent = fileContent + "\n" + tagFlag + " " + newDns[1]
	} else {
		if index > 0 {
			tempStr := fileContent[index+tagFlagLen:]
			index1 := strings.Index(tempStr, tagFlag)
			if index1 > 0 {
				tempStr1 := tempStr[index1:]
				headStr := fileContent[:index]
				tempStr = fileContent[index : index+tagFlagLen+index1]
				tempStr = strings.Replace(tempStr, dns[0], newDns[0], -1)
				if len(newDns) > 1 {
					tempStr1 = strings.Replace(tempStr1, dns[1], newDns[1], -1)
				} else {
					index2 := strings.Index(tempStr1, "\n")
					if index2 > 0 {
						tempStr1 = tempStr1[index2+1:]
					} else {
						tempStr1 = ""
					}
				}
				fileContent = headStr + tempStr + tempStr1
			}
		}
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(fileContent)
	return nil
}

func GetGateway() (string, error) {
	cmd := exec.Command("netstat", "-r")
	buf, err := cmd.Output()
	if err != nil {
		return "", err
	}
	gatewayContent := string(buf)
	tagFlag := "default"
	tagFlagLen := len(tagFlag)
	index := strings.Index(gatewayContent, tagFlag)
	if index > 0 {
		gatewayContent = gatewayContent[index+tagFlagLen : index+tagFlagLen+25]
		gatewayContent = strings.Replace(gatewayContent, " ", "", -1)
		return gatewayContent, nil
	}
	return "", nil
}

func GetTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

/*func SetTime(server string) error {
	//param := "ntp1.aliyun.com"
	command := "/home/updateTime.sh"
	cmd := exec.Command(command, server)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}*/

func UpdateTime(server string) error {
	cmd := exec.Command("ntpdate", "-s \""+server+"\" ")
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	cmd1 := exec.Command("hwclock", "--systohc")
	cmd1.Output()
	return nil
}

func SetHostname(hostname string) error {
	cmd := exec.Command("hostname", hostname)
	_, err := cmd.Output()
	if err != nil {
		//fmt.Println("exec: SetHostname error")
		return err
	}
	file, err := os.OpenFile("/etc/hostname", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		//fmt.Println("/etc/hostname: Open error!")
		return err
	}
	defer file.Close()
	_, err = file.WriteString(hostname + "\n")
	if err != nil {
		//fmt.Println("/etc/hostname: Write error!")
		return err
	}
	return nil
}

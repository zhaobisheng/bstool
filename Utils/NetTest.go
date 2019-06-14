package Utils

import (
	"net"
	"time"
)

func CheckNetwork(ip, port string) (bool,error) {
	conn, err := net.DialTimeout("tcp", ip+":"+port, time.Second)
	if err != nil {
		return false,err
	}
	defer conn.Close()
	return true,nil
}

func CheckServerAlive(addr string) (bool,error) {
	conn, err := net.DialTimeout("tcp", addr, time.Second*3)
	if err != nil {
		return false,err
	}
	defer conn.Close()
	return true,nil
}

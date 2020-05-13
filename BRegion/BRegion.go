package BRegion

import (
	"log"
	"sync"
)

type MyRegion struct {
	RegionDB *Ip2Region
	Lock     sync.RWMutex
	closed   bool
}

var myRegion *MyRegion

func InitIPdb() (err error) {
	myRegion = &MyRegion{}
	myRegion.RegionDB, err = New("ip2region.db")
	defer myRegion.RegionDB.Close()
	if err != nil {
		log.Fatalln("func-InitIPdb-error:", err)
		return err
	}
	return nil
}

func IsChinaIP(ip string) bool {
	return myRegion.CheckChinaIP(ip)
}

func CloseIPdb(ip string) {
	myRegion.RegionDB.Close()
}

func (region *MyRegion) CheckChinaIP(ip string) bool {
	ipInfo, err := region.RegionDB.MemorySearch(ip)
	if err != nil {
		return true
	}
	if ipInfo.Country != "中国" {
		return false
	}
	return true
}

func (region *MyRegion) GetIPDetails(ip string) (IpInfo, error) {
	ipInfo, err := region.RegionDB.MemorySearch(ip)
	if err != nil {
		return ipInfo, err
	}
	return ipInfo, nil
}

func GetCityInfo(ip string) (IpInfo, error) {
	//return region.RegionDB.MemorySearch(ip)
	return myRegion.GetIPDetails(ip)
}

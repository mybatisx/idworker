package util


import (
	"fmt"
	"net"
	"os"
)

type  util struct{
	name string
	age  int
}

func  GetIp() string  {

	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var ip  string
	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip=ipnet.IP.String()
			}

		}
	}
	return ip
}

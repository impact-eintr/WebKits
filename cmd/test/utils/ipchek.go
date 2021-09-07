package utils

import (
	"strconv"
	"strings"
)

func Headcheck(ip string) []string  {

	var ipresult []string
	isexists := strings.Contains(ip, "-")
	if isexists {
		ips := strings.Split(ip, "-")
		ip1 := strings.Split(ips[0], ".")
		ip2, _ := strconv.Atoi(ip1[3])
		ip3, _ := strconv.Atoi(ips[1])
		for i := ip2; i <= ip3; i++ {

			ipresult = append(ipresult, ip1[0]+"."+ip1[1]+"."+ip1[2]+"."+strconv.Itoa(i))
		}
	} else {
		ipresult = append(ipresult,ip)
		return ipresult
	}
	return ipresult
}

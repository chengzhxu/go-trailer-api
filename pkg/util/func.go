package util

import (
	"go-trailer-api/pkg/logging"
	"net"
	"net/http"
	"strings"
	"time"
)

func ExistIntElement(element int64, array []int64) bool {
	for _, e := range array {
		if e == element {
			return true
		}
	}

	return false
}

func GetNowTimeStamp() int {
	return int(time.Now().Unix())
}

//千分位
func NumFormat(str string) string {
	numStr := strings.Split(str, ".")[0] //如果有小数获取整数部分
	length := len(numStr)
	if length < 4 {
		return str
	}
	count := (length - 1) / 3 //取于-有多少组三位数
	for i := 0; i < count; i++ {
		numStr = numStr[:length-(i+1)*3] + "," + numStr[length-(i+1)*3:]
	}
	return numStr
}

//当前时间 Y-M-d
func GetCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

//当前时间 Y-M-d H:i:s
func GetCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//获取服务端IP
func GetServiceIP() string {
	address, err := net.InterfaceAddrs()

	if err != nil {
		logging.Error(err)
	}

	ip := ""
	for _, address := range address {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}

		}
	}

	return ip
}

// ClientIP 尽最大努力实现获取客户端 IP 的算法。
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

// ClientPublicIP 尽最大努力实现获取客户端公网 IP 的算法。
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func ClientPublicIP(r *http.Request) string {
	var ip string
	for _, ip = range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		ip = strings.TrimSpace(ip)
		if ip != "" && !HasLocalIPddr(ip) {
			return ip
		}
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" && !HasLocalIPddr(ip) {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		if !HasLocalIPddr(ip) {
			return ip
		}
	}

	return ""
}

// HasLocalIPddr 检测 IP 地址字符串是否是内网地址
func HasLocalIPddr(ip string) bool {
	return HasLocalIP(net.ParseIP(ip))
}

// HasLocalIP 检测 IP 地址是否是内网地址
func HasLocalIP(ip net.IP) bool {
	//for _, network := range localNetworks {
	//	if network.Contains(ip) {
	//		return true
	//	}
	//}

	return ip.IsLoopback()
}

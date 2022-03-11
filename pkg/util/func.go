package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/ip2region"
	"github.com/nacos-group/nacos-sdk-go/vo"
	ip2region2 "go-trailer-api/pkg/ip2region"
	"go-trailer-api/pkg/logging"
	"go-trailer-api/pkg/nacos"
	"go-trailer-api/pkg/setting"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

var signatureSalt = "trailer_signature_salt" //客户端接口签名 加密盐
var standbyConf setting.StandbyTimeConf      //客户端待机时长配置

//APP 包下载地址
type appPackage struct {
	AppName     string `json:"app_name"`
	AppEnName   string `json:"app_en_name"`
	FileMd5     string `json:"file_md5"`
	PackageName string `json:"package_name"`
	PackageUrl  string `json:"package_url"`
}
type appPackageArray []*appPackage

func GetStandbyTime() (error, int) {
	standbyConf = setting.StandbyTimeSetting

	return nil, standbyConf.Duration
}

//app 包下载地址
func GetAppPackage() appPackageArray {
	conf := appPackageArray{}

	content, err := nacos.NacosClient.GetConfig(vo.ConfigParam{
		DataId: setting.NacosAppPackageConfSetting.DataId,
		Group:  setting.NacosAppPackageConfSetting.Group,
	})
	if err != nil {
		logging.Error(err)
	} else {
		json.Unmarshal([]byte(content), &conf)
	}

	return conf
}

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

// 日期转换为时间戳
func TimeToUnix(t string) int {
	timeLayout := "2006-01-02 15:04:05"          //转化所需模板
	loc, _ := time.LoadLocation("Asia/Shanghai") //设置时
	tt, _ := time.ParseInLocation(timeLayout, t, loc)

	return int(tt.Unix())
}

//检查时间类型
func CheckTime(t string) string {
	dArr := strings.Split(t, " ")
	if len(dArr) == 2 {
		time := dArr[1]
		tArr := strings.Split(time, ":")
		var ntArr []string
		for _, v := range tArr {
			t, _ := strconv.Atoi(v)
			if t > 59 {
				ntArr = append(ntArr, "59")
			} else if t == 0 {
				ntArr = append(ntArr, "00")
			} else {
				ntArr = append(ntArr, strconv.Itoa(t))
			}
		}

		nt := strings.Join(ntArr, ":")
		t = dArr[0] + " " + nt
	}

	return t
}

// 日期转换为时间戳
func DateToUnix(t string) int {
	timeLayout := "2006-01-02"                   //转化所需模板
	loc, _ := time.LoadLocation("Asia/Shanghai") //设置时区
	tt, _ := time.ParseInLocation(timeLayout, t, loc)

	return int(tt.Unix())
}

// 时间戳转换为日期
func UnixToTime(ut int64) string {
	timeLayout := "2006-01-02 15:04:05" //转化所需模板
	tm := time.Unix(ut, 0)

	return tm.Format(timeLayout)
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

// 判断数组中是否存在某字符串
func StrInArray(target string, strArray []string) bool {
	sort.Strings(strArray)
	index := sort.SearchStrings(strArray, target)
	//index的取值：[0,len(str_array)]
	if index < len(strArray) && strArray[index] == target { //需要注意此处的判断，先判断 &&左侧的条件，如果不满足则结束此处判断，不会再进行右侧的判断
		return true
	}
	return false
}

// 截取字符串
func SubStr(str string, start int, length int) string {
	if str == "" || start < 0 || length <= 0 {
		return ""
	}
	rs := []rune(str)
	return string(rs[start:length])
}

func Md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

// 参数签名解析   方式：参数 按字母顺序排序后，用 & 拼接，然后 md5
func CheckParamSignature(params map[string]interface{}, signature string) bool {
	var dataParams string = "?" // 拼接后的参数

	//拼接
	for k, v := range params {
		if k == "signature" {
			signature = fmt.Sprintf("%v", v)
		} else {
			val := fmt.Sprintf("%v", v)
			if len(val) > 0 {
				dataParams += k + "=" + fmt.Sprintf("%v", v) + "&"
			}
		}
	}
	us, _ := url.Parse(dataParams)
	dataParams = us.Query().Encode()
	dataParams += "&" + signatureSalt

	mySignature := Md5V(dataParams)

	if mySignature == signature {
		return true
	}

	return false
}

func GetIpInfoByRequest(r *http.Request) ip2region.IpInfo {
	ip := ClientPublicIP(r)
	if ip == "" {
		ip = ClientIP(r)
	}
	//根据 IP 获取地域信息
	ips, _ := ip2region2.IpRegionDb.MemorySearch(ip)

	return ips
}

//生成随机数
func RandNumber(min, max float64, decimal int) string {
	gap := max - min
	rd := rand.Int63n(int64(gap))
	rd += int64(GetNowTimeStamp())
	ret := float64(rd) / 1.0e6

	if decimal == 0 {
		return strconv.FormatFloat(ret, 'f', decimal, 64)
	}

	return fmt.Sprintf("%."+strconv.Itoa(decimal)+"f\n", ret)
}

func RandIntNumber(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min) + min
}

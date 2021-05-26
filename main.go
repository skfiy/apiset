package main

import (
	destiny "apiset/destinyData"
	tm "apiset/voice"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const GO_URL = "http://localhost:8081"

type MsgInfo struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
	Time string      `json:"time"`
}

func MsgData(code int, msg string, data interface{}) string {
	msgs := &MsgInfo{code, msg, data, time.Now().Format("2006-01-02 15:04:05")}
	Info, _ := json.Marshal(msgs)
	return string(Info)
}

/**
 * 获取命运2周报
 */
func GetDestinyDataSet(w http.ResponseWriter, r *http.Request) {

	var data interface{}

	pipe := make(chan interface{})

	go func() {
		data = destiny.DataInfo()
		pipe <- data
	}()

	Info := <-pipe

	fmt.Fprintf(w, "%s", Info)

}

func TextConversionMp3(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Fatal("系统错误" + err.Error())
	}
	var msg string
	content := r.Form.Get("content")
	s, _ := ioutil.ReadAll(r.Body)
	fmt.Println(s)
	if content == "" {
		msg = MsgData(200, "参数值为空", "")
		fmt.Fprintf(w, "%s", msg)
	} else {
		msg = tm.TextMp3(content)
		fmt.Fprintf(w, "%s", msg)
	}
	fmt.Println(content)

}

func Fanyi(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Fatal("系统错误" + err.Error())
	}
	var msg string
	content := r.Form.Get("data")
	s, _ := ioutil.ReadAll(r.Body)
	fmt.Println(s)
	if content == "" {
		msg = MsgData(200, "参数值为空", "")
		fmt.Fprintf(w, "%s", msg)
	} else {
		msg = Fanyiget(content)
		fmt.Fprintf(w, "%s", msg)
	}
	fmt.Println(content)
}

func Mp3Init(w http.ResponseWriter, r *http.Request) {
	realPath := r.URL.String()
	fmt.Println(realPath)
	file, err := os.Open("." + realPath)
	defer file.Close()
	if err != nil {
		log.Println("static resource:", err)
		w.WriteHeader(404)
	} else {
		bs, _ := ioutil.ReadAll(file)

		w.Write(bs)
	}
}

//首页
func Index(w http.ResponseWriter, r *http.Request) {
	q := "如有问题请联系QQ:1923021"
	s := "------API------"
	z := "\r  地址 api.cloolc.club 以下接口\n" +
		"\r  /text?content=       //文字转语音获取MP3\n" +
		"\r  /destiny             //命运2周报\n" +
		"\r  /getip               //获取IP地址位置\n" +
		"\r  /lunar               //获取农历日\n" +
		"\r  /fanyi?data=         //翻译\n" +
		"\r ------API------"

	fmt.Fprintf(w, "  %s \n %s\n%s\n", q, s, z)

}

func GetIp(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json
	//获取IP地址
	ip := ClientPublicIP(r)
	if ip == "" {
		ip = ClientIP(r)
	}
	json := Ipcity(ip)

	fmt.Fprintf(w, "%s", json)

}

//查询 IP 位置
func Ipcity(ip string) string {
	url := "http://ip.taobao.com/outGetIpInfo?ip=" + ip + "&accessKey=alibaba-inc"
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return string(data)
}

//查询农历日
func Lunar(w http.ResponseWriter, r *http.Request) {
	url := "http://www.skfiy.com/lunar/index.php"
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, "%s", data)
}

//翻译
func Fanyiget(d string) string {
	url := "http://www.skfiy.com/fy/index.php?s=" + d
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return string(data)
}

//获取当前请求用户IP地址
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
		fmt.Println(ip)
		if ip != "" {
			return ip
		}
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

func main() {
	http.HandleFunc("/fanyi", Fanyi)
	http.HandleFunc("/lunar", Lunar)
	http.HandleFunc("/getip", GetIp)
	http.HandleFunc("/mp3/", Mp3Init)
	http.HandleFunc("/text", TextConversionMp3)
	http.HandleFunc("/destiny", GetDestinyDataSet)
	http.HandleFunc("/", Index)

	go func() {
		for {
			time.Sleep(time.Second)
			log.Println("Checking if started...")
			resp, err := http.Get(GO_URL)
			if err != nil {
				log.Println("Failed:", err)
				continue
			}
			resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				log.Println("Not OK:", resp.StatusCode)
				continue
			}
			break
		}
		log.Println("SERVER 启动成功!")
		log.Println("URL：", GO_URL)
	}()
	err := http.ListenAndServe(":8081", nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

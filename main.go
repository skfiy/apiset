package main

import (
	destiny "apiset/destinyData"
	tm "apiset/voice"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const GO_URL = "http://localhost:8081";

type MsgInfo struct {
	 Code int   `json:"code"`
	 Msg  string   `json:"msg"`
	 Data interface{}  `json:"data"`
	 Time string `json:"time"`
}


func  MsgData(code int,msg string,data interface{}) string  {
	  msgs:=  &MsgInfo{code,msg,data,time.Now().Format("2006-01-02 15:04:05")}
	  Info,_ := json.Marshal(msgs)
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

func TextConversionMp3(w http.ResponseWriter, r *http.Request){

	err := r.ParseForm()
	if err != nil {
		log.Fatal("系统错误" + err.Error())
	}
	var  msg string
	content := r.Form.Get("content")
		s,_:= ioutil.ReadAll(r.Body)
		fmt.Println(s)
	if content == ""  {
		msg = MsgData(200,"参数值为空","");
		fmt.Fprintf(w, "%s", msg);
	}else{
		msg = tm.TextMp3(content)
		fmt.Fprintf(w, "%s", msg);
	}
	fmt.Println(content);

}

func Mp3Init(w http.ResponseWriter, r *http.Request){
	realPath := r.URL.String()
	fmt.Println(realPath);
	file,err :=os.Open("."+realPath)
	defer file.Close()
	if err != nil {
		log.Println("static resource:", err)
		w.WriteHeader(404)
	} else {
		bs,_ := ioutil.ReadAll(file)

		w.Write(bs)
	}
}

//首页
func Index(w http.ResponseWriter, r *http.Request){
	q := "QQ:1923021";
	s := "------API------";
	z := "\r  /text?content=      //文字转语音获取MP3\n" +
		 "\r  /destiny            //命运2周报\n" +
		 "\r ------API------";
	fmt.Fprintf(w, "  %s \n %s\n%s", q,s,z);

}

func main() {
	http.HandleFunc("/mp3/",Mp3Init)
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
		log.Println("URL：",GO_URL)
	}()
	err := http.ListenAndServe(":8081", nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

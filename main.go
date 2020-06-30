package main

import (
	destiny "apiset/destinyData"
	"fmt"
	"log"
	"net/http"
	"time"
)
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

func main() {

	http.HandleFunc("/", GetDestinyDataSet)
	go func() {
		for {
			time.Sleep(time.Second)
			log.Println("Checking if started...")
			resp, err := http.Get("http://localhost:8081")
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
	}()
	err := http.ListenAndServe(":8081", nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

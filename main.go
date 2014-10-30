package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/index.html") //解析模板文件
	t.Execute(w, nil)
}

func ChatWith(ws *websocket.Conn) {
	var err error

	buffer, err := ioutil.ReadFile("test.xls")
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	arr := strings.Split(string(buffer), "\n")

	for {
		var page string

		if err = websocket.Message.Receive(ws, &page); err != nil {
			fmt.Println("Can't receive")
			break
		}

		pageCount, _ := strconv.Atoi(page)
		pageSize := 1000
		pageCount = pageCount * pageSize
		i := 0
		data := []string{}
		for {
			if i < pageSize {
				i++
			} else {
				break
			}
			pageCount++
			if pageCount >= len(arr) {
				os.Exit(0)
			}
			data = append(data, arr[pageCount])
		}
		encode, _ := json.Marshal(data)
		if err = websocket.Message.Send(ws, string(encode)); err != nil {
			fmt.Println("Can't send")
			break
		}
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.Handle("/ws", websocket.Handler(ChatWith))
	err := http.ListenAndServe(":8080", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

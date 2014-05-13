package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"strconv"
	"strings"
)

var (
	base_url = "http://data.live.126.net/www/"//45743/live.json"
	room = "45743"
)

func init() {
	flag.StringVar(&room, "room", "", "Room Id")
}

type JsonObject struct {
	Logs []Log
}

type DelId struct {
	id int
}
type Log struct {
	Now          string
	Section      string
	MatchTime    string
	Icon         string
	ImgSrc       string
	ImgHref      string
	ImgDesType   string
	Msg          string
	MsgeHref     string
	MsgFontType  string
	MsgFontColor string
	Score        string
	Commentator  string
	Id           string
}

var prev_id int

func Connect(room string, c chan string) {
	resp, err := http.Get(base_url + room + "/live.json")
	if err != nil {
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var jsonobject JsonObject
	json_body := body[2 : len(body)-2]
	// fmt.Printf("%s\n", json_body)
	err = json.Unmarshal(json_body, &jsonobject)
	if err != nil {
		fmt.Println("unmarshal error:", err)
	}

	// fmt.Printf("Unmarshal: %+v\n", jsonobject)
	logs := jsonobject.Logs
	for i := 0; i < len(logs); i++ {
		id, _ := strconv.Atoi(logs[i].Id)
		if id > prev_id {
			r := strings.NewReplacer("<span>", "", "</span>", "")
			c <- logs[i].Now + " " + logs[i].Section + " " + logs[i].Msg + " \033[31;2m" + r.Replace(logs[i].Score) + "\033[0m"
			// fmt.Printf("Id: %s, Now:%s, Section:%s, Msg:%s, Score:%s\n",
			// 	logs[i].Id, logs[i].Now, logs[i].Section, logs[i].Msg, logs[i].Score)
			prev_id, _ = strconv.Atoi(logs[i].Id)
		}

	}
}

func testSend(c chan string) {
	for {
		c <- "Google China Demlution." + time.Now().String()
		time.Sleep(5 * time.Second)
	}
}

func testConnect(room string, c chan string) {
	for {
		go Connect(room, c)
		time.Sleep(5 * time.Second)
	}

}
func main() {
	flag.Parse()
	c := make(chan string)

	// go testSend(c)
	go testConnect(room, c)

	fmt.Println("Start receive")
	a := "abc"
	b := "abc"
	fmt.Println(a==b)
	for i := range c {
		
		fmt.Println(i)
		fmt.Println(<-c)
	}

}

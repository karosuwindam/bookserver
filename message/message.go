package message

import (
	"encoding/json"
	"time"
)

const (
	LOGOUTPUT_ON  logout = 1
	LOGOUTPUT_OFF logout = 0
)

type logout int

type Message struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Code   int    `json:"code"`
}

type Result struct {
	Name   string      `json:"name"`
	Code   int         `json:"code"`
	Option string      `json:"option"`
	Date   time.Time   `json:"date"`
	Result interface{} `json:"result"`
}

// MessageのJSON変換
func (m *Message) Output() string {
	msg := Message{Name: m.Name, Status: m.Status, Code: m.Code}
	bytes, _ := json.Marshal(msg)
	return string(bytes)

}

//メッセージの設定
func (m *Message) InputMessage(msg string, flag logout) {
	m.Status = msg
	if flag == LOGOUTPUT_ON {
		Println(m.Status)
	}
}

//ResultのJSON変換
func (r *Result) Output() string {
	result := Result{Name: r.Name, Date: r.Date, Result: r.Result, Option: r.Option, Code: r.Code}
	bytes, _ := json.Marshal(result)
	return string(bytes)
}

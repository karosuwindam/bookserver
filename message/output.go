package message

import (
	"fmt"
	"time"
)

const (
	formtdata = "2006-01-02 15:04:05.999"
)

func Println(v ...interface{}) (int, error) {
	now := time.Now()
	output := []interface{}{
		"[" + now.Format(formtdata) + "]",
	}
	for _, vs := range v {
		output = append(output, vs)
	}

	return fmt.Println(output...)
}

func Printf(v ...interface{}) (int, error) {
	now := time.Now()
	output := "[" + now.Format(formtdata) + "] " + v[0].(string)
	return fmt.Printf(output, v[1:]...)

}

func Print(v ...interface{}) (int, error) {
	now := time.Now()
	output := []interface{}{
		"[" + now.Format(formtdata) + "] ",
	}
	for _, vs := range v {
		output = append(output, vs)
	}
	return fmt.Print(output...)
}

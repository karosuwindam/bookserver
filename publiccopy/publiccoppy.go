package publiccopy

import (
	"bookserver/config"
	"bookserver/publiccopy/copycopyfile"
	"bookserver/publiccopy/tablecopyfile"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type CopyFils struct {
	Id   int  `json:"Id"`
	Flag bool `json:"Flag"`
}

var input chan interface{}

// run(data)
//
// dataに入力された変数の型別に処理を実施する。
func run(data interface{}) {
	switch data.(type) {
	case CopyFils:
		if tmp, ok := data.(CopyFils); ok {
			if name, err := tablecopyfile.GetTableByName(tmp.Id, ""); err != nil {
				log.Println(err)
			} else {
				if tmp.Flag {
					if err := copycopyfile.CopyFile(name); err != nil {
						log.Println(err)
					}
				} else {
					if err := copycopyfile.RemoveFile(name); err != nil {
						log.Println(err)
					}
				}
				if err := tablecopyfile.AddCopyFIle(name, tmp.Flag); err != nil {
					log.Println(err)
				}
			}
		}
	default:
		log.Println("input data", data)
	}
}

// loop処理
func Loop(ctx context.Context) {
	fmt.Println("Public Copy Loop Start")
	var wp sync.WaitGroup
loop:
	for {
		select {
		case data := <-input:
			wp.Add(1)
			run(data)
			wp.Done()
		case <-ctx.Done():
			wp.Wait()
			break loop
		case <-time.After(time.Minute * 30): //定期処理
			tablecopyfile.CkCopyFileAll()
		}
	}
	fmt.Println("Public Copy Loop End")
}

// 処理を依頼する処理
func Add(data interface{}) error {
	ctx, cancel := context.WithCancel(context.Background())
	go func(data interface{}) {
		input <- data
		cancel()
	}(data)
	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		cancel()
		return errors.New("Add time out")
	}
	return nil
}

// 終了時の処理
func Wait() error {
	go func() {
		for {
			if len(input) == 0 {
				endch <- nil
			}
			time.Sleep(time.Microsecond * 100)
		}
	}()
	select {
	case err := <-endch:
		return err
	case <-time.After(1 * time.Second):
		close(input)
		return errors.New("Public Copy Loop Time Out")
	}
}

var endch chan error

// セットアップ
func Setup(cfg *config.Config) error {
	endch = make(chan error, 1)
	input = make(chan interface{}, 20)
	if err := copycopyfile.Setup(cfg); err != nil {
		return err
	}
	if err := tablecopyfile.Setup(cfg); err != nil {
		return err
	}
	return nil
}

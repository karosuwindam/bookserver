package readzipfile

import (
	"bookserver/config"
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var zipPath string //アップロードされたzipファイルの保存先
var shutdown chan bool
var shutdown_back chan bool
var loopflag bool

// 初期化処理
func Init() error {
	//zipフォルダの設定
	zipPath = config.BScfg.Zip
	if zipPath[len(zipPath):] != "/" {
		zipPath += "/"
	}
	if err := os.MkdirAll(zipPath, 0777); err != nil {
		return err
	}
	if err := cashInit(); err != nil { //キャッシュ機能の初期化
		return err
	}
	//ループ処理
	shutdown = make(chan bool, 1)
	shutdown_back = make(chan bool, 1)
	return nil
}

// 定期処理
func Run(ctx context.Context) {
	var wg sync.WaitGroup
	slog.InfoContext(ctx, "zip cash loop start")
	loopflag = true
loop:
	for {
		select {
		case <-shutdown:
			shutdown_back <- true
			break loop
		case filename := <-chname: //キャッシュ作成処理
			wg.Add(1)
			go func(filename string) {
				defer wg.Done()
				ctxzipfile, _ := context.WithTimeout(ctx, 1*time.Hour)
				if err := readZipFileAll(filename, ctxzipfile); err != nil {
					slog.ErrorContext(ctxzipfile,
						fmt.Sprintf("readZipFileAll(%v) error", filename),
						"name", filename,
						"Error", err,
					)
				}
			}(filename)
		case <-time.After(1 * time.Second):
			ctxzipfileCash, _ := context.WithTimeout(ctx, 1*time.Hour)
			//１秒ごとの処理
			if err := clearZipFileCash(ctxzipfileCash); err != nil { //キャッシュ定期削除処理
				slog.ErrorContext(ctxzipfileCash,
					"clearZipFileCash error",
					"Error", err,
				)
			}
		}
	}
	wg.Wait()
	loopflag = false
	slog.InfoContext(ctx, "zip cash loop end")

}

// シャットダウン処理
func Stop() error {
	ctx := context.TODO()
	slog.DebugContext(ctx, "zip cash loop stop start")
	if loopflag {
		shutdown <- true
		select {
		case <-shutdown_back:
			break
		case <-time.After(1 * time.Second):
			return errors.New("Stop error timeout")
		}
	}
	slog.DebugContext(ctx, "zip cash loop stop end")
	return nil
}

// Zipファイルの中身をキャッシュへ登録する処理
func AddCash(filename string) error {
	ctx := context.TODO()
	select {
	case chname <- filename:
		slog.DebugContext(ctx,
			fmt.Sprintf("AddCash(%v) ok", filename),
		)
		break
	case <-time.After(1 * time.Second):
		return errors.New("Time Out")
	}

	return nil
}

// Zipファイルの中身を読み取る
func ReadZipfile(filename, zipName string) (*bytes.Buffer, error) {
	ctx := context.TODO()
	buf, err := ReadCashData(filename, zipName)
	if err != nil {
		slog.DebugContext(ctx,
			fmt.Sprintf("ReadZipfile(%v, %v) not cash", filename, zipName),
			"filename", filename,
			"zipName", zipName,
			"Error", err,
		)
		var wg sync.WaitGroup
		wg.Add(1)
		go func(filename string) {
			defer wg.Done()
			if err := AddCash(filename); err != nil {
				slog.ErrorContext(ctx,
					fmt.Sprintf("ReadZipfile AddCash(%v) error", filename),
					"filename", filename,
					"Error", err,
				)
			}
		}(filename)
		buf, err = readZipFileData(filename, zipName)
		wg.Wait()
	}
	return buf, err
}

package listdata

import (
	"bookserver/dirread"
	"bookserver/health/healthmessage"
	"bookserver/table"
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"
)

type ListData struct {
	Name    string `json:"Name"`
	PdfName string `json:"Pdf"`
	PdfFlag bool   `json*"PdfFlag"`
	ZipName string `json:"Zip"`
	ZipFlag bool   `json:"ZipFlag"`
}

var listData []ListData
var listDataTmp []ListData

var listData_mux sync.Mutex

var message healthmessage.HealthMessage

// listDataにデータを登録する
func addListData(list ListData) {
	listData_mux.Lock()
	defer listData_mux.Unlock()
	listDataTmp = append(listDataTmp, list)
}

// 一時保存のリストデータを出力できるようにする
func reNewListData() {
	listData_mux.Lock()
	defer listData_mux.Unlock()
	listData = listDataTmp
}

// listDataの値を読み込む
func readListData() []ListData {
	listData_mux.Lock()
	defer listData_mux.Unlock()
	tmp := listData
	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].Name < tmp[j].Name
	})
	return tmp
}

// 動作チェック確認用
func Health() healthmessage.HealthMessage {
	return message
}

func Loop(ctx context.Context) {
	hMessage := healthmessage.Create(apiname)
	if listData == nil {
		return
	}
	fmt.Println("ListData Loop Start")
	hMessage.ChangeMessage("Create Loop Start", true)
	message = hMessage.ChangeOut()
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-time.After(time.Second):
			hMessage.ChangeMessage("Check FileList Table by folder")
			message = hMessage.ChangeOut()
			ckFileListTableData()
			hMessage.ChangeMessage("OK", true)
			message = hMessage.ChangeOut()
		}
	}
	hMessage.ChangeMessage("ListData Loop End", false)
	message = hMessage.ChangeOut()
	fmt.Println("ListData Loop End")
}

func Wait() error {
	return nil
}

// 登録したテーブルとフォルダ内のデータをチェックして差分がないことを確認する
func ckFileListTableData() {
	var tabledata []table.Filelists
	listDataTmp = []ListData{}
	var ziplist []string
	var pdflist []string
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		if jdata, err := sql.ReadAll(table.FILELIST); err != nil {
			log.Println(err)
			return
		} else {
			if listData, ok := table.JsonToStruct(table.FILELIST, []byte(jdata)).([]table.Filelists); ok {
				tabledata = listData
			}
		}
	}()
	go func() {
		defer wg.Done()
		if f, err := dirread.Setup(pdffolder); err != nil {
			log.Println(err)
			return
		} else {
			if err := f.Read(""); err == nil {
				tmp := []string{}
				for _, str := range f.Data {
					tmp = append(tmp, str.Name)
				}
				pdflist = tmp
			}
		}
	}()
	go func() {
		defer wg.Done()
		if f, err := dirread.Setup(zipfolder); err != nil {
			log.Println(err)
			return
		} else {
			if err := f.Read(""); err == nil {
				tmp := []string{}
				for _, str := range f.Data {
					tmp = append(tmp, str.Name)
				}
				ziplist = tmp
			}
		}
	}()
	wg.Wait()
	rmZip := []string{}
	var zipmu sync.Mutex
	rmPdf := []string{}
	var pdfmu sync.Mutex
	for _, tdata := range tabledata {
		wg.Add(1)
		go func(list ListData) {
			defer wg.Done()
			var wg1 sync.WaitGroup
			wg1.Add(2)
			go func(list *ListData) {
				defer wg1.Done()
				for _, s := range ziplist {
					if list.ZipName == s {
						list.ZipFlag = true
						zipmu.Lock()
						defer zipmu.Unlock()
						rmZip = append(rmZip, s)
						break
					}
				}
			}(&list)
			go func(list *ListData) {
				defer wg1.Done()
				for _, s := range pdflist {
					if list.PdfName == s {
						list.PdfFlag = true
						pdfmu.Lock()
						defer pdfmu.Unlock()
						rmPdf = append(rmPdf, s)
						break
					}
				}
			}(&list)
			wg1.Wait()
			addListData(list)
		}(ListData{
			Name:    tdata.Name,
			ZipName: tdata.Zippass,
			PdfName: tdata.Pdfpass,
		})
	}
	wg.Wait()
	ziplist = rmStringData(ziplist, rmZip)
	pdflist = rmStringData(pdflist, rmPdf)
	ziplist = addStringDataForTable(ziplist, pdflist)
	addZipStringForTable(ziplist)
	reNewListData()
}

// 2つの配列から未登録の配列を作る(pdfベース)
func addStringDataForTable(zip, pdf []string) []string {
	rmtmp := []string{}
	var rmTmp_mu sync.Mutex
	var wg sync.WaitGroup
	for _, str := range pdf {
		if strings.Index(strings.ToLower(str), ".pdf") > 0 {
			wg.Add(1)
			go func(str string) {
				tmp_str := str[:len(str)-4]
				defer wg.Done()
				flag := true
				for i := 0; i < len(zip); i++ {
					if len(tmp_str) >= len(zip[i]) {
						if tmp_str == zip[i][:len(tmp_str)] {
							rmTmp_mu.Lock()
							defer rmTmp_mu.Unlock()
							rmtmp = append(rmtmp, zip[i])
							addListData(ListData{
								PdfName: str,
								ZipName: zip[i],
							})
							flag = false
							break
						}
					}
				}
				if flag {
					addListData(ListData{
						Name:    tmp_str,
						PdfName: str,
						ZipName: tmp_str + ".zip",
					})
				}

			}(str)
		}
	}
	wg.Wait()
	return rmStringData(rmtmp, zip)
}

// zipリストからテーブルを作成する
func addZipStringForTable(zip []string) {
	var wg sync.WaitGroup
	for _, str := range zip {
		if strings.Index(strings.ToLower(str), ".zip") > 0 {
			wg.Add(1)
			go func(str string) {
				tmp_str := str[:len(str)-4]
				addListData(ListData{
					Name:    tmp_str,
					ZipName: str,
				})
			}(str)
		}
	}
	wg.Wait()
}

// ２つの配列から一致しない項目を抜き出す
func rmStringData(a, b []string) []string {
	var tmp1, tmp2 []string
	sort.Slice(b, func(i, j int) bool {
		return b[i] < b[j]
	})
	sort.Slice(a, func(i, j int) bool {
		return a[i] < a[j]
	})
	tmp1 = a
	tmp2 = b
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < len(a); i++ {
			for j := 0; j < len(tmp2); j++ {
				if a[i] == tmp2[j] {
					tmp := []string{}
					tmp = append(tmp, tmp2[:j]...)
					tmp = append(tmp, tmp2[j+1:]...)
					tmp2 = tmp
					break
				}
			}
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < len(b); i++ {
			for j := 0; j < len(tmp1); j++ {
				if a[i] == tmp1[j] {
					tmp := []string{}
					tmp = append(tmp, tmp1[:j]...)
					tmp = append(tmp, tmp1[j+1:]...)
					tmp1 = tmp
					break
				}
			}
		}
	}()
	wg.Wait()
	if len(tmp1) == len(tmp2) && len(tmp1) == 0 {
		return []string{}
	} else {
		output := []string{}
		output = append(output, tmp1...)
		output = append(output, tmp2...)
		return output
	}
}

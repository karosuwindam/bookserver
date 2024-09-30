package tablesearch

import (
	"bookserver/config"
	"bookserver/table/booknames"
	"bookserver/table/copyfiles"
	"bookserver/table/filelists"
	"bookserver/webserver/api/common"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
)

type SearchKey struct {
	Table   string `json:"Table"`
	Keyword string `json:"Keyword"`
}

var keynames []string = []string{
	"today",
	"toweek",
	"tomonth",
	"rand",
} //特殊キーワードリスト

// Getにより受けったURLをベースにデータ検索を実施してJSON形式を返す
func GetSearchTable(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "", "URL", r.URL, "Method", r.Method)
	table := r.PathValue("table")
	keyword := r.PathValue("keyword")
	if checkTableData(table) != nil || keyword == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	}
	jout := SearchKey{
		Table:   table,
		Keyword: keyword,
	}

	if tmp := jout.serchData(); tmp == "" {
		//検索したが失敗したときの処理
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
	} else {
		//検索して成功したときの処理
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(tmp))
	}

}

// Postにより受け取ったJSONデータをベースに検索を実施しての結果をJSON形式で返す
func PostSerchTable(w http.ResponseWriter, r *http.Request) {
	b, _ := io.ReadAll(r.Body)
	slog.InfoContext(r.Context(), "", "URL", r.URL, "Method", r.Method, "body", string(b))
	jout := SearchKey{}
	if err := json.Unmarshal(b, &jout); err != nil || jout.Table == "" {
		//入力データを異常やテーブル指定されていないときの処理
		slog.ErrorContext(r.Context(), "PostSerchTable", "error", err.Error())
	} else {
		if checkTableData(jout.Table) != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("page not found"))
			return
		}
		if tmp := jout.serchData(); tmp == "" {
			//検索したが失敗したときの処理
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("[]"))
		} else {
			//検索して成功したときの処理
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(tmp))
		}
	}
}

// 受け取ったJSONデータをもとに検索結果を返す
func (t *SearchKey) serchData() string {
	output := ""
	switch t.Keyword {
	case "today": //今日更新
		output = toDayGetData(t.Table)
	case "toweek": //今週更新
		output = toWeekGetData(t.Table)
	case "tomonth": //今月更新
		output = toMonthGetData(t.Table)
	case "rand": //ランダムなデータ取得
		output = randGetData(t.Table)
	default:
		output = serchGetData(t.Table, t.Keyword)
	}
	return output
}

// 今日更新したデータを取得
func toDayGetData(table string) string {
	var output string
	switch table {
	case BOOKNAMES:
		if ary, err := booknames.ReadByDataRangeDay(booknames.TODAY); err != nil {
			log.Println(err)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				log.Println(errj)
			} else {
				output = string(b)
			}
		}
	case FILELISTS:
		if ary, err := filelists.ReadByDataRangeDay(filelists.TODAY); err != nil {
			log.Println(err)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				log.Println(errj)
			} else {
				output = string(b)
			}
		}
	}
	return output
}

// 今週更新したデータを取得
func toWeekGetData(table string) string {
	var output string
	switch table {
	case BOOKNAMES:
		if ary, err := booknames.ReadByDataRangeDay(booknames.TOWEEK); err != nil {
			log.Println(err)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				log.Println(errj)
			} else {
				output = string(b)
			}
		}
	case FILELISTS:
		if ary, err := filelists.ReadByDataRangeDay(filelists.TOWEEK); err != nil {
			log.Println(err)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				log.Println(errj)
			} else {
				output = string(b)
			}
		}
	}
	return output
}

// 今月更新したデータを更新
func toMonthGetData(table string) string {
	var output string
	switch table {
	case BOOKNAMES:
		if ary, err := booknames.ReadByDataRangeDay(booknames.TOMONTH); err != nil {
			log.Println(err)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				log.Println(errj)
			} else {
				output = string(b)
			}
		}
	case FILELISTS:
		if ary, err := filelists.ReadByDataRangeDay(filelists.TOMONTH); err != nil {
			log.Println(err)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				log.Println(errj)
			} else {
				output = string(b)
			}
		}
	}
	return output
}

// ランダムなデータを取得
func randGetData(table string) string {
	var output string
	switch table {
	case BOOKNAMES:
		if ary, err := booknames.ReadRandData(config.BScfg.RandamRead); err != nil {
			log.Println(err)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				log.Println(errj)
			} else {
				output = string(b)
			}
		}
	case FILELISTS:
		if ary, err := filelists.ReadRandData(config.BScfg.RandamRead); err != nil {
			log.Println(err)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				log.Println(errj)
			} else {
				output = string(b)
			}
		}
	}
	return output
}

// テーブルを指定して検索を実施
func serchGetData(table string, keyword string) string {
	var output string
	switch table {
	case BOOKNAMES:
		if ary, err := booknames.Search(keyword); err != nil {
			log.Println(err)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				log.Println(errj)
			} else {
				output = string(b)
			}
		}
	case COPYFILES:
		if ary, err := copyfiles.Search(keyword); err != nil {
			log.Println(err)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				log.Println(errj)
			} else {
				output = string(b)
			}
		}
	case FILELISTS:
		if ary, err := filelists.Search(keyword); err != nil {
			log.Println(err)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				log.Println(errj)
			} else {
				output = string(b)
			}
		}
	}
	return output
}

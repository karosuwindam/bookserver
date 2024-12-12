package tablesearch

import (
	"bookserver/config"
	"bookserver/table/booknames"
	"bookserver/table/copyfiles"
	"bookserver/table/filelists"
	"bookserver/webserver/api/common"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	ctx := r.Context()
	slog.InfoContext(ctx,
		fmt.Sprintf("%v %v", r.Method, r.URL),
		"Url", r.URL,
		"Method", r.Method,
	)
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
	ctx := r.Context()
	slog.InfoContext(ctx,
		fmt.Sprintf("%v %v", r.Method, r.URL),
		"Url", r.URL,
		"Method", r.Method,
		"data", string(b),
	)
	jout := SearchKey{}
	if err := json.Unmarshal(b, &jout); err != nil || jout.Table == "" {
		//入力データを異常やテーブル指定されていないときの処理
		slog.ErrorContext(ctx, "PostSerchTable json.Unmarshal error",
			"error", err,
		)
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
	ctx := context.TODO()
	var output string
	switch table {
	case BOOKNAMES:
		if ary, err := booknames.ReadByDataRangeDay(booknames.TODAY); err != nil {
			slog.WarnContext(ctx, "toDayGetData booknames.ReadByDataRangeDay error",
				"error", err,
			)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				slog.WarnContext(ctx, "toDayGetData booknames json.Marshal error",
					"error", errj,
				)
			} else {
				output = string(b)
			}
		}
	case FILELISTS:
		if ary, err := filelists.ReadByDataRangeDay(filelists.TODAY); err != nil {
			slog.WarnContext(ctx, "toDayGetData filelists.ReadByDataRangeDay error",
				"error", err,
			)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				slog.WarnContext(ctx, "toDayGetData filelists json.Marshal error",
					"error", errj,
				)
			} else {
				output = string(b)
			}
		}
	}
	slog.DebugContext(ctx,
		fmt.Sprintf("toDayGetData read table:%v", table),
		"table", table,
		"output", output,
	)
	return output
}

// 今週更新したデータを取得
func toWeekGetData(table string) string {
	ctx := context.TODO()
	var output string
	switch table {
	case BOOKNAMES:
		if ary, err := booknames.ReadByDataRangeDay(booknames.TOWEEK); err != nil {
			slog.WarnContext(ctx, "toWeekGetData booknames.ReadByDataRangeDay error",
				"error", err,
			)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				slog.WarnContext(ctx, "toWeekGetData booknames json.Marshal error",
					"error", errj,
				)
			} else {
				output = string(b)
			}
		}
	case FILELISTS:
		if ary, err := filelists.ReadByDataRangeDay(filelists.TOWEEK); err != nil {
			slog.WarnContext(ctx, "toWeekGetData filelists.ReadByDataRangeDay error",
				"error", err,
			)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				slog.WarnContext(ctx, "toWeekGetData filelists json.Marshal error",
					"error", errj,
				)
			} else {
				output = string(b)
			}
		}
	}
	slog.DebugContext(ctx,
		fmt.Sprintf("toWeekGetData read table:%v", table),
		"table", table,
		"output", output,
	)
	return output
}

// 今月更新したデータを更新
func toMonthGetData(table string) string {
	ctx := context.TODO()
	var output string
	switch table {
	case BOOKNAMES:
		if ary, err := booknames.ReadByDataRangeDay(booknames.TOMONTH); err != nil {
			slog.WarnContext(ctx, "toMonthGetData booknames.ReadByDataRangeDay error",
				"error", err,
			)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				slog.WarnContext(ctx, "toMonthGetData booknames json.Marshal error",
					"error", errj,
				)
			} else {
				output = string(b)
			}
		}
	case FILELISTS:
		if ary, err := filelists.ReadByDataRangeDay(filelists.TOMONTH); err != nil {
			slog.WarnContext(ctx, "toMonthGetData filelists.ReadByDataRangeDay error",
				"error", err,
			)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				slog.WarnContext(ctx, "toMonthGetData filelists json.Marshal error",
					"error", errj,
				)
			} else {
				output = string(b)
			}
		}
	}
	slog.DebugContext(ctx,
		fmt.Sprintf("toMonthGetData read table:%v", table),
		"table", table,
		"output", output,
	)
	return output
}

// ランダムなデータを取得
func randGetData(table string) string {
	ctx := context.TODO()
	var output string
	switch table {
	case BOOKNAMES:
		if ary, err := booknames.ReadRandData(config.BScfg.RandamRead); err != nil {
			slog.WarnContext(ctx, "randGetData booknames.ReadRandData error",
				"error", err,
			)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				slog.WarnContext(ctx, "randGetData booknames json.Marshal error",
					"error", errj,
				)
			} else {
				output = string(b)
			}
		}
	case FILELISTS:
		if ary, err := filelists.ReadRandData(config.BScfg.RandamRead); err != nil {
			slog.WarnContext(ctx, "randGetData filelists.ReadRandData error",
				"error", err,
			)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				slog.WarnContext(ctx, "randGetData filelists json.Marshal error",
					"error", errj,
				)
			} else {
				output = string(b)
			}
		}
	}
	slog.DebugContext(ctx,
		fmt.Sprintf("randGetData read table:%v", table),
		"table", table,
		"output", output,
	)
	return output
}

// テーブルを指定して検索を実施
func serchGetData(table string, keyword string) string {
	ctx := context.TODO()
	var output string
	switch table {
	case BOOKNAMES:
		if ary, err := booknames.Search(keyword); err != nil {
			slog.WarnContext(ctx,
				fmt.Sprintf("serchGetData booknames.Search keyword:%v error", keyword),
				"keyword", keyword,
				"error", err,
			)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				slog.WarnContext(ctx,
					fmt.Sprintf("serchGetData booknames json.Marshal error"),
					"error", errj,
				)
			} else {
				output = string(b)
			}
		}
	case COPYFILES:
		if ary, err := copyfiles.Search(keyword); err != nil {
			slog.WarnContext(ctx,
				fmt.Sprintf("serchGetData copyfiles.Search keyword:%v error", keyword),
				"keyword", keyword,
				"error", err,
			)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				slog.WarnContext(ctx,
					fmt.Sprintf("serchGetData copyfiles json.Marshal error"),
					"error", errj,
				)
			} else {
				output = string(b)
			}
		}
	case FILELISTS:
		if ary, err := filelists.Search(keyword); err != nil {
			slog.WarnContext(ctx,
				fmt.Sprintf("serchGetData filelists.Search keyword:%v error", keyword),
				"keyword", keyword,
				"error", err,
			)
		} else {
			msg := common.Message(ary)
			if b, errj := json.Marshal(&msg); errj != nil {
				slog.WarnContext(ctx,
					fmt.Sprintf("serchGetData filelists json.Marshal error"),
					"error", errj,
				)
			} else {
				output = string(b)
			}
		}
	}
	slog.DebugContext(ctx,
		fmt.Sprintf("serchGetData read table:%v keyword:%v", table, keyword),
		"table", table,
		"keyword", keyword,
		"output", output,
	)
	return output
}

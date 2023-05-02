package health

import (
	"bookserver/api/common"
	"bookserver/api/listdata"
	"bookserver/api/view"
	"bookserver/config"
	"bookserver/health/healthmessage"
	"bookserver/transform"
	"bookserver/webserver"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var apiname = "health" //api名

// health(w, r)
//
// ヘルスチェックの結果表示
func health(w http.ResponseWriter, r *http.Request) {
	var msg common.Result = common.Result{Code: http.StatusOK, Date: time.Now(), Option: r.Method}

	output := []healthmessage.HealthMessage{}
	if tmp := listdata.Health(); &tmp != nil {
		output = append(output, tmp)
	}
	if tmp := view.Health(); &tmp != nil {
		output = append(output, tmp)
	}
	if tmp := transform.Health(); &tmp != nil {
		output = append(output, tmp)
	}
	msg.Name = apiname
	msg.Url = r.URL.Path
	for _, tmp := range output {
		if !tmp.Flag {
			msg.Code = http.StatusServiceUnavailable
			break
		}
	}
	if tmp, err := json.Marshal(output); err != nil {
		log.Println(err)
	} else {
		msg.Result = string(tmp)
	}

	common.Sqlreadmessageback(msg, w)
}

var route []webserver.WebConfig = []webserver.WebConfig{
	{Pass: "/" + apiname, Handler: health},
}

// Setup
func SetUp(cfg *config.Config) ([]webserver.WebConfig, error) {
	return route, nil
}

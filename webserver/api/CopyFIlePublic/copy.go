package copyfilepublic

import (
	"bookserver/controller/copyfile"
	"bookserver/webserver/api/common"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type CopyFileRecv struct {
	Id   int  `json:"Id"`
	Flag bool `json:"Flag"`
}

func PostCopyFile(w http.ResponseWriter, r *http.Request) {
	data, _ := io.ReadAll(r.Body)
	log.Println("info:", r.URL, r.Method, string(data))
	msg := common.Message("NG")
	//処理を後で各
	tmp := CopyFileRecv{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		log.Println("error:", err)
		w.Write(msg.Json())
		return
	}
	if err := copyfile.Add(tmp.Id, tmp.Flag); err != nil {
		log.Println("error:", err)
		w.Write(msg.Json())
		return
	}
	msg = common.Message("OK")
	w.Write(msg.Json())
}

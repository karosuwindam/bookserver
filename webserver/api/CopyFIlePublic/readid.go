package copyfilepublic

import (
	"bookserver/controller/copyfile"
	"bookserver/webserver/api/common"
	"log"
	"net/http"
	"strconv"
)

func GetFileDataFlagById(w http.ResponseWriter, r *http.Request) {
	log.Println("info:", r.URL, r.Method)
	tmpid := r.PathValue("id")
	id, err := strconv.Atoi(tmpid)
	if err != nil {
		log.Println("error:", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	}
	// idに紐づいたファイルが存在するか確認
	if d, err := copyfile.ReadCopyFIleFlagById(id); err != nil {
		log.Println("error:", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	} else {
		// 確認した結果JSON形式でを返す。
		msg := common.Message(d)
		w.Write(msg.Json())
	}

}

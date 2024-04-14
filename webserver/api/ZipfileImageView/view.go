package zipfileimageview

import (
	readzipfile "bookserver/controller/readZipfile"
	"bookserver/table/filelists"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func GetZipFileImageView(w http.ResponseWriter, r *http.Request) {
	log.Println("info:", r.URL, r.Method)

	filename := r.PathValue("filename")
	tmpid := r.PathValue("id")
	id, err := strconv.Atoi(tmpid)
	if err != nil {
		log.Println("error:", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	}
	if d, err := filelists.GetId(id); err != nil {
		log.Println("error:", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	} else {
		if buf, err := readzipfile.ReadZipfile(d.Zippass, filename); err != nil {
			log.Println("error:", err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("page not found"))
			return
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%s", buf)
		}
	}
}

package zipfileimageview

import (
	readzipfile "bookserver/controller/readZipfile"
	"bookserver/table/filelists"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

func GetZipFileImageView(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "", "URL", r.URL, "Method", r.Method)

	filename := r.PathValue("filename")
	tmpid := r.PathValue("id")
	id, err := strconv.Atoi(tmpid)
	if err != nil {
		slog.ErrorContext(r.Context(), "GetZipFileImageView", "error", err.Error())
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	}
	if d, err := filelists.GetId(id); err != nil {
		slog.ErrorContext(r.Context(), "GetZipFileImageView", "error", err.Error())
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	} else {
		if buf, err := readzipfile.ReadZipfile(d.Zippass, filename); err != nil {
			slog.ErrorContext(r.Context(), "GetZipFileImageView", "error", err.Error())
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("page not found"))
			return
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%s", buf)
		}
	}
}

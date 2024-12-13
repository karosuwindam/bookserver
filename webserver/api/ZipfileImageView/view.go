package zipfileimageview

import (
	"bookserver/config"
	readzipfile "bookserver/controller/readZipfile"
	"bookserver/table/filelists"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

func GetZipFileImageView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := config.TracerS(ctx, "GetZipFileImageView", "Get Zip File Image View")
	defer span.End()

	slog.InfoContext(ctx,
		fmt.Sprintf("%v %v", r.Method, r.URL),
		"Url", r.URL,
		"Method", r.Method,
	)

	filename := r.PathValue("filename")
	tmpid := r.PathValue("id")
	id, err := strconv.Atoi(tmpid)
	if err != nil {
		slog.ErrorContext(ctx, "GetZipFileImageView Atoi Error",
			"tmpid", tmpid,
			"error", err,
		)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	}
	if d, err := filelists.GetId(id); err != nil {
		slog.ErrorContext(ctx, "GetZipFileImageView filelists.GetId error",
			"id", id,
			"error", err,
		)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	} else {
		if buf, err := readzipfile.ReadZipfile(d.Zippass, filename); err != nil {
			slog.ErrorContext(ctx, "GetZipFileImageView readzipfile.ReadZipfile error",
				"zippass", d.Zippass,
				"filename", filename,
				"error", err,
			)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("page not found"))
			return
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%s", buf)
		}
	}
}

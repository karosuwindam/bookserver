package filedownload

import (
	"bookserver/config"
	"bookserver/table/filelists"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func GetDownload(w http.ResponseWriter, r *http.Request) {
	ctx, span := config.TracerS(r.Context(), "GetDownload", "Get Download Main")
	defer span.End()
	slog.InfoContext(ctx,
		fmt.Sprintf("%v %v", r.Method, r.URL),
		"Url", r.URL,
		"Method", r.Method,
	)
	ctx, ok := readIdAndType(r)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	}

	ctx, err := readFileName(ctx)
	if err != nil {
		config.TracerError(span, err)
		slog.ErrorContext(ctx,
			"GetDownload readFileName error",
			"error", err,
		)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	}
	if b, err := readFileData(ctx); err != nil {
		config.TracerError(span, err)
		slog.ErrorContext(ctx,
			"GetDownload readFileData error",
			"error", err,
		)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
		return
	} else {
		v, _ := contextReadFileNameFIleType(ctx)
		// ファイル名
		w.Header().Set("Content-Disposition", "attachment; filename="+v.filename)
		// コンテントタイプ
		w.Header().Set("Content-Type", "application/"+v.filetype)
		// ファイルの長さ
		w.Header().Set("Content-Length", strconv.Itoa(len(b)))
		// bodyに書き込み
		w.Write(b)

	}
}

// readIdAndType(r) ctx ok
func readIdAndType(r *http.Request) (context.Context, bool) {
	ctx, span := config.TracerS(r.Context(), "readIdAndType", "read Id And Type")
	defer span.End()
	tmpid := r.PathValue("id")
	filetype := r.PathValue("filetype")
	slog.DebugContext(ctx, "context read data",
		"id", tmpid,
		"filetype", filetype,
	)
	id, err := strconv.Atoi(tmpid)
	if err != nil {
		config.TracerError(span, err)
		slog.ErrorContext(ctx,
			fmt.Sprintf("GetDownload strconv.Atoi id=%v", tmpid),
			"Id", tmpid,
			"Error", err,
		)
		return ctx, false
	}
	return contextWriteIdFIleType(ctx, id, filetype), true

}

func readFileName(ctx context.Context) (context.Context, error) {
	ctx, span := config.TracerS(ctx, "readFileName", "read File Name")
	defer span.End()

	v, ok := contextReadIdFileType(ctx)
	if !ok {
		return ctx, fmt.Errorf("context read err data")
	}
	id := v.id
	filetype := v.filetype

	span.AddEvent("readfile", trace.WithAttributes(
		attribute.Int("id", id),
		attribute.String("filetype", filetype),
	))
	slog.DebugContext(ctx,
		fmt.Sprintf("readFileName id=%v filetype=%v", id, filetype),
		"id", id,
		"filetype", filetype,
	)
	var out string
	d, err := filelists.GetId(id)
	if err != nil {
		config.TracerError(span, err)
		return ctx, errors.Wrap(err, fmt.Sprintf("filelists.GetId(%v)", id))
	}
	switch strings.ToLower(filetype) {
	case "zip":
		out = d.Zippass
	case "pdf":
		out = d.Pdfpass
	default:
		err := errors.New(fmt.Sprintf("eroor file type:%v", filetype))
		config.TracerError(span, err)
		return ctx, err
	}
	return contextWriteFileNameFIleType(ctx, out, filetype), nil
}

func readFileData(ctx context.Context) ([]byte, error) {
	ctx, span := config.TracerS(ctx, "readFileData", "read File Data")
	defer span.End()

	v, ok := contextReadFileNameFIleType(ctx)
	if !ok {
		return []byte{}, fmt.Errorf("error context Read FIledata")
	}
	filename := v.filename
	filetype := v.filetype
	span.AddEvent("readfile", trace.WithAttributes(
		attribute.String("filename", filename),
		attribute.String("filetype", filetype),
	))
	slog.DebugContext(ctx,
		fmt.Sprintf("readFileData filename=%v", filename),
		"filename", filename,
		"fietype", filetype,
	)
	var buffer []byte
	pass, err := createPass(filename, filetype)
	if err != nil {
		config.TracerError(span, err)
		return buffer, errors.Wrap(err, fmt.Sprintf("createPass(%v,%v)", filename, filetype))
	}
	if file, err := os.Open(pass); err != nil {
		config.TracerError(span, err)
		return buffer, errors.Wrap(err, fmt.Sprintf("os.Open(%v)", pass))
	} else {
		defer file.Close()
		buf := make([]byte, 1024)
		for {
			n, err := file.Read(buf)
			if n == 0 {
				break
			}
			if err != nil {
				config.TracerError(span, err)
				slog.ErrorContext(ctx,
					fmt.Sprintf("readFileData file.Read pass=%v", pass),
					"Pass", pass,
					"Error", err,
				)
				break
			}
			buffer = append(buffer, buf[:n]...)
		}
	}

	return buffer, nil
}

func createPass(name, filetype string) (string, error) {
	var out string
	switch strings.ToLower(filetype) {
	case "zip":
		out = zipFolder + name
	case "pdf":
		out = pdfFolder + name
	default:
		return out, errors.New(fmt.Sprintf("eroor file type:%v", filetype))

	}
	return out, nil
}

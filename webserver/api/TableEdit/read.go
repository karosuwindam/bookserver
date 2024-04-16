package tableedit

import (
	tableread "bookserver/webserver/api/TableRead"
	"net/http"
)

// urlでtableとidをしていて読み取る
// tablereadと同じ動きを希望
func GetReadId(w http.ResponseWriter, r *http.Request) {
	tableread.GetReadId(w, r)
}

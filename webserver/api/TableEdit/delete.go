package tableedit

import (
	tabledelete "bookserver/webserver/api/TableDelete"
	"net/http"
)

// Tabledelete内の処理と同じ動き
func DeleteTableById(w http.ResponseWriter, r *http.Request) {
	tabledelete.DeleteTableById(w, r)
}

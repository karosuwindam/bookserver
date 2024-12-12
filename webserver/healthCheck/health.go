package healthcheck

import (
	"bookserver/controller"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type HealthMessage struct {
	WebServer  interface{} `json:Webserver`
	Controller interface{} `json:Controller`
}

func Health(w http.ResponseWriter, r *http.Request) {
	var output HealthMessage
	ctx := r.Context()
	slog.InfoContext(ctx,
		fmt.Sprintf("%v %v", r.Method, r.URL),
		"Url", r.URL,
		"Method", r.Method,
	)
	output.WebServer = bool(true)
	output.Controller = controller.HealthCheck()

	if d, err := json.Marshal(&output); err != nil {
		slog.ErrorContext(ctx, "Health jsonConvert error",
			"error", err,
		)
	} else {
		w.Write(d)
	}
}

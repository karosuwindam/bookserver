package healthcheck

import (
	"bookserver/controller"
	"encoding/json"
	"log"
	"net/http"
)

type HealthMessage struct {
	WebServer  interface{} `json:Webserver`
	Controller interface{} `json:Controller`
}

func Health(w http.ResponseWriter, r *http.Request) {
	var output HealthMessage
	output.WebServer = bool(true)
	output.Controller = controller.HealthCheck()
	if d, err := json.Marshal(&output); err != nil {
		log.Println("error", err)
	} else {
		w.Write(d)
	}
}

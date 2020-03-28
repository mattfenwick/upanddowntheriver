package game

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type Responder interface {
	GetModel() string
	AddPlayer(player string) error
	//StartRound() error
}

func SetupHTTPServer(uiDirectory string, responder Responder) {
	http.Handle("/", http.FileServer(http.Dir(uiDirectory)))
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/model", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			string := responder.GetModel()
			header := w.Header()
			header.Set(http.CanonicalHeaderKey("content-type"), "application/json")
			fmt.Fprint(w, string)
		} else {
			http.NotFound(w, r)
		}
	})

	http.HandleFunc("/player", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Errorf("unable to read body: %+v", err)
				http.Error(w, err.Error(), 400)
				return
			}
			log.Debugf("received POST to /player with body %s", body)
			var bodyParams map[string]string
			err = json.Unmarshal(body, &bodyParams)
			if err != nil {
				log.Errorf("unable to unmarshal json: %+v", err)
				http.Error(w, err.Error(), 400)
				return
			}
			player, ok := bodyParams["Player"]
			if !ok || player == "" {
				log.Errorf("missing or empty Player parameter")
				http.Error(w, "missing or empty Player parameter", 400)
				return
			}
			err = responder.AddPlayer(player)
			if err != nil {
				log.Errorf("unable to add player: %+v", err)
				http.Error(w, err.Error(), 400)
				return
			}
			resp := "{\"Status\": \"Success\"}"
			log.Infof("added player %s", player)
			fmt.Fprint(w, resp)
		} else {
			http.NotFound(w, r)
		}
	})
}

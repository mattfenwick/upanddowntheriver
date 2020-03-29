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
	RemovePlayer(player string) error
	StartRound() error
	MakeWager(player string, hands int) error
	PlayCard(player string, card *Card) error
}

type StartRoundAction struct{}

type MakeWagerAction struct {
	Player string
	Hands  int
}

type PlayCardAction struct {
	Player string
	Card   *Card
}

type PlayerAction struct {
	StartRound *StartRoundAction
	MakeWager  *MakeWagerAction
	PlayCard   *PlayCardAction
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
			log.Errorf("verb %s not supported for /model", r.Method)
			http.NotFound(w, r)
		}
	})

	http.HandleFunc("/player", func(w http.ResponseWriter, r *http.Request) {
		verb := r.Method
		if verb == "POST" || verb == "DELETE" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Errorf("unable to read body: %+v", err)
				http.Error(w, err.Error(), 400)
				return
			}
			log.Debugf("received %s to /player with body %s", verb, body)
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
			if verb == "POST" {
				err = responder.AddPlayer(player)
			} else {
				err = responder.RemovePlayer(player)
			}
			if err != nil {
				log.Errorf("unable to add/remove player: %+v", err)
				http.Error(w, err.Error(), 400)
				return
			}
			resp := "{\"Status\": \"Success\"}"
			log.Infof("added player %s", player)
			fmt.Fprint(w, resp)
		} else {
			log.Errorf("verb %s not supported for /player", r.Method)
			http.NotFound(w, r)
		}
	})

	http.HandleFunc("/action", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Errorf("unable to read body: %+v", err)
				http.Error(w, err.Error(), 400)
				return
			}
			log.Debugf("received POST to /action with body %s", body)
			var action PlayerAction
			err = json.Unmarshal(body, &action)
			if err != nil {
				log.Errorf("unable to unmarshal json: %+v", err)
				http.Error(w, err.Error(), 400)
				return
			}

			var actionErr error
			if action.StartRound != nil {
				actionErr = responder.StartRound()
			} else if action.MakeWager != nil {
				actionErr = responder.MakeWager(action.MakeWager.Player, action.MakeWager.Hands)
			} else if action.PlayCard != nil {
				actionErr = responder.PlayCard(action.PlayCard.Player, action.PlayCard.Card)
			} else {
				http.Error(w, "action must have non-nil for either StartRound, MakeWager, or PlayCard", 400)
				return
			}
			if actionErr != nil {
				log.Errorf("unable to execute action: %+v", actionErr)
				http.Error(w, actionErr.Error(), 400)
				return
			}

			resp := "{\"Status\": \"Success\"}"
			log.Infof("handled action %+v", action)
			fmt.Fprint(w, resp)
		} else {
			log.Errorf("verb %s not supported for /action", r.Method)
			http.NotFound(w, r)
		}
	})
}

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
	GetPlayerModel(player string) (*PlayerModel, error)
	Join(player string) error
	RemovePlayer(player string) error
	SetCardsPerPlayer(count int) error
	StartRound() error
	StartHand() error
	MakeWager(player string, hands int) error
	PlayCard(player string, card *Card) error
	FinishRound() error
}

type GetPlayerModelAction struct{}

type JoinAction struct{}

type RemovePlayerAction struct {
	Player string
}

type MakeWagerAction struct {
	Hands int
}

type SetCardsPerPlayerAction struct {
	Count int
}

type StartHandAction struct{}

type StartRoundAction struct{}

type FinishRoundAction struct{}

type PlayerAction struct {
	Me                string
	GetModel          *GetPlayerModelAction
	Join              *JoinAction
	MakeWager         *MakeWagerAction
	PlayCard          *Card
	RemovePlayer      *RemovePlayerAction
	SetCardsPerPlayer *SetCardsPerPlayerAction
	StartHand         *StartHandAction
	StartRound        *StartRoundAction
	FinishRound       *FinishRoundAction
}

func SetupHTTPServer(uiDirectory string, responder Responder) {
	http.Handle("/", http.FileServer(http.Dir(uiDirectory)))
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/model", func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("receiving %s request to %s", r.Method, r.URL.String())
		if r.Method == "GET" {
			var response string
			var err error
			urlParams := r.URL.Query()
			if players, ok := urlParams["player"]; len(players) > 0 && ok {
				player := players[0]
				var pm *PlayerModel
				pm, err = responder.GetPlayerModel(player)
				if err != nil {
					log.Errorf("unable to get player %s model: %+v", player, err)
					http.Error(w, err.Error(), 400)
					return
				}
				var pmBytes []byte
				pmBytes, err = json.MarshalIndent(pm, "", "  ")
				if err != nil {
					log.Errorf("unable to serialize json: %+v", err)
					http.Error(w, err.Error(), 500)
					return
				}
				response = string(pmBytes)
			} else {
				response = responder.GetModel()
			}
			w.Header().Set(http.CanonicalHeaderKey("content-type"), "application/json")
			fmt.Fprint(w, response)
		} else {
			log.Errorf("verb %s not supported for /model", r.Method)
			http.NotFound(w, r)
		}
	})

	http.HandleFunc("/action", func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("receiving %s request to %s", r.Method, r.URL.String())
		if r.Method == "POST" {
			body, err := ioutil.ReadAll(r.Body)
			log.Debugf("received body %s", string(body))
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
			if action.GetModel != nil {
				actionErr = nil // nothing else to do!
				// just let the playerModel be grabbed down below
			} else if action.Join != nil {
				actionErr = responder.Join(action.Me)
			} else if action.RemovePlayer != nil {
				actionErr = responder.RemovePlayer(action.RemovePlayer.Player)
			} else if action.SetCardsPerPlayer != nil {
				actionErr = responder.SetCardsPerPlayer(action.SetCardsPerPlayer.Count)
			} else if action.StartRound != nil {
				actionErr = responder.StartRound()
			} else if action.MakeWager != nil {
				actionErr = responder.MakeWager(action.Me, action.MakeWager.Hands)
			} else if action.PlayCard != nil {
				actionErr = responder.PlayCard(action.Me, &Card{Suit: action.PlayCard.Suit, Number: action.PlayCard.Number})
			} else if action.StartHand != nil {
				actionErr = responder.StartHand()
			} else if action.FinishRound != nil {
				actionErr = responder.FinishRound()
			} else {
				http.Error(w, "action must have non-nil for one of GetModel, Join, StartRound, MakeWager, RemovePlayer, SetCardsPerPlayer, or PlayCard", 400)
				return
			}
			if actionErr != nil {
				log.Errorf("unable to execute action: %+v", actionErr)
				http.Error(w, actionErr.Error(), 400)
				return
			}

			pm, err := responder.GetPlayerModel(action.Me)
			if err != nil {
				log.Errorf("unable to get player %s model: %+v", action.Me, err)
				http.Error(w, err.Error(), 400)
				return
			}
			pmBytes, err := json.MarshalIndent(pm, "", "  ")
			if err != nil {
				log.Errorf("unable to serialize json: %+v", err)
				http.Error(w, err.Error(), 500)
				return
			}

			log.Infof("handled action %+v", action)
			log.Tracef("response %s", string(pmBytes))
			fmt.Fprint(w, string(pmBytes))
		} else {
			log.Errorf("verb %s not supported for /action", r.Method)
			http.NotFound(w, r)
		}
	})
}

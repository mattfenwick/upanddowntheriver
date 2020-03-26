package game

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Action struct {
	Name  string
	Apply func() error
}

type Game struct {
	Players map[string]bool
	Rounds  []*Round
	Stop    <-chan struct{}
	Actions chan *Action
}

func NewGame(stop <-chan struct{}) *Game {
	game := &Game{
		Players: nil,
		Rounds:  nil,
		Stop:    stop,
		Actions: make(chan *Action),
	}
	go func() {
		game.startActionProcessor()
	}()
	return game
}

func (game *Game) startActionProcessor() {
	for {
		var action *Action
		select {
		case <-game.Stop:
			break
		case action = <-game.Actions:
		}

		err := action.Apply()
		if err != nil {
			log.Errorf("unable to process action type %s: %s", action.Name, err)
		} else {
			log.Infof("successfully processed action type %s", action.Name)
		}
	}
}

// mutators

func (game *Game) AddPlayer(player string) error {
	done := make(chan error)
	game.Actions <- &Action{"addPlayer", func() error {
		var err error
		if game.Players[player] {
			err = errors.New(fmt.Sprintf("can't add player %s, already present", player))
		} else {
			game.Players[player] = true
		}
		go func() {
			done <- err
		}()
		return err
	}}
	return <-done
}

func (game *Game) RemovePlayer(player string) {
	// TODO
}

func (game *Game) StartRound() {
	// TODO
}

func (game *Game) Deal() {
	// TODO
}

func (game *Game) MakeWager(player string, hands int) {
	// TODO
}

func (game *Game) PlayCard(player string, card *Card) {
	// TODO
}

// getters

func (game *Game) GetGameModel() {
	// TODO players
}

func (game *Game) GetRoundModel() {
	// TODO player order, dealer, trump suit, hands, wagers
}

func (game *Game) GetHandModel() {
	// TODO suit, cards played, leader, leader card
}

func (game *Game) GetHandResults() {
	// TODO winner, cards played
}

func (game *Game) GetRoundResults() {
	// TODO wagers, winners, losers
}

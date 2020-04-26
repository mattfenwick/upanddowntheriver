package game

import (
	"encoding/json"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Action struct {
	Name  string
	Apply func() error
}

type GameConcurrencyWrapper struct {
	Game    *Game
	Stop    <-chan struct{}
	Actions chan *Action
}

func NewGameConcurrencyWrapper(game *Game, stop <-chan struct{}) *GameConcurrencyWrapper {
	gcw := &GameConcurrencyWrapper{
		Game:    game,
		Stop:    stop,
		Actions: make(chan *Action),
	}
	go func() {
		gcw.startActionProcessor()
	}()
	return gcw
}

func (gcw *GameConcurrencyWrapper) startActionProcessor() {
	for {
		var action *Action
		select {
		case <-gcw.Stop:
			break
		case action = <-gcw.Actions:
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

func (gcw *GameConcurrencyWrapper) SetDeck() error {
	return errors.New("TODO")
}

func (gcw *GameConcurrencyWrapper) SetCardsPerPlayer(count int) error {
	done := make(chan error)
	gcw.Actions <- &Action{"setCardsPerPlayer", func() error {
		err := gcw.Game.setCardsPerPlayer(count)
		go func() {
			done <- err
		}()
		return err
	}}
	return <-done
}

func (gcw *GameConcurrencyWrapper) SetDeckType(deckType DeckType) error {
	done := make(chan error)
	gcw.Actions <- &Action{"setDeckType", func() error {
		err := gcw.Game.setDeckType(deckType)
		go func() {
			done <- err
		}()
		return err
	}}
	return <-done
}

func (gcw *GameConcurrencyWrapper) Join(player string) (string, error) {
	done := make(chan struct{})
	var err error
	var addedPlayer string
	gcw.Actions <- &Action{"join", func() error {
		addedPlayer, err = gcw.Game.join(player)
		close(done)
		return err
	}}
	<-done
	return addedPlayer, err
}

func (gcw *GameConcurrencyWrapper) RemovePlayer(player string) error {
	done := make(chan error)
	gcw.Actions <- &Action{"removePlayer", func() error {
		err := gcw.Game.removePlayer(player)
		go func() {
			done <- err
		}()
		return err
	}}
	return <-done
}

func (gcw *GameConcurrencyWrapper) StartRound() error {
	done := make(chan error)
	gcw.Actions <- &Action{"startRound", func() error {
		err := gcw.Game.startRound()
		go func() {
			done <- err
		}()
		return err
	}}
	return <-done
}

func (gcw *GameConcurrencyWrapper) FinishRound() error {
	done := make(chan error)
	gcw.Actions <- &Action{"finishRound", func() error {
		err := gcw.Game.finishRound()
		go func() {
			done <- err
		}()
		return err
	}}
	return <-done
}

func (gcw *GameConcurrencyWrapper) MakeWager(player string, hands int) error {
	done := make(chan error)
	gcw.Actions <- &Action{"makeWager", func() error {
		err := gcw.Game.makeWager(player, hands)
		go func() {
			done <- err
		}()
		return err
	}}
	return <-done
}

func (gcw *GameConcurrencyWrapper) PlayCard(player string, card *Card) error {
	done := make(chan error)
	gcw.Actions <- &Action{"playCard", func() error {
		err := gcw.Game.playCard(player, card)
		go func() {
			done <- err
		}()
		return err
	}}
	return <-done
}

// getters

func (gcw *GameConcurrencyWrapper) GetModel() string {
	done := make(chan string)
	gcw.Actions <- &Action{"getJsonModel", func() error {
		bytes, err := json.MarshalIndent(gcw.Game, "", "  ")
		if err != nil {
			panic(err)
		}
		go func() {
			done <- string(bytes)
		}()
		return nil
	}}
	return <-done
}

func (gcw *GameConcurrencyWrapper) GetPlayerModel(player string) *PlayerModel {
	done := make(chan struct{})
	var pm *PlayerModel
	gcw.Actions <- &Action{"getJsonModel", func() error {
		pm = gcw.Game.playerModel(player)
		close(done)
		return nil
	}}
	<-done
	return pm
}

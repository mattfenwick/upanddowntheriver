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

func (gcw *GameConcurrencyWrapper) ReorderPlayers(players []string) error {
	return errors.New("TODO")
}

func (gcw *GameConcurrencyWrapper) SetDeck() error {
	return errors.New("TODO")
}

func (gcw *GameConcurrencyWrapper) SetCardsPerPlayer(count int) error {
	//// 1 = for the trump suit
	//cardsNeeded := cardsPerPlayer*len(players) + 1
	//cardsAvailable := len(Cards(deck))
	//if cardsNeeded > cardsAvailable {
	//	return nil, errors.New(fmt.Sprintf("need %d cards for %d players, a total of %d -- more than the %d available", cardsPerPlayer, len(players), cardsNeeded, cardsAvailable))
	//}
	return errors.New("TODO")
}

func (gcw *GameConcurrencyWrapper) AddPlayer(player string) error {
	done := make(chan error)
	gcw.Actions <- &Action{"addPlayer", func() error {
		err := gcw.Game.addPlayer(player)
		go func() {
			done <- err
		}()
		return err
	}}
	return <-done
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

//type GameModel struct {
//	Players []string
//}
//
//func (gcw *GameConcurrencyWrapper) GetGameModel() *GameModel {
//	done := make(chan *GameModel)
//	gcw.Actions <- &Action{
//		Name: "getGameModel",
//		Apply: func() error {
//			players := []string{}
//			for player := range game.Players {
//				players = append(players, player)
//			}
//			done <- &GameModel{Players: players}
//			return nil
//		},
//	}
//	return <-done
//}
//
//type RoundModel struct {
//	// PlayerOrder implies Dealer -- last player
//	PlayerOrder []string
//	TrumpSuit   string
//	Wagers      map[string]int
//	Hands       []*HandModel
//	CurrentHand *HandModel
//}
//
//func (gcw *GameConcurrencyWrapper) GetRoundModel() {
//	// TODO player order, dealer, trump suit, hands, wagers
//}
//
//type HandModel struct {
//	Suit        string
//	CardsPlayed map[string]*Card
//	Leader      string
//	LeaderCard  *Card
//	NextPlayer  string
//	Hand
//}
//
//func (gcw *GameConcurrencyWrapper) GetHandModel() {
//	// TODO suit, cards played, leader, leader card
//}
//
//type HandResults struct {
//	CardsPlayed map[string]*Card
//	Winners     []string
//	Losers      []string
//}
//
//func (gcw *GameConcurrencyWrapper) GetHandResults() {
//	// TODO winner, cards played
//}
//
//func (gcw *GameConcurrencyWrapper) GetRoundResults() {
//	// TODO wagers, winners, losers
//}

package game

import (
	"fmt"
	"github.com/pkg/errors"
)

type GameState int

const (
	GameStateSetup           GameState = iota
	GameStateRoundInProgress GameState = iota
)

func (g GameState) String() string {
	switch g {
	case GameStateSetup:
		return "GameStateSetup"
	case GameStateRoundInProgress:
		return "GameStateRoundInProgress"
	}
	panic(fmt.Errorf("invalid GameState value: %d", g))
}

func (g GameState) MarshalJSON() ([]byte, error) {
	jsonString := fmt.Sprintf(`"%s"`, g.String())
	return []byte(jsonString), nil
}

func (g GameState) MarshalText() (text []byte, err error) {
	return []byte(g.String()), nil
}

type Game struct {
	Players        []string
	PlayersSet     map[string]bool
	Deck           Deck
	CardsPerPlayer int
	Rounds         []*Round
	State          GameState
}

func NewGame() *Game {
	game := &Game{
		Players:        []string{},
		PlayersSet:     map[string]bool{},
		Deck:           NewStandardDeck(),
		CardsPerPlayer: 1,
		Rounds:         nil,
		State:          GameStateSetup,
	}
	return game
}

// mutators

func (game *Game) addPlayer(player string) error {
	if game.State != GameStateSetup {
		return errors.New(fmt.Sprintf("can't add player, in state %s", game.State.String()))
	} else if game.PlayersSet[player] {
		return errors.New(fmt.Sprintf("can't add player %s, already present", player))
	} else {
		game.Players = append(game.Players, player)
		game.PlayersSet[player] = true
		return nil
	}
}

func (game *Game) removePlayer(player string) error {
	if game.State != GameStateSetup {
		return errors.New(fmt.Sprintf("can't remove player, in state %s", game.State.String()))
	} else if !game.PlayersSet[player] {
		return errors.New(fmt.Sprintf("can't remove player %s, not present", player))
	} else {
		delete(game.PlayersSet, player)
		players := []string{}
		for _, player := range game.Players {
			if _, ok := game.PlayersSet[player]; ok {
				players = append(players, player)
			}
		}
		game.Players = players
		return nil
	}
}

func (game *Game) startRound() error {
	if game.State != GameStateSetup {
		return errors.New(fmt.Sprintf("can't start round, in state %s", game.State.String()))
	}
	playerCount := len(game.Players)
	if playerCount < 2 {
		return errors.New(fmt.Sprintf("can't start game with fewer than 2 players, found %d", playerCount))
	}
	players := append([]string{}, game.Players...)
	game.Rounds = append(game.Rounds, NewRound(players, game.Deck, game.CardsPerPlayer))
	game.State = GameStateRoundInProgress
	return nil
}

func (game *Game) finishRound() error {
	if game.State != GameStateRoundInProgress {
		return errors.New(fmt.Sprintf("can't finish round, in state %s", game.State.String()))
	} else {
		// TODO any other cleanup?
		game.State = GameStateSetup
		return nil
	}
}

func (game *Game) currentRound() (*Round, error) {
	if game.State != GameStateRoundInProgress {
		return nil, errors.New(fmt.Sprintf("can't get current round, game in state %s", game.State.String()))
	}
	return game.Rounds[len(game.Rounds)-1], nil
}

func (game *Game) makeWager(player string, hands int) error {
	round, err := game.currentRound()
	if err != nil {
		return err
	}
	return round.Wager(player, hands)
}

func (game *Game) playCard(player string, card *Card) error {
	round, err := game.currentRound()
	if err != nil {
		return err
	}
	return round.PlayCard(player, card)
}

// getters

//type GameModel struct {
//	Players []string
//}
//
//func (game *Game) GetGameModel() *GameModel {
//	done := make(chan *GameModel)
//	game.actions <- &Action{
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
//func (game *Game) GetRoundModel() {
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
//func (game *Game) GetHandModel() {
//	// TODO suit, cards played, leader, leader card
//}
//
//type HandResults struct {
//	CardsPlayed map[string]*Card
//	Winners     []string
//	Losers      []string
//}
//
//func (game *Game) GetHandResults() {
//	// TODO winner, cards played
//}
//
//func (game *Game) GetRoundResults() {
//	// TODO wagers, winners, losers
//}

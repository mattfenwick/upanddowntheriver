package game

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
	FinishedRounds []*Round
	CurrentRound   *Round
	State          GameState
}

func NewGame() *Game {
	game := &Game{
		Players:        []string{},
		PlayersSet:     map[string]bool{},
		Deck:           NewStandardDeck(),
		CardsPerPlayer: 1,
		FinishedRounds: []*Round{},
		CurrentRound:   nil,
		State:          GameStateSetup,
	}
	return game
}

// mutators

func (game *Game) addPlayer(player string) error {
	if game.State != GameStateSetup {
		return errors.New(fmt.Sprintf("can't add player %s, in state %s", player, game.State.String()))
	} else if game.PlayersSet[player] {
		return errors.New(fmt.Sprintf("can't add player %s, already present", player))
	} else {
		game.Players = append(game.Players, player)
		game.PlayersSet[player] = true
		maxCardsPerPlayer := len(Cards(game.Deck)) / len(game.Players)
		if game.CardsPerPlayer > maxCardsPerPlayer {
			game.CardsPerPlayer = maxCardsPerPlayer
		}
		return nil
	}
}

func (game *Game) join(player string) (string, error) {
	if player == "" {
		return "", errors.New("invalid name: empty")
	}
	if len(player) > 20 {
		// just take the first 20 characters so as not to get overwhelmed by excessively long names
		shortName := player[:20]
		log.Infof("player name <%s> too long, truncating to <%s>", player, shortName)
		player = shortName
	}
	// if player's already in the game, nothing to do!
	if game.PlayersSet[player] {
		return player, nil
	}
	if game.State != GameStateSetup {
		return "", errors.New(fmt.Sprintf("can't join as %s, in state %s", player, game.State.String()))
	}
	return player, game.addPlayer(player)
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

func (game *Game) playerModel(player string) *PlayerModel {
	return newPlayerModel(game, player)
}

func (game *Game) setCardsPerPlayer(count int) error {
	if game.State != GameStateSetup {
		return errors.New(fmt.Sprintf("can't set cards per player, in state %s", game.State.String()))
	}
	maxCardsPerPlayer := len(Cards(game.Deck)) / len(game.Players)
	if count > maxCardsPerPlayer {
		return errors.New(fmt.Sprintf("requested cardsPerPlayer of %d, which is greater than the maxCardsPerPlayer of %d", count, maxCardsPerPlayer))
	}
	game.CardsPerPlayer = count
	return nil
}

func (game *Game) setDeckType(deckType DeckType) error {
	if game.State != GameStateSetup {
		return errors.New(fmt.Sprintf("can't set deck type, in state %s", game.State.String()))
	}
	switch deckType {
	case DeckTypeStandard:
		game.Deck = NewStandardDeck()
	case DeckTypeDoubleStandard:
		game.Deck = NewDoubleStandardDeck()
	case DeckTypeDeterministicStandard:
		game.Deck = NewDeterministicShuffleDeck()
	default:
		return errors.New(fmt.Sprintf("invalid deck type %s", deckType))
	}
	return nil
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
	game.CurrentRound = NewRound(players, game.Deck, game.CardsPerPlayer)
	game.State = GameStateRoundInProgress
	return nil
}

func (game *Game) finishRound() error {
	if game.State != GameStateRoundInProgress {
		return errors.New(fmt.Sprintf("can't finish round, in state %s", game.State.String()))
	} else {
		game.FinishedRounds = append(game.FinishedRounds, game.CurrentRound)
		game.CurrentRound = nil
		game.State = GameStateSetup
		// move the first player to the end
		game.Players = append(game.Players[1:], game.Players[0])
		return nil
	}
}

func (game *Game) makeWager(player string, hands int) error {
	if game.State != GameStateRoundInProgress {
		return errors.New(fmt.Sprintf("can't make wager, game in state %s", game.State.String()))
	}
	return game.CurrentRound.Wager(player, hands)
}

func (game *Game) playCard(player string, card *Card) error {
	if game.State != GameStateRoundInProgress {
		return errors.New(fmt.Sprintf("can't play card, game in state %s", game.State.String()))
	}
	return game.CurrentRound.PlayCard(player, card)
}

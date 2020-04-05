package game

import (
	"fmt"
	"github.com/pkg/errors"
	"sort"
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

func (game *Game) join(player string) error {
	if game.State != GameStateSetup {
		return errors.New(fmt.Sprintf("can't join as %s, in state %s", player, game.State.String()))
	}
	// if player's already in the game, nothing to do!
	if game.PlayersSet[player] {
		return nil
	}
	return game.addPlayer(player)
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

func (game *Game) setCardsPerPlayer(count int) error {
	maxCardsPerPlayer := len(Cards(game.Deck)) / len(game.Players)
	if count > maxCardsPerPlayer {
		return errors.New(fmt.Sprintf("requested cardsPerPlayer of %d, which is greater than the maxCardsPerPlayer of %d", count, maxCardsPerPlayer))
	}
	game.CardsPerPlayer = count
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

func (game *Game) startHand() error {
	if game.State != GameStateRoundInProgress {
		return errors.New(fmt.Sprintf("can't start hand, in state %s", game.State.String()))
	}
	return game.CurrentRound.StartHand()
}

func (game *Game) finishRound() error {
	if game.State != GameStateRoundInProgress {
		return errors.New(fmt.Sprintf("can't finish round, in state %s", game.State.String()))
	} else {
		game.FinishedRounds = append(game.FinishedRounds, game.CurrentRound)
		game.CurrentRound = nil
		game.State = GameStateSetup
		return nil
	}
}

func (game *Game) makeWager(player string, hands int) error {
	if game.State != GameStateRoundInProgress {
		return errors.New(fmt.Sprintf("can't get current round, game in state %s", game.State.String()))
	}
	return game.CurrentRound.Wager(player, hands)
}

func (game *Game) playCard(player string, card *Card) error {
	if game.State != GameStateRoundInProgress {
		return errors.New(fmt.Sprintf("can't get current round, game in state %s", game.State.String()))
	}
	return game.CurrentRound.PlayCard(player, card)
}

// getters

type PlayerState int

const (
	PlayerStateGameWaitingForPlayers PlayerState = iota
	PlayerStateGameReady             PlayerState = iota
	PlayerStateRoundWagerTurn        PlayerState = iota
	PlayerStateRoundHandReady        PlayerState = iota
	PlayerStateHandPlayTurn          PlayerState = iota
	PlayerStateRoundFinished         PlayerState = iota
)

func (p PlayerState) JSONString() string {
	switch p {
	case PlayerStateGameWaitingForPlayers:
		return "WaitingForPlayers"
	//case PlayerStateGameReady:
	//	return "Ready"
	case PlayerStateRoundWagerTurn:
		return "RoundWagerTurn"
	case PlayerStateRoundHandReady:
		return "RoundHandReady"
	case PlayerStateHandPlayTurn:
		return "HandPlayTurn"
	case PlayerStateRoundFinished:
		return "RoundFinished"
	}
	panic(fmt.Errorf("invalid PlayerState value: %d", p))
}

func (p PlayerState) MarshalJSON() ([]byte, error) {
	jsonString := fmt.Sprintf(`"%s"`, p.JSONString())
	return []byte(jsonString), nil
}

func (p PlayerState) MarshalText() (text []byte, err error) {
	return []byte(p.JSONString()), nil
}

type PlayerGame struct {
	Players        []string
	CardsPerPlayer int
}

type PlayedCard struct {
	Player string
	Card   *Card
}

type PlayerHand struct {
	Cards       []*Card
	Suit        string
	Leader      string
	LeaderCard  *Card
	CardsPlayed []*PlayedCard
	NextPlayer  string
}

type PlayerWager struct {
	Player   string
	Count    *int
	HandsWon *int
}

type PlayerRound struct {
	Cards           []*Card
	Wagers          []*PlayerWager
	TrumpSuit       string
	NextWagerPlayer string
	WagerSum        int
}

type PlayerModel struct {
	Me    string
	State PlayerState
	Game  *PlayerGame
	Round *PlayerRound
	Hand  *PlayerHand
}

func (game *Game) playerModel(player string) (*PlayerModel, error) {
	if _, ok := game.PlayersSet[player]; !ok {
		return nil, errors.New(fmt.Sprintf("player %s not found", player))
	}

	var state PlayerState
	var round *PlayerRound
	var hand *PlayerHand
	switch game.State {
	case GameStateSetup:
		state = PlayerStateGameWaitingForPlayers
		break
	case GameStateRoundInProgress:
		cards := []*Card{}
		for _, pc := range game.CurrentRound.Players[player] {
			if !pc.IsPlayed {
				cards = append(cards, pc.Card)
			}
		}
		// let's sort the cards numerically ascending, then break ties with suits
		sort.Slice(cards, func(i, j int) bool {
			return game.Deck.Compare(cards[i], cards[j]) < 0
		})
		playerWins := map[string]int{}
		for _, hand := range game.CurrentRound.FinishedHands {
			if _, ok := playerWins[hand.Leader]; !ok {
				playerWins[hand.Leader] = 0
			}
			playerWins[hand.Leader]++
		}
		wagers := []*PlayerWager{}
		for _, p := range game.CurrentRound.PlayersOrder {
			var wager *int
			count, ok := game.CurrentRound.Wagers[p]
			if ok {
				wager = &count
			}
			var handsWon *int
			if won, ok := playerWins[p]; ok {
				handsWon = &won
			}
			wagers = append(wagers, &PlayerWager{
				Player:   p,
				Count:    wager,
				HandsWon: handsWon,
			})
		}
		round = &PlayerRound{
			Cards:           cards,
			Wagers:          wagers,
			TrumpSuit:       game.CurrentRound.TrumpSuit,
			NextWagerPlayer: "",
			WagerSum:        game.CurrentRound.WagerSum,
		}
		for _, player := range game.CurrentRound.PlayersOrder {
			if _, ok := game.CurrentRound.Wagers[player]; !ok {
				round.NextWagerPlayer = player
				break
			}
		}
		switch game.CurrentRound.State {
		case RoundStateCardsDealt:
			state = PlayerStateRoundWagerTurn
			break
		case RoundStateWagersMade:
			state = PlayerStateRoundHandReady
			break
		case RoundStateHandInProgress:
			state = PlayerStateHandPlayTurn
			ch := game.CurrentRound.CurrentHand
			cardsPlayed := []*PlayedCard{}
			nextPlayer := ""
			for _, p := range ch.PlayersOrder {
				pc := &PlayedCard{Player: p, Card: nil}
				card, ok := ch.CardsPlayed[p]
				if ok {
					pc.Card = card
				} else if !ok && nextPlayer == "" {
					nextPlayer = p
				}
				cardsPlayed = append(cardsPlayed, pc)
			}
			hand = &PlayerHand{
				Cards:       cards,
				Suit:        ch.Suit,
				Leader:      ch.Leader,
				LeaderCard:  ch.LeaderCard,
				CardsPlayed: cardsPlayed,
				NextPlayer:  nextPlayer,
			}
			break
		case RoundStateFinished:
			state = PlayerStateRoundFinished
			break
		}
	}
	model := &PlayerModel{
		Me:    player,
		State: state,
		Game: &PlayerGame{
			Players:        game.Players,
			CardsPerPlayer: game.CardsPerPlayer,
		},
		Round: round,
		Hand:  hand,
	}
	return model, nil
}

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

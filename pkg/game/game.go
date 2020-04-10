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

// getters

func (game *Game) playerModel(player string) (*PlayerModel, error) {
	if _, ok := game.PlayersSet[player]; !ok && player != "" {
		return nil, errors.New(fmt.Sprintf("player %s not found", player))
	}

	var state PlayerState
	var status *Status
	var myCards []*Card
	if player != "" {
		switch game.State {
		case GameStateSetup:
			state = PlayerStateWaitingForPlayers
			break
		case GameStateRoundInProgress:
			state, status, myCards = game.playerStatusAndCards(player)
		}
	} else {
		state = PlayerStateNotJoined
	}
	model := &PlayerModel{
		Me:    player,
		State: state,
		Game: &PlayerGame{
			Players:        game.Players,
			CardsPerPlayer: game.CardsPerPlayer,
		},
		Status:  status,
		MyCards: myCards,
	}
	return model, nil
}

func (game *Game) playerStatusAndCards(player string) (PlayerState, *Status, []*Card) {
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
	playerStatuses := []*PlayerStatus{}
	var prevHand *Hand
	if len(game.CurrentRound.FinishedHands) > 0 {
		prevHand = game.CurrentRound.FinishedHands[len(game.CurrentRound.FinishedHands)-1]
	}
	currHand := game.CurrentRound.CurrentHand
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
		ps := &PlayerStatus{
			Player:       p,
			Wager:        wager,
			HandsWon:     handsWon,
			PreviousCard: nil,
			CurrentCard:  nil,
		}
		if prevHand != nil {
			ps.PreviousCard = prevHand.CardsPlayed[p]
		}
		if currHand != nil {
			ps.CurrentCard = currHand.CardsPlayed[p]
		}
		playerStatuses = append(playerStatuses, ps)
	}
	status := &Status{
		PlayerStatuses: playerStatuses,
		TrumpSuit:      game.CurrentRound.TrumpSuit,
		WagerSum:       game.CurrentRound.WagerSum,
	}
	if prevHand != nil {
		status.PreviousHand = &PreviousHand{
			Suit:   prevHand.Suit,
			Winner: prevHand.Leader,
		}
	}
	for _, player := range game.CurrentRound.PlayersOrder {
		if _, ok := game.CurrentRound.Wagers[player]; !ok {
			status.NextWagerPlayer = player
			break
		}
	}

	var state PlayerState
	switch game.CurrentRound.State {
	case RoundStateWagers:
		state = PlayerStateWagerTurn
		break
	case RoundStateHandInProgress:
		state = PlayerStatePlayCardTurn
		ch := game.CurrentRound.CurrentHand
		nextPlayer := ""
		for _, p := range ch.PlayersOrder {
			_, ok := ch.CardsPlayed[p]
			if !ok && nextPlayer == "" {
				nextPlayer = p
				break
			}
		}
		status.CurrentHand = &CurrentHand{
			Suit:       ch.Suit,
			Leader:     ch.Leader,
			LeaderCard: ch.LeaderCard,
			NextPlayer: nextPlayer,
		}
		break
	case RoundStateFinished:
		state = PlayerStateRoundFinished
		break
	}

	if game.CurrentRound.State == RoundStateHandInProgress || game.CurrentRound.State == RoundStateFinished {
		zero := 0
		for _, s := range status.PlayerStatuses {
			if s.HandsWon == nil {
				s.HandsWon = &zero
			}
		}
	}

	return state, status, cards
}

package game

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"sort"
)

type PlayerState int

const (
	PlayerStateNotJoined         PlayerState = iota
	PlayerStateWaitingForPlayers PlayerState = iota
	PlayerStateWagerTurn         PlayerState = iota
	PlayerStatePlayCardTurn      PlayerState = iota
	PlayerStateRoundFinished     PlayerState = iota
)

func (p PlayerState) JSONString() string {
	switch p {
	case PlayerStateNotJoined:
		return "NotJoined"
	case PlayerStateWaitingForPlayers:
		return "WaitingForPlayers"
	case PlayerStateWagerTurn:
		return "WagerTurn"
	case PlayerStatePlayCardTurn:
		return "PlayCardTurn"
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

func parsePlayerState(text string) (PlayerState, error) {
	switch text {
	case "NotJoined":
		return PlayerStateNotJoined, nil
	case "WaitingForPlayers":
		return PlayerStateWaitingForPlayers, nil
	case "WagerTurn":
		return PlayerStateWagerTurn, nil
	case "PlayCardTurn":
		return PlayerStatePlayCardTurn, nil
	case "RoundFinished":
		return PlayerStateRoundFinished, nil
	}
	return PlayerStateWaitingForPlayers, errors.New(fmt.Sprintf("unable to parse player state %s", text))
}

func (p *PlayerState) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	status, err := parsePlayerState(str)
	if err != nil {
		return err
	}
	*p = status
	return nil
}

func (p *PlayerState) UnmarshalText(text []byte) (err error) {
	status, err := parsePlayerState(string(text))
	if err != nil {
		return err
	}
	*p = status
	return nil
}

type PlayerGame struct {
	Players        []string
	CardsPerPlayer int
}

type CurrentHand struct {
	Suit       string
	Leader     string
	LeaderCard *Card
	NextPlayer string
}

type PlayerStatus struct {
	Player           string
	IsNextWagerer    bool
	IsNextPlayer     bool
	IsCurrentLeader  bool
	IsPreviousWinner bool
	Wager            *int
	HandsWon         *int
	PreviousCard     *Card
	CurrentCard      *Card
}

type PreviousHand struct {
	Suit   string
	Winner string
}

type Status struct {
	PlayerStatuses  []*PlayerStatus
	TrumpSuit       string
	NextWagerPlayer string
	WagerSum        int
	PreviousHand    *PreviousHand
	CurrentHand     *CurrentHand
}

type PlayerModel struct {
	Me      string
	State   PlayerState
	Game    *PlayerGame
	Status  *Status
	MyCards []*Card
}

func newPlayerModel(game *Game, player string) (*PlayerModel, error) {
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
			state, status, myCards = playerStatusAndCards(game, player)
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

func playerStatusAndCards(game *Game, player string) (PlayerState, *Status, []*Card) {
	// get my cards
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
	var prevHand *Hand
	if len(game.CurrentRound.FinishedHands) > 0 {
		prevHand = game.CurrentRound.FinishedHands[len(game.CurrentRound.FinishedHands)-1]
	}
	currHand := game.CurrentRound.CurrentHand

	var nextWagerPlayer string
	for _, player := range game.CurrentRound.PlayersOrder {
		if _, ok := game.CurrentRound.Wagers[player]; !ok {
			nextWagerPlayer = player
			break
		}
	}
	playerStatuses := []*PlayerStatus{}
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
			Player:        p,
			IsNextWagerer: p == nextWagerPlayer,
			Wager:         wager,
			HandsWon:      handsWon,
			PreviousCard:  nil,
			CurrentCard:   nil,
		}
		if prevHand != nil {
			ps.PreviousCard = prevHand.CardsPlayed[p]
			ps.IsPreviousWinner = prevHand.Leader == p
		}
		if currHand != nil {
			ps.CurrentCard = currHand.CardsPlayed[p]
			ps.IsCurrentLeader = currHand.Leader == p
		}
		playerStatuses = append(playerStatuses, ps)
	}
	status := &Status{
		PlayerStatuses:  playerStatuses,
		TrumpSuit:       game.CurrentRound.TrumpSuit,
		WagerSum:        game.CurrentRound.WagerSum,
		NextWagerPlayer: nextWagerPlayer,
	}
	if prevHand != nil {
		status.PreviousHand = &PreviousHand{
			Suit:   prevHand.Suit,
			Winner: prevHand.Leader,
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
		for _, ps := range status.PlayerStatuses {
			if ps.Player == nextPlayer {
				ps.IsNextPlayer = true
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

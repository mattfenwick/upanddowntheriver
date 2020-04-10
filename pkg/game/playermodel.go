package game

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
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
	Player       string
	Wager        *int
	HandsWon     *int
	PreviousCard *Card
	CurrentCard  *Card
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

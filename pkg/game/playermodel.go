package game

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

type PlayerState int

const (
	PlayerStateNotJoined             PlayerState = iota
	PlayerStateGameWaitingForPlayers PlayerState = iota
	//PlayerStateGameReady             PlayerState = iota
	PlayerStateRoundWagerTurn PlayerState = iota
	PlayerStateRoundHandReady PlayerState = iota
	PlayerStateHandPlayTurn   PlayerState = iota
	PlayerStateHandFinished   PlayerState = iota
	PlayerStateRoundFinished  PlayerState = iota
)

func (p PlayerState) JSONString() string {
	switch p {
	case PlayerStateNotJoined:
		return "NotJoined"
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
	case PlayerStateHandFinished:
		return "HandFinished"
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
		return PlayerStateGameWaitingForPlayers, nil
	//case PlayerStateGameReady:
	//	return "Ready"
	case "RoundWagerTurn":
		return PlayerStateRoundWagerTurn, nil
	case "RoundHandReady":
		return PlayerStateRoundHandReady, nil
	case "HandPlayTurn":
		return PlayerStateHandPlayTurn, nil
	case "HandFinished":
		return PlayerStateHandFinished, nil
	case "RoundFinished":
		return PlayerStateRoundFinished, nil
	}
	return PlayerStateGameWaitingForPlayers, errors.New(fmt.Sprintf("unable to parse player state %s", text))
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

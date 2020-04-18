package game

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
	IsMe             bool
	IsNextWagerer    bool
	IsNextPlayer     bool
	IsCurrentLeader  bool
	IsPreviousWinner bool
	Mood             PlayerMood
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

func newPlayerModel(game *Game, player string) *PlayerModel {
	pg := &PlayerGame{
		Players:        game.Players,
		CardsPerPlayer: game.CardsPerPlayer,
	}
	// empty player, or player not found?  we'll only let them see who's playing and the game config
	if _, ok := game.PlayersSet[player]; !ok {
		return &PlayerModel{
			State: PlayerStateNotJoined,
			Game:  pg,
		}
	}

	var state PlayerState
	var status *Status
	var myCards []*Card
	switch game.State {
	case GameStateSetup:
		state = PlayerStateWaitingForPlayers
		break
	case GameStateRoundInProgress:
		state, status, myCards = playerStatusAndCards(game, player)
	}
	return &PlayerModel{
		Me:      player,
		State:   state,
		Game:    pg,
		Status:  status,
		MyCards: myCards,
	}
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
			IsMe:          p == player,
			IsNextWagerer: p == nextWagerPlayer,
			Wager:         wager,
			HandsWon:      handsWon,
			Mood:          PlayerMoodNone,
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
			// figure out the mood
			handsFinished := len(game.CurrentRound.FinishedHands)
			totalHands := game.CurrentRound.CardsPerPlayer
			isWinningCurrentHand := (game.CurrentRound.CurrentHand != nil) && (game.CurrentRound.CurrentHand.Leader == s.Player)
			hasPlayedForCurrentHand := s.CurrentCard != nil
			mood := playerMood(*s.Wager, *s.HandsWon, handsFinished, totalHands, hasPlayedForCurrentHand, isWinningCurrentHand)
			log.Debugf("player %s mood for wager %d, hands won %d, hands finished %d, total hands %d, has played %t, is winning current hand %t: %s", s.Player, *s.Wager, *s.HandsWon, handsFinished, totalHands, hasPlayedForCurrentHand, isWinningCurrentHand, mood.JSONString())
			s.Mood = mood
		}
	}

	return state, status, cards
}

func lostMood(miss int) PlayerMood {
	if miss > 3 || miss < -3 {
		return PlayerMoodLostReallyBadly
	} else if miss > 1 || miss < -1 {
		return PlayerMoodLostBadly
	} else if miss > 0 || miss < 0 {
		return PlayerMoodLost
	}
	panic("can't handle 0s")
}

func playerMood(wager int, handsWon int, handsFinished int, totalHands int, hasPlayedForCurrentHand bool, isWinningCurrentHand bool) PlayerMood {
	diff := wager - handsWon
	// game over?
	if handsFinished == totalHands {
		if diff == 0 {
			if wager == 0 {
				return PlayerMoodPotato
			}
			return PlayerMoodWon
		}
		return lostMood(diff)
	}

	handsRemaining := totalHands - handsFinished - 1
	maxWins := handsWon + handsRemaining
	if !hasPlayedForCurrentHand || (hasPlayedForCurrentHand && isWinningCurrentHand) {
		maxWins++
	}
	minWins := handsWon
	log.Debugf("max, min wins: %d, %d", maxWins, minWins)
	if minWins == wager && isWinningCurrentHand {
		return PlayerMoodScared
	}
	if maxWins == wager || minWins == wager {
		return PlayerMoodBarelyWinnable
	}
	if maxWins >= wager && minWins <= wager {
		return PlayerMoodWinnable
	}
	if maxWins < wager {
		return lostMood(maxWins - wager)
	}
	if minWins > wager {
		return lostMood(minWins - wager)
	}
	panic("not sure how this happened")
}

type PlayerMood int

const (
	PlayerMoodNone            PlayerMood = iota
	PlayerMoodLost            PlayerMood = iota
	PlayerMoodLostBadly       PlayerMood = iota
	PlayerMoodLostReallyBadly PlayerMood = iota
	PlayerMoodScared          PlayerMood = iota
	PlayerMoodWinnable        PlayerMood = iota
	PlayerMoodBarelyWinnable  PlayerMood = iota
	PlayerMoodPotato          PlayerMood = iota
	PlayerMoodWon             PlayerMood = iota
)

func (p PlayerMood) JSONString() string {
	switch p {
	case PlayerMoodNone:
		return "None"
	case PlayerMoodLost:
		return "Lost"
	case PlayerMoodLostBadly:
		return "LostBadly"
	case PlayerMoodLostReallyBadly:
		return "LostReallyBadly"
	case PlayerMoodScared:
		return "Scared"
	case PlayerMoodWinnable:
		return "Winnable"
	case PlayerMoodBarelyWinnable:
		return "BarelyWinnable"
	case PlayerMoodPotato:
		return "Potato"
	case PlayerMoodWon:
		return "Won"
	}
	panic(fmt.Errorf("invalid PlayerMood value: %d", p))
}

func (p PlayerMood) MarshalJSON() ([]byte, error) {
	jsonString := fmt.Sprintf(`"%s"`, p.JSONString())
	return []byte(jsonString), nil
}

func (p PlayerMood) MarshalText() (text []byte, err error) {
	return []byte(p.JSONString()), nil
}

func parsePlayerMood(text string) (PlayerMood, error) {
	switch text {
	case "None":
		return PlayerMoodNone, nil
	case "Lost":
		return PlayerMoodLost, nil
	case "LostBadly":
		return PlayerMoodLostBadly, nil
	case "LostReallyBadly":
		return PlayerMoodLostReallyBadly, nil
	case "Scared":
		return PlayerMoodScared, nil
	case "Winnable":
		return PlayerMoodWinnable, nil
	case "BarelyWinnable":
		return PlayerMoodBarelyWinnable, nil
	case "Potato":
		return PlayerMoodPotato, nil
	case "Won":
		return PlayerMoodWon, nil
	}
	return PlayerMoodLost, errors.New(fmt.Sprintf("unable to parse player state %s", text))
}

func (p *PlayerMood) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	status, err := parsePlayerMood(str)
	if err != nil {
		return err
	}
	*p = status
	return nil
}

func (p *PlayerMood) UnmarshalText(text []byte) (err error) {
	status, err := parsePlayerMood(string(text))
	if err != nil {
		return err
	}
	*p = status
	return nil
}

package game

import (
	"errors"
	"fmt"
)

type RoundState int

const (
	RoundStateNothingDoneYet RoundState = iota
	RoundStateCardsDealt     RoundState = iota
	RoundStateWagersMade     RoundState = iota
	RoundStateHandInProgress RoundState = iota
	RoundStateFinished       RoundState = iota
)

func (r RoundState) String() string {
	switch r {
	case RoundStateNothingDoneYet:
		return "RoundStateNothingDoneYet"
	case RoundStateCardsDealt:
		return "RoundStateCardsDealt"
	case RoundStateWagersMade:
		return "RoundStateWagersMade"
	case RoundStateHandInProgress:
		return "RoundStateHandInProgress"
	case RoundStateFinished:
		return "RoundStateFinished"
	}
	panic(fmt.Errorf("invalid RoundState value: %d", r))
}

type PlayerCard struct {
	Card     *Card
	IsPlayed bool
}

type Round struct {
	CardsPerPlayer int
	Deck           Deck
	// Players are ordered
	PlayersOrder []string
	Players      map[string]map[string]*PlayerCard
	TrumpSuit    string
	Wagers       map[string]int
	WagerSum     int
	Hands        []*Hand
	//
	State RoundState
}

func NewRound(players []string, deck Deck, cardsPerPlayer int) *Round {
	playersMap := map[string]map[string]*PlayerCard{}
	for _, player := range players {
		playersMap[player] = map[string]*PlayerCard{}
	}
	return &Round{
		CardsPerPlayer: cardsPerPlayer,
		Deck:           deck,
		PlayersOrder:   players,
		Players:        playersMap,
		TrumpSuit:      "",
		Wagers:         map[string]int{},
		WagerSum:       0,
		Hands:          []*Hand{},
		State:          RoundStateNothingDoneYet,
	}
}

func (round *Round) Deal() error {
	if round.State != RoundStateNothingDoneYet {
		return errors.New(fmt.Sprintf("expected state RoundStateNothingDoneYet for deal, found %s", round.State.String()))
	}
	cards := round.Deck.Shuffle()
	j := 0
	for i := 0; i < round.CardsPerPlayer; i++ {
		for _, player := range round.PlayersOrder {
			round.Players[player][cards[j].Key()] = &PlayerCard{Card: cards[j], IsPlayed: false}
			j++
		}
	}
	round.TrumpSuit = cards[j].Suit
	round.State = RoundStateCardsDealt
	return nil
}

func (round *Round) Wager(player string, hands int) error {
	if round.State != RoundStateCardsDealt {
		return errors.New(fmt.Sprintf("expected state RoundStateCardsDealt for wager, found %s", round.State.String()))
	}
	if hands > round.CardsPerPlayer {
		return errors.New(fmt.Sprintf("%d cards per player, but wager was %d", round.CardsPerPlayer, hands))
	}
	// players must make wagers in order
	nextPlayer := round.PlayersOrder[len(round.Wagers)]
	if nextPlayer != player {
		return errors.New(fmt.Sprintf("it is player %s's turn to wager, but got %s", nextPlayer, player))
	}
	playerCount, wagerCount := len(round.PlayersOrder), len(round.Wagers)
	// on the last (i.e. dealer) wager?
	if playerCount == wagerCount+1 {
		// then can't add up to the number of cards
		if hands+round.WagerSum == round.CardsPerPlayer {
			// TODO distinguish between violations of game rules (like this) and something else unexpected going
			// wrong -- like above, where a player has already made a wager or where a player is unrecognized
			return errors.New(fmt.Sprintf("dealer's wager can't add up to %d (had %d already, wagered %d)", round.CardsPerPlayer, round.WagerSum, hands))
		}
		round.State = RoundStateWagersMade
	}
	round.Wagers[player] = hands
	round.WagerSum += hands
	return nil
}

func (round *Round) StartHand() error {
	if round.State != RoundStateWagersMade {
		return errors.New(fmt.Sprintf("expected state RoundStateWagersMade for starting a hand, found %s", round.State.String()))
	}
	round.State = RoundStateHandInProgress
	var players []string
	if len(round.Hands) == 0 {
		// first hand?  start with the first player
		players = append([]string{}, round.PlayersOrder...)
	} else {
		// not the first hand?  start with the previous winner, otherwise continue in the same order
		prevHand := round.Hands[len(round.Hands)-1]
		var i int
		var player string
		for i, player = range round.PlayersOrder {
			if player == prevHand.Leader {
				break
			}
		}
		for j := 0; j < len(round.PlayersOrder); i, j = i+1, j+1 {
			players = append(players, round.PlayersOrder[i%len(round.PlayersOrder)])
		}
	}
	round.Hands = append(round.Hands, NewHand(round.Deck, round.TrumpSuit, players))
	return nil
}

func (round *Round) CurrentHand() (*Hand, error) {
	if round.State != RoundStateHandInProgress {
		return nil, errors.New(fmt.Sprintf("expected state RoundStateHandInProgress, found %s", round.State.String()))
	}
	return round.Hands[len(round.Hands)-1], nil
}

func (round *Round) playerHasCard(player string, card *Card) bool {
	_, ok := round.Players[player][card.Key()]
	return ok
}

func (round *Round) PlayCard(player string, card *Card) error {
	hand, err := round.CurrentHand()
	if err != nil {
		return err
	}
	// is this the right next player?
	nextPlayer := hand.PlayersOrder[len(hand.CardsPlayed)]
	if nextPlayer != player {
		return errors.New(fmt.Sprintf("expected player %s, got %s", nextPlayer, player))
	}
	// is this a card they have?
	if !round.playerHasCard(player, card) {
		return errors.New(fmt.Sprintf("player %s can't play card %+v: does not have it", player, card))
	}
	// have they already played this card?
	if round.Players[player][card.Key()].IsPlayed {
		return errors.New(fmt.Sprintf("player %s can't play card %+v: already played", player, card))
	}
	//is this a card they can legally play?
	if len(hand.CardsPlayed) > 0 {
		// must follow suit if possible, otherwise anything goes
		mustFollowSuit := false
		for _, pc := range round.Players[player] {
			if !pc.IsPlayed && pc.Card.Suit == hand.Suit {
				mustFollowSuit = true
				break
			}
		}
		if mustFollowSuit && card.Suit != hand.Suit {
			return errors.New(fmt.Sprintf("player %s must follow suit %s, but did not", player, hand.Suit))
		}
	}
	hand.PlayCard(player, card)
	round.Players[player][card.Key()].IsPlayed = true

	// have we finished the hand?
	if len(hand.CardsPlayed) == len(round.PlayersOrder) {
		// have we finished the round?
		if len(round.Hands) == round.CardsPerPlayer {
			round.State = RoundStateFinished
		} else {
			round.State = RoundStateWagersMade
		}
	}

	return nil
}

type Hand struct {
	Deck         Deck
	TrumpSuit    string
	CardsPlayed  map[string]*Card
	PlayersOrder []string
	Suit         string
	Leader       string
	LeaderCard   *Card
}

func NewHand(deck Deck, trumpSuit string, playersOrder []string) *Hand {
	return &Hand{
		Deck:         deck,
		TrumpSuit:    trumpSuit,
		CardsPlayed:  map[string]*Card{},
		PlayersOrder: playersOrder,
		Suit:         "",
		Leader:       "",
		LeaderCard:   nil,
	}
}

func (hand *Hand) PlayCard(player string, card *Card) {
	cardsPlayed := len(hand.CardsPlayed)
	hand.CardsPlayed[player] = card
	if cardsPlayed == 0 {
		hand.Suit = card.Suit
		hand.Leader = player
		hand.LeaderCard = card
	} else {
		// which suit is better?  trump > following suit > something else
		if card.Suit == hand.TrumpSuit && hand.LeaderCard.Suit == hand.TrumpSuit {
			// 1. both trumps -- use numbers
			if hand.Deck.Compare(hand.LeaderCard.Number, card.Number) < 0 {
				hand.Leader = player
				hand.LeaderCard = card
			}
		} else if card.Suit == hand.TrumpSuit && hand.LeaderCard.Suit != hand.TrumpSuit {
			// 2. new card is a trump, old one isn't
			hand.Leader = player
			hand.LeaderCard = card
		} else if card.Suit != hand.TrumpSuit && hand.LeaderCard.Suit == hand.TrumpSuit {
			// 3. old card is a trump, new one isn't
			// nothing to do
		} else if card.Suit == hand.Suit && hand.LeaderCard.Suit == hand.Suit {
			// 4. both following suit
			if hand.Deck.Compare(hand.LeaderCard.Number, card.Number) < 0 {
				hand.Leader = player
				hand.LeaderCard = card
			}
		} else if card.Suit == hand.Suit && hand.LeaderCard.Suit != hand.Suit {
			// 5. new card follows suit, old one doesn't
			hand.Leader = player
			hand.LeaderCard = card
		} else if card.Suit != hand.Suit && hand.LeaderCard.Suit == hand.Suit {
			// 6. old card follows suit, new one doesn't
			// nothing to do
		} else {
			// 7. new card can't possibly be better
			// nothing to do
		}
	}
}

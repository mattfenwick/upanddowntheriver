package game

import (
	"fmt"
	"github.com/pkg/errors"
)

type RoundState int

const (
	RoundStateWagers         RoundState = iota
	RoundStateHandInProgress RoundState = iota
	RoundStateFinished       RoundState = iota
)

func (r RoundState) String() string {
	switch r {
	case RoundStateWagers:
		return "RoundStateWagers"
	case RoundStateHandInProgress:
		return "RoundStateHandInProgress"
	case RoundStateFinished:
		return "RoundStateFinished"
	}
	panic(fmt.Errorf("invalid RoundState value: %d", r))
}

func (r RoundState) MarshalJSON() ([]byte, error) {
	jsonString := fmt.Sprintf(`"%s"`, r.String())
	return []byte(jsonString), nil
}

func (r RoundState) MarshalText() (text []byte, err error) {
	return []byte(r.String()), nil
}

type PlayerCard struct {
	Card  *Card
	Count int
}

type CardBag struct {
	Cards map[string]*PlayerCard
}

func NewCardBag(cards []*Card) *CardBag {
	cb := &CardBag{Cards: map[string]*PlayerCard{}}
	for _, card := range cards {
		cb.add(card)
	}
	return cb
}

func (cb *CardBag) add(card *Card) {
	key := card.Key()
	if _, ok := cb.Cards[key]; !ok {
		cb.Cards[key] = &PlayerCard{
			Card:  card,
			Count: 0,
		}
	}
	cb.Cards[key].Count++
}

func (cb *CardBag) remove(card *Card) error {
	key := card.Key()
	if _, ok := cb.Cards[key]; !ok {
		return errors.New(fmt.Sprintf("can't remove card %+v, not found", card))
	}
	cb.Cards[key].Count--
	if cb.Cards[key].Count == 0 {
		delete(cb.Cards, key)
	}
	return nil
}

func (cb *CardBag) has(card *Card) bool {
	_, ok := cb.Cards[card.Key()]
	return ok
}

// does return dupes if necessary
func (cb *CardBag) cards() []*Card {
	cards := []*Card{}
	for _, pc := range cb.Cards {
		for i := 0; i < pc.Count; i++ {
			cards = append(cards, pc.Card)
		}
	}
	return cards
}

type Round struct {
	Guid           string
	CardsPerPlayer int
	Deck           Deck
	// Players are ordered
	PlayersOrder  []string
	PlayerCards   map[string]*CardBag
	TrumpSuit     string
	Wagers        map[string]int
	WagerSum      int
	FinishedHands []*Hand
	CurrentHand   *Hand
	//
	State RoundState
}

func NewRound(players []string, deck Deck, cardsPerPlayer int) *Round {
	playerCards := map[string]*CardBag{}
	for _, player := range players {
		playerCards[player] = NewCardBag([]*Card{})
	}
	round := &Round{
		Guid:           NewGuid(),
		CardsPerPlayer: cardsPerPlayer,
		Deck:           deck,
		PlayersOrder:   players,
		PlayerCards:    playerCards,
		TrumpSuit:      "",
		Wagers:         map[string]int{},
		WagerSum:       0,
		FinishedHands:  []*Hand{},
		CurrentHand:    nil,
		State:          RoundStateWagers,
	}
	round.deal()
	return round
}

func (round *Round) deal() {
	cards := Shuffle(round.Deck)
	j := 0
	for i := 0; i < round.CardsPerPlayer; i++ {
		for _, player := range round.PlayersOrder {
			round.PlayerCards[player].add(cards[j])
			j++
		}
	}
	// instead of reserving a card to choose as the trump suit, we'll just randomly pick a suit
	// meaning that every single card could be dealt to players
	// idk, it just seems like this should be fine
	round.TrumpSuit = RandomSuit(round.Deck)
}

func (round *Round) Wager(player string, hands int) error {
	if round.State != RoundStateWagers {
		return errors.New(fmt.Sprintf("expected state RoundStateWagers for wager, found %s", round.State.String()))
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
		round.startHand()
	}
	round.Wagers[player] = hands
	round.WagerSum += hands
	return nil
}

func (round *Round) startHand() {
	round.State = RoundStateHandInProgress
	var players []string
	if len(round.FinishedHands) == 0 {
		// first hand?  start with the first player
		players = append([]string{}, round.PlayersOrder...)
	} else {
		// not the first hand?  start with the previous winner, otherwise continue in the same order
		prevHand := round.FinishedHands[len(round.FinishedHands)-1]
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
	round.CurrentHand = NewHand(round.Deck, round.TrumpSuit, players)
}

func (round *Round) PlayCard(player string, card *Card) error {
	if round.State != RoundStateHandInProgress {
		return errors.New(fmt.Sprintf("expected state RoundStateHandInProgress, found %s", round.State.String()))
	}

	// is this the right next player?
	hand := round.CurrentHand
	nextPlayer := hand.PlayersOrder[len(hand.CardsPlayed)]
	if nextPlayer != player {
		return errors.New(fmt.Sprintf("expected player %s, got %s", nextPlayer, player))
	}
	// is this a card they have?
	if !round.PlayerCards[player].has(card) {
		return errors.New(fmt.Sprintf("player %s can't play card %+v: does not have it", player, card))
	}
	//is this a card they can legally play?
	if len(hand.CardsPlayed) > 0 {
		// must follow suit if possible, otherwise anything goes
		mustFollowSuit := false
		for _, card := range round.PlayerCards[player].cards() {
			if card.Suit == hand.Suit {
				mustFollowSuit = true
				break
			}
		}
		if mustFollowSuit && card.Suit != hand.Suit {
			return errors.New(fmt.Sprintf("player %s must follow suit %s, but did not", player, hand.Suit))
		}
	}
	hand.PlayCard(player, card)
	err := round.PlayerCards[player].remove(card)
	if err != nil {
		return errors.WithMessagef(err, "unable to remove card")
	}

	// have we finished the hand?
	if len(hand.CardsPlayed) == len(round.PlayersOrder) {
		round.finishHand()
	}

	return nil
}

func (round *Round) finishHand() {
	round.FinishedHands = append(round.FinishedHands, round.CurrentHand)
	round.CurrentHand = nil

	// have we finished the round?
	if len(round.FinishedHands) == round.CardsPerPlayer {
		round.State = RoundStateFinished
	} else {
		round.startHand()
	}
}

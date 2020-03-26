package game

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Deck interface {
	Suits() []string
	Numbers() []string
	Compare(l string, r string) int
}

type Card struct {
	Suit   string
	Number string
}

func Cards(deck Deck) []*Card {
	cards := []*Card{}
	for _, suit := range deck.Suits() {
		for _, num := range deck.Numbers() {
			cards = append(cards, &Card{
				Suit:   suit,
				Number: num,
			})
		}
	}
	return cards
}

func Shuffle(cards []*Card) []*Card {
	cardsCopy := make([]*Card, len(cards))
	for i, c := range cards {
		cardsCopy[i] = c
	}
	swap := func(i int, j int) {
		cardsCopy[i], cardsCopy[j] = cardsCopy[j], cardsCopy[i]
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cardsCopy), swap)
	return cardsCopy
}

type StandardDeck struct {
	DeckSuits     []string
	DeckNumbers   []string
	NumberRatings map[string]int
}

func NewStandardDeck() *StandardDeck {
	numbers := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "K", "Q", "A"}
	ratings := map[string]int{}
	for i, num := range numbers {
		ratings[num] = i
	}
	return &StandardDeck{
		DeckSuits:     []string{"Clubs", "Diamonds", "Hearts", "Spades"},
		DeckNumbers:   numbers,
		NumberRatings: ratings,
	}
}

func (sd *StandardDeck) Suits() []string {
	return sd.DeckSuits
}

func (sd *StandardDeck) Numbers() []string {
	return sd.DeckNumbers
}

func (sd *StandardDeck) Compare(l string, r string) int {
	lRating, ok := sd.NumberRatings[l]
	if !ok {
		panic(fmt.Sprintf("invalid value %s, expected one of %+v", l, sd.DeckNumbers))
	}
	rRating, ok := sd.NumberRatings[r]
	if !ok {
		panic(fmt.Sprintf("invalid value %s, expected one of %+v", r, sd.DeckNumbers))
	}
	return lRating - rRating
}

type RoundState int

const (
	RoundStateNothingDoneYet = iota
	RoundStateCardsDealt     = iota
	RoundStateWagersMade     = iota
)

func (r RoundState) String() string {
	switch r {
	case RoundStateNothingDoneYet:
		return "RoundStateNothingDoneYet"
	case RoundStateCardsDealt:
		return "RoundStateCardsDealt"
	case RoundStateWagersMade:
		return "RoundStateWagersMade"
	}
	panic(fmt.Errorf("invalid RoundState value: %d", r))
}

type Round struct {
	// Dealer must be in Players
	Dealer         string
	CardsPerPlayer int
	Deck           Deck
	// Players are ordered
	PlayersOrder []string
	Players      map[string][]*Card
	TrumpSuit    string
	Wagers       map[string]int
	WagerSum     int
	//
	State RoundState
}

func NewRound(dealer string, players []string, deck Deck, cardsPerPlayer int) (*Round, error) {
	playersMap := map[string][]*Card{}
	for _, player := range players {
		if _, ok := playersMap[player]; ok {
			return nil, errors.New(fmt.Sprintf("duplicate player name: %s", player))
		}
		playersMap[player] = []*Card{}
	}
	if _, ok := playersMap[dealer]; !ok {
		return nil, errors.New(fmt.Sprintf("invalid dealer name %s, not found in %+v", dealer, players))
	}
	// 1 = for the trump suit
	cardsNeeded := cardsPerPlayer*len(players) + 1
	cardsAvailable := len(Cards(deck))
	if cardsNeeded > cardsAvailable {
		return nil, errors.New(fmt.Sprintf("need %d cards for %d players, a total of %d -- more than the %d available", cardsPerPlayer, len(players), cardsNeeded, cardsAvailable))
	}
	return &Round{
		Dealer:         dealer,
		CardsPerPlayer: cardsPerPlayer,
		Deck:           deck,
		PlayersOrder:   players,
		Players:        playersMap,
		TrumpSuit:      "",
		Wagers:         map[string]int{},
		WagerSum:       0,
		State:          RoundStateNothingDoneYet,
	}, nil
}

func (round *Round) Deal() error {
	if round.State != RoundStateNothingDoneYet {
		return errors.New(fmt.Sprintf("expected state RoundStateNothingDoneYet for deal, found %s", round.State.String()))
	}
	cards := Shuffle(Cards(round.Deck))
	j := 0
	for i := 0; i < round.CardsPerPlayer; i++ {
		for _, player := range round.PlayersOrder {
			round.Players[player] = append(round.Players[player], cards[j])
			j++
		}
	}
	round.TrumpSuit = cards[j].Suit
	round.State = RoundStateCardsDealt
	return nil
}

func (round *Round) Wager(player string, hands int) error {
	if hands > round.CardsPerPlayer {
		return errors.New(fmt.Sprintf("%d cards per player, but wager was %d", round.CardsPerPlayer, hands))
	}
	if _, ok := round.Players[player]; !ok {
		return errors.New(fmt.Sprintf("unrecognized player name <%s>", player))
	}
	if _, ok := round.Wagers[player]; ok {
		return errors.New(fmt.Sprintf("player <%s> has already made a wager", player))
	}
	playerCount, wagerCount := len(round.PlayersOrder), len(round.Wagers)
	// on the last (i.e. dealer) wager?  then can't add up to the number of cards
	if (playerCount == wagerCount+1) && (hands+round.WagerSum == round.CardsPerPlayer) {
		// TODO distinguish between violations of game rules (like this) and something else unexpected going
		// wrong -- like above, where a player has already made a wager or where a player is unrecognized
		return errors.New(fmt.Sprintf("dealer's wager can't add up to %d (had %d already, wagered %d)", round.CardsPerPlayer, round.WagerSum, hands))
	}
	round.Wagers[player] = hands
	round.WagerSum += hands
	return nil
}

type Hand struct {
	Deck        Deck
	TrumpSuit   string
	CardsPlayed map[string]*Card
	Suit        string
	Leader      string
	LeaderCard  *Card
}

func NewHand(deck Deck, trumpSuit string) *Hand {
	return &Hand{
		Deck:        deck,
		TrumpSuit:   trumpSuit,
		CardsPlayed: map[string]*Card{},
		Suit:        "",
		Leader:      "",
		LeaderCard:  nil,
	}
}

func (hand *Hand) PlayCard(player string, card *Card) error {
	// TODO check to make sure same card hasn't already been played
	// TODO check to make sure same player hasn't already played
	// TODO check to make sure players play in right order
	// TODO need to know all the card's in `player`s hand in order to make sure they followed suit appropriately
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
	return nil
}

package game

import (
	"fmt"
	"math/rand"
	"time"
)

type Deck interface {
	Suits() []string
	Numbers() []string
	Compare(l string, r string) int
	Shuffle() []*Card
}

type Card struct {
	Suit   string
	Number string
}

func (card *Card) Key() string {
	return fmt.Sprintf("%s-%s", card.Suit, card.Number)
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
	numbers := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
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

func (sd *StandardDeck) Shuffle() []*Card {
	return Shuffle(Cards(sd))
}

// for testing purposes:

type DeterministicShuffleDeck struct {
	StandardDeck *StandardDeck
}

func NewDeterministicShuffleDeck() *DeterministicShuffleDeck {
	return &DeterministicShuffleDeck{StandardDeck: NewStandardDeck()}
}

func (dsd *DeterministicShuffleDeck) Suits() []string {
	return dsd.StandardDeck.Suits()
}

func (dsd *DeterministicShuffleDeck) Numbers() []string {
	return dsd.StandardDeck.Numbers()
}

func (dsd *DeterministicShuffleDeck) Compare(l string, r string) int {
	return dsd.StandardDeck.Compare(l, r)
}

func (dsd *DeterministicShuffleDeck) Shuffle() []*Card {
	return Cards(dsd.StandardDeck)
}

package game

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"math/rand"
	"sort"
	"time"
)

type Deck interface {
	Suits() []string
	Numbers() []string
	CompareNumbers(l string, r string) int
	CompareSuits(l string, r string) int
	Compare(l *Card, r *Card) int
	Shuffle() []*Card
	DeckType() DeckType
	Size() int
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

func RandomSuit(deck Deck) string {
	suits := deck.Suits()
	suitsCopy := make([]string, len(suits))
	for i, c := range suits {
		suitsCopy[i] = c
	}
	swap := func(i int, j int) {
		suitsCopy[i], suitsCopy[j] = suitsCopy[j], suitsCopy[i]
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(suitsCopy), swap)
	return suitsCopy[0]
}

type StandardDeck struct {
	DeckSuits     []string
	SuitRatings   map[string]int
	DeckNumbers   []string
	NumberRatings map[string]int
}

func NewStandardDeck() *StandardDeck {
	numbers := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
	ratings := map[string]int{}
	for i, num := range numbers {
		ratings[num] = i
	}
	suits := []string{"Clubs", "Diamonds", "Hearts", "Spades"}
	suitRatings := map[string]int{}
	for i, suit := range suits {
		suitRatings[suit] = i
	}
	return &StandardDeck{
		DeckSuits:     suits,
		SuitRatings:   suitRatings,
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

func (sd *StandardDeck) CompareNumbers(l string, r string) int {
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

func (sd *StandardDeck) CompareSuits(l string, r string) int {
	lRating, ok := sd.SuitRatings[l]
	if !ok {
		panic(fmt.Sprintf("invalid value %s, expected one of %+v", l, sd.DeckNumbers))
	}
	rRating, ok := sd.SuitRatings[r]
	if !ok {
		panic(fmt.Sprintf("invalid value %s, expected one of %+v", r, sd.DeckNumbers))
	}
	return lRating - rRating
}

func (sd *StandardDeck) Compare(l *Card, r *Card) int {
	num := sd.CompareNumbers(l.Number, r.Number)
	if num != 0 {
		return num
	}
	return sd.CompareSuits(l.Suit, r.Suit)
}

func (sd *StandardDeck) Shuffle() []*Card {
	return Shuffle(Cards(sd))
}

func (sd *StandardDeck) DeckType() DeckType {
	return DeckTypeStandard
}

func (sd *StandardDeck) Size() int {
	return len(sd.DeckSuits) * len(sd.DeckNumbers)
}

// for larger games

type DoubleStandardDeck struct {
	StandardDeck *StandardDeck
}

func NewDoubleStandardDeck() *DoubleStandardDeck {
	return &DoubleStandardDeck{
		StandardDeck: NewStandardDeck(),
	}
}

func (dsd *DoubleStandardDeck) Suits() []string {
	return dsd.StandardDeck.DeckSuits
}

func (dsd *DoubleStandardDeck) Numbers() []string {
	nums := append(dsd.StandardDeck.Numbers(), dsd.StandardDeck.Numbers()...)
	sort.Slice(nums, func(i, j int) bool {
		return dsd.StandardDeck.NumberRatings[nums[i]] < dsd.StandardDeck.NumberRatings[nums[j]]
	})
	return nums
}

func (dsd *DoubleStandardDeck) CompareNumbers(l string, r string) int {
	return dsd.StandardDeck.CompareNumbers(l, r)
}

func (dsd *DoubleStandardDeck) CompareSuits(l string, r string) int {
	return dsd.StandardDeck.CompareSuits(l, r)
}

func (dsd *DoubleStandardDeck) Compare(l *Card, r *Card) int {
	return dsd.StandardDeck.Compare(l, r)
}

func (dsd *DoubleStandardDeck) Shuffle() []*Card {
	return Shuffle(Cards(dsd))
}

func (dsd *DoubleStandardDeck) DeckType() DeckType {
	return DeckTypeDoubleStandard
}

func (dsd *DoubleStandardDeck) Size() int {
	return 2 * dsd.StandardDeck.Size()
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

func (dsd *DeterministicShuffleDeck) CompareNumbers(l string, r string) int {
	return dsd.StandardDeck.CompareNumbers(l, r)
}

func (dsd *DeterministicShuffleDeck) CompareSuits(l string, r string) int {
	return dsd.StandardDeck.CompareSuits(l, r)
}

func (dsd *DeterministicShuffleDeck) Compare(l *Card, r *Card) int {
	return dsd.StandardDeck.Compare(l, r)
}

func (dsd *DeterministicShuffleDeck) Shuffle() []*Card {
	return Cards(dsd.StandardDeck)
}

func (dsd *DeterministicShuffleDeck) DeckType() DeckType {
	return DeckTypeDeterministicStandard
}

func (dsd *DeterministicShuffleDeck) Size() int {
	return dsd.StandardDeck.Size()
}

// deck type

type DeckType string

const (
	DeckTypeStandard              DeckType = "DeckTypeStandard"
	DeckTypeDoubleStandard        DeckType = "DeckTypeDoubleStandard"
	DeckTypeDeterministicStandard DeckType = "DeckTypeDeterministicStandard"
)

func (d DeckType) JSONString() string {
	switch d {
	case DeckTypeStandard:
		return "Standard"
	case DeckTypeDoubleStandard:
		return "DoubleStandard"
	case DeckTypeDeterministicStandard:
		return "DeterministicStandard"
	}
	panic(fmt.Errorf("invalid DeckType value: %s", d))
}

func (d DeckType) MarshalJSON() ([]byte, error) {
	jsonString := fmt.Sprintf(`"%s"`, d.JSONString())
	return []byte(jsonString), nil
}

func (d DeckType) MarshalText() (text []byte, err error) {
	return []byte(d.JSONString()), nil
}

func parseDeckType(text string) (DeckType, error) {
	switch text {
	case "Standard":
		return DeckTypeStandard, nil
	case "DoubleStandard":
		return DeckTypeDoubleStandard, nil
	case "DeterministicStandard":
		return DeckTypeDeterministicStandard, nil
	}
	return DeckTypeStandard, errors.New(fmt.Sprintf("unable to parse deck type %s", text))
}

func (d *DeckType) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	status, err := parseDeckType(str)
	if err != nil {
		return err
	}
	*d = status
	return nil
}

func (d *DeckType) UnmarshalText(text []byte) (err error) {
	status, err := parseDeckType(string(text))
	if err != nil {
		return err
	}
	*d = status
	return nil
}

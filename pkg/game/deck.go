package game

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"math/rand"
	"sort"
	"time"
)

type shuffler func([]*Card) []*Card

type Deck interface {
	Suits() []string
	Numbers() []string
	CompareNumbers(l string, r string) int
	CompareSuits(l string, r string) int
	Compare(l *Card, r *Card) int
	Shuffle() shuffler
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

func Shuffle(deck Deck) []*Card {
	return deck.Shuffle()(Cards(deck))
}

func RandomShuffle(cards []*Card) []*Card {
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

func NoShuffle(cards []*Card) []*Card {
	return cards
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

type SimpleDeck struct {
	DeckSuits     []string
	SuitRatings   map[string]int
	DeckNumbers   []string
	NumberRatings map[string]int
	Type          DeckType
	shuffle       shuffler
}

func NewSimpleDeck(numbers []string, suits []string, deckType DeckType, shuffle shuffler) *SimpleDeck {
	ratings := map[string]int{}
	for i, num := range numbers {
		ratings[num] = i
	}
	suitRatings := map[string]int{}
	for i, suit := range suits {
		suitRatings[suit] = i
	}
	return &SimpleDeck{
		DeckSuits:     suits,
		SuitRatings:   suitRatings,
		DeckNumbers:   numbers,
		NumberRatings: ratings,
		Type:          deckType,
		shuffle:       shuffle,
	}
}

func NewMiniDeckWithShuffle(shuffle shuffler) *SimpleDeck {
	numbers := []string{"J", "Q", "K", "A"}
	suits := []string{"Clubs", "Diamonds", "Hearts", "Spades"}
	return NewSimpleDeck(numbers, suits, DeckTypeMini, shuffle)
}

func NewStandardDeckWithShuffle(shuffle shuffler) *SimpleDeck {
	numbers := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
	suits := []string{"Clubs", "Diamonds", "Hearts", "Spades"}
	return NewSimpleDeck(numbers, suits, DeckTypeStandard, shuffle)
}

func NewStandardDeck() *SimpleDeck {
	return NewStandardDeckWithShuffle(RandomShuffle)
}

func NewDeterministicShuffleDeck() *SimpleDeck {
	return NewStandardDeckWithShuffle(NoShuffle)
}

func (sd *SimpleDeck) Suits() []string {
	return sd.DeckSuits
}

func (sd *SimpleDeck) Numbers() []string {
	return sd.DeckNumbers
}

func (sd *SimpleDeck) CompareNumbers(l string, r string) int {
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

func (sd *SimpleDeck) CompareSuits(l string, r string) int {
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

func (sd *SimpleDeck) Compare(l *Card, r *Card) int {
	num := sd.CompareNumbers(l.Number, r.Number)
	if num != 0 {
		return num
	}
	return sd.CompareSuits(l.Suit, r.Suit)
}

func (sd *SimpleDeck) Shuffle() shuffler {
	return sd.shuffle
}

func (sd *SimpleDeck) DeckType() DeckType {
	return sd.Type
}

func (sd *SimpleDeck) Size() int {
	return len(sd.DeckSuits) * len(sd.DeckNumbers)
}

// DoubleDeck can be used for larger games to increase the number of cards.
// It has two cards for each number/suit, and otherwise works identically
// to its underlying deck.
type DoubleDeck struct {
	UnderlyingDeck Deck
	Type           DeckType
}

func NewDoubleDeck(underlyingDeck Deck, deckType DeckType) *DoubleDeck {
	return &DoubleDeck{UnderlyingDeck: underlyingDeck, Type: deckType}
}

func NewDoubleStandardDeck() *DoubleDeck {
	return NewDoubleDeck(NewStandardDeck(), DeckTypeDoubleStandard)
}

func (dd *DoubleDeck) Suits() []string {
	return dd.UnderlyingDeck.Suits()
}

func (dd *DoubleDeck) Numbers() []string {
	nums := append(dd.UnderlyingDeck.Numbers(), dd.UnderlyingDeck.Numbers()...)
	sort.Slice(nums, func(i, j int) bool {
		return dd.CompareNumbers(nums[i], nums[j]) < 0
	})
	return nums
}

func (dd *DoubleDeck) CompareNumbers(l string, r string) int {
	return dd.UnderlyingDeck.CompareNumbers(l, r)
}

func (dd *DoubleDeck) CompareSuits(l string, r string) int {
	return dd.UnderlyingDeck.CompareSuits(l, r)
}

func (dd *DoubleDeck) Compare(l *Card, r *Card) int {
	return dd.UnderlyingDeck.Compare(l, r)
}

func (dd *DoubleDeck) Shuffle() shuffler {
	return dd.UnderlyingDeck.Shuffle()
}

func (dd *DoubleDeck) DeckType() DeckType {
	return dd.Type
}

func (dd *DoubleDeck) Size() int {
	return 2 * dd.UnderlyingDeck.Size()
}

// deck type

type DeckType string

const (
	// use Custom if your deck isn't one of the predefined types
	DeckTypeCustom                DeckType = "DeckTypeCustom"
	DeckTypeMini                  DeckType = "DeckTypeMini"
	DeckTypeDoubleMini            DeckType = "DeckTypeDoubleMini"
	DeckTypeStandard              DeckType = "DeckTypeStandard"
	DeckTypeDoubleStandard        DeckType = "DeckTypeDoubleStandard"
	DeckTypeDeterministicStandard DeckType = "DeckTypeDeterministicStandard"
)

func (d DeckType) JSONString() string {
	switch d {
	case DeckTypeCustom:
		return "DeckTypeCustom"
	case DeckTypeMini:
		return "DeckTypeMini"
	case DeckTypeDoubleMini:
		return "DeckTypeDoubleMini"
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
	case "Custom":
		return DeckTypeCustom, nil
	case "Mini":
		return DeckTypeMini, nil
	case "DoubleMini":
		return DeckTypeDoubleMini, nil
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

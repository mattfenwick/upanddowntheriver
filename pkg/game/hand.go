package game

type Hand struct {
	Guid         string
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
		Guid:         NewGuid(),
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
			if hand.Deck.CompareNumbers(hand.LeaderCard.Number, card.Number) < 0 {
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
			if hand.Deck.CompareNumbers(hand.LeaderCard.Number, card.Number) < 0 {
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

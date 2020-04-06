package game

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func RunHandTests() {
	//twoOfClubs := &Card{Suit: "Clubs", Number: "2"}
	threeOfClubs := &Card{Suit: "Clubs", Number: "3"}
	//fourOfClubs := &Card{Suit: "Clubs", Number: "4"}
	//fiveOfClubs := &Card{Suit: "Clubs", Number: "5"}
	//sixOfClubs := &Card{Suit: "Clubs", Number: "6"}
	//sevenOfClubs := &Card{Suit: "Clubs", Number: "7"}
	jackOfClubs := &Card{Suit: "Clubs", Number: "J"}
	//sixOfSpades := &Card{Suit: "Spades", Number: "6"}
	//aceOfSpades := &Card{Suit: "Spades", Number: "A"}
	//twoOfDiamonds := &Card{Suit: "Diamonds", Number: "2"}
	//threeOfDiamonds := &Card{Suit: "Diamonds", Number: "3"}
	//fourOfDiamonds := &Card{Suit: "Diamonds", Number: "4"}
	//fiveOfDiamonds := &Card{Suit: "Diamonds", Number: "5"}
	nineOfDiamonds := &Card{Suit: "Diamonds", Number: "9"}
	kingOfHearts := &Card{Suit: "Hearts", Number: "K"}

	Describe("Hand", func() {
		deck := NewDeterministicShuffleDeck()
		players := []string{"ned", "homer", "karina"}

		Describe("initialization", func() {
			It("should have the right state", func() {
				hand := NewHand(deck, "Diamonds", players)
				Expect(hand.TrumpSuit).To(Equal("Diamonds"))
			})

			It("should have other fields empty", func() {
				hand := NewHand(deck, "Clubs", players)
				Expect(hand.Suit).To(Equal(""))
				Expect(hand.CardsPlayed).To(Equal(map[string]*Card{}))
				Expect(hand.Leader).To(Equal(""))
				Expect(hand.LeaderCard).To(BeNil())
			})
		})

		Describe("Play", func() {
			It("Should have the first card determine the suit", func() {
				hand := NewHand(deck, "Clubs", players)

				hand.PlayCard("abc", nineOfDiamonds)

				Expect(hand.Leader).To(Equal("abc"))
				Expect(hand.LeaderCard).To(Equal(nineOfDiamonds))
				Expect(hand.CardsPlayed).To(Equal(map[string]*Card{
					"abc": nineOfDiamonds,
				}))
				Expect(hand.Suit).To(Equal("Diamonds"))
			})

			It("Should treat a trump card as better", func() {
				hand := NewHand(deck, "Clubs", players)

				hand.PlayCard("abc", nineOfDiamonds)
				hand.PlayCard("def", threeOfClubs)

				Expect(hand.Leader).To(Equal("def"))
				Expect(hand.LeaderCard).To(Equal(threeOfClubs))
				Expect(hand.CardsPlayed).To(Equal(map[string]*Card{
					"abc": nineOfDiamonds,
					"def": threeOfClubs,
				}))
				Expect(hand.Suit).To(Equal("Diamonds"))
			})

			It("Should treat a higher card in the same suit as better", func() {
				hand := NewHand(deck, "Diamonds", players)

				hand.PlayCard("abc", threeOfClubs)
				hand.PlayCard("ghi", jackOfClubs)

				Expect(hand.Leader).To(Equal("ghi"))
				Expect(hand.LeaderCard).To(Equal(jackOfClubs))
				Expect(hand.CardsPlayed).To(Equal(map[string]*Card{
					"abc": threeOfClubs,
					"ghi": jackOfClubs,
				}))
				Expect(hand.Suit).To(Equal("Clubs"))
			})

			It("Should treat a following-suit card as better", func() {
				hand := NewHand(deck, "Diamonds", players)

				hand.PlayCard("abc", threeOfClubs)
				hand.PlayCard("def", kingOfHearts)

				Expect(hand.Leader).To(Equal("abc"))
				Expect(hand.LeaderCard).To(Equal(threeOfClubs))
				Expect(hand.CardsPlayed).To(Equal(map[string]*Card{
					"abc": threeOfClubs,
					"def": kingOfHearts,
				}))
				Expect(hand.Suit).To(Equal("Clubs"))
			})
		})
	})
}

package game

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func RunRoundTests() {
	twoOfClubs := &Card{Suit: "Clubs", Number: "2"}
	threeOfClubs := &Card{Suit: "Clubs", Number: "3"}
	fourOfClubs := &Card{Suit: "Clubs", Number: "4"}
	fiveOfClubs := &Card{Suit: "Clubs", Number: "5"}
	sixOfClubs := &Card{Suit: "Clubs", Number: "6"}
	sevenOfClubs := &Card{Suit: "Clubs", Number: "7"}
	//jackOfClubs := &Card{Suit: "Clubs", Number: "J"}
	//sixOfSpades := &Card{Suit: "Spades", Number: "6"}
	//aceOfSpades := &Card{Suit: "Spades", Number: "A"}
	twoOfDiamonds := &Card{Suit: "Diamonds", Number: "2"}
	threeOfDiamonds := &Card{Suit: "Diamonds", Number: "3"}
	fourOfDiamonds := &Card{Suit: "Diamonds", Number: "4"}
	fiveOfDiamonds := &Card{Suit: "Diamonds", Number: "5"}
	//nineOfDiamonds := &Card{Suit: "Diamonds", Number: "9"}
	//kingOfHearts := &Card{Suit: "Hearts", Number: "K"}

	Describe("Round", func() {
		players := []string{"player1", "jimbo", "alfonso"}
		deck := NewDeterministicShuffleDeck()

		Describe("initialization", func() {
			It("should have the right state and cards per player", func() {
				round := NewRound(players, deck, 3)
				Expect(round.State).To(Equal(RoundStateWagers))
				Expect(round.CardsPerPlayer).To(Equal(3))
			})

			It("should have the right wagersum, wagers, cards, and hands", func() {
				round := NewRound(players, deck, 3)
				Expect(round.WagerSum).To(Equal(0))
				Expect(round.PlayersOrder).To(Equal(players))
				Expect(round.Wagers).To(Equal(map[string]int{}))
				for _, player := range round.PlayersOrder {
					Expect(len(round.PlayerCards[player].cards())).To(Equal(3))
				}
				Expect(round.FinishedHands).To(Equal([]*Hand{}))
			})
		})

		Describe("Deal", func() {
			It("should deal the right number of cards", func() {
				round := NewRound(players, deck, 5)

				Expect(len(round.PlayerCards)).To(Equal(3))
				Expect(round.State).To(Equal(RoundStateWagers))

				for _, cardBag := range round.PlayerCards {
					Expect(len(cardBag.cards())).To(Equal(5))
				}
			})

			It("should not deal the same card multiple times from a deck with no duplicates", func() {
				round := NewRound(players, deck, 17)

				keys := map[string]bool{}
				for _, cardBag := range round.PlayerCards {
					for _, card := range cardBag.cards() {
						key := card.Key()
						Expect(keys[key]).To(BeFalse())
						keys[key] = true
					}
				}
				Expect(len(keys)).To(Equal(51))
			})
		})

		Describe("Make wagers", func() {
			It("Requires wagers to be made in the right order", func() {
				round := NewRound(players, deck, 3)

				// "jimbo" is in position 2 -- not good
				Expect(round.Wager("jimbo", 2)).ToNot(BeNil())

				// "player1" is in position 1 -- okay
				Expect(round.Wager("player1", 2)).Should(Succeed())
			})

			It("Doesn't let the dealer (last player) make a wager for a total equal to the number of cards per player", func() {
				round := NewRound(players, deck, 13)

				Expect(round.Wager("player1", 7)).Should(Succeed())
				Expect(len(round.Wagers)).To(Equal(1))
				Expect(round.WagerSum).To(Equal(7))

				Expect(round.Wager("jimbo", 2)).Should(Succeed())
				Expect(len(round.Wagers)).To(Equal(2))
				Expect(round.WagerSum).To(Equal(9))

				Expect(round.Wager("alfonso", 4)).ToNot(BeNil())
				Expect(len(round.Wagers)).To(Equal(2))
				Expect(round.WagerSum).To(Equal(9))
				Expect(round.State).To(Equal(RoundStateWagers))

				Expect(round.Wager("alfonso", 2)).Should(Succeed())
				Expect(len(round.Wagers)).To(Equal(3))
				Expect(round.WagerSum).To(Equal(11))
				Expect(round.State).To(Equal(RoundStateHandInProgress))
			})

			It("Doesn't allow wagers higher than the number of cards per player", func() {
				round := NewRound(players, deck, 3)

				Expect(round.Wager("player1", 4)).ToNot(BeNil())
				Expect(len(round.Wagers)).To(Equal(0))
				Expect(round.WagerSum).To(Equal(0))
				Expect(round.State).To(Equal(RoundStateWagers))

				Expect(round.Wager("player1", 2)).Should(Succeed())
				Expect(len(round.Wagers)).To(Equal(1))
				Expect(round.WagerSum).To(Equal(2))
				Expect(round.State).To(Equal(RoundStateWagers))
			})

			It("Finishes wagers after all players have wagered in turn", func() {
				round := NewRound(players, deck, 5)

				Expect(round.Wager("player1", 1)).Should(Succeed())
				Expect(round.Wager("jimbo", 2)).Should(Succeed())
				Expect(round.Wager("alfonso", 3)).Should(Succeed())

				Expect(len(round.Wagers)).To(Equal(3))
				Expect(round.WagerSum).To(Equal(6))
				Expect(round.State).To(Equal(RoundStateHandInProgress))
			})
		})

		wageredRound := func() *Round {
			round := NewRound(players, deck, 17)

			round.Wager("player1", 1)
			round.Wager("jimbo", 2)
			round.Wager("alfonso", 3)

			return round
		}

		smallWageredRound := func() *Round {
			round := NewRound(players, deck, 2)

			round.Wager("player1", 2)
			round.Wager("jimbo", 1)
			round.Wager("alfonso", 0)

			return round
		}

		Describe("Play", func() {
			It("Has players go in order", func() {
				round := wageredRound()

				Expect(round.PlayCard("jimbo", threeOfClubs)).ToNot(BeNil())

				Expect(round.PlayCard("player1", twoOfClubs)).Should(Succeed())
			})

			It("Should require players to only play cards that they have", func() {
				round := wageredRound()

				Expect(round.PlayCard("player1", threeOfClubs)).ToNot(BeNil())
				Expect(round.PlayCard("player1", twoOfClubs)).Should(Succeed())
			})

			It("Should require players to only play cards that have not already been played", func() {
				round := wageredRound()

				Expect(round.PlayCard("player1", fiveOfClubs)).Should(Succeed())
				Expect(round.PlayCard("jimbo", threeOfClubs)).Should(Succeed())
				Expect(round.PlayCard("alfonso", fourOfClubs)).Should(Succeed())

				Expect(round.State).To(Equal(RoundStateHandInProgress))
				Expect(len(round.FinishedHands)).To(Equal(1))
				Expect(round.CurrentHand).ToNot(BeNil())

				Expect(round.PlayCard("player1", fiveOfClubs)).ToNot(BeNil())
				Expect(round.PlayCard("player1", twoOfClubs)).Should(Succeed())
				Expect(round.State).To(Equal(RoundStateHandInProgress))
				Expect(len(round.FinishedHands)).To(Equal(1))
				Expect(round.CurrentHand).ToNot(BeNil())
			})

			It("Has players follow suit", func() {
				round := wageredRound()

				Expect(round.PlayCard("player1", fiveOfClubs)).Should(Succeed())
				Expect(round.PlayCard("jimbo", twoOfDiamonds)).ToNot(BeNil())
				Expect(round.PlayCard("jimbo", threeOfClubs)).Should(Succeed())
			})

			It("Has the first player start the first hand", func() {
				round := wageredRound()

				Expect(round.PlayCard("jimbo", twoOfDiamonds)).ToNot(BeNil())
				Expect(round.PlayCard("alfonso", threeOfDiamonds)).ToNot(BeNil())

				Expect(round.PlayCard("player1", fiveOfClubs)).Should(Succeed())
				Expect(round.PlayCard("jimbo", threeOfClubs)).Should(Succeed())
				Expect(round.PlayCard("alfonso", fourOfClubs)).Should(Succeed())
			})

			It("Uses the winner of the previous hand as the starter of the next", func() {
				round := wageredRound()

				Expect(round.PlayCard("player1", twoOfClubs)).Should(Succeed())
				Expect(round.PlayCard("jimbo", threeOfClubs)).Should(Succeed())
				Expect(round.PlayCard("alfonso", fourOfClubs)).Should(Succeed())

				Expect(round.State).To(Equal(RoundStateHandInProgress))
				Expect(len(round.FinishedHands)).To(Equal(1))
				Expect(round.CurrentHand).ToNot(BeNil())

				hand := round.CurrentHand
				Expect(hand.PlayersOrder).To(Equal([]string{"alfonso", "player1", "jimbo"}))

				Expect(round.PlayCard("alfonso", threeOfDiamonds)).Should(Succeed())
				Expect(round.PlayCard("player1", fourOfDiamonds)).Should(Succeed())
				Expect(round.PlayCard("jimbo", fiveOfDiamonds)).Should(Succeed())

				hand3 := round.CurrentHand
				Expect(hand3.PlayersOrder).To(Equal([]string{"jimbo", "alfonso", "player1"}))
			})

			It("Should report cards as having been played, as they are played", func() {
				round := wageredRound()

				Expect(round.PlayCard("player1", twoOfClubs)).Should(Succeed())
				Expect(round.PlayCard("jimbo", threeOfClubs)).Should(Succeed())
				Expect(round.PlayCard("alfonso", fourOfClubs)).Should(Succeed())

				Expect(round.PlayerCards["player1"].has(twoOfClubs)).To(BeFalse())
				Expect(round.PlayerCards["player1"].has(fiveOfClubs)).To(BeTrue())
				Expect(round.PlayerCards["jimbo"].has(threeOfClubs)).To(BeFalse())
				Expect(round.PlayerCards["jimbo"].has(sixOfClubs)).To(BeTrue())
				Expect(round.PlayerCards["alfonso"].has(fourOfClubs)).To(BeFalse())
				Expect(round.PlayerCards["alfonso"].has(sevenOfClubs)).To(BeTrue())
			})
		})

		Describe("Finish", func() {
			It("Has a number of hands equal to the number of cards per player", func() {
				round := smallWageredRound()

				Expect(round.PlayCard("player1", twoOfClubs)).Should(Succeed())
				Expect(round.PlayCard("jimbo", threeOfClubs)).Should(Succeed())
				Expect(round.PlayCard("alfonso", fourOfClubs)).Should(Succeed())

				Expect(round.State).To(Equal(RoundStateHandInProgress))

				Expect(round.PlayCard("alfonso", sevenOfClubs)).Should(Succeed())
				Expect(round.PlayCard("player1", fiveOfClubs)).Should(Succeed())
				Expect(round.PlayCard("jimbo", sixOfClubs)).Should(Succeed())

				Expect(round.State).To(Equal(RoundStateFinished))
			})
		})
	})
}

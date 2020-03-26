package game

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func RunGameTests() {
	Describe("Game", func() {
		Describe("Initialization", func() {
			It("should have the right number of cards per player", func() {
				game := NewGame()
				Expect(game.CardsPerPlayer).To(Equal(1))
			})
			It("should have a default deck", func() {
				game := NewGame()
				Expect(game.Deck.Numbers()).To(Equal([]string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}))
				Expect(game.Deck.Suits()).To(Equal([]string{"Clubs", "Diamonds", "Hearts", "Spades"}))
			})
		})

		Describe("Setup", func() {
			It("should add players", func() {
				game := NewGame()
				Expect(game.addPlayer("abc")).To(BeNil())
				Expect(game.addPlayer("def")).To(BeNil())
				Expect(game.Players).To(Equal([]string{"abc", "def"}))
			})

			It("should not add the same player twice", func() {
				game := NewGame()
				Expect(game.addPlayer("abc")).To(BeNil())
				Expect(game.addPlayer("def")).To(BeNil())
				Expect(game.Players).To(Equal([]string{"abc", "def"}))

				Expect(game.addPlayer("def")).ToNot(BeNil())
				Expect(game.addPlayer("abc")).ToNot(BeNil())
				Expect(game.addPlayer("ghi")).To(BeNil())
				Expect(game.Players).To(Equal([]string{"abc", "def", "ghi"}))
			})

			// TODO remove player

			It("should start a round", func() {
				game := NewGame()
				Expect(game.addPlayer("abc")).To(BeNil())
				Expect(game.addPlayer("def")).To(BeNil())
				Expect(game.addPlayer("ghi")).To(BeNil())

				Expect(game.State).To(Equal(GameStateSetup))
				Expect(game.startRound()).To(BeNil())
				Expect(game.State).To(Equal(GameStateRoundInProgress))
			})

			It("shouldn't start a round with fewer than 2 players", func() {
				game := NewGame()
				Expect(game.addPlayer("abc")).To(BeNil())

				Expect(game.State).To(Equal(GameStateSetup))
				Expect(game.startRound()).ToNot(BeNil())
				Expect(game.State).To(Equal(GameStateSetup))
			})
		})
	})

	Describe("Round", func() {
		players := []string{"player1", "jimbo", "alfonso"}
		deck := NewStandardDeck()

		Describe("initialization", func() {
			It("should have the right state", func() {
				round := NewRound(players, deck, 3)
				Expect(round.State).To(Equal(RoundStateNothingDoneYet))
				Expect(round.CardsPerPlayer).To(Equal(3))
			})

			It("should have cards unitialized", func() {
				round := NewRound(players, deck, 3)
				Expect(round.WagerSum).To(Equal(0))
				Expect(round.PlayersOrder).To(Equal(players))
				Expect(round.Players).To(Equal(map[string]map[string]*PlayerCard{
					"player1": {},
					"jimbo":   {},
					"alfonso": {},
				}))
				Expect(round.Hands).To(Equal([]*Hand{}))
			})
		})

		Describe("Deal", func() {

		})

		Describe("Make wagers", func() {

		})

		Describe("Play", func() {

		})

		Describe("Finish", func() {

		})
	})

	Describe("Hand", func() {
		deck := NewStandardDeck()

		threeOfClubs := &Card{Suit: "Clubs", Number: "3"}
		jackOfClubs := &Card{Suit: "Clubs", Number: "J"}
		//sixOfSpades := &Card{Suit: "Spades", Number: "6"}
		//aceOfSpades := &Card{Suit: "Spades", Number: "A"}
		nineOfDiamonds := &Card{Suit: "Diamonds", Number: "9"}
		kingOfHearts := &Card{Suit: "Hearts", Number: "K"}

		Describe("initialization", func() {
			It("should have the right state", func() {
				hand := NewHand(deck, "Diamonds")
				Expect(hand.TrumpSuit).To(Equal("Diamonds"))
			})

			It("should have other fields empty", func() {
				hand := NewHand(deck, "Clubs")
				Expect(hand.Suit).To(Equal(""))
				Expect(hand.CardsPlayed).To(Equal(map[string]*Card{}))
				Expect(hand.Leader).To(Equal(""))
				Expect(hand.LeaderCard).To(BeNil())
			})
		})

		Describe("Play", func() {
			It("Should have the first card determine the suit", func() {
				hand := NewHand(deck, "Clubs")

				hand.PlayCard("abc", nineOfDiamonds)

				Expect(hand.Leader).To(Equal("abc"))
				Expect(hand.LeaderCard).To(Equal(nineOfDiamonds))
				Expect(hand.CardsPlayed).To(Equal(map[string]*Card{
					"abc": nineOfDiamonds,
				}))
				Expect(hand.Suit).To(Equal("Diamonds"))
			})

			It("Should treat a trump card as better", func() {
				hand := NewHand(deck, "Clubs")

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
				hand := NewHand(deck, "Diamonds")

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
				hand := NewHand(deck, "Diamonds")

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

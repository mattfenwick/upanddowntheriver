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
				Expect(game.addPlayer("abc")).Should(Succeed())
				Expect(game.addPlayer("def")).Should(Succeed())
				Expect(game.Players).To(Equal([]string{"abc", "def"}))
			})

			It("should not add the same player twice", func() {
				game := NewGame()
				Expect(game.addPlayer("abc")).Should(Succeed())
				Expect(game.addPlayer("def")).Should(Succeed())
				Expect(game.Players).To(Equal([]string{"abc", "def"}))

				Expect(game.addPlayer("def")).ToNot(BeNil())
				Expect(game.addPlayer("abc")).ToNot(BeNil())
				Expect(game.addPlayer("ghi")).Should(Succeed())
				Expect(game.Players).To(Equal([]string{"abc", "def", "ghi"}))
			})

			It("should handle setCardsPerPlayer to max", func() {
				game := NewGame()
				Expect(game.addPlayer("abc")).Should(Succeed())
				Expect(game.addPlayer("def")).Should(Succeed())
				Expect(game.addPlayer("ghi")).Should(Succeed())
				Expect(game.addPlayer("jkl")).Should(Succeed())
				Expect(game.Players).To(Equal([]string{"abc", "def", "ghi", "jkl"}))

				Expect(game.setCardsPerPlayer(13)).To(Succeed())
				Expect(game.startRound()).To(Succeed())
			})

			// TODO remove player

			It("should start a round", func() {
				game := NewGame()
				Expect(game.addPlayer("abc")).Should(Succeed())
				Expect(game.addPlayer("def")).Should(Succeed())
				Expect(game.addPlayer("ghi")).Should(Succeed())

				Expect(game.State).To(Equal(GameStateSetup))
				Expect(game.startRound()).Should(Succeed())
				Expect(game.State).To(Equal(GameStateRoundInProgress))
			})

			It("shouldn't start a round with fewer than 2 players", func() {
				game := NewGame()
				Expect(game.addPlayer("abc")).Should(Succeed())

				Expect(game.State).To(Equal(GameStateSetup))
				Expect(game.startRound()).ToNot(BeNil())
				Expect(game.State).To(Equal(GameStateSetup))
			})
		})

		getFirstCard := func(cards map[string]*PlayerCard) *Card {
			for _, c := range cards {
				return c.Card
			}
			panic("no cards found")
		}

		Describe("Round", func() {
			It("should rotate players after a round", func() {
				game := NewGame()
				Expect(game.addPlayer("abc")).Should(Succeed())
				Expect(game.addPlayer("def")).Should(Succeed())
				Expect(game.addPlayer("ghi")).Should(Succeed())
				Expect(game.setCardsPerPlayer(1)).Should(Succeed())

				Expect(game.Players).To(Equal([]string{"abc", "def", "ghi"}))

				Expect(game.startRound()).Should(Succeed())

				Expect(game.makeWager("abc", 0)).Should(Succeed())
				Expect(game.makeWager("def", 0)).Should(Succeed())
				Expect(game.makeWager("ghi", 0)).Should(Succeed())

				Expect(game.playCard("abc", getFirstCard(game.CurrentRound.Players["abc"]))).Should(Succeed())
				Expect(game.playCard("def", getFirstCard(game.CurrentRound.Players["def"]))).Should(Succeed())
				Expect(game.playCard("ghi", getFirstCard(game.CurrentRound.Players["ghi"]))).Should(Succeed())

				Expect(game.finishRound()).Should(Succeed())

				Expect(game.Players).To(Equal([]string{"def", "ghi", "abc"}))
			})
		})
	})
}

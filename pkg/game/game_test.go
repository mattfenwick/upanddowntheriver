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
	})
}

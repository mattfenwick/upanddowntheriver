package game

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

func joinGame(game *Game, player string) error {
	addedPlayer, err := game.join(player)
	if err != nil {
		return err
	}
	if addedPlayer != player {
		return errors.New(fmt.Sprintf("attempted to add player %s, got player %s added instead", player, addedPlayer))
	}
	return nil
}

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
				Expect(joinGame(game, "abc")).Should(Succeed())
				Expect(joinGame(game, "def")).Should(Succeed())
				Expect(game.Players).To(Equal([]string{"abc", "def"}))
			})

			It("should reject empty string for player name", func() {
				game := NewGame()
				Expect(joinGame(game, "abc")).Should(Succeed())
				Expect(joinGame(game, "")).ShouldNot(Succeed())
				Expect(joinGame(game, "def")).Should(Succeed())
				Expect(game.Players).To(Equal([]string{"abc", "def"}))
			})

			It("should not add the same player twice; however, doing so is not an error", func() {
				game := NewGame()
				Expect(joinGame(game, "abc")).Should(Succeed())
				Expect(joinGame(game, "def")).Should(Succeed())
				Expect(game.Players).To(Equal([]string{"abc", "def"}))

				Expect(joinGame(game, "def")).Should(Succeed())
				Expect(joinGame(game, "abc")).Should(Succeed())
				Expect(joinGame(game, "ghi")).Should(Succeed())
				Expect(game.Players).To(Equal([]string{"abc", "def", "ghi"}))
			})

			It("should truncate names longer than 20 characters", func() {
				game := NewGame()
				Expect(joinGame(game, "abc")).Should(Succeed())
				abcs := "abcdefghijklmnopqrstuvwxyz"
				err := joinGame(game, abcs)
				Expect(err).ShouldNot(Succeed())
				msg := "attempted to add player abcdefghijklmnopqrstuvwxyz, got player abcdefghijklmnopqrst added instead"
				Expect(err.Error()).Should(Equal(msg))

				Expect(game.Players).To(Equal([]string{"abc", "abcdefghijklmnopqrst"}))
			})

			It("should handle setCardsPerPlayer to max", func() {
				game := NewGame()
				Expect(joinGame(game, "abc")).Should(Succeed())
				Expect(joinGame(game, "def")).Should(Succeed())
				Expect(joinGame(game, "ghi")).Should(Succeed())
				Expect(joinGame(game, "jkl")).Should(Succeed())
				Expect(game.Players).To(Equal([]string{"abc", "def", "ghi", "jkl"}))

				Expect(game.setCardsPerPlayer(13)).To(Succeed())
				Expect(game.startRound()).To(Succeed())
			})

			It("should handle removing players that exist, and fail to remove players that don't exist", func() {
				game := NewGame()

				Expect(joinGame(game, "abc")).Should(Succeed())
				Expect(joinGame(game, "def")).Should(Succeed())
				Expect(joinGame(game, "ghi")).Should(Succeed())

				Expect(game.Players).To(Equal([]string{"abc", "def", "ghi"}))

				Expect(game.removePlayer("jkl")).ShouldNot(Succeed())
				Expect(game.removePlayer("def")).Should(Succeed())
				Expect(game.removePlayer("def")).ShouldNot(Succeed())

				Expect(game.Players).To(Equal([]string{"abc", "ghi"}))
			})

			emptyPm := &PlayerModel{
				Me:    "",
				State: PlayerStateNotJoined,
				Game: &PlayerGame{
					Players:           []string{"abc", "def", "ghi"},
					MaxCardsPerPlayer: 17,
					CardsPerPlayer:    1,
					DeckType:          DeckTypeStandard,
				},
			}

			It("should return an 'empty' model for an 'empty' player, no matter whether the game's in progress", func() {
				game := NewGame()
				Expect(joinGame(game, "abc")).Should(Succeed())
				Expect(joinGame(game, "def")).Should(Succeed())
				Expect(joinGame(game, "ghi")).Should(Succeed())

				pm := game.playerModel("")
				Expect(pm).To(Equal(emptyPm))

				Expect(game.startRound()).Should(Succeed())

				pm2 := game.playerModel("")
				Expect(pm2).To(Equal(emptyPm))
			})

			It("should return an 'empty' player model for a nonexisting player", func() {
				game := NewGame()
				Expect(joinGame(game, "abc")).Should(Succeed())
				Expect(joinGame(game, "def")).Should(Succeed())
				Expect(joinGame(game, "ghi")).Should(Succeed())

				pm := game.playerModel("jkl")
				Expect(pm).To(Equal(emptyPm))
			})

			It("should start a round", func() {
				game := NewGame()
				Expect(joinGame(game, "abc")).Should(Succeed())
				Expect(joinGame(game, "def")).Should(Succeed())
				Expect(joinGame(game, "ghi")).Should(Succeed())

				Expect(game.State).To(Equal(GameStateSetup))
				Expect(game.startRound()).Should(Succeed())
				Expect(game.State).To(Equal(GameStateRoundInProgress))
			})

			It("shouldn't start a round with fewer than 2 players", func() {
				game := NewGame()
				Expect(joinGame(game, "abc")).Should(Succeed())

				Expect(game.State).To(Equal(GameStateSetup))
				Expect(game.startRound()).ToNot(BeNil())
				Expect(game.State).To(Equal(GameStateSetup))
			})

			It("should calculate the right maxCardsPerPlayer number for different deck sizes", func() {
				game := NewGame()
				game.Deck = NewDoubleStandardDeck()
				Expect(joinGame(game, "abc")).Should(Succeed())
				Expect(joinGame(game, "def")).Should(Succeed())
				Expect(joinGame(game, "ghi")).Should(Succeed())

				Expect(game.playerModel("abc").Game.MaxCardsPerPlayer).To(Equal(34))
			})
		})

		getFirstCard := func(cardBag *CardBag) *Card {
			for _, c := range cardBag.cards() {
				return c
			}
			panic("no cards found")
		}

		Describe("Round", func() {
			It("should rotate players after a round", func() {
				game := NewGame()
				Expect(joinGame(game, "abc")).Should(Succeed())
				Expect(joinGame(game, "def")).Should(Succeed())
				Expect(joinGame(game, "ghi")).Should(Succeed())
				Expect(game.setCardsPerPlayer(1)).Should(Succeed())

				Expect(game.Players).To(Equal([]string{"abc", "def", "ghi"}))

				Expect(game.startRound()).Should(Succeed())

				Expect(game.makeWager("abc", 0)).Should(Succeed())
				Expect(game.makeWager("def", 0)).Should(Succeed())
				Expect(game.makeWager("ghi", 0)).Should(Succeed())

				Expect(game.playCard("abc", getFirstCard(game.CurrentRound.PlayerCards["abc"]))).Should(Succeed())
				Expect(game.playCard("def", getFirstCard(game.CurrentRound.PlayerCards["def"]))).Should(Succeed())
				Expect(game.playCard("ghi", getFirstCard(game.CurrentRound.PlayerCards["ghi"]))).Should(Succeed())

				Expect(game.finishRound()).Should(Succeed())

				Expect(game.Players).To(Equal([]string{"def", "ghi", "abc"}))
			})

			It("should calculate player moods according to their chances of winning or how bad they lost", func() {
				game := NewGame()
				game.Deck = NewDeterministicShuffleDeck()
				Expect(joinGame(game, "abc")).Should(Succeed())
				Expect(joinGame(game, "def")).Should(Succeed())
				Expect(joinGame(game, "ghi")).Should(Succeed())
				Expect(game.setCardsPerPlayer(4)).Should(Succeed())

				Expect(game.startRound()).Should(Succeed())

				Expect(game.makeWager("abc", 4)).Should(Succeed())
				Expect(game.makeWager("def", 2)).Should(Succeed())
				Expect(game.makeWager("ghi", 0)).Should(Succeed())

				pm1 := game.playerModel("abc").Status.PlayerStatuses
				Expect(pm1[0].Mood).To(Equal(PlayerMoodBarelyWinnable))
				Expect(pm1[1].Mood).To(Equal(PlayerMoodWinnable))
				Expect(pm1[2].Mood).To(Equal(PlayerMoodBarelyWinnable))

				// TODO track moods through more hands, to the end of the round
			})
		})
	})
}

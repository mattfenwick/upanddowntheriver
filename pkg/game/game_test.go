package game

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func RunGameTests() {
	twoOfClubs := &Card{Suit: "Clubs", Number: "2"}
	threeOfClubs := &Card{Suit: "Clubs", Number: "3"}
	fourOfClubs := &Card{Suit: "Clubs", Number: "4"}
	fiveOfClubs := &Card{Suit: "Clubs", Number: "5"}
	sixOfClubs := &Card{Suit: "Clubs", Number: "6"}
	sevenOfClubs := &Card{Suit: "Clubs", Number: "7"}
	jackOfClubs := &Card{Suit: "Clubs", Number: "J"}
	//sixOfSpades := &Card{Suit: "Spades", Number: "6"}
	//aceOfSpades := &Card{Suit: "Spades", Number: "A"}
	twoOfDiamonds := &Card{Suit: "Diamonds", Number: "2"}
	threeOfDiamonds := &Card{Suit: "Diamonds", Number: "3"}
	fourOfDiamonds := &Card{Suit: "Diamonds", Number: "4"}
	fiveOfDiamonds := &Card{Suit: "Diamonds", Number: "5"}
	nineOfDiamonds := &Card{Suit: "Diamonds", Number: "9"}
	kingOfHearts := &Card{Suit: "Hearts", Number: "K"}

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

	Describe("Round", func() {
		players := []string{"player1", "jimbo", "alfonso"}
		deck := NewDeterministicShuffleDeck()

		Describe("initialization", func() {
			It("should have the right state and cards per player", func() {
				round := NewRound(players, deck, 3)
				Expect(round.State).To(Equal(RoundStateCardsDealt))
				Expect(round.CardsPerPlayer).To(Equal(3))
			})

			It("should have the right wagersum, wagers, cards, and hands", func() {
				round := NewRound(players, deck, 3)
				Expect(round.WagerSum).To(Equal(0))
				Expect(round.PlayersOrder).To(Equal(players))
				Expect(round.Wagers).To(Equal(map[string]int{}))
				for _, player := range round.PlayersOrder {
					Expect(len(round.Players[player])).To(Equal(3))
				}
				Expect(round.Hands).To(Equal([]*Hand{}))
			})
		})

		Describe("Deal", func() {
			It("should deal the right number of cards", func() {
				round := NewRound(players, deck, 5)

				Expect(len(round.Players)).To(Equal(3))
				Expect(round.State).To(Equal(RoundStateCardsDealt))

				for _, cards := range round.Players {
					Expect(len(cards)).To(Equal(5))
				}
			})

			It("should not deal the same card multiple times", func() {
				round := NewRound(players, deck, 17)

				keys := map[string]bool{}
				for _, cards := range round.Players {
					for key, _ := range cards {
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
				Expect(round.State).To(Equal(RoundStateCardsDealt))

				Expect(round.Wager("alfonso", 2)).Should(Succeed())
				Expect(len(round.Wagers)).To(Equal(3))
				Expect(round.WagerSum).To(Equal(11))
				Expect(round.State).To(Equal(RoundStateWagersMade))
			})

			It("Doesn't allow wagers higher than the number of cards per player", func() {
				round := NewRound(players, deck, 3)

				Expect(round.Wager("player1", 4)).ToNot(BeNil())
				Expect(len(round.Wagers)).To(Equal(0))
				Expect(round.WagerSum).To(Equal(0))
				Expect(round.State).To(Equal(RoundStateCardsDealt))

				Expect(round.Wager("player1", 2)).Should(Succeed())
				Expect(len(round.Wagers)).To(Equal(1))
				Expect(round.WagerSum).To(Equal(2))
				Expect(round.State).To(Equal(RoundStateCardsDealt))
			})

			It("Finishes wagers after all players have wagered in turn", func() {
				round := NewRound(players, deck, 5)

				Expect(round.Wager("player1", 1)).Should(Succeed())
				Expect(round.Wager("jimbo", 2)).Should(Succeed())
				Expect(round.Wager("alfonso", 3)).Should(Succeed())

				Expect(len(round.Wagers)).To(Equal(3))
				Expect(round.WagerSum).To(Equal(6))
				Expect(round.State).To(Equal(RoundStateWagersMade))
			})
		})

		wageredRound := func() *Round {
			round := NewRound(players, deck, 17)

			round.Wager("player1", 1)
			round.Wager("jimbo", 2)
			round.Wager("alfonso", 3)

			Expect(round.StartHand()).Should(Succeed())

			return round
		}

		smallWageredRound := func() *Round {
			round := NewRound(players, deck, 2)

			round.Wager("player1", 2)
			round.Wager("jimbo", 1)
			round.Wager("alfonso", 0)

			Expect(round.StartHand()).Should(Succeed())

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

				Expect(round.State).To(Equal(RoundStateWagersMade))
				Expect(len(round.Hands)).To(Equal(1))

				Expect(round.StartHand()).Should(Succeed())
				Expect(round.PlayCard("player1", fiveOfClubs)).ToNot(BeNil())
				Expect(round.PlayCard("player1", twoOfClubs)).Should(Succeed())
				Expect(round.State).To(Equal(RoundStateHandInProgress))
				Expect(len(round.Hands)).To(Equal(2))
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

				Expect(round.State).To(Equal(RoundStateWagersMade))
				Expect(len(round.Hands)).To(Equal(1))

				Expect(round.StartHand()).To(Succeed())

				hand, err := round.CurrentHand()
				Expect(err).To(Succeed())
				Expect(hand.PlayersOrder).To(Equal([]string{"alfonso", "player1", "jimbo"}))

				Expect(round.PlayCard("alfonso", threeOfDiamonds)).Should(Succeed())
				Expect(round.PlayCard("player1", fourOfDiamonds)).Should(Succeed())
				Expect(round.PlayCard("jimbo", fiveOfDiamonds)).Should(Succeed())

				Expect(round.StartHand()).To(Succeed())

				hand3, err := round.CurrentHand()
				Expect(err).To(Succeed())
				Expect(hand3.PlayersOrder).To(Equal([]string{"jimbo", "alfonso", "player1"}))
			})

			It("Should report cards as having been played, as they are played", func() {
				round := wageredRound()

				Expect(round.PlayCard("player1", twoOfClubs)).Should(Succeed())
				Expect(round.PlayCard("jimbo", threeOfClubs)).Should(Succeed())
				Expect(round.PlayCard("alfonso", fourOfClubs)).Should(Succeed())

				Expect(round.Players["player1"]["Clubs-2"].IsPlayed).To(BeTrue())
				Expect(round.Players["player1"]["Clubs-5"].IsPlayed).To(BeFalse())
				Expect(round.Players["jimbo"]["Clubs-3"].IsPlayed).To(BeTrue())
				Expect(round.Players["jimbo"]["Clubs-6"].IsPlayed).To(BeFalse())
				Expect(round.Players["alfonso"]["Clubs-4"].IsPlayed).To(BeTrue())
				Expect(round.Players["alfonso"]["Clubs-7"].IsPlayed).To(BeFalse())
			})
		})

		Describe("Finish", func() {
			It("Has a number of hands equal to the number of cards per player", func() {
				round := smallWageredRound()

				Expect(round.PlayCard("player1", twoOfClubs)).Should(Succeed())
				Expect(round.PlayCard("jimbo", threeOfClubs)).Should(Succeed())
				Expect(round.PlayCard("alfonso", fourOfClubs)).Should(Succeed())

				Expect(round.State).To(Equal(RoundStateWagersMade))
				Expect(round.StartHand()).Should(Succeed())

				Expect(round.State).To(Equal(RoundStateHandInProgress))

				Expect(round.PlayCard("alfonso", sevenOfClubs)).Should(Succeed())
				Expect(round.PlayCard("player1", fiveOfClubs)).Should(Succeed())
				Expect(round.PlayCard("jimbo", sixOfClubs)).Should(Succeed())

				Expect(round.State).To(Equal(RoundStateFinished))
			})
		})
	})

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

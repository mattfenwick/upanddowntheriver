package game

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func RunPlayerModelTests() {
	Describe("PlayerModel", func() {
		It("should calculate won/lost when a round is over", func() {
			Expect(playerMood(3, 3, 6, 6, false, false)).To(Equal(PlayerMoodWon))
			Expect(playerMood(0, 0, 6, 6, false, false)).To(Equal(PlayerMoodWon))

			Expect(playerMood(1, 0, 6, 6, false, false)).To(Equal(PlayerMoodLost))
			Expect(playerMood(2, 0, 6, 6, false, false)).To(Equal(PlayerMoodLostBadly))
			Expect(playerMood(3, 0, 6, 6, false, false)).To(Equal(PlayerMoodLostBadly))
			Expect(playerMood(4, 0, 6, 6, false, false)).To(Equal(PlayerMoodLostReallyBadly))
			Expect(playerMood(5, 0, 6, 6, false, false)).To(Equal(PlayerMoodLostReallyBadly))

			Expect(playerMood(2, 3, 9, 9, false, false)).To(Equal(PlayerMoodLost))
			Expect(playerMood(2, 4, 9, 9, false, false)).To(Equal(PlayerMoodLostBadly))
			Expect(playerMood(2, 5, 9, 9, false, false)).To(Equal(PlayerMoodLostBadly))
			Expect(playerMood(2, 6, 9, 9, false, false)).To(Equal(PlayerMoodLostReallyBadly))
			Expect(playerMood(2, 7, 9, 9, false, false)).To(Equal(PlayerMoodLostReallyBadly))
		})

		It("should calculate lost for in-progress round if hands won > wager", func() {
			Expect(playerMood(2, 3, 7, 8, false, false)).To(Equal(PlayerMoodLost))
			Expect(playerMood(2, 4, 7, 8, false, false)).To(Equal(PlayerMoodLostBadly))
			Expect(playerMood(2, 5, 7, 8, false, false)).To(Equal(PlayerMoodLostBadly))
			Expect(playerMood(2, 6, 7, 8, false, false)).To(Equal(PlayerMoodLostReallyBadly))
			Expect(playerMood(2, 7, 7, 8, false, false)).To(Equal(PlayerMoodLostReallyBadly))

			Expect(playerMood(2, 3, 7, 11, false, false)).To(Equal(PlayerMoodLost))
			Expect(playerMood(2, 4, 7, 11, false, false)).To(Equal(PlayerMoodLostBadly))
			Expect(playerMood(2, 5, 7, 11, false, false)).To(Equal(PlayerMoodLostBadly))
			Expect(playerMood(2, 6, 7, 11, false, false)).To(Equal(PlayerMoodLostReallyBadly))
			Expect(playerMood(2, 7, 7, 11, false, false)).To(Equal(PlayerMoodLostReallyBadly))
		})

		It("should calculate lost for in-progress round if (hands won + hands remaining < wager) ", func() {
			Expect(playerMood(8, 4, 4, 8, true, false)).To(Equal(PlayerMoodLost))
			Expect(playerMood(8, 3, 4, 8, true, false)).To(Equal(PlayerMoodLostBadly))
			Expect(playerMood(8, 2, 4, 8, true, false)).To(Equal(PlayerMoodLostBadly))
			Expect(playerMood(8, 1, 4, 8, true, false)).To(Equal(PlayerMoodLostReallyBadly))
			Expect(playerMood(8, 0, 4, 8, true, false)).To(Equal(PlayerMoodLostReallyBadly))

			Expect(playerMood(8, 4, 4, 8, true, true)).To(Equal(PlayerMoodBarelyWinnable))
			Expect(playerMood(8, 3, 4, 8, true, true)).To(Equal(PlayerMoodLost))
			Expect(playerMood(8, 2, 4, 8, true, true)).To(Equal(PlayerMoodLostBadly))
			Expect(playerMood(8, 1, 4, 8, true, true)).To(Equal(PlayerMoodLostBadly))
			Expect(playerMood(8, 0, 4, 8, true, true)).To(Equal(PlayerMoodLostReallyBadly))

			Expect(playerMood(4, 0, 0, 4, true, false)).To(Equal(PlayerMoodLost))
			Expect(playerMood(4, 0, 1, 4, false, false)).To(Equal(PlayerMoodLost))
		})

		// barely winnable: every remaining trick must go a specific way
		//  - eg: wager 0, won 0, 5 hands remaining
		//  - eg: wager 3, won 0, 3 hands remaining
		//  - eg: wager 3, won 3, 3 hands remaining
		It("should calculate barely winnable if every remaining hand must go a specific way for the player to win", func() {
			Expect(playerMood(0, 0, 0, 5, false, false)).To(Equal(PlayerMoodBarelyWinnable))
			Expect(playerMood(3, 0, 0, 4, true, false)).To(Equal(PlayerMoodBarelyWinnable))
			Expect(playerMood(3, 0, 0, 3, true, true)).To(Equal(PlayerMoodBarelyWinnable))
			Expect(playerMood(4, 0, 0, 4, false, false)).To(Equal(PlayerMoodBarelyWinnable))
			Expect(playerMood(0, 0, 0, 2, false, false)).To(Equal(PlayerMoodBarelyWinnable))
			Expect(playerMood(6, 2, 2, 7, true, false)).To(Equal(PlayerMoodBarelyWinnable))
			Expect(playerMood(6, 2, 2, 6, true, true)).To(Equal(PlayerMoodBarelyWinnable))
		})

		// winnable: there's some slack, not every single trick has to go a certain way
		//  - eg: wager 2, won 0, 4 hands remaining
		//  - eg: wager 2, won 1, 2 hands remaining
		It("should calculate winnable, there's some slack for the player to win", func() {
			Expect(playerMood(2, 0, 0, 4, false, false)).To(Equal(PlayerMoodWinnable))
			Expect(playerMood(2, 1, 2, 4, false, false)).To(Equal(PlayerMoodWinnable))
		})

		It("should take into account whether the player has played a card in the current hand, and is winning the current hand", func() {
			Expect(playerMood(4, 2, 2, 4, false, false)).To(Equal(PlayerMoodBarelyWinnable))
			Expect(playerMood(4, 2, 2, 4, true, true)).To(Equal(PlayerMoodBarelyWinnable))
			Expect(playerMood(4, 2, 2, 4, true, false)).To(Equal(PlayerMoodLost))

			// this is an interesting one: need to win 0 more hands, but leading current one.  what to do?
			Expect(playerMood(0, 0, 2, 4, true, true)).To(Equal(PlayerMoodScared))
		})

	})
}

package game

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func RunDeckTests() {
	Describe("Deck", func() {
		It("should order cards correctly", func() {
			deck := NewStandardDeck()
			Expect(deck.CompareNumbers("A", "Q") > 0).To(BeTrue())
			Expect(deck.CompareNumbers("Q", "A") < 0).To(BeTrue())
		})
	})
}

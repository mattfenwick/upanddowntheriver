package game

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func RunPlayerStateTests() {
	Describe("PlayerState", func() {
		It("should serialize to json", func() {
			val := []PlayerState{PlayerStateHandFinished, PlayerStateHandPlayTurn}
			bytes, err := json.Marshal(val)
			Expect(err).Should(Succeed())
			result := []PlayerState{}
			err = json.Unmarshal(bytes, &result)
			Expect(err).Should(Succeed())
			Expect(result).To(Equal(val))
		})
	})
}

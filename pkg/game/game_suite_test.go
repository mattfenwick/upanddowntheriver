package game

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestModel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunGameTests()
	RunRoundTests()
	RunDeckTests()
	RunHandTests()
	RunPlayerStateTests()
	RunSpecs(t, "game suite")
}

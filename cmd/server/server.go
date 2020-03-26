package main

import (
	"github.com/mattfenwick/upanddowntheriver/pkg/game"
	"os"
)

func main() {
	game.Run(os.Args[1])
}

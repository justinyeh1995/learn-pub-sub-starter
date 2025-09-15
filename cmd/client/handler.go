package main

import (
	"fmt"

	gamelogic "github.com/justinyeh1995/learn-pub-sub-starter/internal/gamelogic"
	routing "github.com/justinyeh1995/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(ps routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
	}
}

package main

import (
	gamelogic "github.com/justinyeh1995/learn-pub-sub-starter/internal/gamelogic"
	routing "github.com/justinyeh1995/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState)

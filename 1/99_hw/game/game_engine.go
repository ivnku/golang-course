package main

import "strings"

type GameEngine struct {
	user           User
	world          map[string]Room
	commandHandler CommandHandler
}

func NewGameEngine() *GameEngine {
	ge := new(GameEngine)
	ge.user = User{}
	ge.world = make(map[string]Room)
	ge.commandHandler = NewCommandHandler()
	return ge
}

func (ge *GameEngine) AddRoom(room Room) {
	ge.world[room.name] = room
}

func (ge *GameEngine) AddCommand(name string, foo func(ge *GameEngine, args ...string) string) {
	ge.commandHandler.commands[name] = foo
}

func (ge *GameEngine) HandleCommand(command string) string {
	splitted := strings.Split(command, " ")
	if len(splitted) == 0 {
		splitted[0] = ""
	}
	return ge.commandHandler.basicAct(splitted[0], ge, splitted[1:]...)
}

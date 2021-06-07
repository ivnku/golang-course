package main

import "strings"

type GameEngine struct {
	player         *Player
	world          map[string]*Room
	commandHandler CommandHandler
}

func NewGameEngine() *GameEngine {
	ge := new(GameEngine)
	ge.player = &Player{}
	ge.world = make(map[string]*Room)
	ge.commandHandler = NewCommandHandler()
	return ge
}

func (ge *GameEngine) addRoom(room *Room) {
	ge.world[room.name] = room
}

func (ge *GameEngine) addCommand(name string, foo func(ge *GameEngine, args ...string) string) {
	ge.commandHandler.commands[name] = foo
}

func (ge *GameEngine) addPlayer(player *Player) {
	ge.player = player
}

func (ge *GameEngine) HandleCommand(command string) string {
	splitted := strings.Split(command, " ")
	if len(splitted) == 0 {
		splitted[0] = ""
	}
	return ge.commandHandler.basicAct(splitted[0], ge, splitted[1:]...)
}

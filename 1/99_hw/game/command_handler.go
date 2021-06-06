package main

import (
	"os"
	"strings"
)

type CommandHandler struct {
	commands map[string]func(ge *GameEngine, args ...string) string
}

func NewCommandHandler() CommandHandler {
	commands := make(map[string]func(ge *GameEngine, args ...string) string)

	commands["осмотреться"] = func(ge *GameEngine, args ...string) string {
		return ge.player.currentRoom.getLookAroundText()
	}

	commands["идти"] = func(ge *GameEngine, args ...string) string {
		if (len(args) == 0) {
			return "такого пути нет"
		}
		targetRoom := args[0]
		if room, isExist := ge.world[targetRoom]; isExist && containString(ge.player.currentRoom.routes, targetRoom) {
			// check if player can go to the room, if not - return message of restriction
			for restrictiontext, checkFunc := range ge.player.currentRoom.conditions {
				if !checkFunc(ge.player) {
					return restrictiontext
				}
			}
			ge.player.currentRoom = room
			return room.getEntryText()
		}
	
		return "невозможно пройти в " + targetRoom
	}

	commands["применить"] = func(ge *GameEngine, args ...string) string {
		return "Primenit'"
	}

	commands["взять"] = func(ge *GameEngine, args ...string) string {
		return "Vzyat'"
	}

	commands["надеть"] = func(ge *GameEngine, args ...string) string {
		return "Nadet'"
	}
	
	commands["гдея"] = func(ge *GameEngine, args ...string) string {
		return ge.player.currentRoom.getEntryText()
	}

	commands["default"] = func(ge *GameEngine, args ...string) string {
		return "неизвестная команда"
	}

	return CommandHandler{
		commands: commands,
	}

}

func (ch *CommandHandler) basicAct(command string, ge *GameEngine, args ...string) string {
	command = strings.ToLower(command)
	_, isExist := ch.commands[command]
	if isExist {
		return ch.commands[command](ge, args...)
	} else if (command == "стоп" || command == "") {
		os.Exit(0)
		return ""
	} else {
		return ch.commands["default"](ge, args...)
	}
}

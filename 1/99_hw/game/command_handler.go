package main

import (
	"strings"
	"os"
)

type CommandHandler struct {
	commands map[string]func(ge *GameEngine, args ...string) string
}

func NewCommandHandler() CommandHandler {
	commands := make(map[string]func(ge *GameEngine, args ...string) string)

	commands["осмотреться"] = func(ge *GameEngine, args ...string) string {
		return ge.user.currentRoom.getLookAroundText()
	}

	commands["идти"] = func(ge *GameEngine, args ...string) string {
		return "Idti"
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

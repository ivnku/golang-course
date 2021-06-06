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
		lookAroundText := ge.player.currentRoom.getLookAroundText()
		if (lookAroundText != "") {
			return lookAroundText
		} else {
			return ge.player.currentRoom.getEntryText()
		}
	}

	commands["идти"] = func(ge *GameEngine, args ...string) string {
		if len(args) == 0 {
			return "путь не указан"
		}
		targetRoom := args[0]
		isRoutePossible, _ := containString(ge.player.currentRoom.routes, targetRoom)
		if room, isExist := ge.world[targetRoom]; isExist && isRoutePossible {
			// check if player can go to the room, if not - return message of restriction
			for restrictiontext, checkFunc := range room.conditions {
				if !checkFunc(ge) {
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
		if len(args) == 0 {
			return "предмет не указан"
		}
		itemName := args[0]
		if item, isEmpty := ge.player.currentRoom.getItem(itemName); !isEmpty {
			if ge.player.inventoryStorage.name != "" {
				ge.player.addInventoryItem(item.name, item)
				ge.player.currentRoom.removeItem(item.name)
				return "предмет добавлен в инвентарь: " + item.name
			} else {
				return "некуда класть"
			}
		}
		return "нет такого"
	}

	commands["надеть"] = func(ge *GameEngine, args ...string) string {
		if len(args) == 0 {
			return "предмет не указан"
		}
		itemName := args[0]
		if item, isEmpty := ge.player.currentRoom.getItem(itemName); !isEmpty {
			if item.isWearable {
				if item.isStorage {
					ge.player.inventoryStorage = item
				}
				ge.player.currentRoom.removeItem(item.name)
				return "вы надели: " + item.name
			} else {
				return "невозможно надеть " + item.name
			}
		}
		return "нет такого"
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
	} else if command == "стоп" || command == "" {
		os.Exit(0)
		return ""
	} else {
		return ch.commands["default"](ge, args...)
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

var gameEngine = NewGameEngine()

func main() {
	initGame()
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		fmt.Println(handleCommand(input.Text()))
	}
}

func initGame() {
	kitchen := Room{
		name:           "кухня",
		entryText:      "кухня, ничего интересного. :routes",
		lookAroundText: "ты находишься на кухне, :items, надо собрать рюкзак и идти в универ. :routes",
	}
	kitchen.addItem("на столе", Item{name: "чай"})
	kitchen.addRoutes([]string{"коридор"})

	corridor := Room{
		name:      "коридор",
		entryText: "ничего интересного. :routes",
	}
	corridor.addRoutes([]string{"комната", "кухня", "улица"})

	room := Room{
		name:           "комната",
		entryText:      "ты в своей комнате. :routes",
		lookAroundText: ":items. :routes",
	}
	room.addItem("на столе", Item{name: "конспекты"})
	room.addItem("на столе", Item{name: "ключи"})
	room.addItem("на стуле", Item{name: "рюкзак", isStorage: true, isWearable: true})
	room.addRoutes([]string{"коридор"})

	street := Room{
		name:      "улица",
		entryText: "на улице весна. :routes",
	}
	street.addRoutes([]string{"домой"})
	street.addCondition("дверь закрыта", func(ge *GameEngine) bool {
		return ge.player.hasInventoryItem("ключ")
	})

	home := Room{name: "домой", entryText: "дом милый дом. :routes"}
	home.addRoutes([]string{"коридор"})

	gameEngine.addRoom(kitchen)
	gameEngine.addRoom(corridor)
	gameEngine.addRoom(room)
	gameEngine.addRoom(street)
	gameEngine.addRoom(home)

	player := Player{currentRoom: kitchen, inventory: make(map[string]Item)}
	gameEngine.addPlayer(player)
}

func handleCommand(command string) string {
	return gameEngine.HandleCommand(command)
}

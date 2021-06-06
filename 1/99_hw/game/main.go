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
		name:           "Кухня",
		lookAroundText: "ты находишься на кухне, :items, надо собрать рюкзак и идти в универ. :routes",
	}
	kitchen.addItem("на столе", Item{name: "чай"})
	kitchen.addRoutes([]string{"коридор"})

	corridor := Room{
		name:      "Коридор",
		entryText: "ничего интересного. :routes",
	}
	corridor.addRoutes([]string{"комната", "кухня", "улица"})

	room := Room{
		name:           "Комната",
		entryText:      "ты в своей комнате. можно пройти - коридор",
		lookAroundText: "на столе: ключи, конспекты, на стуле: рюкзак. можно пройти - коридор",
	}
	room.addItem("на столе", Item{name: "конспекты"})
	room.addItem("на столе", Item{name: "ключи"})
	room.addItem("на стуле", Item{name: "рюкзак", isStorage: true})
	room.addRoutes([]string{"коридор"})

	street := Room{
		name:      "Улица",
		entryText: "на улице весна. :routes",
	}
	street.addRoutes([]string{"домой"})

	home := Room{name: "домой", entryText: "дом милый дом. :routes"}
	home.addRoutes([]string{"коридор"})

	gameEngine.addRoom(kitchen)
	gameEngine.addRoom(corridor)
	gameEngine.addRoom(room)
	gameEngine.addRoom(street)

	player := Player{currentRoom: kitchen, inventory: make(map[string]Item)}
	gameEngine.addPlayer(player)
}

func handleCommand(command string) string {
	return gameEngine.HandleCommand(command)
}

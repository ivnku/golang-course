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
	kitchen.AddItem("на столе", Item{name: "чай"})
	kitchen.AddRoutes([]string{"коридор"})

	corridor := Room{
		name:      "Коридор",
		entryText: "ничего интересного. :routes",
	}
	corridor.AddRoutes([]string{"комната", "кухня", "улица"})

	room := Room{
		name:           "Комната",
		entryText:      "ты в своей комнате. можно пройти - коридор",
		lookAroundText: "на столе: ключи, конспекты, на стуле: рюкзак. можно пройти - коридор",
	}
	room.AddItem("на столе", Item{name: "конспекты"})
	room.AddItem("на столе", Item{name: "ключи"})
	room.AddItem("на стуле", Item{name: "рюкзак", isStorage: true})
	room.AddRoutes([]string{"коридор"})

	street := Room{
		name:      "Улица",
		entryText: "на улице весна. :routes",
	}
	street.AddRoutes([]string{"домой"})

	home := Room{name: "домой", entryText: "дом милый дом. :routes"}
	home.AddRoutes([]string{"коридор"})

	gameEngine.AddRoom(kitchen)
	gameEngine.AddRoom(corridor)
	gameEngine.AddRoom(room)
	gameEngine.AddRoom(street)

	user := User{currentRoom: kitchen, inventory: make(map[string]Item)}
	gameEngine.AddUser(user)
}

func handleCommand(command string) string {
	return gameEngine.HandleCommand(command)
}

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
	kitchen := Room{name: "кухня", entryText: "кухня, ничего интересного. :routes", lookAroundText: "ты находишься на кухне, :items, надо :goals. :routes"}
	kitchen.addItem("на столе", Item{name: "чай"})
	kitchen.addRoutes([]string{"коридор"})

	corridor := Room{name: "коридор", entryText: "ничего интересного. :routes"}
	corridor.addRoutes([]string{"кухня", "комната", "улица"})
	corridor.addCondition("дверь", &Condition{
		state: false,
		check: func(player Player, targetRoom Room, condition *Condition) (bool, string) {
			if targetRoom.name == "улица" && condition.state != true {
				return false, "дверь закрыта"
			}
			return true, ""
		},
	})

	room := Room{name: "комната", entryText: "ты в своей комнате. :routes", lookAroundText: ":items. :routes"}
	room.addItem("на столе", Item{name: "ключи", affectOn: map[string]string{"дверь": "дверь открыта"}})
	room.addItem("на столе", Item{name: "конспекты"})
	room.addItem("на стуле", Item{name: "рюкзак", isStorage: true, isWearable: true})
	room.addRoutes([]string{"коридор"})

	street := Room{name: "улица", entryText: "на улице весна. :routes"}
	street.addRoutes([]string{"домой"})

	home := Room{name: "домой", entryText: "дом милый дом. :routes"}
	home.addRoutes([]string{"коридор"})

	gameEngine.addRoom(kitchen)
	gameEngine.addRoom(corridor)
	gameEngine.addRoom(room)
	gameEngine.addRoom(street)
	gameEngine.addRoom(home)

	player := Player{currentRoom: kitchen, inventory: make(map[string]Item)}
	player.goals = []Goal{
		{text: "собрать рюкзак", check: func(player *Player) bool {
			return player.inventoryStorage.name == "" && len(player.inventory) == 0
		}},
		{text: "идти в универ", check: func(player *Player) bool {
			return true
		}},
	}
	gameEngine.addPlayer(player)
}

func handleCommand(command string) string {
	return gameEngine.HandleCommand(command)
}

package main

import (
	"strings"
)

type Room struct {
	name           string
	entryText      string // text which displayed when user entered a room
	lookAroundText string // basic formatted text for lookAround command
	itemsList      map[string][]Item
	itemsPlaces    []string // necessary for saving the order of printing due to unordered nature of map
	routes         []string
	conditions     map[string]*Condition
}

type Condition struct {
	state bool
	check func(player Player, targetRoom Room, condition *Condition) (bool, string)
}

func (condition *Condition) checkCondition(player Player, targetRoom Room) (bool, string) {
	return condition.check(player, targetRoom, condition)
}

func (room *Room) addItem(place string, item Item) {
	if room.itemsList == nil && room.itemsPlaces == nil {
		room.itemsList = make(map[string][]Item)
		room.itemsPlaces = make([]string, 0)
	}

	room.itemsList[place] = append(room.itemsList[place], item)

	if isExist, _ := containString(room.itemsPlaces, place); !isExist {
		room.itemsPlaces = append(room.itemsPlaces, place)
	}
}

func (room *Room) getItem(name string) (Item, bool) {
	for _, items := range room.itemsList {
		for _, item := range items {
			if item.name == name {
				return item, false
			}
		}
	}
	return Item{}, true
}

func (room *Room) removeItem(name string) {
finish:
	for place, items := range room.itemsList {
		for index, item := range items {
			if item.name != name {
				continue
			}
			ret := make([]Item, 0)
			if len(items) <= 1 { // if this is the last item - remove key from the map and remove itemsPlaces item
				isExist, placeIndex := containString(room.itemsPlaces, place)
				if isExist {
					room.itemsPlaces = removeStringFromSlice(room.itemsPlaces, placeIndex)
					delete(room.itemsList, place)
				}
			} else {
				ret = append(ret, items[:index]...)
				room.itemsList[place] = append(ret, items[index+1:]...)
			}
			break finish
		}
	}
}

func (room *Room) addRoutes(routes []string) {
	for _, route := range routes {
		room.routes = append(room.routes, route)
	}
}

func getFormattedText(str string, player Player) string {
	itemsText := player.currentRoom.getItemsText()
	
	if itemsText == "" {
		itemsText = "пустая комната"
	}

	result := strings.Replace(str, ":items", itemsText, -1)
	result = strings.Replace(result, ":routes", player.currentRoom.getRoutesText(), -1)
	result = strings.Replace(result, ":goals", player.currentRoom.getGoalsText(player), -1)

	return result
}

func (room *Room) getLookAroundText(player Player) string {
	return getFormattedText(room.lookAroundText, player)
}

func (room *Room) getEntryText(player Player) string {
	return getFormattedText(room.entryText, player)
}

func (room *Room) getGoalsText(player Player) string {
	goalsText := []string{}
	for _, goal := range player.goals {
		if goal.check(&player) {
			goalsText = append(goalsText, goal.text)
		}
	}
	return strings.Join(goalsText, " и ")
}

func (room *Room) getItemsText() string {
	itemsText := []string{}
	for _, place := range room.itemsPlaces {
		text := place + ": "
		items := []string{}
		for _, item := range room.itemsList[place] {
			items = append(items, item.name)
		}
		text += strings.Join(items, ", ")
		itemsText = append(itemsText, text)
	}
	return strings.Join(itemsText, ", ")
}

func (room *Room) getRoutesText() string {
	routesText := "можно пройти - "
	routes := []string{}
	for _, route := range room.routes {
		routes = append(routes, route)
	}
	routesText += strings.Join(routes, ", ")
	return routesText
}

func (room *Room) addCondition(name string, condition *Condition) {
	if room.conditions == nil {
		room.conditions = make(map[string]*Condition)
	}
	room.conditions[name] = condition
}

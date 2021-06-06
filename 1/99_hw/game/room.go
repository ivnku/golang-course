package main

import (
	"regexp"
	"strings"
)

type Room struct {
	name           string
	entryText      string // text which displayed when user entered a room
	lookAroundText string // basic formatted text for lookAround command
	itemsList      map[string][]Item
	itemsPlaces    []string
	routes         []string
	conditions     map[string]func(user Player) bool
}

func (room *Room) addItem(place string, item Item) {
	if room.itemsList == nil && room.itemsPlaces == nil {
		room.itemsList = make(map[string][]Item)
		room.itemsPlaces = make([]string, 0)
	}

	room.itemsList[place] = append(room.itemsList[place], item)

	if !containString(room.itemsPlaces, place) {
		room.itemsPlaces = append(room.itemsPlaces, place)
	}
}

func (room *Room) addRoutes(routes []string) {
	for _, route := range routes {
		room.routes = append(room.routes, route)
	}
}

func getFormattedItemsAndRoutes(str string, room *Room) string {
	itemsRegexp := regexp.MustCompile(`:items`)
	routesRegexp := regexp.MustCompile(`:routes`)

	itemsText := room.getItemsText()
	routesText := room.getRoutesText()

	result := itemsRegexp.ReplaceAll([]byte(str), []byte(itemsText))
	result = routesRegexp.ReplaceAll([]byte(result), []byte(routesText))
	return string(result)
}

func (room *Room) getLookAroundText() string {
	return getFormattedItemsAndRoutes(room.lookAroundText, room)
}

func (room *Room) getEntryText() string {
	return getFormattedItemsAndRoutes(room.entryText, room)
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

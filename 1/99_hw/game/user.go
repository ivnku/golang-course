package main

type Player struct {
	currentRoom      Room
	inventoryStorage Item
	inventory        map[string]Item
}

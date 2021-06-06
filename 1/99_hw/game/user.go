package main

type User struct {
	currentRoom      Room
	inventoryStorage Item
	inventory        map[string]Item
}

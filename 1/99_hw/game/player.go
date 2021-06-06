package main

type Player struct {
	currentRoom      Room
	inventoryStorage Item
	inventory        map[string]Item
}

func (player *Player) addInventoryItem(name string, item Item) {
	if player.inventory == nil {
		player.inventory = make(map[string]Item)
	}

	player.inventory[name] = item
}

func (player *Player) hasInventoryItem(itemName string) bool {
	_, isExist := player.inventory[itemName]
	return isExist
}
package main

type Item struct {
	name           string
	affectOn       []string
	onAffectedText map[string]string // text which is displayed when item is affected by another item
	isStorage      bool
	isWearable     bool
}

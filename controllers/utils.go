package controllers

import (
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

/*
Extracts the topic from an mqtt message and returns the last token of its topic
e.g: '/users/mike' -> 'mike'
*/
type UpdateRequest struct {
	OldName  string
	Username string
	Password string
}

func GetPath(message mqtt.Message) string {
	tokens := strings.Split(message.Topic(), "/")
	result := tokens[len(tokens)-1]
	return result
}

// Adds mqtt code to stringified json
func AddCodeStringJson(json string, code string) string {
	var newJson string
	length := len(json)
	index := 0

	runes := []rune(json)

	for index >= 0 && index < (length-1) {
		newJson = newJson + string(runes[index])
		index++
	}
	newJson = newJson + ",\"Code\": \"" + code + "\"}"
	return newJson
}

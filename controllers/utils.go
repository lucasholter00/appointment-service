package controllers

import (

)

func AddCodeStringJson (json string, code string) string {
    var newJson string
    length := len(json)
    index := 0

    runes := []rune(json)

    for index >= 0 && index < (length - 1) {
        newJson = newJson + string(runes[index])
        index++
    }
    newJson = newJson + "\",\"Code\": \"" + code + "\"}"
    return newJson
}

